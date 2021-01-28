# Tanzu Community Edition

A modular application platform built atop Kubernetes.

## Overview

Tanzu Community Edition (TCE) enables the creation of application platforms.
Leveraging [Cluster API](https://cluster-api.sigs.k8s.io/), Kubernetes is used
as the foundational way to schedule and orchestrate workloads. With Kubernetes
in place, TCE enables the installation of platform extensions that support
software running in clusters. While an opinionated set of extensions is offered,
TCE is modular and enables you to bring your own.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

To use TCE, infrastructure where Kubernetes clusters can be bootstrapped is
required. TCE supports vSphere 6.7.X, AWS, and Azure.

### Installing

Step by step guide on how to install TCE

## Extensions

Extensions provide the additional functionality necessary to build an application platform atop Kubernetes. We follow a modular approach in which operators building a platform can deploy the extensions they need to fulfill their requirements.

| Name | Description | Documentation |
|------|-------------|---------------|
| Velero | Provides disaster recovery capabilities | [Velero extension docs](./extensions/velero) |

