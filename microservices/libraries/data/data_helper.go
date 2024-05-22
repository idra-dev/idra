package data

import (
	"database/sql"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"gorm.io/gorm"
	"microservices/libraries/custom_errors"
	"time"
)

func LastOffsetId(connector cdc_shared.Connector, results []map[string]interface{}) int64 {
	var offset int64 = -1
	if len(results) > 0 {
		value := results[len(results)-1][connector.IdField]

		switch value.(type) {
		case int16:
			offset = int64(value.(int16))
		case int32:
			offset = int64(value.(int32))
		default:
			offset = value.(int64)
		}
	}
	return offset
}

func GetQueryRows(rows *sql.Rows, db *gorm.DB, results []map[string]interface{}) []map[string]interface{} {
	for rows.Next() {
		var result map[string]interface{}
		db.ScanRows(rows, &result)
		results = append(results, result)
	}
	return results
}

func RetrieveMaxTimestamp(db *gorm.DB, query string) (time.Time, error) {
	var timestamp sql.NullTime
	db.Raw(query).Scan(&timestamp)
	dbConnection, err := db.DB()
	defer dbConnection.Close()
	return timestamp.Time, err
}

func RetrieveMaxId(db *gorm.DB, query string) int64 {
	var offset sql.NullInt64
	tx := db.Raw(query).Scan(&offset)
	if tx.Error != nil {
		handler := custom_errors.CdcErrorHandler{}
		handler.SaveExecutionError(fmt.Sprintf("%v", tx.Error.Error()))
		return -1
	}
	return offset.Int64
}


