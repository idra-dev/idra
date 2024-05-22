---
title: API Server
description: >
  #### Data Management Rest API Server
date: 2017-01-05
categories: [Architecture, Docs]
tags: [api]
weight: 4
---

The Web server allows access to all synchronization information present in ETCD via API. 
The API server is written in Golang using the Gin framework, and this server is used by the Web client UI. GIN is a lightweight and fast web framework written in Go that enables the creation of scalable and high-performance web applications. Here are some of its key features:

Routing: 

GIN offers a flexible and easy-to-use routing system, allowing for efficient handling of HTTP requests. You can define routes, manage route parameters, use middleware to filter requests, and more.

Middleware: 

GIN supports the use of middleware to modularly handle HTTP requests. There are many middleware available, including logging middleware, error handling middleware, security middleware, and more.

Binding: 

GIN offers a binding system that automatically binds HTTP request data to your application's data types. You can easily handle form data, JSON data, XML data, and more.

Rendering: 

GIN provides a flexible and easy-to-use rendering system, allowing for easy generation of HTML, JSON, XML, and other formats.

Testing: 

GIN provides a great testing experience, with features such as integration test support and the ability to easily and intuitively test HTTP calls.

Performance: 

GIN is known for its high performance and ability to easily handle high-intensity workloads. You can use GIN to create high-performance web applications, even in high concurrency environments.

In summary, GIN is an extremely useful web framework for creating web applications in Go. Thanks to its flexibility, high performance, and wide range of features.
