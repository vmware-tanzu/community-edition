## Create Standalone Azure Clusters

This section covers setting up a standalone cluster in Azure. A standalone cluster provides a workload cluster that is **not** managed by a centralized management cluster.

1. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu standalone-cluster create --ui
    ```

1. Complete the configuration steps in the installer interface and create the standalone cluster. The following configuration settings are recommended:

   * Use the Instance type drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs. The minimum configuration is 2 CPUs and 8 GB memory.

   * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments. For more information about enabling Identity Management, see [Identity Management](../azure-install-mgmt/#step-5-identity-management).

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context <STANDALONE-CLUSTER-NAME>-admin@<STANDALONE-CLUSTER-NAME>
    ```

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get nodes
    ```

    The output will look similar to the following:

    ```sh
    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```
