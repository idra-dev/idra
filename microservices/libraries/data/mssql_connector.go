package data

import (
	"github.com/antrad1978/cdc_shared"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"microservices/libraries/custom_errors"
	"microservices/libraries/models"
	"strconv"
	"time"
)

type MssqlManager struct{
	ConnectorId string
	ClauseColumn string
}

const MssqlGORM = "MssqlGORM"

func (rdb MssqlManager) Name() string {
	return MssqlGORM
}

func (MssqlManager) Modes() []string {
	return []string{models.Id, models.Timestamp, models.LastDestinationId, models.LastDestinationTimestamp, models.FullWithId}
}

func (rdb MssqlManager) MoveData(sourceConnector cdc_shared.Connector, destinationConnector cdc_shared.Connector, mode string){

}

func GetMssqlDatabase(dsn string) (*gorm.DB, error) {
	// github.com/microsoft/go-mssqldb
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	return db, err
}

func (rdb MssqlManager) GetMaxTableId(connector cdc_shared.Connector) int64 {
	db,err := GetMssqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	query := "SELECT MAX(\""+connector.IdField+"\") FROM \"" + connector.Table+"\""
	return RetrieveMaxId(db, query)
}

func (rdb MssqlManager) GetMaxTimestamp(connector cdc_shared.Connector) (time.Time,error) {
	db,err := GetMssqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	query := "SELECT MAX(\""+connector.TimestampField+"\") FROM \"" + connector.Table+"\""
	return RetrieveMaxTimestamp(db, query)
}

func (rdb MssqlManager) GetRowsById(connector cdc_shared.Connector, lastId int64) ([]map[string]interface{}, int64){
	db,err := GetMssqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	var results []map[string]interface{}
	if connector.Query==""{
		db.Table(connector.Table).Where(" \"" + connector.IdField+"\">"+ strconv.FormatInt(lastId, 10), nil).Order("\"" + connector.IdField+"\""+" ASC").Limit(connector.MaxRecordBatchSize).Find(&results)
	}else{
		rows, err := db.Raw(connector.Query+" WHERE "+connector.IdField+" > "+ strconv.FormatInt(lastId, 10)).Rows()
		custom_errors.CdcLog(connector, err)
		defer rows.Close()
		results = GetQueryRows(rows, db, results)
	}
	offset := LastOffsetId(connector, results)
	if lastId >= offset{
		offset=lastId
	}
	return results, offset
}

func (rdb MssqlManager) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int{
	db,err := GetMssqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	i:= 0
	for ; i < len(rows); i++ {
		row := rows[i]
		tx := db.Table(connector.Table).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: rdb.ClauseColumn}},
			UpdateAll: true,
		}).Create(row)
		tx.Commit()
	}
	return i
}

func (rdb MssqlManager) GetRecordsByTimestamp(connector cdc_shared.Connector, lastTimestamp time.Time) ([]map[string]interface{}, time.Time){
	db,err := GetMssqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	var results []map[string]interface{}
	if connector.Query==""{
		db.Table(connector.Table).Where(" \"" + connector.TimestampField +"\">"+ lastTimestamp.String(), nil).Order("\"" + connector.TimestampField +"\""+" ASC").Limit(connector.MaxRecordBatchSize).Find(&results)
	}else{
		rows, err := db.Raw(connector.Query+" WHERE "+connector.TimestampField+" > "+ lastTimestamp.String() + " ORDER BY "+connector.TimestampField + " ASC").Rows()
		custom_errors.CdcLog(connector, err)
		defer rows.Close()
		results = GetQueryRows(rows, db, results)
	}
	var res time.Time
	if len(results) > 0{
		res = results[len(results)-1][connector.TimestampField].(time.Time)
	}else{
		res = lastTimestamp
	}
	return results, res
}



