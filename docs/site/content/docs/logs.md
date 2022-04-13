# Access local and kubectl-based Logs

Tanzu Community Edition retains logs for management and workload cluster deployment and operation.

## Using Local Cluster Logs

Local cluster logs capture the management and workload cluster creation, upgrade, and deletion activity in the `~/.config/tanzu/tkg/logs/<CLUSTER-NAME>` file. These logs can be used to troubleshoot the cluster creation activity and other failures. Successfully deleting a cluster automatically deletes its logs file, but if the delete action fails, the logs file remain.

## Using kubectl-based Cluster Logs

If local cluster logs do not provide sufficient information, you can retrieve logs using `kubectl` as follows:

### Access Management Cluster Deployment Logs

To monitor and troubleshoot management cluster deployments, review:

* The log file listed in the terminal output **Logs of the command execution can also be found at...**

* The log from your cloud provider module for Cluster API. Retrieve the most recent one as follows:
    1. Search your `tanzu management-cluster create` output for **Bootstrapper created. Kubeconfig:** and copy the `kubeconfig` file path listed. The file is in `~/.kube-tkg/tmp/`.
    1. Run the following, based on your cloud provider:
        * **vSphere**: `kubectl logs deployment.apps/capv-controller-manager -n capv-system manager --kubeconfig <PATH-TO-KUBECONFIG>`
        * **Amazon Web Services (AWS)**: `kubectl logs deployment.apps/capa-controller-manager -n capa-system manager --kubeconfig <PATH-TO-KUBECONFIG>`
        * **Azure**: `kubectl logs deployment.apps/capz-controller-manager -n capz-system manager --kubeconfig <PATH-TO-KUBECONFIG>`

### Access Workload Deployments Logs

After running `tanzu cluster create`, you can monitor the deployment process in the Cluster API logs on the management cluster.

To access these logs, follow the steps below:

1. Set `kubeconfig` to your management cluster. For example:

    ```sh
    kubectl config use-context my-management-cluster-admin@my-management-cluster
    ```

1. Run the following:

    * `capi` logs:
  
       ```sh
       kubectl logs deployments/capi-controller-manager -n capi-system manager
       ```

    * IaaS-specific logs:

        * **vSphere**: `kubectl logs deployments/capv-controller-manager -n capv-system manager`
        * **AWS**: `kubectl logs deployments/capa-controller-manager -n capa-system manager`
        * **Azure**: `kubectl logs deployments/capz-controller-manager -n capz-system manager`
