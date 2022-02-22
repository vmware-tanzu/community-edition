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

    **Note for AWS**: If the cluster you are deleting is deployed on AWS, you must precede the delete command 
    with the region. For example,

    ```sh
    AWS_REGION=us-west-2 tanzu management-cluster delete my-mgmt-cluster
    ```

    For more information on deleting clusters, see [Delete Management Clusters](../delete-mgmt/), and [Delete Workload Clusters](../delete-cluster/).
