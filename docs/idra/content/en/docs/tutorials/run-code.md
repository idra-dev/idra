---
title: Run Code
description: >
  Tutorial to run code
date: 2017-01-05
weight: 4
---

### Fix etcd

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