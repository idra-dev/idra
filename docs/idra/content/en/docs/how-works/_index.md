---
title: How Work
description: How works a sync
weight: 1
---

A sync is an entity that describes all informations that Idra needs to process data from a source to a destination.

A typical record that describe a sync has this form:

```json
{"id":"28955df5-3f4a-48f5-b60e-cf898d01cbe6","sync_name":"Sincro_postgres_mysql",
"source_connector":
    {"id":"id","connector_name":"con1","connector_source_type":"PostgresGORM","connection_string":"host=localhost user=postgres password=postgres dbname=school port=5432 sslmode=disable TimeZone=Europe/Rome","database":"","query":"","table":"bimbi","polling_time":60,"timestamp_field":"ts","timestamp_field_format":"","max_record_batch_size":50000,"save_mode":"Insert","start_offset":0,"token":""},
"destination_connector":
    {"id":"id","connector_name":"con2","connector_source_type":"MysqlGORM","connection_string":"root:cacata12@tcp(127.0.0.1:3306)/school?charset=utf8mb4\u0026parseTime=True\u0026loc=Local","database":"school","query":"","table":"bimbi","polling_time":60,"timestamp_field":"ts","timestamp_field_format":"","max_record_batch_size":50000,"save_mode":"Insert","start_offset":0,"token":""},
"mode":"LastDestinationTimestamp","disabled":false}
```

### Save Mode (mode parameter):
This parameter describes the strategy used to update data on destination.

#### FullWithId
This mode uses an identifier key to track the last record inserted. It employs an insert-only approach, meaning it does not update existing records but adds new ones sequentially. The offset, which is the last used ID, is stored in the Idra database to ensure that new records are inserted correctly and in order.
#### LastDestinationId
This strategy checks the last registered ID in the destination system, retrieves it, and transfers all new records to the data source. No separate offset needs to be maintained in the Idra database, as the most recent ID is directly consulted to determine which records need to be transferred. This approach is useful for ensuring data completeness without manually managing offset information.
#### LastDestinationTimestamp
This strategy is similar to LastDestinationId but uses a timestamp field to identify recent records. In this case, you can choose between an insert strategy or a merge strategy. Utilizing a timestamp allows for more flexible data management by synchronizing records based on their creation or update time.
#### Timestamp
This mode is similar to FullWithId but uses a timestamp field to track records. As with the previous strategy, you can choose whether to use an append strategy or a merge strategy. Using a timestamp provides finer control over inserted and updated records, making it easier to manage changes over time.

A connector, on the other hand, is a data reference. For more details, please refer to the dedicated page in the documentation.


