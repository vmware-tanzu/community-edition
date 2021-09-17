# Deploy a Management Cluster to Azure

{{% include "/docs/assets/step1.md" %}}


### Step 1: IaaS Provider

1. In the **IaaS Provider** section, enter the **Tenant ID**, **Client ID**, **Client Secret**, and **Subscription ID** values for your Azure account.  You recorded these values when you registered an Azure app and created a secret for it using the Azure Portal. For more information, see the [Register Tanzu Community Edition as an Azure Client App](azure-mgmt/#a-idtkg-appa-register-tanzu-community-edition-as-an-azure-client-app) topic.


1. Select your **Azure Environment**, either **Public Cloud** or **US Government Cloud**. You can specify other environments by deploying from a configuration file and setting `AZURE_ENVIRONMENT`.
1. Click **Connect**. The installer verifies the connection and changes the button label to **Connected**.
1. Select the Azure region in which to deploy the management cluster.
1. Paste the contents of your SSH public key, such as `.ssh/id_rsa.pub`, into the text box.
1. Under **Resource Group**, select either the **Select an existing resource group** or the **Create a new resource group** radio button.

    - If you select **Select an existing resource group**, use the drop-down menu to select the group, then click **Next**.
    - If you select **Create a new resource group**, enter a name for the new resource group and then click **Next**.

1. In the **VNET for Azure** section, select either the **Create a new VNET on Azure** or the **Select an existing VNET** radio button.
    - If you select **Create a new VNET on Azure**, use the drop-down menu to select the resource group in which to create the VNET and provide the following:
       - A name and a CIDR block for the VNET. The default is `10.0.0.0/16`.
       - A name and a CIDR block for the control plane subnet. The default is `10.0.0.0/24`.
       - A name and a CIDR block for the worker node subnet. The default is `10.0.1.0/24`.
    - If you select **Select an existing VNET**, use the drop-down menus to select the resource group in which the VNET is located, the VNET name, the control plane and worker node subnets, and then click **Next**.
    - To make the management cluster private, enable the **Private Azure Cluster** checkbox. By default, Azure management and workload clusters are public. But you can also configure them to be private, which means their API server uses an Azure internal load balancer (ILB) and is therefore only accessible from within the cluster’s own VNET or peered VNETs. For more information, see the [Azure Private Clusters](azure-wl-template/#a-idprivatea-azure-private-clusters) topic.

### Step 2: Management Cluster Settings

1. In the **Management Cluster Settings** section, select an instance size for either **Development** or **Production**. If you select **Development**, the installer deploys a management cluster with a single control plane node. If you select **Production**, the installer deploys a highly available management cluster with three control plane nodes. Use the **Instance type** drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs.  The minimum configuration is 2 CPUs and 8 GB memory. The list of compatible instance types varies in different regions. For information about the configurations of the different sizes of node instances for Azure, see [Sizes for virtual machines in Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/sizes).
2. (Optional) Enter a name for your management cluster. If you do not specify a name, the installer generates a unique name. The name must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements as outlined in [RFC 952](https://tools.ietf.org/html/rfc952) and amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).
3. Under **Worker Node Instance Type**, select the configuration for the worker node VM.  If you select an instance type in the **Production** tile, the instance type that you select is automatically selected for the **Worker Node Instance Type**. If necessary, you can change this.
4. The MachineHealthCheck option provides node health monitoring and node auto-repair on the clusters that you deploy with this management cluster. [MachineHealthCheck](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/machine-health-check.html#machinehealthcheck) is enabled by default. You can enable or disable MachineHealthCheck on clusters after deployment by using the CLI. For instructions, see [Configure Machine Health Checks for Tanzu Kubernetes Clusters](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3/vmware-tanzu-kubernetes-grid-13/GUID-cluster-lifecycle-configure-health-checks.html).

### Step 3: Metadata
{{% include "/docs/assets/metadata.md" %}}

### Step 4: Kubernetes Network

1.  Review the **Cluster Service CIDR** and **Cluster Pod CIDR** ranges. If the recommended CIDR ranges of `100.64.0.0/13` and `100.96.0.0/11` are unavailable, update the values under **Cluster Service CIDR** and **Cluster Pod CIDR**.

2. (Optional) To send outgoing HTTP(S) traffic from the management cluster to a proxy, toggle **Enable Proxy Settings** and follow the instructions below to enter your proxy information. Tanzu applies these settings to kubelet, containerd, and the control plane. You can choose to use one proxy for HTTP traffic and another proxy for HTTPS traffic or to use the same proxy for both HTTP and HTTPS traffic.

    - To add your HTTP proxy information: Under **HTTP Proxy URL**, enter the URL of the proxy that handles HTTP requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`.  If the proxy requires authentication, under **HTTP Proxy Username** and **HTTP Proxy Password**, enter the username and password to use to connect to your HTTP proxy.

    - To add your HTTPS proxy information: If you want to use the same URL for both HTTP and HTTPS traffic, select **Use the same configuration for https proxy**.  If you want to use a different URL for HTTPS traffic, enter the URL of the proxy that handles HTTPS requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`. If the proxy requires authentication, under **HTTPS Proxy Username** and **HTTPS Proxy Password**, enter the username and password to use to connect to your HTTPS proxy.

    - Under **No proxy**, enter a comma-separated list of network CIDRs or hostnames that must bypass the HTTP(S) proxy. For example, `noproxy.yourdomain.com,192.168.0.0/24`. Internally, Tanzu appends `localhost`, `127.0.0.1`, your VPC CIDR, **Cluster Pod CIDR**, and **Cluster Service CIDR**, `.svc`, `.svc.cluster.local`, and `169.254.0.0/16` to the list that you enter in this field.


    **Important:** If the management cluster VMs need to communicate with external services and infrastructure endpoints in your Tanzu environment, ensure that those endpoints are reachable by the proxies that you configured above or add them to **No proxy**. Depending on your environment configuration, this may include, but is not limited to, your OIDC or LDAP server, and Harbor.

### Step 5: Identity Management
{{% include "/docs/assets/identity-management.md" %}}

### Step 7: Register TMC
{{% include "/docs/assets/register_tmc.md" %}}

{{% include "/docs/assets/final-step.md" %}}


