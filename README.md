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

Read our [Getting Started guide](docs/getting-started.md).

## Designs

The TCE project operates off design proposals.

* [Tanzu Add-on Packaging](./docs/designs/tanzu-addon-packaging.md)
  * How add-ons are packaged and distributed in TCE.

## Add-Ons

Add-ons provide the additional functionality necessary to build an application platform atop Kubernetes. We follow a modular approach in which operators building a platform can deploy the extensions they need to fulfill their requirements.

| Name | Description | Documentation |
|------|-------------|---------------|
| Velero | Provides disaster recovery capabilities | [Velero extension docs](./extensions/velero) |
| Gatekeeper | Provides policy enforcement within clusters | [Gatekeeper extension docs](./extensions/gatekeeper) |
| Contour | Provides ingress support to workloads | [Contour extension docs](./extensions/contour) |
| Cert Manager | Provides certificate management provisioning within the cluster | [Cert Manager extension docs](./extensions/cert-manager) |
