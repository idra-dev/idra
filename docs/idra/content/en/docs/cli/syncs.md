---
title: Syncs
description: >
  #### Commands related to syncs
date: 2024-10-18
categories: [Docs]
tags: [cli, syncs]
weight: 4
---

# CLI Syncs Commands

The `webcdc-cli syncs` commands provide a comprehensive set of functionalities for managing syncs related to the current user. You can [ list, import, export, view details, and delete syncs ] with ease.

## Available Commands

### 1. List Syncs

To retrieve and display all syncs associated with the current user, use the following command:

```bash
webcdc-cli syncs list
```

This command will get a list of all syncs.

![](/images/cli/syncs_list.png)

### 2. Import Syncs

To import sync data from a JSON file, use the following command:

```bash
webcdc-cli syncs import <filename.json>
```

File format need to be:

```json
[
	{
	"id": "",
    "mode": string with posible values ["FullWithId", "LastDestinationId", "LastDestinationTimestamp", "Timestamp"],
    "sync_name": string,
    "source_connector": {
      "id": "",
      "query": string,
      "table": string,
      "polling_time": number,
      "connector_name": string,
      "timestamp_field": string,
      "connection_string": string,
      "connector_source_type": string with posible values ["PostgresGORM", "MysqlGORM", "MssqlGORM", "KafkaConnector", "MongodbManager", "S3"],
      "max_record_batch_size": number,
      "timestamp_field_format": string,
	  "attributes": {
	 	"anyKey": string
		...
	  }
    },
    "destination_connector": {
      "id": "",
      "table": string,
      "save_mode": string with posible values ["Insert", "Upsert"],
      "connector_name": string,
      "timestamp_field": string,
      "connection_string": string,
      "attributes": {
        "anyKey": string
		...
      },
      "connector_source_type": string with posible values ["PostgresGORM", "MysqlGORM", "MssqlGORM", "KafkaConnector", "MongodbManager", "S3"],
      "max_record_batch_size": number,
      "timestamp_field_format": string
    },

	...
]
```

You can add more than one sync in a single import, which simplifies the process of creating syncs.

#### syncs.json

```json
[
  {
    "id": "",
    "mode": "LastDestinationId",
    "sync_name": "TestImportCLI",
    "source_connector": {
      "id": "TestCLI",
      "query": "Test",
      "table": "Test",
      "save_mode": "Upsert",
      "polling_time": 60,
      "connector_name": "Test",
      "timestamp_field": "Test",
      "connection_string": "ConnStringCLI",
      "connector_source_type": "MysqlGORM",
      "max_record_batch_size": 50000,
      "timestamp_field_format": "Test"
    },
    "destination_connector": {
      "id": "TestCLI",
      "query": "Test",
      "table": "Test",
      "save_mode": "Upsert",
      "polling_time": 60,
      "connector_name": "Test",
      "timestamp_field": "Test",
      "connection_string": "ConnStringCLI",
      "attributes": {
        "s": "s"
      },
      "connector_source_type": "MysqlGORM",
      "max_record_batch_size": 50000,
      "timestamp_field_format": "Test"
    }
  }
]
```

![](/images/cli/syncs_import.png)

### 3. View Sync

To view detailed information about a sync, use the following command:

```bash
webcdc-cli syncs view <ID>
```

Replace <ID> with the identifier of the sync you wish to view.

![](/images/cli/syncs_view.png)

This command will display all relevant details associated with the specified sync.

### 4. Export Syncs

To export all sync data in JSON format, use the following command:

```bash
webcdc-cli syncs export <Optional: Filename.json -> default: backup_currentTime.json>
```

You can specify a filename by replacing <Optional: Filename.json> with your desired filename. If you do not provide a filename, the default export filename will be backup_currentTime.json, where currentTime represents the current timestamp.

![](/images/cli/syncs_export.png)

### 5. Delete Sync

To remove a specific sync from the system, use the following command:

```bash
webcdc-cli syncs delete <ID>
```

Replace <ID> with the identifier of the sync you want to delete. This command will permanently remove the specified sync from the system.

![](/images/cli/syncs_delete.png)
