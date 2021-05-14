{{% include "/docs/assets/step-one.md" %}}

### <a id="iaas-aws"></a> Configure a Amazon EC2 Infrastructure Provider

1. In the **IaaS Provider** section, enter credentials for your Amazon EC2 account. You have two options:
    - In the **AWS Credential Profile** drop-down, you can select an already existing AWS credential profile. If you select a profile, the access key and session token information configured for your profile are passed to the Installer without displaying actual values in the UI. For information about setting up credential profiles, see [Credential Files and Profiles](aws.md#profiles).
    - Alternately, enter AWS account credentials directly in the **Access Key ID** and **Secret Access Key** fields for your Amazon EC2 account. Optionally specify an AWS session token in **Session Token** if your AWS account is configured to require temporary credentials. For more information on acquiring session tokens, see [Using temporary credentials with AWS resources](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html).
1. In **Region**, select the AWS region in which to deploy the management cluster. If you intend to deploy a production management cluster, this region must have at least three availability zones. This region must also be registered with the SSH key entered in the next field.
1. In **SSH Key Name**, specify the name of an SSH key that is already registered with your Amazon EC2 account and in the region where you are deploying the management cluster. You may have set this up in [Configure AWS Account Credentials and SSH Key](aws.md#register-ssh).
1. If this is the first time that you are deploying a management cluster, select the **Automate creation of AWS CloudFormation Stack** checkbox, and click **Connect**.

   This CloudFormation stack creates the identity and access management (IAM) resources that Tanzu Kubernetes Grid needs to deploy and run clusters on Amazon EC2. For more information, see [Required IAM Resources](aws.md#iam-permissions) in _Prepare to Deploy Management Clusters to Amazon EC2_.

   **IMPORTANT:** The **Automate creation of AWS CloudFormation Stack** checkbox replaces the `clusterawsadm` command line utility that existed in Tanzu Kubernetes Grid v1.1.x and earlier. For existing management and Tanzu Kubernetes clusters initially deployed with v1.1.x or earlier, continue to use the CloudFormation stack that was created by running the `clusterawsadm alpha bootstrap create-stack` command.

   ![Configure the connection to AWS](../images/connect-to-aws.png)
1. If the connection is successful, click **Next**.
1. In the **VPC for AWS** section, do one of the following:

    - To create a new VPC, select **Create new VPC on AWS**, check that the pre-filled CIDR block is available, and click **Next**. If the recommended CIDR block is not available, enter a new IP range in CIDR format for the management cluster to use. The recommended CIDR block for **VPC CIDR** is 10.0.0.0/16.

       ![Create a new VPC](../images/aws-new-vpc.png)

    - To use an existing VPC, select **Select an existing VPC** and select the **VPC ID** from the drop-down menu. The **VPC CIDR** block is filled in automatically when you select the VPC.

       ![Use and existing VPC](../images/aws-existing-vpc.png)

## Configure the Management Cluster Settings

This section applies to all infrastructure providers.

1. In the **Management Cluster Settings** section, select the **Development** or **Production** tile.

   - If you select **Development**, the installer deploys a management cluster with a single control plane node.
   - If you select **Production**, the installer deploys a highly available management cluster with three control plane nodes.

1. In either of the **Development** or **Production** tiles, use the **Instance type** drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs.

   Choose the configuration for the control plane node VMs depending on the expected workloads that it will run. For example, some workloads might require a large compute capacity but relatively little storage, while others might require a large amount of storage and less compute capacity. If you select an instance type in the **Production** tile, the instance type that you selected is automatically selected for the **Worker Node Instance Type**. If necessary, you can change this.

   <!--If you plan on registering the management cluster with Tanzu Mission Control, ensure that your Tanzu Kubernetes clusters meet the requirements listed in [Requirements for Registering a Tanzu Kubernetes Cluster with Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-3AE5F733-7FA7-4B34-8935-C25D41D15EF9.html) in the Tanzu Mission Control documentation.-->

   - **vSphere**: Select a size from the predefined CPU, memory, and storage configurations. The minimum configuration is 2 CPUs and 4 GB memory.
   - **Amazon EC2**: Select an instance size. The drop-down menu lists choices alphabetically, not by size. The minimum configuration is 2 CPUs and 8 GB memory. The list of compatible instance types varies in different regions. For information about the configuration of the different sizes of instances, see [Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/).
   
   ![Select the control plane node configuration](../images/configure-control-plane.png)

1. (Optional) Enter a name for your management or stand-alone cluster.

   If you do not specify a name, the installer generates a unique name. If you do specify a name, that name must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements as outlined in [RFC 952](https://tools.ietf.org/html/rfc952) and amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).

1. Under **Worker Node Instance Type**, select the configuration for the worker node VM.
1. Deselect the **Machine Health Checks** checkbox if you want to
disable [`MachineHealthCheck`](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/machine-health-check.html#machinehealthcheck).

   `MachineHealthCheck` provides node health monitoring and node auto-repair on the clusters that you deploy with this management cluster. You can enable or disable
   `MachineHealthCheck` on clusters after deployment by using the CLI. For instructions, see [Configure Machine Health Checks for Tanzu Kubernetes Clusters](../cluster-lifecycle/configure-health-checks.md).
1. **(vSphere Only)** Under **Control Plane Endpoint**, enter a static virtual IP address or FQDN for API requests to the management cluster.

   Ensure that this IP address is not in your DHCP range, but is in the same subnet as the DHCP range. If you mapped an FQDN to the VIP address, you can specify the FQDN instead of the VIP address. For more information, see [Static VIPs and Load Balancers for vSphere](vsphere.md#load-balancer).

   ![Select the cluster configuration](../images/configure-cluster.png)
1. **(Amazon EC2 only)** Optionally, disable the **Bastion Host** checkbox if a bastion host already exists in the availability zone(s) in which you are deploying the management cluster.

   If you leave this option enabled, Tanzu Kubernetes Grid creates a bastion host for you.

1. **(Amazon EC2 only)** Configure Availability Zones

    1. From the **Availability Zone 1** drop-down menu, select an availability zone for the management cluster. You can select only one availability zone in the **Development** tile. See the image below.

        ![Configure the cluster](../images/aws-az.png)

        If you selected the **Production** tile above, use the **Availability Zone 1**, **Availability Zone 2**, and **Availability Zone 3** drop-down menus to select three unique availability zones for the management cluster. When Tanzu Kubernetes Grid deploys the management cluster, which includes three control plane nodes, it distributes the control plane nodes across these availability zones.

    1. To complete the configuration of the **Management Cluster Settings** section, do one of the following:

        - If you created a new VPC in the **VPC for AWS** section, click **Next**.
        - If you selected an existing VPC in the **VPC for AWS** section, use the **VPC public subnet** and **VPC private subnet** drop-down menus to select existing subnets on the VPC and click **Next**. The image below shows the **Development** tile.

        ![Set the VPC subnets](../images/aws-subnets.png)

1. Click **Next**.

   - If you are deploying the management cluster to vSphere, go to [Configure VMware NSX Advanced Load Balancer](#nsx-adv-lb).
   - If you are deploying the management cluster to Amazon EC2 or Azure, go to [Configure Metadata](#metadata).

