# Use an Existing Bootstrap Cluster to Deploy Management Clusters

By default, when you deploy a management cluster by running `tanzu management-cluster create`, Tanzu Kubernetes Grid creates a temporary `kind` cluster on the bootstrap machine, that it uses to provision the final management cluster. This temporary cluster is removed after the deployment of the final management cluster to vSphere, Amazon EC2, or Azure completes successfully. The same process of creating a temporary `kind` cluster also applies when you run `tanzu management-cluster delete` to remove a management cluster.

In some circumstances, it might be desirable to keep the local bootstrap cluster after deploying or deleting a management cluster. For example, you might want to examine the objects in the cluster or review its logs. In this case, you can skip the creation of the `kind` cluster and use any Kubernetes cluster that already exists on your bootstrap machine as the local bootstrap cluster.

- Using an existing bootstrap cluster is an advanced use case that is for experienced Kubernetes users. If possible, it is **strongly recommended** to use the default `kind` cluster that Tanzu Kubernetes Grid provides to bootstrap your management clusters.
- If you have used an existing cluster to bootstrap a management cluster, you cannot use that same cluster to bootstrap another management cluster. The same applies to deleting management clusters.

**IMPORTANT**:

The May 2021 Linux security patch causes kind clusters to fail during management cluster creation. If you run Tanzu CLI commands on a machine with a recent Linux kernel, for example Linux 5.11 and 5.12 with Fedora, `kind` clusters do not operate. This happens because `kube-proxy` attempts to change `nf_conntrack_max sysctl`, which was made read-only in the May 2021 Linux security patch, and `kube-proxy` enters a `CrashLoopBackoff state`. The security patch is being backported to all LTS kernels from 4.9 onwards, so as  operating system updates are shipped, including for Docker Machine on Mac OS and Windows Subsystem for Linux, `kind` clusters will fail, resulting in management cluster deployment failure. In this case you must do the following:

1. Download and install a version of kind that is at least version v0.11.

  For information about how to download and install kind, see the [kind documentation](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).
1. Create a `kind` cluster.

  ```
  kind create cluster
  ```

1. Follow the procedure below to run `tanzu management-cluster create` with the `--use-existing-bootstrap-cluster` option.

## Procedure

1. Set the context of `kubectl` to the local Kubernetes cluster that you want to use as a bootstrap cluster.

   ```
   kubectl config use-context my-bootstrap-cluster-admin@my-bootstrap-cluster
   ```

1. To create a management cluster, run the `tanzu management-cluster create` command and specify the `--use-existing-bootstrap-cluster` option.

   ```
   tanzu management-cluster create --file vsphere-mc.yaml --use-existing-bootstrap-cluster my-bootstrap-cluster
   ```

1. To delete a management cluster, run the `tanzu management-cluster delete` command and specify the `--use-existing-cleanup-cluster` option.

   ```
   tanzu management-cluster delete --use-existing-cleanup-cluster my-cleanup-cluster
   ```
