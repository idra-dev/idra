package main

import (
	"github.com/antrad1978/cdc_shared"
	data2 "microservices/libraries/data"
	"os"
	"testing"
)

func TestInsertRowsKafka(t *testing.T){
	manager := data2.PostgresGormManager{}
	connector := cdc_shared.Connector{}
	connector.Table = "table"
	connector.IdField = "id"
	connector.ConnectionString = os.Getenv("KAFKA_URL")
	rows, _ := manager.GetRowsById(connector, 0)

	kafka := data2.KafkaConnector{}
	kafka.Brokers = "localhost:9092"
	kafka.ClientName = "UnitTest"
	kafka.Acks = "all"
	connector2 := cdc_shared.Connector{}
	connector2.Table = "table"
	connector2.IdField = "id"
	kafka.InsertRows(connector2, rows)
}
