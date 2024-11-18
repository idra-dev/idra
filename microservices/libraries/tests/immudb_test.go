package tests

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"testing"
)

func TestImmudbMax(t *testing.T) {
	manager := data.ImmudbDriver{}
	connector := cdc_shared.Connector{}
	attributes := getAttributesImmudb()
	connector.Attributes = attributes
	connector.Table = "products"
	connector.IdField = "id"
	id := manager.GetMaxTableId(connector)
	fmt.Println(id)
}

func getAttributesImmudb() map[string]string {
	attributes := make(map[string]string)
	attributes["usernama"] = "immudb"
	attributes["password"] = "immudb"
	attributes["port"] = "5432"
	attributes["database"] = "defaultdb"
	return attributes
}

func TestImmudblInsert(t *testing.T) {
	manager := data.ImmudbDriver{}
	connector := cdc_shared.Connector{}
	attributes := getAttributesImmudb()
	connector.Attributes = attributes
	connector.Table = "products"
	connector.IdField = "id"

	items := []map[string]interface{}{}

	myMap := make(map[string]interface{})

	// Aggiunta di valori alla mappa
	myMap["id"] = 102
	myMap["price"] = "25"
	myMap["product"] = "caca"

	items = append(items, myMap)

	res := manager.InsertRows(connector, items)
	fmt.Println(res)
}

func TestImmudbDBConnector(t *testing.T) {
	manager := data.ImmudbDriver{}
	connector := cdc_shared.Connector{}
	attributes := getAttributesImmudb()
	connector.Attributes = attributes
	connector.Table = "products"
	connector.IdField = "id"
	rows, res := manager.GetRowsById(connector, 0)
	fmt.Println(res, len(rows))
}
