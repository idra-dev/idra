package data

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"io"
	"log"
	"math"
	"microservices/libraries/custom_errors"
	"microservices/libraries/models"
	"net/http"
	"strconv"
	"time"
)

type RestConnector struct{}

func (rdb RestConnector) MoveData(sync cdc_shared.Sync) {
}

func (RestConnector) Modes() []string {
	return []string{models.Id, models.Timestamp}
}

func (rdb RestConnector) Name() string {
	return "RestConnector"
}

func (rdb RestConnector) GetMaxTableId(connector cdc_shared.Connector) int64 {
	return -1
}

func (rdb RestConnector) GetMaxTimestamp(connector cdc_shared.Connector) (int64, error) {
	return -1, nil
}

func (rdb RestConnector) GetRowsById(connector cdc_shared.Connector, lastId int64) ([]map[string]interface{}, int64) {
	data, err := getData(connector, lastId, time.Time{})
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return nil, -1
	}
	length := int64(len(data))
	if len(data) > 0 {
		currentId, err := toInt64(data[length-1][connector.IdField])
		if err != nil {
			custom_errors.CdcLog(connector, err)
			return data, -1
		}
		return data, currentId
	}
	return data, -1
}

func (rdb RestConnector) GetRecordsByTimestamp(connector cdc_shared.Connector, lastTimestamp time.Time) ([]map[string]interface{}, time.Time) {
	data, err := getData(connector, -1, lastTimestamp)
	if err != nil {
		return nil, time.Time{}
	}
	length := int64(len(data))
	parsedTimeCustom, err := time.Parse(data[length-1][connector.TimestampFieldFormat].(string), data[length-1][connector.TimestampField].(string))
	if err != nil {
		fmt.Println("Error during conversion:", err)
		return nil, time.Time{}
	}
	return data, parsedTimeCustom
}

func (rdb RestConnector) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	postData(connector, rows)
	return len(rows)
}

func getData(connector cdc_shared.Connector, lastId int64, lastTimestamp time.Time) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", connector.ConnectionString, nil)
	q := req.URL.Query()
	if connector.IdField != "" {
		q.Set(connector.IdField, string(lastId))
	} else if connector.TimestampField != "" {
		q.Set(connector.TimestampField, string(lastId))
	}
	for key, value := range connector.Attributes {
		q.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()
	addAuth(connector, req)
	response, err := executeRequest(req, connector)
	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return nil, err
	}
	fmt.Println(string(responseData))
	var result []map[string]interface{}
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return result, err
	}

	fmt.Println(result)
	return result, err
}

func postData(connector cdc_shared.Connector, payload []map[string]interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return err
	}
	req, err := http.NewRequest("POST", connector.ConnectionString, bytes.NewBuffer(jsonData))
	addAuth(connector, req)
	response, err := executeRequest(req, connector)
	if err != nil {
		log.Println(response.StatusCode)
		custom_errors.CdcLog(connector, err)
		return err
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		custom_errors.CdcLog(connector, err)
		return err
	}
	// Manage errors
	if response.StatusCode != http.StatusOK {
		log.Println(response.StatusCode)
		custom_errors.CdcLog(connector, fmt.Errorf(string(responseData)))
		return fmt.Errorf(string(responseData))
	}
	fmt.Println(string(responseData))
	return err
}

func executeRequest(req *http.Request, connector cdc_shared.Connector) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(req)
	if response.StatusCode != http.StatusOK {
		log.Println(response.StatusCode)
		responseData, _ := io.ReadAll(response.Body)
		custom_errors.CdcLog(connector, fmt.Errorf(string(responseData)))
	}
	return response, err
}

func addAuth(connector cdc_shared.Connector, req *http.Request) {
	_, exists := connector.Attributes["username"]

	if exists {
		addBasicAuth(req, connector.Attributes["username"], connector.Attributes["password"])
	}
	if connector.Token != "" {
		setToken(req, connector.Token)
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func addBasicAuth(req *http.Request, username string, password string) {
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
}

func setToken(req *http.Request, token string) {
	req.Header.Add("Authorization", "Bearer "+token)
}

func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, fmt.Errorf("overflow: value %d is too large to fit in an int64", v)
		}
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid string value: %v", v)
			}
			return int64(f), nil
		}
		return i, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
