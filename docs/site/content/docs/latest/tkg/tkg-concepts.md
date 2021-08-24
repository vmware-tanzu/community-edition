# Tanzu Kubernetes Grid Concepts

This topic describes the key elements and concepts of a Tanzu Kubernetes Grid deployment.

## <a id="mgmt-cluster"></a> Management Cluster

A management cluster is the first element that you deploy when you create a Tanzu Kubernetes Grid instance. The management cluster is a Kubernetes cluster that performs the role of the primary management and operational center for the Tanzu Kubernetes Grid instance. This is where Cluster API runs to create the Tanzu Kubernetes clusters in which your application workloads run, and where you configure the shared and in-cluster services that the clusters use.

**NOTE**: On vSphere 7, it is recommended to use a built-in supervisor cluster from vSphere with Tanzu instead of deploying a Tanzu Kubernetes Grid management cluster. Deploying a Tanzu Kubernetes Grid management cluster to vSphere 7 when vSphere with Tanzu is not enabled is supported, but the preferred option is to enable vSphere with Tanzu and use the Supervisor Cluster. For details, see [vSphere with Tanzu Provides Management Cluster](mgmt-clusters/vsphere.md#mc-vsphere7).

When you deploy a management cluster, networking with [Antrea](https://antrea.io/) is automatically enabled in the management cluster.  The management cluster is purpose-built for operating the platform and managing the lifecycle of Tanzu Kubernetes clusters.  As such, the management cluster should not be used as a general purpose compute environment for end-user workloads.

## <a id="clusters"></a> Tanzu Kubernetes Clusters

After you have deployed a management cluster, you use the Tanzu CLI to deploy CNCF conformant Kubernetes clusters and manage their lifecycle. These clusters, known as Tanzu Kubernetes clusters, are the clusters that handle your application workloads, that you manage through the management cluster. Tanzu Kubernetes clusters can run different versions of Kubernetes, depending on the needs of the applications they run. You can manage the entire lifecycle of Tanzu Kubernetes clusters by using the Tanzu CLI. Tanzu Kubernetes clusters implement Antrea for pod-to-pod networking by default.

## <a id="plans"></a> Tanzu Kubernetes Cluster Plans

A cluster plan is the blueprint that describes the configuration with which to deploy a Tanzu Kubernetes cluster. It provides a set of configurable values that describe settings like the number of control plane machines, worker machines, VM types, and so on.

This release of Tanzu Kubernetes Grid provides two default templates, `dev` and `prod`.

## <a id="services"></a> Shared and In-Cluster Services

Shared and in-cluster services are services that run in the Tanzu Kubernetes Grid instance, to provide authentication and authorization of Tanzu Kubernetes clusters, logging, and ingress control.

## <a id="instance"></a> Tanzu Kubernetes Grid Instance

A Tanzu Kubernetes Grid instance is a full deployment of Tanzu Kubernetes Grid, including the management cluster, the deployed Tanzu Kubernetes clusters, and the shared and in-cluster services that you configure. You can operate many instances of Tanzu Kubernetes Grid, for different environments, such as production, staging, and test; for different IaaS providers, such as vSphere, Azure, and Amazon EC2; and for different failure domains, for example Datacenter-1, AWS us-east-2, or AWS us-west-2.

## <a id="bootstrap"></a> Bootstrap Machine

The bootstrap machine is the laptop, host, or server on which you download and run the Tanzu CLI. This is where the initial bootstrapping of a management cluster occurs, before it is pushed to the platform where it will run.

## <a id="installer"></a> Tanzu Kubernetes Grid Installer

The Tanzu Kubernetes Grid installer is a graphical wizard that you start up by running the `tanzu management-cluster create --ui` command. The installer wizard runs locally on the bootstrap machine, and provides a user interface to guide you through the process of deploying a management cluster.

## <a id="upgrades"></a> Tanzu Kubernetes Grid and Cluster Upgrades

Upgrading a Tanzu Kubernetes Grid release means upgrading the management clusters created by the CLI version of that release.

Upgrading a management or Tanzu Kubernetes (workload) cluster in Tanzu Kubernetes Grid means migrating its nodes to run on a base VM image with a newer version of Kubernetes:

  - **Management clusters** upgrade to the latest available version of Kubernetes.
  - **Workload clusters** upgrade by default to the current Kubernetes version of their management cluster.  Or you can specify other, non-default Kubernetes versions to upgrade workload clusters to.

To find out which Kubernetes versions are available in Tanzu Kubernetes Grid, see [List Available Versions](tanzu-k8s-clusters/k8s-versions.html#k8s-vers-list).
