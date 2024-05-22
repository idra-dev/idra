package tests

import (
	"log"
	"microservices/libraries/custom_errors"
	"testing"
)

func TestErrors(t *testing.T){
	handler := custom_errors.CdcErrorHandler{}
	error := custom_errors.ExecutionError{}
	error.ErrorTimestamp=7586575
	error.ConnectorName="connector"
	handler.SaveError(error)
	res, _ := handler.GetErrors("")
	if len(res)!=1{
		log.Fatalln("Test failed")
	}
	handler.DeleteError(res[0].Id)
	res, _ = handler.GetErrors("")
	if len(res)!=0{
		log.Fatalln("Test failed")
	}
}
