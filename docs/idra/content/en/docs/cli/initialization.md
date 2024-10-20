---
title: Initialization
description: >
  #### Initialization of CLI
date: 2024-10-18
categories: [Docs, Setup]
tags: [cli]
weight: 2
---

# Initialization of the CLI

To ensure the proper functioning of the IDRA Command Line Interface (CLI), you need to perform the following initialization steps.

## Step 1: Initialize the CLI with API URL

Run the following command to initialize the CLI with your API URL:

```bash
webcdc-cli init http://your.api.url:port/
```

This command will store the provided API URL in temporary files. The CLI will use this URL for all future commands, allowing you to interact with the IDRA system seamlessly.

## Step 2: Authenticate with the IDRA System

After initializing the CLI, you need to authenticate yourself in the IDRA system. Use the following command to log in:

```bash
webcdc-cli login username password
```

Replace username and password with your actual credentials. This command will authenticate your session and store user information locally in the same temporary file where the API URL is stored.

## After Initialization

Once you have completed these steps, you can use all the functionalities of the CLI to manage Status, Assignments, Agents, Syncs, Users, Offsets and DataSources. The CLI will reference the stored API URL and user information for all subsequent commands.
By following these initialization steps, you will be set up to efficiently manage your IDRA system through the CLI.
