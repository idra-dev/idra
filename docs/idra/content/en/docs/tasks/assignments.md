---
title: Assignments
description: >
  #### Assignments and agents
date: 2024-01-05
categories: [Architecture, Docs]
tags: [cdc, agent]
weight: 2
---

## Assignment
An assignment is an association between a sync and an agent that is processsing that sync.

![](/images/assignments.png)

#### How syncs are processed
In a cluster of agents with more than one member, synchronization work is balanced among all elements within the cluster. Each synchronization task is handled by a single agent at a time. If an agent is added to the cluster or if an agent crashes, a rebalancing process is triggered, redistributing all assigned synchronization tasks.

When an agent crashes, the synchronization tasks assigned to the failed agent are reassigned to other agents within the cluster. This mechanism is somewhat similar to Kafka's rebalancing process, which relies on Zookeeper. Also it is similar to shard management in some databases.

![](/images/shards.png)