package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"microservices/libraries"
	"time"
)

type Manager struct {
}

type RunFunc func()

func (m *Manager) startWorker(f RunFunc) {
	go f()
}

// ListenSyncEvents If syncs changed this routine is called
func (m *Manager) ListenSyncEvents(session *concurrency.Session) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		channel := session.Client().Watch(ctx, "/syncs", clientv3.WithPrefix())

		for resp := range channel {
			if resp.Err() != nil {
				fmt.Printf("Error watching /syncs: %v\n", resp.Err())
				fmt.Println("Reconnecting to /syncs...")
				time.Sleep(3 * time.Second)
				break
			}
			locking := Locking{}
			for _, event := range resp.Events {
				fmt.Println("Event received:", event)

				switch event.Type {
				case mvccpb.PUT:
					syncId := string(event.Kv.Key)[7:]
					var sync cdc_shared.Sync
					lockKey := "/lock/" + syncId
					err := json.Unmarshal(event.Kv.Value, &sync)
					if err != nil {
						log.Fatal(err)
					}
					executionValue, exists := SyncExecutions[sync.Id]
					if exists && executionValue.cancel != nil {
						executionValue.cancel()
						client, err := libraries.GetClient()
						defer client.Close()
						if err != nil {
							log.Fatal("")
						}
						if !sync.Disabled {
							err := locking.AcquireLock(context.Background(), client, lockKey)
							if err == nil {
								go ExecuteSync(sync)
							} else {
								isOwner := locking.VerifyOwnerLock(lockKey)
								if isOwner {
									go ExecuteSync(sync)
								}
							}
						} else {
							locking.ReleaseLock(lockKey)
						}
					}
					break
				case mvccpb.DELETE:
					syncId := string(event.Kv.Key)[7:]
					executionValue, exists := SyncExecutions[syncId]
					if exists && executionValue.cancel != nil {
						executionValue.cancel()
						delete(SyncExecutions, syncId)
						locking.ReleaseLock(string(event.Kv.Key))
					}
					break
				}
			}
			fmt.Println("Exit inner for...")
		}

		fmt.Println("Watch channel closed, restarting watch...")
		time.Sleep(1 * time.Second)
	}
}

// ListenGloballyBalanceEvent Observe if a sync is added or and agent is added or removed
func (m *Manager) ListenGloballyBalanceEvent(session *concurrency.Session, keyPrefix string) {
	go func() {
		lm := libraries.LeaseManager{}
		sessionSyncEvents := lm.GetLeasedSession()
		m.ListenSyncEvents(sessionSyncEvents)
	}()
}
