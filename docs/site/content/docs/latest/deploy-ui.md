# Deploy Management and Stand-alone Clusters with the Installer Interface

This topic describes how to use the Tanzu Kubernetes Grid installer interface to deploy a management or stand-alone cluster. The installer interface launches in a browser and takes you through steps to configure the management or stand-alone cluster. The input values are saved in a cluster configuration file. After you confirm your input values, the installer saves them to: `~/.tanzu/tkg/clusterconfigs/cluster-config.yaml`. 

## Before you begin

- Make sure that you have met all of the requirements and followed all of the procedures in [Install the Tanzu CLI](../installation-cli). 

- Make sure that you have met all of the requirements listed [Prepare to Deploy Management Clusters to Amazon EC2](../prepare-deployment).

- You have met the following installer prerequisites:

   - NTP is running on the bootstrap machine on which you are running `tanzu management-cluster create` and on the hypervisor.
   - A DHCP server is available.
   - The CLI can connect to the location from which it pulls the required images.
   - Docker is running.

- By default Tanzu Kubernetes Grid saves the `kubeconfig` for all management clusters in the `~/.kube-tkg/config` file. If you want to save the `kubeconfig` file to a different location, set the `KUBECONFIG` environment variable before running the installer, for example:
  ```
   KUBECONFIG=/path/to/mc-kubeconfig.yaml
   ```

<!--- For production deployments, it is strongly recommended to enable identity management for your clusters. For information about the preparatory steps to perform before you deploy a management cluster, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).
- If you want to register your management cluster with Tanzu Mission Control, follow the procedure in [Register Your Management Cluster with Tanzu Mission Control](register_tmc.md).
- If you are deploying clusters in an internet-restricted environment to either vSphere or Amazon EC2, you must also perform the steps in [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md).-->

<!--- **NOTE**: On vSphere with Tanzu, you do not need to deploy a management cluster. See [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).-->



## Procedure

1. On the machine on which you downloaded and installed the Tanzu CLI, run the `tanzu management-cluster create` command with the `--ui` option.

   ```
   tanzu management-cluster create --ui
   ```
   If the prerequisites are met, the installer interface opens locally, at http://127.0.0.1:8080 in your default browser. To change where the installer interface runs, including running it on a different machine from the `tanzu` CLI, use the following parameters.

   - `--browser` specifies the local browser to open the interface in. Supported values are `chrome`, `firefox`, `safari`, `ie`, `edge`, or `none`. Use `none` with `--bind` to run the interface on a different machine.
   - `--bind` specifies the IP address and port to serve the interface from. For example, if another process is already using http://127.0.0.1:8080, use `--bind` to serve the interface from a different local port.
   
   Example:  
   ```
   tanzu management-cluster create --ui --bind 192.168.1.87:5555 --browser none
   ```  



1. Click the **Deploy** button for **VMware vSphere**, **Amazon EC2**, or **??Stand-alone??**.

   ![Tanzu Kubernetes Grid installer interface welcome page with Deploy to vSphere button](../images/deploy-management-cluster.png)


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
