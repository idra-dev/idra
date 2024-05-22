---
title: Getting Started
description: What does your user need to know to try your project?
categories: [Setup, Placeholders]
tags: [cdc,etcd,rest]
weight: 2
---


## Local Setup

### Prerequisites

Setup ETCD

* [ETCD](https://etcd.io/docs/): Installing ETCD

* [Install dependencies](https://go.dev/doc/modules/managing-dependencies): go build or go mod download 

### Fix GRPC version error:

go get github.com/coreos/etcd/clientv3

Add to mod file:

google.golang.org/grpc v1.26.0

go mod download google.golang.org/grpc


first: open the go.mod, add this line :

replace ( google.golang.org/grpc => google.golang.org/grpc v1.26.0)

then:

go get -u -v go.etcd.io/etcd

go mod download google.golang.org/grpc                                                                                                                                                                   

go mod tidy

go get google.golang.org/grpc@v1.26.0

### Run CDC and Web REST API

* Run CDC app from code:

Run main.go in cdc_agent folder using "go run main.go"

* Run Web app from code:

Run main.go in web folder using "go run main.go"

## Docker

* ...

