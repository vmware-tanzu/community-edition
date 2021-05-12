WIP - DRAFT - WIP - DRAFT - WIP - DRAFT


# Deploy Management Clusters with the Installer Interface

This topic describes how to use the Tanzu Kubernetes Grid installer interface to deploy a management cluster to Amazon EC2, and Microsoft Azure. The Tanzu Kubernetes Grid installer interface guides you through the deployment of the management cluster, and provides different configurations for you to select or reconfigure. If this is the first time that you are deploying a management cluster to a given infrastructure provider, it is recommended to use the installer interface.

## <a id="prereqs"></a> Prerequisites

Before you can deploy a management cluster, you must make sure that your environment meets the requirements for the target infrastructure provider.

### General Prerequisites

- Ensure that you have met all of the requirements and followed all of the procedures in [Install the Tanzu CLI](../installation-cli.md).
- Ensure you have completed the steps in [Prepare to Deploy Clusters to Amazon EC2](../aws.md)
<!--- For production deployments, it is strongly recommended to enable identity management for your clusters. For information about the preparatory steps to perform before you deploy a management cluster, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).
- If you want to register your management cluster with Tanzu Mission Control, follow the procedure in [Register Your Management Cluster with Tanzu Mission Control](register_tmc.md).
- If you are deploying clusters in an internet-restricted environment to either vSphere or Amazon EC2, you must also perform the steps in [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md).-->

### Amazon EC2 Prerequisites

