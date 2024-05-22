package tests

import (
	"context"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	var name = flag.String("name", "", "give a name")
	flag.Parse()
	// Create a etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// create a sessions to aqcuire a lock
	s, _ := concurrency.NewSession(cli)
	defer s.Close()
	l := concurrency.NewMutex(s, "/distributed-lock/")
	ctx := context.Background()
	// acquire lock (or wait to have it)
	if err := l.Lock(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for ", *name)
	fmt.Println("Do some work in", *name)
	time.Sleep(5 * time.Second)
	if err := l.Unlock(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for ", *name)
}

func TestLeader(t *testing.T) {
	var name = flag.String("name", "", "give a name")
	flag.Parse()
	// Create a etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// create a sessions to elect a Leader
	s, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	e := concurrency.NewElection(s, "/leader-election/")
	ctx := context.Background()
	// Elect a leader (or wait that the leader resign)
	if err := e.Campaign(ctx, "e"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("leader election for ", *name)
	fmt.Println("Do some work in", *name)
	time.Sleep(5 * time.Second)
	if err := e.Resign(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("resign ", *name)
}
