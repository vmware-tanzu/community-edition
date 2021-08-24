# Access the Tanzu Kubernetes Grid Logs

Tanzu Kubernetes Grid retains logs for management cluster and Tanzu Kubernetes cluster deployment and operation.

## <a id="mgmt-logs"></a> Access Management Cluster Deployment Logs

To monitor and troubleshoot management cluster deployments, review:

* The log file listed in the terminal output **Logs of the command execution can also be found at...**

* The log from your cloud provider module for Cluster API. Retrieve the most recent one as follows:
    1. Search your `tanzu management-cluster create` output for **Bootstrapper created. Kubeconfig:** and copy the `kubeconfig` file path listed. The file is in `~/.kube-tkg/tmp/`.
    1. Run the following, based on your cloud provider:
        * **vSphere**: `kubectl logs deployment.apps/capv-controller-manager -n capv-system manager --kubeconfig </path/to/kubeconfig>`
        * **Amazon EC2**: `kubectl logs deployment.apps/capa-controller-manager -n capa-system manager --kubeconfig </path/to/kubeconfig>`
        * **Azure**: `kubectl logs deployment.apps/capz-controller-manager -n capz-system manager --kubeconfig </path/to/kubeconfig>`

## <a id="workload-logs"></a> Monitor Tanzu Kubernetes Cluster Deployments in Cluster API Logs

After running `tanzu cluster create`,
you can monitor the deployment process in the Cluster API logs on the management cluster.

To access these logs, follow the steps below:

1. Set `kubeconfig` to your management cluster. For example:

    ```
    kubectl config use-context my-management-cluster-admin@my-management-cluster
    ```

1. Run the following:

    * `capi` logs:

       ```
       kubectl logs deployments/capi-controller-manager -n capi-system manager
       ```

    * IaaS-specific logs:

        * **vSphere**: `kubectl logs deployments/capv-controller-manager -n capv-system manager`
        * **Amazon EC2**: `kubectl logs deployments/capa-controller-manager -n capa-system manager`
        * **Azure**: `kubectl logs deployments/capz-controller-manager -n capz-system manager`
