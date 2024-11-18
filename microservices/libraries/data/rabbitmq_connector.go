package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"microservices/libraries/custom_errors"
	"sync"
	"time"
)

type RabbitMQConnector struct{}

func (RabbitMQConnector) Name() string {
	return "RabbitMQConnector"
}

func (RabbitMQConnector) Modes() []string {
	return []string{"AutoAck", "ManualAck"}
}

func (RabbitMQConnector) GetRecords(dataSync cdc_shared.Sync) {
	conn, err := amqp.Dial(dataSync.SourceConnector.ConnectionString)
	failOnError(err, dataSync.SourceConnector)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, dataSync.SourceConnector)
	defer ch.Close()

	autoAck := true
	if dataSync.Mode == "ManualAck" {
		autoAck = false
	}

	messages, err := ch.Consume(
		dataSync.SourceConnector.Table,         // queue
		dataSync.SourceConnector.ConnectorName, // consumer
		autoAck,                                // auto-ack
		false,                                  // exclusive
		false,                                  // no-local
		false,                                  // no-wait
		nil,                                    // args
	)
	failOnError(err, dataSync.SourceConnector)

	provider := RetrieveProvider(dataSync.DestinationConnector.ConnectorType)
	var result map[string]interface{}

	var wg sync.WaitGroup

	wg.Add(1)

	stopCh := make(chan struct{})
	go func(stopCh chan struct{}) {
		for {
			select {
			case <-stopCh:
				fmt.Println(dataSync.SyncName + " is stopping...")
				return // exit the goroutine
			default:
				for message := range messages {
					err := json.Unmarshal(message.Body, &result)
					if err != nil {
						custom_errors.CdcLog(dataSync.DestinationConnector, err)
					}
					provider.InsertRows(dataSync.DestinationConnector, []map[string]interface{}{result})
					if !autoAck {
						ch.Ack(message.DeliveryTag, true)
					}
				}
			}
		}
	}(stopCh)
	wg.Wait()
}

func (reader RabbitMQConnector) MoveData(sync cdc_shared.Sync) {
	reader.GetRecords(sync)
}

func (RabbitMQConnector) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	conn, err := amqp.Dial(connector.ConnectionString)
	failOnError(err, connector)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, connector)
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for _, row := range rows {
		body, err := json.Marshal(row)
		if err != nil {
			log.Fatalf("Errore nella serializzazione: %v", err)
		}
		err = ch.PublishWithContext(ctx,
			connector.Table,             // exchange
			connector.Attributes["key"], // routing key
			false,                       // mandatory
			false,                       // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(body),
			})
		failOnError(err, connector)
		log.Printf(" [x] Sent %s\n", body)
	}

	return len(rows)
}

func failOnError(err error, connector cdc_shared.Connector) {
	if err != nil {
		custom_errors.CdcLog(connector, err)
	}
}
