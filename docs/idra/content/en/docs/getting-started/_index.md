---
title: Getting Started
description: What does your user need to know to try your project?
categories: [Setup, Placeholders]
tags: [cdc,etcd,rest]
weight: 2
---


## Local Setup

### Prerequisites

Setup ETCD

* [ETCD](https://etcd.io/docs/): Installing ETCD

* [Install dependencies](https://go.dev/doc/modules/managing-dependencies): go build or go mod download 

### Run CDC and Web REST API

* Run CDC app from code:

Run main.go in cdc_agent folder using "go run main.go"

* Run Web app from code:

Run main.go in web folder using "go run main.go"

## Docker

Every application contains a Docker file that permits to build and run the application without to install any Golang environment.

## Kubernetes

It is possible to deploy applications using Helm charts in Kubernetes. Idra is written and inspired by a Cloud Native Approach.

