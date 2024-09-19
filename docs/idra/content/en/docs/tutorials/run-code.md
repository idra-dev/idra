---
title: Run Code
description: >
  Tutorial to run code
date: 2017-01-05
weight: 4
---

### Install etcd

https://etcd.io/docs/v3.5/install/

### Run Code

go get .

In folder cdc_agent there is main where you can decide if use ETCD for manage scalable configuration. If you need something simple you can simply pass a json filte with connectors. 

```go
func main() {
	err := godotenv.Load("../cdc_agent/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	if custom_errors.IsStaticRunMode() {
		processing.ProcessStatic()
	}else{
		processing.StartWorkerNode()
	}
}
...
func IsStaticRunMode() bool {
	value := os.Getenv(models.Static)
	return value == "true"
}
...
unc ProcessStatic() {
	var syncs []cdc_shared.Sync
	path, err := os.Getwd()
	custom_errors.LogAndDie(err)
	staticFilePath := os.Getenv(models.StaticFilePath)
	dat, err2 := os.ReadFile(staticFilePath)
	custom_errors.LogAndDie(err2)

	json.Unmarshal(dat, &syncs)
	for {
		for _, sync := range syncs {
			fmt.Println(sync.SyncName)
			data.SyncData(sync, sync.Mode)
		}
		time.Sleep(1 * 3600 * time.Second)
	}
}
```
