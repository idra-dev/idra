---
title: Workers, Jargon and Distributed Locks
description: >
  #### A short description about some concepts that are part of Idra.
date: 2017-01-05
weight: 5
categories: [Architecture, Docs]
tags: [workers, providers, concepts]
---
Worker
Each worker node is responsible for processing one or more syncs. A sync is an object that contains a source connector, from which data is retrieved, and a destination connector, where the data is written.
In its simplest configuration, a worker can use a JSON file and be launched without the support of ETCD in single mode. Idra can also be launched in cluster mode (multiple instances are run to increase computing capacity).
The supported connectors at the moment are:

Postgresql

Mysql-Mariadb

Sqlite

Microsoft SQL Server

MongoDB

Apache Kafka

Amazon S3

Custom API

Here are some concepts present in Idra:

##### Sync: Data synchronization process consisting of a source and a destination

##### Connector: Source or destination provider that connects to a database, sensor, middleware, etc.
##### Agent: Instance of Idra responsible for executing syncs and connectors
##### ETCD: Distributed database based on the key-value paradigm.

Each worker, besides being responsible for processing synchronizations, also implements specific algorithms for distributed concurrency. By using leader election,
the system implements the ability to distribute the load and redistribute computation if a worker fails or a new worker is started.
The leader election algorithm, or distributed consensus algorithm, is a mechanism used by distributed systems to select a node within the system to act as a leader.


## Distributed Lock

Each synchronization process is guaranteed to process a single synchronization process and uses a distributed lock to achieve this result. A distributed lock is a mechanism used in distributed systems to coordinate concurrent access to shared resources by multiple nodes. Essentially, a distributed lock functions as a global semaphore that ensures only one entity at a time can access a particular resource.

The idea behind the distributed lock is to use a distributed coordination system, in this case, we use ETCD, to allow nodes to compete for control of the shared resource. This coordination system can be implemented using a variety of techniques, including election algorithms, communication protocols, and other mechanisms.

When a node requests control of a resource, it sends a request to acquire the distributed lock to the distributed coordination system. If the lock is available, the node acquires the lock and can access the shared resource. If the lock is not available, the node waits until it becomes available.

It is important to note that a distributed lock can be implemented in different modes. For example, a distributed lock can be exclusive, meaning that only one node at a time can acquire it, or it can be shared, meaning that multiple nodes can acquire it simultaneously. The choice of the type of distributed lock depends on the specific requirements of the distributed system in which it is used.
