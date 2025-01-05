package processing

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"microservices/libraries"
	"microservices/libraries/etcd"
	"microservices/libraries/models"
	"time"
)

func BalanceSyncs(session *concurrency.Session) {
	fmt.Println("Started balance syncs")
	agents := etcd.GetAgents()
	e := concurrency.NewElection(session, "/rebalance/")
	ctx := context.Background()
	// Elect a leader (or wait that the leader resign)
	if err := e.Campaign(ctx, "e"); err != nil {
		return
	}
	fmt.Println("leader election for me")
	leader, errLeader := e.Leader(ctx)
	if errLeader != nil {
		return
	}
	//Check if really we are leader
	if e.Key() == string(leader.Kvs[0].Key) {
		syncs := etcd.GetSyncs()
		lb := CreateLoadBalancer(agents)
		assignments := lb.AddTasks(syncs)
		//Delete all assignments
		executed := libraries.DeleteKeys(models.AssignmentsPath)
		if executed {
			for k := range assignments {
				fmt.Printf("key[%s] value[%s]\n", k, assignments[k])
				libraries.SaveKey(models.AssignmentsPath+assignments[k]+"/"+k, []byte(assignments[k]))
			}
		}
		fmt.Println("Load balancing executed")
	}

	if err := e.Resign(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Resign ")
}

func AllocateSyncs(session *concurrency.Session) {
	fmt.Println("Started balance syncs")
	agents := etcd.GetAgents()
	syncs := etcd.GetSyncs()
	lb := CreateLoadBalancer(agents)
	assignments := lb.AddTasks(syncs)
	//Delete all assignments
	executed := libraries.DeleteKeys(models.AssignmentsPath)
	if executed {
		for k := range assignments {
			fmt.Printf("key[%s] value[%s]\n", k, assignments[k])
			libraries.SaveKey(models.AssignmentsPath+assignments[k]+"/"+k, []byte(assignments[k]))
		}
	}
	fmt.Println("Load balancing executed")
	fmt.Println("Resign ")
}

func ClusterNodeLeader() bool {
	client, err := libraries.GetClient()
	if err != nil {
		return false
	}
	electionSession, err := concurrency.NewSession(client, concurrency.WithTTL(1))
	if err != nil {
		log.Fatal(err)
	}
	defer electionSession.Close() // cleanup

	election := concurrency.NewElection(electionSession, "/election-prefix")
	ctx := context.Background()

	fmt.Println("Attempting to become a leader")
	// start leader election
	if err := election.Campaign(ctx, "value"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Became a leader!")
	go func() {
		for {
			if err := election.(); err != nil {
				log.Printf("Errore durante KeepAlive: %v", err)
			}
			time.Sleep(5 * time.Second) // Mantiene viva l'elezione ogni 5 secondi
		}
	}()

	return true
}
