package processing

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"microservices/cdc_agent/cobra"
	"microservices/libraries"
	"strconv"
)

type Locking struct{}

func (Locking) AcquireLock(ctx context.Context, client *clientv3.Client, lockKey string) error {
	txn := client.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)). // If key doesn't exist
		Then(clientv3.OpPut(lockKey, strconv.Itoa(cobra.CusterNodeId))).
		Else(clientv3.OpGet(lockKey))

	txnResp, err := txn.Commit()
	if err != nil {
		return fmt.Errorf("Error during transaction: %v", err)
	}

	if !txnResp.Succeeded {
		return fmt.Errorf("Already locked!")
	}

	fmt.Println("Lock acquired!")
	return nil
}

func (Locking) ReleaseLock(lockKey string) error {
	client, err := libraries.GetClient()
	defer client.Close()
	if err != nil {
		log.Fatal("")
	}
	_, err2 := client.Delete(context.Background(), lockKey)
	if err2 != nil {
		return fmt.Errorf("Error during lock release: %v", err)
	}

	fmt.Println("Lock released!")
	return nil
}

func (Locking) GetValue(ctx context.Context, client *clientv3.Client, lockKey string) string {
	resp, err := client.Get(ctx, lockKey)
	if err != nil {
		fmt.Println("Error during lock release: %v" + err.Error())
		return "error"
	}
	for _, kv := range resp.Kvs {
		return string(kv.Value)
	}

	return ""
}

func (Locking) VerifyOwnerLock(lockKey string) bool {
	locking := Locking{}
	client, err := libraries.GetClient()
	defer client.Close()
	if err != nil {
		log.Fatal("")
	}
	res := locking.GetValue(context.Background(), client, lockKey)
	nodeId, _ := strconv.Atoi(res)
	if nodeId == cobra.CusterNodeId {
		locking.AcquireLock(context.Background(), client, lockKey)
		return true
	}
	return false
}
