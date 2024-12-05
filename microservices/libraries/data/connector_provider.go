package data

import (
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/models"
	"time"
)

func checkProviderTypeIsDatabase(i interface{}) bool {
	switch i.(type) {
	case cdc_shared.DatabaseConnectorProvider:
		return true
	default:
		return false
	}
}

func SyncData(sync cdc_shared.Sync) {
	providerSource := RetrieveProvider(sync.SourceConnector.ConnectorType)
	if checkProviderTypeIsDatabase(providerSource) {
		providerDestination := RetrieveProvider(sync.SourceConnector.ConnectorType)
		if providerSource != nil && providerDestination != nil {
			ProcessRDBMSProvider(sync, providerSource.(cdc_shared.DatabaseConnectorProvider), providerDestination.(cdc_shared.DatabaseConnectorProvider))
		}
	} else {
		providerSource := RetrieveProvider(sync.SourceConnector.ConnectorType)
		providerSource.MoveData(sync)
	}
}

func ProcessRDBMSProvider(sync cdc_shared.Sync, providerSource cdc_shared.DatabaseConnectorProvider, providerDestination cdc_shared.DatabaseConnectorProvider) {
	switch {
	case sync.Mode == models.Id:
		SyncById(sync, providerSource, providerDestination.(cdc_shared.ConnectorProvider), sync.Id)
	case sync.Mode == models.Timestamp:
		SyncByTimestamp(sync, providerSource, providerDestination.(cdc_shared.ConnectorProvider), sync.Id)
	case sync.Mode == models.LastDestinationId:
		SyncByLastDestinationId(sync, providerDestination, providerSource)
	case sync.Mode == models.FullWithId:
		rows, _ := providerSource.GetRowsById(sync.SourceConnector, -1)
		providerDestination.InsertRows(sync.DestinationConnector, rows)
	case sync.Mode == models.LastDestinationTimestamp:
		SyncByLastDestinationTimestamp(sync.SourceConnector, sync.DestinationConnector, providerDestination, providerSource)
	}
	time.Sleep(time.Duration(time.Millisecond.Milliseconds() * int64(sync.SourceConnector.PollingTime)))
}

func RetrieveProvider(name string) cdc_shared.ConnectorProvider {
	switch {
	case name == "PostgresGORM":
		return PostgresGormManager{}
	case name == "MysqlGORM":
		return MysqlConnector{}
	case name == "MssqlGORM":
		return MssqlManager{}
	case name == "KafkaConnector":
		return KafkaConnector{}
	case name == "MongodbConnector":
		return MongodbConnector{}
	case name == "s3JsonConnector":
		return S3JsonConnector{}
	case name == "ImmudbDriver":
		return ImmudbDriver{}
	case name == "ChromaDbConnector":
		return ChromaDbConnector{}
	case name == "RestConnector":
		return RestConnector{}
	case name == "RabbitMQStreamConnector":
		return RabbitMQStreamConnector{}
	case name == "RabbitMQConnector":
		return RabbitMQConnector{}
	}
	//TODO: If missing search in plugins
	return nil
}

func GetProviders() []string {
	return []string{"PostgresGORM", "MysqlGORM", "MssqlGORM", "KafkaConnector", "MongodbConnector", "S3", "Immudb", "ChromaDbConnector", "RestConnector", "RabbitMQStreamConnector"}
}
