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
	for {
		// Creazione di un contesto per timeout (se necessario) e watch su /syncs
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		channel := session.Client().Watch(ctx, "/syncs", clientv3.WithPrefix())

		for resp := range channel {
			if resp.Err() != nil {
				fmt.Printf("Error watching /syncs: %v\n", resp.Err())
				fmt.Println("Reconnecting to /syncs...")
				time.Sleep(2 * time.Second)
				break
			}

			for _, event := range resp.Events {
				fmt.Println("Event received:", event)

				switch event.Type {
				case mvccpb.PUT:
					value := event.Kv.Value
					var sync cdc_shared.Sync
					err := json.Unmarshal(value, &sync)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(event.Kv.Key)
					executionValue, exists := SyncExecutions[sync.Id]
					if exists && executionValue.cancel != nil {
						executionValue.cancel()
					}

					time.Sleep(2 * time.Second)
					if !sync.Disabled {
						go ExecuteSync(sync)
					}
					break
				case mvccpb.DELETE:
					fmt.Println(event.Kv.Key)
					syncId := string(event.Kv.Key)[7:]
					executionValue, exists := SyncExecutions[syncId]
					if exists && executionValue.cancel != nil {
						executionValue.cancel()
					}
					delete(SyncExecutions, syncId)
					time.Sleep(2 * time.Second)
					break
				}
			}
			fmt.Println("Exit inner for...")
		}

		fmt.Println("Watch channel closed, restarting watch...")
		time.Sleep(1 * time.Second) // Aggiungi un ritardo prima di riavviare il watch
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
