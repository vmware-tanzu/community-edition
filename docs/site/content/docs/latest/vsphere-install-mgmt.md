# Deploy a Management Cluster to vSphere

{{% include "/docs/assets/step1.md" %}}
### Step 1: IaaS Provider

1. Enter the IP address or fully qualified domain name (FQDN) for the vCenter Server instance on which to deploy the management cluster.  Tanzu does not support IPv6 addresses. This is because upstream Kubernetes only provides alpha support for IPv6. 
2. Enter the vCenter Single Sign On username and password for a user account that has the required privileges for Tanzu operation, and click **Connect**.
3. Verify the SSL thumbprint of the vCenter Server certificate and click **Continue** if it is valid.    For information about how to obtain the vCenter Server certificate thumbprint, see [Obtain vSphere Certificate Thumbprints](ref-vsphere.md#certificates).
4. Select the datacenter in which to deploy the management cluster from the **Datacenter** drop-down menu.
5. Paste the contents of your SSH public key into the text box and click **Next**.
## Step 2: Management Cluster Settings

1. In the **Management Cluster Settings** section, select an instance size for either **Development** or
   **Production**. If you select **Development**, the installer deploys a management cluster with a single control
   plane node. If you select **Production**, the installer deploys a highly available management cluster with three
   control plane nodes. Use the **Instance type** drop-down menu to select from different combinations of CPU, RAM,
   and storage for the control plane node VM or VMs.
1. (Optional) Enter a name for your management cluster under **Management Cluster Name**. If you do not specify a
   name, the installer generates a unique name. The name must end with a letter, not a numeric character, and must
   be compliant with DNS hostname requirements as outlined in [RFC 952](https://tools.ietf.org/html/rfc952) and
   amended in [RFC 1123](https://tools.ietf.org/html/rfc1123).
1. The **Machine Health Check** option provides node health monitoring and node auto-repair on the clusters that you
   deploy with this management cluster. [Machine Health Checks](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/machine-health-check.html#machinehealthcheck)
   are enabled by default. You can enable or disable Machine Health Checks on clusters after deployment by using the
   CLI. For instructions, see [Configure Machine Health Checks for Tanzu Kubernetes Clusters](../cluster-lifecycle/configure-health-checks.md).

1. Select the **Control Plane Endpoint Provider**. This can be either the default [kube-vip](https://kube-vip.io/),
   or if available, you may use an
   [NSX Advanced Load Balancer](https://www.vmware.com/products/nsx-advanced-load-balancer.html).
1. Under **Control Plane Endpoint**, enter a static virtual IP address or FQDN for API requests to the management
   cluster. Ensure that this IP address is not in your DHCP range, but is in the same subnet as the DHCP range. If
   you mapped an FQDN to the VIP address, you can specify the FQDN instead of the VIP address. For more information,
   see [Static VIPs and Load Balancers for vSphere](vsphere.md#load-balancer).
1. Under **Worker Node Instance Type**, select the configuration for the worker node VM. If you select an instance
   type in the **Production** tile, the instance type that you select is automatically selected for the
   **Worker Node Instance Type**. If necessary, you can change this.
1. Checking the **Enable Audit Logging** checkbox will enable additional audit logging to be captured.

### Step 3: VMware NSX Advanced Load Balancer

VMware NSX Advanced Load Balancer provides an L4 load balancing solution for vSphere. NSX Advanced Load Balancer includes a Kubernetes operator that integrates with the Kubernetes API to manage the lifecycle of load balancing and ingress resources for workloads. To use NSX Advanced Load Balancer, you must first deploy it in your vSphere environment. For information, see [Install VMware NSX Advanced Load Balancer on a vSphere Distributed Switch](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3/vmware-tanzu-kubernetes-grid-13/GUID-mgmt-clusters-install-nsx-adv-lb.html).

In the optional **VMware NSX Advanced Load Balancer** section, you can configure Tanzu to use NSX Advanced Load Balancer. By default all workload clusters will use the load balancer.

1. For **Controller Host**, enter the IP address or FQDN of the Controller VM.
1. Enter the username and password that you set for the Controller host when you deployed it, and click **Verify Credentials**.
1. Use the **Cloud Name** drop-down menu to select the cloud that you created in your NSX Advanced Load Balancer deployment.

   For example, `Default-Cloud`.
1. Use the **Service Engine Group Name** drop-down menu to select a Service Engine Group.

   For example, `Default-Group`.
1. For **VIP Network Name**, use the drop-down menu to select the name of the network where the load balancer floating IP Pool resides.

   The VIP network for NSX Advanced Load Balancer must be present in the same vCenter Server instance as the Kubernetes network that Tanzu uses. This allows NSX Advanced Load Balancer to discover the Kubernetes network in vCenter Server and to deploy and configure Service Engines. The drop-down menu is present in Tanzu v1.3.1 and later. In v1.3.0, you enter the name manually.

   You can see the network in the **Infrastructure** > **Networks** view of the NSX Advanced Load Balancer interface.
1. For **VIP Network CIDR**, use the drop-down menu to select the CIDR of the subnet to use for the load balancer VIP.

   This comes from one of the VIP Network's configured subnets. You can see the subnet CIDR for a particular network in the **Infrastructure** > **Networks** view of the NSX Advanced Load Balancer interface.

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
### Step 4: Metadata
{{% include "/docs/assets/metadata.md" %}}

### Step 5: Resources
In the **Resources** section, select vSphere resources for the management cluster to use, and click **Next**.

   * Select the VM folder in which to place the management cluster VMs.
   * Select a vSphere datastore for the management cluster to use.
   * Select the cluster, host, or resource pool in which to place the management cluster.

   If appropriate resources do not already exist in vSphere, without quitting the Tanzu installer, go to vSphere to create them. Then click the refresh button so that the new resources can be selected.

### Step 6: Kubernetes Network

1.   Under **Network Name**, select a vSphere network to use as the Kubernetes service network.  

2. (Optional) To send outgoing HTTP(S) traffic from the management cluster to a proxy, toggle **Enable Proxy Settings** and follow the instructions below to enter your proxy information. Tanzu applies these settings to kubelet, containerd, and the control plane. You can choose to use one proxy for HTTP traffic and another proxy for HTTPS traffic or to use the same proxy for both HTTP and HTTPS traffic.  

    - To add your HTTP proxy information: Under **HTTP Proxy URL**, enter the URL of the proxy that handles HTTP requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`.  If the proxy requires authentication, under **HTTP Proxy Username** and **HTTP Proxy Password**, enter the username and password to use to connect to your HTTP proxy.

    - To add your HTTPS proxy information: If you want to use the same URL for both HTTP and HTTPS traffic, select **Use the same configuration for https proxy**.  If you want to use a different URL for HTTPS traffic, enter the URL of the proxy that handles HTTPS requests. The URL must start with `http://`. For example, `http://myproxy.com:1234`. If the proxy requires authentication, under **HTTPS Proxy Username** and **HTTPS Proxy Password**, enter the username and password to use to connect to your HTTPS proxy.

    - Under **No proxy**, enter a comma-separated list of network CIDRs or hostnames that must bypass the HTTP(S) proxy. For example, `noproxy.yourdomain.com,192.168.0.0/24`. You must enter the CIDR of the vSphere network that you selected under **Network Name**. The vSphere network CIDR includes the IP address of your **Control Plane Endpoint**. If you entered an FQDN under **Control Plane Endpoint**, add both the FQDN and the vSphere network CIDR to **No proxy**. Internally, Tanzu appends `localhost`, `127.0.0.1`, the values of **Cluster Pod CIDR** and **Cluster Service CIDR**, `.svc`, and `.svc.cluster.local` to the list that you enter in this field.
      

    **Important:** If the management cluster VMs need to communicate with external services and infrastructure endpoints in your Tanzu environment, ensure that those endpoints are reachable by the proxies that you configured above or add them to **No proxy**. Depending on your environment configuration, this may include, but is not limited to, your OIDC or LDAP server, and Harbor.

### Step 7: Identity Management
{{% include "/docs/assets/identity-management.md" %}}

### Step 8: OS Image

In the **OS Image** section, use the drop-down menu to select the OS and Kubernetes version image template to use for deploying Tanzu VMs, and click **Next**.

The drop-down menu includes all of the image templates that are present in your vSphere instance that meet the criteria for use as Tanzu base images. The image template must include the correct version of Kubernetes for this release of Tanzu. If you have not already imported a suitable image template to vSphere, you can do so now without quitting the Tanzu installer. After you import it, use the Refresh button to make it available in the drop-down menu.

### Step 7: Register TMC
{{% include "/docs/assets/register_tmc.md" %}}

{{% include "/docs/assets/final-step.md" %}}





