## Create Standalone Azure Clusters

This section covers setting up a standalone cluster in Azure. A standalone cluster provides a workload cluster that is **not** managed by a centralized management cluster.

1. Initialize the Tanzu kickstart UI.

    ```sh
    tanzu standalone-cluster create --ui
    ```

1. Go through the configuration steps, considering the following.


   * In Management Cluster Settings, use the Instance type drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs. The minimum configuration is 2 CPUs and 8 GB memory

   * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments. For more information about enabling Identity Management, see [Identity Management](latest/azure-install-mgmt.md/#step-5-identity-management).

1. At the end of the UI, deploy the cluster.

1. Store the name of your cluster (set during configuration or automatically generated) to a
   `WORKLOAD_CLUSTER_NAME` environment variable.

    ```sh
    export WORKLOAD_CLUSTER_NAME="<INSERT_WORKLOAD_CLUSTER_NAME_HERE>"
    ```

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context ${WORKLOAD_CLUSTER_NAME}-admin@${WORKLOAD_CLUSTER_NAME}
    ```

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get nodes

    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```
