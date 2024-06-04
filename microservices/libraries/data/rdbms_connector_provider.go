package data

import (
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/custom_errors"
	"microservices/libraries/etcd"
)

var RDBMS_PROVIDERS = [3]string{"PostgresGORM", "MysqlGORM", "MssqlGORM"}

func SyncByLastDestinationTimestamp(connectorSource cdc_shared.Connector, connectorDestination cdc_shared.Connector, providerDestination cdc_shared.DatabaseConnectorProvider, providerSource cdc_shared.DatabaseConnectorProvider) {
	max, err := providerDestination.GetMaxTimestamp(connectorDestination)
	if err != nil {
		custom_errors.CdcLog(connectorDestination, err)
	}
	rows, _ := providerSource.GetRecordsByTimestamp(connectorSource, max)
	providerDestination.InsertRows(connectorDestination, rows)
}

func SyncByLastDestinationId(sync cdc_shared.Sync, providerDestination cdc_shared.DatabaseConnectorProvider, providerSource cdc_shared.DatabaseConnectorProvider) {
	max := providerDestination.GetMaxTableId(sync.DestinationConnector)
	if max > -1 {
		rows, _ := providerSource.GetRowsById(sync.SourceConnector, max)
		providerDestination.InsertRows(sync.DestinationConnector, rows)
	}
}

func SyncByTimestamp(sync cdc_shared.Sync, providerSource cdc_shared.DatabaseConnectorProvider, providerDestination cdc_shared.ConnectorProvider, syncId string) {
	manager := etcd.IntOffset{}
	max := manager.GetOffsetId(sync.DestinationConnector.ConnectorName)
	timeStamp := etcd.GetTimeFromInt64(max)
	rows, offset := providerSource.GetRecordsByTimestamp(sync.SourceConnector, timeStamp)
	if !offset.IsZero() {
		providerDestination.InsertRows(sync.DestinationConnector, rows)
		offsetManager := etcd.IntOffset{}
		offsetInt := etcd.GetInt64FromTime(offset)
		offsetManager.SetOffsetId(syncId, offsetInt)
	}
}

func SyncById(sync cdc_shared.Sync, providerSource cdc_shared.DatabaseConnectorProvider, providerDestination cdc_shared.ConnectorProvider, syncId string) {
	offsetSource := etcd.IntOffset{}
	max := offsetSource.GetOffsetId(sync.SourceConnector.ConnectorName)
	rows, offset := providerSource.GetRowsById(sync.SourceConnector, max)
	if offset > 0 {
		providerDestination.InsertRows(sync.DestinationConnector, rows)
		offsetSource.SetOffsetId(syncId, offset)
	}
}

func RetrieveDatabaseProvider(name string) cdc_shared.DatabaseConnectorProvider {
	switch {
	case name == "PostgresGORM":
		return PostgresGormManager{}
	case name == "MysqlGORM":
		return MysqlConnector{}
	case name == "MssqlGORM":
		return MssqlManager{}
	}
	return nil
}

func IsInDBMS(target string) bool {
	for _, element := range RDBMS_PROVIDERS {
		if element == target {
			return true
		}
	}
	return false
}
