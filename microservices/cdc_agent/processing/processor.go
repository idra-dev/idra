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
	startChan chan bool
	stopChan  chan bool
	status    string
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

	manager := Manager{}
	manager.startWorker(RunSyncs)
	if session == nil {
		log.Printf("session nil'\n")
	}
	manager.ObserveDiedAgentEvent(session, keyPrefix)
	manager.ListenGloballyBalanceEvent(session, keyPrefix)

	key := keyPrefix + strconv.FormatInt(int64(lease.ID), 10)
	data, _ := json.Marshal(agentData)
	value := string(data)
	_, err = session.Client().Put(context.Background(), key, value, clientv3.WithLease(lease.ID))

	RenewLease(session, lease)
	time.Sleep(2 * time.Second)
	fmt.Println("Continue...")
	wg.Wait()
	session.Done()
	session.Close()
	fmt.Println("Terminating...")
}

func RunSyncs(stop chan bool) {
	for {
		select {
		case <-stop:
			fmt.Println("Force stop for rebalance event...")
			time.Sleep(5 * time.Second)
		default:
			ExecuteSyncs()
		}
	}
}

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

func ExecuteSync(sync cdc_shared.Sync, startChan chan bool, stopChan chan bool, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error in ProcessSync goroutine:", r)
			time.Sleep(30 * time.Second)
		}
	}()
	defer wg.Done()
	name := sync.Id
	cli, _ := libraries.GetClient()
	defer cli.Close()
	// create a sessions to acquire a lock
	s, err := concurrency.NewSession(cli)
	custom_errors.LogAndDie(err)

	defer s.Close()
	mutex := concurrency.NewMutex(s, "/distributed-locks/"+name)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := mutex.Lock(ctx); err != nil {
		fmt.Println("Lock failed")
		return
	}
	fmt.Println("Acquired lock for ", name)
	for {
		select {
		case <-startChan:
			// Start
			fmt.Printf("Goroutine %s started\n", sync.SyncName)
			for {
				select {
				case <-stopChan:
					// Stop
					fmt.Printf("Goroutine %s stopped\n", sync.SyncName)
					SyncExecutions[sync.Id] = struct {
						startChan chan bool
						stopChan  chan bool
						status    string
					}{startChan, stopChan, "stopped"}
					StopSync(sync)
					return
				default:
					data.SyncData(sync)
					SyncExecutions[sync.Id] = struct {
						startChan chan bool
						stopChan  chan bool
						status    string
					}{startChan, stopChan, "running"}
					time.Sleep(30 * time.Second)
				}
			}
		}
	}

	s.Orphan()
	fmt.Println("Data processed for sync: ", name+" "+sync.SyncName)
	if err := mutex.Unlock(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Released lock for ", name)
}

func ExecuteSyncs() {
	var wg sync.WaitGroup
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
			go StartSync(sync, &wg)
		}
		wg.Wait()
		time.Sleep(5 * time.Second)
	}
}

func StartSync(sync cdc_shared.Sync, wg *sync.WaitGroup) {
	startChan := make(chan bool)
	stopChan := make(chan bool)
	SyncExecutions[sync.Id] = struct {
		startChan chan bool
		stopChan  chan bool
		status    string
	}{startChan, stopChan, "stopped"}

	wg.Add(1)
	go ExecuteSync(sync, startChan, stopChan, wg)
	startChan <- true
}

func StopSync(sync cdc_shared.Sync) {
	if goroutine, exists := SyncExecutions[sync.Id]; exists {
		goroutine.stopChan <- true
		fmt.Printf("Goroutine ID %s stopped\n", sync.Id)
	}
}

func StopAllSyncs() {
	for key := range SyncExecutions {
		if goroutine, exists := SyncExecutions[key]; exists {
			goroutine.stopChan <- true
			fmt.Printf("Goroutine ID %s stopped\n", key)
		}
	}
}
