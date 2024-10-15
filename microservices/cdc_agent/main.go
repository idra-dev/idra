package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"microservices/cdc_agent/processing"
)

func main() {
	err := godotenv.Load("/Users/antonioradesca/Code/idra/microservices/cdc_agent/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	//if custom_errors.IsStaticRunMode() {
	//	processing.ProcessStatic()
	//} else {
	processing.StartWorkerNode()
	//}
}
