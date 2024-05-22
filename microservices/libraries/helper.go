package libraries

import (
	"context"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type Single struct {
	Value string
}

var singleInstance *Single

var (
	dialTimeout    = 3 * time.Second
	requestTimeout = 10 * time.Second
)

func GetContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	return ctx
}

func GetClient() (*clientv3.Client, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
	}
	etcd := os.Getenv("ETCD")
	if etcd == "" {
		etcd = "127.0.0.1:2379"
	}
	cluster := []string{etcd}

	username := os.Getenv("ETCD_USERNAME")
	password := os.Getenv("ETCD_PASSWORD")
	if username == "" || password == "" {
		return clientv3.New(clientv3.Config{
			DialTimeout:          dialTimeout,
			Endpoints:            cluster,
			DialKeepAliveTimeout: 5 * time.Second,
			DialOptions:          dialOpts,
		})
	} else {
		return clientv3.New(clientv3.Config{
			DialTimeout:          dialTimeout,
			Endpoints:            cluster,
			DialKeepAliveTimeout: 5 * time.Second,
			DialOptions:          dialOpts,
			Username:             username,
			Password:             password,
		})
	}
}

func GetEtcd() (context.Context, *clientv3.Client) {
	ctx := GetContext()
	cli, err := GetClient()
	if err != nil {
		log.Fatalln("Cannot connect to ETCD:", err)
	}
	return ctx, cli
}

func GetOptionsRead() []clientv3.OpOption {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(10000),
	}
	return opts
}

func SaveKey(key string, data []byte) {
	ctx, cli := GetEtcd()
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	_, _ = kv.Put(ctx, key, string(data))
}

func SetIntKey(key string, offset int64) int64 {
	ctx, cli := GetEtcd()
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	_, err := kv.Put(ctx, key, strconv.FormatInt(offset, 10))
	if err != nil {
		return -1
	}
	return offset
}

func DeleteKey(key string) bool {
	ctx, cli := GetEtcd()
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	response, err := kv.Delete(ctx, key)
	if err != nil {
		log.Println(err)
	}
	if response.Deleted == 1 {
		return true
	} else {
		return false
	}
}

func DeleteKeys(path string) bool {
	cli, err := GetClient()
	if err != nil {
		fmt.Println("Failed to get keys:", err)
		return false
	}
	defer cli.Close()
	// Get all keys under the path
	resp, err := cli.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("Failed to get keys:", err)
		return false
	}

	// Delete all keys under the path
	for _, kv := range resp.Kvs {
		_, err = cli.Delete(context.Background(), string(kv.Key))
		if err != nil {
			fmt.Println("Failed to delete key:", err)
			return false
		}
	}
	return true
}

func GetKeys(prefix string) *clientv3.GetResponse {
	ctx, cli := GetEtcd()
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	opts := GetOptionsRead()

	gr, _ := kv.Get(ctx, prefix, opts...)
	return gr
}

func GetKey(key string) *clientv3.GetResponse {
	ctx, cli := GetEtcd()
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	gr, _ := kv.Get(ctx, key)
	return gr
}

func GetMachineId() string {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			singleInstance = &Single{}
			singleInstance.Value = getMachineId()
		}
	}

	return singleInstance.Value
}

func getMachineId() string {
	const DisableMachineId = "DisableMachineId"
	value := os.Getenv(DisableMachineId)
	if value != "true" {
		id, err := machineid.ID()
		if err != nil {
			log.Fatal(err)
		}
		return id
	} else {
		fmt.Print("Use Fake Machine Id")
		return uuid.New().String()
	}
}
