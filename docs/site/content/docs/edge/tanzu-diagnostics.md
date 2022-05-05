# Troubleshoot Clusters with Tanzu Diagnostics

Tanzu Community Edition comes with a diagnostics CLI plugin that helps with the task of collecting diagnostics data when debugging installation issues.  The plugin (based on the [Crash-Diagnostics project](https://github.com/vmware-tanzu/crash-diagnostics)) can automatically collect diagnostics data from either the bootstrap, management, unmanaged, or a workload cluster (or all four) using the `tanzu diagnostics collect` command.

When debugging a failed management cluster deployment, the kind bootstrap cluster is the key to being able to introspect and understand what is happening. Any issues or errors that occur in the kind bootstrap cluster during bootstrapping will provide information about potential problems in the final management cluster on the target provider. If the bootstrap process is stuck or did not finish successfully, use the diagnostics plugin to collect logs and cluster information.

## Prerequisites

Before collecting diagnostics data, your local machine must have the following programs in its $PATH:

* Tanzu CLI
* kind (for troubleshooting bootstrap cluster issues)
* kubectl

Other requirements include:

* Access to a bootstrap cluster machine (if diagnosing bootstrap cluster)
* Access to a management cluster and its managed workload clusters (if needed)
* Access to a local unmanaged cluster

### Before you Begin

If you want to debug the kind bootstrap cluster, first, determine the name. The kind bootstrap cluster name and the management cluster name are not related. If there are multiple kind bootstrap clusters present, before you begin debugging, you should determine the name of the kind cluster.

If the current context is set to the kind bootstrap cluster, run:

```sh
kubectl config view --minify
```

If current context is not set to the kind bootstrap cluster, examine the contexts to determine the name:

```sh
 kubectl config get-contexts
```

Alternatively, run

```sh
kind get clusters
```

Determine the name in one or both of the following files:

* `~/.kube-tkg/tmp/`
* `~/.config/tanzu/tkg/config.yaml` -  find the cluster name in the `name` parameter and the corresponding kind bootstrap cluster name in the `context` parameter prefixed by `kind-tkg-kind`

## Collecting Diagnostics

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

## Collecting Unmanaged Cluster Diagnostics

The tanzu diagnostics plugin also supports collecting diagnostics from [unmanaged clusters](glossary/#unmanaged-cluster).
Because unmanaged clusters do _not_ have a bootstrap cluster and are _not_ managed by a Tanzu management cluster, you must provide the name of the unmanaged cluster and the context:

```sh
tanzu diagnostics collect --unmanaged-cluster-name kind-my-unmanaged-cluster --unmanaged-cluster-context kind-my-unmanaged-cluster
```

_Note:_ The default provider for unmanaged clusters is Kind. This means that by default, the name of the unmanaged cluster will have `kind-` prefixed.

If you need to see what the name of your cluster is and its context, while targeting your unmanaged cluster, run:

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

## Tanzu Diagnostics Examples

* If the cluster bootstrap process is stuck or did not finish successfully, use the diagnostics plugin to collect logs and cluster information:

    ```sh
    tanzu diagnostics collect
    ```

    This command searches for Tanzu bootstrap kind clusters (with tkg-kind-* prefix) and collects logs and API objects (excluding [Secrets](https://kubernetes.io/docs/concepts/configuration/secret)) that can help diagnose the bootstrap issues.

* To collect diagnostics for a specific bootstrap cluster, first, determine the name of the bootstrap cluster.

    ```sh
    kind get clusters
    ```

    Next, run the `tanzu diagnostics` command to collect diagnostics data to help troubleshoot the cluster identified in the step above:

    ```sh
    tanzu diagnostics collect --boostrap-cluster-name=<BOOTSTRAP-CLUSTER-NAME>
    ```

* If you do not have a management cluster yet (or do not need to collect management cluster diagnostics), you can modify the `tanzu diagnostics` command to only collect diagnostics for the bootstrap cluster as follows:

    ```sh
    tanzu diagnostics collect --management-cluster-skip
    ```

* To collect diagnostics for a workload cluster:

  ```sh
  tanzu diagnostics collect --workload-cluster-name=<WORKLOAD_CLUSTER_NAME>
  ```

  This command will attempt to collect diagnostics data for any bootstrap cluster that still exists, diagnostics for the current management cluster, and diagnostics for the named workload cluster.  It is possible to skip collection for both the bootstrap and the management cluster as follows:

  ```sh
  tanzu diagnostics collect --bootstrap-cluster-skip --management-cluster-skip --workload-cluster-name=<WORKLOAD_CLUSTER_NAME>

## Next Steps

If the Tanzu diagnostics command did not yield the information you needed, see:

* [Troubleshooting a Bootstrap Cluster Manually](ts-manually)
* [Access local and kubectl-based Logs](logs)
