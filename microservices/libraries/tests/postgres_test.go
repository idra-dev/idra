package tests

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"os"
	"testing"
)

var ConnectionString = os.Getenv("POSTGRES_DB_SOURCE")
var ConnectionStringDestination = os.Getenv("POSTGRES_DB_DESTINATION")

func TestMax(t *testing.T){
	manager := data.PostgresGormManager{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = ConnectionString
	connector.Table = "table"
	connector.IdField = "id"
	id := manager.GetMaxTableId(connector)
	if id != 1 {
		t.Errorf("got %q, wanted %q", id, 2)
	}
}

func TestMaxTimestamp(t *testing.T){
	manager := data.PostgresGormManager{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = ConnectionString
	connector.Table = "table"
	connector.TimestampField = "time"
	ts,_ := manager.GetMaxTimestamp(connector)
	fmt.Println(" "+ts.String())
}

func TestDBConnector(t *testing.T){
	manager := data.PostgresGormManager{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = ConnectionString
	connector.Table = "table"
	connector.IdField = "id"
	rows, _ := manager.GetRowsById(connector, 0)
	manager2 := data.PostgresGormManager{}
	connector2 := cdc_shared.Connector{}
	connector2.ConnectionString = ConnectionStringDestination
	connector2.Table = "table"
	manager2.InsertRows(connector2, rows)
}

