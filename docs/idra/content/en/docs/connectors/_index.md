---
title: Connectors
weight: 4
description: >
  Data Connectors
categories: [Architecture, Docs]
tags: [concepts]
---

### GORM Postgres
Idra allows the integration of applications, data, sensors, and messaging platforms using connectors that enable interaction with all the data sources we want. The solution is based on the concept of a connector, which is nothing more than a pluggable library (implementing a specific interface) in Golang.
### Monitoring
Monitoring is done using a dedicated web application that allows you to see the connectors that are running and to conveniently define new ones using a user-friendly interface.


### Architecture
The architecture of Idra is based on a series of workers responsible for executing various synchronization processes. 
Alternatively, Idra can also be run in single-worker mode. 

Additionally, it's possible to run it without relying on a supporting database simply by using a configuration file. 
Conceptually, it moves data from a source to a destination. The source can be a sensor exposing data, a Cloud API, a database, a message middleware such as Apache Kafka, any Storage (currently supporting S3), while the destination can be another system of the same type (database, message middleware, sensor, API, etc.).