---
title: Change Data Capture
description: What is a Change Data capture
weight: 1
---

This is an introduction to Change Data Capture

## What is it?

Change Data Capture (CDC) is a set of techniques or tools used to identify and capture changes made to data in a database or data source. This is often used in data integration, data warehousing, and real-time analytics to ensure that changes in the source system are reflected in the target system or application.

## How works?

* **Monitoring**: CDC involves monitoring the data source for changes. This can be done in various ways depending on the database system and the CDC method used.

* **Capture**: Once changes are detected, CDC captures these changes. This can include new records, updates to existing records, and deletions.

* **Transformation**: Sometimes, captured changes need to be transformed or formatted to match the requirements of the target system.

* **Loading**: The transformed changes are then loaded into the target system, which could be another database, a data warehouse, or an analytics platform.

* **Application**: The target system applies these changes to keep its data synchronized with the source.

## Benefits
Real-Time Data Synchronization: Keeps the target systems up-to-date with minimal delay, which is critical for real-time analytics and decision-making.

Efficient Data Integration: Reduces the need for full data extracts and transfers, which can be resource-intensive.

Improved Data Accuracy: Ensures that changes are consistently and accurately reflected across systems.

Reduced Latency: Helps in minimizing the lag between when a change occurs and when it is reflected in the target system.

## Use Cases

* **Data Warehousing**: To keep the data warehouse synchronized with operational databases.

* **Real-Time Analytics**: To ensure analytics platforms reflect the latest data.

* **ETL Processes**: For efficient Extract, Transform, Load (ETL) processes by capturing only changed data.

* **Data Replication**: To replicate changes from one database to another.

* **Microservices**: CDC is well known patern in Microservices Oriented Architectures

* [Getting Started](/docs/getting-started/): Get started with Idra


