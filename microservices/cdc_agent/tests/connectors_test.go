package tests

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"microservices/libraries/data"
	"microservices/libraries/models"
	"os"
	"shared"
	"testing"
)

const QuerySelect = "SELECT * FROM `pratica`"
var MysqlConnectionString = os.Getenv("SAMPLE_DB_MYSQL")
var ConnectionStringDestination = os.Getenv("LOCAL_MYSQL_SAMPLE2")
var MysqlConnectionStringDestination = os.Getenv("DESTINATION_SAMPLE_DB_MYSQL")
var MysqlConnectionStringSource = os.Getenv("SOURCE_SAMPLE_DB_MYSQL")
var MysqlLocalConnection = os.Getenv("LOCAL_MYSQL_SAMPLE")
var MysqlTestConnectionString = os.Getenv("TEST_DB_MYSQL")


func TestPostgresConnectorLastDestinationId(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "PostgresGORM"
	connector.ConnectionString = MysqlLocalConnection
	connector.Table = "table"
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "PostgresGORM"
	connector2.ConnectionString = ConnectionStringDestination
	connector2.Table = "table"
	connector2.IdField = "id"

	data.SyncData(connector,connector2, models.LastDestinationId)
}

func TestMysqlConnectorLastDestinationId(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "MysqlGORM"
	connector.ConnectionString = MysqlTestConnectionString
	connector.Table = "comune"
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "MysqlGORM"
	connector2.ConnectionString = MysqlLocalConnection
	connector2.Table = "comune"
	connector2.IdField = "id"
	connector2.SaveMode = models.Insert

	data.SyncData(connector,connector2, models.LastDestinationId)
}

func TestPostgresConnectorLastOffsetId(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "PostgresGORM"
	connector.ConnectionString = MysqlLocalConnection
	connector.Table = "table"
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "PostgresGORM"
	connector2.ConnectionString = ConnectionStringDestination
	connector2.Table = "table"
	connector2.IdField = "id"

	data.SyncData(connector,connector2, models.Id)
}

func TestPostgresConnectorLastTimestamp(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "PostgresGORM"
	connector.ConnectionString = MysqlLocalConnection
	connector.Table = "table"
	connector.TimestampField = "time"
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "PostgresGORM"
	connector2.ConnectionString = ConnectionStringDestination
	connector2.Table = "table"
	connector2.TimestampField = "time"
	connector2.IdField = "id"

	data.SyncData(connector,connector2, models.LastDestinationTimestamp)
}

func TestPostgresConnectorQueryLastOffsetId(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "PostgresGORM"
	connector.ConnectionString = MysqlLocalConnection
	connector.Query = "select * from \"table\""
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "PostgresGORM"
	connector2.ConnectionString = ConnectionStringDestination
	connector2.Table = "table"
	connector2.IdField = "id"

	data.SyncData(connector,connector2, models.Id)
}


func TestPostgresConnectorTimestamp(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "PostgresGORM"
	connector.ConnectionString = MysqlLocalConnection
	connector.Table = "table"
	connector.TimestampField = "time"
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "PostgresGORM"
	connector2.ConnectionString = ConnectionStringDestination
	connector2.Table = "table"
	connector2.TimestampField = "time"
	connector2.IdField = "id"

	data.SyncData(connector,connector2, models.Timestamp)
}


func TestMysqlConnectorQueryLastDestinationId(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "MysqlGORM"
	connector.ConnectionString = MysqlTestConnectionString
	connector.Query = "SELECT id,data_nascita,sesso,id_comune_nascita,id_comune,marketing_privacy,data_inserimento,id_provincia,cessione_a_terzi_privacy,id_storico,last_update FROM `anagrafica`"
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "MysqlGORM"
	connector2.ConnectionString = MysqlLocalConnection
	connector2.Table = "anagrafica"
	connector2.IdField = "id"
	connector2.SaveMode = models.Insert

	data.SyncData(connector,connector2, models.LastDestinationId)
}


