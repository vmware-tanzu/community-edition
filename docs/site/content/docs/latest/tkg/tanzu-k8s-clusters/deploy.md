# Deploy Tanzu Kubernetes Clusters

After you have deployed a management cluster to vSphere, Amazon EC2, or Azure, or you have connected the Tanzu CLI to a vSphere with Tanzu Supervisor Cluster, you can use the Tanzu CLI to deploy Tanzu Kubernetes clusters.

To deploy a Tanzu Kubernetes cluster, you create a configuration file that specifies the different options with which to deploy the cluster. You then run the `tanzu cluster create` command, specifying the configuration file in the `--file` option.

This topic describes the most basic configuration options for Tanzu Kubernetes clusters.

## <a id="prereqs"></a> Prerequisites for Cluster Deployment

- You have followed the procedures in [Install the Tanzu CLI and Other Tools](../install-cli.md) and [Deploying Management Clusters](../mgmt-clusters/deploy-management-clusters.md) to deploy a management cluster to vSphere, Amazon EC2, or Azure.
- You have already upgraded the management cluster to the version that corresponds with the Tanzu CLI version. If you attempt to deploy a Tanzu Kubernetes cluster with an updated CLI without upgrading the management cluster first, the Tanzu CLI returns the error `Error: validation failed: version mismatch between management cluster and cli version. Please upgrade your management cluster to the latest to continue.` For instructions on how to upgrade management clusters, see [Upgrade Management Clusters](../upgrade-tkg/management-cluster.md).
- Alternatively, you have a vSphere 7 instance on which a vSphere with Tanzu Supervisor Cluster is running. To deploy clusters to a vSphere 7 instance on which the vSphere with Tanzu feature is enabled, you must connect the Tanzu CLI to the vSphere with Tanzu Supervisor Cluster. For information about how to do this, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](connect-vsphere7.md).
- **vSphere**: If you are deploying Tanzu Kubernetes clusters to vSphere, each cluster requires one static virtual IP address to provide a stable endpoint for Kubernetes. Make sure that this IP address is not in the DHCP range, but is in the same subnet as the DHCP range.
- **Azure**: If you are deploying Tanzu Kubernetes clusters to Azure, each cluster requires a Network Security Group (NSG) for its worker nodes named `CLUSTER-NAME-node-nsg`, where `CLUSTER-NAME` is the name of the cluster. For more information, see [Network Security Groups on Azure](../mgmt-clusters/azure.md#nsgs).
- Configure Tanzu Kubernetes cluster node size depending on cluster complexity and expected demand.
For more information, see [Minimum VM Sizes for Cluster Nodes](../mgmt-clusters/vsphere.md#vsphere-vm-sizes).

## <a id="config"></a> Create a Tanzu Kubernetes Cluster Configuration File

When you deploy a Tanzu Kubernetes cluster, most of the configuration for the cluster is the same as the configuration of the management cluster that you use to deploy it. Because most of the configuration is the same, the easiest way to obtain an initial configuration file for a Tanzu Kubernetes cluster is to make a copy of the management cluster configuration file and to update it.

1. Locate the YAML configuration file for the management cluster.

  - If you deployed the management cluster from the installer interface and you did not specify the `--file` option when you ran `tanzu management-cluster create --ui`, the configuration file is saved in `~/.tanzu/tkg/clusterconfigs/`. The file has a randomly generated name, for example, `bm8xk9bv1v.yaml`.
  - If you deployed the management cluster from the installer interface and you did specify the `--file` option, the management cluster configuration is taken from in the file that you specified.
  - If you deployed the management cluster from the Tanzu CLI without using the installer interface, the management cluster configuration is taken from either a file that you specified in the `--file` option, or from the default location, `~/.tanzu/tkg/cluster-config.yaml`.
  
1. Make a copy of the management cluster configuration file and save it with a new name.

   For example, save the file as `my-aws-tkc.yaml`, `my-azure-tkc.yaml` or `my-vsphere-tkc.yaml`.

**IMPORTANT**: The recommended practice is to use a dedicated configuration file for every Tanzu Kubernetes cluster that you deploy.

## <a id="deploy"></a> Deploy a Tanzu Kubernetes Cluster with Minimum Configuration

The simplest way to deploy a Tanzu Kubernetes cluster is to specify a configuration that is identical to that of the management cluster. In this case, you only need to specify a name for the cluster. If you are deploying the cluster to vSphere, you must also specify an IP address or FQDN for the Kubernetes API endpoint.

**Note**: To configure a workload cluster to use an OS other than the default Ubuntu 20.04, you must set the `OS_NAME` and `OS_VERSION` values in the cluster configuration file.
The installer interface does not include node VM OS values in the management cluster configuration files that it saves to `~/.tanzu/tkg/clusterconfigs`.

1. Open the new YAML cluster configuration file in a text editor.
1. Optionally set a name for the cluster in the `CLUSTER_NAME` variable.

   For example, if you are deploying the cluster to vSphere, set the name to `my-vsphere-tkc`.

   ```
   CLUSTER_NAME: my-vsphere-tkc
   ```

   If you do not specify a `CLUSTER_NAME` value in the cluster configuration file or as an environment variable,
   you must pass it as the first argument to the `tanzu cluster create` command.
   The `CLUSTER_NAME` value passed to `tanzu cluster create` overrides the name you set in the configuration file.<br />
   Workload cluster names must be must be 42 characters or less, and must comply with DNS hostname requirements as amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).

1. If you are deploying the cluster to vSphere, specify a static virtual IP address or FQDN in the `VSPHERE_CONTROL_PLANE_ENDPOINT` variable.

   No two clusters, including any management cluster and workload cluster, can have the same `VSPHERE_CONTROL_PLANE_ENDPOINT` address.

  - Ensure that this IP address is not in the DHCP range, but is in the same subnet as the DHCP range.
  - If you mapped a fully qualified domain name (FQDN) to the VIP address, you can specify the FQDN instead of the VIP address.

   ```
   VSPHERE_CONTROL_PLANE_ENDPOINT: 10.90.110.100
   ```

1. Save the configuration file.
1. Run the `tanzu cluster create` command, specifying the path to the configuration file in the `--file` option.

   If you saved the Tanzu Kubernetes cluster configuration file `my-vsphere-tkc.yaml` in the default `clusterconfigs` folder, run the following command to create a cluster with a name that you specified in the configuration file:

   ```
   tanzu cluster create --file .tanzu/tkg/clusterconfigs/my-vsphere-tkc.yaml
   ```

   If you did not specify a name in the configuration file, or to create a cluster with a different name to the one that you specified, specify the cluster name in the `tanzu cluster create` command. For example, to create a cluster named `another-vsphere-tkc` from the configuration file `my-vsphere-tkc.yaml`, run the following command:

   ```
   tanzu cluster create another-vsphere-tkc --file .tanzu/tkg/clusterconfigs/my-vsphere-tkc.yaml
   ```

   Any name that you specify in the `tanzu cluster create` command will override the name you set in the configuration file.

1. To see information about the cluster, run the `tanzu cluster get` command, specifying the cluster name.

   ```
   tanzu cluster get my-vsphere-tkc
   ```

   The output lists information about the status of the control plane and worker nodes, the Kubernetes version that the cluster is running, and the names of the nodes.

    ```
    NAME             NAMESPACE  STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
    my-vsphere-tkc   default    running  1/1           1/1      v1.20.5+vmware.2  <none>

    Details:

    NAME                                                                READY  SEVERITY  REASON  SINCE  MESSAGE
    /my-vsphere-tkc                                                    True                     17m
    ├─ClusterInfrastructure - VSphereCluster/my-vsphere-tkc            True                     19m
    ├─ControlPlane - KubeadmControlPlane/my-vsphere-tkc-control-plane  True                     17m
    │ └─Machine/my-vsphere-tkc-control-plane-ss9rt                     True                     17m
    └─Workers
      └─MachineDeployment/my-vsphere-tkc-md-0
        └─Machine/my-vsphere-tkc-md-0-657958d58-mgtpp                  True                     8m33s
    
    ```

The cluster runs the default version of Kubernetes for this Tanzu Kubernetes Grid release, which in Tanzu Kubernetes Grid v1.3.1 is v1.20.5.

## <a id="deploy-counts"></a> Deploy a Cluster with Different Numbers of Control Plane and Worker Nodes

In the preceding example, because you did not change any of the node settings in the Tanzu Kubernetes cluster configuration file, the resulting Tanzu Kubernetes cluster has the same numbers of control plane and worker nodes as the management cluster. The nodes have the same CPU, memory, and disk configuration as the management cluster nodes.

- If you selected **Development** in the **Management Cluster Settings** section of the installer interface, or specified `CLUSTER_PLAN: dev` and the default numbers of nodes in the management cluster configuration, the Tanzu Kubernetes cluster consists of the following VMs or instances:

    - vSphere:
        - One control plane node, with a name similar to `my-dev-cluster-control-plane-nj4z6`.
        - One worker node, with a name similar to `my-dev-cluster-md-0-6ff9f5cffb-jhcrh`.
    - Amazon EC2:
        - One control plane node, with a name similar to `my-dev-cluster-control-plane-d78t5`.
        - One EC2 bastion node, with the name  `my-dev-cluster-bastion`.
        - One worker node, with a name similar to `my-dev-cluster-md-0-2vsr4`.
    - Azure:
        - One control plane node, with a name similar to `my-dev-cluster-20200902052434-control-plane-4d4p4`.
        - One worker node, with a name similar to `my-dev-cluster-20200827115645-md-0-rjdbr`.
- If you selected **Production** in the **Management Cluster Settings** section of the installer interface, or specified `CLUSTER_PLAN: prod` and the default numbers of nodes in the management cluster configuration file, Tanzu CLI deploys a cluster with three control plane nodes and automatically implements stacked control plane HA for the cluster. The Tanzu Kubernetes cluster consists of the following VMs or instances:

    - vSphere
        - Three control plane nodes, with names similar to `my-prod-cluster-control-plane-nj4z6`.
        - Three worker nodes, with names similar to `my-prod-cluster-md-0-6ff9f5cffb-jhcrh`.
    - Amazon EC2:
        - Three control plane nodes, with names similar to `my-prod-cluster-control-plane-d78t5`.
        - One EC2 bastion node, with the name  `my-prod-cluster-bastion`.
        - Three worker nodes, with names similar to `my-prod-cluster-md-0-2vsr4`.
    - Azure:
        - Three control plane nodes, with names similar to `my-prod-cluster-20200902052434-control-plane-4d4p4`.
        - Three worker nodes, with names similar to `my-prod-cluster-20200827115645-md-0-rjdbr`.

If you copied the cluster configuration from the management cluster, you can update the `CLUSTER_PLAN` variable in the configuration file to deploy a Tanzu Kubernetes cluster that uses the `prod` plan, even if the management cluster was deployed with the `dev` plan, and the reverse.

```
CLUSTER_PLAN: prod
```

To deploy a Tanzu Kubernetes cluster with more control plane nodes than the `dev` and `prod` plans define by default, specify the `CONTROL_PLANE_MACHINE_COUNT` variable in the cluster configuration file. The number of control plane nodes that you specify in `CONTROL_PLANE_MACHINE_COUNT` must be uneven.

```
CONTROL_PLANE_MACHINE_COUNT: 5
```

Specify the number of worker nodes for the cluster in the `WORKER_MACHINE_COUNT` variable.

```
WORKER_MACHINE_COUNT: 10
```

How you configure the size and resource configurations of the nodes depends on whether you are deploying clusters to vSphere, Amazon EC2, or Azure. For information about how to configure the nodes, see the appropriate topic for each provider:

- [Deploy Tanzu Kubernetes Clusters to vSphere](vsphere.md)
- [Deploy Tanzu Kubernetes Clusters to Amazon EC2](aws.md)
- [Deploy Tanzu Kubernetes Clusters to Azure](azure.md).

## <a id="common-config"></a> Configure Common Settings

You configure proxies, Machine Health Check, private registries, and Antrea  on Tanzu Kubernetes Clusters in the same way as you do for management clusters. For information, see [Create a Management Cluster Configuration File](../mgmt-clusters/create-config-file.md).

## <a id="namespace"></a> Deploy a Cluster in a Specific Namespace

If you have created namespaces in your Tanzu Kubernetes Grid instance, you can deploy Tanzu Kubernetes clusters to those namespaces by specifying the `NAMESPACE` variable. If you do not specify the the `NAMESPACE` variable, Tanzu Kubernetes Grid places clusters in the `default` namespace. Any namespace that you identify in the `NAMESPACE` variable must exist in the management cluster before you run the command. For example, you might want to create different types of clusters in dedicated namespaces. For information about creating namespaces in the management cluster, see [Create Namespaces in the Management Cluster](../cluster-lifecycle/multiple-management-clusters.md#create-namespaces).

```
NAMESPACE: production
```

**NOTE**: If you have created namespaces, you must provide a unique name for all Tanzu Kubernetes clusters across all namespaces. If you provide a cluster name that is in use in another namespace in the same instance, the deployment fails with an error.

## <a id="manifest"></a> Create Tanzu Kubernetes Cluster Manifest Files

You can use the Tanzu CLI to create cluster manifest files for Tanzu Kubernetes clusters without actually creating the clusters.
To generate a cluster manifest YAML file that you can pass to `kubectl apply -f`, run the `tanzu cluster create` command with the `--dry-run` option and save the output to a file. Use the same options and configuration `--file` that you would use if you were creating the cluster, for example:

```
tanzu cluster create my-cluster --file my-cluster-config.yaml --dry-run > my-cluster-manifest.yaml
```

### <a id="deploy-manifest"></a> Deploy a Cluster from a Saved Manifest File

To deploy a cluster from the saved manifest file, pass it to the `kubectl apply -f` command.
For example:

```
kubectl config use-context my-mgmt-context-admin@my-mgmt-context
```

```
kubectl apply -f my-cluster-manifest.yaml
```

## <a id="advanced"></a> Advanced Configuration of Tanzu Kubernetes Clusters

If you need to deploy a Tanzu Kubernetes cluster with more advanced configuration, rather than copying the configuration file of the management cluster, see the topics that describe the options that are specific to each infrastructure provider.

- [Deploy Tanzu Kubernetes Clusters to vSphere](vsphere.md)
- [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](connect-vsphere7.md)
- [Deploy Tanzu Kubernetes Clusters to Amazon EC2](aws.md)
- [Deploy Tanzu Kubernetes Clusters to Azure](azure.md)

Each of the topics on deployment to vSphere, Amazon EC2, and Azure include Tanzu Kubernetes cluster templates, that contain all of the options that you can use for each provider.

You can further customize the configuration of your Tanzu Kubernetes clusters by performing the following types of operations:

- [Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions](k8s-versions.md)
- [Customize Tanzu Kubernetes Cluster Networking](networking.md)
- [Create Persistent Volumes with Storage Classes](storage.md)

## What to Do Next

After you have deployed Tanzu Kubernetes clusters, the Tanzu CLI provides commands and options to perform the following cluster lifecycle management operations. See [Managing Cluster Lifecycles](../cluster-lifecycle/index.md).
