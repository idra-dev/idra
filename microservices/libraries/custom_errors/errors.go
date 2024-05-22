package custom_errors

import (
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"microservices/libraries"
	"microservices/libraries/models"
	"os"
	"time"
)

type CdcErrorHandler struct{}

type ExecutionError struct {
	Id             string `json:"id"`
	ConnectorName  string `json:"ip_address" binding:"required"`
	ConnectorType  string `json:"connector_type"`
	ErrorTimestamp int64  `json:"error_timestamp"`
	StackCall      string `json:"stack_call"`
	ErrorText      string `json:"error_text"`
}

func (handler CdcErrorHandler) GetErrors(lastKey string) ([]ExecutionError, string) {
	var items []ExecutionError
	ctx, cli := libraries.GetEtcd()
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	opts := libraries.GetOptionsRead()
	if lastKey != "" {
		opts = append(opts, clientv3.WithFromKey())
	} else {
		lastKey = models.ErrorsPath
	}
	gr, _ := kv.Get(ctx, lastKey, opts...)

	if gr.Kvs != nil {
		for _, kv := range gr.Kvs {
			error := ExecutionError{}
			json.Unmarshal(kv.Value, &error)
			items = append(items, error)
		}
	}
	if len(gr.Kvs) > 0 {
		lastKey = string(gr.Kvs[len(gr.Kvs)-1].Key)
	}
	return items, lastKey
}

func (handler CdcErrorHandler) DeleteErrors(identifiers []string) {
	for _, identifier := range identifiers {
		handler.DeleteError(identifier)
	}
}

func (handler CdcErrorHandler) DeleteError(id string) {
	key := models.ErrorsPath + id
	libraries.DeleteKey(key)
}

func (handler CdcErrorHandler) SaveError(error ExecutionError) {
	error.Id = uuid.New().String()
	key := models.ErrorsPath + error.Id
	data, _ := json.Marshal(error)
	libraries.SaveKey(key, data)
}

func (handler CdcErrorHandler) SaveCdcInstanceError(connector cdc_shared.Connector, errorText string) {
	error := ExecutionError{}
	error.ErrorTimestamp = time.Now().UnixMilli()
	error.ConnectorName = connector.ConnectorName
	error.StackCall = ""
	error.ErrorText = errorText
	if IsStaticRunMode() {
		loggedError, _ := json.Marshal(error)
		fmt.Println(string(loggedError))
	} else {
		handler.SaveError(error)
	}
}

func (handler CdcErrorHandler) SaveExecutionError(errorText string) {
	error := ExecutionError{}
	error.ErrorTimestamp = time.Now().UnixMilli()
	error.ErrorText = errorText
	if IsStaticRunMode() {
		loggedError, _ := json.Marshal(error)
		fmt.Println(string(loggedError))
	} else {
		handler.SaveError(error)
		time.Sleep(30 * time.Second)
	}
}

func CdcLog(connector cdc_shared.Connector, err error) {
	if err != nil {
		handler := CdcErrorHandler{}
		handler.SaveCdcInstanceError(connector, err.Error())
	}
}

func IsStaticRunMode() bool {
	value := os.Getenv(models.Static)
	return value == "true"
}

func LogAndDie(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
