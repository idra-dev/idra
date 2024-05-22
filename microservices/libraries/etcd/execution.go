package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"microservices/libraries"
	"microservices/libraries/models"
	"path"
)

func GetAgents() []models.CdcAgent {
	var items []models.CdcAgent
	cli, err := libraries.GetClient()
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	leases := GetLeases()

	for _, lease := range leases {
		prefix := fmt.Sprintf("agents/%d", lease.ID)

		gr, err := cli.Get(context.Background(), prefix, clientv3.WithLease(lease.ID), clientv3.WithPrefix())
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range gr.Kvs {
			agent := models.CdcAgent{}
			json.Unmarshal(item.Value, &agent)
			items = append(items, agent)
		}
	}
	return items
}

func GetLeases() []clientv3.LeaseStatus {
	leasesId := []clientv3.LeaseStatus{}
	cli, err := libraries.GetClient()
	if err != nil {
		fmt.Println("Failed to get keys:", err)
		log.Fatal(err)
	}
	leaseResp, err := cli.Leases(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Stampa i lease ID
	for _, leaseID := range leaseResp.Leases {
		leasesId = append(leasesId, leaseID)
	}
	return leasesId
}

func GetSyncs() []cdc_shared.Sync {
	prefix := "/syncs"
	gr := libraries.GetKeys(prefix)
	var items []cdc_shared.Sync
	if gr != nil {
		for _, item := range gr.Kvs {
			syncItem := cdc_shared.Sync{}
			json.Unmarshal(item.Value, &syncItem)
			items = append(items, syncItem)
		}
	}
	return items
}

func GetAssignments() map[string]string {
	gr := libraries.GetKeys(models.AssignmentsPath)
	m := make(map[string]string)
	if gr != nil {
		for _, item := range gr.Kvs {
			m[path.Base(string(item.Key))] = string(item.Value)
		}
	}
	return m
}

func GetSync(key string) (cdc_shared.Sync, error) {
	prefix := "/syncs"
	gr := libraries.GetKey(prefix + "/" + key)
	if len(gr.Kvs) > 0 {
		syncItem := cdc_shared.Sync{}
		json.Unmarshal(gr.Kvs[0].Value, &syncItem)
		return syncItem, nil
	} else {
		return cdc_shared.Sync{}, errors.New("Cannot find key.")
	}
}
