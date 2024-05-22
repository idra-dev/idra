package tests


import (
	"fmt"
	"microservices/libraries/etcd"
	"testing"
	"time"
)

func TestIntFromTime(t *testing.T){
	startTime := time.Now()
	res := etcd.GetInt64FromTime(startTime)
	fmt.Println(res)
	time := etcd.GetTimeFromInt64(res)
	fmt.Println(time)
	if res!= etcd.GetInt64FromTime(time){
		t.Fail()
	}
}