---
title: ETCD
description: >
  #### About ETCD
date: 2017-01-05
categories: [Architecture, Docs]
tags: [cdc, etcd]
weight: 2
---

ETCD plays an important role in the application. ETCD is a highly reliable distributed key-value database designed to be used as a coordination data store for distributed applications. 

Here are some of its key features:

Distributed architecture: 

ETCD is designed to operate in a distributed environment and to be able to scale horizontally. It can run on a cluster of machines working together to provide a reliable service.

Distributed consensus:

ETCD uses a distributed consensus algorithm to ensure that all machines within the cluster have a consistent copy of the data. This distributed consensus algorithm is called Raft.

RESTful API: 

ETCD provides a RESTful API that allows applications to access the data stored in it easily and conveniently. ETCD's RESTful API is designed to be simple and intuitive to use.

Data consistency: 

ETCD ensures that data is always consistent and correct. This means that all changes made to the data are quickly and reliably propagated to all machines within the cluster.

Security: 

ETCD provides a range of security mechanisms to protect the data. This includes authentication and authorization, encryption, and key management.

Open source: 

ETCD is an open-source project that is available for free use and modification. This means that developers can contribute to the code and improve it to meet their specific needs.

ETCD is used in this application to ensure that one and only one agent performs data synchronization. It allows for the election of a leader who is responsible for rebalancing the work of the agents when a new agent is added and is no longer available (due to deletion or crash), and when something is changed at the sync level such as the addition or removal of a sync.

The agent is written in Golang to simplify the process of managing the code that handles concurrency. In fact, it makes heavy use of Goroutines, which simplify the writing and management of concurrency.

Code can also use syntax highlighting.

```go
func main() {
  input := `var foo = "bar";`

  lexer := lexers.Get("javascript")
  iterator, _ := lexer.Tokenise(nil, input)
  style := styles.Get("github")
  formatter := html.New(html.WithLineNumbers())

  var buff bytes.Buffer
  formatter.Format(&buff, style, iterator)

  fmt.Println(buff.String())
}
```


