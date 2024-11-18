package data

import (
	"context"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"github.com/codenotary/immudb/pkg/api/schema"
	immudb "github.com/codenotary/immudb/pkg/client"
	"log"
	"microservices/libraries/models"
	"strconv"
	"time"
)

type ImmudbDriver struct {
	ip       string
	username string
	password string
	port     int
}

func (rdb ImmudbDriver) MoveData(sync cdc_shared.Sync) {
}

func (rdb ImmudbDriver) Name() string {
	return "ImmudbDriver"
}

func (ImmudbDriver) Modes() []string {
	return []string{models.Id, models.Timestamp}
}

func (rdb ImmudbDriver) GetMaxTableId(connector cdc_shared.Connector) int64 {
	rdb.ip = connector.Attributes["Host"]
	client := immudb.NewClient().WithOptions(rdb.InitOptions())
	err := client.OpenSession(context.Background(), []byte(connector.Attributes["username"]), []byte(connector.Attributes["password"]), connector.Attributes["database"])
	if err != nil {
		log.Fatal(err)
	}
	defer client.CloseSession(context.Background())
	query := "SELECT MAX(\"" + connector.IdField + "\") FROM \"" + connector.Table + "\""
	res, err := client.SQLQuery(context.Background(), query, nil, true)
	if err != nil {
		log.Fatal(err)
	}
	return res.GetRows()[0].Values[0].GetN()
}

func (rdb ImmudbDriver) InitOptions() *immudb.Options {
	opts := immudb.DefaultOptions().WithAddress(rdb.ip).WithPort(rdb.port)
	return opts
}

func (rdb ImmudbDriver) GetMaxTimestamp(connector cdc_shared.Connector) (int64, error) {
	client := immudb.NewClient().WithOptions(rdb.InitOptions())
	err := client.OpenSession(context.Background(), []byte(connector.Attributes["usernama"]), []byte(connector.Attributes["password"]), connector.Attributes["database"])
	if err != nil {
		log.Fatal(err)
	}
	log.Println(connector, err)
	query := "SELECT MAX(\"" + connector.TimestampField + "\") FROM \"" + connector.Table + "\""
	res, err := client.SQLQuery(context.Background(), query, nil, true)
	if err != nil {
		log.Fatal(err)
	}
	return res.GetRows()[0].Values[0].GetTs(), nil
}

func (rdb ImmudbDriver) GetRowsById(connector cdc_shared.Connector, lastId int64) ([]map[string]interface{}, int64) {
	client := immudb.NewClient().WithOptions(rdb.InitOptions())
	err := client.OpenSession(context.Background(), []byte(connector.Attributes["usernama"]), []byte(connector.Attributes["password"]), connector.Attributes["database"])
	if err != nil {
		log.Fatal(err)
	}
	defer client.CloseSession(context.Background())
	var result *schema.SQLQueryResult
	if connector.Query == "" {
		result, err = client.SQLQuery(context.Background(), "SELECT * FROM "+connector.Table+" WHERE "+connector.IdField+" > "+strconv.FormatInt(lastId, 10)+" ORDER BY "+connector.IdField+" ASC", nil, true)
	} else {
		result, err = client.SQLQuery(context.Background(), connector.Query+" WHERE "+connector.IdField+" > "+strconv.FormatInt(lastId, 10)+" ORDER BY "+connector.IdField+" ASC", nil, true)
	}
	var items []map[string]interface{}
	rows := result.GetRows()
	for _, rowValue := range rows {
		row := make(map[string]interface{})
		for idx, colName := range result.GetColumns() {
			row[colName.Name] = rowValue.Values[idx]
		}
		items = append(items, row)
	}
	max := int64(0)
	if len(items) > 0 {
		lastRow := rows[len(rows)-1]
		lastRow.Values[rdb.GetColumnIndex(lastRow.GetColumns(), "("+connector.Table+"."+connector.IdField+")")].GetN()
	}
	return items, max
}

func (rdb ImmudbDriver) GetRecordsByTimestamp(connector cdc_shared.Connector, lastTimestamp time.Time) ([]map[string]interface{}, int64) {
	client := immudb.NewClient().WithOptions(rdb.InitOptions())
	err := client.OpenSession(context.Background(), []byte(connector.Attributes["usernama"]), []byte(connector.Attributes["password"]), connector.Attributes["database"])
	if err != nil {
		log.Fatal(err)
	}
	result, err := client.SQLQuery(context.Background(), connector.Query+" WHERE "+connector.TimestampField+" > "+strconv.FormatInt(lastTimestamp.UnixMilli(), 10), nil, true)
	if err != nil {
		log.Fatal(err)
	}
	var items []map[string]interface{}
	rows := result.GetRows()
	for _, rowValue := range rows {
		row := make(map[string]interface{})
		for idx, colName := range result.GetColumns() {
			row[colName.Name] = rowValue.Values[idx]
		}
		items = append(items, row)
	}
	client.CloseSession(context.Background())
	index := rdb.GetColumnIndex(rows[len(rows)-1].Columns, connector.TimestampField)
	return items, rows[len(rows)-1].Values[index].GetTs()
}

func (rdb ImmudbDriver) GetColumnIndex(columns []string, desiredColumnName string) int {
	var desiredColumnIndex int = -1
	for i, colName := range columns {
		if colName == desiredColumnName {
			desiredColumnIndex = i
			break
		}
	}

	if desiredColumnIndex == -1 {
		log.Fatalf("Columns '%s' not found", desiredColumnName)
	}
	return desiredColumnIndex
}

func (rdb ImmudbDriver) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	client := immudb.NewClient().WithOptions(rdb.InitOptions())
	err := client.OpenSession(context.Background(), []byte(connector.Attributes["usernama"]), []byte(connector.Attributes["password"]), connector.Attributes["database"])
	if err != nil {
		log.Fatal(err)
	}
	defer client.CloseSession(context.Background())
	for _, item := range rows {
		columns := ""
		values := ""

		for key, _ := range item {
			columns += fmt.Sprintf("%s,", key)
			values += fmt.Sprintf("@%s,", key)
		}
		columns = columns[:len(columns)-1]
		values = values[:len(values)-1]

		query := "UPSERT INTO " + connector.Table + " (" + columns + ") VALUES (" + values + ");"

		_, err = client.SQLExec(context.Background(), query, item)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successfully row insertion\n")
	}
	return len(rows)
}
