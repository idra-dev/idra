package data

import (
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"microservices/libraries/custom_errors"
	"microservices/libraries/models"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type KafkaConnector struct {
	ClientName    string
	ConsumerGroup string
	Offset        int64
	Acks          string
}

func (KafkaConnector) Name() string {
	return "KafkaConnector"
}

func (KafkaConnector) Modes() []string {
	return []string{models.Default}
}

func (reader KafkaConnector) MoveData(sourceConnector cdc_shared.Connector, destinationConnector cdc_shared.Connector, mode string) {
	destinationProvider := RetrieveProvider(destinationConnector.ConnectorType)
	reader.GetRecords(sourceConnector, destinationProvider, destinationConnector)
}

func (writer KafkaConnector) InsertRows(connector cdc_shared.Connector, records []map[string]interface{}) int {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": connector.ConnectionString,
		"client.id":         writer.ClientName,
		"acks":              writer.Acks})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	i := 0
	for _, record := range records {
		recordValue, err := json.Marshal(record)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			recordKey := fmt.Sprintf("%v", record[connector.IdField])
			fmt.Printf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &connector.Table, Partition: kafka.PartitionAny},
				Key:            []byte(recordKey),
				Value:          []byte(recordValue),
			}, nil)
			i++
		} else {
			break
		}
	}

	// Wait for all messages to be delivered
	p.Flush(5 * 1000)

	fmt.Printf("%d messages were produced to topic %s!", i, connector.Table)

	p.Close()
	return len(records)
}

func (reader KafkaConnector) GetRecords(connector cdc_shared.Connector, destinationProvider cdc_shared.ConnectorProvider, destinationConnector cdc_shared.Connector) {
	// Create Consumer instance
	c, err := reader.getConsumer(connector)
	custom_errors.CdcLog(connector, err)
	// Subscribe to topic
	err = c.SubscribeTopics([]string{connector.Table}, nil)
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	var tp kafka.TopicPartition
	tp.Topic = &connector.Table
	tp.Offset = kafka.Offset(connector.StartOffset)
	//partitionNumberParameter, ok := connector.Parameters["PartitionNumber"]
	// If the key exists
	//if ok {
	//	partitionNumber, _ := strconv.Atoi(partitionNumberParameter.ParameterValue)
	//	tp.Partition = int32(partitionNumber)
	//}else{
	tp.Partition = 0
	//}

	c.Seek(tp, 1000)

	// Process messages
	totalCount := 0
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:

			msg, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				fmt.Println(string(err.Error()))
				// Errors are informational and automatically handled by the consumer
				continue
			}
			recordKey := string(msg.Key)
			partition := msg.TopicPartition.Partition
			offset := int64(msg.TopicPartition.Offset)
			atomic.StoreInt64(&reader.Offset, offset)
			fmt.Printf("%d partition, %d offset\n", partition, offset)
			recordValue := msg.Value
			data := make(map[string]interface{})
			err = json.Unmarshal(recordValue, &data)
			if err != nil {
				fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
				continue
			}
			totalCount += 1
			fmt.Printf("Consumed record with key %s and value %s, and updated total count to %d\n", recordKey, recordValue, totalCount)
			var items []map[string]interface{}
			items = append(items, data)
			destinationProvider.InsertRows(destinationConnector, items)
			c.CommitMessage(msg)
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}

func (reader KafkaConnector) getConsumer(connector cdc_shared.Connector) (*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": connector.ConnectionString,
		"group.id":          reader.ConsumerGroup,
		"auto.offset.reset": "earliest",
	}
	authAttributes := []string{"sasl.mechanisms", "security.protocol", "sasl.username", "sasl.password"}
	setAttributesConfigMap(connector, config, authAttributes)
	c, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}
	return c, err
}

func setAttributesConfigMap(connector cdc_shared.Connector, config *kafka.ConfigMap, authAttributes []string) {
	for _, key := range authAttributes {
		value, exists := connector.Attributes[key]
		if exists {
			config.SetKey("sasl.mechanisms", value)
		}
	}
}

func (reader KafkaConnector) GetCGOffset(connector cdc_shared.Connector, topic string, partitionNumber int) int64 {
	var tp kafka.TopicPartition
	tp.Topic = &topic
	tp.Partition = int32(partitionNumber)
	c, err := reader.getConsumer(connector)
	res, err := c.Committed([]kafka.TopicPartition{tp}, partitionNumber)
	if err != nil {
		fmt.Printf("Failed to get offset: %s", err)
		return -1
	}
	return int64(res[0].Offset)
}
