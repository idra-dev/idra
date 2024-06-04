package tests

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"testing"
)

func TestRest(t *testing.T) {
	manager := data.RestConnector{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = "https://jsonplaceholder.typicode.com/posts"
	connector.IdField = "id"

	res, lastId := manager.GetRowsById(connector, -1)
	fmt.Println(res)
	fmt.Println(lastId)
	manager.InsertRows(connector, res)
}
