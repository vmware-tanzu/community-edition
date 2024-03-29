# Unmanaged

`unmanaged-cluster` provides single-node, static, Tanzu clusters. It is ideal for local
workstation or single-node deployments. It is not meant for production workloads
and does not offer cluster lifecycle capabilities. For needs involving
cluster lifecycle, use the `tanzu management-cluster` feature.

[unmanaged-video](https://user-images.githubusercontent.com/23109390/140801679-813a9f72-ded9-4453-b46b-75adb61e7ce2.mp4)

This feature is currently
[under proposal](https://github.com/vmware-tanzu/community-edition/issues/2266) and users should
expect instability. We appreciate all usage feedback to be added to [the
issue](https://github.com/vmware-tanzu/community-edition/issues/2266).

## Usage

### Setup

> _**Pre-reqs**: Currently, you must have docker available on the workstation you're
using `unmanaged-cluster`. This could be in the form of the docker daemon (Linux) or Docker Desktop
(Mac/Windows)._

* Will be included in the official release, remove this section when ready.

### Create a cluster

* **Mac/Linux: Create a cluster.**

    ⚠️: Depending on CDN caching, your first create may take many minutes to
    download the image(s). Subsequent creates should be significantly faster.

    ```sh
    tanzu unmanaged-cluster create hello
    ```

    > `hello` is the cluster name.

* **Windows: Create a cluster.**

    ```sh
    tanzu unmanaged-cluster create hello --cni=calico
    ```

    > `hello` is the cluster name.
    > The current version of antrea we ship cannot work on the standard WSL
    > kernel.

### Bring your own cluster

   ```sh
   tanzu unmanaged-cluster create my-cluster --config=/path/to/config
   ```

   Within the config file, you can specify a path to an existing kubeconfig file.
   This will not create a cluster but instead, will simply deploy the Tanzu bits
   ontop of your existing cluster from your existing kubeconfig.

   Warning: Providing a kubeconfig does not support setting a context.
   It is best to provide a _single cluster_ kubeconfig.

### List clusters

```sh
tanzu unmanaged-cluster ls
```

### Delete a cluster

```sh
tanzu unmanaged-cluster delete hello
```

> `hello` is the cluster name.

### Specifying a TKR

The Tanzu Kubernetes Release (TKR) identifies which packages will be run in the
cluster. You can specify your own TKRs or point to different versions available
by TCE.

For example, to track with the rest of `tanzu`, at this time, we ship on
Kubernetes 1.21.x, which also uses `kapp-controller:v0.23.0`.

To create an `unmanaged-cluster` running the newer Kubernetes `1.22.2` and
`kapp-controller:v0.25.0`, run the following:

```sh
tanzu unmanaged-cluster create hello --tkr projects.registry.vmware.com/tce/tkr:v1.22.2
```

### Provide Custom Configuration

1. Generate a config file with defaults

    ```sh
    tanzu unmanaged-cluster configure hello
    ```

    > `hello` is the cluster name and will generate `./hello.yaml`.

1. Modify the configuration (`hello.yaml`) as desired.

1. Create the cluster with the custom configuration.

    ```sh
    tanzu unmanaged-cluster create hello -f hello.yaml
    ```

### Interacting with Clusters

Upon successful bootstrap, we automatically set your default kube context to the
newly created cluster. This means all operations you're used to should work as
is. For example:

* List running pods

    ```sh
    kubectl get po -A
    ```

* List available packages

    ```sh
    tanzu package available list
    ```

## Unmanaged as an API

While `unmanaged-cluster` provides cluster creation ability via CLI, it can also be
called programmatically to install unmanaged on most arbitrary clusters.
This can be especially compelling for projects that handle the underlying
provisioning of the VM and container runtime but are looking for the higher-level
Tanzu bits to be installed atop. To get started, try:

1. Import the `tanzu` package of unmanaged-cluster to your project.

    ```sh
    go get -d github.com/vmware-tanzu/community-edition/standalone-overhaul/cli/cmd/plugin/unmanaged-cluster@standalone-overhaul
    ```

1. Setup your project to use the manager instance.

    ```go
    package main

    import (
        /* your deps */
        "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/config"
        "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/log"
        "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster/tanzu"
    )

    func main() {
        // provides stylized logging. First argument is whether to disable
        // stylization (tty: false). Second argument is to set the displayed log
        // level.
        log := logger.NewLogger(true, 0)

        // New tanzu manager
        tm := tanzu.New(log)

        // settings for how to create the cluster
        clusterConfig := config.UnmanagedClusterConfig{}

        // deploy the cluster, by default using kind
        err = tm.Deploy(clusterConfig)
        if err != nil {
          return err
        }

        // list clusters
        err = tm.List()
        if err != nil {
          return err
        }

        // delete clusters
        err = tm.Delete("cluster-name")
        if err != nil {
          return err
        }
    }
    ```

    > If you do not want to use our `log` package, you can implement the
    > `log.Logger` interface.

In the above example, the `config.UnmanagedClusterConfig{}` struct determines how
the Tanzu cluster is deployed. It features options for specifying how the
Kubernetes cluster should be created to the CNI that runs on top of it. For
example, if you wish to pre-create the Kubernetes environment before calling
deploy, you could do so with the following config:

```go
// settings for how to create the cluster
clusterConfig := config.UnmanagedClusterConfig{
    ExistingKubeconfig: kubeConfigByteArray,
    Provider: "none",
}
```

> By setting the `Provider` value to `none`, you are instructing `unmanaged-cluster`
> that it should not concern itself with creating the underlying Kubernetes
> environment. Instead, it should interact with an existing API server to
> bootstrap the Tanzu components.

## Design Details

This section contains design details, package architecture, and historical
context around the why of this feature.

### Philosophy and history

> _note: much of this is covered [in our
proposal](https://github.com/vmware-tanzu/community-edition/issues/2266)_

`unmanaged-cluster`, previously `standalone-cluster` has existed since the public
release of Tanzu Community Edition (TCE) and used internally at VMware for
several months prior. The original intent was to provide a quick and easy time
to workload cluster, which would lower the barrier of entry to trying Tanzu. Our
original implementation attempted to re-use [Cluster
API](https://cluster-api.sigs.k8s.io/) and
[tanzu-framework](https://github.com/vmware-tanzu/tanzu-framework) in the same
way as the capable `management-cluster` functionality. Over time, we've learned
that even with changes to these dependencies, the need for a bootstrap cluster,
to create a management cluster, which eventually processes a Tanzu Kubernetes
Release (TKR) to create a workload cluster is far too heavy weight for users
looking to get started[0]. These users are often looking to deploy to a single node
or on their local workstation. Additionally, they are rarely concerned with
long-running clusters, hosting production workload, or simulating cluster
lifecycle. When this is required, it is a far better option to use the
`management-cluster` model anyway.

With the above, the new implementation of `unmanaged-cluster` aims to solve:

* How can we lower the barrier of entry to using Tanzu?
* How can we provide an exceptional on-ramp to the user persona that is not
  running production workloads tomorrow.
* How can we provide a quicker feedback cycle for developers, package authors,
  and TKR creators?

> [0]: This is backed up by the countless instances of troubleshooting users
> deploying local (CAPD-backed) unmanaged-clusters. We found the issue was
> almost always sourced in system resource constraints, and with so many moving
> parts to get a single-node cluster, the troubleshooting was far too complex
> for a new user.

### The new standalone

The new `unmanaged-cluster` implementation creates single node, static, environments of
Tanzu. This enables users to run Tanzu environments on single VMs or local
workstations.

Existing, managed, Tanzu environments rely on the processing of a Tanzu
Kubernetes Release (TKR). This can be thought of as a Bill of Materials (BOM)
that provides instructions for how a target workload cluster should be created.
Using Tanzu's robust package model, this largely comes down to specifying a base
image and a set of packages to manage on top. Tanzu's managed clusters process
all of this in a long-running management cluster. Which provides extremely
capable multi-cluster management.

`unmanaged-cluster` does **not** utilize a long-running management cluster. Instead, it
processes the TKR client-side and bootstraps the cluster based on interacting
with its API server. This works well with single-node use cases as the
underlying cluster can be bootstrapped in a variety of ways, through kind (our
default), minikube, or even a cluster-api provider.

This provides a great toolset for users looking to:

* Get started with Tanzu
* Run single-node experiments with Tanzu
* Author packages
* Author TKRs
* Integrate an unmanaged-Tanzu environment into their CI/CD

### Cluster Infrastructure

By default `unmanaged-cluster` calls to [kind](https://kind.sigs.k8s.io) as the
lower-level Kubernetes subsystem. We embed kind in our `unmanaged-cluster` binary.
In order to use this plugin, pre-reqs must be setup such as Docker and, if not
running Linux, a relevant Linux VM where Docker and its transitive dependencies
can run. Historically, many accomplished this using Docker Desktop.

It is important to note `unmanaged-cluster` **is not** in the business of providing
cluster infrastructure or becoming a desktop implementation of Tanzu. What
unmanaged does care about is having an API server which it can satisfy the TKR
against. This leaves the door open to use as both CLI and API for
satisfying the single-node, static, and quick-bootstrap use case. We expect
cluster providers, say VMware tools that provision Kubernetes to be able to offer this
single-node use case by provisioning a VM, container runtime, then calling our
API with the Provider config set to `none`. This will instruct `unmanaged-cluster` to
do no infrastructure boostrapping and instead interact with pre-provisioned
assets.

### Configuration

Simple, robust, configuration is a primary concern of `unmanaged-cluster` along with
consistency for CLI and API consumers. As such, the (evolving) configuration can
be found in our `config` package. We believe in sensible defaults that can be
easily overwritten and understood. With this, no configuration is need to run
`tanzu unmanaged-cluster create`. However, a configuration file with all the defaults
can be generated by running `tanzu unmanaged-cluster configure`. For every field in the
configuration struct, we support a correlated environment variable and CLI flag.
Users can set configuration in a config file, environment variables, and flags.
To accomplish this, we offer the following configuration precedence:

![Unmanaged configuration
precedence](../../../../docs/images/stanadlone-config-precedence.png)

As seen above, we also persist the **rendered** configuration to
`~/.config/tanzu/tkg/unmanaged/${CLUSTER_NAME}/config.yaml`. This provides
users with a concrete way to understand the end-state configuration that was
used.

Along with the configuration file, we also store multiple assets in this
directory, including:

* The `kubeconfig` file.
* The bootstrapping logs.

### Package Architecture

At a code-level, there are multiple packages that make `unmanaged-cluster` possible.

![Unmanaged Package
Architecture](../../../../docs/images/stanadlone-package-arch.png)

* `cmd`: Cobra commands for interacting with `unmanaged-cluster` via CLI.
* `tanzu`: Orchestrator of operations. The Manager interface provides the API
  for which `unmanaged-cluster` can be called programmatically.
* `tkr`: Responsible for processing a TKR and resolving the various packages
  (OCI bundles) that it makes up.
* `cluster`: Manages underlying Kubernetes clusters to deploy `unmanaged-cluster` atop.
  The `ClusterManager` interface enables creation of multiple providers. By
  default, we use a `kind` implementation. The cluster package can also be
  bypassed when configuration is provider specifying the provider is `none`.
* `kapp`: Facilitates the deployment of `kapp-controller` into the cluster. This
  facilitates the management of all packages in the cluster. It is always
  deployed first as a Pod that can tolerate non-ready nodes and run on the host
  network, enabling it to run before a CNI so it can deploy a CNI package.
* `packages`: Facilitates deployments of packages (OCI bundles) and package
  repositories. It is leveraged for installing the CNI package along with
  higher-level packages that are desired.
* `kubeconfig`: Facilitates the management of kubeconfigs. Largely is used to
  manage the kubeconfig while `unmanaged-cluster` bootstraps. It also, as a final step,
  is used to add the cluster record to the default `~/.kube/config` and
  automatically switch the user's context.
* Utility packages:
  * `config`: Contains the configuration types and various config operations.
    The configuration type in the package is used by the CLI and API callers to
    feed the `tanzu` package.
  * `log`: Logging utilities that provide detailed bootstrap logs and
    user-friendly CLI logs.

### Deprecation of Existing Unmanaged Clusters

For existing `standalone-cluster` users running locally (CAPD), due to the
limitations of `standalone-cluster`, you should be able to start using the new
`unmanaged-cluster` feature today. You'll likely find the bootstrapping experience
better, faster, and the feature set to be far more capable and rich.

Those pointing `standalone-cluster` at a infrastructure provider such as AWS,
vSphere, or Azure can get nearly the same functionality (and much more) by
creating `management-cluster`(s) going forward. In fact, the old
`standalone-cluster` feature was calling almost the exact same code path as the
`management-cluster` plugin.
