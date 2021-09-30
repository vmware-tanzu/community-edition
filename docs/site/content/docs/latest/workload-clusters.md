# Deploying a workload cluster

After you have deployed a management cluster to Amazon Web Services (AWS), Microsoft Azure, or vSphere you can use the Tanzu CLI to deploy workload clusters. In Tanzu Community Edition, your application workloads run on workload clusters.

Tanzu Community Edition automatically deploys workload clusters to the platform on which you deployed the management cluster. For example, you cannot deploy workload clusters to AWS from a management cluster that is running in vSphere. It is not possible to use shared services between the different providers because each provider uses different systems.

Tanzu Community Edition automatically deploys clusters from whichever management cluster you have set as the context for the CLI by using the `tanzu login` command.

To deploy a workload cluster, you create a configuration file. You then run the `tanzu cluster create` command, specifying the configuration file in the `--file` option. To see an example of a workload cluster configuration file template, see [AWS Workload Cluster Template](../aws-wl-template), [Azure Workload Cluster Template](../azure-wl-template), or  [vSphere Workload Cluster Template](../vsphere-wl-template).

For specific configuration parameters, see:

* [AWS Workload Cluster Template](aws-wl-template)
* [vSphere Workload Cluster Template](vsphere-wl-template)
* [Microsoft Azure Workload Cluster Template](azure-wl-template)

## Before You Begin

* Copy the configuration file: When you deploy a workload cluster, most of the configuration for the cluster is the same as the configuration of the management cluster. The easiest way to obtain an initial configuration file for a workload cluster is to make a copy of the management cluster configuration file and to update it. Locate the YAML configuration file for the management cluster.
  * If you deployed the management cluster from the installer interface, the configuration file is here: `~/.config/tanzu/tkg/clusterconfigs/<MGMT-CLUSTER-NAME>.yaml`
      where `<MGMT-CLUSTER-NAME>` is either the name you specified in the installer interface or else the randomly generated name if you didn’t specify a name.

  * If you deployed the management cluster from the Tanzu CLI, the configuration file is in the default location (`~/.config/tanzu/tkg/cluster-config.yaml`) or in the location you specified in the `--file` parameter.
* vSphere: If you are deploying workload clusters to vSphere, each cluster requires one static virtual IP address to provide a stable endpoint for Kubernetes. Make sure that this IP address is not in the DHCP range, but is in the same subnet as the DHCP range.
* Create namespaces: To help you to organize and manage your development projects, you can optionally divide the management cluster into Kubernetes namespaces. You can then use Tanzu CLI to deploy workload clusters to specific namespaces in your management cluster. For example, you might want to create different types of workload clusters in dedicated namespaces. If you do not create additional namespaces, Tanzu Community Edition creates all workload clusters in the `default` namespace. Complete the following steps:

    1. Run `kubectl config current-context` to make sure that `kubectl` is connected to the correct management cluster context.
    1. Run `kubectl get namespaces` to list the namespaces that are currently present in the management cluster.
    1. Use `kubectl create -foo` to create new namespaces, for example for development and production.
    1. Run `kubectl get namespaces --show-labels` to see the new namespaces.

## Deploying a Workload Cluster Procedure {#procedure}

1. Make a copy of the management cluster configuration file and save it with a new name.
1. Open the new YAML cluster configuration file in a text editor.
1. Set a name for the cluster in the `CLUSTER_NAME` variable. If you do not specify a `CLUSTER_NAME` value in the cluster configuration file you must pass it as the first argument in the `tanzu cluster create` command. The `CLUSTER_NAME` value passed to the `tanzu cluster create` command overrides the `CLUSTER_NAME` in the configuration file. Workload cluster names must be must be 42 characters or less, and must comply with DNS hostname requirements as amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).
1. vSphere: Specify a static virtual IP address or FQDN in the `VSPHERE_CONTROL_PLANE_ENDPOINT` variable. No two clusters, including any management cluster and workload cluster, can have the same `VSPHERE_CONTROL_PLANE_ENDPOINT` address. Ensure that this IP address is not in the DHCP range, but is in the same subnet as the DHCP range. If you mapped a fully qualified domain name (FQDN) to the VIP address, you can specify the FQDN instead of the VIP address.
1. Optional: Update the `CLUSTER_PLAN` variable to deploy a workload cluster that uses the `prod` plan, even if the management cluster was deployed with the `dev` plan, and the reverse.
1. Optional: To deploy a workload cluster with more control plane nodes than the `dev` and `prod` plans define by default, edit the `CONTROL_PLANE_MACHINE_COUNT` variable. The number of control plane nodes that you specify in `CONTROL_PLANE_MACHINE_COUNT` must be uneven.
1. Optional: Specify the number of worker nodes for the cluster in the `WORKER_MACHINE_COUNT` variable.
1. Optional: Deploy a workload cluster in a specific namespace.
If you have created namespaces, you can deploy workload clusters to those namespaces by specifying the `NAMESPACE` variable. If you do not specify the `NAMESPACE` variable, Tanzu Community Edition places clusters in the `default` namespace. Any namespace that you identify in the `NAMESPACE` variable must exist in the management cluster before you run the command. You must provide a unique name for all workload clusters across all namespaces. If you provide a workload cluster name that is already in use in another namespace in the same instance, the deployment fails with an error.
1. Optional: To configure a workload cluster to use an OS other than the default Ubuntu v21.04, you must set the `OS_NAME` and `OS_VERSION` values in the cluster configuration file. The installer interface does not include node VM OS values in the management cluster configuration files that it saves to `~/.config/tanzu/tkg/clusterconfigs`.
1. Run the following command to deploy the workload cluster:

   ```sh
   tanzu cluster create <WORKLOAD-CLUSTER-NAME> --file <CONFIG-FILE>
   ```

    where

   `<WORKLOAD-CLUSTER-NAME>` is the name you want to assign to the workload cluster, if you don't specify this the name is taken from the config file

   `<CONFIG-FILE>` is the duplicated file you created in the previous step. To see an example of a workload cluster template, see [AWS Workload Cluster Template](../aws-wl-template), [Azure Workload Cluster Template](../azure-wl-template), or  [vSphere Workload Cluster Template](../vsphere-wl-template).

