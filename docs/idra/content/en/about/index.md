---
title: About Idra
linkTitle: About
menu: {main: {weight: 20}}
---

{{% blocks/cover title="About Idra" image_anchor="bottom" height="auto" %}}
A Golang based Change Data Capture platform
{.mt-5}

{{% /blocks/cover %}}

{{% blocks/lead %}}

The main goal of CDC is to efficiently and reliably capture real-time data changes so that they can be used to feed analysis processes or data replication.

CDC works by constantly monitoring the data source to detect any changes. When a change is detected, CDC records the details of the modification (such as which data was modified, when the modification was made, and what the previous and subsequent values were) and stores them in a dedicated table, so they can be used later for analysis or replication purposes.

CDC is a useful technique for many business applications, such as data integration, data replication, change management, and data security. 

Thanks to its ability to capture changes in real-time, CDC is particularly useful in situations where the timeliness of information is essential, such as in environmental or healthcare monitoring.

The Idra platform enables data synchronization from various sources such as APIs, sensors, message middleware, and databases. 
The platform includes a frontend for managing the connectors to be synchronized, a backend that handles data synchronization, and the Open Source ETCD database used to ensure data consistency in distributed processes.

{{% /blocks/lead %}}

