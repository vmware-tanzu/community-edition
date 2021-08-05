## Cleaning up

After going through this guide, the following enables you to clean-up resources.

1. Delete any deployed workload clusters.

    ```sh
    tanzu cluster delete ${WORKLOAD_CLUSTER_NAME}
    ```

1. Once all workload clusters have been deleted, the management cluster can
   then be removed as well.

    ```sh
    tanzu management-cluster get

    NAME                         NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
    tkg-mgmt-aws-20210226062452  tkg-system  running  1/1           1/1      v1.20.1+vmware.2  management
    ```

    ```sh
    tanzu management-cluster delete ${MGMT_CLUSTER_NAME}
    ```
