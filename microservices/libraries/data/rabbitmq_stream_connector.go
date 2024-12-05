package data

import (
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"microservices/libraries"
	"microservices/libraries/custom_errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type RabbitMQStreamConnector struct{}

func (RabbitMQStreamConnector) Name() string {
	return "RabbitMQStreamConnector"
}

func (RabbitMQStreamConnector) Modes() []string {
	return []string{"Last", "First", "Next"}
}

func (RabbitMQStreamConnector) GetRecords(sync cdc_shared.Sync) {
	//destinationProvider := RetrieveProvider(sync.DestinationConnector.ConnectorType)
	env, err := getEnv(sync.SourceConnector)
	if err != nil {
		custom_errors.CdcLog(sync.SourceConnector, err)
		return
	}
	var dataBatch []map[string]interface{}
	batchSize := 1
	value, exists := sync.SourceConnector.Attributes["batch"]
	if exists {
		batchSize, _ = strconv.Atoi(value)
	}

	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		processMessages(consumerContext, message, dataBatch, batchSize, sync)
	}

	offset := setOffsetStrategy(sync)

	consumer, err := env.NewConsumer(sync.SourceConnector.Table, messagesHandler,
		stream.NewConsumerOptions().SetOffset(offset))

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	run := true
	for run == true {
		select {
		case sig := <-sigChannel:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			err = consumer.Close()
			if err != nil {
				panic(err)
			}
			run = false
		}
	}

}

func setOffsetStrategy(sync cdc_shared.Sync) stream.OffsetSpecification {
	var offset stream.OffsetSpecification
	switch sync.Mode {
	case "Next":
		offset = stream.OffsetSpecification{}.Next()
	case "First":
		offset = stream.OffsetSpecification{}.First()
	case "Offset":
		offsetManager := libraries.IntOffset{}
		value := offsetManager.GetOffsetId(sync.Mode)
		res, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			res = 0
		}
		offset = stream.OffsetSpecification{}.Offset(res)
	default:
		offset = stream.OffsetSpecification{}.Last()
	}
	return offset
}

func processMessages(consumerContext stream.ConsumerContext, message *amqp.Message, dataBatch []map[string]interface{}, batchSize int, sync cdc_shared.Sync) {
	var dataValue map[string]interface{}

	if err := json.Unmarshal(message.Data[0], &dataValue); err != nil {
		custom_errors.CdcLog(sync.SourceConnector, err)
		return
	}
	dataBatch = append(dataBatch, dataValue)
	if len(dataBatch) >= batchSize {
		fmt.Printf("Stream: %s - Received message: %s\n", consumerContext.Consumer.GetStreamName(), message.Data)
		provider := RetrieveProvider(sync.DestinationConnector.ConnectorType)

		provider.InsertRows(sync.DestinationConnector, dataBatch)
		dataBatch = dataBatch[:0]
		if sync.Mode == "Offset" {
			offset := consumerContext.Consumer.GetOffset()
			libraries.IntOffset{}.SetOffsetId(sync.Id, offset)
		}
	}
}

func (reader RabbitMQStreamConnector) MoveData(sync cdc_shared.Sync) {
	fmt.Println("Started SYnc RabbitMQ streaming connector " + sync.SyncName)
	reader.GetRecords(sync)
}

func (RabbitMQStreamConnector) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	env, err := getEnv(connector)
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return -1
	}
	defer env.Close()
	producer, err := env.NewProducer(connector.Table, stream.NewProducerOptions())
	for _, row := range rows {
		byteArray, err := json.Marshal(row)
		if err != nil {
			fmt.Println("Error during serialization:", err)
			return -1
		}
		err = producer.Send(amqp.NewMessage(byteArray))
	}
	return 1
}

func getEnv(connector cdc_shared.Connector) (*stream.Environment, error) {
	port, err := strconv.Atoi(connector.Attributes["port"])
	if err != nil {
		port = 5552
	}
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost(connector.Attributes["host"]).
			SetPort(port).
			SetUser(connector.Attributes["user"]).
			SetPassword(connector.Attributes["password"]))
	return env, err
}
