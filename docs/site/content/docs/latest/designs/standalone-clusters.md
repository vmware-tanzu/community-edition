# Standalone Workload Cluster

This design document details our desired implementation to support *Standalone workload clusters*. This design document is
a work in progress and will be updated as we implement.

## Why Standalone Workload Clusters?

Tanzu Community Edition users need a means to bootstrap clusters in a fast manner (~10 minutes or less) with minimal
resources. This enables our users to try out many projects and technology in the Tanzu portfolio with a reduced barrier
of entry.

## Current Bootstrapping Flow

Using Tanzu Kubernetes Grid (TKG), users initialize their clusters using the following commands.

```shell
# create a management cluster
tanzu management-cluster create --file ${HOME}/.config/tanzu/cluster-configs/some-config.yaml

## < ... logs for bootstrap cluster and management cluster ... >

# create workload cluster
tanzu standalone-cluster create --file ${HOME}/.config/tanzu/cluster-configs/some-config.yaml

## < ... logs for bootstrap cluster and management cluster ... >
```

When these commands are run, the following flows are triggered.

![tkg current flow](/docs/img/ttwc-current-flow.png)

As seen above, we start with a bootstrap cluster that runs in the same location as the `tanzu` CLI. This bootstrap
cluster leverages Docker and runs a
[kind](https://kind.sigs.k8s.io/) cluster. The kind cluster is injected with provider details for creating the
management cluster.

The management cluster is created in your target infrastructure (e.g. vSphere or AWS). Initially, the management cluster
is like any other cluster. This is until a **pivot** occurs, where the cluster is initialized with management components
and resources are copied over from the bootstrap cluster to this cluster. This enables the management cluster to create
and manage workloads clusters. This pivot process
is [described in the Cluster API book](https://cluster-api.sigs.k8s.io/clusterctl/commands/move.html#bootstrap--pivot).

The workload cluster is a Kubernetes cluster that does not host any provider
management capabilities. It is where Kubernetes consumers (often referred to as _application developers_) deploy their
workloads.

## Standalone Workload Cluster (SC) Creation

While the aforementioned process provides a production-capable, multi-cluster, platform, it also requires non-trivial
resources to get bootstrapped. In order to minimize resources and time required to achieve an eventual workload cluster,
TCE implements features to stop or delay the initialization of the management cluster, between bootstrap and cluster-a,
described in the previous section. The implementation of these flows would be done in a **new CLI plugin** leveraging
as much of the existing **TKG (cli) library** as possible. The command would be as follows.

> **Design note:** The TCE team has intentionally
> decided to implement this functionality in a new plugin rather than adding it to the existing cluster plugin. By implementing it in a new plugin, we are free to experiment without
> impacting the existing cluster plugin, which features an entirely different flow and set of assumptions (e.g. the presence of a management cluster).
> After this plugin sees more usage, we can re-evaluate the value of integrating its functionality into the existing cluster plugin.

```sh
tanzu standalone-cluster create --ui
```

This would trigger a flow that looks as follows.

![Standalone cluster flow](/docs/img/ttwc-minimal-flow.png)

At a lower-level, the flow triggered from the CLI would look as follows.

![Standalone cluster internal flow](/docs/img/ttwc-minimal-internal-flow.png)

Once this process is complete, the bootstrap cluster is killed and the user is left with a workload cluster, which is
not actively managed. This cluster contains `kapp-controller` and is still able to have packages installed using
the `tanzu` CLI.

## Standalone Workload Cluster (SC) Management

Eventually the SC will need to be managed again. Reasons could include:

* **scaling**: The user wants to scale the SC up or down.
* **deleting**: The user wants to delete the SC cluster.

> NOTE: At this time we do not plan to support **upgrading** of Kubernetes in this model.

In order to manage the SC, we must re-initialize the original bootstrap/management cluster to control the SC. In order
to do this efficiently, the following must be in place.

1. Create an image that contains all the management components such that they do not need to be pulled.
    * cert-manager container images loaded
    * CAPI management container images loaded
    * CAPI providers used in TCE (CAPA, CAPV, CAPD, etc) container images loaded
1. Persist all provider configuration details in `~/.config/tanzu/cluster-config` until the user deletes the SC.
1. Upon a management request of the SC, start the fully-baked image and apply the provider configuration.

Assuming an SC pre-exists, a **scaling** request would look as follows.

```shell
$ tanzu standalone-cluster scale ${SC_CLUSTER_NAME} --worker-machine-count 2

starting management components...
started management components
scaling ${SC_CLUSTER_NAME} worker nodes from 1 to 2...
scaled ${SC_CLUSTER_NAME} worker nodes to 2
stopping management components...
stopped management components
```

The flow of the above interaction would look as follows.

![SC scale flow](/docs/img/ttwc-scale-flow.png)

Assuming an MC pre-exists, a **deleting** request would look as follows.

```shell
$ tanzu standalone-cluster delete ${SC_CLUSTER_NAME}

starting management components...
started management components
deleting ${SC_CLUSTER_NAME}...
deleted ${SC_CLUSTER_NAME}
cleaned ${SC_CLUSTER_NAME} local configuration
stopping management components...
stopped management components
```

![SC delete flow](/docs/img/ttwc-delete-flow.png)

## Minimizing the Bill of Materials (BOM)

SCs will not have their core packages managed via a management cluster's `kapp-controller`. Instead, packages such as
CNI or CSI will be deployed at initialization then unmanaged (aside from standard Kubernetes reconciliation) thereafter.

Additionally, the BOM that composes a SC should be as minimal as possible. The primary packages that should be installed
as part of a SC instantiation should be:

* CNI plugin: provides networking
* CSI plugin (when available): provides storage integration
* kapp-controller: provides package management and installation features

This minimal BOM will not be reconciled after the fact. Users of the SC that install packages will have their
packages managed via the kapp-controller instance.

The flow of realizing the BOM for SC is as follows.

![BOM flow](/docs/img/ttwc-bom-flow.png)

It is TBD how the SC injector should be implemented. In theory it could be satisfied by `kapp-controller`, as long
as `kapp-controller` will not try to reconcile the things it injects into the SC if the bootstrap cluster needs to
come back up for a **scale** or **delete** operation. See last section.

## To be designed

The following items are important but not designed or prioritized for initial implementation.

* [ ] Further optimization of CAPD: While this proposal does work for CAPD, there are more optimizations we could consider.
  Namely, in the CAPD model, the bootstrap cluster already exists on the same host as the eventual SC. There is likely room to run one hybrid cluster that can self manage.

* [ ] Decision and eventual design of delayed pivot: With this Standalone model in place, we could offer a flow where users can pivot into the more production-ready model of running a dedicated management cluster.

* [ ] Minimize bootstrap tooling: Currently, users must create a cluster to get an eventual management/workload cluster. This has the upside of re-using existing controller and tooling, but downside of requiring non-trivial resources to initiate a cluster. Minimizing this would be an ideal long-term goal. Perhaps in the form a static binary.
