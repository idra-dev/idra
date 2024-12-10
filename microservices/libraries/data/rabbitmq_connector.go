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

// GetRecords è stato aggiornato per accettare un context.
func (RabbitMQConnector) GetRecords(dataSync cdc_shared.Sync, ctx context.Context) {
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

	stopCh := make(chan struct{})

	wg.Add(1)
	go func(stopCh chan struct{}) {
		defer wg.Done() // Ensure WaitGroup is marked as done when the goroutine finishes
		for {
			select {
			case <-ctx.Done(): // Check if context is canceled
				fmt.Println(dataSync.SyncName + " is stopping due to context cancellation...")
				close(stopCh)
				return // Exit the goroutine
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

	// Wait for the goroutine to finish or be stopped by the context
	wg.Wait()
}

func (reader RabbitMQConnector) MoveData(sync cdc_shared.Sync, ctx context.Context) {
	fmt.Println("Start RabbitMQ sync " + sync.SyncName)
	reader.GetRecords(sync, ctx)
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