func TestMysqlConnectorQueryLastDestinationTimestamp(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "MysqlGORM"
	connector.ConnectionString = MysqlConnectionString
	connector.Query = QuerySelect
	connector.TimestampField = "last_update"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "MysqlGORM"
	connector2.ConnectionString = MysqlConnectionStringDestination
	connector2.Table = "pratica"
	connector2.TimestampField = "last_update"
	connector2.SaveMode = models.Upsert
	connector2.IdField = "id"

	data.SyncData(connector,connector2, models.LastDestinationTimestamp)
}

func TestMysqlConnectorFullById(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "MysqlGORM"
	connector.ConnectionString = MysqlConnectionStringSource
	connector.Table = "contatto_status"
	connector.MaxRecordBatchSize = math.MaxInt
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "MysqlGORM"
	connector2.ConnectionString = MysqlConnectionStringDestination
	connector2.Table = "contatto_status"
	connector2.IdField = "id"
	connector2.SaveMode = models.Insert

	data.SyncData(connector,connector2, models.LastDestinationId)
}


func TestMysqlConnectorJson(t *testing.T){
	sync := shared.Sync{}
	connector := shared.Connector{}
	connector.ConnectorName = "sample"
	connector.ConnectorType = "MysqlGORM"
	connector.ConnectionString = MysqlConnectionStringSource
	connector.Table = "campagna"
	connector.MaxRecordBatchSize = math.MaxInt
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "sample"
	connector2.ConnectorType = "MysqlGORM"
	connector2.ConnectionString = MysqlLocalConnection
	connector2.Table = "campagna"
	connector2.IdField = "id"
	connector2.SaveMode = models.Upsert

	sync.SourceConnector = connector
	sync.DestinationConnector = connector2
	sync.SyncName="Campagna"
	sync.Id = uuid.New().String()
	sync.Mode = models.FullWithId

	b, err := json.Marshal(sync)
	if err != nil {
		fmt.Println(sync)
		return
	}
	fmt.Println(string(b))
}

func TestMysqlConnectorFromJson(t *testing.T){
	var sync []shared.Sync
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	dat, err2 := os.ReadFile("../static_configs/static.json")
	if err2 != nil {
		log.Fatal(err2)
	}

	// `&myStoredVariable` is the address of the variable we want to store our
	// parsed data in
	json.Unmarshal([]byte(string(dat)), &sync)
	fmt.Println(sync[0])
	fmt.Println(len(sync))
}

func TestMysqlConnectorTableLastDestinationId(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "pratica"
	connector.ConnectorType = "MysqlGORM"
	connector.ConnectionString = MysqlConnectionStringSource
	connector.Table = QuerySelect
	connector.TimestampField = "last_update"
	connector.IdField = "id"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "pratica"
	connector2.ConnectorType = "MysqlGORM"
	connector2.ConnectionString = MysqlConnectionStringDestination
	connector2.Table = "pratica"
	connector2.SaveMode = models.Upsert
	connector2.TimestampField = "last_update"
	connector2.IdField = "id"

	data.SyncData(connector,connector2, models.LastDestinationTimestamp)
}

func TestMysqlConnectorTableJson(t *testing.T){
	connector := shared.Connector{}
	connector.ConnectorName = "pratica_source"
	connector.ConnectorType = "MysqlGORM"
	connector.ConnectionString = MysqlConnectionStringSource
	connector.Query = QuerySelect
	connector.TimestampField = "last_update"


	connector2 := shared.Connector{}
	connector2.ConnectorName = "pratica_dest"
	connector2.ConnectorType = "MysqlGORM"
	connector2.ConnectionString = MysqlConnectionStringDestination
	connector2.Table = "pratica"
	connector2.TimestampField = "last_update"
	connector2.SaveMode = models.Upsert
	connector2.IdField = "id"

	sync := shared.Sync{}
	sync.SourceConnector = connector
	sync.DestinationConnector = connector2
	sync.SyncName="pratica"
	sync.Id = uuid.New().String()
	sync.Mode = models.LastDestinationTimestamp

	b, err := json.Marshal(sync)
	if err != nil {
		fmt.Println(sync)
		return
	}
	fmt.Println(string(b))
}


func TestGetAgentInfo(t *testing.T){
	agent := models.GetCurrentAgentInfo()
	fmt.Println(agent.AgentId)
}