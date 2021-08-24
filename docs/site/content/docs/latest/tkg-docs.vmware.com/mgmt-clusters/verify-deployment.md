# Examine the Management Cluster Deployment

During the deployment of the management cluster, either from the installer interface or the CLI, Tanzu Kubernetes Grid creates a temporary management cluster using a [Kubernetes in Docker](https://kind.sigs.k8s.io/), `kind`, cluster on the bootstrap machine. Then, Tanzu Kubernetes Grid uses it to provision the final management cluster on the platform of your choice, depending on whether you are deploying to vSphere, Amazon EC2, or Microsoft Azure. After the deployment of the management cluster finishes successfully, Tanzu Kubernetes Grid deletes the temporary `kind` cluster.

When Tanzu Kubernetes Grid creates a management cluster for the first time, it also creates a folder `~/.tanzu/tkg/providers` that contains all of the files required by Cluster API to create the management cluster.

The Tanzu Kubernetes Grid installer interface saves the settings for the management cluster that it creates into a cluster configuration file `~/.tanzu/tkg/clusterconfigs/UNIQUE-ID.yaml`, where `UNIQUE-ID` is a generated filename.

**IMPORTANT**: By default, unless you set the `KUBECONFIG` environment variable to save the `kubeconfig` for a cluster to a specific file, all clusters that you deploy from the Tanzu CLI are added to a shared `.kube-tkg/config` file. If you delete the shared `.kube-tkg/config` file, all management clusters become orphaned and thus unusable.

## <a id="networking"></a> Management Cluster Networking

When you deploy a management cluster, pod-to-pod networking with [Antrea](https://antrea.io/) is automatically enabled in the management cluster.

## <a id="dhcp"></a> Configure DHCP Reservations for the Control Plane Nodes (vSphere Only)

After you deploy a cluster to vSphere, each control plane node requires a static IP address. This includes both management and Tanzu Kubernetes clusters. These static IP addresses are required in addition to the static IP address that you assigned to Kube-VIP when you deploy a management cluster.

To make the IP addresses that your DHCP server assigned to the control plane nodes static, you can configure a DHCP reservation for each control plane node in the cluster. For instructions on how to configure DHCP reservations, see your DHCP server documentation.

## <a id="verify-deployment"></a>Verify the Deployment of the Management Cluster

After the deployment of the management cluster completes successfully, you can obtain information about your management cluster by:

* Locating the management cluster objects in vSphere, Amazon EC2, or Azure
* Using the Tanzu CLI and `kubectl`

### <a id="infrastructure"></a>View Management Cluster Objects in vSphere, Amazon EC2, or Azure

To view the management cluster objects in vSphere, Amazon EC2, or Azure, do the following:

   - If you deployed the management cluster to vSphere, go to the resource pool that you designated when you deployed the management cluster.
   - If you deployed the management cluster to Amazon EC2, go to the **Instances** view of your EC2 dashboard.
   - If you deployed the management cluster to Azure, go to the resource group that you designated when you deployed the management cluster.

   You should see the following VMs or instances.

   - **vSphere**:
       - One or three control plane VMs, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-sx5rp`
       - A worker node VM with a name similar to `CLUSTER-NAME-md-0-6b8db6b59d-kbnk4`
   - **Amazon EC2**:
       - One or three control plane VM instances, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-bcpfp`
       - A worker node instance with a name similar to `CLUSTER-NAME-md-0-dwfnm`
       - An EC2 bastion host instance with the name `CLUSTER-NAME-bastion`
   - **Azure**:
       - One or three control plane VMs, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-rh7xv`
       - A worker node VMs with a name similar to `CLUSTER-NAME-md-0-rh7xv`
       - Disk and Network Interface resources for the control plane and worker node VMs, with names based on the same name patterns.

   If you did not specify a name for the management cluster, `CLUSTER-NAME` is something similar to `tkg-mgmt-vsphere-20200323121503` or `tkg-mgmt-aws-20200323140554`.

### <a id="cli"></a> View Management Cluster Details With Tanzu CLI and `kubectl`

Tanzu CLI provides commands that facilitate many of the operations that you can perform with your management cluster. However, for certain operations, you still need to use `kubectl`.

Tanzu Kubernetes Grid provides two access levels for every management cluster and Tanzu Kubernetes cluster:

- The `admin` context of a cluster gives you full access to that cluster.
    - If you implemented identity management on the cluster, using the `admin` context allows you to run `kubectl` operations without requiring authentication with your identity provider (IDP).
    - If you did not implement identity management on the management cluster, you must use the `admin` context to run `kubectl` operations.
- If you implemented identity management on the cluster, using the regular context requires you to authenticate with your IDP before you can run `kubectl` operations on the cluster.

When you deploy a management cluster, the `kubectl` context is not automatically set to the management cluster.

Before you can run `kubectl` operations on a management cluster, you must obtain its `kubeconfig`.

1. On the bootstrap machine, run the `tanzu login` command to see the available management clusters and which one is the current login context for the CLI.

   For more information, see [List Management Clusters and Change Context](../cluster-lifecycle/multiple-management-clusters.md#login).

1. To see the details of the management cluster, run `tanzu management-cluster get`.  

   For more information, see [See Management Cluster Details](../cluster-lifecycle/multiple-management-clusters.md#list-mc).

1. To retrieve a `kubeconfig` for the management cluster, run the `tanzu management-cluster kubeconfig get` command as described in [Retrieve Management Cluster `kubeconfig`](#kubeconfig).

1. Set the context of `kubectl` to the management cluster.

   ```
   kubectl config use-context my-mgmnt-cluster-admin@my-mgmnt-cluster
   ```

1. Use `kubectl` commands to examine the resources of the management cluster.

   For example, run `kubectl get nodes`, `kubectl get pods`, or `kubectl get namespaces` to see the nodes, pods, and namespaces running in the management cluster.

## <a id="kubeconfig"></a> Retrieve Management Cluster `kubeconfig`

The `tanzu management-cluster kubeconfig get` command retrieves `kubeconfig` configuration information for the current management cluster, with options as follows:

- `--export-file FILE`
  - **Without option**: Add the retrieved cluster configuration information to the `kubectl` CLI's current `kubeconfig` file, whether it is the default `~/.kube/config` or set by the `KUBECONFIG` environment variable.
  - **With option**: Write the cluster configuration to a standalone `kubeconfig` file `FILE` that you can share with others.

- `--admin`
  - **Without option**: Generate a _regular `kubeconfig`_ that requires the user to authenticate with an external identity provider, and grants them access to cluster resources based on their assigned roles. To generate a regular  `kubeconfig`, identity management must be configured on the cluster.
      - The context name for this `kubeconfig` includes a `tanzu-cli-` prefix. For example, `tanzu-cli-id-mgmt-test@id-mgmt-test`.
  - **With option**: Generate an _administrator `kubeconfig`_ containing embedded credentials that lets the user access the cluster without logging in to an identity provider, and grants full access to the cluster's resources. If identity management is not configured on the cluster, you must specify the `--admin` option.
      - The context name for this `kubeconfig` includes an `-admin` suffix. For example, `id-mgmt-test-admin@id-mgmt-test`.

For example, to generate a standalone `kubeconfig` file to share with someone to grant them full access to your current management cluster:

   ```
   tanzu management-cluster kubeconfig get --admin --export-file MC-ADMIN-KUBECONFIG
   ```

To retrieve a `kubeconfig` for a workload cluster, run `tanzu cluster kubeconfig get` as described in [Retrieve Tanzu Kubernetes Cluster `kubeconfig`](../cluster-lifecycle/connect.md#kubeconfig).

## <a id="what-next"></a> What to Do Next

You can now use Tanzu Kubernetes Grid to start deploying Tanzu Kubernetes clusters. For information, see [Deploying Tanzu Kubernetes Clusters](../tanzu-k8s-clusters/index.md).

If you need to deploy more than one management cluster, on any or all of vSphere, Azure, and Amazon EC2, see [Manage Your Management Clusters](../cluster-lifecycle/multiple-management-clusters.md). This topic also provides information about how to add existing management clusters to your CLI instance, obtain credentials, scale and delete management clusters, and how to opt in or out of the CEIP.
