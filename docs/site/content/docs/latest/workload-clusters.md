# Deploy workload cluster

After you have deployed a management cluster to vSphere, or Amazon EC2, you can use the Tanzu CLI to deploy workload clusters.

To deploy a workload cluster, you create a configuration file. You then run the `tanzu cluster create` command, specifying the configuration file in the `--file` option.


## Before You Begin

-  When you deploy a workload cluster, most of the configuration for the cluster is the same as the configuration of the management cluster. The easiest way to obtain an initial configuration file for a workload cluster is to make a copy of the management cluster configuration file and to update it. Locate the YAML configuration file for the management cluster
    - If you deployed the management cluster from the installer interface, the configuration file is here: ``~/.tanzu/tkg/clusterconfigs/<MGMT-CLUSTER-NAME>.yaml``
  where ``<MGMT-CLUSTER-NAME>`` is the either the name you specified in the installer interface or else the randomly generated name if you didn’t specify a name

    - If you deployed the management cluster from the Tanzu CLI, the configuration file is in the default location (`~/.tanzu/tkg/cluster-config.yaml`) or in the location you specified in the –file parameter.
- **vSphere**: If you are deploying workload clusters to vSphere, each cluster requires one static virtual IP address to provide a stable endpoint for Kubernetes. Make sure that this IP address is not in the DHCP range, but is in the same subnet as the DHCP range.
- Configure workload cluster node size depending on cluster complexity and expected demand.
For more information, see [Minimum VM Sizes for Cluster Nodes](../mgmt-clusters/vsphere.md#vsphere-vm-sizes).

## Procedure


1. Make a copy of the management cluster configuration file and save it with a new name. 


<!--The simplest way to deploy a workload cluster is to specify a configuration that is identical to that of the management cluster. In this case, you only need to specify a name for the cluster. If you are deploying the cluster to vSphere, you must also specify an IP address or FQDN for the Kubernetes API endpoint.-->

1. Open the new YAML cluster configuration file in a text editor.
1. Optionally set a name for the cluster in the `CLUSTER_NAME` variable. If you do not specify a `CLUSTER_NAME` value in the cluster configuration file you must pass it as the first argument in the `tanzu cluster create` command.    The `CLUSTER_NAME` value passed to `tanzu cluster create` command overrides the `CLUSTER_NAME in the configuration file. Workload cluster names must be must be 42 characters or less, and must comply with DNS hostname requirements as amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).

1. vSphere: Specify a static virtual IP address or FQDN in the `VSPHERE_CONTROL_PLANE_ENDPOINT` variable.

   No two clusters, including any management cluster and workload cluster, can have the same `VSPHERE_CONTROL_PLANE_ENDPOINT` address.

  - Ensure that this IP address is not in the DHCP range, but is in the same subnet as the DHCP range.
  - If you mapped a fully qualified domain name (FQDN) to the VIP address, you can specify the FQDN instead of the VIP address.
   
   ```
   VSPHERE_CONTROL_PLANE_ENDPOINT: 10.90.110.100
   ```

1. Optional: You can update the `CLUSTER_PLAN` variable in the configuration to deploy a workload cluster that uses the `prod` plan, even if the management cluster was deployed with the `dev` plan, and the reverse.

1. To deploy a workload cluster with more control plane nodes than the `dev` and `prod` plans define by default, specify the `CONTROL_PLANE_MACHINE_COUNT` variable in the cluster configuration file. The number of control plane nodes that you specify in `CONTROL_PLANE_MACHINE_COUNT` must be uneven.

1. Specify the number of worker nodes for the cluster in the `WORKER_MACHINE_COUNT` variable.

1. Deploy a workload cluster in a specific Namespace

If you have created namespaces, you can deploy workload clusters to those namespaces by specifying the `NAMESPACE` variable. If you do not specify the `NAMESPACE` variable, Tanzu Community Edition places clusters in the `default` namespace. Any namespace that you identify in the `NAMESPACE` variable must exist in the management cluster before you run the command. For example, you might want to create different types of clusters in dedicated namespaces. For information about creating namespaces in the management cluster, see [Create Namespaces in the Management Cluster](../cluster-lifecycle/multiple-management-clusters.md#create-namespaces). **NOTE**: If you have created namespaces, you must provide a unique name for all workload cluster across all namespaces. If you provide a cluster name that is in use in another namespace in the same instance, the deployment fails with an error.


1. You can configure proxies, Machine Health Check, private registries, and Antrea on workload clusters in the same way as you do for management clusters. For information, see [Create a Management Cluster Configuration File](../mgmt-clusters/create-config-file.md).
1. To configure a workload cluster to use an OS other than the default Ubuntu v20.0.4, you must set the `OS_NAME` and `OS_VERSION` values in the cluster configuration file. The installer interface does not include node VM OS values in the management cluster configuration files that it saves to `~/.tanzu/tkg/clusterconfigs`.
1. Run the following command to deploy the workload cluster:

tanzu cluster create <WORKLOAD-CLUSTER-NAME> --file <CONFIG-FILE>
where
<WORKLOAD-CLUSTER-NAME> if you don't specify this the name is taken from the config file
<CONFIG-FILE> is the duplicated file you created in the previous steps

Any name that you specify in the `tanzu cluster create` command will override the name you set in the configuration file.
1. To see information about the cluster, run the `tanzu cluster get` command, specifying the cluster name.

   ```
   tanzu cluster get <WORKLOAD-CLUSTER>
   ```

   The output lists information about the status of the control plane and worker nodes, the Kubernetes version that the cluster is running, and the names of the nodes.

    ```   
    NAME             NAMESPACE  STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
    my-vsphere-tkc   default    running  1/1           1/1      v1.20.5+vmware.1  <none>

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

- A  **Development** workload cluster consists of the following VMs or instances:

    - vSphere:
        - One control plane node, with a name similar to `my-dev-cluster-control-plane-nj4z6`.
        - One worker node, with a name similar to `my-dev-cluster-md-0-6ff9f5cffb-jhcrh`. 
    - Amazon EC2: 
        - One control plane node, with a name similar to `my-dev-cluster-control-plane-d78t5`.
        - One EC2 bastion node, with the name  `my-dev-cluster-bastion`.
        - One worker node, with a name similar to `my-dev-cluster-md-0-2vsr4`. 

- A **Production** workload cluster consists of the following VMs or instances:

    - vSphere
        - Three control plane nodes, with names similar to `my-prod-cluster-control-plane-nj4z6`.
        - Three worker nodes, with names similar to `my-prod-cluster-md-0-6ff9f5cffb-jhcrh`.
    - Amazon EC2: 
        - Three control plane nodes, with names similar to `my-prod-cluster-control-plane-d78t5`.
        - One EC2 bastion node, with the name  `my-prod-cluster-bastion`.
        - Three worker nodes, with names similar to `my-prod-cluster-md-0-2vsr4`.














## <a id="manifest"></a> Deploy a Cluster from a Saved Manifest File

You can use the Tanzu CLI to create cluster manifest files for workload clusters without actually creating the clusters.
1. To generate a cluster manifest YAML file that you can pass to `kubectl apply -f`, run the `tanzu cluster create` command with the `--dry-run` option and save the output to a file. Use the same options and configuration `--file` that you would use if you were creating the cluster, for example:

```
tanzu cluster create my-cluster --file my-cluster-config.yaml --dry-run > my-cluster-manifest.yaml
```

2. To deploy a cluster from the saved manifest file, pass it to the `kubectl apply -f` command.
For example:

```
kubectl config use-context my-mgmt-context-admin@my-mgmt-context
```
```
kubectl apply -f my-cluster-manifest.yaml
```





