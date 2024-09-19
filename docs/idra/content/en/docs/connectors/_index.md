---
title: Connectors
weight: 4
description: >
  Data Connectors
categories: [Architecture, Docs]
tags: [concepts]
---

## A connector is a reference to a data source or destination used by Idra to trasfer data from a source to a destination. Idra supports multiple connector types that we describe here.

We have some DBMS based connectors and some other. Every connector is based on an interface. So eventually to add a new connector we need to implement just this interface.

#### GORM Postgres
Postgres connector based on GORM.
Connection String sample: 
host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
### GORM Mysql
Mysql connector based on GORM
Connection String Sample: 
user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local


### GORM SQL Server
SQL Server connector based on GORM
Connection String Sample: 
sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm

### SQLite GORM
SQLite connector based on GORM

### REST Connector
REST Connector that uses GET request for read data and POST to push data.
URL Sample: https://jsonplaceholder.typicode.com/posts

### S3 JSON Connector
Connector that send data to AWS S3 Bucket using JSON format
```go
os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET")
...
manager2 := data2.S3JsonConnector{}
manager2.ConnectionString = "eu-west-1"
manager2.SaveData("eu-west-1", "samplengt1", rows, "Account_Charges")
```

#### MongoDB
MongoDB connector based on Mongo Stream technology
Connection String sample: 
"mongodb+srv://username:password@cluster1.rl9dgsm.mongodb.net/?retryWrites=true&w=majority&appName=Cluster1"

#### Kafka (Under development)
Kafka connector
Connection String sample: 
ConnectionString="127.0.0.1:9092"
Attributes:
Username="kafkauser"
Password="kafkapassword"
ConsumerGroup="consumer_group_id"
ClientName="myapp"
Offset="0"
Acks="0|1|-1"

#### Immudb
Connection String sample: 
ConnectionString="127.0.0.1"
Attributes:
username=user
password=password
database=mydb"

#### ChromaDB (Under development, Insert Only)
Connection String sample: 
http://localhost:8000

Table = "products"
IdField = "id"






