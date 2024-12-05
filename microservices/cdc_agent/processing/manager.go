package processing

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"microservices/libraries"
	"time"
)

type Manager struct {
	stopChan chan bool
}

type RunFunc func(chan bool)

func (m *Manager) startWorker(f RunFunc) {
	m.stopChan = make(chan bool)

	// Start the worker goroutine
	go f(m.stopChan)
}

func (m *Manager) stopWorker() {
	// Send the stop signal to the worker
	m.stopChan <- true

	// Wait for the worker to terminate
	time.Sleep(5 * time.Second)
}

// ListenBalanceEvents If syncs or agents changed this routine is called
func (m *Manager) ListenBalanceEvents(session *concurrency.Session, keyPrefix string) {
	channel := session.Client().Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for resp := range channel {
		for _, event := range resp.Events {
			fmt.Println(event)
			switch event.Type {
			case mvccpb.PUT:
				m.BalanceAndStop(session)
			case mvccpb.DELETE:
				m.BalanceAndStop(session)
			}
		}
	}
}

func (m *Manager) BalanceAndStop(session *concurrency.Session) {
	BalanceSyncs(session)
	m.stopWorker()
}

// ListenGloballyBalanceEvent Observe if a sync is added or and agent is added or removed
func (m *Manager) ListenGloballyBalanceEvent(session *concurrency.Session, keyPrefix string) {
	m.ObserveBalanceProcess(session)
	go func() {
		m.ListenBalanceEvents(session, keyPrefix)
	}()

	go func() {
		m.ListenBalanceEvents(session, "syncs")
	}()
}

// ObserveBalanceProcess Observe balance-rebalnce syncs process and reactto this event
func (m *Manager) ObserveBalanceProcess(session *concurrency.Session) {
	go func() {
		watcher := clientv3.NewWatcher(session.Client())
		watchChan := watcher.Watch(context.Background(), "/rebalance/", clientv3.WithPrefix())

		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				if event.Type == clientv3.EventTypeDelete {
					m.stopWorker()
				}
			}
		}
	}()
}

// ObserveDiedAgentEvent Observe if an agent die and start a rebalance process
func (m *Manager) ObserveDiedAgentEvent(session *concurrency.Session, keyPrefix string) {
	go func() {
		select {
		case ev := <-libraries.WatchExpiredLease(session.Client(), keyPrefix):
			if ev != nil {
				log.Printf("ObserveDiedAgentEvent'\n")
				log.Printf("Expiry event for key: '%s' and value: '%s'\n", ev.PrevKv.Key, ev.PrevKv.Value)
				if session == nil {
					log.Printf("session nil ObserveDiedAgentEvent'\n")
				}
				BalanceSyncs(session)
				m.stopWorker()
			}
		}
	}()
}
