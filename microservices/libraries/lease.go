package libraries

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
)

type LeaseManager struct {
}

func (manager LeaseManager) GetLeasedSession() *concurrency.Session {
	_, cli := GetEtcd()
	session, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal("Error to create session...")
	}
	return session
}

func (manager LeaseManager) Leases(session *concurrency.Session) []clientv3.LeaseStatus {
	leases, err := session.Client().Leases(context.Background())
	if err != nil {
		log.Panic("Error during reading leases")
	}
	for _, lease := range leases.Leases {
		fmt.Printf("ID: %d\n", lease.ID)
	}
	return leases.Leases
}

func WatchExpiredLease(cli *clientv3.Client, keyPrefix string) <-chan *clientv3.Event {
	events := make(chan *clientv3.Event)

	wCh := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix(), clientv3.WithPrevKV(), clientv3.WithFilterPut())

	go func() {
		defer close(events)
		for wResp := range wCh {
			for _, ev := range wResp.Events {
				expired, err := isExpired(cli, ev)
				if err != nil {
					log.Println("Error when checking expiry")
				} else if expired {
					events <- ev
				}
			}
		}
	}()

	return events
}

func isExpired(cli *clientv3.Client, ev *clientv3.Event) (bool, error) {
	if ev.PrevKv == nil {
		return false, nil
	}

	leaseID := clientv3.LeaseID(ev.PrevKv.Lease)
	if leaseID == clientv3.NoLease {
		return false, nil
	}

	ttlResponse, err := cli.TimeToLive(context.Background(), leaseID)
	if err != nil {
		return false, err
	}

	return ttlResponse.TTL == -1, nil
}
