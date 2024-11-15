package tests

import (
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"testing"
)

func TestInsertRowsRabbitMQStreaming(t *testing.T) {
	sync := cdc_shared.Sync{}
	sync.SyncName = "Test"
	kafka := data.RabbiMQStreamConnector{}
	//Producer
	connector := cdc_shared.Connector{}
	connector.ConnectorType = "RabbiMQStreamConnector"
	connector.Table = "hello-go-stream"
	connector.ConnectionString = "127.0.0.1"
	connector.Attributes = map[string]string{}
	connector.Attributes["username"] = "guest"
	connector.Attributes["password"] = "guest"
	connector.Attributes["port"] = "5552"
	//Consumer
	connector2 := cdc_shared.Connector{}
	connector2.ConnectorType = "RabbiMQStreamConnector"
	connector2.Attributes = map[string]string{}
	connector2.Attributes["username"] = "guest"
	connector2.Attributes["password"] = "guest"
	connector2.Attributes["port"] = "5552"
	connector2.ConnectionString = "127.0.0.1"
	connector2.Table = "sample_stream"

	sync.SourceConnector = connector
	sync.DestinationConnector = connector2
	sync.Mode = "Last"

	kafka.MoveData(sync)
}