## Viewing Information about the Workload Cluster {#view}

To see information about the cluster, run the `tanzu cluster get` command, specifying the cluster name.

```sh
tanzu cluster get <WORKLOAD-CLUSTER>
```

The output lists information about the status of the control plane and worker nodes, the Kubernetes version that the cluster is running, and the names of the nodes.

```sh
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

* A **Development** workload cluster consists of the following VMs or instances:
  * vSphere:
    * One control plane node, with a name similar to `my-dev-cluster-control-plane-nj4z6`.
    * One worker node, with a name similar to `my-dev-cluster-md-0-6ff9f5cffb-jhcrh`.
  * AWS:
    * One control plane node, with a name similar to `my-dev-cluster-control-plane-d78t5`.
    * One EC2 bastion node, with the name  `my-dev-cluster-bastion`.
    * One worker node, with a name similar to `my-dev-cluster-md-0-2vsr4`.
  * Azure:
    * One control plane node, with a name similar to `my-dev-cluster-20200902052434-control-plane-4d4p4`.
    * One worker node, with a name similar to `my-dev-cluster-20200827115645-md-0-rjdbr`.

* A **Production** workload cluster consists of the following VMs or instances:
  * vSphere
    * Three control plane nodes, with names similar to `my-prod-cluster-control-plane-nj4z6`.
    * Three worker nodes, with names similar to `my-prod-cluster-md-0-6ff9f5cffb-jhcrh`.
  * AWS:
    * Three control plane nodes, with names similar to `my-prod-cluster-control-plane-d78t5`.
    * One EC2 bastion node, with the name  `my-prod-cluster-bastion`.
    * Three worker nodes, with names similar to `my-prod-cluster-md-0-2vsr4`.
  * Azure:
    * Three control plane nodes, with names similar to `my-prod-cluster-20200902052434-control-plane-4d4p4`.
    * Three worker nodes, with names similar to `my-prod-cluster-20200827115645-md-0-rjdbr`.

## Set the Kubectl Context to the Workload Cluster {#context}

Tanzu Community Edition does not automatically set the kubectl context to a workload cluster when you create it. You must set the kubectl context to the workload cluster manually by using the `kubectl config use-context` command. Complete the following steps:

1. Capture the workload cluster’s kubeconfig.

   ```sh
   tanzu cluster kubeconfig get <WORKLOAD-CLUSTER-NAME> --admin
   ```

1. Set your kubectl context accordingly.

   ```sh
   kubectl config use-context <WORKLOAD-CLUSTER-NAME>-admin@<WORKLOAD-CLUSTER-NAME>
   ```

## Deploying a Workload Cluster from a Saved Manifest File {#manifest}

You can use the Tanzu CLI to create cluster manifest files for workload clusters without actually creating the clusters.

1. To generate a cluster manifest YAML file that you can pass to `kubectl apply -f`, run the `tanzu cluster create` command with the `--dry-run` option and save the output to a file. Use the same options and configuration `--file` that you would use if you were creating the cluster, for example:

   ```sh
   tanzu cluster create my-cluster --file my-cluster-config.yaml --dry-run > my-cluster-manifest.yaml
   ```

1. To deploy a cluster from the saved manifest file, pass it to the `kubectl apply -f` command.

   For example:

   ```sh
   kubectl config use-context my-mgmt-context-admin@my-mgmt-context
   ```

   ```sh
   kubectl apply -f my-cluster-manifest.yaml
   ```

## (Optional) Attach a workload cluster to Tanzu Mission Control

You can optionally attach a workload cluster to VMware Tanzu Mission Control, for more information, see [Attach a Cluster](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-6DF2CE3E-DD07-499B-BC5E-6B3B2E02A070.html?hWord=N4IghgNiBc4C5zAYwBYAIxqRArgZzgFMAnEAXyA) in the Tanzu Mission Control documentation.
