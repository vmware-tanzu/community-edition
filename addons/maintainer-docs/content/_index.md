---
title: Package Ownership and Maintenance
weight: 1
---

## What is a Package

As defined in a Carvel [blog](https://carvel.dev/blog/introduction-to-carvel-package-manager-for-kubernetes/), a package is versioned metadata which informs kapp-controller how to fetch, template, and install the underlying software contents. These contents usually consist of configuration and container images which have been bundled together and stored in some location.

Tanzu Community Edition makes use of kapp-controller to install its various components and is an important part of the extensibility of the system. This documentation will provide a guide to creating a package.

Tanzu Community Edition classifies packages as either Core or User Managed. Core packages are those that are necessary to bootstrap clusters and that are required to actually run the management and workload clusters. They are defined by the upstream Tanzu Framework package and cannot be changed by end users. Examples are Antrea, kapp-controller, metrics-server, etc.

User managed packages are those that are available to be installed by an end user to a running cluster. These packages are optional installs, not required by TCE to manage or run clusters. Examples of user-managed packages are Contour, Harbor, Grafana, etc.
