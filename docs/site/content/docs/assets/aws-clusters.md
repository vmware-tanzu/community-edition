## Create Amazon EC2 Clusters

This section describes setting up management and workload/guest clusters for
Amazon EC2.

1. Initialize the Tanzu kickstart UI.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Go through the installation process for Amazon EC2. With the following
   considerations:


   * Check the "Automate creation of AWS CloudFormation Stack" box if you do not have an existing TKG CloudFormation stack. This stack is used to created IAM resources that TCE clusters use in Amazon EC2.
     You only need 1 TKG CloudFormation stack per AWS account. CloudFormation is global and not locked to a region. For more information, see [Required IAM resources](../ref-aws/#permissions)

   * Set the instance type size to m5.xlarge or larger for both the control plane node and worker node.

   * Disable OIDC configuration.

    > Until we have more TCE documentation, you can find the full TKG docs
    > [here](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-mgmt-clusters-deploy-management-clusters.html).
    > We will have more complete `tanzu` cluster bootstrapping documentation available here in the near future.
1. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get

    NAME  NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
    mtce  tkg-system  running  1/1           1/1      v1.20.1+vmware.2  management

    Details:

    NAME                                                     READY  SEVERITY  REASON  SINCE  MESSAGE
    /mtce                                                    True                     113m
    ├─ClusterInfrastructure - AWSCluster/mtce                True                     113m
    ├─ControlPlane - KubeadmControlPlane/mtce-control-plane  True                     113m
    │ └─Machine/mtce-control-plane-r7k52                     True                     113m
    └─Workers
      └─MachineDeployment/mtce-md-0
        └─Machine/mtce-md-0-fdfc9f766-6n6lc                  True                     113m

    Providers:

    NAMESPACE                          NAME                   TYPE                    PROVIDERNAME  VERSION  WATCHNAMESPACE
    capa-system                        infrastructure-aws     InfrastructureProvider  aws           v0.6.4
    capi-kubeadm-bootstrap-system      bootstrap-kubeadm      BootstrapProvider       kubeadm       v0.3.14
    capi-kubeadm-control-plane-system  control-plane-kubeadm  ControlPlaneProvider    kubeadm       v0.3.14
    capi-system                        cluster-api            CoreProvider            cluster-api   v0.3.14
    ```

1. Create a cluster names that will be used throughout this Getting Started guide. This instance of `MGMT_CLUSTER_NAME` should be set to whatever value is returned by `tanzu management-cluster get` above.

    ```sh
    export MGMT_CLUSTER_NAME="<INSERT_MGMT_CLUSTER_NAME_HERE>"
    export GUEST_CLUSTER_NAME="<INSERT_GUEST_CLUSTER_NAME_HERE>"
    ```

1. Capture the management cluster's kubeconfig.

    ```sh
    tanzu management-cluster kubeconfig get ${MGMT_CLUSTER_NAME} --admin

    Credentials of workload cluster 'mtce' have been saved
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

    > Note the context name `${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}`, you'll use the above command in
    > future steps.

1. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context ${MGMT_CLUSTER_NAME}-admin@${MGMT_CLUSTER_NAME}
    ```

1. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes

    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```

1. Next you will create a guest cluster. First, setup a guest cluster config file.

    ```sh
    cp  ~/.tanzu/tkg/clusterconfigs/xw6nt8jduy.yaml ~/.tanzu/tkg/clusterconfigs/guest1.yaml
    ```

   > This step duplicates the configuration file that was created when you deployed your management cluster. This duplicated file will be used as the configuration file for your guest cluster. You can edit the parameters in this new  file as required. For an example of a guest cluster template, see  [Amazon EC2 Guest Cluster Template](../aws-wl-template).

   [](ignored)

   > In the next two steps you will edit the parameters in this new file (`guest1`) and then use the file to deploy a guest cluster.

   [](ignored)

1. In the new guest cluster file (`~/.tanzu/tkg/clusterconfigs/guest1.yaml`), edit the CLUSTER_NAME parameter to assign a name to your guest cluster. For example,

   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-guest-cluster
   CLUSTER_PLAN: dev
   ```
   #### Note
   * If you did not specify a name for your management cluster, the installer generated a unique name, in this case, you must manually add the CLUSTER_NAME parameter and assign a guest cluster name.
   * If you specified a name for your management cluster, the CLUSTER_NAME parameter is present and needs to be changed to the new guest cluster name.
   > The other parameters in ``guest1.yaml`` are likely fine as-is. However, you can change them as required. Reference an example configuration template here:  [Amazon EC3 Guest Cluster Template](../aws-wl-template).

   > Validation is performed on the file prior to applying it, so the `tanzu` command will return a message if something necessary is omitted.

1. Create your guest cluster.

    ```sh
    tanzu cluster create ${GUEST_CLUSTER_NAME} --file ${HOME}/.tanzu/tkg/clusterconfigs/guest1.yaml
    ```

1. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

1. Capture the guest cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get ${GUEST_CLUSTER_NAME} --admin
    ```

1. Set your `kubectl` context accordingly.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
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
