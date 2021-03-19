# Addons

The TCE project's goal is to provide a usable Kubernetes deployment out of the box using open-source components. A usable Kubernetes deployment requires certain capabilities such as storage and networking plugins. Additionally, there are other capabilities that are not strictly required for a Kubernetes cluster to run, but are typically expected in order to operate and maintain a cluster or run basic workloads.

TCE uses an `addon` framework to package and install these key software components. The TCE project welcomes pull-requests to add new useful addons or extend existing addons. Instructions for building and packaging an addon can be found [here](./designs/tanzu-addon-packaging.md).

## What should be an `addon`?

As general guiding principles, TCE addons should have one or more of the following attributes:

* required for a functioning cluster, such as [CNI plugins](https://github.com/containernetworking/cni) or [CSI drivers](https://kubernetes-csi.github.io/docs/).
* not strictly required, but highly desirable for a cluster operator:
  * Monitoring (`prometheus`, `grafana`)
  * Logging (`fluent-bit`)
  * Backup/restore (`velero`)
* not strictly required, but typically expected for running basic workloads:
  * Ingress Providers (`contour-operator`)

In constrast, the following types of software are probably not good fits for TCE addons:

* generic application dependencies, like relational databases (`mysql`), or object storage (`minio`)
* open-source projects used as utilities or building blocks for other applications (e.g. `nginx` or `busybox`)
