package tests

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"os"
	"testing"
)

var LocalMysqlConnectionString = os.Getenv("MYSQL_DB")


func TestMysqlMax(t *testing.T){
	manager := data.MysqlConnector{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = LocalMysqlConnectionString
	//connector.Table = "table2"
	connector.Table = "table"
	//connector.IdField = "id"
	connector.IdField = "id"
	id := manager.GetMaxTableId(connector)
	fmt.Println(id)
}

func TestMysqlMaxTimestamp(t *testing.T){
	manager := data.MysqlConnector{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = LocalMysqlConnectionString
	connector.Table = "table2"
	connector.TimestampField = "times"
	ts,_ := manager.GetMaxTimestamp(connector)
	fmt.Println(" "+ts.String())
}

func TestMysqlDBConnector(t *testing.T){
	manager := data.MysqlConnector{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = LocalMysqlConnectionString
	connector.IdField = "id"
	connector.Table = "table2"
	rows, res := manager.GetRowsById(connector, 0)
	fmt.Println(res, len(rows))
	manager2 := data.MysqlConnector{}
	connector2 := cdc_shared.Connector{}
	connector2.ConnectionString = LocalMysqlConnectionString
	connector2.IdField = "id"
	connector2.Table = "table2"
	manager2.InsertRows(connector2, rows)
}

