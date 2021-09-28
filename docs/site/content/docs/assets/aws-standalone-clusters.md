## Create Standalone Clusters in AWS

This section covers setting up a standalone cluster in Amazon Web Services (AWS). A standalone cluster provides a workload cluster that is **not** managed by a centralized management cluster.

Ensure that you have set up your AWS account to be ready to deploy Tanzu clusters.
Refer to the [Prepare to Deploy a Management or Standalone Cluster to AWS](../aws) docs for instructions on deploying an SSH key-pair and preparing your AWS account.

1. Initialize the Tanzu Community Edition Installer Interface.

    ```sh
    tanzu standalone-cluster create --ui
    ```
1. Choose Amazon from the provider tiles.

    ![kickstart amazon tile](/docs/img/kickstart-amazon-tile.png)

1. Fill out the IaaS Provider section.

    ![kickstart vsphere iaas](/docs/img/kickstart-amazon-iaas.png)

    * `A`: Whether to use AWS named profiles or provide static
      credentials. It is **highly** recommended you use profiles. This can be
      setup by installing the AWS CLI on the bootstrap machine.
    * `B`: If using profiles, the name of the profile (credentials) you'd like
      to use. By default, profiles are stored in `${HOME}/.aws/credentials`.
    * `C`: [The region of
      AWS](https://aws.amazon.com/about-aws/global-infrastructure/regions_az/)
      you'd like all networking, compute, etc to be created within.

1. Fill out the VPC settings.

    ![kickstart aws iaas](/docs/img/kickstart-amazon-vpc.png)

    * `A`: Whether to create a new Virtual Private Cloud in AWS or use an existing
      one. If using an existing one, you must provide its VPC ID. For initial
      deployments, it is recomended to create a new Virtual Private Cloud. This will
      ensure the installer takes care of all networking creation and configuration.
    * `B`: If creating a new VPC, the CIDR range or IPs to use for hosts (EC2
      VMs).

1. Fill out the Standalone Cluster Settings.

    ![kickstart aws standalone cluster settings](/docs/img/kickstart-amazon-sa-cluster.png)

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
      stack](https://aws.amazon.com/cloudformation/) expected by Tanzu. Checking
      this box is recommended. If the stack pre-exists, this step will be skipped.
    * `H`: The AWS availability zone in your chosen region to create control
      plane node(s) in. If the Production profile was chosen, you'll have 3
      options of zones, one for each host.

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

1. Click the Review Configuration button.

    > For your record, the configuration settings have been saved to
    > `${HOME}/.config/tanzu/tkg/clusterconfigs`.

1. Deploy the cluster.

    > If you experience issues deploying your cluster, visit the [Troubleshooting
    > documentation](../tsg-bootstrap).

1. Set your kubectl context to the cluster.

    ```sh
    kubectl config use-context <STANDALONE-CLUSTER-NAME>-admin@<STANDALONE-CLUSTER-NAME>
    ```

    Where `<STANDALONE-CLUSTER-NAME>` is the name of the standalone cluster that you specified or if you didn't specify a name, it's the randomly generated name.

1. Validate you can access the cluster's API server.

    ```sh
    kubectl get nodes
    ```

    The output will look similar to the following:

    ```sh
    NAME                                       STATUS   ROLES                  AGE    VERSION
    ip-10-0-1-133.us-west-2.compute.internal   Ready    <none>                 123m   v1.20.1+vmware.2
    ip-10-0-1-76.us-west-2.compute.internal    Ready    control-plane,master   125m   v1.20.1+vmware.2
    ```
