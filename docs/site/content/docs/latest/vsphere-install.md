WIP DRAFT WIP DRAFT

# Deploy a cluster to vSphere

{{% include "/docs/assets/step1.md" %}}
## Step 1: IaaS Provider

1. Enter the IP address or fully qualified domain name (FQDN) for the vCenter Server instance on which to deploy the management cluster.  Tanzu Kubernetes Grid does not support IPv6 addresses. This is because upstream Kubernetes only provides alpha support for IPv6. 
2. Enter the vCenter Single Sign On username and password for a user account that has the required privileges for Tanzu Kubernetes Grid operation, and click **Connect**.
<!--  ![Configure the connection to vSphere](../images/install-v-1iaas.png)--> 
3. Verify the SSL thumbprint of the vCenter Server certificate and click **Continue** if it is valid.    For information about how to obtain the vCenter Server certificate thumbprint, see [Obtain vSphere Certificate Thumbprints](ref-vsphere.md#certificates).
<!--[Verify vCenter Server certificate thumbprint](../images/vsphere-thumprint.png)-->
<!--1. If you are deploying a management cluster to a vSphere 7 instance, confirm whether or not you want to proceed with the deployment.   

   On vSphere 7, the vSphere with Tanzu option includes a built-in supervisor cluster that works as a management cluster and provides a better experience than a separate management cluster deployed by Tanzu Kubernetes Grid.  Deploying a Tanzu Kubernetes Grid management cluster to vSphere 7 when vSphere with Tanzu is not enabled is supported, but the preferred option is to enable vSphere with Tanzu and use the Supervisor Cluster. VMware Cloud on AWS and Azure VMware Solution do not support a supervisor cluster, so you need to deploy a management cluster.
   For information, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).

   To reflect the recommendation to use vSphere with Tanzu when deploying to vSphere 7, the Tanzu Kubernetes Grid installer behaves as follows:

      - **If vSphere with Tanzu is enabled**, the installer informs you that deploying a management cluster is not possible, and exits.
      - **If vSphere with Tanzu is not enabled**, the installer informs you that deploying a Tanzu Kubernetes Grid management cluster is possible but not recommended, and presents a choice:
          - **Configure vSphere with Tanzu** opens the vSphere Client so you can configure your Supervisor Cluster as described in [Configuring and Managing a Supervisor Cluster](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-21ABC792-0A23-40EF-8D37-0367B483585E.html) in the vSphere documentation.
          - **Deploy TKG Management Cluster** allows you to continue deploying a management cluster, against recommendation for vSphere 7, but as required for VMware Cloud on AWS and Azure VMware Solution. When using vSphere 7, the preferred option is to enable vSphere with Tanzu and use the built-in Supervisor Cluster instead of deploying a Tanzu Kubernetes Grid management cluster.

   ![Deploy management cluster to vSphere 7](../images/vsphere7-detected.png)-->
4. Select the datacenter in which to deploy the management cluster from the **Datacenter** drop-down menu.
5. Paste the contents of your SSH public key into the text box and click **Next**.
<!--   ![Select datacenter and provide SSH public key](../images/dc-ssh-vsphere.png)-->

## Step 2: Management Cluster

