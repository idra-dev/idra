package data

import (
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"microservices/libraries/models"
	"os"
	"strings"
	"time"
)

type S3JsonConnector struct {
	ConnectorId      string
	ConnectionString string
}

func (s3JsonConnector S3JsonConnector) InsertRows(connector cdc_shared.Connector, rows []map[string]interface{}) int {
	return -1
}

func (s3JsonConnector S3JsonConnector) MoveData(sync cdc_shared.Sync) {
	return
}

func (s3JsonConnector S3JsonConnector) Name() string {
	return "s3JsonConnector"
}

func (S3JsonConnector) Modes() []string {
	return []string{models.Default}
}

func GetSession(region string) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET"), ""),
	})
	if err != nil {
		log.Fatalln("S3 access error")
	}
	return sess
}

func (s3JsonConnector S3JsonConnector) GetData(bucketName string, lastChange int64) []map[string]interface{} {
	var results []map[string]interface{}
	return results
}

func (s3JsonConnector S3JsonConnector) SaveData(regionName string, bucketName string, rows []map[string]interface{}, destination string) int64 {
	session := GetSession(regionName)
	jsonString, err := json.Marshal(rows)
	if err != nil {
		fmt.Println(err)
	}
	reader := strings.NewReader(string(jsonString))

	uploader := s3manager.NewUploader(session)
	fileName := destination + time.Now().Format("20060102150405") + ".json"
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   reader,
	})
	if err != nil {
		fmt.Printf("Unable to upload %q to %q, %v", fileName, bucketName, err)
		return 0
	}
	return int64(len(rows))
}
