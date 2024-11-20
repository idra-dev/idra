package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"microservices/cdc_agent/processing"
	"microservices/libraries/custom_errors"
	"net/http"
)

func main() {
	err := godotenv.Load("/Users/antonioradesca/Code/idra/microservices/cdc_agent/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if custom_errors.IsStaticRunMode() {
		processing.ProcessStatic()
	} else {
		processing.StartWorkerNode()
	}
}
