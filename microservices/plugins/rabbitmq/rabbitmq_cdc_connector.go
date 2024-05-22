package main

import (
	"github.com/antrad1978/cdc_shared"
)

type RabbitmqConnector struct {
	ConnectorId string
	Brokers string
	ClientName string
	ConsumerGroup string
	Offset int64
	ConnectionString string
	Acks string
}


func (RabbitmqConnector) Name() string {
	return "RabbitMQConnector"
}

func (RabbitmqConnector) Modes() []string {
	return []string{"Default"}
}

func (reader RabbitmqConnector) MoveData(sourceConnector cdc_shared.Connector, destinationConnector cdc_shared.Connector, mode string){

}

func (writer RabbitmqConnector) InsertRows(connector cdc_shared.Connector, records []map[string]interface{}) int{
	return 1
}

var Connector RabbitmqConnector






