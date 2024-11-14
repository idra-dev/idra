package data

import (
	"context"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"microservices/libraries/etcd"
	"microservices/libraries/models"
	"reflect"
	"sync"
	"time"
)

type MongodbManager struct {
	Token       string
	ConnectorId string
}

func (MongodbManager) Modes() []string {
	return []string{models.Default}
}

func (mdb MongodbManager) MoveData(sync cdc_shared.Sync) {
	mdb.GetRowsByToken(sync.SourceConnector, sync.DestinationConnector)
}

func (mdb MongodbManager) Name() string {
	return "MongoDB"
}

func GetMongoClient(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}
	return client, err
}

func (mdb MongodbManager) GetRowsByToken(connector cdc_shared.Connector, destinationConnector cdc_shared.Connector) ([]map[string]interface{}, string) {
	client, _ := GetMongoClient(connector.ConnectionString)
	mdb.ConnectorId = connector.IdField
	mdb.Token = etcd.GetOffsetToken(mdb.ConnectorId)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalln(err)
		}
	}()
	collection := client.Database(connector.Database).Collection(connector.Table)
	var waitGroup sync.WaitGroup
	cso := options.ChangeStream()
	cso.SetResumeAfter(bson.D{{"_data", mdb.Token}})
	var changeStream *mongo.ChangeStream
	var err interface{}
	if mdb.Token != "" {
		changeStream, err = collection.Watch(context.TODO(), mongo.Pipeline{}, cso)
	} else {
		changeStream, err = collection.Watch(context.TODO(), mongo.Pipeline{})
	}

	if err != nil {
		log.Fatalln(err)
	}
	waitGroup.Add(1)
	routineCtx, _ := context.WithCancel(context.Background())
	go mdb.IterateChangeStream(routineCtx, waitGroup, changeStream, destinationConnector)

	waitGroup.Wait()
	return nil, ""
}

func (mdb MongodbManager) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := GetMongoClient(connector.ConnectionString)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalln(err)
		}
	}()
	collection := client.Database(connector.Database).Collection(connector.Table)
	var insertedIDs []interface{}
	for _, document := range rows {
		insertResult, err := collection.InsertOne(ctx, document)
		if err != nil {
			log.Fatal(err)
		}
		insertedIDs = append(insertedIDs, insertResult.InsertedID)
	}
	return 1
}

func (mdb MongodbManager) IterateChangeStream(routineCtx context.Context, waitGroup sync.WaitGroup, stream *mongo.ChangeStream, destinationConnector cdc_shared.Connector) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()
	provider := RetrieveProvider(destinationConnector.ConnectorType)
	for stream.Next(routineCtx) {
		change := stream.Current
		fmt.Printf("%+v\n", change)
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%v\n", data)
		mdb.Token = fmt.Sprintf("%v", data["_id"])
		operation := data["operationType"]
		res := make(map[string]interface{})
		res["operationType"] = operation
		if operation == "insert" || operation == "update" {
			iter := reflect.ValueOf(data["fullDocument"]).MapRange()
			for iter.Next() {
				key := iter.Key().String()
				value := iter.Value().Interface()
				res[key] = value
			}
		} else {
			res["_id"] = data["_id"]
			fmt.Printf("%v\n", data)
		}
		rows := make([]map[string]interface{}, 1)
		rows[0] = res
		inserted := provider.InsertRows(destinationConnector, rows)
		fmt.Println(inserted)
		etcd.SetOffsetToken(mdb.ConnectorId, mdb.Token)
	}
}
