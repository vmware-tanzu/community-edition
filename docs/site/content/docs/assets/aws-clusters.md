## Create Amazon EC2 Clusters

This section describes setting up management and workload clusters for
Amazon EC2.

Ensure that you have set up your AWS account to be ready to deploy Tanzu clusters.
Refer to the [Prepare to Deploy a Management or Standalone Cluster to Amazon EC2](../aws) docs for instructions on deploying an SSH key-pair and preparing your AWS account.

1. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Complete the configuration steps in the installer interface and create the management cluster. The following configuration settings are recommended:


   *  If you do not specify a name, the installer interface for Amazon EC2 generates a unique name. If you do specify a name, the name must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements described here: [RFC 1123](https://tools.ietf.org/html/rfc1123).
   * Check the "Automate creation of AWS CloudFormation Stack" box if you do not have an existing CloudFormation stack. This stack is used to created IAM resources that Tanzu Community Edition clusters use in Amazon EC2.
     You only need 1 CloudFormation stack per AWS account. CloudFormation is global and not locked to a region. For more information, see [Required IAM resources](../ref-aws/#permissions).

   * Set the instance type size to m5.xlarge or larger for both the control plane node and worker node.

   * Disable **Enable Identity Management Settings**. You can disable identity management for proof-of-concept/development deployments, but it is strongly recommended to implement identity management in production deployments. For more information about enabling Identity Management, see [Identity Management](../aws-install-mgmt/#step-6-identity-management).

2. Validate the management cluster started successfully.

    ```sh
    tanzu management-cluster get
    ```
    The output will look similar to the following:

    ```sh
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

3. Capture the management cluster's kubeconfig and take note of the command for accessing the cluster in the output message, as you will use this for setting the context in the next step.

    ```sh
    tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin
    ```
    Where <``MGMT-CLUSTER-NAME>`` should be set to the name returned by `tanzu management-cluster get`.  <br><br>
    For example, if your management cluster is called 'mtce', you will see a message similar to:
    ```sh
    Credentials of workload cluster 'mtce' have been saved.
    You can now access the cluster by running 'kubectl config use-context mtce-admin@mtce'
    ```

4. Set your kubectl context to the management cluster.

    ```sh
    kubectl config use-context <MGMT-CLUSTER-NAME>-admin@<MGMT-CLUSTER-NAME>
    ```
    Where <``MGMT-CLUSTER-NAME>`` should be set to the name returned by `tanzu management-cluster get`.
5. Validate you can access the management cluster's API server.

    ```sh
    kubectl get nodes
    ```
    The output will look similar to the following:
    ```sh
    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```

6. Next, you will create a workload cluster. First, create a workload cluster configuration file by taking a copy of the management cluster YAML configuration file that was created when you deployed your management cluster. This example names the workload cluster configuration file ``workload1.yaml``.


    ```sh
    cp  ~/.config/tanzu/tkg/clusterconfigs/<MGMT-CONFIG-FILE> ~/.config/tanzu/tkg/clusterconfigs/workload1.yaml
    ```

   * Where ``<MGMT-CONFIG-FILE>`` is the name of the management cluster YAML configuration file. The management cluster YAML configuration file will either have the name you assigned to the management cluster, or if no name was assigned, it will be a randomly generated name.

   * The duplicated file (``workload1.yaml``) will be used as the configuration file for your workload cluster. You can edit the parameters in this new  file as required. For an example of a workload cluster template, see  [Amazon EC2 Workload Cluster Template](../aws-wl-template).


   * In the next two steps you will edit the parameters in this new file (`workload1.yaml`) and then use the file to deploy a workload cluster.


7. In the new workload cluster file (`~/.config/tanzu/tkg/clusterconfigs/workload1.yaml`), edit the CLUSTER_NAME parameter to assign a name to your workload cluster. For example,


   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-workload-cluster
   CLUSTER_PLAN: dev
   ```

   * If you did not specify a name for your management cluster, the installer generated a random unique name. In this case, you must manually add the CLUSTER_NAME parameter and assign a workload cluster name. The workload cluster names must be must be 42 characters or less and must comply with DNS hostname requirements as described here: [RFC 1123](https://tools.ietf.org/html/rfc1123)
   * If you specified a name for your management cluster, the CLUSTER_NAME parameter is present and needs to be changed to the new workload cluster name.
   * The other parameters in ``workload1.yaml`` are likely fine as-is. However, you can change them as required. Validation is performed on the file prior to applying it, so the `tanzu` command will return a message if something necessary is omitted. Reference an example configuration template here:  [Amazon EC2 Workload Cluster Template](../aws-wl-template).

8. Create your workload cluster.

    ```sh
    tanzu cluster create <WORKLOAD-CLUSTER-NAME> --file ~/.config/tanzu/tkg/clusterconfigs/workload1.yaml
    ```

9. Validate the cluster starts successfully.

    ```sh
    tanzu cluster list
    ```

10. Capture the workload cluster's kubeconfig.

    ```sh
    tanzu cluster kubeconfig get <WORKLOAD-CLUSTER-NAME> --admin
    ```

11. Set your `kubectl` context to the workload cluster.

    ```sh
    kubectl config use-context <WORKLOAD-CLUSTER-NAME>-admin@<WORKLOAD-CLUSTER-NAME>
    ```

12. Verify you can see pods in the cluster.

    ```sh
    kubectl get pods --all-namespaces
    ```
    The output will look similar to the following:

    ```sh
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
