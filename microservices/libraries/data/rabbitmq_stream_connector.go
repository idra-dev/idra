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

type RabbiMQStreamConnector struct{}

func (RabbiMQStreamConnector) Name() string {
	return "RabbiMQStreamConnector"
}

func (RabbiMQStreamConnector) Modes() []string {
	return []string{"Last", "First", "Next"}
}

func (RabbiMQStreamConnector) GetRecords(connector cdc_shared.Connector, destinationProvider cdc_shared.ConnectorProvider, destinationConnector cdc_shared.Connector, mode string) {
	env, err := getEnv(connector)
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return
	}
	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		fmt.Printf("Stream: %s - Received message: %s\n", consumerContext.Consumer.GetStreamName(), message.Data)
	}

	var offset stream.OffsetSpecification

	switch mode {
	case "Next":
		offset = stream.OffsetSpecification{}.Next()
	case "First":
		offset = stream.OffsetSpecification{}.First()
	case "Offset":
		offsetManager := libraries.IntOffset{}
		value := offsetManager.GetOffsetId("")
		res, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			res = 0
		}
		offset = stream.OffsetSpecification{}.Offset(res)
	default:
		offset = stream.OffsetSpecification{}.Last()
	}

	consumer, err := env.NewConsumer(connector.Table, messagesHandler,
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
		default:
			fmt.Print("Process " + connector.ConnectorName)
		}
	}

}

func (reader RabbiMQStreamConnector) MoveData(sourceConnector cdc_shared.Connector, destinationConnector cdc_shared.Connector, mode string) {
	destinationProvider := RetrieveProvider(destinationConnector.ConnectorType)
	reader.GetRecords(sourceConnector, destinationProvider, destinationConnector, mode)
}

func (RabbiMQStreamConnector) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	env, err := getEnv(connector)
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return -1
	}
	producer, err := env.NewProducer(connector.Table, stream.NewProducerOptions())
	for _, row := range rows {
		byteArray, err := json.Marshal(row)
		if err != nil {
			fmt.Println("Error during serialization:", err)
			return -1
		}
		err = producer.Send(amqp.NewMessage(byteArray))
	}

	err = producer.Close()
	if err != nil {
		panic(err)
	}
	return 1
}

func getEnv(connector cdc_shared.Connector) (*stream.Environment, error) {
	port, err := strconv.Atoi(connector.Attributes["port"])
	if err != nil {
		port = 5672
	}
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost(connector.Attributes["host"]).
			SetPort(port).
			SetUser(connector.Attributes["user"]).
			SetPassword(connector.Attributes["password"]))
	return env, err
}
