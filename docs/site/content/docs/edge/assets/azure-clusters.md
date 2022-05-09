## Create Microsoft Azure Clusters

This section describes setting up management and workload clusters for Microsoft Azure.

### Deploy a Management Cluster

There are some prerequisites this process will assume. Refer to the [Prepare to Deploy a Cluster to Azure](../azure-mgmt) docs for instructions on accepting image licenses and preparing your Azure account.

1. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu management-cluster create --ui
    ```

    Note: If you are bootstrapping from a Windows machine and encounter an `unable to ensure prerequisites` error, see the following  [troubleshooting topic](../faq-cluster-bootstrapping/#x509-certificate-signed-by-unknown-authority-when-deploying-management-cluster-from-windows).

1. Choose Azure from the provider tiles.

    ![kickstart azure tile](/docs/img/kickstart-azure-tile.png)

1. Fill out the IaaS Provider section.

    ![kickstart azure iaas](/docs/img/kickstart-azure-iaas.png)

    * `A`: Your account's Tenant ID.
    * `B`: Your Client ID.
    * `C`: Your Client secret.
    * `D`: Your Subscription ID.
    * `E`: The Azure Cloud in which to deploy. For example, "Public Cloud", "US
      Government Cloud", etc.
    * `F`: [The region of
      Azure](https://azure.microsoft.com/en-us/global-infrastructure/geographies/#geographies)
      you'd like all networking, compute, etc to be created within.
    * `G`: The public key you'd like to use for your VM instances. This is how
      you'll SSH into control plane and worker nodes.
    * `H`: Whether to use an existing
      [resource group](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/manage-resource-groups-portal#what-is-a-resource-group)
      or create a new one.
    * `I`: The existing resource group, or the name to provide the new resource group.

1. Fill out the VNET settings.

    ![kickstart azure vnet](/docs/img/kickstart-azure-vnet.png)

    * `A`: Whether to create a new
      [Virtual Network in Azure](https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-overview)
      or use an existing one.
      If using an existing one, you must provide its VNET name. For initial
      deployments, it is recommended to create a new Virtual Network. This will
      ensure the installer takes care of all networking creation and configuration.
    * `B`: The Resource Group under which to create the VNET.
    * `C`: The name to use when creating a new VNET.
    * `D`: The CIDR block to use for this VNET.
    * `E`: The name for the control plane subnet.
    * `F`: The CIDR block to use for the control plane subnet. This range
      should be within the VNET CIDR.
    * `G`: The name for the worker node subnet.
    * `H`: The CIDR block to use for the worker node subnet. This range should
      be within the VNET CIDR and not overlap with the control plane CIDR.
    * `I`: Whether to deploy without a publicly accessible IP address.
      Access to the cluster will be limited to your Azure private network only.
      Various ways for connecting to your private cluster
      [can be found in the Azure private cluster
      documentation](https://docs.microsoft.com/en-us/azure/aks/private-clusters#options-for-connecting-to-the-private-cluster).

1. Fill out the Management Cluster Settings.

    ![kickstart azure management cluster settings](/docs/img/kickstart-azure-mgmt-cluster.png)

    * `A`: Choose between Development profile with one control plane node, or
      Production, which features a highly-available three node control plane.
      Additionally, choose the instance type you'd like to use for control plane nodes.
    * `B`: Name the cluster. This is a friendly name that will be used to
      reference your cluster in the Tanzu CLI and `kubectl`.
    * `C`: The instance type to be used for each node creation. See the
      instances types documentation to understand trade-offs between CPU,
      memory, pricing and more.
    * `D`: Whether to enable [Cluster API's machine health
      checks](https://cluster-api.sigs.k8s.io/tasks/healthcheck.html).
    * `E`: Choose whether you'd like to enable [Kubernetes API server
      auditing](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/).

1. If you would like additional metadata to be tagged in your soon-to-be-created
   Azure infrastructure, fill out the Metadata section.

1. Fill out the Kubernetes Network section.

    ![kickstart kubernetes networking](/docs/img/kickstart-azure-network.png)

    * `A`: Set the CIDR for Kubernetes [Services (Cluster
      IPs)](https://kubernetes.io/docs/concepts/services-networking/service/).
These are internal IPs that, by default, are only exposed and routable within
Kubernetes.
    * `B`: Set the CIDR range for Kubernetes Pods. These are internal IPs that, by
      default, are only exposed and routable within Kubernetes.
    * `C`: Set a network proxy that internal traffic should egress through to
      access external network(s).

1. Fill out the Identity Management section.

    ![kickstart identity management](/docs/img/kickstart-identity.png)

    * `A`: Select whether you want to enable identity management. If this is
      off, certificates (via kubeconfig) are used to authenticate users. For
      most development scenarios, it is preferred to keep this off.
    * `B`: If identity management is on, choose whether to authenticate using
      [OIDC](https://openid.net/connect/) or [LDAPS](https://ldap.com).
    * `C`: Fill out connection details for identity management.

1. Fill out the OS Image section.

    ![kickstart azure os](/docs/img/kickstart-azure-os.png)

    * `A`: The Azure image to use for Kubernetes host VMs. This list should
      populate based on known images uploaded by VMware.
      These images are publicly accessible for your use. Choose
      based on your preferred Linux distribution.

1. Skip the TMC Registration section.

1. Click the Review Configuration button.

    > For your record, the configuration settings have been saved to
    > `${HOME}/.config/tanzu/tkg/clusterconfigs`.

1. Deploy the cluster.

    > If you experience issues deploying your cluster, visit the [Troubleshooting
    > documentation](../tanzu-diagnostics).

1. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get
    ```

    The output will look similar to the following:

    ```sh
    NAME         NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES       
    mgmt         tkg-system  running  1/1           1/1      v1.21.2+vmware.1  management

    Details:

    NAME                                                            READY  SEVERITY  REASON  SINCE  MESSAGE
    /mgmt                                                           True                     5m38s
    ├─ClusterInfrastructure - AzureCluster/mgmt                     True                     5m42s
    ├─ControlPlane - KubeadmControlPlane/mgmt-control-plane         True                     5m38s
    │ └─Machine/mgmt-control-plane-d99g5                            True                     5m41s
    └─Workers
      └─MachineDeployment/mgmt-md-0
        └─Machine/mgmt-md-0-bc94f54b4-tgr9h                         True                     5m41s

    Providers:

    NAMESPACE                          NAME                   TYPE                    PROVIDERNAME  VERSION  WATCHNAMESPACE
    capi-kubeadm-bootstrap-system      bootstrap-kubeadm      BootstrapProvider       kubeadm       v0.3.23
    capi-kubeadm-control-plane-system  control-plane-kubeadm  ControlPlaneProvider    kubeadm       v0.3.23
    capi-system                        cluster-api            CoreProvider            cluster-api   v0.3.23
    capz-system                        infrastructure-azure   InfrastructureProvider  azure         v0.4.15

    ```

