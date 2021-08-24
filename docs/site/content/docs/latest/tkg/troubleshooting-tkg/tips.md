# Troubleshooting Tips for Tanzu Kubernetes Grid

This section includes tips to help you to troubleshoot common problems that you might encounter when installing Tanzu Kubernetes Grid and deploying Tanzu Kubernetes clusters.

Many of these procedures use the `kind` CLI on your bootstrap machine. To install `kind`, see [Installation](https://kind.sigs.k8s.io/docs/user/quick-start#installation) in the `kind` documentation.

## <a id="cleanup"></a> Clean Up After an Unsuccessful Management Cluster Deployment

**Problem**

An unsuccessful attempt to deploy a Tanzu Kubernetes Grid management cluster leaves orphaned objects in your cloud infrastructure and on your bootstrap machine.

**Solution**

1. Monitor your `tanzu management-cluster create` command output either in the terminal or Tanzu Kubernetes Grid installer interface. If the command fails, it prints a help message that includes the following: "**Failure while deploying management cluster... To clean up the resources created by the management cluster: tkg delete mc...**."
1. Run `tanzu management-cluster delete YOUR-CLUSTER-NAME`. This command removes the objects that it created in your infrastructure and locally.

You can also use the alternative methods described below:

* Bootstrap machine cleanup:

   * To remove a `kind` cluster, use the `kind` CLI. For example:

     ```
     kind get clusters
     kind delete cluster --name tkg-kind-example1234567abcdef
     ```

   * To remove Docker objects, use the `docker` CLI. For example, `docker rm`, `docker rmi`, and `docker system prune`.

     **CAUTION**: If you are running Docker processes that are not related to Tanzu Kubernetes Grid on your system, remove unneeded Docker objects individually.

* Infrastructure provider cleanup:

    * **vSphere**: Locate, power off, and delete the VMs and other resources that were created by Tanzu Kubernetes Grid.
    * **AWS**: Log in to your Amazon EC2 dashboard and delete the resources manually or use an automated solution.
    * **Azure**: In **Resource Groups**, open your `AZURE_RESOURCE_GROUP`. Use checkboxes to select and **Delete** the resources that were created by Tanzu Kubernetes Grid, which contain a timestamp in their names.

## <a id="kubectl-delete"></a> Delete Users, Contexts, and Clusters with `kubectl`

To clean up your `kubectl` state by deleting some or all of its users, contexts, and clusters:

1. Open your `~/.kube/config` and `~/.kube-tkg/config` files.

1. For the `user` objects that you want to delete, run:

    ```
    kubectl config unset users.USER-NAME
    kubectl config unset users.USER-NAME --kubeconfig ~/.kube-tkg/config
    ```

    Where `USER-NAME` is the `name` property of each top-level `user` object, as listed in the `config` files.

1. For the `context` objects that you want to delete, run:

    ```
    kubectl config unset contexts.CONTEXT-NAME
    kubectl config unset contexts.CONTEXT-NAME --kubeconfig ~/.kube-tkg/config
    ```

    Where `CONTEXT-NAME` is the `name` property of each top-level `context` object, as listed in the `config` files, typically of the form `contexts.mycontext-admin@mycontext`.

1. For the `cluster` objects that you want to delete, run:

    ```
    kubectl config unset clusters.CLUSTER-NAME
    kubectl config unset clusters.CLUSTER-NAME --kubeconfig ~/.kube-tkg/config
    ```

    Where `CLUSTER-NAME` is the `name` property of each top-level `cluster` object, as listed in the `config` files.

1. If the `config` files list the current context as a cluster that you deleted, unset the context:

    ```
    kubectl config unset current-context
    kubectl config unset current-context --kubeconfig ~/.kube-tkg/config
    ```

1. If you deleted management clusters that are tracked by the `tanzu` CLI, delete them from the `tanzu` CLI's state by running `tanzu config server delete` as described in [Delete Management Clusters from Your Tanzu CLI Configuration](../cluster-lifecycle/multiple-management-clusters.md#delete-mc-config).
  - To see the management clusters that the `tanzu` CLI tracks, run `tanzu login`.

## <a id="remove-kind"></a> Kind Cluster Remains after Deleting Management Cluster

**Problem**

Running `tanzu management-cluster delete` removes the management cluster, but fails to delete the local `kind` cluster from the bootstrap machine.

**Solution**

1. List all running `kind` clusters and remove the one that looks like `tkg-kind-unique_ID`

   ```
   kind delete cluster --name tkg-kind-unique_ID
   ```

1. List all running clusters and identify the `kind` cluster.

   ```
   docker ps -a
   ```

1. Copy the container ID of the `kind` cluster and remove it.

   ```
   docker kill container_ID
   ```

## <a id="aws-credentials"></a> Failed Validation, Credentials Error on Amazon EC2

**Problem**

Running `tanzu management-cluster create` fails with an error similar to the following:

```
Validating the pre-requisites...
Looking for AWS credentials in the default credentials provider chain

Error: : Tkg configuration validation failed: failed to get AWS client: NoCredentialProviders: no valid providers in chain
caused by: EnvAccessKeyNotFound: AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY not found in environment
SharedCredsLoad: failed to load shared credentials file
caused by: FailedRead: unable to open file
caused by: open /root/.aws/credentials: no such file or directory
EC2RoleRequestError: no EC2 instance role found
caused by: EC2MetadataError: failed to make EC2Metadata request
```

**Solution**

Tanzu Kubernetes Grid uses the default AWS credentials provider chain. Before creating a management or a workload cluster on Amazon EC2, you must configure your AWS account credentials as described in [Configure AWS Credentials](../mgmt-clusters/aws.md#account-setup).

## <a id="azure-license"></a> Failed Validation, Legal Terms Error on Azure

Before creating a management or workload cluster on Azure, you must accept the legal terms that cover the VM image used by cluster nodes.
Running `tanzu management-cluster create` or `tanzu cluster create` without having accepted the license fails with an error like:

```
User failed validation to purchase resources. Error message: 'You have not accepted the legal terms on this subscription: '*********' for this plan. Before the subscription can be used, you need to accept the legal terms of the image.
```

If this happens, accept the legal terms and try again:

* **Management Cluster**: See [Accept the Base Image License](../mgmt-clusters/azure.md#license)
* **Workload Cluster**: See the **Azure** instructions in [Deploy a Cluster with a Non-Default Kubernetes Version](../tanzu-k8s-clusters/k8s-versions.md#non-default).

## <a id="cluster-timeout"></a> Deploying a Tanzu Kubernetes Cluster Times Out, but the Cluster Is Created

**Problem**

Running `tanzu cluster create` fails with a timeout error similar to the following:

```
I0317 11:11:16.658433 clusterclient.go:341] Waiting for resource my-cluster of type *v1alpha3.Cluster to be up and running
E0317 11:26:16.932833 common.go:29]
Error: unable to wait for cluster and get the cluster kubeconfig: error waiting for cluster to be provisioned (this may take a few minutes): cluster control plane is still being initialized
E0317 11:26:16.933251 common.go:33]
Detailed log about the failure can be found at: /var/folders/_9/qrf26vd5629_y5vgxc1vjk440000gp/T/tkg-20200317T111108811762517.log
```

However, if you run `tanzu cluster list`, the cluster appears to have been created.

```
---------------------------------+

NAME          STATUS
---------------------------------+

my-cluster    Provisioned
---------------------------------+
```

**Solution**

1. Use the `tanzu cluster kubeconfig get CLUSTER-NAME --admin` command to add the cluster credentials to your `kubeconfig`.

   ```
   tanzu cluster kubeconfig get my-cluster --admin
   ```

1. Set `kubectl` to the cluster's context.

   ```
   kubectl config set-context my-cluster@user
   ```

1. Check whether the cluster nodes are all in the ready state.

   ```
   kubectl get nodes
   ```

1. Check whether all of the pods are up and running.

   ```
   kubectl get pods -A
   ```

1. If all of the nodes and pods are running correctly, your Tanzu Kubernetes cluster has been created successfully and you can ignore the error.
1. If the nodes and pods are not running correctly, attempt to delete the cluster.

   ```
   tanzu cluster delete my-cluster
   ```

1. If `tanzu cluster delete` fails, use `kubectl` to delete the cluster manually, as described in [Delete Users, Contexts, and Clusters with `kubectl`](#kubectl-delete).

## <a id="pods-pending"></a>Pods Are Stuck in Pending on Cluster Due to vCenter Connectivity

**Problem**

When you run `kubectl get pods -A` on the created cluster, some pods remain in pending.

You run `kubectl describe pod -n pod-namespace pod-name` on an affected pod and review
events and see the following event:

```
n node(s) had taint {node.cloudprovider.kubernetes.io/uninitialized: true}, that the pod didn't tolerate
```

**Solution**

Ensure there is connectivity and firewall rules in place to ensure communication between the cluster and vCenter.
For firewall ports and protocols requirements, see [Ports and Protocols](../tkg-security.md#ports)

## <a id="windows-ui"></a> Tanzu Kubernetes Grid UI Does Not Display Correctly on Windows

**Problem**

When you run the `tanzu management-cluster create --ui` command on a Windows system, the UI opens in your default browser, but the graphics and styling are not applied. This happens because a Windows registry is set to `application/x-css`.

**Solution**

1. In Windows search, enter `regedit` to open the Registry Editor utility.
1. Expand `HKEY_CLASSES_ROOT` and select `.css`.
1. Right-click **Content Type** and select **Modify**.
1. Set the Value to `text/css` and click **OK**.
1. Run the `tanzu management-cluster create --ui` command again to relaunch the UI.

## <a id="macos-kubectl"></a> Running `tanzu management-cluster create` on macOS Results in `kubectl` Version Error

**Problem**

If you run the `tanzu management-cluster create` command on macOS with the latest stable version of Docker Desktop, `tanzu management-cluster create` fails with the error message:

```
Error: : kubectl prerequisites validation failed: kubectl client version v1.15.5 is less than minimum supported kubectl client version 1.17.0
```

This happens because Docker Desktop symlinks `kubectl` 1.15 into the path.

**Solution**

Place a newer supported version of `kubectl` in the path before Docker's version.

## <a id="connect-ssh"></a> Connect to Cluster Nodes with SSH

You can use SSH to connect to individual nodes of management clusters or Tanzu Kubernetes clusters. To do so, the SSH key pair that you created when you deployed the management cluster must be available on the machine on which you run the SSH command. Consquently, you must run `ssh` commands on the machine on which you run `tanzu` commands.

The SSH keys that you register with the management cluster, and consequently that are used by any Tanzu Kubernetes clusters that you deploy from the management cluster, are associated with the following user accounts:

- vSphere management cluster and Tanzu Kubernetes nodes running on both Photon OS and Ubuntu: `capv`
- Amazon EC2 bastion nodes: `ubuntu`
- Amazon EC2 management cluster and Tanzu Kubernetes nodes running on Ubuntu: `ubuntu`
- Amazon EC2 management cluster and Tanzu Kubernetes nodes running on Amazon Linux: `ec2-user`
- Azure management cluster and Tanzu Kubernetes nodes (always Ubuntu): `capi`

To connect to a node by using SSH, run one of the following commands from the machine that you use as the bootstrap machine:

- vSphere nodes: <code>ssh capv@<em>node_address</em></code>
- Amazon EC2 bastion nodes and management cluster and workload nodes on Ubuntu: <code>ssh ubuntu@<em>node_address</em></code>
- Amazon EC2 management cluster and Tanzu Kubernetes nodes running on Amazon Linux: <code>ssh ec2-user@<em>node_address</em></code>
- Azure nodes: <code>ssh capi@<em>node_address</em></code>

Because the SSH key is present on the system on which you are running the `ssh` command, no password is required.

## <a id="recover-mc-creds"></a> Recover Management Cluster Credentials

If you have lost the credentials for a management cluster, for example by inadvertently deleting the `.kube-tkg/config` file on the system on which you run `tanzu` commands, you can recover the credentials from the management cluster control plane node.

1. Run `tanzu management-cluster create` to recreate the `.kube-tkg/config` file.
1. Obtain the public IP address of the management cluster control plane node, from vSphere, Amazon EC2, or Azure.
1. Use SSH to log in to the management cluster control plane node.

   See [Connect to Cluster Nodes with SSH](#connect-ssh) above for the credentials to use for each infrastructure provider.

1. Access the `admin.conf` file for the management cluster.

   ```
   sudo vi /etc/kubernetes/admin.conf
   ```

   The `admin.conf` file contains the cluster name, the cluster user name, the cluster context, and the client certificate data.
1. Copy the cluster name, the cluster user name, the cluster context, and the client certificate data into the `.kube-tkg/config` file on the system on which you run `tanzu` commands.

## <a id="restore-dir"></a> Restore ~/.tanzu Directory

**Problem**

The `~/.tanzu` directory on the bootstrap machine has been accidentally deleted or corrupted.
The `tanzu` CLI creates and uses this directory, and cannot function without it.

**Solution**

To restore the contents of the `~/.tanzu` directory:

1. To identify existing Tanzu Kubernetes Grid management clusters, run:

  ```
  kubectl --kubeconfig ~/.kube-tkg/config config get-contexts
  ```

  The command output lists names and contexts of all management clusters created or added by the v1.2 `tkg` or v1.3 `tanzu` CLI.

1. For each management cluster listed in the output, restore it to the `~/.tanzu` directory and CLI by running:

  ```
  tanzu login --kubeconfig ~/.kube-tkg/config --context MGMT-CLUSTER-CONTEXT --name MGMT-CLUSTER
  ```

## <a id="nfs-utils"></a> Disable `nfs-utils` on Photon OS Nodes

**Problem**

In Tanzu Kubernetes Grid v1.1.2 and later, `nfs-utils` is enabled by default. If you do not require `nfs-utils`, you can remove it from cluster node VMs.

**Solution**

To disable `nfs-utils` on clusters that you deploy with Tanzu Kubernetes Grid v1.1.2 or later, [use SSH to log in to the cluster node VMs](#connect-ssh) and run the following command:

   ```
   tdnf erase nfs-utils
   ```

For information about using `nfs-utils` on clusters deployed with Tanzu Kubernetes Grid v1.0 or 1.1.0, see [Enable or Disable nfs-utils on Photon OS Nodes](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.1/vmware-tanzu-kubernetes-grid-11/GUID-troubleshooting-tkg-tips.html#nfs-utils)
in the VMware Tanzu Kubernetes Grid 1.1.x Documentation.

## <a id="nsxadvlb-no-route-to-hosts"></a> Requests to NSX Advanced Load Balancer VIP fail with the message no route to host

**Problem**

If the total number of `LoadBalancer` type Service is large, and if all of the Service Engines are deployed in the same L2 network, requests to the NSX Advanced Load Balancer VIP can fail with the message `no route to host`.

This occurs because the default ARP rate limit on Service Engines is 100.

**Solution**

Set the ARP rate limit to a larger number. This parameter is not tunable in NSX Advanced Load Balancer Essentials, but it is tunable in NSX Advanced Load Balancer Enterprise Edition.
