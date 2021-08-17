




test
test
test


## <a id="network"></a> Configure the Kubernetes Network and Proxies
<!-- note to self: right now I can't figure a good way to turn this into an include that could be reused across amazon and vsphere, so it will be  added manually to each and cleaned up appropriately - so this will need to be copied into both vsphere and amazon topics-->
This section applies to all infrastructure providers.

1. In the **Kubernetes Network** section, configure the networking for Kubernetes services and click **Next**.

   * **(vSphere only)** Under **Network Name**, select a vSphere network to use as the Kubernetes service network.
   * Review the **Cluster Service CIDR** and **Cluster Pod CIDR** ranges. If the recommended CIDR ranges of `100.64.0.0/13` and `100.96.0.0/11` are unavailable, update the values under **Cluster Service CIDR** and **Cluster Pod CIDR**.

<!--   ![Configure the Kubernetes service network](../images/install-v-6k8snet.png)-->

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


      **Important:** If the management cluster VMs need to communicate with external services and infrastructure endpoints in your Tanzu Kubernetes Grid environment, ensure that those endpoints are reachable by the proxies that you configured above or add them to **No proxy**. Depending on your environment configuration, this may include, but is not limited to, your OIDC or LDAP server, Harbor, and in the case of vSphere, NSX-T and NSX Advanced Load Balancer.





## <a id="id-mgmt"></a> Configure Identity Management

For information about how Tanzu Kubernetes Grid implements identity management, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).
<!-- ??I don't know if this is something we want to reference or if we need to supply our own??? I presume this full section needs to be reworked for TCE -->

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

   When you click **Review Configuration**, Tanzu Kubernetes Grid populates the cluster configuration file, which is located in the `~/.config/tanzu/tkg/clusterconfigs` subdirectory, with the settings that you specified in the interface. You can optionally copy the cluster configuration file without completing the deployment. You can copy the cluster configuration file to another bootstrap machine and deploy the management cluster from that machine. For example, you might do this so that you can deploy the management cluster from a bootstrap machine that does not have a Web browser.

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
