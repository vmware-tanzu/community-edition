## Create Managed Clusters in AWS

This section describes setting up management and workload clusters in Amazon Web
Services (AWS).

There are some prerequisites the installation process will assume.  Refer to the [Prepare to Deploy a Management or Standalone Cluster to
AWS](../aws) docs for instructions on deploying an SSH key-pair and preparing your AWS account.

1. Initialize the Tanzu Community Edition installer interface.

    ```sh
    tanzu management-cluster create --ui
    ```

1. Choose Amazon from the provider tiles.

    ![kickstart amazon tile](/docs/img/kickstart-amazon-tile.png)

1. Fill out the IaaS Provider section.

    ![kickstart amazon iaas](/docs/img/kickstart-amazon-iaas.png)

    * `A`: Whether to use AWS named profiles or provide static
      credentials. It is **highly** recommended you use profiles. This can be
      setup by installing the AWS CLI on the bootstrap machine.
    * `B`: If using profiles, the name of the profile (credentials) you'd like
      to use. By default, profiles are stored in `${HOME}/.aws/credentials`.
    * `C`: [The region of
      AWS](https://aws.amazon.com/about-aws/global-infrastructure/regions_az/)
      you'd like all networking, compute, etc to be created within.

1. Fill out the VPC settings.

    ![kickstart aws vpc](/docs/img/kickstart-amazon-vpc.png)

    * `A`: Whether to create a new Virtual Private Cloud in AWS or use an existing
      one. If using an existing one, you must provide its VPC ID. For initial
      deployments, it is recomended to create a new Virtual Private Cloud. This will
      ensure the installer takes care of all networking creation and configuration.
    * `B`: If creating a new VPC, the CIDR range or IPs to use for hosts (EC2
      VMs).

1. Fill out the Management Cluster Settings.

    ![kickstart aws management cluster settings](/docs/img/kickstart-amazon-mgmt-cluster.png)

    * `A`: Choose between Development profile, with 1 control plane node or
      Production, which features a highly-available three node control plane.
      Additionally, choose the instance type you'd like to use for control plane nodes.
    * `B`: Name the cluster. This is a friendly name that will be used to
      reference your cluster in the Tanzu CLI and `kubectl`.
    * `C`: Choose an SSH key to use for accessing control plane and workload
      nodes. This SSH key must be accessible in the AWS region chosen in a
      previous step. See the [AWS
      documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html#having-ec2-create-your-key-pair) 
      for instructions on creating a key pair.
    * `D`: Whether to enable [Cluster API's machine health
      checks](https://cluster-api.sigs.k8s.io/tasks/healthcheck.html).
    * `E`: Whether to create a bastion host in your VPC. This host will be
      publicly accessible via your SSH key. All Kubernetes-related hosts will
      **not** be accessible without SSHing into this host. If preferred, you can create
      a bastion host independent of the installation process.
    * `F`: Choose whether you'd like to enable [Kubernetes API server
      auditing](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/).
    * `G`: Choose whether you'd like to create the [CloudFormation
      stack](https://aws.amazon.com/cloudformation/) expected by Taznu. Checking
      this box is recommended. If the stack pre-exists, this step will be skipped.
    * `H`: The AWS availability zone in your chosen region to create control
      plane node(s) in. If the Production profile was chosen, you'll have 3
      options of zones, one for each host.
    * `I`: The AWS EC2 instance type to be used for each node creation. See the
      instances types documentation to understand trade-offs between CPU,
      memory, pricing and more.

1. If you would like additional metadata to be tagged in your soon-to-be-created
   AWS infrastructure, fill out the Metadata section.

1. Fill out the Kubernetes Network section.

    ![kickstart kubernetes networking](/docs/img/kickstart-amazon-network.png)

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

    ![kickstart aws os](/docs/img/kickstart-amazon-os.png)

    * `A`: The [Amazon Machine Image
      (AMI)](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) to
      use for Kubernetes host VMs. This list should populate based on known AMIs
      uploaded by VMware. These AMIs are publicly accessible for your use. Choose
      based on your preferred Linux distribution.

1. Skip the TMC Registration section.

1. Click the Review Configuration button.

    > For your record, the configuration settings have been saved to
    > `${HOME}/.config/tanzu/tkg/clusterconfigs`.

1. Deploy the cluster.

    > If you experience issues deploying your cluster, visit the [Troubleshooting
    > documentation](../tsg-bootstrap).

1. Validate the management cluster started successfully.

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

1. Capture the management cluster's kubeconfig and take note of the command for accessing the cluster in the output message, as you will use this for setting the context in the next step.

    ```sh
    tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin
    ```

    Where `<MGMT-CLUSTER-NAME>` should be set to the name returned by `tanzu management-cluster get`.  <br><br>
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
    ```

    The output will look similar to the following:

    ```sh
    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```

1. Next, you will create a workload cluster. First, create a workload cluster configuration file by taking a copy of the management cluster YAML configuration file that was created when you deployed your management cluster. This example names the workload cluster configuration file `workload1.yaml`.

    ```sh
    cp  ~/.config/tanzu/tkg/clusterconfigs/<MGMT-CONFIG-FILE> ~/.config/tanzu/tkg/clusterconfigs/workload1.yaml
    ```

   * Where ``<MGMT-CONFIG-FILE>`` is the name of the management cluster YAML configuration file. The management cluster YAML configuration file will either have the name you assigned to the management cluster, or if no name was assigned, it will be a randomly generated name.

   * The duplicated file (``workload1.yaml``) will be used as the configuration file for your workload cluster. You can edit the parameters in this new  file as required. For an example of a workload cluster template, see  [AWS Workload Cluster Template](../aws-wl-template).

   * In the next two steps you will edit the parameters in this new file (`workload1.yaml`) and then use the file to deploy a workload cluster.

1. In the new workload cluster file (`~/.config/tanzu/tkg/clusterconfigs/workload1.yaml`), edit the CLUSTER_NAME parameter to assign a name to your workload cluster. For example,

   ```yaml
   CLUSTER_CIDR: 100.96.0.0/11
   CLUSTER_NAME: my-workload-cluster
   CLUSTER_PLAN: dev
   ```

   * If you did not specify a name for your management cluster, the installer generated a random unique name. In this case, you must manually add the CLUSTER_NAME parameter and assign a workload cluster name. The workload cluster names must be must be 42 characters or less and must comply with DNS hostname requirements as described here: [RFC 1123](https://tools.ietf.org/html/rfc1123)
   * If you specified a name for your management cluster, the CLUSTER_NAME parameter is present and needs to be changed to the new workload cluster name.
   * The other parameters in ``workload1.yaml`` are likely fine as-is. However, you can change them as required. Validation is performed on the file prior to applying it, so the `tanzu` command will return a message if something necessary is omitted. Reference an example configuration template here:  [AWS Workload Cluster Template](../aws-wl-template).

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
    kubectl config use-context <WORKLOAD-CLUSTER-NAME>-admin@<WORKLOAD-CLUSTER-NAME>
    ```

1. Verify you can see pods in the cluster.

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
