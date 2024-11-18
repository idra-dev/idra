package tests

import (
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"testing"
	"time"
)

// var LocalConnectionString = os.Getenv("MONGODB")
var LocalConnectionString = "mongodb+srv://eyespot:DpZq55euNwc0aaJw@cluster1.rl9dgsm.mongodb.net/?retryWrites=true&w=majority&appName=Cluster1"
var token = ""

func TestMongoDBConnector(t *testing.T) {
	manager := data.MongodbConnector{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = LocalConnectionString
	connector.ConnectorType = "MongodbConnector"
	connector.Database = "sample"
	connector.Table = "movies"

	connector2 := cdc_shared.Connector{}
	connector2.ConnectionString = LocalConnectionString
	connector2.ConnectorType = "MongodbConnector"
	connector2.Database = "data"
	connector2.Table = "movies"
	manager.GetRowsByToken(connector, connector2)
	time.Sleep(1000 * time.Second)
}
