---
title: "Tutorial: Setup a data connector"
description: >
 
  The **Idra** platform enables _data synchronization_ from various sources such as DBs, APIs, sensors, and message middleware. To enable it, needs setup a _data connector_. 
  <br>Below we will analyze all the above cases.
 

date: 2017-01-05
weight: 5
---

### Data synchronization between DBs
**1)** (_Postgres/MySql_ with _Timestamp_)
<br>Suppose we have two databases _Postgres_ and _MySql_, the first one is the _source_, the second one is the _destination_. 
<br>The first thing to do is to log in so as explained in the section **_[Web UI](http://localhost:1313/docs/tasks/webui/)_**.
<br>Once logged in, you need to click on the 'Syncs' menu item and then on the 'Create Sync' button, where we setup 3 sections: <br>_Sync Details_, _Source Connector_ and _Destination Connector_.
<br>In the first section, you need to assign a suitable _Sync Name_ to the synchronization, and also choose a synchronization _Mode_. 
<br>With regard to the last 2 sections, where the _data connector_ is configured, the first one is about configuring the _Source Connector_, the second one further down is about the _Destination Connector_.
<br>As you can see in the 2 images below, mode _Timestamp_ has been selected.
<br>In addition, the connection strings for _Postgres_ and _MySQL_ have been properly configured respectively, following the rules indicated here:
[GORM](https://gorm.io/docs/connecting_to_the_database.html).
<br>Moreover, it's important to create the _Timestamp Field_ in both tables of the 2 databases with the exact same name!
<br>Finally, once all these things have been done, you can click on the 'Create' button (see the second image below).
<br>
<br>![](/images/7.png)
<br>Connection Strings is: _host=localhost user=postgres password=mimmo dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Rome_
<br><br>![](/images/8.png)
<br>Connection Strings is: _root:@tcp(127.0.0.1:3306)/scuola?charset=utf8mb4&parseTime=True&loc=Local_
<br><br><br>
**2)** (_Postgres/MySql_ with _ConnectorId_)

<br><br><br><br><br>




