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

func SyncData(sync cdc_shared.Sync, mode string) {
	providerSource := RetrieveProvider(sync.SourceConnector.ConnectorType)
	if checkProviderTypeIsDatabase(providerSource) {
		providerDestination := RetrieveProvider(sync.SourceConnector.ConnectorType)
		if providerSource != nil && providerDestination != nil {
			ProcessRDBMSProvider(sync, mode, providerSource.(cdc_shared.DatabaseConnectorProvider), providerDestination.(cdc_shared.DatabaseConnectorProvider))
		}
	} else {
		providerSource := RetrieveProvider(sync.SourceConnector.ConnectorType)
		providerSource.MoveData(sync)
	}
}

func ProcessRDBMSProvider(sync cdc_shared.Sync, mode string, providerSource cdc_shared.DatabaseConnectorProvider, providerDestination cdc_shared.DatabaseConnectorProvider) {
	switch {
	case mode == models.Id:
		SyncById(sync, providerSource, providerDestination.(cdc_shared.ConnectorProvider), sync.Id)
	case mode == models.Timestamp:
		SyncByTimestamp(sync, providerSource, providerDestination.(cdc_shared.ConnectorProvider), sync.Id)
	case mode == models.LastDestinationId:
		SyncByLastDestinationId(sync, providerDestination, providerSource)
	case mode == models.FullWithId:
		rows, _ := providerSource.GetRowsById(sync.SourceConnector, -1)
		providerDestination.InsertRows(sync.DestinationConnector, rows)
	case mode == models.LastDestinationTimestamp:
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
	case name == "MongodbManager":
		return MongodbManager{}
	case name == "S3":
		return S3JsonConnector{}
	case name == "Immudb":
		return ImmudbIdraDriver{}
	case name == "ChromaDb":
		return ChromaDb{}
	case name == "RestConnector":
		return RestConnector{}
	case name == "RabbiMQStreamConnector":
		return RabbiMQStreamConnector{}
	case name == "RabbiMQConnector":
		return RabbiMQConnector{}
	}
	//TODO: If missing search in plugins
	return nil
}

func GetProviders() []string {
	return []string{"PostgresGORM", "MysqlGORM", "MssqlGORM", "KafkaConnector", "MongodbManager", "S3", "Immudb", "ChromaDb", "RestConnector", "RabbiMQStreamConnector"}
}
