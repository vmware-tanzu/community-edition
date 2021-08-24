# Scale Tanzu Kubernetes Clusters

This topic explains how to scale a Tanzu Kubernetes cluster in three ways:

* **Autoscale**: Enable Cluster Autoscaler, which scales the number of worker nodes. See
[Scale Worker Nodes with Cluster Autoscaler](#enable-autoscaler) below.

* **Scale horizontally**: Run the `tanzu cluster scale` command with the `--controlplane-machine-count` and `--worker-machine-count` options, which scale the number of control plane and worker nodes.
See [Scale Clusters Horizontally With the Tanzu CLI](#horizontal-cli) below.

* **Scale vertically**: Change the cluster's machine template to increase the size of the control plane and worker nodes.
See [Scale Clusters Vertically With kubectl](#vertical-kubectl) below.

## <a id="enable-autoscaler"></a> Scale Worker Nodes with Cluster Autoscaler

To enable [Cluster Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) in a Tanzu Kubernetes cluster, you set the `AUTOSCALER_` options in the configuration file that you use to deploy the cluster or as environment variables before running `tanzu cluster create --file`. For information about the default configuration of Cluster Autoscaler in Tanzu Kubernetes Grid and the Cluster Autoscaler options, see [Cluster Autoscaler](../tanzu-config-reference.md#autoscaler) in the *Tanzu CLI Configuration File Variable Reference*.

1. Set the following configuration parameters.

   * For clusters with a single machine deployment such as `dev` clusters on vSphere, Amazon EC2, or Azure and `prod` clusters on vSphere or Azure, set `AUTOSCALER_MIN_SIZE_0` and `AUTOSCALER_MAX_SIZE_0`.
   * For clusters with multiple machine deployments such as `prod` clusters on Amazon EC2, set:
     * `AUTOSCALER_MIN_SIZE_0` and `AUTOSCALER_MAX_SIZE_0`
     * `AUTOSCALER_MIN_SIZE_1` and `AUTOSCALER_MAX_SIZE_1`
     * `AUTOSCALER_MIN_SIZE_2` and `AUTOSCALER_MAX_SIZE_2`

   You cannot modify these values after you deploy the cluster.

    ```
    #! ---------------------------------------------------------------------
    #! Autoscaler related configuration
    #! ---------------------------------------------------------------------
    ENABLE_AUTOSCALER: false
    AUTOSCALER_MAX_NODES_TOTAL: "0"
    AUTOSCALER_SCALE_DOWN_DELAY_AFTER_ADD: "10m"
    AUTOSCALER_SCALE_DOWN_DELAY_AFTER_DELETE: "10s"
    AUTOSCALER_SCALE_DOWN_DELAY_AFTER_FAILURE: "3m"
    AUTOSCALER_SCALE_DOWN_UNNEEDED_TIME: "10m"
    AUTOSCALER_MAX_NODE_PROVISION_TIME: "15m"
    AUTOSCALER_MIN_SIZE_0:
    AUTOSCALER_MAX_SIZE_0:
    AUTOSCALER_MIN_SIZE_1:
    AUTOSCALER_MAX_SIZE_1:
    AUTOSCALER_MIN_SIZE_2:
    AUTOSCALER_MAX_SIZE_2:
    ```

1. Create the cluster. For example:

   <pre>
   tanzu cluster create example-cluster
   </pre>

For each Tanzu Kubernetes cluster that you create with Autoscaler enabled, Tanzu Kubernetes Grid creates a Cluster Autoscaler deployment in the management cluster. To disable Cluster Autoscaler, delete the Cluster Autoscaler deployment associated with your Tanzu Kubernetes cluster.

## <a id="horizontal"></a> Scale a Cluster Horizontally With the Tanzu CLI

To horizontally scale a Tanzu Kubernetes cluster, use the `tanzu cluster scale` command.
You change the number of control plane nodes by specifying the `--controlplane-machine-count` option.
You change the number of worker nodes by specifying the `--worker-machine-count` option.

**NOTE**: On clusters that run in vSphere with Tanzu, you can only run either 1 control plane node or 3 control plane nodes. You can scale up the number of control plane nodes from 1 to 3, but you cannot scale down the number from 3 to 1.

- To scale a cluster that you originally deployed with 3 control plane nodes and 5 worker nodes to 5 and 10 nodes respectively, run the following command:

    <pre>tanzu cluster scale <em>cluster_name</em> --controlplane-machine-count 5 --worker-machine-count 10</pre>

   If you initially deployed a cluster with `--controlplane-machine-count 1` and then you scale it up to 3 control plane nodes, Tanzu Kubernetes Grid automatically enables stacked HA on the control plane.

- If the cluster in running in a namespace other than the `default` namespace, you must specify the `--namespace` option to scale that cluster.

    <pre>tanzu cluster scale <em>cluster_name</em> --controlplane-machine-count 5 --worker-machine-count 10 --namespace=my-namespace</pre>

**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu Kubernetes Grid operations are running.

## <a id="vertical-kubectl"></a> Scale a Cluster Vertically With kubectl

To vertically scale a Tanzu Kubernetes cluster, follow the [Changing Infrastructure Machine Templates](https://cluster-api.sigs.k8s.io/tasks/change-machine-template.html) procedure in _The Cluster API Book_, which changes the cluster's machine template.

The procedure downloads the cluster's existing machine template, with a `kubectl get` command that you can construct as follows:

  ```
  kubectl get MACHINE-TEMPLATE-TYPE MACHINE-TEMPLATE-NAME -o yaml
  ```

  Where:

  - `MACHINE-TEMPLATE-TYPE` is:
     - `VsphereMachineTemplate` on vSphere
     - `AWSMachineTemplate` on Amazon EC2
     - `AzureMachineTemplate` on Azure
  - `MACHINE-TEMPLATE-NAME` is the name of the machine template for the cluster nodes that you are scaling, which follows the form:
     - `CLUSTER-NAME-control-plane` for control plane nodes
     - `CLUSTER-NAME-worker` for worker nodes

For example:

  ```
  kubectl get VsphereMachineTemplate monitoring-cluster-worker -o yaml
  ```
