package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"microservices/libraries"
	"microservices/libraries/custom_errors"
	"microservices/libraries/data"
	"microservices/libraries/etcd"
	"microservices/libraries/models"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

var SyncExecutions = make(map[string]struct {
	cancel context.CancelFunc
	status string
})

func StartWorkerNode() {
	var wg sync.WaitGroup
	lm := libraries.LeaseManager{}
	session := lm.GetLeasedSession()
	wg.Add(1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			wg.Done()
		}
	}()

	fmt.Println("Start processing")
	agentData := models.GetCurrentAgentInfo()
	lease, err := session.Client().Grant(context.Background(), 30)
	custom_errors.LogAndDie(err)
	keyPrefix := "agents/"

	key := keyPrefix + strconv.FormatInt(int64(lease.ID), 10)
	data, _ := json.Marshal(agentData)
	value := string(data)
	_, err = session.Client().Put(context.Background(), key, value, clientv3.WithLease(lease.ID))

	RenewLease(session, lease)

	AllocateSyncs(session)

	manager := Manager{}
	manager.startWorker(RunSyncs)
	if session == nil {
		log.Printf("session nil'\n")
	}
	manager.ListenGloballyBalanceEvent(session, keyPrefix)

	time.Sleep(2 * time.Second)
	fmt.Println("Continue...")
	wg.Wait()
	session.Done()
	session.Close()
	fmt.Println("Terminating...")
}

func RunSyncs() {
	go CheckExecutions()
	ExecuteSyncs()
}

// RenewLease Renew periodically lease to show that agent is running correctly
func RenewLease(session *concurrency.Session, lease *clientv3.LeaseGrantResponse) {
	go func() {
		for {
			_, err := session.Client().KeepAliveOnce(context.Background(), lease.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Lease renew executed.")
			time.Sleep(20 * time.Second)
		}
	}()
}

func CreateLoadBalancer(agents []models.CdcAgent) *etcd.LoadBalancer {
	var nodes []etcd.WorkerNode
	for _, agent := range agents {
		node := etcd.WorkerNode{}
		node.Name = agent.AgentId
		node.Load = 0
		node.Capacity = 100
		nodes = append(nodes, node)
	}
	lb := etcd.NewLoadBalancer(nodes)
	syncs := etcd.GetSyncs()
	lb.AddTasks(syncs)
	return lb
}

func ExecuteSync(sync cdc_shared.Sync) {
	cli, _ := libraries.GetClient()
	defer cli.Close()

	fmt.Println("Executing sync " + sync.SyncName)
	ctx, cancel := data.SyncData(sync)
	SyncExecutions[sync.Id] = struct {
		cancel context.CancelFunc
		status string
	}{cancel, "running"}
	if ctx == nil && cancel == nil {
		time.Sleep(30 * time.Second)
	}
}

func ExecuteSyncs() {
	var syncs []cdc_shared.Sync
	id := libraries.GetMachineId()
	syncsIds := libraries.GetKeys(models.AssignmentsPath + id)
	if len(syncsIds.Kvs) > 0 {
		for _, syncId := range syncsIds.Kvs {
			parts := strings.Split(string(syncId.Key), "/")
			sync, err := etcd.GetSync(parts[3])

			if err == nil && !sync.Disabled {
				syncs = append(syncs, sync)
			}
		}
	}
	if len(syncs) > 0 {
		for _, sync := range syncs {
			go ExecuteSync(sync)
		}
		time.Sleep(5 * time.Second)
	}
}

func CheckExecutions() {
	for {
		for key := range SyncExecutions {
			if goroutine, exists := SyncExecutions[key]; exists {
				fmt.Printf("Goroutine ID %s\n", key)
				fmt.Printf("Goroutine status %s\n", goroutine.status)
				fmt.Printf("startChan %s\n", goroutine.cancel)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
