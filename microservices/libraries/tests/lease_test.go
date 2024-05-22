package tests

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"microservices/libraries"
	"testing"
	"time"
)

func TestLease(t *testing.T) {
	var sampleManager libraries.LeaseManager
	session := sampleManager.GetLeasedSession()
	res := sampleManager.Leases(session)
	fmt.Println(len(res))

	res2 := sampleManager.Leases(session)
	if len(res2) == 0 {
		t.Fail()
	}
	time.Sleep(10 * time.Second)
	res2 = sampleManager.Leases(session)
	if len(res2) == 0 {
		t.Fail()
	}
	session.Done()
	session.Close()
}

func TestWatchExpiredEventsReceiveExpiry(t *testing.T) {
	keyPrefix := "agents"
	key := keyPrefix + "my-test-key"
	value := "my value that expired"

	cfg := clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}}
	cli, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	go func() {
		lease, err := cli.Grant(context.Background(), 30)
		if err != nil {
			log.Fatal(err)
		}
		_, err = cli.Put(context.Background(), key, value, clientv3.WithLease(lease.ID))
		if err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case ev := <-libraries.WatchExpiredLease(cli, keyPrefix):
		log.Printf("Expiry event for key: '%s' and value: '%s'\n", ev.PrevKv.Key, ev.PrevKv.Value)
	case <-time.After(50 * time.Second):
		t.Fatalf("Expiry event did not fire")
	}
	time.Sleep(300000)
}

func TestRenewLease(t *testing.T) {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
