package data

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"microservices/libraries/custom_errors"
	"microservices/libraries/models"
	"os"
	"strconv"
	"time"
)

const AscLimit = " ASC LIMIT "

type MysqlConnector struct {
	ConnectorId string
}

func (rdb MysqlConnector) Name() string {
	return "MysqlGORM"
}

func (MysqlConnector) Modes() []string {
	return []string{models.Id, models.Timestamp, models.LastDestinationId, models.LastDestinationTimestamp, models.FullWithId}
}

func (rdb MysqlConnector) MoveData(sync cdc_shared.Sync) {

}

func GetMysqlDatabase(dsn string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})
	db.Exec("SET sql_mode='IGNORE_SPACE,STRICT_TRANS_TABLES,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';")
	db.Logger = newLogger
	return db, err
}

func (rdb MysqlConnector) GetMaxTableId(connector cdc_shared.Connector) int64 {
	db, err := GetMysqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	query := "SELECT MAX(`" + connector.IdField + "`) FROM `" + connector.Table + "`"
	sqlDB := getDB(db)
	defer sqlDB.Close()
	return RetrieveMaxId(db, query)
}

func (rdb MysqlConnector) GetMaxTimestamp(connector cdc_shared.Connector) (time.Time, error) {
	db, err := GetMysqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	query := "SELECT MAX(`" + connector.TimestampField + "`) FROM `" + connector.Table + "`"
	sqlDB := getDB(db)
	defer sqlDB.Close()
	return RetrieveMaxTimestamp(db, query)
}

func (rdb MysqlConnector) GetRowsById(connector cdc_shared.Connector, lastId int64) ([]map[string]interface{}, int64) {
	db, err := GetMysqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	var results []map[string]interface{}
	if connector.Query == "" {
		if connector.MaxRecordBatchSize == 0 {
			db.Table(connector.Table).Where(" `" + connector.IdField + "` > " + strconv.FormatInt(lastId, 10)).Order("`" + connector.IdField + "`" + " ASC").Limit(models.MaxBatchSizeDefault).Find(&results)
		} else {
			db.Table(connector.Table).Where(" `" + connector.IdField + "` > " + strconv.FormatInt(lastId, 10)).Order("`" + connector.IdField + "`" + " ASC").Limit(connector.MaxRecordBatchSize).Find(&results)
		}
	} else {
		var query string
		if connector.MaxRecordBatchSize == 0 {
			query = connector.Query + " WHERE " + connector.IdField + " > " + strconv.FormatInt(lastId, 10) + " ORDER BY " + connector.IdField + AscLimit + strconv.FormatInt(models.MaxBatchSizeDefault, 10)
		} else {
			query = connector.Query + " WHERE " + connector.IdField + " > " + strconv.FormatInt(lastId, 10) + " ORDER BY " + connector.IdField + AscLimit + strconv.Itoa(connector.MaxRecordBatchSize)
		}
		rows, err := db.Raw(query).Rows()
		custom_errors.CdcLog(connector, err)
		defer rows.Close()
		results = GetQueryRows(rows, db, results)
	}
	offset := LastOffsetId(connector, results)
	if lastId >= offset {
		offset = lastId
	}
	sqlDB := getDB(db)
	defer sqlDB.Close()
	return results, offset
}

func (rdb MysqlConnector) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	db, err := GetMysqlDatabase(connector.ConnectionString)
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return -1
	}
	sqlDB := getDB(db)
	defer sqlDB.Close()
	return SaveData(connector, rows, db)
}

func (rdb MysqlConnector) GetRecordsByTimestamp(connector cdc_shared.Connector, lastTimestamp time.Time) ([]map[string]interface{}, time.Time) {
	db, err := GetMysqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	sqlDB := getDB(db)
	defer sqlDB.Close()
	var results []map[string]interface{}
	if connector.Query == "" {
		db.Table(connector.Table).Where(" \""+connector.TimestampField+"\">"+lastTimestamp.String(), nil).Order("\"" + connector.TimestampField + "\"" + " ASC").Limit(connector.MaxRecordBatchSize).Find(&results)
	} else {
		lastTimestampFormatted := fmt.Sprintf("%v", lastTimestamp.Format("2006-01-02 15:04:05.999"))
		query := connector.Query + " WHERE `" + connector.TimestampField + "` > '" + lastTimestampFormatted + "' ORDER BY " + connector.TimestampField + AscLimit + strconv.FormatInt(models.MaxBatchSizeDefault, 10)
		rows, err := db.Raw(query).Rows()
		custom_errors.CdcLog(connector, err)
		defer rows.Close()
		results = GetQueryRows(rows, db, results)
	}
	var res time.Time
	if len(results) > 0 {
		res = results[len(results)-1][connector.TimestampField].(time.Time)
	} else {
		res = lastTimestamp
	}
	return results, res
}
