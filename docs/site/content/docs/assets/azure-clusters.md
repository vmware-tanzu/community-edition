## Create Microsoft Azure Clusters

This section describes setting up management and workload clusters for
Microsoft Azure.

1. Initialize the Tanzu kickstart UI.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Go through the installation process for Azure. With the following
   considerations:


   * In Management Cluster Settings, use the Instance type drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs. The minimum configuration is 2 CPUs and 8 GB memory

   * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments. For more information about enabling Identity Management, see [Identity Management](latest/azure-install-mgmt/#step-5-identity-management).

2. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get
    ```
3. Create a cluster name that will be used throughout this Getting Started guide. This instance of `MGMT_CLUSTER_NAME` should be set to whatever value is returned by `tanzu management-cluster get` in the previous step.

    ```sh
    export MGMT_CLUSTER_NAME="<INSERT_MGMT_CLUSTER_NAME_HERE>"
    export WORKLOAD_CLUSTER_NAME="<INSERT_WORKLOAD_CLUSTER_NAME_HERE>"
    ```

4. Capture the management cluster's kubeconfig.

    ```sh
    tanzu management-cluster kubeconfig get ${MGMT_CLUSTER_NAME} --admin

    Credentials of workload cluster 'mtce' have been saved
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

    > Note the context name `${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}`, you'll use the above command in
    > future steps.

5. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context ${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}
    ```

6. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes
    ```
7. Next you will create a workload cluster. First, setup a workload cluster config file.

    ```sh
    cp  ~/.tanzu/tkg/clusterconfigs/<MGMT-CONFIG-FILE> ~/.tanzu/tkg/clusterconfigs/workload1.yaml
    ```

   > ``<MGMT-CONFIG-FILE>`` is the name of the management cluster YAML config file

   > This step duplicates the configuration file that was created when you deployed your management cluster. The configuration file will either have the name you assigned to the management cluster, or if no name was assigned, it will be a randomly generated name.

   > This duplicated file will be used as the configuration file for your workload cluster. You can edit the parameters in this new  file as required. For an example of a workload cluster template, see  [Azure Workload Cluster Template](../azure-wl-template).

   [](ignored)

   > In the next two steps you will edit the parameters in this new file (`workload1`) and then use the file to deploy a workload cluster.

   [](ignored)

8. In the new workload cluster file (`~/.config/tanzu/tkg/clusterconfigs/workload1.yaml`), edit the CLUSTER_NAME parameter to assign a name to your workload cluster. For example,

   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-workload-cluster
   CLUSTER_PLAN: dev
   ```
   #### Note
   * If you did not specify a name for your management cluster, the installer generated a random unique name. In this case, you must manually add the CLUSTER_NAME parameter and assign a workload cluster name.
   * If you specified a name for your management cluster, the CLUSTER_NAME parameter is present and needs to be changed to the new workload cluster name.
   > The other parameters in ``workload1.yaml`` are likely fine as-is. However, you can change them as required. Reference an example configuration template here:  [Amazon EC3 Workload Cluster Template](../aws-wl-template).

   > Validation is performed on the file prior to applying it, so the `tanzu` command will return a message if something necessary is omitted.

9. Create your workload cluster.

    ```sh
    tanzu cluster create ${WORKLOAD_CLUSTER_NAME} --file ~/.tanzu/tkg/clusterconfigs/workload1.yaml
    ```

10. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

11. Capture the workload cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get ${WORKLOAD_CLUSTER_NAME} --admin
    ```

12. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context ${WORKLOAD_CLUSTER_NAME}-admin@${WORKLOAD_CLUSTER_NAME}
    ```

13. Verify you can see pods in the cluster.

    ```sh
    kubectl get pods --all-namespaces

    NAMESPACE     NAME                                                    READY   STATUS    RESTARTS   AGE
    kube-system   antrea-agent-9d4db                                      2/2     Running   0          3m42s
    kube-system   antrea-agent-vkgt4                                      2/2     Running   1          5m48s
    kube-system   antrea-controller-5d594c5cc7-vn5gt                      1/1     Running   0          5m49s
    kube-system   coredns-5d6f7c958-hs6vr                                 1/1     Running   0          5m49s
    kube-system   coredns-5d6f7c958-xf6cl                                 1/1     Running   0          5m49s
    kube-system   etcd-tce-guest-control-plane-b2wsf                      1/1     Running   0          5m56s
    kube-system   kube-apiserver-tce-guest-control-plane-b2wsf            1/1     Running   0          5m56s
    kube-system   kube-controller-manager-tce-guest-control-plane-b2wsf   1/1     Running   0          5m56s
    kube-system   kube-proxy-9825q                                        1/1     Running   0          5m48s
    kube-system   kube-proxy-wfktm                                        1/1     Running   0          3m42s
    kube-system   kube-scheduler-tce-guest-control-plane-b2wsf            1/1     Running   0          5m56s
    ```
