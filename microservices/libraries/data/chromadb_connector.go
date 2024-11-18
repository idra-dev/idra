package data

import (
	"context"
	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/hf"
	"github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/antrad1978/cdc_shared"
	"log"
	"microservices/libraries/models"
	"time"
)

type ChromaDbConnector struct {
	token            string
	tenant           string
	embeddedFunction string //hash:NewConsistentHashEmbeddingFunction, huggingface:hf.NewHuggingFaceEmbeddingFunction, openai: openai.NewOpenAIEmbeddingFunction
	model            string
	providerApiKey   string
	metaDataKey      string
	metaDataValue    string
	column           string
	distanceFunction string //L2 DistanceFunction = "l2", COSINE DistanceFunction = "cosine", IP DistanceFunction = "ip"
}

func (rdb ChromaDbConnector) MoveData(sync cdc_shared.Sync) {
}

func (ChromaDbConnector) Modes() []string {
	return []string{models.Id, models.Timestamp}
}

func (rdb ChromaDbConnector) Name() string {
	return "ChromaDbConnector"
}

func (rdb ChromaDbConnector) GetMaxTableId(connector cdc_shared.Connector) int64 {
	return -1
}

func (rdb ChromaDbConnector) GetMaxTimestamp(connector cdc_shared.Connector) (int64, error) {
	return -1, nil
}

func (rdb ChromaDbConnector) GetRowsById(connector cdc_shared.Connector, lastId int64) ([]map[string]interface{}, int64) {
	return nil, -1
}

func (rdb ChromaDbConnector) GetRecordsByTimestamp(connector cdc_shared.Connector, lastTimestamp time.Time) ([]map[string]interface{}, int64) {
	return nil, -1
}

func (rdb ChromaDbConnector) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	var err error
	ef, err := rdb.getEmbeddingFunction(connector, err)
	tenant := func() string {
		if connector.Attributes["Tenant"] != "" {
			return connector.Attributes["Tenant"]
		}
		return types.DefaultTenant
	}()
	database := func() string {
		if connector.Database != "" {
			return connector.Database
		}
		return types.DefaultDatabase
	}()

	distanceFunction := func() types.DistanceFunction {
		if connector.Attributes["DistanceFunction"] != "" {
			return types.DistanceFunction(connector.Attributes["DistanceFunction"])
		}
		return types.L2
	}()

	client, err := rdb.getClient(connector, err)
	if client != nil {
		client.SetTenant(tenant)
		client.SetDatabase(database)
	} else {
		log.Fatalf("Nil client: %s \n", err)
	}
	if err != nil {
		log.Fatalf("Error creating client: %s \n", err)
	}
	client.Heartbeat(context.Background())

	collections, err := client.ListCollections(context.Background())
	if err != nil {
		log.Fatalf("Error creating collection: %s \n", err)
	}
	var newCollection *chroma.Collection = nil

	if !exists(connector.Table, collections) {
		newCollection, err = client.NewCollection(
			context.TODO(),
			collection.WithName(connector.Table),
			collection.WithMetadata(connector.Attributes["Key"], connector.Attributes["Value"]),
			collection.WithEmbeddingFunction(ef),
			collection.WithHNSWDistanceFunction(distanceFunction),
		)
		if err != nil {
			log.Fatalf("Error creating collection: %s \n", err)
		}
	} else {
		newCollection, err = client.GetCollection(
			context.TODO(),
			connector.Table,
			ef,
		)
		if err != nil {
			log.Fatalf("Error creating collection: %s \n", err)
		}
	}
	// Create a new record set with to hold the records to insert
	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(ef),
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}
	// Add a few records to the record set
	for _, row := range rows {
		value := row[connector.Attributes["Column"]]
		rs.WithRecord(types.WithDocument(value.(string)), types.WithMetadata(connector.Attributes["Key"], connector.Attributes["Value"]))
	}

	// Build and validate the record set (this will create embeddings if not already present)
	_, err = rs.BuildAndValidate(context.TODO())
	if err != nil {
		log.Fatalf("Error validating record set: %s \n", err)
	}

	// Add the records to the collection
	_, err = newCollection.AddRecords(context.Background(), rs)
	if err != nil {
		log.Fatalf("Error adding documents: %s \n", err)
	}
	return len(rows)
}

func (rdb ChromaDbConnector) getClient(connector cdc_shared.Connector, err error) (*chroma.Client, error) {
	var client *chroma.Client
	if connector.Attributes["Token"] != "" {
		var defaultHeaders = map[string]string{"Authorization": "Bearer " + connector.Attributes["Token"]}
		client, err = chroma.NewClient(connector.ConnectionString, chroma.WithDefaultHeaders(defaultHeaders))
	} else if connector.Attributes["Username"] != "" {
		client, err = chroma.NewClient(connector.ConnectionString, chroma.WithAuth(types.NewBasicAuthCredentialsProvider(connector.Attributes["Username"], connector.Attributes["Password"])))
	} else {
		client, err = chroma.NewClient(connector.ConnectionString)
	}
	return client, err
}

func (rdb ChromaDbConnector) getEmbeddingFunction(connector cdc_shared.Connector, err error) (types.EmbeddingFunction, error) {
	var ef types.EmbeddingFunction

	switch connector.Attributes["ProviderApiKey"] {
	case "huggingface":
		ef = hf.NewHuggingFaceEmbeddingFunction(connector.Attributes["ApiKey"], connector.Attributes["Model"])
	case "openapi":
		ef, err = openai.NewOpenAIEmbeddingFunction(connector.Attributes["ProviderApiKey"])
		if err != nil {
			log.Fatalf("Error creating OpenAI embedding function: %s \n", err)
		}
	default:
		ef = types.NewConsistentHashEmbeddingFunction()
	}
	return ef, err
}

func exists(value string, data []*chroma.Collection) (exists bool) {
	for _, search := range data {
		if search.Name == value {
			return true
		}
	}
	return false
}
