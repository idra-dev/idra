package data

import (
	"database/sql"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"microservices/libraries/custom_errors"
	"microservices/libraries/models"
	"strconv"
	"strings"
	"time"
)

type MssqlManager struct {
	ConnectorId  string
	ClauseColumn string
}

const MssqlGORM = "MssqlGORM"

const MergeInstruction = " " +
	" MERGE %s AS T USING (VALUES ( %s )) AS TMP (%s) " +
	" ON ( %s ) " +
	" WHEN MATCHED THEN " +
	" UPDATE SET %s " +
	" WHEN NOT MATCHED BY TARGET THEN " +
	" INSERT (%s) VALUES(%s);"

func (rdb MssqlManager) Name() string {
	return MssqlGORM
}

func (MssqlManager) Modes() []string {
	return []string{models.Id, models.Timestamp, models.LastDestinationId, models.LastDestinationTimestamp, models.FullWithId}
}

func (rdb MssqlManager) MoveData(sync cdc_shared.Sync) {

}

func GetMssqlDatabase(dsn string) (*gorm.DB, error) {
	// github.com/microsoft/go-mssqldb
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{PrepareStmt: true})
	return db, err
}

func (rdb MssqlManager) GetMaxTableId(connector cdc_shared.Connector) int64 {
	db, sqlDB := getSQLServerDB(connector)
	defer sqlDB.Close()
	query := "SELECT MAX(\"" + connector.IdField + "\") FROM \"" + connector.Table + "\""
	return RetrieveMaxId(db, query)
}

func (rdb MssqlManager) GetMaxTimestamp(connector cdc_shared.Connector) (time.Time, error) {
	db, sqlDB := getSQLServerDB(connector)
	defer sqlDB.Close()
	query := "SELECT MAX(" + connector.TimestampField + ") FROM " + connector.Table
	return RetrieveMaxTimestamp(db, query)
}

func getSQLServerDB(connector cdc_shared.Connector) (*gorm.DB, *sql.DB) {
	db, err := GetMssqlDatabase(connector.ConnectionString)
	sqlDB := getDB(db)
	custom_errors.CdcLog(connector, err)
	return db, sqlDB
}

func (rdb MssqlManager) GetRowsById(connector cdc_shared.Connector, lastId int64) ([]map[string]interface{}, int64) {
	db, err := GetMssqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	sqlDB := getDB(db)
	defer sqlDB.Close()
	var results []map[string]interface{}
	if connector.Query == "" {
		db.Table(connector.Table).Where(" \""+connector.IdField+"\">"+strconv.FormatInt(lastId, 10), nil).Order("\"" + connector.IdField + "\"" + " ASC").Limit(connector.MaxRecordBatchSize).Find(&results)
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

func (rdb MssqlManager) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	db, err := GetMssqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	sqlDB := getDB(db)
	defer sqlDB.Close()
	i := 0
	for ; i < len(rows); i++ {
		row := rows[i]
		if rdb.ClauseColumn == "" {
			rdb.ClauseColumn = "id"
		}
		tx := db.Begin()
		if err := tx.Exec("SET IDENTITY_INSERT " + connector.Table + " ON;").Error; err != nil {
			return -1
		}
		tx.PrepareStmt = true
		query := getMergeQuery(connector, row)
		if err := tx.Exec(query).Error; err != nil {
			custom_errors.CdcLog(connector, err)
			return -1
		}
		if err := tx.Exec("SET IDENTITY_INSERT " + connector.Table + " OFF;").Error; err != nil {
			return -1
		}

		tx.Commit()
		if tx.Error != nil {
			return -1
		}
	}
	return i
}

func (rdb MssqlManager) GetRecordsByTimestamp(connector cdc_shared.Connector, lastTimestamp time.Time) ([]map[string]interface{}, time.Time) {
	db, err := GetMssqlDatabase(connector.ConnectionString)
	custom_errors.CdcLog(connector, err)
	var results []map[string]interface{}
	if connector.Query == "" {
		if lastTimestamp.IsZero() {
			db.Table(connector.Table).Order("\"" + connector.TimestampField + "\"" + " ASC").Limit(connector.MaxRecordBatchSize).Find(&results)
		} else {
			lastTimestampValue := ""
			if connector.TimestampFieldFormat == "" {
				lastTimestampValue = lastTimestamp.Format("2024-07-17 13:22:40.983")
			} else {
				lastTimestampValue = lastTimestamp.Format(connector.TimestampFieldFormat)
			}
			db.Debug().Table(connector.Table).Where(" \""+connector.TimestampField+"\">'"+lastTimestampValue+"'", nil).Order("\"" + connector.TimestampField + "\"" + " ASC").Limit(connector.MaxRecordBatchSize).Find(&results)
		}

	} else {
		rows, err := db.Raw(connector.Query + " WHERE " + connector.TimestampField + " > " + lastTimestamp.String() + " ORDER BY " + connector.TimestampField + " ASC").Rows()
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

func getMergeQuery(connector cdc_shared.Connector, values map[string]interface{}) string {
	keys := getKeys(values)
	identifiers := strings.Split(connector.IdField, ",")
	updateParameters := composeUpdateParameters(keys, identifiers, connector)
	queryValues := composeValues(keys, values, connector)
	insertNames := composeNamesInsert(keys)
	return fmt.Sprintf(MergeInstruction, connector.Table, queryValues, strings.Join(keys, ", "), composeKey(connector.IdField), updateParameters, strings.Join(keys, ", "), insertNames)
}

func composeKey(identifier string) string {
	values := strings.Split(identifier, ",")
	var builder strings.Builder
	for i, current := range values {
		if i > 0 {
			builder.WriteString(" AND ")
		}
		builder.WriteString("T." + current + " = TMP." + current)
	}
	return builder.String()
}

func composeUpdateParameters(parameters []string, identifiers []string, connector cdc_shared.Connector) string {
	result := ""
	length := len(parameters)
	for i := 0; i < length; i++ {
		current := parameters[i]
		if !stringInSlice(current, strings.Split(connector.IdField, ",")) {
			if length == 1 || i == len(parameters)-1 {
				result += current + "=" + "TMP." + current
			} else {
				result += current + "=" + "TMP." + current + ", "
			}
		}
	}
	return result
}

func composeValues(params []string, values map[string]interface{}, connector cdc_shared.Connector) string {
	var builder strings.Builder
	for i, param := range params {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(convertToString(values[param], connector))
	}
	return builder.String()
}

func composeNamesInsert(params []string) string {
	names := make([]string, len(params))
	for i, param := range params {
		names[i] = "TMP." + param
	}
	return strings.Join(names, ", ")
}

func convertToString(value interface{}, connector cdc_shared.Connector) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case string:
		return "'" + strings.Replace(v, "'", "''", -1) + "'"
	case time.Time:
		return "'" + value.(time.Time).Format(connector.TimestampFieldFormat) + "'"
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.Itoa(int(v))
	case int16:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case complex64:
		return fmt.Sprintf("%g", v)
	case complex128:
		return fmt.Sprintf("%g", v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

func getKeys(m map[string]interface{}) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func stringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
