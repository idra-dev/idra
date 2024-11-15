package tests

import (
	"github.com/antrad1978/cdc_shared"
	data2 "microservices/libraries/data"
	"testing"
)

func TestInsertRowsKafka(t *testing.T) {
	sync := cdc_shared.Sync{}
	kafka := data2.KafkaConnector{}
	//Producer
	connector := cdc_shared.Connector{}
	connector.ConnectorType = "KafkaConnector"
	connector.Table = "topic1"
	connector.IdField = "id"
	connector.ConnectionString = "127.0.0.1:9092"
	connector.Attributes = map[string]string{}
	connector.Attributes["client.id"] = "client1"
	connector.Attributes["acks"] = "all"
	//Consumer
	connector2 := cdc_shared.Connector{}
	connector2.Attributes = map[string]string{}
	connector2.ConnectorType = "KafkaConnector"
	connector2.Attributes["group.id"] = "me"
	connector2.Attributes["auto.offset.reset"] = "smallest"
	connector2.ConnectionString = "127.0.0.1:9092"
	connector2.Table = "topic2"
	connector2.IdField = "id"

	sync.SourceConnector = connector
	sync.DestinationConnector = connector2
	sync.Mode = "Default"

	kafka.MoveData(sync)
}
