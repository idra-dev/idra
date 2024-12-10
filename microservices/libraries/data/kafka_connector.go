package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"microservices/libraries/custom_errors"
	"microservices/libraries/models"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

type KafkaConnector struct {
	Offset int64
}

func (KafkaConnector) Name() string {
	return "KafkaConnector"
}

func (KafkaConnector) Modes() []string {
	return []string{models.Default}
}

func (connector KafkaConnector) MoveData(sync cdc_shared.Sync, ctx context.Context) {
	destinationProvider := RetrieveProvider(sync.DestinationConnector.ConnectorType)
	connector.GetRecords(sync.SourceConnector, destinationProvider, sync.DestinationConnector)
}

func (writer KafkaConnector) InsertRows(connector cdc_shared.Connector, records []map[string]interface{}) int {
	config := getConfigMap(connector)
	p, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	i := 0
	for _, record := range records {
		recordValue, err := json.Marshal(record)
		if err == nil {
			recordKey := fmt.Sprintf("%v", record[connector.IdField])
			fmt.Printf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &connector.Table, Partition: kafka.PartitionAny},
				Key:            []byte(recordKey),
				Value:          []byte(recordValue),
			}, nil)
			i++
		} else {
			fmt.Printf("Error: %s", err.Error())
			break
		}
	}

	// Wait for all messages to be delivered
	p.Flush(5 * 1000)

	fmt.Printf("%d messages were produced to topic %s!", i, connector.Table)

	p.Close()
	return len(records)
}

func getConfigMap(connector cdc_shared.Connector) *kafka.ConfigMap {
	config := &kafka.ConfigMap{
		"bootstrap.servers": connector.ConnectionString,
	}
	for key, value := range connector.Attributes {
		config.SetKey(key, value)
	}
	return config
}

func (reader KafkaConnector) GetRecords(connector cdc_shared.Connector, destinationProvider cdc_shared.ConnectorProvider, destinationConnector cdc_shared.Connector) {
	// Create Consumer instance
	c, err := reader.getConsumer(connector)
	custom_errors.CdcLog(connector, err)
	// Subscribe to topic
	err = c.SubscribeTopics([]string{connector.Table}, nil)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	var tp kafka.TopicPartition
	tp.Topic = &connector.Table
	tp.Offset = kafka.Offset(connector.StartOffset)
	pn, ok := connector.Attributes["PartitionNumber"]
	// If the key exists
	if ok {
		partitionNumber, _ := strconv.Atoi(pn)
		tp.Partition = int32(partitionNumber)
	} else {
		tp.Partition = 0
	}

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
	config := getConfigMap(connector)
	c, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}
	return c, err
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
