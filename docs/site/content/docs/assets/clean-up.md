## Cleaning up

After going through this guide, the following enables you to clean-up resources.

1. Delete any deployed workload clusters.

    ```sh
    tanzu cluster delete <WORKLOAD-CLUSTER-NAME>
    ```

1. Once all workload clusters have been deleted, the management cluster can
   then be removed as well. Run the following commands to get the name of the cluster and delete the cluster

    ```sh
    tanzu management-cluster get
    tanzu management-cluster delete <MGMT-CLUSTER-NAME>
    ```
