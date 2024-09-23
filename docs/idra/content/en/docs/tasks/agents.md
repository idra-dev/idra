---
title: Agents
description: >
  #### About Agents
date: 2024-01-05
categories: [Architecture, Docs]
tags: [cdc, etcd, clustering]
weight: 2
---

An agent is simply a running instance of Idra.
Idra is designed to run in cluster mode, which enhances its ability to scale effectively. In this mode, all agents within the system connect to a shared ETCD instance. ETCD serves as a distributed key-value store that helps manage configuration data and state across multiple instances of the application.

By having all agents share the same ETCD instance, Idra ensures that they can communicate and coordinate their activities seamlessly. This shared architecture allows the system to scale horizontally, meaning that you can add more agents to handle increased loads without sacrificing performance.

Moreover, using a centralized ETCD instance is crucial for implementing locks that prevent concurrent processing on the same data sources. When multiple agents attempt to access the same resource simultaneously, it can lead to data inconsistencies and processing errors. The locking mechanism provided by ETCD ensures that only one agent can process a given data source at any time. This prevents conflicts and guarantees that the integrity of the data is maintained throughout the processing cycle. Overall, this architecture not only enhances scalability but also improves the reliability and efficiency of data handling within the system.