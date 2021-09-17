## Create Microsoft Azure Clusters

This section describes setting up management and workload clusters for
Microsoft Azure.

1. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Complete the configuration steps in the installer interface for Azure and create the management cluster. The following configuration settings are recommended:

   * If you do not specify a name, the installer interface generates a unique name. If you do specify a name, the name must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements described here: [RFC 1123](https://tools.ietf.org/html/rfc1123).
   * In Management Cluster Settings, use the Instance type drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs. The minimum configuration is 2 CPUs and 8 GB memory

   * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments. For more information about enabling Identity Management, see [Identity Management](../azure-install-mgmt/#step-5-identity-management).

1. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get
    ```

    The output will look similar to the following:

    ```sh
    NAME            NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
    mgmtclusterone  tkg-system  running  1/1           1/1      v1.21.2+vmware.1  management

    Details:

    NAME                                                               READY  SEVERITY  REASON  SINCE  MESSAGE
    /mgmtclusterone                                                    True                     67s
    ├─ClusterInfrastructure - AzureCluster/mgmtclusterone              True                     69s
    ├─ControlPlane - KubeadmControlPlane/mgmtclusterone-control-plane  True                     67s
    │ └─Machine/mgmtclusterone-control-plane-4hszz                     True                     68s
    └─Workers
    └─MachineDeployment/mgmtclusterone-md-0
        └─Machine/mgmtclusterone-md-0-85b4bc7c6d-mbj7j                 True                     68s


    Providers:

    NAMESPACE                          NAME                   TYPE                    PROVIDERNAME  VERSION  WATCHNAMESPACE
    capi-kubeadm-bootstrap-system      bootstrap-kubeadm      BootstrapProvider       kubeadm       v0.3.22
    capi-kubeadm-control-plane-system  control-plane-kubeadm  ControlPlaneProvider    kubeadm       v0.3.22
    capi-system                        cluster-api            CoreProvider            cluster-api   v0.3.22
    capz-system                        infrastructure-azure   InfrastructureProvider  azure         v0.4.15
    ```

1. Capture the management cluster's kubeconfig.

    ```sh
    tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin
    ```

    Where `<MGMT-CLUSTER-NAME>` should be set to the name returned by `tanzu management-cluster get` above.  <br><br>
    For example, if your management cluster is called 'mtce', you will see a message similar to:

    ```sh
    Credentials of workload cluster 'mtce' have been saved.
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

1. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context <MGMT-CLUSTER-NAME>-admin@<MGMT-CLUSTER-NAME>
    ```

    Where `<MGMT-CLUSTER-NAME>` should be set to the name returned by `tanzu management-cluster get`.

1. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes

    NAME                                   STATUS   ROLES                  AGE     VERSION
    standalonedelete-control-plane-9ndzx   Ready    control-plane,master   3m36s   v1.21.2+vmware.1
    standalonedelete-md-0-7hvll            Ready    <none>                 113s    v1.21.2+vmware.1
    ```

1. Next you will create a workload cluster. First, setup a workload cluster configuration file.

    ```sh
    cp  ~/.tanzu/tkg/clusterconfigs/<MGMT-CONFIG-FILE> ~/.tanzu/tkg/clusterconfigs/workload1.yaml
    ```

   * Where ``<MGMT-CONFIG-FILE>`` is the name of the management cluster YAML configuration file

   * This step duplicates the configuration file that was created when you deployed your management cluster. The configuration file will either have the name you assigned to the management cluster, or if no name was assigned, it will be a randomly generated name.

   * This duplicated file will be used as the configuration file for your workload cluster. You can edit the parameters in this new  file as required. For an example of a workload cluster template, see  [Azure Workload Cluster Template](../azure-wl-template).

   * In the next two steps you will edit the parameters in this new file (`workload1`) and then use the file to deploy a workload cluster.

1. In the new workload cluster file (`~/.config/tanzu/tkg/clusterconfigs/workload1.yaml`), edit the CLUSTER_NAME parameter to assign a name to your workload cluster. For example,

   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-workload-cluster
   CLUSTER_PLAN: dev
   ```

   * If you did not specify a name for your management cluster, the installer generated a random unique name. In this case, you must manually add the CLUSTER_NAME parameter and assign a workload cluster name. The workload cluster names must be must be 42 characters or less and must comply with DNS hostname requirements as described here: [RFC 1123](https://tools.ietf.org/html/rfc1123)
   * If you specified a name for your management cluster, the CLUSTER_NAME parameter is present and needs to be changed to the new workload cluster name.
   * The other parameters in ``workload1.yaml`` are likely fine as-is. Validation is performed on the file prior to applying it, so the `tanzu` command will return a message if something necessary is omitted. However, you can change paramaters as required. Reference an example configuration template here:  [Azure Workload Cluster Template](../azure-wl-template).

1. Create your workload cluster.

    ```sh
    tanzu cluster create <WORKLOAD-CLUSTER-NAME> --file ~/.tanzu/tkg/clusterconfigs/workload1.yaml
    ```

1. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

1. Capture the workload cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get <WORKLOAD-CLUSTER-NAME> --admin
    ```

1. Set your `kubectl` context to the workload cluster.

    ```sh
    kubectl config use-context <WORKLOAD-CLUSTER_NAME>-admin@<WORKLOAD-CLUSTER-NAME>
    ```

1. Verify you can see pods in the cluster.

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
