# Architecture

Tanzu Community Edition consists of a variety of components that enable the
bootstrapping and management of Kubernetes clusters and the various platform
services ran atop. This page details the architecture of:

* Tanzu CLI
* Managed Clusters
* Standalone Clusters
* Package Management

## Tanzu CLI

The `tanzu` CLI command exposes multiple subcommands.

```sh
$ tanzu

Tanzu CLI

Usage:
  tanzu [command]

Available command groups:

  Admin
    builder                 Build Tanzu components
    codegen                 Tanzu code generation tool
    test                    Test the CLI

  Run
    cluster                 Kubernetes cluster operations
    kubernetes-release      Kubernetes release operations
    management-cluster      Kubernetes management cluster operations
    package                 Tanzu package management
    standalone-cluster      Create clusters without a dedicated management cluster

  System
    completion              Output shell completion code
    config                  Configuration for the CLI
    init                    Initialize the CLI
    login                   Login to the platform
    plugin                  Manage CLI plugins
    update                  Update the CLI
    version                 Version information


Flags:
  -h, --help   help for tanzu

Use "tanzu [command] --help" for more information about a command.

Not logged in
```

Each subcommand provides functionality for Tanzu Community Edition. This
functionality can range from creating clusters to managing the software running
in clusters. Subcommands in `tanzu` are independent static binaries hosted on a
client system. This enables a pluggable architecture where plugins can be added,
removed, and updated independent of each other. The `tanzu` command is expected
to be installed in machine's path. Each subcommand (binary) is expected to be
installed in `${XDG_DATA_HOME}/tanzu-cli`. This relationship is demonstrated
below.

![CLI Architecture](../../img/cli-arch.png)

> [Click here to see where $XDG_DATA_HOME resolves
> to.](https://github.com/adrg/xdg#xdg-base-directory)

Tanzu Community Edition ships with the `tanzu` CLI and a select set of plugins.
Some plugins may live in the vmware-tanzu/community-edition repository while
others live in vmware-tanzu/tanzu-framework. Plugins that live in
vmware-tanzu/tanzu-framework may be used in multiple Tanzu Editions. Plugins
that live in vmware-tanzu/community-edition are used exclusively in Tanzu
Community Edition. Plugins in vmware-tanzu/community-edition may be promoted
(moved) to vmware-taznu/tanzu-framework. This move would not impact users of
Tanzu Community Edition; it would only impact contributors to the plugin.

## Managed Clusters

Clusters that are deployed and managed using centralized management clusters are
considered managed clusters. This is the primary deployment model for clusters
in the Tanzu ecosystem and is recommended for production scenarios. To bootstrap
managed clusters, you first need a management cluster.  This is done using the
`tanzu management-cluster create` command. When running this command, a
bootstrap cluster is created locally and is used to then create the management
cluster. The following diagram shows this flow.

![Bootstrap cluster create](../../img/bootstrap-cluster-create.png)

Once the management cluster has been created, the bootstrap cluster will perform
a move (aka pivot) of all management objects to the management cluster. From
this point forward, the management cluster is responsible for managing itself
and any new clusters you create. These new clusters, managed by the management
cluster, are referred to as workload clusters. The following diagram shows this
relationship end-to-end.

![Management cluster bootstrapping](../../img/management-cluster-flow.png)

## Standalone Clusters 

Clusters that run without a long-running management cluster are considered
standalone clusters. This is an experimental cluster deployment model currently
being iterated on by the Tanzu Community Edition team. This model provides a few
benefits including:

* Faster time to cluster (relative to managed clusters)
* Reduced system requirements

As such, this model is not recommended for production workloads.

Creating a standalone cluster is done using the `tanzu standalone-cluster
create` command. When running this command, a bootstrap cluster is created
locally and is then used to create the standalone cluster. After successful
bootstrapping, the bootstrap cluster is deleted. Management resources are
**not** moved into the standalone cluster. This newly created cluster can be
referred to as a workload cluster.  The following diagram shows this
relationship.

![Standalone cluster flow](../../img/standalone-cluster-flow.png)

When you'd like to delete or scale the workload cluster, a new bootstrap
cluster is created and the workload cluster is deleted or scaled. This bootstrap
cluster can be thought of as a _temporary_ management cluster. Users should
expect a delay in the operation as there will be time lost to re-creating the
boostrap cluster. The following diagram shows this relationship.

![Standalone scale example flow](../../img/flow-for-standalone-mutation.png)

## Package Management

Tanzu Community Edition provides package management to users via the `tanzu`
CLI. Package management is defined as the discovery, installation, upgrading,
and deletion of software that runs on Tanzu clusters. Each package is created
using [carvel tools](https://carvel.dev/) and following our [packaging
process](designs/package-process). Packages are put into a single bundle,
called a package repository and pushed to an OCI-compliant registry. In Tanzu
clusters, [kapp-controller](https://carvel.dev/kapp-controller). When a cluster
is told about this package repository (likely via the `tanzu package repository`
command), kapp-controller can pull down that repository and make all the packages
available to the cluster. This relationship is shown below.

![kapp-controller repo read](../../img/tanzu-carvel-new-apis.png)

With the packages available in the cluster, users of `tanzu` can install various
packages. Under the hood, the creates an [PackageInstall](https://carvel.dev/kapp-controller/docs/latest/packaging/#packageinstall) resource that
instructs `kapp-controller` to download the package and install the software in
your cluster. This flow is shown below.

![tanzu package install](../../img/tanzu-package-install-2.png)
