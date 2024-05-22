package tests

import (
	"fmt"
	"microservices/cdc_agent/processing"
	"microservices/libraries"
	"microservices/libraries/etcd"
	"microservices/libraries/models"
	"strconv"
	"testing"
)

func TestBalance(t *testing.T){
	var sampleManager libraries.LeaseManager
	session := sampleManager.GetLeasedSession()
	res := sampleManager.Leases(session)
	fmt.Println(len(res))
	var nodes []etcd.WorkerNode
	for _, item := range res {
		node := etcd.WorkerNode{}
		node.Name = strconv.FormatInt(int64(item.ID), 10)
		node.Load = 0
		node.Capacity = 999
		nodes = append(nodes, node)
	}
	fmt.Println(len(res))
}

func TestBalanceNode(t *testing.T){
	agent := models.CdcAgent{}
	agent.AgentId = "2334153"
	agent2 := models.CdcAgent{}
	agent2.AgentId = "11111"
	agents := []models.CdcAgent{ agent, agent2 }
	_ := processing.CreateLoadBalancer(agents)
}

