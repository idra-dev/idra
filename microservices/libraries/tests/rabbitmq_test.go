package tests

import (
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"testing"
	"time"
)

func TestInsertRowsRabbitMQ(t *testing.T) {
	sync := cdc_shared.Sync{}
	sync.SyncName = "Test"
	rabbit := data.RabbitMQConnector{}
	//Producer
	connector := cdc_shared.Connector{}
	connector.ConnectorType = "RabbitMQConnector"
	connector.Table = "hello-go"
	connector.ConnectionString = "amqp://guest:guest@localhost:5672/"
	connector.Attributes = map[string]string{}
	connector.Attributes["username"] = "guest"
	//Consumer
	connector2 := cdc_shared.Connector{}
	connector2.ConnectorType = "RabbitMQConnector"
	connector2.Attributes = map[string]string{}
	connector2.Attributes["username"] = "guest"
	connector2.ConnectionString = "amqp://guest:guest@localhost:5672/"
	connector2.Table = "sample"

	sync.SourceConnector = connector
	sync.DestinationConnector = connector2
	sync.Mode = "Last"

	rabbit.MoveData(sync)

	time.Sleep(30 * time.Second)
}
