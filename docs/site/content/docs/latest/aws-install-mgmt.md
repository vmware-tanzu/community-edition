# Deploy a Management Cluster to Amazon EC2

{{% include "/docs/assets/step1.md" %}}

### Step 1: IaaS Provider

1. In the **IaaS Provider** section, enter credentials for your Amazon EC2 account. You have two options:  
    - In the **AWS Credential Profile** drop-down, you can select an already existing AWS credential profile. If you select a profile, the access key and session token information configured for your profile are passed to the Installer without displaying actual values in the UI. 
    - Alternately, enter AWS account credentials directly in the **Access Key ID** and **Secret Access Key** fields for your Amazon EC2 account. For information about setting up credential profiles, see [Prepare to Deploy a Management or Standalone Cluster to Amazon EC2](aws).
    - Optionally, specify an AWS session token in **Session Token** if your AWS account is configured to require temporary credentials. For more information on acquiring session tokens, see [Using temporary credentials with AWS resources](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html).  
2. In **Region**, select the AWS region in which to deploy the cluster. If you intend to deploy a production management cluster, this region must have at least three availability zones. This region must also be registered with the SSH key entered in the next field.  
3. In **SSH Key Name**, specify the name of an SSH key that is already registered with both your Amazon EC2 account and the region where you are deploying the cluster. For information about setting up credential profiles, see [Prepare to Deploy a Management or Stand-alone Cluster to Amazon EC2](aws.md#profiles).  
4. If this is the first time deploying a cluster, select the **Automate creation of AWS CloudFormation Stack** checkbox, and click **Connect**.  

   The CloudFormation stack creates the identity and access management (IAM) resources that Tanzu Community Edition needs to deploy and run clusters on Amazon EC2. For more information, see [Required IAM Resources](ref-aws.md#permissions). 
<!--   **IMPORTANT:** The **Automate creation of AWS CloudFormation Stack** checkbox replaces the `clusterawsadm` command line utility that existed in Tanzu Kubernetes Grid v1.1.x and earlier. For existing management and Tanzu Kubernetes clusters initially deployed with v1.1.x or earlier, continue to use the CloudFormation stack that was created by running the `clusterawsadm alpha bootstrap create-stack` command.-->

<!--![Configure the connection to AWS](../images/connect-to-aws.png)-->
5. If the connection is successful, click **Next**.

### Step 2: VPC for AWS
In the **VPC for AWS** section, do one of the following:  
    - To create a new VPC, select **Create new VPC on AWS**, check that the pre-filled CIDR block is available, and click **Next**. If the recommended CIDR block is not available, enter a new IP range in CIDR format for the management cluster to use. The recommended CIDR block for **VPC CIDR** is 10.0.0.0/16.  
    - To use an existing VPC, select **Select an existing VPC** and select the **VPC ID** from the drop-down menu. The **VPC CIDR** block is filled in automatically when you select the VPC.

For more information about VPC, see [Virtual Private Clouds and NAT Gateway Limits](ref-aws.md/#vpc).
<!--![Create a new VPC](../images/aws-new-vpc.png)
![Use and existing VPC](../images/aws-existing-vpc.png)-->

### Step 3: Management Cluster Settings

1. In the **Management Cluster Settings** section, select an instance size for either **Development** or **Production**. If you select **Development**, the installer deploys a management cluster with a single control plane node. If you select **Production**, the installer deploys a highly available management cluster with three control plane nodes. Use the **Instance type** drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs.  Choices are listed alphabetically, not by size. The minimum configuration is 2 CPUs and 8 GB memory. The list of compatible instance types varies in different regions. For information about the configuration of the different sizes of instances, see [Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/). 

<!--Choose the configuration for the control plane node VMs depending on the expected workloads that it will run. For example, some workloads might require a large compute capacity but relatively little storage, while others might require a large amount of storage and less compute capacity. -->    


<!--If you plan on registering the management cluster with Tanzu Mission Control, ensure that your Tanzu Kubernetes clusters meet the requirements listed in [Requirements for Registering a Tanzu Kubernetes Cluster with Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-3AE5F733-7FA7-4B34-8935-C25D41D15EF9.html) in the Tanzu Mission Control documentation.
    - **vSphere**: Select a size from the predefined CPU, memory, and storage configurations. The minimum configuration is 2 CPUs and 4 GB memory.-->
    
<!--![Select the control plane node configuration](../images/configure-control-plane.png)-->
2. (Optional) Enter a name for your management cluster. If you do not specify a name, the installer generates a unique name. If you do specify a name, the name must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements as outlined in [RFC 952](https://tools.ietf.org/html/rfc952) and amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).
3. Under **Worker Node Instance Type**, select the configuration for the worker node VM.  If you select an instance type in the **Production** tile, the instance type that you select is automatically selected for the **Worker Node Instance Type**. If necessary, you can change this.   
4. [`MachineHealthCheck`](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/machine-health-check.html#machinehealthcheck) is enabled by default. `MachineHealthCheck` provides node health monitoring and node auto-repair on the clusters that you deploy with this management cluster. You can enable or disable `MachineHealthCheck` on clusters after deployment by using the CLI. For instructions, see [Configure Machine Health Checks for Tanzu Kubernetes Clusters](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3/vmware-tanzu-kubernetes-grid-13/GUID-cluster-lifecycle-configure-health-checks.html).
5. (Optional) Disable the **Bastion Host** checkbox if a bastion host already exists in the availability zone(s) in which you are deploying the management cluster.  
6. Configure Availability Zones. From the **Availability Zone 1** drop-down menu, select an availability zone for the management cluster. You can select only one availability zone in the **Development** tile.  If you selected the **Production** tile, use the **Availability Zone 1**, **Availability Zone 2**, and **Availability Zone 3** drop-down menus to select three unique availability zones for the management cluster. When Tanzu deploys the management cluster, which includes three control plane nodes, it distributes the control plane nodes across these availability zones.  
<!--![Configure the cluster](../images/aws-az.png)-->        
7. To complete the configuration of the **Management Cluster Settings** section, do one of the following:
   - If you created a new VPC in the **VPC for AWS** section, click **Next**.
   - If you selected an existing VPC in the **VPC for AWS** section, use the **VPC public subnet** and **VPC private subnet** drop-down menus to select existing subnets on the VPC and click **Next**. 

<!--![Set the VPC subnets](../images/aws-subnets.png)-->

### Step 4: Metadata
{{% include "/docs/assets/metadata.md" %}}


### Step 5: Kubernetes Network
<!-- note to self: right now I can't figure a good way to turn this into an include that could be reused across amazon and vsphere as there is too much mixed up information about both in it, so it will be added manually to each and cleaned up appropriately - so this will need to be copied into both vsphere and amazon topics-->

1.  Review the **Cluster Service CIDR** and **Cluster Pod CIDR** ranges. If the recommended CIDR ranges of `100.64.0.0/13` and `100.96.0.0/11` are unavailable, update the values under **Cluster Service CIDR** and **Cluster Pod CIDR**.  

<!--![Configure the Kubernetes service network](../images/install-v-6k8snet.png) -->   

2. (Optional) To send outgoing HTTP(S) traffic from the management cluster to a proxy, toggle **Enable Proxy Settings** and follow the instructions below to enter your proxy information. Tanzu applies these settings to kubelet, containerd, and the control plane. You can choose to use one proxy for HTTP traffic and another proxy for HTTPS traffic or to use the same proxy for both HTTP and HTTPS traffic.  

    - To add your HTTP proxy information: Under **HTTP Proxy URL**, enter the URL of the proxy that handles HTTP requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`.  If the proxy requires authentication, under **HTTP Proxy Username** and **HTTP Proxy Password**, enter the username and password to use to connect to your HTTP proxy.

    - To add your HTTPS proxy information: If you want to use the same URL for both HTTP and HTTPS traffic, select **Use the same configuration for https proxy**.  If you want to use a different URL for HTTPS traffic, enter the URL of the proxy that handles HTTPS requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`. If the proxy requires authentication, under **HTTPS Proxy Username** and **HTTPS Proxy Password**, enter the username and password to use to connect to your HTTPS proxy.

    - Under **No proxy**, enter a comma-separated list of network CIDRs or hostnames that must bypass the HTTP(S) proxy. For example, `noproxy.yourdomain.com,192.168.0.0/24`. Internally, Tanzu appends `localhost`, `127.0.0.1`, your VPC CIDR, **Cluster Pod CIDR**, and **Cluster Service CIDR**, `.svc`, `.svc.cluster.local`, and `169.254.0.0/16` to the list that you enter in this field.
      

    **Important:** If the management cluster VMs need to communicate with external services and infrastructure endpoints in your Tanzu environment, ensure that those endpoints are reachable by the proxies that you configured above or add them to **No proxy**. Depending on your environment configuration, this may include, but is not limited to, your OIDC or LDAP server, and Harbor.

### Step 6: Identity Management
{{% include "/docs/assets/identity-management.md" %}}

{{% include "/docs/assets/final-step.md" %}}


