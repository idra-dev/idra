package tests

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"testing"
)

func getAttributes() map[string]string {
	attributes := make(map[string]string)
	attributes["Key"] = "key"
	attributes["Value"] = "value"
	attributes["Column"] = "data"
	return attributes
}

func TestChromadbInsert(t *testing.T) {
	manager := data.ChromaDbConnector{}
	connector := cdc_shared.Connector{}
	attributes := getAttributes()
	connector.Attributes = attributes
	connector.Table = "products"
	connector.IdField = "id"
	connector.ConnectionString = "http://localhost:8000"

	items := []map[string]interface{}{}

	myMap := make(map[string]interface{})

	myMap["data"] = "GitHub allows to customize how issues and pull requests are presented to the public. Custom templates encourage collaboration and maintainability."

	items = append(items, myMap)

	res := manager.InsertRows(connector, items)
	fmt.Println(res)
}
