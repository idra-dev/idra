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

var SyncsWaitGroup sync.WaitGroup

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
	// Create a session to acquire the lock
	s, err := concurrency.NewSession(cli)
	custom_errors.LogAndDie(err)
	defer s.Close()

	// Create a mutex for locking
	mutex := concurrency.NewMutex(s, "/distributed-locks/"+name)

	// Retry logic to acquire lock
	for {
		// Create a context with a timeout for the lock acquisition
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel() // ensure cancel is called once
		if err := mutex.Lock(ctx); err != nil {
			fmt.Println("Lock failed, retrying in 30 seconds...")
			time.Sleep(30 * time.Second) // Wait before retrying
			continue                     // Retry acquiring the lock
		}
		// Lock acquired, break out of retry loop
		fmt.Println("Acquired lock for ", name)
		break
	}

	// Main logic when lock is acquired
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
					StopSync(sync.Id)
					// Release the lock and exit the function
					if err := mutex.Unlock(context.Background()); err != nil {
						log.Fatal(err)
					}
					fmt.Println("Released lock for ", name)
					return
				default:
					SyncExecutions[sync.Id] = struct {
						startChan chan bool
						stopChan  chan bool
						status    string
					}{startChan, stopChan, "running"}
					// Perform sync operation
					fmt.Println("Executed sync " + sync.SyncName)
					data.SyncData(sync)
					time.Sleep(30 * time.Second) // Simulate work
				}
			}
		}
	}

	// In case the goroutine ends without a stop signal, release the lock.
	if err := mutex.Unlock(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Released lock for ", name)
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
			startChan := make(chan bool)
			stopChan := make(chan bool)
			SyncExecutions[sync.Id] = struct {
				startChan chan bool
				stopChan  chan bool
				status    string
			}{startChan, stopChan, "stopped"}

			SyncsWaitGroup.Add(1)
			go ExecuteSync(sync, startChan, stopChan, &SyncsWaitGroup)
			startChan <- true
		}
		SyncsWaitGroup.Wait()
		time.Sleep(5 * time.Second)
	}
}

func StopSync(syncId string) {
	if goroutine, exists := SyncExecutions[syncId]; exists {
		goroutine.stopChan <- true
		fmt.Printf("Goroutine ID %s stopped\n", syncId)
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

func CheckExecutions() {
	for {
		for key := range SyncExecutions {
			if goroutine, exists := SyncExecutions[key]; exists {
				fmt.Printf("Goroutine ID %s\n", key)
				fmt.Printf("Goroutine status %s\n", goroutine.status)
				fmt.Printf("startChan %s\n", goroutine.startChan)
				fmt.Printf("stopChan %s\n", goroutine.stopChan)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
