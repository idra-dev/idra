package models

import (
	"microservices/libraries"
	"net"
	"os"
	"runtime"
)

type CdcAgent struct {
	AgentId string `json:"agent_id" binding:"required"`
	HostName string `json:"host_name" binding:"required"`
	IpAddress string `json:"ip_address" binding:"required"`
	WorkingConnectorId string `json:"working_connector_id"`
	LastErrorTimestamp string `json:"last_error_timestamp"`
	LastError string `json:"last_error"`
	AllocatedMemory uint64 `json:"allocated_memory"`
	TotalAllocatedMemory uint64 `json:"total_allocated_memory"`
	SystemdMemory uint64 `json:"system_memory"`
}


func GetCurrentAgentInfo() CdcAgent {
	agent := CdcAgent{}
	hostname, _ := os.Hostname()
	agent.HostName = hostname
	id := libraries.GetMachineId()
	agent.AgentId = id
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	agent.AllocatedMemory = bToMb(m.Alloc)
	agent.TotalAllocatedMemory = bToMb(m.TotalAlloc)
	agent.SystemdMemory = bToMb(m.Sys)
	//fmt.Printf("\tNumGC = %v\n", m.NumGC)

	// perform lookup of the hostname
	addresses, _ := net.LookupHost(hostname)

	// fetch ip address details
	for _, addr := range addresses {
		agent.IpAddress = addr
		return agent
	}
	return agent
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}