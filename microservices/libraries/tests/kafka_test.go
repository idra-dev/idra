package tests

import (
	"github.com/antrad1978/cdc_shared"
	data2 "microservices/libraries/data"
	"os"
	"testing"
)

func TestInsertRowsKafka(t *testing.T) {
	manager := data2.PostgresGormManager{}
	connector := cdc_shared.Connector{}
	connector.Table = "topic1"
	connector.IdField = "id"
	connector.ConnectionString = os.Getenv("KAFKA_URL")
	rows, _ := manager.GetRowsById(connector, 0)

	kafka := data2.KafkaConnector{}
	kafka.ClientName = "UnitTest"
	kafka.Acks = "all"
	connector2 := cdc_shared.Connector{}
	connector2.Table = "topic2"
	connector2.IdField = "id"
	kafka.InsertRows(connector2, rows)
}

func TestGetRowsKafka(t *testing.T) {
	kafka := data2.KafkaConnector{}
	kafka.ClientName = "UnitTest"
	kafka.ConsumerGroup = "cg1"
	connector2 := cdc_shared.Connector{}
	connector2.Table = "topic2"
	connector2.IdField = "id"
	connector2.ConnectionString = os.Getenv("KAFKA_URL")
	//kafka.GetRecords(connector2)
}
