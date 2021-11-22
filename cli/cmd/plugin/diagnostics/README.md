# Tanzu Diagnostics Plugin

The Tanzu Diagnostics Plugin allows operators to collect cluster diagnostics data. The plugin uses the [Crashd API](https://github.com/vmware-tanzu/crash-diagnostics) to run internal script that automates the collection of diagnostics information for cluster troubleshooting.

## Collecting diagnostics

By default, the diagnostics plugin will attempt to collect information from:

* Any bootstrap cluster (with name `tkg-kind-*`) in kind
* The current Tanzu management cluster
* Any specified workload cluster
* Any specified standalone cluster

For instance, the following command will collect logs, API objects, and other API server info:

```shell
tanzu diagnostics collect
2021/09/02 08:01:02 Collecting bootstrap cluster diagnostics
2021/09/02 08:01:03 Info: Found bootstrap cluster(s): ["tkg-kind-b4o9sn5948199qbgca8d"]
2021/09/02 08:01:03 Info: Bootstrap cluster: tkg-kind-b4o9sn5948199qbgca8d: capturing node logs
2021/09/02 08:01:10 Info: Capturing pod logs: cluster=tkg-kind-b4o9sn5948199qbgca8d
2021/09/02 08:01:10 Info: Capturing API objects: cluster=tkg-kind-b4o9sn5948199qbgca8d
2021/09/02 08:01:13 Info: Archiving: bootstrap.tkg-kind-b4o9sn5948199qbgca8d.diagnostics.tar.gz
2021/09/02 08:01:13 Info: Capturing management cluster diagnostics
...
```

The `collect` command collects diagnostics data and creates tarball file for each cluster type.

### Collecting workload cluster diagnostics

To include diagnostics data from a managed workload cluster, simply provide its cluster name:

```shell
tanzu diagnostics collect --workload-cluster-name=wc-webtier-1
```

Specify the workload cluster's namespace if needed:

```shell
tanzu diagnostics collect --workload-cluster-name=wc-webtier-1 --workload-cluster-namespace="ns-webtier"
```

### Skipping bootstrap and management clusters

In certain instances, it may be useful to skip collection of the either the bootstrap or the management cluster. This can be done as follows:

```bash
tanzu diagnostics collect --bootstrap-cluster-skip=true --management-cluster-skip=true --workload-cluster-name=wc-webtier-1
```

The command above will collect diagnostics only for the specified workload cluster.

## Command arguments

The following shows a list of command arguments that can be used to override default values when collecting diagnostics.

```shell
tanzu diagnostics collect --help
Collect cluster diagnostics for the specified cluster

Usage:
  tanzu diagnostics collect [flags]

Flags:
      --bootstrap-cluster-name string          A specific bootstrap cluster to diagnose
      --bootstrap-cluster-skip                 If true, skips bootstrap cluster diagnostics
  -h, --help                                   help for collect
      --management-cluster-context string      The context name of the management cluster
      --management-cluster-kubeconfig string   The management cluster config file (required)
      --management-cluster-name string         The name of the management cluster (required)
      --management-cluster-skip                If true, skips management cluster diagnostics
      --output-dir string                      Output directory for collected bundle (default "./")
      --standalone-cluster-context string      The context name of the standalone cluster
      --standalone-cluster-kubeconfig string   The standalone cluster config file (required) (default "${HOME}/.kube/config")
      --standalone-cluster-name string         The name for the standalone cluster (required)
      --work-dir string                        Working directory for collected data (default "${HOME}/.config/tanzu/diagnostics")
      --workload-cluster-context string        The context name of the workload cluster
      --workload-cluster-infra string          Overrides the infrastructure type for the managed cluster (i.e. aws, azure, vsphere, etc) (default "docker")
      --workload-cluster-kubeconfig string     The workload cluster config file
      --workload-cluster-name string           The name of the managed cluster for which to collect diagnostics (required)
      --workload-cluster-namespace string      The namespace where managed workload resources are stored (default "default")
```
