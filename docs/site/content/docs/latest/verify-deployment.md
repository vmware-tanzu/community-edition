# Examine the Management Cluster Deployment

During the deployment of the management cluster, either from the installer interface or the CLI, Tanzu Community Edition creates a temporary management cluster using a [Kubernetes in Docker](https://kind.sigs.k8s.io/), `kind`, cluster on the bootstrap machine. Then, Tanzu Community Edition uses it to provision the final management cluster on the platform of your choice, depending on whether you are deploying to vSphere or Amazon EC2. After the deployment of the management cluster finishes successfully, Tanzu Community Edition deletes the temporary `kind` cluster.

1. Run the following command to verify that management cluster started successfully. If you did not specify a name for the management cluster, it will be something similar to `tkg-mgmt-vsphere-20200323121503` or `tkg-mgmt-aws-20200323140554`.
<!--add content for docker here -what will docker file name be-->
```sh
tanzu management-cluster get
```

2. Examine the folder structure. When Tanzu creates a management cluster for the first time, it creates a folder `~/.config/tanzu/tkg/providers` that contains all of the files required by Cluster API to create the management cluster.
The Tanzu installer interface saves the settings for the management cluster that it creates into a cluster configuration file `~/.config/tanzu/tkg/clusterconfigs/UNIQUE-ID.yaml`, where `UNIQUE-ID` is a generated filename.

3. To view the management cluster objects in vSphere, or Amazon EC2, do the following:
   * If you deployed the management cluster to vSphere, go to the resource pool that you designated when you deployed the management cluster. You should see:

      * One or three control plane VMs, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-sx5rp`
      * A worker node VM with a name similar to `CLUSTER-NAME-md-0-6b8db6b59d-kbnk4`
   * If you deployed the management cluster to Amazon EC2, go to the **Instances** view of your EC2 dashboard. You should see the following VMs or instances.
      * One or three control plane VM instances, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-bcpfp`
      * A worker node instance with a name similar to `CLUSTER-NAME-md-0-dwfnm`
      * An EC2 bastion host instance with the name `CLUSTER-NAME-bastion`

## Verify the Deployment of the Management Cluster

After the deployment of the management cluster completes successfully, you can obtain information about your management cluster by:

* Locating the management cluster objects in vSphere or Amazon EC2
* Using the Tanzu CLI and `kubectl`

### View Management Cluster Objects in vSphere or Amazon EC2

To view the management cluster objects in vSphere or Amazon EC2, do the following:

   - If you deployed the management cluster to vSphere, go to the resource pool that you designated when you deployed the management cluster.
   - If you deployed the management cluster to Amazon EC2, go to the **Instances** view of your EC2 dashboard.

   You should see the following VMs or instances.

   - **vSphere**:
       - One or three control plane VMs, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-sx5rp`
       - A worker node VM with a name similar to `CLUSTER-NAME-md-0-6b8db6b59d-kbnk4`
   - **Amazon EC2**:
       - One or three control plane VM instances, for development or production control plane, respectively, with names similar to `CLUSTER-NAME-control-plane-bcpfp`
       - A worker node instance with a name similar to `CLUSTER-NAME-md-0-dwfnm`
       - An EC2 bastion host instance with the name `CLUSTER-NAME-bastion`

   If you did not specify a name for the management cluster, `CLUSTER-NAME` is something similar to `tkg-mgmt-vsphere-20200323121503` or `tkg-mgmt-aws-20200323140554`.

### View Management Cluster Details With Tanzu CLI and `kubectl`

Tanzu CLI provides commands that facilitate many of the operations that you can perform with your management cluster. However, for certain operations, you still need to use `kubectl`.

 Tanzu Community Edition provides two access level contexts for every management and workload cluster:

- `admin` context gives you full access to a cluster.
    - If you implemented identity management on the cluster, using  `admin` context allows you to run `kubectl` operations without requiring authentication with your identity provider (IDP).
    - If you did not implement identity management on the management cluster, you must use the `admin` context to run `kubectl` operations.
- `regular` If you implemented identity management on the cluster, using the regular context, you must authenticate with your IDP before you can run `kubectl` operations on the cluster.

When you deploy a management cluster, the `kubectl` context is not automatically set to the context of the management cluster.

Before you can run `kubectl` operations on a management cluster, you must obtain its `kubeconfig` and set the context to the management cluster.

1. On the bootstrap machine, run the `tanzu login` command to see the available management clusters and which one is the current login context for the CLI.

1. To see the details of the management cluster, run `tanzu management-cluster get`.
1. To retrieve the `kubeconfig` for the management cluster, run the `tanzu management-cluster kubeconfig get` command with the following options:
   - `--export-file FILE` <br>
       - **Without option**: Add the retrieved cluster configuration information to the `kubectl` CLI's current `kubeconfig` file, whether it is the default `~/.kube/config` or set by the `KUBECONFIG` environment variable.
       - **With option**: Write the cluster configuration to a standalone `kubeconfig` file `FILE` that you can share with others.
   - `--admin`
       - **Without option**: Generate a _regular `kubeconfig`_ that requires the user to authenticate with an external identity provider, and grants them access to cluster resources based on their assigned roles. To generate a regular  `kubeconfig`, identity management must be configured on the cluster.<br>
       The context name for this `kubeconfig` includes a `tanzu-cli-` prefix. For example, `tanzu-cli-id-mgmt-test@id-mgmt-test`.
       - **With option**: Generate an _administrator `kubeconfig`_ containing embedded credentials that lets the user access the cluster without logging in to an identity provider, and grants full access to the cluster's resources. If identity management is not configured on the cluster, you must specify the `--admin` option. <br>
       The context name for this `kubeconfig` includes an `-admin` suffix. For example, `id-mgmt-test-admin@id-mgmt-test`.<br>
   For example, to generate a standalone `kubeconfig` file to share with someone to grant them full access to your current management cluster:

   ```sh
   tanzu management-cluster kubeconfig get --admin --export-file MC-ADMIN-KUBECONFIG
   ```
1. Set the context of `kubectl` to the management cluster.

   ```sh
   kubectl config use-context <MGMT-CLUSTER>-admin@<MGMT-CLUSTER>
   ```
   where ``<MGMT-CLUSTER>`` is the name of the management cluster
1. Use `kubectl` commands to examine the resources of the management cluster.

   For example, run `kubectl get nodes`, `kubectl get pods`, or `kubectl get namespaces` to see the nodes, pods, and namespaces running in the management cluster.


## Configure DHCP Reservations for the Control Plane Nodes (vSphere Only)

After you deploy a cluster to vSphere, each control plane node requires a static IP address. This includes both management and workload clusters. These static IP addresses are required in addition to the static IP address that you assigned to Kube-VIP when you deploy a management cluster.

To make the IP addresses that your DHCP server assigned to the control plane nodes static, you can configure a DHCP reservation for each control plane node in the cluster. For instructions on how to configure DHCP reservations, see your DHCP server documentation.

## Management Cluster Networking

When you deploy a management cluster, pod-to-pod networking with [Antrea](https://antrea.io/) is automatically enabled in the management cluster.