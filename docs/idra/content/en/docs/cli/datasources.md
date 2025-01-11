---
title: Data Sources
description: >
  #### Commands related to data sources
date: 2024-10-18
categories: [Docs]
tags: [cli, datasources]
weight: 8
---

# CLI Data sources Commands

The `webcdc-cli datasources` commands provide a comprehensive set of functionalities for managing datasources related to the current user. You can [ list, create, edit and delete datasource ] with ease.

## Available Commands

### 1. List Data Sources

To retrieve and display all data sources associated with the current user, use the following command:

```bash
webcdc-cli datasources list
```

This command will get a list of all datasources.

![](/images/cli/datasources_list.png)

### 2. Create Data source

To create data source, use the following command:

```bash
webcdc-cli datasources create <ConnectionString> <DataSourceName>
```

Replace <ConnectionString> with the needed connection string and <DataSourceName> with name for data source you want to create.

![](/images/cli/datasources_create.png)

### 3. Edit Data source

To edit data source, use the following command:

```bash
webcdc-cli datasources edit <ConnectionString> <DataSourceName> <DataSourceID>
```

Replace <ConnectionString> with the needed connection string and <DataSourceName> with name and <DataSourceID> with identifier for data source you want to edit.

![](/images/cli/datasources_edit.png)

### 4. Delete Data source

To delete data source, use the following command:

```bash
webcdc-cli datasources delete <DataSourceID>
```

Replace <DataSourceID> with identifier for data source you want to delete.

![](/images/cli/datasources_delete.png)
