---
title: Overview
description: >
  #### Overview
date: 2024-10-18
categories: [Docs]
tags: [cli]
weight: 1
---

# CLI Overview for Managing IDRA

The IDRA Command Line Interface (CLI) is a powerful tool designed to facilitate the management and interaction with various aspects of the IDRA system directly from the terminal. It offers streamlined, efficient operations for handling Assignments, Agents, Syncs, Users, Offsets, DataSources, and the general status of the application.

## Key Functionalities

### 1. Status

- **Show General Application Status**: Display the overall status of the application, including system health and ETCD health.

### 2. Assignments

- **Show List**: Retrieve and display a list of all assignments.

### 3. Agents

- **Show List**: Retrieve and display a list of all agents.

### 4. Syncs

- **Show List**: Display a comprehensive list of all syncs.
- **View Details**: Fetch and display detailed information about a specific sync.
- **Delete**: Remove a specific sync from the system.
- **Export in JSON Format**: Export all sync data in JSON format for external use or backup.
- **Import JSON Elements**: Import sync data from a JSON file to create new sync entries.

### 5. Users

- **Create User from JSON File**: Create a new user by importing details from a JSON file.
- **Show List**: Retrieve and display a list of all users.
- **View Details**: Fetch and display detailed information about a specific user.
- **Delete**: Remove a specific user from the system.

### 6. Offsets

- **Show List**: Retrieve and display a list of all offsets.

### 7. DataSources

- **Show List**: Retrieve and display a list of all data sources.
- **Create**: Add a new data source to the system.
- **Delete**: Remove a specific data source.
- **Edit**: Modify an existing data source.

---

This CLI enables seamless management of IDRA system components, allowing users to perform essential operations like creating, viewing, editing, deleting, and exporting data with ease. Through a clear and structured command set, administrators can efficiently manage the platform without relying on the UI.
