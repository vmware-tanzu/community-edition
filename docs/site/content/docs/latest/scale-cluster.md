# Scale workloads

This topic explains how to scale a workload cluster in three ways:

* **Autoscale**: Enable Cluster Autoscaler, which scales the number of worker nodes. See
[Scale Worker Nodes with Cluster Autoscaler](#enable-autoscaler).

* **Scale horizontally**: Run the `tanzu cluster scale` command with the `--controlplane-machine-count` and `--worker-machine-count` options, which scale the number of control plane and worker nodes.
See [Scale Clusters Horizontally With the Tanzu CLI](#horizontal-cli).

* **Scale vertically**: Change the cluster's machine template to increase the size of the control plane and worker nodes.
See [Scale Clusters Vertically With kubectl](#vertical-kubectl).

## <a id="enable-autoscaler"></a> Scale Worker Nodes with Cluster Autoscaler

Cluster Autoscaler is a Kubernetes program that automatically scales Kubernetes clusters depending on the demands on the workload clusters. For more information about Cluster Autoscaler, see the following documentation in GitHub:

* [Cluster Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) introduction and documentation.
* [Frequently Asked Questions](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md) about Cluster Autoscaler and how it relates to alternative autoscaling approaches.
* [Cluster Autoscaler on Cluster API](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler/cloudprovider/clusterapi) for information about cluster-api provider implementation for cluster autoscaler.

By default, Cluster Autoscaler is disabled. To enable Cluster Autoscaler in a workload cluster, set the `ENABLE_AUTOSCALER` to `true` and set the `AUTOSCALER_` options in the cluster configuration file or as environment variables before running `tanzu cluster create --file`.

Each Cluster Autoscaler configuration variable in a cluster configuration file corresponds to a parameter in the Cluster Autoscaler tool.
For a list of these variables see [Cluster Autoscaler](../tanzu-config-reference.md#autoscaler) in the *Tanzu CLI Configuration File Variable Reference*.

The `AUTOSCALER_*_SIZE` settings limit the number of worker nodes in a cluster, while `AUTOSCALER_MAX_NODES_TOTAL` limits the count of all nodes, both worker and control plane.

Set `AUTOSCALER_*_SIZE` values depending the number of worker nodes in the cluster:

   * For clusters with a single worker node, such as `dev` clusters, set `AUTOSCALER_MIN_SIZE_0` and `AUTOSCALER_MAX_SIZE_0`.
   * For clusters with multiple worker nodes, such as `prod` clusters, set:
     * `AUTOSCALER_MIN_SIZE_0` and `AUTOSCALER_MAX_SIZE_0`
     * `AUTOSCALER_MIN_SIZE_1` and `AUTOSCALER_MAX_SIZE_1`
     * `AUTOSCALER_MIN_SIZE_2` and `AUTOSCALER_MAX_SIZE_2`

Below is an example of Cluster Autoscaler settings in a cluster configuration file.
You cannot modify these values after you deploy the cluster.

  ```sh
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

For each workload that you create with Cluster Autoscaler enabled, Tanzu creates a Cluster Autoscaler deployment in the management cluster. To disable Cluster Autoscaler, delete the Cluster Autoscaler deployment associated with your workload.

## <a id="horizontal"></a> Scale a Cluster Horizontally With the Tanzu CLI

To horizontally scale a workload, use the `tanzu cluster scale` command.
You change the number of control plane nodes by specifying the `--controlplane-machine-count` option.
You change the number of worker nodes by specifying the `--worker-machine-count` option.

- To scale a cluster that you originally deployed with 3 control plane nodes and 5 worker nodes to 5 and 10 nodes respectively, run the following command:

    ```sh
    tanzu cluster scale cluster_name --controlplane-machine-count 5 --worker-machine-count 10
    ```

   If you initially deployed a cluster with `--controlplane-machine-count 1` and then you scale it up to 3 control plane nodes, Tanzu Kubernetes Grid automatically enables stacked HA on the control plane.

- If the cluster is running in a namespace other than the `default` namespace, you must specify the `--namespace` option to scale that cluster.
    ```sh
    tanzu cluster scale <em>cluster_name --controlplane-machine-count 5 --worker-machine-count 10 --namespace=my-namespace
    ```

**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu operations are running.

## <a id="vertical-kubectl"></a> Scale a Cluster Vertically With kubectl

To vertically scale a workload, follow the [Updating Infrastructure Machine Templates](https://cluster-api.sigs.k8s.io/tasks/updating-machine-templates.html#updating-infrastructure-machine-templates) procedure in _The Cluster API Book_, which changes the cluster's machine template.

The procedure downloads the cluster's existing machine template, with a `kubectl get` command that you can construct as follows:

  ```sh
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
