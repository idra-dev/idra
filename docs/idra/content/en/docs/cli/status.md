---
title: Status
description: >
  #### Visualization Service and ETCD health
date: 2024-10-18
categories: [Docs]
tags: [cli, etcd]
weight: 3
---

# CLI Status Command

The `webcdc-cli status` command allows you to visualize the health of both the Service and ETCD components of the IDRA system. This command provides valuable insights into the operational status of these critical components, helping you maintain system reliability.

```bash
webcdc-cli status
```

## Service Health

When you execute the `webcdc-cli status` command, you can retrieve information about the Service health, which includes:

- **Region**

## ETCD Health

In addition to Service health, the command also provides details about the ETCD component. The ETCD health information includes:

- **ID**
- **Name**
- **PeerURLs**
- **ClientURLs**

## Overall Status

The output of the `webcdc-cli status` command provides an overall status for both the Service and ETCD components. This summary allows you to quickly assess the health of the system at a glance and take necessary actions if any issues are detected.

---

By utilizing the `webcdc-cli status` command, you gain immediate access to critical health metrics for the Service and ETCD components, enabling proactive monitoring and management of your IDRA system.

![](/images/cli/status.png)
