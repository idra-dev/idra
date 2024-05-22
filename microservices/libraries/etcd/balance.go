package etcd

import (
	"github.com/antrad1978/cdc_shared"
	"sort"
)


type WorkerNode struct {
	Name     string
	Capacity int
	Load     int
}

type ByCapacity []WorkerNode

func (s ByCapacity) Len() int {
	return len(s)
}

func (s ByCapacity) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCapacity) Less(i, j int) bool {
	return s[i].Capacity < s[j].Capacity
}

type LoadBalancer struct {
	servers []WorkerNode
	config map[string][]string
}

func NewLoadBalancer(servers []WorkerNode) *LoadBalancer {
	lb := &LoadBalancer{
		servers: servers,
	}
	sort.Sort(ByCapacity(lb.servers))
	return lb
}

func (lb *LoadBalancer) AddTasks(syncs []cdc_shared.Sync) map[string]string {
	result := make(map[string]string)
	for _, sync := range syncs {
		minLoad := lb.servers[0].Load
		minIndex := 0
		for j := 1; j < len(lb.servers); j++ {
			if lb.servers[j].Load < minLoad {
				minLoad = lb.servers[j].Load
				minIndex = j
			}
		}
		result[sync.Id] = lb.servers[minIndex].Name
		lb.servers[minIndex].Load++
	}
	sort.Sort(ByCapacity(lb.servers))
	lb.config = make(map[string][]string)
	for key, value := range result {
		_, ok := lb.config[value]
		if !ok {
			lb.config[value] = []string{}
		}
		lb.config[value] = append(lb.config[value], key)
	}
	return result
}
