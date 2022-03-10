# Troubleshoot Clusters with Tanzu Diagnostics

Tanzu Community Edition comes with a diagnostics CLI plugin that helps with the task of collecting diagnostics data when debugging installation issues.  The plugin (based on the [Crash-Diagnostics project](https://github.com/vmware-tanzu/crash-diagnostics)) can automatically collect diagnostics data from either the boostrap, management, unmanaged, or a workload cluster (or all three) using the `tanzu diagnostics collect` command.

## Pre-requisites

Prior to collecting diagnostics data, your local machine must have the following programs in its $PATH:

* Tanzu CLI
* kind (for troubleshooting bootstrap cluster issues)
* kubectl

Other requirements include:

* Access to a bootstrap cluster machine (if diagnosing bootstrap cluster)
* Access to a management cluster and its managed workload clusters (if needed)
* Access to a local unmanaged cluster

## Collecting diagnostics

The diagnostics tool can automatically collect logs, and API resources data, for all cluster types including bootstrap, management, and workload clusters.  By default, the diagnostics plugin will attempt to automatically collect data in the following order:

```text
[bootstrap cluster] --> [management cluster] --> [workload cluster]
```

The diagnostics plugin will follow the order above unless the user specifically skips either the bootstrap or the management cluster using a CLI argument. If a cluster type is not available (i.e. bootstrap for instance) the tool will simply skip it.

### What is collected?

The diagnostics plugin collects a set of known resources to help troubleshoot bootstrap cluster problems, including:

* Kind cluster logs (obtained with command kind export logs)
* Pod logs from `capi*, capv*, capa*, tkg-system` namespaces (if they exist)
* Pods, services, deployments, apps (from `capi*, capv*, capa*, tkg-system` namespaces)
* Any other server resources in the `cluster-api` resource category

No other cluster objects are collected.

The `tanzu diagnostics` command automatically generates a tar file which you can unpack to analyze and troubleshoot your cluster.

### The `tanzu diagnostics collect` command

The diagnostics command makes it easy for users to collect diagnostics data by supporting sensible default values. However, the command allows its default values to be overridden using CLI argument flags as listed below:

```sh
tanzu diagnostics collect --help
Collect cluster diagnostics for the specified cluster

Usage:
  diagnostics collect [flags]

Flags:
      --bootstrap-cluster-name string          A specific bootstrap cluster to diagnose
      --bootstrap-cluster-skip                 If true, skips bootstrap cluster diagnostics
  -h, --help                                   help for collect
      --management-cluster-context string      The context name of the management cluster
      --management-cluster-kubeconfig string   The management cluster config file (required)
      --management-cluster-name string         The name of the management cluster (required)
      --management-cluster-skip                If true, skips management cluster diagnostics
      --output-dir string                      Output directory for collected bundle (default "./")
      --unmanaged-cluster-context string       The context name of the unmanaged cluster
      --unmanaged-cluster-kubeconfig string    The unmanaged cluster config file (required) (default "${HOME}/.kube/config")
      --unmanaged-cluster-name string          The name for the unmanaged cluster (required)
      --work-dir string                        Working directory for collected data (default "${HOME}/.config/tanzu/diagnostics")
      --workload-cluster-context string        The context name of the workload cluster
      --workload-cluster-infra string          Overrides the infrastructure type for the managed cluster (i.e. aws, azure, vsphere, etc) (default "docker")
      --workload-cluster-kubeconfig string     The workload cluster config file
      --workload-cluster-name string           The name of the managed cluster for which to collect diagnostics (required)
      --workload-cluster-namespace string      The namespace where managed workload resources are stored (default "default")
```

## Collecting boostrap cluster diagnostics

When setting up a management cluster, for the first time, it is possible for things to go wrong during the bootstrap stage or while pivoting to the management cluster stage.  If the bootstrap process is stuck or did not finish successfully, use the diagnostics plugin to collect logs and cluster information:

```sh
tanzu diagnostics collect
```

This command searches for Tanzu bootstrap kind clusters (with tkg-kind-* prefix) and collects logs and API objects that can help diagnose the bootstrap issues. You should see output similar to the following:

```sh
tanzu diagnostics collect
2021/09/20 11:05:17 Collecting bootstrap cluster diagnostics
2021/09/20 11:05:17 Info: Found bootstrap cluster(s): ["tkg-kind-b4o9sn5948199qbgca8d"]
2021/09/20 11:05:17 Info: Bootstrap cluster: tkg-kind-b4o9sn5948199qbgca8d: capturing node logs
2021/09/20 11:05:22 Info: Capturing pod logs: cluster=tkg-kind-b4o9sn5948199qbgca8d
2021/09/20 11:05:22 Info: Capturing API objects: cluster=tkg-kind-b4o9sn5948199qbgca8d
2021/09/20 11:05:25 Info: Archiving: bootstrap.tkg-kind-b4o9sn5948199qbgca8d.diagnostics.tar.gz
```

For more bootstrap cluster debugging, see additional instructions [here](./tsg-bootstrap).

## Collecting management cluster diagnostics

When your management cluster is not deployed properly, you can use the diagnostics plugin to collect data to help debug the issue. As before, use the following command collect diagnostics data:

```sh
tanzu diagnostics collect
```

The plugin will attempt to collect diagnostics for any boostrap cluster that still exists and then continues to collect diagnostics for the current management cluster.  You can specifically skip collection of your bootstrap cluster (if it exists) using the following command:

```shell
tanzu diagnostics collect --bootstrap-cluster-skip
```

## Collecting workload cluster diagnostics

The plugin can also automate the collection of diagnostics data to debug issues with a workload cluster setup. This can be done by invoking the `diagnostics collect` command with the name of the workload cluster:

```sh
tanzu diagnostics collect --workload-cluster-name=<WORKLOAD_CLUSTER_NAME>
```

This command will attempt to collect diagnostics data for any boostrap cluster that still exists, diagnostics for the current management cluster, and diagnostics for the named workload cluster.  It is possible to skip collection for both the bootstrap and the management cluster as follows:

```sh
tanzu diagnostics collect --bootstrap-cluster-skip --management-cluster-skip --workload-cluster-name=<WORKLOAD_CLUSTER_NAME>
```

## Collecting unmanaged cluster diagnostics

Finally, the plugin also supports collecting diagnostics from unmanaged-clusters, a lite weight, local deployment model.
Because unmanaged-clusters do _not_ have a bootstrap cluster and are _not_ managed by a Tanzu management cluster,
you must provide the name of the unmanaged cluster and the context:

```sh
tanzu diagnostics collect --unmanaged-cluster-name kind-my-unmanaged-cluster --unmanaged-cluster-context kind-my-unmanaged-cluster
```

_Note:_ The default provider for unmanaged-clusters is Kind. This means that by default, the name of the unmanaged-cluster will
have `kind-` prefixed.

If you need to see what the name of your cluster is and its context, while targeting your unmanged cluster, run:

```sh
kubectl config view --minify
```

Make sure you have a context set first! You can see your current context with

```sh
kubectl config get-contexts
```

and set your context with

```sh
kubectl config set-context <name-of-context>
```