- Make sure that you have met the all of the requirements listed [Prepare to Deploy Management Clusters to Amazon EC2](aws.md).
- For information about the configurations of the different sizes of node instances, for example `t3.large`, or `t3.xlarge`, see [Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/).
- For information about when to create a Virtual Private Cloud (VPC) and when to reuse an existing VPC, see [Resource Usage in Your Amazon Web Services Account](aws.md#aws-resources).

### Microsoft Azure Prerequisites

- Make sure that you have met the requirements listed in [Prepare to Deploy Management Clusters to Microsoft Azure](azure.md).
- For information about the configurations of the different sizes of node instances for Azure, for example, `Standard_D2s_v3` or `Standard_D4s_v3`, see [Sizes for virtual machines in Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/sizes).

## Start the Installer Interface

<p class="note warning"><strong>Warning</strong>: The <code>tanzu management-cluster create</code> command takes time to complete.
While <code>tanzu management-cluster create</code> is running, do not run additional invocations of <code>tanzu management-cluster create</code> on the same bootstrap machine to deploy multiple management clusters, change context, or edit <code>~/.kube-tkg/config</code>.</p>

1. On the machine on which you downloaded and installed the Tanzu CLI, run the `tanzu management-cluster create` command with the `--ui` option.

   ```
   tanzu management-cluster create --ui
   ```

   The installer interface launches in a browser and takes you through steps to configure the management cluster.
   The `tanzu management-cluster create --ui` command saves the settings from your installer input in a cluster configuration file.
   After you confirm your input values on the last pane of the installer interface, the installer saves them to `~/.tanzu/tkg/clusterconfigs` with a generated filename of the form `UNIQUE-ID.yaml`.

   By default Tanzu Kubernetes Grid saves the `kubeconfig` for all management clusters in the `~/.kube-tkg/config` file. If you want to save the `kubeconfig` file for your management cluster to a different location, set the `KUBECONFIG` environment variable before running `tanzu management-cluster create`.

   ```
   KUBECONFIG=/path/to/mc-kubeconfig.yaml
   ```

   When you run the `tanzu management-cluster create --ui` command, it validates that your system meets the prerequisites:

   - NTP is running on the bootstrap machine on which you are running `tanzu management-cluster create` and on the hypervisor.
   - A DHCP server is available.
   - The CLI can connect to the location from which it pulls the required images.
   - Docker is running.

   If the prerequisites are met, `tanzu management-cluster create --ui` launches the Tanzu Kubernetes Grid installer interface.

   By default, `tanzu management-cluster create --ui` opens the installer interface locally, at http://127.0.0.1:8080 in your default browser.
   The [Installer Interface Options](#ui-options) section below explains how you can change where the installer interface runs, including running it on a different machine from the `tanzu` CLI.

1. Click the **Deploy** button for **VMware vSphere**, **Amazon EC2**, or **Microsoft Azure**.

   ![Tanzu Kubernetes Grid installer interface welcome page with Deploy to vSphere button](../images/deploy-management-cluster.png)

### <a id="ui-options"></a> Installer Interface Options

By default, `tanzu management-cluster create --ui` opens the installer interface locally, at http://127.0.0.1:8080 in your default browser.
You can use the `--browser` and `--bind` options to control where the installer interface runs:

- `--browser` specifies the local browser to open the interface in.
   - Supported values are `chrome`, `firefox`, `safari`, `ie`, `edge`, or `none`.
   - Use `none` with `--bind` to run the interface on a different machine, as described below.
- `--bind` specifies the IP address and port to serve the interface from.

<p class="note warning"><strong>Warning</strong>: Serving the installer interface from a non-default IP address and port could expose the <code>tanzu</code> CLI to a potential security risk while the interface is running. VMware recommends passing in to the <code>--bind</code> option an IP and port on a secure network.</p>

Use cases for `--browser` and `--bind` include:

- If another process is already using http://127.0.0.1:8080, use `--bind` to serve the interface from a different local port.
- To run the `tanzu` CLI and create management clusters on a remote machine, and run the installer interface locally or elsewhere:
  1. On the remote bootstrap machine, run `tanzu management-cluster create --ui` with the following options and values:
      - `--bind`: an IP address and port for the remote machine
      - `--browser`: `none`
        ```
        tanzu management-cluster create --ui --bind 192.168.1.87:5555 --browser none
        ```  
  1. On the local UI machine, browse to the remote machine's IP address to access the installer interface.

## Configure the Infrastructure Provider

The options to configure the infrastructure provider section of the installer interface depend on which provider you are using.

- [Configure a vSphere Infrastructure Provider](#iaas-vsphere)
- [Configure an Amazon EC2 Provider](#iaas-aws)
- [Configure a Microsoft Azure Infrastructure Provider](#iaas-azure)

### <a id="iaas-vsphere"></a> Configure a vSphere Infrastructure Provider

1. In the **IaaS Provider** section, enter the IP address or fully qualified domain name (FQDN) for the vCenter Server instance on which to deploy the management cluster.

   Tanzu Kubernetes Grid does not support IPv6 addresses. This is because upstream Kubernetes only provides alpha support for IPv6. Always provide IPv4 addresses in the procedures in this topic.
1. Enter the vCenter Single Sign On username and password for a user account that has the required privileges for Tanzu Kubernetes Grid operation, and click **Connect**.

   ![Configure the connection to vSphere](../images/install-v-1iaas.png)

1. Verify the SSL thumbprint of the vCenter Server certificate and click **Continue** if it is valid.

   For information about how to obtain the vCenter Server certificate thumbprint, see [Obtain vSphere Certificate Thumbprints](vsphere.md#vc-thumbprint).

   ![Verify vCenter Server certificate thumbprint](../images/vsphere-thumprint.png)

1. If you are deploying a management cluster to a vSphere 7 instance, confirm whether or not you want to proceed with the deployment.   

   On vSphere 7, the vSphere with Tanzu option includes a built-in supervisor cluster that works as a management cluster and provides a better experience than a separate management cluster deployed by Tanzu Kubernetes Grid.  Deploying a Tanzu Kubernetes Grid management cluster to vSphere 7 when vSphere with Tanzu is not enabled is supported, but the preferred option is to enable vSphere with Tanzu and use the Supervisor Cluster. VMware Cloud on AWS and Azure VMware Solution do not support a supervisor cluster, so you need to deploy a management cluster.
   For information, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).

   To reflect the recommendation to use vSphere with Tanzu when deploying to vSphere 7, the Tanzu Kubernetes Grid installer behaves as follows:

      - **If vSphere with Tanzu is enabled**, the installer informs you that deploying a management cluster is not possible, and exits.
      - **If vSphere with Tanzu is not enabled**, the installer informs you that deploying a Tanzu Kubernetes Grid management cluster is possible but not recommended, and presents a choice:
          - **Configure vSphere with Tanzu** opens the vSphere Client so you can configure your Supervisor Cluster as described in [Configuring and Managing a Supervisor Cluster](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-21ABC792-0A23-40EF-8D37-0367B483585E.html) in the vSphere documentation.
          - **Deploy TKG Management Cluster** allows you to continue deploying a management cluster, against recommendation for vSphere 7, but as required for VMware Cloud on AWS and Azure VMware Solution. When using vSphere 7, the preferred option is to enable vSphere with Tanzu and use the built-in Supervisor Cluster instead of deploying a Tanzu Kubernetes Grid management cluster.

   ![Deploy management cluster to vSphere 7](../images/vsphere7-detected.png)

1. Select the datacenter in which to deploy the management cluster from the **Datacenter** drop-down menu.

1. Paste the contents of your SSH public key into the text box and click **Next**.

   ![Select datacenter and provide SSH public key](../images/dc-ssh-vsphere.png)

For the next steps, go to [Configure the Management Cluster Settings](#config-mgmt-cluster).

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

For the next steps, go to [Configure the Management Cluster Settings](#config-mgmt-cluster).

### <a id="iaas-azure"></a> Configure a Microsoft Azure Infrastructure Provider

**IMPORTANT**: If this is the first time that you are deploying a management cluster to Azure with a new version of Tanzu Kubernetes Grid, for example v1.3.1, make sure that you have accepted the base image license for that version. For information, see [Accept the Base Image License](azure.md#license) in *Prepare to Deploy Management Clusters to Microsoft Azure*.

1. In the **IaaS Provider** section, enter the Tenant ID, Client ID, Client Secret, and Subscription ID for your Azure account and click **Connect**.  You recorded these values when you registered an Azure app and created a secret for it using the Azure Portal.

   ![Configure the connection to Azure](../images/connect-to-azure.png)
1. Select the Azure region in which to deploy the management cluster.
1. Paste the contents of your SSH public key, such as `.ssh/id_rsa.pub`, into the text box.

1. Under **Resource Group**, select either the **Select an existing resource group** or the **Create a new resource group** radio button.

    - If you select **Select an existing resource group**, use the drop-down menu to select the group, then click **Next**.

       ![Select existing resource group](../images/select-azure-resource-group.png)

    - If you select **Create a new resource group**, enter a name for the new resource group and then click **Next**.

       ![Create new resource group](../images/create-azure-resource-group.png)

1. In the **VNET for Azure** section, select either the **Create a new VNET on Azure** or the **Select an existing VNET** radio button.

    - If you select **Create a new VNET on Azure**, use the drop-down menu to select the resource group in which to create the VNET and provide the following:

       - A name and a CIDR block for the VNET. The default is `10.0.0.0/16`.
       - A name and a CIDR block for the control plane subnet. The default is `10.0.0.0/24`.
       - A name and a CIDR block for the worker node subnet. The default is `10.0.1.0/24`.

       ![Create a new VNET on Azure](../images/create-vnet-azure.png)

       After configuring these fields, click **Next**.

    - If you select **Select an existing VNET**, use the drop-down menus to select the resource group in which the VNET is located, the VNET name, the control plane and worker node subnets, and then click **Next**.

       ![Select an existing VNET](../images/select-vnet-azure.png)

## <a id="config-mgmt-cluster"></a> Configure the Management Cluster Settings

This section applies to all infrastructure providers.

1. In the **Management Cluster Settings** section, select the **Development** or **Production** tile.

   - If you select **Development**, the installer deploys a management cluster with a single control plane node.
   - If you select **Production**, the installer deploys a highly available management cluster with three control plane nodes.

1. In either of the **Development** or **Production** tiles, use the **Instance type** drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs.

   Choose the configuration for the control plane node VMs depending on the expected workloads that it will run. For example, some workloads might require a large compute capacity but relatively little storage, while others might require a large amount of storage and less compute capacity. If you select an instance type in the **Production** tile, the instance type that you selected is automatically selected for the **Worker Node Instance Type**. If necessary, you can change this.

   If you plan on registering the management cluster with Tanzu Mission Control, ensure that your Tanzu Kubernetes clusters meet the requirements listed in [Requirements for Registering a Tanzu Kubernetes Cluster with Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-3AE5F733-7FA7-4B34-8935-C25D41D15EF9.html) in the Tanzu Mission Control documentation.

   - **vSphere**: Select a size from the predefined CPU, memory, and storage configurations. The minimum configuration is 2 CPUs and 4 GB memory.
   - **Amazon EC2**: Select an instance size. The drop-down menu lists choices alphabetically, not by size. The minimum configuration is 2 CPUs and 8 GB memory. The list of compatible instance types varies in different regions. For information about the configuration of the different sizes of instances, see [Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/).
   - **Microsoft Azure**: Select an instance size. The minimum configuration is 2 CPUs and 8 GB memory. The list of compatible instance types varies in different regions. For information about the configurations of the different sizes of node instances for Azure, see [Sizes for virtual machines in Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/sizes).

   ![Select the control plane node configuration](../images/configure-control-plane.png)

1. Optionally enter a name for your management cluster.

   If you do not specify a name, Tanzu Kubernetes Grid automatically generates a unique name. If you do specify a name, that name must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements as outlined in [RFC 952](https://tools.ietf.org/html/rfc952) and amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).

1. Under **Worker Node Instance Type**, select the configuration for the worker node VM.
1. Deselect the **Machine Health Checks** checkbox if you want to
disable [`MachineHealthCheck`](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/machine-health-check.html#machinehealthcheck).

   `MachineHealthCheck` provides node health monitoring and node auto-repair on the clusters that you deploy with this management cluster. You can enable or disable
   `MachineHealthCheck` on clusters after deployment by using the CLI. For instructions, see [Configure Machine Health Checks for Tanzu Kubernetes Clusters](../cluster-lifecycle/configure-health-checks.md).
1. **(Azure Only)** If you are deploying the management cluster to Azure, click **Next**.

   For the next steps for an Azure deployment, go to [Configure Metadata](#metadata).
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

## <a id="nsx-adv-lb"></a> (vSphere Only) Configure VMware NSX Advanced Load Balancer

VMware NSX Advanced Load Balancer provides an L4 load balancing solution for vSphere. NSX Advanced Load Balancer includes a Kubernetes operator that integrates with the Kubernetes API to manage the lifecycle of load balancing and ingress resources for workloads. To use NSX Advanced Load Balancer, you must first deploy it in your vSphere environment. For information, see [Install VMware NSX Advanced Load Balancer on a vSphere Distributed Switch](install-nsx-adv-lb.md).

In the optional **VMware NSX Advanced Load Balancer** section, you can configure Tanzu Kubernetes Grid to use NSX Advanced Load Balancer. By default all workload clusters will use the load balancer.

1. For **Controller Host**, enter the IP address or FQDN of the Controller VM.
1. Enter the username and password that you set for the Controller host when you deployed it, and click **Verify Credentials**.
1. Use the **Cloud Name** drop-down menu to select the cloud that you created in your NSX Advanced Load Balancer deployment.

   For example, `Default-Cloud`.
1. Use the **Service Engine Group Name** drop-down menu to select a Service Engine Group.

   For example, `Default-Group`.
1. For **VIP Network Name**, use the drop-down menu to select the name of the network where the load balancer floating IP Pool resides.

   The VIP network for NSX Advanced Load Balancer must be present in the same vCenter Server instance as the Kubernetes network that Tanzu Kubernetes Grid uses. This allows NSX Advanced Load Balancer to discover the Kubernetes network in vCenter Server and to deploy and configure Service Engines. The drop-down menu is present in Tanzu Kubernetes Grid v1.3.1 and later. In v1.3.0, you enter the name manually.

   You can see the network in the **Infrastructure** > **Networks** view of the NSX Advanced Load Balancer interface.
1. For **VIP Network CIDR**, use the drop-down menu to select the CIDR of the subnet to use for the load balancer VIP.

   This comes from one of the VIP Network's configured subnets. You can see the subnet CIDR for a particular network in the **Infrastructure** > **Networks** view of the NSX Advanced Load Balancer interface. The drop-down menu is present in Tanzu Kubernetes Grid v1.3.1 and later. In v1.3.0, you enter the CIDR manually.

1. Paste the contents of the Certificate Authority that is used to generate your Controller Certificate into the **Controller Certificate Authority** text box.

   If you have a self-signed Controller Certificate, the Certificate Authority is the same as the Controller Certificate.
1. (Optional) Enter one or more cluster labels to identify clusters on which to selectively enable NSX Advanced Load Balancer or to customize NSX Advanced Load Balancer Settings per group of clusters.

   By default, all clusters that you deploy with this management cluster will enable NSX Advanced Load Balancer. All clusters will share the same VMware NSX Advanced Load Balancer Controller, Cloud, Service Engine Group, and VIP Network as you entered previously. This cannot be changed later. To only enable the load balancer on a subset of clusters, or to preserve the ability to customize NSX Advanced Load Balancer settings for a group of clusters, add labels in the format `key: value`. For example `team: tkg`.

   This is useful in the following scenarios:

   - You want to configure different sets of workload clusters to different Service Engine Groups to implement isolation or to support more Service type Load Balancers than one Service Engine Group's capacity.
   - You want to configure different sets of workload clusters to different Clouds because they are deployed in separate sites.

   **NOTE**: Labels that you define here will be used to create a label selector. Only workload cluster `Cluster` objects that have the matching labels will have the load balancer enabled. As a consequence, you are responsible for making sure that the workload cluster's `Cluster` object has the corresponding labels. For example, if you use `team: tkg`, to enable the load balancer on a workload cluster, you will need to perform the following steps after deployment of the management cluster:

   1. Set `kubectl` to the management cluster's context.

      ```
      kubectl config set-context management-cluster@admin
      ```

   1. Label the `Cluster` object of the corresponding workload cluster with the labels defined. If you define multiple key-values, you need to apply all of them.     

      ```
      kubectl label cluster <cluster-name> team=tkg
      ```      

   ![Configure NSX Advanced Load Balancer](../images/install-v-3nsx.png)
1. Click **Next** to configure metadata.

## <a id="metadata"></a> Configure Metadata

This section applies to all infrastructure providers.

In the optional **Metadata** section, optionally provide descriptive information about this management cluster.

Any metadata that you specify here applies to the management cluster and to the Tanzu Kubernetes clusters that it manages, and can be accessed by using the cluster management tool of your choice.

- **Location**: The geographical location in which the clusters run.
- **Description**: A description of this management cluster. The description has a maximum length of 63 characters and must start and end with a letter. It can contain only lower case letters, numbers, and hyphens, with no spaces.
- **Labels**: Key/value pairs to help users identify clusters, for example `release : beta`, `environment : staging`, or `environment : production`. For more information, see [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) in the Kubernetes documentation.<br />
You can click **Add** to apply multiple labels to the clusters.

![Add cluster metadata](../images/install-v-4metadata.png)

If you are deploying to vSphere, click **Next** to go to [Configure Resources](#resources). If you are deploying to Amazon EC2 or Azure, click **Next** to go to [Configure the Kubernetes Network and Proxies](#network).

## <a id="resources"></a> (vSphere Only) Configure Resources

1. In the **Resources** section, select vSphere resources for the management cluster to use, and click **Next**.

   - Select the VM folder in which to place the management cluster VMs.
   - Select a vSphere datastore for the management cluster to use.
   - Select the cluster, host, or resource pool in which to place the management cluster.

   If appropriate resources do not already exist in vSphere, without quitting the Tanzu Kubernetes Grid installer, go to vSphere to create them. Then click the refresh button so that the new resources can be selected.

   ![Select vSphere resources](../images/install-v-5resources.png)

## <a id="network"></a> Configure the Kubernetes Network and Proxies

This section applies to all infrastructure providers.

1. In the **Kubernetes Network** section, configure the networking for Kubernetes services and click **Next**.

   * **(vSphere only)** Under **Network Name**, select a vSphere network to use as the Kubernetes service network.
   * Review the **Cluster Service CIDR** and **Cluster Pod CIDR** ranges. If the recommended CIDR ranges of `100.64.0.0/13` and `100.96.0.0/11` are unavailable, update the values under **Cluster Service CIDR** and **Cluster Pod CIDR**.

   ![Configure the Kubernetes service network](../images/install-v-6k8snet.png)

1. (Optional) To send outgoing HTTP(S) traffic from the management cluster to a proxy, toggle **Enable Proxy Settings** and follow the instructions below to enter your proxy information. Tanzu Kubernetes Grid applies these settings to kubelet, containerd, and the control plane.

   You can choose to use one proxy for HTTP traffic and another proxy for HTTPS traffic or to use the same proxy for both HTTP and HTTPS traffic.

   1. To add your HTTP proxy information:

       1. Under **HTTP Proxy URL**, enter the URL of the proxy that handles HTTP requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`.
       1. If the proxy requires authentication, under **HTTP Proxy Username** and **HTTP Proxy Password**, enter the username and password to use to connect to your HTTP proxy.

   1. To add your HTTPS proxy information:

       * If you want to use the same URL for both HTTP and HTTPS traffic, select **Use the same configuration for https proxy**.
       * If you want to use a different URL for HTTPS traffic, do the following:

         1. Under **HTTPS Proxy URL**, enter the URL of the proxy that handles HTTPS requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`.
         1. If the proxy requires authentication, under **HTTPS Proxy Username** and **HTTPS Proxy Password**, enter the username and password to use to connect to your HTTPS proxy.

   1. Under **No proxy**, enter a comma-separated list of network CIDRs or hostnames that must bypass the HTTP(S) proxy.

      For example, `noproxy.yourdomain.com,192.168.0.0/24`.

      - **vSphere**: You must enter the CIDR of the vSphere network that you selected under **Network Name**. The vSphere network CIDR includes the IP address of your **Control Plane Endpoint**. If you entered an FQDN under **Control Plane Endpoint**, add both the FQDN and the vSphere network CIDR to **No proxy**. Internally, Tanzu Kubernetes Grid appends `localhost`, `127.0.0.1`, the values of **Cluster Pod CIDR** and **Cluster Service CIDR**, `.svc`, and `.svc.cluster.local` to the list that you enter in this field.
      - **Amazon EC2**: Internally, Tanzu Kubernetes Grid appends `localhost`, `127.0.0.1`, your VPC CIDR, **Cluster Pod CIDR**, and **Cluster Service CIDR**, `.svc`, `.svc.cluster.local`, and `169.254.0.0/16` to the list that you enter in this field.
      - **Azure**:  Internally, Tanzu Kubernetes Grid appends `localhost`, `127.0.0.1`, your VNET CIDR, **Cluster Pod CIDR**, and **Cluster Service CIDR**, `.svc`, `.svc.cluster.local`, `169.254.0.0/16`, and `168.63.129.16` to the list that you enter in this field.

      **Important:** If the management cluster VMs need to communicate with external services and infrastructure endpoints in your Tanzu Kubernetes Grid environment, ensure that those endpoints are reachable by the proxies that you configured above or add them to **No proxy**. Depending on your environment configuration, this may include, but is not limited to, your OIDC or LDAP server, Harbor, and in the case of vSphere, NSX-T and NSX Advanced Load Balancer.

## <a id="id-mgmt"></a> Configure Identity Management

This section applies to all infrastructure providers. For information about how Tanzu Kubernetes Grid implements identity management, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).

1. In the **Identity Management** section, optionally disable **Enable Identity Management Settings** .

   ![Configure external Identity Provider](../images/install-v-7id.png)

   You can disable identity management for proof-of-concept deployments, but it is strongly recommended to implement identity management in production deployments. If you disable identity management, you can reenable it later.   
1. If you enable identity management, select **OIDC** or **LDAPS**.

   **OIDC**:

   Provide details of your OIDC provider account, for example, Okta.

   - **Issuer URL**: The IP or DNS address of your OIDC server.
   - **Client ID**: The `client_id` value that you obtain from your OIDC provider. For example, if your provider is Okta, log in to Okta, create a Web application, and select the **Client Credentials** options in order to get a `client_id` and `secret`.
   - **Client Secret**: The `secret` value that you obtain from your OIDC provider.
   - **Scopes**: A comma separated list of additional scopes to request in the token response. For example, `openid,groups,email`.
   - **Username Claim**: The name of your username claim. This is used to set a user's username in the JSON Web Token (JWT) claim. Depending on your provider, enter claims such as `user_name`, `email`, or `code`.
   - **Groups Claim**: The name of your groups claim. This is used to set a user's group in the JWT claim. For example, `groups`.

   ![Configure external Identity Provider](../images/install-v-7id-oidc.png)

   **LDAPS**:

   Provide details of your company's LDAPS server. All settings except for **LDAPS Endpoint** are optional.

   - **LDAPS Endpoint**: The IP or DNS address of your LDAPS server. Provide the address and port of the LDAP server, in the form `host:port`.
   - **Bind DN**: The DN for an application service account. The connector uses these credentials to search for users and groups. Not required if the LDAP server provides access for anonymous authentication.
   - **Bind Password**: The password for an application service account, if **Bind DN** is set.

   Provide the user search attributes.

   - **Base DN**: The point from which to start the LDAP search. For example, `OU=Users,OU=domain,DC=io`.
   - **Filter**: An optional filter to be used by the LDAP search.
   - **Username**: The LDAP attribute that contains the user ID. For example, `uid, sAMAccountName`.

   Provide the group search attributes.

   - **Base DN**: The point from which to start the LDAP search. For example, `OU=Groups,OU=domain,DC=io`.
   - **Filter**: An optional filter to be used by the LDAP search.
   - **Name Attribute**: The LDAP attribute that holds the name of the group. For example, `cn`.
   - **User Attribute**: The attribute of the user record that is used as the value of the membership attribute of the group record. For example, `distinguishedName, dn`.
   - **Group Attribute**:  The attribute of the group record that holds the user/member information. For example, `member`.

   Paste the contents of the LDAPS server CA certificate into the **Root CA** text box.

   ![Configure external Identity Provider](../images/install-v-7id-ldap.png)

1. If you are deploying to vSphere, click **Next** to go to [Select the Base OS Image](#base-os). If you are deploying to Amazon EC2 or Azure, click **Next** to go to [Register with Tanzu Mission Control](#register-tmc).

## <a id="base-os"></a> (vSphere Only) Select the Base OS Image

In the **OS Image** section, use the drop-down menu to select the OS and Kubernetes version image template to use for deploying Tanzu Kubernetes Grid VMs, and click **Next**.

The drop-down menu includes all of the image templates that are present in your vSphere instance that meet the criteria for use as Tanzu Kubernetes Grid base images. The image template must include the correct version of Kubernetes for this release of Tanzu Kubernetes Grid. If you have not already imported a suitable image template to vSphere, you can do so now without quitting the Tanzu Kubernetes Grid installer. After you import it, use the Refresh button to make it available in the drop-down menu.

   ![Select the base image template](../images/install-v-8image.png)

## <a id="register-tmc"></a> Register with Tanzu Mission Control

This section applies to all infrastructure providers, however the functionality described in this section is being rolled out in Tanzu Mission Control.

**Note** At time of publication, you can only register Tanzu Kubernetes Grid management clusters that are deployed on vSphere 6.7U3, vSphere 7.0 without vSphere with Tanzu enabled, and VMware Cloud on AWS with SDDC v1.12. You cannot register management clusters that are deployed on Azure VMware Solution, Amazon EC2, or Microsoft Azure.

For more information about registering your Tanzu Kubernetes Grid management cluster with Tanzu Mission Control, see [Register Your Management Cluster with Tanzu Mission Control](register_tmc.md).

1. In the **Registration URL** field, copy and paste the registration URL you obtained from Tanzu Mission Control.

   ![Register with Tanzu Mission Control](../images/aws-tmc-register.png)

1. If the connection is successful, you can review the configuration YAML retrieved from the URL.

1. Click **Next**.

## <a id="finalize-deployment"></a> Finalize the Deployment

This section applies to all infrastructure providers.

1. In the **CEIP Participation** section, optionally deselect the check box to opt out of the VMware Customer Experience Improvement Program.

   You can also opt in or out of the program after the deployment of the management cluster. For information about the CEIP, see [Opt in or Out of the VMware CEIP](../cluster-lifecycle/multiple-management-clusters.md#ceip) and [https://www.vmware.com/solutions/trustvmware/ceip.html](https://www.vmware.com/solutions/trustvmware/ceip.html).
1. Click **Review Configuration** to see the details of the management cluster that you have configured.

   The image below shows the configuration for a deployment to vSphere.

   ![Review the management cluster configuration](../images/review-settings-vsphere.png)

   When you click **Review Configuration**, Tanzu Kubernetes Grid populates the cluster configuration file, which is located in the `~/.tanzu/tkg/clusterconfigs` subdirectory, with the settings that you specified in the interface. You can optionally copy the cluster configuration file without completing the deployment. You can copy the cluster configuration file to another bootstrap machine and deploy the management cluster from that machine. For example, you might do this so that you can deploy the management cluster from a bootstrap machine that does not have a Web browser.

1. (Optional) Under **CLI Command Equivalent**, click the **Copy** button to copy the CLI command for the configuration that you specified.

   Copying the CLI command allows you to reuse the command at the command line to deploy management clusters with the configuration that you specified in the interface. This can be useful if you want to automate management cluster deployment.

1. (Optional) Click **Edit Configuration** to return to the installer wizard to modify your configuration.
1. Click **Deploy Management Cluster**.

   Deployment of the management cluster can take several minutes. The first run of `tanzu management-cluster create` takes longer than subsequent runs because it has to pull the required Docker images into the image store on your bootstrap machine. Subsequent runs do not require this step, so are faster. You can follow the progress of the deployment of the management cluster in the installer interface or in the terminal in which you ran `tanzu management-cluster create --ui`. If the machine on which you run `tanzu management-cluster create` shuts down or restarts before the local operations finish, the deployment will fail. If you inadvertently close the browser or browser tab in which the deployment is running before it finishes, the deployment continues in the terminal.
   
   **NOTE**: The screen capture below shows the deployment status page in Tanzu Kubernetes Grid v1.3.1.

   ![Monitor the management cluster deployment](../images/mgmt-cluster-deployment.png)

## <a id="what-next"></a> What to Do Next

- If you enabled identity management on the management cluster, you must perform post-deployment configuration steps to allow users to access the management cluster. For more information, see [Configure Identity Management After Management Cluster Deployment](configure-id-mgmt.md).
- For information about what happened during the deployment of the management cluster and how to connect `kubectl` to the management cluster, see [Examine the Management Cluster Deployment](verify-deployment.md).
- If you need to deploy more than one management cluster, on any or all of vSphere, Azure, and Amazon EC2, see [Manage Your Management Clusters](../cluster-lifecycle/multiple-management-clusters.md). This topic also provides information about how to add existing management clusters to your CLI instance, obtain credentials, scale and delete management clusters, add namespaces, and how to opt in or out of the CEIP.