1. Capture the management cluster's kubeconfig.

    ```sh
    tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin
    ```

    Where `<MGMT-CLUSTER-NAME>` should be set to the name returned by `tanzu management-cluster get` above.  <br><br>
    For example, if your management cluster is called 'mgmt', you will see a message similar to:

    ```sh
    Credentials of workload cluster 'mgmt' have been saved.
    You can now access the cluster by running 'kubectl config use-context mgmt-admin@mgmt'
    ```

1. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context <MGMT-CLUSTER-NAME>-admin@<MGMT-CLUSTER-NAME>
    ```

    Where `<MGMT-CLUSTER-NAME>` should be set to the name returned by `tanzu management-cluster get`.

1. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes

    NAME                       STATUS   ROLES                  AGE    VERSION
    mgmt-control-plane-vkpsm   Ready    control-plane,master   111m   v1.21.2+vmware.1
    mgmt-md-0-qbbhk            Ready    <none>                 109m   v1.21.2+vmware.1
    ```

### Deploy a Workload Cluster

1. Next you will create a workload cluster. First, setup a workload cluster configuration file.

    ```sh
    cp  ~/.config/tanzu/tkg/clusterconfigs/<MGMT-CONFIG-FILE> ~/.config/tanzu/tkg/clusterconfigs/workload1.yaml
    ```

   * Where ``<MGMT-CONFIG-FILE>`` is the name of the management cluster YAML configuration file

   * This step duplicates the configuration file that was created when you deployed your management cluster. The configuration file will either have the name you assigned to the management cluster, or if no name was assigned, it will be a randomly generated name.

   * This duplicated file will be used as the configuration file for your workload cluster. You can edit the parameters in this new  file as required. For an example of a workload cluster template, see  [Azure Workload Cluster Template](../azure-wl-template).

   * In the next two steps you will edit the parameters in this new file (`workload1`) and then use the file to deploy a workload cluster.

1. In the new workload cluster file (`~/.config/tanzu/tkg/clusterconfigs/workload1.yaml`),
   edit the `CLUSTER_NAME` parameter to assign a name to your workload cluster. For example,

   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-workload-cluster
   CLUSTER_PLAN: dev
   ```

   * If you did not specify a name for your management cluster, the installer generated a random unique name. In this case, you must manually add the CLUSTER_NAME parameter and assign a workload cluster name. The workload cluster names must be must be 42 characters or less and must comply with DNS hostname requirements as described here: [RFC 1123](https://tools.ietf.org/html/rfc1123)
   * If you specified a name for your management cluster, the CLUSTER_NAME parameter is present and needs to be changed to the new workload cluster name.
   * The other parameters in ``workload1.yaml`` are likely fine as-is. Validation is performed on the file prior to applying it, so the `tanzu` command will return a message if something necessary is omitted. However, you can change parameters as required. Reference an example configuration template here:  [Azure Workload Cluster Template](../azure-wl-template).
   * To deploy a workload cluster with a non-default version of Kubernetes, use the `--tkr` option. For more information, see [Deploy Clusters with Different Kubernetes Versions](../tkr-managed-cluster).

1. Create your workload cluster.

    ```sh
    tanzu cluster create <WORKLOAD-CLUSTER-NAME> --file ~/.config/tanzu/tkg/clusterconfigs/workload1.yaml
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