1. In the **Management Cluster Settings** section, select an instance size for either **Development** or **Production**. If you select **Development**, the installer deploys a management cluster with a single control plane node. If you select **Production**, the installer deploys a highly available management cluster with three control plane nodes. Use the **Instance type** drop-down menu to select from different combinations of CPU, RAM, and storage for the control plane node VM or VMs. 
<!--Choose the configuration for the control plane node VMs depending on the expected workloads that it will run. For example, some workloads might require a large compute capacity but relatively little storage, while others might require a large amount of storage and less compute capacity.-->    
<!--If you plan on registering the management cluster with Tanzu Mission Control, ensure that your Tanzu Kubernetes clusters meet the requirements listed in [Requirements for Registering a Tanzu Kubernetes Cluster with Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-3AE5F733-7FA7-4B34-8935-C25D41D15EF9.html) in the Tanzu Mission Control documentation.
![Select the control plane node configuration](../images/configure-control-plane.png)-->
2. (Optional) Enter a name for your management or stand-alone cluster. If you do not specify a name, the installer generates a unique name. If you do specify a name, that name must end with a letter, not a numeric character, and must be compliant with DNS hostname requirements as outlined in [RFC 952](https://tools.ietf.org/html/rfc952) and amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).
3. Under **Worker Node Instance Type**, select the configuration for the worker node VM.  If you select an instance type in the **Production** tile, the instance type that you select is automatically selected for the **Worker Node Instance Type**. If necessary, you can change this. 
4. The `MachineHealthCheck` option provides node health monitoring and node auto-repair on the clusters that you deploy with this management cluster. [`MachineHealthCheck`](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/machine-health-check.html#machinehealthcheck) is enabled by default. You can enable or disable `MachineHealthCheck` on clusters after deployment by using the CLI. For instructions, see [Configure Machine Health Checks for Tanzu Kubernetes Clusters](../cluster-lifecycle/configure-health-checks.md). **<ENG TEAM  - DO WE NEED TO BRING IN THIS TOPIC FROM TKG DOCS>**
5. Under **Control Plane Endpoint**, enter a static virtual IP address or FQDN for API requests to the management cluster. Ensure that this IP address is not in your DHCP range, but is in the same subnet as the DHCP range. If you mapped an FQDN to the VIP address, you can specify the FQDN instead of the VIP address. For more information, see [Static VIPs and Load Balancers for vSphere](vsphere.md#load-balancer). <!--different for vspehere and aws--> **ENG TEAM  - DO WE NEED TO BRING IN THIS TOPIC FROM TKG DOCS**
<!--![Select the cluster configuration](../images/configure-cluster.png)-->
6. To complete the configuration of the **Management Cluster Settings** section, do one of the following:
   * If you created a new VPC in the **VPC for AWS** section, click **Next**.
   * If you selected an existing VPC in the **VPC for AWS** section, use the **VPC public subnet** and **VPC private subnet** drop-down menus to select existing subnets on the VPC and click **Next**. 
<!--        ![Set the VPC subnets](../images/aws-subnets.png)-->

## Step 3: Configure VMware NSX Advanced Load Balancer

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

   This comes from one of the VIP Network's configured subnets. You can see the subnet CIDR for a particular network in the **Infrastructure** > **Networks** view of the NSX Advanced Load Balancer interface. The drop-down menu is present in Tanzu Kubernetes Grid v1.3.1 and later. In v1.3.0, you enter the CIDR manually.<!--Which version are we consuming?>

1. Paste the contents of the Certificate Authority that is used to generate your Controller Certificate into the **Controller Certificate Authority** text box.

   If you have a self-signed Controller Certificate, the Certificate Authority is the same as the Controller Certificate.
1. (Optional) Enter one or more cluster labels to identify clusters on which to selectively enable NSX Advanced Load Balancer or to customize NSX Advanced Load Balancer Settings per group of clusters.

   By default, all clusters that you deploy with this management cluster will enable NSX Advanced Load Balancer. All clusters will share the same VMware NSX Advanced Load Balancer Controller, Cloud, Service Engine Group, and VIP Network as you entered previously. This cannot be changed later. To only enable the load balancer on a subset of clusters, or to preserve the ability to customize NSX Advanced Load Balancer settings for a group of clusters, add labels in the format `key: value`. For example `team: tkg`.

   This is useful in the following scenarios:

   - You want to configure different sets of workload clusters to different Service Engine Groups to implement isolation or to support more Service type Load Balancers than one Service Engine Group's capacity.
   - You want to configure different sets of workload clusters to different Clouds because they are deployed in separate sites.

   **NOTE**: Labels that you define here will be used to create a label selector. Only workload cluster `Cluster` objects that have the matching labels will have the load balancer enabled. As a consequence, you are responsible for making sure that the workload cluster's `Cluster` object has the corresponding labels. For example, if you use `team: tkg`, to enable the load balancer on a workload cluster, you will need to perform the following steps after deployment of the management cluster:

   1. Set `kubectl` to the management cluster's context.

      ```sh
      kubectl config set-context management-cluster@admin
      ```

   1. Label the `Cluster` object of the corresponding workload cluster with the labels defined. If you define multiple key-values, you need to apply all of them.     

      ```sh
      kubectl label cluster <cluster-name> team=tkg
      ```

<!--![Configure NSX Advanced Load Balancer](../images/install-v-3nsx.png)-->

1. Click **Next** to configure metadata.

## Step 4: Metadata
{{% include "/docs/assets/metadata.md" %}}

## Step 6: Kubernetes Network
<!-- note to self: right now I can't figure a good way to turn this into an include that could be reused across amazon and vsphere as there is too much mixed up information about both in it, so it will be added manually to each and cleaned up appropriately - so this will need to be copied into both vsphere and amazon topics-->

1.   Under **Network Name**, select a vSphere network to use as the Kubernetes service network.  <!--different for vspehere and aws-->
<!--![Configure the Kubernetes service network](../images/install-v-6k8snet.png) -->   
2. (Optional) To send outgoing HTTP(S) traffic from the management cluster to a proxy, toggle **Enable Proxy Settings** and follow the instructions below to enter your proxy information. Tanzu Kubernetes Grid applies these settings to kubelet, containerd, and the control plane. You can choose to use one proxy for HTTP traffic and another proxy for HTTPS traffic or to use the same proxy for both HTTP and HTTPS traffic.  

    - To add your HTTP proxy information: Under **HTTP Proxy URL**, enter the URL of the proxy that handles HTTP requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`.  If the proxy requires authentication, under **HTTP Proxy Username** and **HTTP Proxy Password**, enter the username and password to use to connect to your HTTP proxy.

    - To add your HTTPS proxy information: If you want to use the same URL for both HTTP and HTTPS traffic, select **Use the same configuration for https proxy**.  If you want to use a different URL for HTTPS traffic, enter the URL of the proxy that handles HTTPS requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`. If the proxy requires authentication, under **HTTPS Proxy Username** and **HTTPS Proxy Password**, enter the username and password to use to connect to your HTTPS proxy.

    - Under **No proxy**, enter a comma-separated list of network CIDRs or hostnames that must bypass the HTTP(S) proxy. For example, `noproxy.yourdomain.com,192.168.0.0/24`. <!--different for vspehere and aws-->You must enter the CIDR of the vSphere network that you selected under **Network Name**. The vSphere network CIDR includes the IP address of your **Control Plane Endpoint**. If you entered an FQDN under **Control Plane Endpoint**, add both the FQDN and the vSphere network CIDR to **No proxy**. Internally, Tanzu Kubernetes Grid appends `localhost`, `127.0.0.1`, the values of **Cluster Pod CIDR** and **Cluster Service CIDR**, `.svc`, and `.svc.cluster.local` to the list that you enter in this field.
      

    **Important:** If the management cluster VMs need to communicate with external services and infrastructure endpoints in your Tanzu Kubernetes Grid environment, ensure that those endpoints are reachable by the proxies that you configured above or add them to **No proxy**. Depending on your environment configuration, this may include, but is not limited to, your OIDC or LDAP server, and Harbor.

3.  In the **Resources** section, select vSphere resources for the management cluster to use, and click **Next**. <!--different for vspehere and aws--> 

   - Select the VM folder in which to place the management cluster VMs.
   - Select a vSphere datastore for the management cluster to use.
   - Select the cluster, host, or resource pool in which to place the management cluster.

   If appropriate resources do not already exist in vSphere, without quitting the Tanzu Kubernetes Grid installer, go to vSphere to create them. Then click the refresh button so that the new resources can be selected.

<!--[Select vSphere resources](../images/install-v-5resources.png)-->

## Step 7: Identity Management
{{% include "/docs/assets/identity-management.md" %}}

## Step 8: OS Image

In the **OS Image** section, use the drop-down menu to select the OS and Kubernetes version image template to use for deploying Tanzu Kubernetes Grid VMs, and click **Next**.

The drop-down menu includes all of the image templates that are present in your vSphere instance that meet the criteria for use as Tanzu Kubernetes Grid base images. The image template must include the correct version of Kubernetes for this release of Tanzu Kubernetes Grid. If you have not already imported a suitable image template to vSphere, you can do so now without quitting the Tanzu Kubernetes Grid installer. After you import it, use the Refresh button to make it available in the drop-down menu.

<!--   ![Select the base image template](../images/install-v-8image.png)-->



{{% include "/docs/assets/final-step.md" %}}
