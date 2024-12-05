package data

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"microservices/libraries/custom_errors"
	"microservices/libraries/models"
	"os"
	"strconv"
	"time"
)

type PostgresGormManager struct {
	ConnectorId string
}

func (rdb PostgresGormManager) Name() string {
	return "PostgresGORM"
}

func (PostgresGormManager) Modes() []string {
	return []string{models.Id, models.Timestamp, models.LastDestinationId, models.LastDestinationTimestamp, models.FullWithId}
}

func (rdb PostgresGormManager) MoveData(sync cdc_shared.Sync) {

}

func GetDatabase(dsn string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true, SkipDefaultTransaction: true, Logger: newLogger})
	return db, err
}

func (rdb PostgresGormManager) GetMaxTableId(connector cdc_shared.Connector) int64 {
	db, err := GetDatabase(connector.ConnectionString)
	sqlDB := getDB(db)
	defer sqlDB.Close()
	custom_errors.CdcLog(connector, err)
	query := "SELECT MAX(\"" + connector.IdField + "\") FROM \"" + connector.Table + "\""
	offset := RetrieveMaxId(db, query)
	return offset
}

func (rdb PostgresGormManager) GetMaxTimestamp(connector cdc_shared.Connector) (time.Time, error) {
	db, err := GetDatabase(connector.ConnectionString)
	sqlDB := getDB(db)
	defer sqlDB.Close()
	custom_errors.CdcLog(connector, err)
	query := "SELECT MAX(\"" + connector.TimestampField + "\") FROM " + connector.Table
	return RetrieveMaxTimestamp(db, query)
}

func (rdb PostgresGormManager) GetRowsById(connector cdc_shared.Connector, lastId int64) ([]map[string]interface{}, int64) {
	db, err := GetDatabase(connector.ConnectionString)
	sqlDB := getDB(db)
	defer sqlDB.Close()
	custom_errors.CdcLog(connector, err)
	var results []map[string]interface{}
	if connector.Query == "" {
		tx := db.Debug().Table(connector.Table).Where(connector.IdField+">?", strconv.FormatInt(lastId, 10)).Order(connector.IdField + " ASC").Limit(models.MaxBatchSizeDefault).Find(&results)
		if tx.Error != nil {
			custom_errors.CdcLog(connector, err)
		}
	} else {
		rows, err := db.Raw(connector.Query + " WHERE " + connector.IdField + " > " + strconv.FormatInt(lastId, 10)).Rows()
		custom_errors.CdcLog(connector, err)
		defer rows.Close()
		results = GetQueryRows(rows, db, results)
	}

	offset := LastOffsetId(connector, results)
	if lastId >= offset {
		offset = lastId
	}
	return results, offset
}

func (rdb PostgresGormManager) GetRecordsByTimestamp(connector cdc_shared.Connector, lastTimestamp time.Time) ([]map[string]interface{}, time.Time) {
	db, err := GetDatabase(connector.ConnectionString)
	sqlDB := getDB(db)
	defer sqlDB.Close()
	var res time.Time
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return nil, time.Time{}
	}
	var results []map[string]interface{}
	if connector.Query == "" {
		lastTimestampValue := ""
		if connector.TimestampFieldFormat == "" {
			lastTimestampValue = lastTimestamp.Format("2006-01-02 15:04:05.99999")
		} else {
			lastTimestampValue = lastTimestamp.Format(connector.TimestampFieldFormat)
		}
		tx := db.Debug().Table(connector.Table).Where(" \"" + connector.TimestampField + "\">'" + lastTimestampValue + "'").Order("\"" + connector.TimestampField + "\"" + " ASC").Limit(connector.MaxRecordBatchSize).Find(&results)
		fmt.Println(tx.Error)
	} else {
		lastTimestampFormatted := fmt.Sprintf("%v", lastTimestamp.Format("2006-01-02 15:04:05.99999"))
		query := connector.Query + " WHERE `" + connector.TimestampField + "` > '" + lastTimestampFormatted + "' ORDER BY " + connector.TimestampField + AscLimit + strconv.FormatInt(models.MaxBatchSizeDefault, 10)
		rows, err := db.Raw(query).Rows()
		custom_errors.CdcLog(connector, err)
		defer rows.Close()
		results = GetQueryRows(rows, db, results)
	}
	if len(results) > 0 {
		res = results[len(results)-1][connector.TimestampField].(time.Time)
	} else {
		res = lastTimestamp
	}
	return results, res
}

func (rdb PostgresGormManager) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	db, err := GetDatabase(connector.ConnectionString)
	sqlDB := getDB(db)
	defer sqlDB.Close()
	custom_errors.CdcLog(connector, err)
	return SaveData(connector, rows, db)
}
