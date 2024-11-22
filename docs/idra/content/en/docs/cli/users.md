---
title: Users
description: >
  #### Commands related to users
date: 2024-10-18
categories: [Docs]
tags: [cli, users]
weight: 4
---

# CLI Users Commands

The `webcdc-cli users` commands provide a comprehensive set of functionalities for managing users. You can [ list, import, view details, and delete ] users.

## Available Commands

### 1. List Users

To retrieve and display all users, use the following command:

```bash
webcdc-cli users list
```

This command will get a list of all users.

![](/images/cli/users_list.png)

### 2. Import Users

To import users data from a JSON file, use the following command:

```bash
webcdc-cli users import <filename.json>
```

File format need to be:

```json
[
  {
    "name": string,
    "email": valid Email string,
    "notes": string,
    "password": string,
    "username": string
  },
  ...
]
```

You can add more than one user in a single import

#### user.json

```json
[
  {
    "name": "CLI",
    "email": "cliImport@gmail.com",
    "notes": "cliImportNotes",
    "password": "secretPassword",
    "username": "ImportUser"
  }
]
```

![](/images/cli/users_import.png)

### 3. View User

To view detailed information about a sync, use the following command:

```bash
webcdc-cli users view <Username>
```

Replace <Username> with the User Name of the user you wish to view.

![](/images/cli/users_details.png)

This command will display all relevant details associated with the specified user.

### 4. Delete Sync

To remove a specific user from the system, use the following command:

```bash
webcdc-cli users delete <Username>
```

Replace <Username> with the User Name of the user you want to delete. This command will permanently remove the specified user from the system.

![](/images/cli/users_delete.png)
