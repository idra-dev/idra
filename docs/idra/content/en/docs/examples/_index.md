---
title: Setup
weight: 3
date: 2017-01-05
description: Instructions
---


### Docker

* Cdc app

See deploy script

docker build -t agent-cdc -f ./cdc_agent/Dockerfile .


* Web app

docker build -t web-cdc .

docker run -p 8080:8080/tcp web-cdc


* Command to deploy etcd in kubernetes

helm install my-release bitnami/etcd --set auth.rbac.create=false

To connect to your etcd server from outside the cluster execute the following commands:

    kubectl port-forward --namespace default svc/my-release-etcd 2379:2379 &
    echo "etcd URL: http://127.0.0.1:2379"

* Force new deploy on helm

helm upgrade web-cdc chart-web-cdc

* Run etcd

./bin/etcd

* Rebalance process

Use Leases to monitor cluster nodes
In case a node is added or removed we elect a leader that will manage rebalance process assigning
syncs to cluster nodes


* Run

go run main.go

### Environment Variables

* Static

This environment variable is used to run agent in static mode using a JSON file or run it using ETCD

* ETCD

Contains value for ETCD database server url

* Domain

Contains value for domain to assign to GinSwager to expose Rest API using Kubernetes deployment

var urlPath = os.Getenv("DOMAIN")
var url func(config *ginSwagger.Config)
if urlPath == "" {
	url = ginSwagger.URL("http://0.0.0.0:8080/swagger/doc.json")
} else {
	url = ginSwagger.URL(urlPath + "/swagger/doc.json")
}
