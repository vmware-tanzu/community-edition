# Tanzu CLI Configuration File Variable Reference

This reference lists all the variables that you can specify to provide configuration options to the Tanzu CLI.

To set these variables in a YAML configuration file, leave a space between the colon (:) and the variable value. For example:

```
CLUSTER_NAME: my-cluster
```

Line order in the configuration file does not matter. Options are presented here in alphabetical order.

## <a id="all_iaases"></a> Common Variables for All Infrastructure Providers

This section lists variables that are common to all infrastructure providers. These variables may apply to management clusters, Tanzu Kubernetes clusters, or both. For more information, see [Configure Basic Management Cluster Creation Information](mgmt-clusters/create-config-file.md#basic) in *Create a Management Cluster Configuration File*. For the variables that are specific to workload clusters, see [Deploy Tanzu Kubernetes Clusters](tanzu-k8s-clusters/deploy.md).

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
 <tbody>
  <tr>
    <td><code>CLUSTER_CIDR</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional, set if you want to override the default value</strong>. The CIDR range to use for pods. By default, this range is set to <code>100.96.0.0/11</code>. Change the default value only if the recommended range is unavailable.</td>
  </tr>  
  <tr>
    <td><code>CLUSTER_NAME</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td>This name must comply with DNS hostname requirements as outlined in <a href="https://tools.ietf.org/html/rfc952">RFC 952</a> and amended in <a href="https://tools.ietf.org/html/rfc1123">RFC 1123</a>, and must be 42 characters or less.<br />
    For workload clusters, this setting is overridden by the <code>CLUSTER_NAME</code> argument passed to <code>tanzu cluster create</code>.<br />
    For management clusters, if you do not specify <code>CLUSTER_NAME</code>, a unique name is generated.</td>
  </tr>
  <tr>
    <td><code>CLUSTER_PLAN</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. Set to <code>dev</code>, <code>prod</code>, or a custom plan as exemplified in <a href="tanzu-k8s-clusters/config-plans.md#nginx">New Plan <code>nginx</code></a>.<br />
    The <code>dev</code> plan deploys a cluster with a single control plane node. The <code>prod</code> plan deploys a highly available cluster with three control plane nodes.</td>
  </tr>
  <tr>
    <td><code>CNI</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Optional, set if you want to override the default value. Do not override the default value for management clusters</strong>. Container network interface. By default, <code>CNI</code> is set to <code>antrea</code>. If you want to customize your Antrea configuration, see <a href="#antrea">Antrea CNI Configuration</a> below. For Tanzu Kubernetes clusters, you can set <code>CNI</code> to <code>antrea</code>, <code>calico</code>, or <code>none</code>. Setting <code>none</code> allows you to provide your own CNI. For more information about CNI options, see <a href="tanzu-k8s-clusters/networking.md#nondefault-cni.md">Deploy a Cluster with a Non-Default CNI</a>.</td>
  </tr>
  <tr>
      <td><code>ENABLE_AUDIT_LOGGING</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to override the default value</strong>. Audit logging for the Kubernetes API server. The default value is <code>false</code>. To enable audit logging, set the variable to <code>true</code>. Tanzu Kubernetes Grid writes these logs to <code>/var/log/kubernetes/audit.log</code>. For more information, see <a href="troubleshooting-tkg/audit-logging.md">Audit Logging</a>.</td>
  </tr>
  <tr>
    <td><code>ENABLE_AUTOSCALER</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>false</code>. If set to <code>true</code>, you must include additional variables. </td>
  </tr>
  <tr>
    <td><code>ENABLE_CEIP_PARTICIPATION</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>true</code>. <code>false</code> opts out of the VMware Customer Experience Improvement Program. You can also opt in or out of the program after deploying  the management cluster. For information, see <a href="cluster-lifecycle/multiple-management-clusters.md#ceip">Opt in or Out of the VMware CEIP</a> and <a href="https://www.vmware.com/solutions/trustvmware/ceip.html">Customer Experience Improvement Program ("CEIP")</a>.</td>
  </tr>
  <tr>
    <td><code>ENABLE_DEFAULT_STORAGE_CLASS</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>true</code>. For information about storage classes, see [Create Persistent Volumes with Storage Classes](tanzu-k8s-clusters/storage.md).</td>
  </tr>
  <tr>
    <td><code>ENABLE_MHC</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>true</code>. See <a href="#mhc">Machine Health Checks</a> below.</td>
  </tr>
  <tr>
    <td><code>IDENTITY_MANAGEMENT_TYPE</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. Set either <code>oidc</code> or <code>ldap</code>. Additional OIDC or LDAP settings are required. For more information, see <a href="#identity-management">Identity Providers</a> below. Set <code>none</code> to disable identity management. It is strongly recommended to enable identity management for production deployments. </td>
  </tr>
  <tr>
    <td><code>INFRASTRUCTURE_PROVIDER</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. Set to <code>vsphere</code>, <code>aws</code>, or <code>azure</code>.</td>
  </tr>
  <tr>
    <td><code>NAMESPACE</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Optional, set if you want to override the default value</strong>. By default, Tanzu Kubernetes Grid deploys Tanzu Kubernetes clusters to the <code>default</code> namespace.</td>
  </tr>
  <tr>
    <td><code>SERVICE_CIDR</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional, set if you want to override the default value</strong>. The CIDR range to use for the Kubernetes services. By default, this range is set to <code>100.64.0.0/13</code>. Change this value only if the recommended range is unavailable.</td>
  </tr>
   <tr>
    <td><code>TMC_REGISTRATION_URL</code></td>
      <td>&#10004;</td>
      <td>&#10006;</td>
      <td><strong>Optional</strong>. Set if you want to register your management cluster with Tanzu Mission Control. For more information, see <a href="mgmt-clusters/register_tmc.md">Register Your Management Cluster with Tanzu Mission Control</a>.</td>
  </tr>
</tbody>
</table>

### <a id="identity-management-oidc"></a> Identity Providers - OIDC

If you set `IDENTITY_MANAGEMENT_TYPE: oidc`, set the following variables to configure an OIDC identity provider. For more information, see [Configure Identity Management](mgmt-clusters/create-config-file.md#identity-mgmt) in *Create a Management Cluster Configuration File*.

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
<tbody>
  <tr>
    <td><code>IDENTITY_MANAGEMENT_TYPE</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td>Enter <code>oidc</code>.</td>
  </tr>
 <tr>
    <td><code>CERT_DURATION</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. Default <code>2160h</code>. Set this variable if you configure Pinniped and Dex to use self-signed certificates managed by <code>certifcate-manager</code>.</td>
  </tr>
 <tr>
    <td><code>CERT_RENEW_BEFORE</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. Default <code>360h</code>. Set this variable if you configure Pinniped and Dex to use self-signed certificates managed by <code>certifcate-manager</code>.</td>
  </tr>
  <tr>
    <td><code>OIDC_IDENTITY_PROVIDER_CLIENT_ID</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The <code>client_id</code> value that you obtain from your OIDC provider. For example, if your provider is Okta, log in to Okta, create a Web application, and select the Client Credentials options in order to get a <code>client_id</code> and <code>secret</code>.</td>
  </tr>
  <tr>
    <td><code>OIDC_IDENTITY_PROVIDER_CLIENT_SECRET</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The Base64 <code>secret</code> value that you obtain from your OIDC provider.</td>
  </tr>
  <tr>
    <td><code>OIDC_IDENTITY_PROVIDER_GROUPS_CLAIM</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The name of your groups claim. This is used to set a user&rsquo;s group in the JSON Web Token (JWT) claim. The default value is <code>groups</code>.</td>
  </tr>
  <tr>
    <td><code>OIDC_IDENTITY_PROVIDER_ISSUER_URL</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The IP or DNS address of your OIDC server.</td>
  </tr>
 <tr>
    <td><code>OIDC_IDENTITY_PROVIDER_SCOPES</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. A comma separated list of additional scopes to request in the token response. For example, <code>&quot;email,offline_access&quot;</code>.</td>
 </tr>
<tr>
  <td><code>OIDC_IDENTITY_PROVIDER_USERNAME_CLAIM</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Required</strong>. The name of your username claim. This is used to set a user&rsquo;s username in the JWT claim. Depending on your provider,  enter claims such as <code>user_name</code>, <code>email</code>, or <code>code</code>.</td>
</tr>
<tr>
  <td><code>SUPERVISOR_ISSUER_URL</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Do not modify</strong>. This variable is automatically updated in the configuration file when you run the <code>tanzu cluster create command</code> command. </td>
</tr>
<tr>
  <td><code>SUPERVISOR_ISSUER_CA_BUNDLE_DATA_B64</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Do not modify</strong>. This variable is automatically updated in the configuration file when you run the <code>tanzu cluster create command</code> command.</td>
</tr>
</tbody>
</table>

### <a id="identity-management-ldap"></a> Identity Providers - LDAP

If you set `IDENTITY_MANAGEMENT_TYPE: ldap`, set the following variables to configure an LDAP identity provider. For more information, see [Enabling Identity Management in Tanzu Kubernetes Grid](mgmt-clusters/enabling-id-mgmt.md) and [Configure Identity Management](mgmt-clusters/create-config-file.md#identity-mgmt) in *Create a Management Cluster Configuration File*.

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
  <tbody>
<tr>
  <td><code>LDAP_BIND_DN</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The DN for an application service account. The connector uses these credentials to search for users and groups. Not required if the LDAP server provides access for anonymous authentication.</td>
</tr>
<tr>
  <td><code>LDAP_BIND_PASSWORD</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The password for an application service account, if  <code>LDAP_BIND_DN</code> is set.</td>
</tr>
<tr>
  <td><code>LDAP_GROUP_SEARCH_BASE_DN</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The point from which to start the LDAP search.<strong> </strong>For example, <code>OU=Groups,OU=domain,DC=io</code>.</td>
</tr>
<tr>
  <td><code>LDAP_GROUP_SEARCH_FILTER</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. An optional filter to be used by the LDAP search</td>
</tr>
<tr>
  <td><code>LDAP_GROUP_SEARCH_GROUP_ATTRIBUTE</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The attribute of the group record that holds the user/member information. For example, <code>member</code>.</td>
</tr>
<tr>
  <td><code>LDAP_GROUP_SEARCH_NAME_ATTRIBUTE</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The LDAP attribute that holds the name of the group. For example, <code>cn</code>.</td>
</tr>
<tr>
  <td><code>LDAP_GROUP_SEARCH_USER_ATTRIBUTE</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The attribute of the user record that is used as the value of the membership attribute of the group record. For example, <code>distinguishedName</code>, <code>dn</code>.</td>
</tr>
<tr>
  <td><code>LDAP_HOST</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Required</strong>. The IP or DNS address of your LDAP server. If the LDAP server is listening on the default port 636, which is the secured configuration, you do not need to specify the port. If the LDAP server is listening on a different port, provide the address and port of the LDAP server, in the form <code>"host:port"</code>.</td>
</tr>
<tr>
  <td><code>LDAP_ROOT_CA_DATA_B64</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. If you are using an LDAPS endpoint, paste the base64 encoded contents of the LDAP server certificate.</td>
</tr>
<tr>
  <td><code>LDAP_USER_SEARCH_BASE_DN</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The point from which to start the LDAP search.<strong> </strong>For example, <code>OU=Users,OU=domain,DC=io</code>. </td>
</tr>
<tr>
  <td><code>LDAP_USER_SEARCH_EMAIL_ATTRIBUTE</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The LDAP attribute that holds the email address. For example, <code>email</code>, <code>userPrincipalName</code>.</td>
</tr>
<tr>
  <td><code>LDAP_USER_SEARCH_FILTER</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>.  An optional filter to be used by the LDAP search.</td>
</tr>
<tr>
  <td><code>LDAP_USER_SEARCH_ID_ATTRIBUTE</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The LDAP attribute that contains the user ID<strong>.</strong> Similar to <code>LDAP_USER_SEARCH_USERNAME</code>.</td>
</tr>
<tr>
  <td><code>LDAP_USER_SEARCH_NAME_ATTRIBUTE</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The LDAP attribute that holds the given name of the user. For example, <code>givenName</code>.</td>
</tr>
<tr>
  <td><code>LDAP_USER_SEARCH_USERNAME</code></td>
  <td>&#10004;</td>
  <td>&#10006;</td>
  <td><strong>Optional</strong>. The  LDAP attribute that contains the user ID. For example, <code>uid</code>, <code>sAMAccountName</code>.</td>
</tr>
</tbody>
</table>

### <a id="size"></a> Node Configuration

Configure the size and number of control plane and worker nodes, and the operating system that the node instances run. For more information, see [Configure Node Settings](mgmt-clusters/create-config-file.md#nodes) in *Create a Management Cluster Configuration File*.

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
  <tbody>
 <tr>
    <td><code>CONTROL_PLANE_MACHINE_COUNT</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Deploy a Tanzu Kubernetes cluster with more control plane nodes than the <code>dev</code> and <code>prod</code> plans define by default. The number of control plane nodes that you specify must be odd. </td>
  </tr>
  <tr>
    <td><code>CONTROLPLANE_SIZE</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Size for control plane node VMs. Overrides the <code>VSPHERE_CONTROL_PLANE_*</code> parameters. See <code>SIZE</code> for possible values.</td>
  </tr>
    <tr>
    <td><code>NODE_STARTUP_TIMEOUT</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>20m</code>.</td>
  </tr>
  <tr>
   <td><code>OS_ARCH</code></td>
   <td>&#10004;</td>
   <td>&#10004;</td>
   <td><strong>Optional</strong>. Architecture for node VM OS. Default and only current choice is <code>amd64</code>.</td>
  </tr>
  <tr>
   <td><code>OS_NAME</code></td>
   <td>&#10004;</td>
   <td>&#10004;</td>
   <td><strong>Optional</strong>. Node VM OS. Defaults to <code>ubuntu</code> for <a href="https://ubuntu.com">Ubuntu LTS</a>. Can also be
   <code>photon</code> for <a href="https://vmware.github.io/photon/assets/files/html/3.0/">Photon OS</a> on vSphere or
   <code>amazon</code> for <a href="https://aws.amazon.com/amazon-linux-2/">Amazon Linux</a> on Amazon EC2.</td>
  </tr>
  <tr>
   <td><code>OS_VERSION</code></td>
   <td>&#10004;</td>
   <td>&#10004;</td>
   <td><strong>Optional</strong>. Version for <code>OS_NAME</code> OS, above. Defaults to <code>20.04</code> for Ubuntu. Can be <code>3</code> for Photon on vSphere and <code>2</code> for Amazon Linux on Amazon EC2.</td>
  </tr>
  <tr>
    <td><code>SIZE</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Size for both control plane and worker node VMs. Overrides the <code>CONTROLPLANE_SIZE</code> and <code>WORKER_SIZE</code> parameters. For vSphere, set <code>small</code>, <code>medium</code>, <code>large</code>, or <code>extra-large</code>. For Amazon EC2, set an instance type, for example, <code>t3.small</code>. For Azure, set an instance type, for example, <code>Standard_D2s_v3</code>.</td>
  </tr>
   <tr>
    <td><code>WORKER_MACHINE_COUNT</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Deploy a Tanzu Kubernetes cluster with more worker nodes than the <code>dev</code> and <code>prod</code> plans define by default.</td>
  </tr>
    <tr>
    <td><code>WORKER_SIZE</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Size for worker node VMs. Overrides the <code>VSPHERE_WORKER_*</code> parameters. See <code>SIZE</code> for possible values.</td>
  </tr>
  </tbody>
</table>

### <a id="autoscaler"></a> Cluster Autoscaler

Additonal variables to set if `ENABLE_AUTOSCALER` is set to `true`. For information about Cluster Autoscaler, [Scale Tanzu Kubernetes Clusters](cluster-lifecycle/scale-cluster.md).

<table class="table">
  <col width="29%">
  <col width="18%">
  <col width="18%">
  <col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
 <tbody>
  <tr>
    <td><code>AUTOSCALER_MAX_NODES_TOTAL</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td>Maximum number of worker nodes in the cluster. Cluster Autoscaler does not attempt to scale your cluster beyond this limit. If set to <code>0</code>, Cluster Autoscaler makes scaling decisions based on the minimum and maximum values that you configure for each machine deployment. Default <code>0</code>. See below.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_SCALE_DOWN_DELAY_AFTER_ADD</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td>Amount of time that Cluster Autoscaler waits after a scale-up operation and then resumes scale-down scans. Default <code>10m</code>.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_SCALE_DOWN_DELAY_AFTER_DELETE</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td>Amount of time that Cluster Autoscaler waits after deleting a node and then resumes scale-down scans. Default <code>10s</code>.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_SCALE_DOWN_DELAY_AFTER_FAILURE</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td>Amount of time that Cluster Autoscaler waits after a scale-down failure and then resumes scale-down scans. Default <code>3m</code>.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_SCALE_DOWN_UNNEEDED_TIME</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td>Amount of time that Cluster Autoscaler must wait before scaling down an eligible node. Default <code>10m</code>.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_MAX_NODE_PROVISION_TIME</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td>Maximum amount of time Cluster Autoscaler waits for a node to be provisioned. Default <code>15m</code>.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_MIN_SIZE_0</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td> <strong>Required</strong>, all IaaSes. Minimum number of worker nodes. Cluster Autoscaler does not attempt to scale down the nodes below this limit. For <code>prod</code> clusters on Amazon EC2, <code>AUTOSCALER_MIN_SIZE_0</code> sets the minimum number of worker nodes in the first AZ. If not set, defaults to the value of <code>WORKER_MACHINE_COUNT</code> for clusters with a single machine deployment or <code>WORKER_MACHINE_COUNT_0</code> for clusters with multiple machine deployments.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_MAX_SIZE_0</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>, all IaaSes. Maximum number of worker nodes. Cluster Autoscaler does not attempt to scale up the nodes beyond this limit. For <code>prod</code> clusters on Amazon EC2, <code>AUTOSCALER_MAX_SIZE_0</code> sets the maximum number of worker nodes in the first AZ. If not set, defaults to the value of <code>WORKER_MACHINE_COUNT</code> for clusters with a single machine deployment or <code>WORKER_MACHINE_COUNT_0</code> for clusters with multiple machine deployments.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_MIN_SIZE_1</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>, use only for <code>prod</code> clusters on Amazon EC2. Minimum number of worker nodes in the second AZ. Cluster Autoscaler does not attempt to scale down the nodes below this limit. If not set, defaults to the value of <code>WORKER_MACHINE_COUNT_1</code>.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_MAX_SIZE_1</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>, use only for <code>prod</code> clusters on Amazon EC2. Maximum number of worker nodes nodes in the second AZ. Cluster Autoscaler does not attempt to scale up the nodes beyond this limit. If not set, defaults to the value of <code>WORKER_MACHINE_COUNT_1</code>.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_MIN_SIZE_2</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>, use only for <code>prod</code> clusters on Amazon EC2. Minimum number of worker nodes in the third AZ. Cluster Autoscaler does not attempt to scale down the nodes below this limit. If not set, defaults to the value of <code>WORKER_MACHINE_COUNT_2</code>.</td>
  </tr>
  <tr>
    <td><code>AUTOSCALER_MAX_SIZE_2</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>, use only for <code>prod</code> clusters on Amazon EC2. Maximum number of worker nodes in the third AZ. Cluster Autoscaler does not attempt to scale up the nodes beyond this limit. If not set, defaults to the value of <code>WORKER_MACHINE_COUNT_2</code>.</td>
  </tr>
 </tbody>
</table>

### <a id="proxies"></a> Proxy Configuration

If your environment is internet-restricted or otherwise includes proxies, you can optionally configure Tanzu Kubernetes Grid to send outgoing HTTP and HTTPS traffic from `kubelet`, `containerd`, and the control plane to your proxies.

Tanzu Kubernetes Grid allows you to enable proxies for any of the following:

* For both the management cluster and one or more Tanzu Kubernetes clusters
* For the management cluster only
* For one or more Tanzu Kubernetes clusters

For more information, see [Configure Proxies](mgmt-clusters/create-config-file.md#proxies) in *Create a Management Cluster Configuration File*.

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
  <thead>
    <tr>
      <th rowspan=2>Variable</th>
      <th colspan=2>Can be set in...</th>
      <th rowspan=2>Description</th>
    </tr>
    <tr>
      <th>Management cluster YAML</th>
      <th>Tanzu Kubernetes cluster YAML</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>TKG_HTTP_PROXY</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><p><strong>Optional, set if you want to configure a proxy; to disable your proxy configuration for an individual cluster, set this to <code>""</code></strong>. The URL of your HTTP proxy, formatted as follows: <code>PROTOCOL://USERNAME:PASSWORD@FQDN-OR-IP:PORT</code>, where:</p>
        <ul>
          <li><strong>Required</strong>. <code>PROTOCOL</code> is <code>http</code>.</li>
          <li><strong>Optional</strong>. <code>USERNAME</code> and <code>PASSWORD</code> are your HTTP proxy username and password. Include these if the proxy requires authentication.</li>
          <li><strong>Required</strong>. <code>FQDN-OR-IP</code> and <code>PORT</code> are the FQDN or IP address and port number of your HTTP proxy.</li>
        </ul>
        <p>For example, <code>http://user:password@myproxy.com:1234</code> or <code>http://myproxy.com:1234</code>. If you set <code>TKG_HTTP_PROXY</code>, you must also set <code>TKG_HTTPS_PROXY</code>.</p></td>
    </tr>
    <tr>
      <td><code>TKG_HTTPS_PROXY</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to configure a proxy</strong>. The URL of your HTTPS proxy. You can set this variable to the same value as <code>TKG_HTTP_PROXY</code> or provide a different value. The URL must start with <code>http://</code>. If you set <code>TKG_HTTPS_PROXY</code>, you must also set <code>TKG_HTTP_PROXY</code>.</td>
    </tr>
    <tr>
      <td><code>TKG_NO_PROXY</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><p><strong>Optional</strong>. One or more network CIDRs or hostnames, comma-separated and without spaces, that must bypass the HTTP(S) proxy.
      For example, <code>noproxy.yourdomain.com,192.168.0.0/24</code>.</p>
      <p>Internally, Tanzu Kubernetes Grid appends <code>localhost</code>, <code>127.0.0.1</code>, the values of <code>CLUSTER_CIDR</code> and <code>SERVICE_CIDR</code>, <code>.svc</code>, and <code>.svc.cluster.local</code> to the value that you set in <code>TKG_NO_PROXY</code>. It also appends your AWS VPC CIDR and <code>169.254.0.0/16</code> for deployments to Amazon EC2 and your Azure VNET CIDR, <code>169.254.0.0/16</code>, and <code>168.63.129.16</code> for deployments to Azure. For vSphere, you must manually add the CIDR of <code>VSPHERE_NETWORK</code>, which includes the IP address of your control plane endpoint, to <code>TKG_NO_PROXY</code>. If you set <code>VSPHERE_CONTROL_PLANE_ENDPOINT</code> to an FQDN, add both the FQDN and <code>VSPHERE_NETWORK</code> to <code>TKG_NO_PROXY</code>.</p>
      <p><strong>Important:</strong> In environments where Tanzu Kubernetes Grid runs behind a proxy, <code>TKG_NO_PROXY</code> lets cluster VMs communicate directly with infrastructure that runs the same network, behind the same proxy. This may include, but is not limited to, your infrastructure, OIDC or LDAP server, Harbor, NSX-T and NSX Advanced Load Balancer (vSphere), and AWS VPC CIDRs (AWS). Set <code>TKG_NO_PROXY</code> to include all such endpoints that clusters must access but that are not reachable by your proxies.</p></td>
    </tr>
  </tbody>
</table>

### <a id="antrea"></a> Antrea CNI Configuration

Additonal optional variables to set if `CNI` is set to `antrea`. For more information, see [Configure Antrea CNI](mgmt-clusters/create-config-file.md#antrea) in *Create a Management Cluster Configuration File*.

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
  <tbody>
<tr>
    <td><code>ANTREA_NO_SNAT</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Default <code>false</code>. Set to <code>true</code> to disable Source Network Address Translation (SNAT).</td>
  </tr>
  <tr>
    <td><code>ANTREA_TRAFFIC_ENCAP_MODE</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Default <code>"encap"</code>. Set to either <code>noEncap</code>, <code>hybrid</code>, or <code>NetworkPolicyOnly</code>. For information about using NoEncap or Hybrid traffic modes, see <a href="https://antrea.io/docs/noencap-hybrid-modes/">NoEncap and Hybrid Traffic Modes of Antrea</a> in the Antrea documentation.</td>
  </tr>
  <tr>
    <td><code>ANTREA_PROXY</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Default <code>false</code>.  Enables or disables <code>AntreaProxy</code>, to replace <code>kube-proxy</code> for pod-to-ClusterIP Service traffic, for better performance and lower latency. Note that <code>kube-proxy</code> is still used for other types of Service traffic. </td>
  </tr>
  <tr>
    <td><code>ANTREA_POLICY</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Default <code>true</code>. Enables or disables the Antrea-native policy API, which are policy CRDs specific to Antrea. Also, the implementation of Kubernetes Network Policies remains active when this variable is enabled. For information about using network policies, see <a href="https://antrea.io/docs/v0.11.3/antrea-network-policy/">Antrea Network Policy CRDs</a> in the Antrea documentation.</td>
  </tr>
  <tr>
    <td><code>ANTREA_TRACEFLOW</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Default <code>false</code>. Set to <code>true</code> to enable Traceflow. For information about using Traceflow, see the <a href="https://antrea.io/docs/v0.13.1/traceflow-guide/">Traceflow User Guide</a> in the Antrea documentation.</td>
  </tr>
</tbody>
</table>

### <a id="mhc"></a> Machine Health Checks

If you want to configure machine health checks for management and Tanzu Kubernetes clusters, set the following variables. For more information, see [Configure Machine Health Checks](mgmt-clusters/create-config-file.md#mhc) in *Create a Management Cluster Configuration File*. For information about how to perform Machine Health Check operations after cluster deployment, see [Configure Machine Health Checks for Tanzu Kubernetes Clusters](cluster-lifecycle/configure-health-checks.md).

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
  <thead>
    <tr>
      <th rowspan=2>Variable</th>
      <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
    </tr>
    <tr>
      <th>Management cluster YAML</th>
      <th>Tanzu Kubernetes cluster YAML</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>ENABLE_MHC</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>true</code>. This variable enables or disables the <a href="https://cluster-api.sigs.k8s.io/developer/architecture/controllers/machine-health-check.html#machinehealthcheck"><code>MachineHealthCheck</code></a> controller, which provides node health monitoring and node auto-repair for worker nodes in management and Tanzu Kubernetes clusters. You can also enable or disable <code>MachineHealthCheck</code> after deployment by using the CLI. For instructions, see <a href="cluster-lifecycle/configure-health-checks.md">Configure Machine Health Checks for Tanzu Kubernetes Clusters</a>.</td>
    </tr>
    <tr>
      <td><code>MHC_UNKNOWN_STATUS_TIMEOUT</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>5m</code>. By default, if the <code>Ready</code> condition of a node remains <code>Unknown</code> for longer than <code>5m</code>, <code>MachineHealthCheck</code> considers the machine unhealthy and recreates it.</td>
    </tr>
    <tr>
      <td><code>MHC_FALSE_STATUS_TIMEOUT</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>5m</code>. By default, if the <code>Ready</code> condition of a node remains <code>False</code> for longer than <code>5m</code>, <code>MachineHealthCheck</code> considers the machine unhealthy and recreates it.</td>
    </tr>
  </tbody>
</table>

### <a id="mhc"></a> Private Image Repository Configuration

If you deploy Tanzu Kubernetes Grid management clusters and Kubernetes clusters in environments that are not connected to the Internet, you need to set up a private image repository within your firewall and populate it with the Tanzu Kubernetes Grid images. For information about setting up a private image repository, see [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](mgmt-clusters/airgapped-environments.md) and [Deploy Harbor Registry as a Shared Service](extensions/harbor-registry.md).

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
  <thead>
    <tr>
      <th rowspan=2>Variable</th>
      <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
    </tr>
    <tr>
      <th>Management cluster YAML</th>
      <th>Tanzu Kubernetes cluster YAML</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>TKG_CUSTOM_IMAGE_REPOSITORY</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong> if you deploy Tanzu Kubernetes Grid in an Internet-restricted environment. Provide the IP address or FQDN of your private registry. For example, <code>custom-image-repository.io/yourproject</code>.</td>
    </tr>
    <tr>
      <td><code>TKG_CUSTOM_IMAGE_REPOSITORY_SKIP_TLS_VERIFY</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td>Not currently implemented. Do not use. </td>
    </tr>
    <tr>
      <td><code>TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Set if your private image registry uses a self-signed certificate. Provide the CA certificate in base64 encoded format, for example <code>TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE: "LS0t[...]tLS0tLQ==""</code>.</td>
    </tr>
  </tbody>
</table>

## <a id="vsphere"></a> vSphere

The options in the table below are the minimum options that you specify in the cluster configuration file when deploying Tanzu Kubernetes clusters to vSphere. Most of these options are the same for both the Tanzu Kubernetes cluster and the management cluster that you use to deploy it.

For more information about the configuration files for vSphere, see [Management Cluster Configuration for vSphere](mgmt-clusters/config-vsphere.md) and [Deploy Tanzu Kubernetes Clusters to vSphere](tanzu-k8s-clusters/vsphere.md).

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
 <tbody>
<tr>
    <td><code>DEPLOY_TKG_ON_VSPHERE7</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. If deploying to vSphere 7, set to <code>true</code> to skip the prompt about deployment on vSphere 7, or <code>false</code>. See <a href="mgmt-clusters/config-vsphere.md#vsphere-7">Management Clusters on vSphere with Tanzu</a>.</td>
  </tr>
  <tr>
    <td><code>ENABLE_TKGS_ON_VSPHERE7</code></td>
        <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong> if deploying to vSphere 7. Set to <code>true</code> to be redirected to the vSphere with Tanzu enablement UI page, or <code>false</code>. See <a href="mgmt-clusters/config-vsphere.md#vsphere-7">Management Clusters on vSphere with Tanzu</a>.</td>
  </tr>
<tr>
     <td><code>VIP_NETWORK_INTERFACE</code></td>
        <td>&#10004;</td>
    <td>&#10004;</td>
     <td><strong>Optional</strong>. Set to <code>eth0</code>, <code>eth1</code>, etc. Network interface name, for example an Ethernet interface. Defaults to <code>eth0</code>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_CONTROL_PLANE_DISK_GIB</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The size in gigabytes of the disk for the control plane node VMs. Include the quotes (<code>""</code>). For example, <code>&quot;30&quot;</code>.</td>
  </tr>
<tr>
    <td><code>VSPHERE_CONTROL_PLANE_ENDPOINT</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. Static virtual IP address for API requests to the Tanzu Kubernetes cluster. If you mapped a fully qualified domain name (FQDN) to the VIP address, you can specify the FQDN instead of the VIP address.</td>
  </tr>
    <tr>
    <td><code>VSPHERE_CONTROL_PLANE_MEM_MIB</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The amount of memory in megabytes for the control plane node VMs. Include the quotes (<code>""</code>). For example, <code>&quot;2048&quot;</code>.</td>
   </tr>
  <tr>
    <td><code>VSPHERE_CONTROL_PLANE_NUM_CPUS</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The number of CPUs for the control plane node VMs. Include the quotes (<code>""</code>). Must be at least 2. For example, <code>&quot;2&quot;</code>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_DATACENTER</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The name of the datacenter in which to deploy the cluster, as it appears in the vSphere inventory. For example, <code>/MY-DATACENTER</code>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_DATASTORE</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The name of the vSphere datastore for the cluster to use, as it appears in the vSphere inventory. For example, <code>/MY-DATACENTER/datastore/MyDatastore</code>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_FOLDER</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The name of an existing VM folder in which to place Tanzu Kubernetes Grid VMs, as it appears in the vSphere inventory. For example, if you created a folder named <code>TKG</code>, the path is <code>/MY-DATACENTER/vm/TKG</code>.</td>
  </tr>
<tr>
    <td><code>VSPHERE_INSECURE</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Set to <code>true</code> or <code>false</code> to bypass thumbprint verification. If false, set <code>VSPHERE_TLS_THUMBPRINT</code>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_NETWORK</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The name of an existing vSphere network to use as the Kubernetes service network, as it appears in the vSphere inventory. For example, <code>VM Network</code>.</td>
  </tr>
<tr>
    <td><code>VSPHERE_PASSWORD</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The password for the vSphere user account. This value is base64-encoded when you run <code>tanzu cluster create</code>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_RESOURCE_POOL</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The name of an existing resource pool in which to place this Tanzu Kubernetes Grid instance, as it appears in the vSphere inventory. To use the root resource pool for a cluster, enter the full path, for example for a cluster named <code>cluster0</code> in datacenter <code>MY-DATACENTER</code>, the full path is <code>/MY-DATACENTER/host/cluster0/Resources</code>.</td>
  </tr>
<tr>
    <td><code>VSPHERE_SERVER</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The IP address or FQDN of the vCenter Server instance on which to deploy the Tanzu Kubernetes cluster.</td>
  </tr>
<tr>
    <td><code>VSPHERE_SSH_AUTHORIZED_KEY</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. Paste in the contents of the SSH public key that you created in <a href="mgmt-clusters/vsphere.md">Deploy a Management Cluster to vSphere</a>. For example, <code>&quot;ssh-rsa
      NzaC1yc2EA [...] hnng2OYYSl+8ZyNz3fmRGX8uPYqw==
      email@example.com&quot;.</code></td>
  </tr>
  <tr>
    <td><code>VSPHERE_STORAGE_POLICY_ID</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The name of a VM storage policy for the management cluster, as it appears in <strong>Policies and Profiles</strong> > <strong>VM Storage Policies</strong>.<br />
    If <code>VSPHERE_DATASTORE</code> is set, the storage policy must include it. Otherwise, the cluster creation process chooses a datastore that compatible with the policy.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_TEMPLATE</code></td>
    <td>&#10006;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Specify the path to an OVA file if you are using multiple custom OVA images for the same Kubernetes version, in the format <code>/MY-DC/vm/MY-FOLDER-PATH/MY-IMAGE</code>. For more information, see <a href="tanzu-k8s-clusters/vsphere.md#custom-ova">Deploy a Cluster with a Custom OVA Image</a>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_TLS_THUMBPRINT</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong> if <code>VSPHERE_INSECURE</code> is <code>false</code>. The thumbprint of the vCenter Server certificate. For information about how to obtain the vCenter Server certificate thumbprint, see <a href="mgmt-clusters/vsphere.md#vc-thumbprint">Obtain vSphere Certificate Thumbprints</a>. This value can be skipped if user wants to use insecure connection by setting `VSPHERE_INSECURE` to `true`.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_USERNAME</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. A vSphere user account with the required privileges for Tanzu Kubernetes Grid operation. For example, <code>tkg-user@vsphere.local</code>.</td>
  </tr>
   <tr>
    <td><code>VSPHERE_WORKER_DISK_GIB</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The size in gigabytes of the disk for the worker node VMs. Include the quotes (<code>""</code>). For example, <code>&quot;50&quot;</code>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_WORKER_MEM_MIB</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The amount of memory in megabytes for the worker node VMs. Include the quotes (<code>""</code>). For example, <code>&quot;4096&quot;</code>.</td>
  </tr>
  <tr>
    <td><code>VSPHERE_WORKER_NUM_CPUS</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The number of CPUs for the worker node VMs. Include the quotes (<code>""</code>). Must be at least 2. For example, <code>&quot;2&quot;</code>.</td>
  </tr>
  </tbody>
</table>

### <a id="nsx-lb"></a> NSX Advanced Load Balancer

For information about how to deploy NSX Advanced Load Balancer, see [Install VMware NSX Advanced Load Balancer on a vSphere Distributed Switch](mgmt-clusters/install-nsx-adv-lb.md).

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
  <tbody>
     <tr>
    <td><code>AVI_ENABLE</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. Set to <code>true</code> or <code>false</code>. Enables NSX Advanced Load Balancer. If <code>true</code>, you must set the required variables listed in <a href="#nsx-lb">NSX Advanced Load Balancer</a> below. Defaults to <code>false</code>.</td>
  </tr>
 <tr>
    <td><code>AVI_ADMIN_CREDENTIAL_NAME</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. The name of the Kubernetes Secret that contains the NSX Advanced Loader Balancer controller admin username and password. Default <code>avi-controller-credentials</code>.</td>
  </tr>
<tr>
    <td><code>AVI_AKO_IMAGE_PULL_POLICY</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. Default <code>IfNotPresent</code>.</td>
  </tr>
    <tr>
    <td><code>AVI_CA_DATA_B64</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The contents of the Controller Certificate Authority that is used to sign the Controller certificate. It must be base64 encoded.</td>
  </tr>
    <tr>
    <td><code>AVI_CA_NAME</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. The name of the Kubernetes Secret that holds the NSX Advanced Loader Balancer Controller Certificate Authority. Default <code>avi-controller-ca</code>.</td>
  </tr>
    <tr>
    <td><code>AVI_CLOUD_NAME</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The cloud that you created in your NSX Advanced Load Balancer deployment. For example, <code>Default-Cloud</code>.</td>
  </tr>
    <tr>
    <td><code>AVI_CONTROLLER</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The IP or hostname of the NSX Advanced Loader Balancer controller.</td>
  </tr>
  <tr>
    <td><code>AVI_DATA_NETWORK</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The name of the Network on which the Load Balancer floating IP subnet or IP Pool is configured. This Network must be present in the same vCenter Server instance as the Kubernetes network that Tanzu Kubernetes Grid uses, that you specify in the <code>SERVICE_CIDR</code> variable. This allows NSX Advanced Load Balancer to discover the Kubernetes network in vCenter Server and to deploy and configure Service Engines.</td>
  </tr>
    <tr>
    <td><code>AVI_DATA_NETWORK_CIDR</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>.The CIDR of the subnet to use for the load balancer VIP. This comes from one of the VIP network's configured subnets. You can see the subnet CIDR for a particular Network in the <b>Infrastructure</b> - <b>Networks</b> view of the NSX Advanced Load Balancer interface.</td>
  </tr>
  <tr>
    <td><code>AVI_DISABLE_INGRESS_CLASS</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. Disable Ingress Class. Default <code>false</code>.</td>
  </tr>
    <tr>
    <td><code>AVI_INGRESS_DEFAULT_INGRESS_CONTROLLER</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. Use AKO as the default Ingress Controller. Default <code>false</code>.</td>
  </tr>
  <tr>
    <td><code>AVI_INGRESS_SERVICE_TYPE</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. Specifies whether the AKO functions in <code>ClusterIP</code> mode or <code>NodePort</code> mode. Defaults to <code>NodePort</code>.
</td>
  </tr>
  <tr>
    <td><code>AVI_INGRESS_SHARD_VS_SIZE</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. AKO uses a sharding logic for Layer 7 ingress objects.
 A sharded VS involves hosting multiple insecure or secure ingresses hosted by one virtual IP or VIP. Set to <code>LARGE</code>, <code>MEDIUM</code>, or <code>SMALL</code>. Default <code>SMALL</code>. Use this to control the layer 7 VS numbers. This applies to both secure/insecure VSes but does not apply for passthrough.</td>
  </tr>
    <tr>
    <td><code>AVI_LABELS</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. Optional labels in the format <code>key: value</code>. When set, NSX Advanced Load Balancer is enabled only on workload clusters that have this label. For example <code>team: tkg</code>.</td>
  </tr>
  <tr>
    <td><code>AVI_NAMESPACE</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Optional</strong>. The namespace for AKO operator. Default <code>"tkg-system-networking"</code>.</td>
  </tr>
  <tr>
    <td><code>AVI_PASSWORD</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The password that you set for the Controller admin when you deployed it.</td>
  </tr>
    <tr>
    <td><code>AVI_SERVICE_ENGINE_GROUP</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. Name of the Service Engine Group. For example, <code>Default-Group</code>. </td>
  </tr>
  <tr>
    <td><code>AVI_USERNAME</code></td>
    <td>&#10004;</td>
    <td>&#10006;</td>
    <td><strong>Required</strong>. The admin username that you set for the Controller host when you deployed it.</td>
  </tr>
 </tbody>
</table>

### <a id="nsxt-pod-routing"></a> NSX-T Pod Routing

These variables configure routable-IP address workload pods, as described in [Deploy a Cluster with Routable-IP Pods](tanzu-k8s-clusters/networking.md#routable-config).
All variables are strings in double-quotes, for example `"true"`.

<table class="table">
  <col width="29%">
  <col width="18%">
  <col width="18%">
  <col width="35%">
  <thead>
    <tr>
      <th rowspan=2>Variable</th>
      <th colspan=2>Can be set in...</th>
      <th rowspan=2>Description</th>
    </tr>
    <tr>
      <th>Management cluster YAML</th>
      <th>Tanzu Kubernetes cluster YAML</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>NSXT_POD_ROUTING_ENABLED</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. <code>"true"</code> enables NSX-T routable pods with the variables below.
      Default is <code>"false"</code>. See <a href="tanzu-k8s-clusters/networking.md#routable-ip">Deploy a Cluster with Routable-IP Pods</a>.</td>
    </tr>
    <tr>
      <td><code>NSXT_MANAGER_HOST</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong> if <code>NSXT_POD_ROUTING_ENABLED= "true"</code>.<br />
      IP address of NSX-T Manager.</td>
    </tr>
    <tr>
      <td><code>NSXT_ROUTER_PATH</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong> if <code>NSXT_POD_ROUTING_ENABLED= "true"</code>. T1 router path shown in NSX-T Manager.</td>
    </tr>
    <tr>
      <th colspan=4>For username/password authentication to NSX-T:</th>
    </tr>
    <tr>
      <td><code>NSXT_USERNAME</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>Username for logging in to NSX-T Manager.</td>
    </tr>
    <tr>
      <td><code>NSXT_PASSWORD</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>Password for logging in to NSX-T Manager.</td>
    </tr>
    <tr>
      <th colspan=4>For authenticating to NSX-T using credentials and storing them in a Kubernetes secret (also set <code>NSXT_USERNAME</code> and <code>NSXT_PASSWORD</code> above):</th>
    </tr>
    <tr>
      <td><code>NSXT_SECRET_NAMESPACE</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>The namespace with the secret containing NSX-T username and password. Default is <code>"kube-system"</code>.</td>
    </tr>
    <tr>
      <td><code>NSXT_SECRET_NAME</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>The name of the secret containing NSX-T username and password. Default is <code>"cloud-provider-vsphere-nsxt-credentials"</code>.</td>
    </tr>
    <tr>
      <th colspan=4>For certificate authentication to NSX-T:</th>
    </tr>
    <tr>
      <td><code>NSXT_ALLOW_UNVERIFIED_SSL</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>Set this to <code>"true"</code> if NSX-T uses a self-signed certificate. Default is <code>false</code>.</td>
    </tr>
    <tr>
      <td><code>NSXT_ROOT_CA_DATA_B64</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong> if <code>NSXT_ALLOW_UNVERIFIED_SSL= "false"</code>.<br />
      Base64-encoded Certificate Authority root certificate string that NSX-T uses for LDAP authentication.</td>
    </tr>
    <tr>
      <td><code>NSXT_CLIENT_CERT_KEY_DATA</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>Base64-encoded cert key file string for local client certificate.</td>
    </tr>
    <tr>
      <td><code>NSXT_CLIENT_CERT_DATA</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>Base64-encoded cert file string for local client certificate.</td>
    </tr>
    <tr>
      <th colspan=4>For remote authentication to NSX-T with VMware Identity Manager, on VMware Cloud (VMC):</th>
    </tr>
    <tr>
      <td><code>NSXT_REMOTE_AUTH</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>Set this to <code>"true"</code> for remote authentication to NSX-T with VMware Identity Manager, on VMware Cloud (VMC). Default is <code>"false"</code>.</td>
    </tr>
    <tr>
      <td><code>NSXT_VMC_AUTH_HOST</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>VMC authentication host. Default is empty.</td>
    </tr>
    <tr>
      <td><code>NSXT_VMC_ACCESS_TOKEN</code></td>
      <td>&#10006;</td>
      <td>&#10004;</td>
      <td>VMC authentication access token. Default is empty.</td>
    </tr>
  </tbody>
</table>

## <a id="aws"></a> Amazon EC2

The variables in the table below are the options that you specify in the cluster configuration file when deploying Tanzu Kubernetes clusters to Amazon EC2. Many of these options are the same for both the Tanzu Kubernetes cluster and the management cluster that you use to deploy it.

For more information about the configuration files for Amazon EC2, see [Management Cluster Configuration for Amazon EC2](mgmt-clusters/config-aws.md) and [Deploy Tanzu Kubernetes Clusters to Amazon EC2](tanzu-k8s-clusters/aws.md).

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
 <thead>
  <tr>
    <th rowspan=2>Variable</th>
    <th colspan=2>Can be set in...</th>
    <th rowspan=2>Description</th>
  </tr>
  <tr>
    <th>Management cluster YAML</th>
    <th>Tanzu Kubernetes cluster YAML</th>
  </tr>
 </thead>
  <tbody>
  <tr>
    <td><code>AWS_ACCESS_KEY_ID</code></td>
      <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The access key ID for your AWS account. Alternatively, you can specify account credentials as a local environment variables or in your AWS default credential provider chain.</td>
  </tr>

  <tr>
    <td><code>AWS_NODE_AZ</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The name of the AWS availability zone in your chosen region that you want use as the availability zone for this management cluster. Availability zone names are the same as the AWS region name, with a single lower-case letter suffix, such as <code>a</code>, <code>b</code>, <code>c</code>. For example,
      <code>us-west-2a</code>. To deploy a <code>prod</code> management cluster with three control plane nodes, you must also set
      <code>AWS_NODE_AZ_1</code> and <code>AWS_NODE_AZ_2</code>. The letter suffix in each of these availability zones must be unique. For example, <code>us-west-2a</code>,
      <code>us-west-2b</code>, and <code>us-west-2c</code>.</td>
  </tr>
  <tr>
    <td><code>AWS_NODE_AZ_1</code> and <code>AWS_NODE_AZ_2</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Set these variables if you want to deploy a <code>prod</code> management cluster with three control plane nodes. Both availability zones must be in the same region as <code>AWS_NODE_AZ</code>. See <code>AWS_NODE_AZ</code> above for more information.
   For example, <code>us-west-2a</code>, <code>ap-northeast-2b</code>, etc.</td>
  <tr>
    <td><code>AWS_PRIVATE_NODE_CIDR</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Set this variable if you set <code>AWS_VPC_CIDR</code>. If the recommended range of 10.0.0.0/24 is not available, enter a different IP range in CIDR format for private nodes to use. When Tanzu Kubernetes Grid deploys your management cluster, it creates this subnetwork in
      <code>AWS_NODE_AZ</code>. To deploy a <code>prod</code> management cluster with three control plane nodes, you must also set
      <code>AWS_PRIVATE_NODE_CIDR_1</code> and
      <code>AWS_PRIVATE_NODE_CIDR_2</code>. For example, <code>10.0.0.0/24<code>.</td></td>
  </tr>
  <tr>
    <td><code>AWS_PRIVATE_NODE_CIDR_1</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. If the recommended range of 10.0.2.0/24 is not available, enter a different IP range in CIDR format. When Tanzu Kubernetes Grid deploys your management cluster, it creates this subnetwork in <code>AWS_NODE_AZ_1</code>. See
    <code>AWS_PRIVATE_NODE_CIDR</code> above for more information.</td>
  <tr>
    <td><code>AWS_PRIVATE_NODE_CIDR_2</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. If the recommended range of 10.0.4.0/24 is not available, enter a different IP range in CIDR format. When Tanzu Kubernetes Grid deploys your management cluster, it creates this subnetwork in <code>AWS_NODE_AZ_2</code>. See
      <code>AWS_PRIVATE_NODE_CIDR</code> above for more information.</td>
  <tr>
  <tr>
    <td><code>AWS_PUBLIC_NODE_CIDR</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Set this variable if you set <code>AWS_VPC_CIDR</code>. If the recommended range of 10.0.1.0/24 is not available, enter a different IP range in CIDR format for public nodes to use. When Tanzu Kubernetes Grid deploys your management cluster, it creates this subnetwork in
      <code>AWS_NODE_AZ</code>. To deploy a <code>prod</code> management cluster with three control plane nodes, you must also set
      <code>AWS_PUBLIC_NODE_CIDR_1</code> and
      <code>AWS_PUBLIC_NODE_CIDR_2</code>.</td>
  </tr>
  <tr>
    <td><code>AWS_PUBLIC_NODE_CIDR_1</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. If the recommended range of 10.0.3.0/24 is not available, enter a different IP range in CIDR format. When Tanzu Kubernetes Grid deploys your management cluster, it creates this subnetwork in <code>AWS_NODE_AZ_1</code>.
    See <code>AWS_PUBLIC_NODE_CIDR</code> above for more information.</td>
  </tr>
  <tr>
    <td><code>AWS_PUBLIC_NODE_CIDR_2</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. If the recommended range of 10.0.5.0/24 is not available, enter a different IP range in CIDR format. When Tanzu Kubernetes Grid deploys your management cluster, it creates this subnetwork in <code>AWS_NODE_AZ_2</code>. See <code>AWS_PUBLIC_NODE_CIDR</code> above for more information.</td>
  </tr>
  <tr>
    <td><code>AWS_PRIVATE_SUBNET_ID</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. If you set <code>AWS_VPC_ID</code> to use an existing VPC, enter the ID of a private subnet that already exists in <code>AWS_NODE_AZ</code>. This setting is optional. If you do not set it, <code>tanzu management-cluster create</code> identifies the private subnet automatically. To deploy a <code>prod</code> management cluster with three control plane nodes, you must also set
      <code>AWS_PRIVATE_SUBNET_ID_1</code> and
      <code>AWS_PRIVATE_SUBNET_ID_2</code>.</td>
  </tr>
  <tr>
    <td><code>AWS_PRIVATE_SUBNET_ID_1</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The ID of a private subnet that exists in <code>AWS_NODE_AZ_1</code>. If you do not set this variable, <code>tanzu management-cluster create</code> identifies the private subnet automatically. See
      <code>AWS_PRIVATE_SUBNET_ID</code> above for more information.</td>
  </tr>
  <tr>
    <td><code>AWS_PRIVATE_SUBNET_ID_2</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The ID of a private subnet that exists in <code>AWS_NODE_AZ_2</code>. If you do not set this variable, <code>tanzu management-cluster create</code> identifies the private subnet automatically. See
      <code>AWS_PRIVATE_SUBNET_ID</code> above for more information.</td>
  </tr>
  <tr>
    <td><code>AWS_PUBLIC_SUBNET_ID</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. If you set <code>AWS_VPC_ID</code> to use an existing VPC, enter the ID of a public subnet that already exists in <code>AWS_NODE_AZ</code>. This setting is optional. If you do not set it, <code>tanzu management-cluster create</code> identifies the public subnet automatically. To deploy a <code>prod</code> management cluster with three control plane nodes, you must also set
      <code>AWS_PUBLIC_SUBNET_ID_1</code> and
      <code>AWS_PUBLIC_SUBNET_ID_2</code>.</td>
  </tr>
  <tr>
    <td><code>AWS_PUBLIC_SUBNET_ID_1</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The ID of a public subnet that exists in <code>AWS_NODE_AZ_1</code>. If you do not set this variable, <code>tanzu management-cluster create</code> identifies the public subnet automatically. See
      <code>AWS_PUBLIC_SUBNET_ID</code> above for more information.</td>
  </tr>
  <tr>
    <td><code>AWS_PUBLIC_SUBNET_ID_2</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. The ID of a public subnet that exists in <code>AWS_NODE_AZ_2</code>. If you do not set this variable, <code>tanzu management-cluster create</code> identifies the public subnet automatically. See
      <code>AWS_PUBLIC_SUBNET_ID</code> above for more information.</td>
  </tr>
  <tr>
    <td><code>AWS_REGION</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The name of the AWS region in which to deploy the cluster. For example, <code>us-west-2</code>. You can also specify the
      <code>us-gov-east</code> and <code>us-gov-west</code> regions in AWS GovCloud. If you have already set a different region as an environment variable, for example, in <a href="mgmt-clusters/aws.md">Deploy Management Clusters to Amazon EC2</a>, you must unset that environment variable. For example, <code>us-west-2</code>, <code>ap-northeast-2</code>, etc.</td>
  </tr>

  <tr>
    <td><code>AWS_SECRET_ACCESS_KEY</code></td>
      <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The secret access key for your AWS account. Alternatively, you can specify account credentials as an environment variable with the same name or in your AWS default credential provider chain. </td>
  </tr>
  <tr>
    <td><code>AWS_SESSION_TOKEN</code></td>
      <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. Provide the AWS session token granted to your account if you are required to use a temporary access key. For more information about using temporary access keys, see [Understanding and getting your AWS credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#temporary-access-keys). provide the session token for your AWS account. Alternatively, you can specify account credentials as a local environment variables or in your AWS default credential provider chain. </td>
  </tr>
    <tr>
    <td><code>AWS_SSH_KEY_NAME</code></td>
     <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong>. The name of the SSH private key that you registered with your AWS account.</td>
  </tr>
  <tr>
    <td><code>AWS_VPC_ID</code></td>
      <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. To use a VPC that already exists in your selected AWS region, enter the ID of the VPC and then set <code>AWS_PUBLIC_SUBNET_ID</code> and
      <code>AWS_PRIVATE_SUBNET_ID</code>. Set either <code>AWS_VPC_ID</code> or <code>AWS_VPC_CIDR</code>, but not both.</td>
  </tr>
  <tr>
    <td><code>AWS_VPC_CIDR</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. <code>10.0.0.0/16</code>. If you want Tanzu Kubernetes Grid to create a new VPC in the selected region, set the <code>AWS_VPC_CIDR</code>,
      <code>AWS_PUBLIC_NODE_CIDR</code>, and
      <code>AWS_PRIVATE_NODE_CIDR</code> variables. If the recommended range of 10.0.0.0/16 is not available, enter a different IP range in CIDR format in <code>AWS_VPC_CIDR</code> for the management cluster to use. Set either <code>AWS_VPC_CIDR</code> or <code>AWS_VPC_ID</code>, but not both.</td>
  </tr>
  <tr>
    <td><code>BASTION_HOST_ENABLED</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Optional</strong>. By default this option is set to <code>&quot;true&quot;</code> in the global Tanzu Kubernetes Grid configuration. Specify <code>&quot;true&quot;</code> to deploy an AWS bastion host or <code>&quot;false&quot;</code> to reuse an existing bastion host. If no bastion host exists in your availability zone(s) and you set <code>AWS_VPC_ID</code> to use an existing VPC, set <code>BASTION_HOST_ENABLED</code> to <code>&quot;true&quot;</code>.</td>
  </tr>
  <tr>
    <td><code>CONTROL_PLANE_MACHINE_TYPE</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong> if cloud-agnostic <code>SIZE</code> or <code>CONTROLPLANE_SIZE</code> are not set. The Amazon EC2 instance type to use for cluster control plane nodes, for example <code>t3.small</code> or <code>m5.large</code>.</td>
  </tr>
  <tr>
    <td><code>NODE_MACHINE_TYPE</code></td>
    <td>&#10004;</td>
    <td>&#10004;</td>
    <td><strong>Required</strong> if cloud-agnostic <code>SIZE</code> or <code>WORKER_SIZE</code> are not set. The Amazon EC2, instance type to use for cluster worker nodes, for example <code>t3.small</code> or <code>m5.large</code>.</td>
  </tr>
  </tbody>
</table>

## <a id="azure"></a> Microsoft Azure

The variables in the table below are the options that you specify in the cluster configuration file when deploying Tanzu Kubernetes clusters to Azure. Many of these options are the same for both the Tanzu Kubernetes cluster and the management cluster that you use to deploy it.

For more information about the configuration files for Azure, see [Management Cluster Configuration for Azure](mgmt-clusters/config-azure.md) and [Deploy Tanzu Kubernetes Clusters to Azure](tanzu-k8s-clusters/azure.md).

<table class="table">
<col width="29%">
<col width="18%">
<col width="18%">
<col width="35%">
  <thead>
    <tr>
      <th rowspan=2>Variable</th>
      <th colspan=2>Can be set in...</th>
      <th rowspan=2>Description</th>
    </tr>
    <tr>
      <th>Management cluster YAML</th>
      <th>Tanzu Kubernetes cluster YAML</th>
    </tr>
    </thead>
  <tbody>
    <tr>
      <td><code>AZURE_CLIENT_ID</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong>. The client ID of the app for Tanzu Kubernetes Grid that you registered with Azure.</td>
    </tr>
    <tr>
      <td><code>AZURE_CLIENT_SECRET</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong>. Your Azure client secret from <a href="mgmt-clusters/azure.md#register-app">Register a Tanzu Kubernetes Grid App on Azure</a>.</td>
    </tr>
    <tr>
      <td><code>AZURE_CUSTOM_TAGS</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Comma-separated list of tags to apply to Azure resources created for the cluster. A tag is a key-value pair, for example, <code>"foo=bar, plan=prod"</code>. For more information about tagging Azure resources, see <a href="https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json">Use tags to organize your Azure resources and management hierarchy</a> and <a href="https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/tag-support">Tag support for Azure resources</a> in the Microsoft Azure documentation.</td>
    </tr>
    <tr>
      <td><code>AZURE_ENVIRONMENT</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to override the default value</strong>. The default value is <code>AzurePublicCloud</code>. Supported clouds are <code>AzurePublicCloud</code>, <code>AzureChinaCloud</code>, <code>AzureGermanCloud</code>, <code>AzureUSGovernmentCloud</code>.</td>
    </tr>
    <tr>
      <td><code>AZURE_LOCATION</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong>. The name of the Azure region in which to deploy the cluster. For example, <code>eastus</code>.</td>
    </tr>
    <tr>
      <td><code>AZURE_RESOURCE_GROUP</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. The name of the Azure resource group that you want to use for the cluster. Defaults to the <code>CLUSTER_NAME</code>. Must be unique to each cluster. <code>AZURE_RESOURCE_GROUP</code> and <code>AZURE_VNET_RESOURCE_GROUP</code> are the same by default.</td>
    </tr>
    <tr>
      <td><code>AZURE_SSH_PUBLIC_KEY_B64</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong>. Your SSH public key, created in <a href="mgmt-clusters/azure.md">Deploy a Management Cluster to Microsoft Azure</a>, converted into base64 with newlines removed. For example, <code>c3NoLXJzYSBB [...] vdGFsLmlv</code>.</td>
    </tr>
    <tr>
      <td><code>AZURE_SUBSCRIPTION_ID</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong>. The subscription ID of your Azure subscription.</td>
    </tr>
    <tr>
      <td><code>AZURE_TENANT_ID</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Required</strong>. The tenant ID of your Azure account.</td>
    </tr>
    <tr>
      <th colspan=4>Networking</th>
    </tr>
    <tr>
      <td><code>AZURE_ENABLE_ACCELERATED_NETWORKING</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Reserved for future use.</strong>. Set to <code>true</code> to enable Azure accelerated networking on VMs based on compatible Azure Tanzu Kubernetes release (TKr) images. Currently, Azure TKr do not support Azure accelerated networking.</td>
    </tr>
    <tr>
      <td><code>AZURE_ENABLE_PRIVATE_CLUSTER</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Set this to <code>true</code> to configure the cluster as private and use an Azure Internal Load Balancer (ILB) for its incoming traffic. For more information, see <a href="tanzu-k8s-clusters/azure.md#private">Azure Private Clusters</a>.</td>
    </tr>
    <tr>
      <td><code>AZURE_FRONTEND_PRIVATE_IP</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Set this if <code>AZURE_ENABLE_PRIVATE_CLUSTER</code> is <code>true</code> and you want to override the default internal load balancer address of <code>10.0.0.100</code>.</td>
    </tr>
    <tr>
      <td><code>AZURE_VNET_CIDR</code></td>
      <td rowspan=3>&#10004;</td>
      <td rowspan=3>&#10004;</td>
      <td rowspan=3><strong>Optional, set if you want to deploy the cluster to a new VNET and subnets and override the default values</strong>. By default, <code>AZURE_VNET_CIDR</code> is set to <code>10.0.0.0/16</code>, <code>AZURE_CONTROL_PLANE_SUBNET_CIDR</code> to <code>10.0.0.0/24</code>, and <code>AZURE_NODE_SUBNET_CIDR</code> to <code>10.0.1.0/24</code>.</td>
    </tr>
    <tr>
      <td><code>AZURE_CONTROL_PLANE_SUBNET_CIDR</code></td>
    </tr>
    <tr>
      <td><code>AZURE_NODE_SUBNET_CIDR</code></td>
    </tr>
    <tr>
      <td><code>AZURE_VNET_NAME</code></td>
      <td rowspan=3>&#10004;</td>
      <td rowspan=3>&#10004;</td>
      <td rowspan=3><strong>Optional, set if you want to deploy the cluster to an existing VNET and subnets or assign names to a new VNET and subnets</strong>.</td>
    </tr>
    <tr>
      <td><code>AZURE_CONTROL_PLANE_SUBNET_NAME</code></td>
    </tr>
    <tr>
      <td><code>AZURE_NODE_SUBNET_NAME</code></td>
    </tr>
    <tr>
      <td><code>AZURE_VNET_RESOURCE_GROUP</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to override the default value</strong>. The default value is set to the value of <code>AZURE_RESOURCE_GROUP</code>.</td>
    </tr>
    <tr>
      <th colspan=4>Control Plane VMs</th>
    </tr>
    <tr>
      <td><code>AZURE_CONTROL_PLANE_DATA_DISK_SIZE_GIB</code></td>
      <td rowspan=2>&#10004;</td>
      <td rowspan=2>&#10004;</td>
      <td rowspan=2><strong>Optional</strong>. Size of data disk and OS disk, as described in Azure documentation <a href="https://docs.microsoft.com/en-us/azure/virtual-machines/managed-disks-overview#data-disk">Disk roles</a>, for control plane VMs, in GB. Examples: <code>128</code>, <code>256</code>. Control plane nodes are always provisioned with a data disk.</td>
    </tr>
    <tr>
      <td><code>AZURE_CONTROL_PLANE_OS_DISK_SIZE_GIB</code></td>
    </tr>
    <tr>
      <td><code>AZURE_CONTROL_PLANE_MACHINE_TYPE</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to override the default value</strong>. An Azure VM size for the control plane node VMs, chosen to fit expected workloads. The default value is <code>Standard_D2s_v3</code>. The minimum requirement for Azure instance types is 2 CPUs and 8 GB memory.  For possible values, see the Tanzu Kubernetes Grid installer interface.</td>
    </tr>
    <tr>
      <td><code>AZURE_CONTROL_PLANE_OS_DISK_STORAGE_ACCOUNT_TYPE</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Type of Azure storage account for control plane VM disks. Example: <code>Premium_LRS</code>.</td>
    </tr>
    <tr>
      <th colspan=4>Worker Node VMs</th>
    </tr>
    <tr>
      <td><code>AZURE_ENABLE_NODE_DATA_DISK</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Set to <code>true</code> to provision a data disk for each worker node VM, as described in Azure documentation <a href="https://docs.microsoft.com/en-us/azure/virtual-machines/managed-disks-overview#data-disk">Disk roles</a>. Default: <code>false</code>.</td>
    </tr>
    <tr>
      <td><code>AZURE_NODE_DATA_DISK_SIZE_GIB</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Set this variable if <code>AZURE_ENABLE_NODE_DATA_DISK</code> is <code>true</code>.  Size of data disk, as described in Azure documentation <a href="https://docs.microsoft.com/en-us/azure/virtual-machines/managed-disks-overview#data-disk">Disk roles</a>, for worker VMs, in GB. Examples: <code>128</code>, <code>256</code>.</td>
    </tr>
    <tr>
      <td><code>AZURE_NODE_OS_DISK_SIZE_GIB</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Size of OS disk, as described in Azure documentation <a href="https://docs.microsoft.com/en-us/azure/virtual-machines/managed-disks-overview#data-disk">Disk roles</a>, for worker VMs, in GB. Examples: <code>128</code>, <code>256</code>.</td>
    </tr>
    <tr>
      <td><code>AZURE_NODE_MACHINE_TYPE</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional, set if you want to override the default value</strong>. An Azure VM size for the worker node VMs, chosen to fit expected workloads. The default value is <code>Standard_D2s_v3</code>. For possible values, see the Tanzu Kubernetes Grid installer interface.</td>
    </tr>
    <tr>
      <td><code>AZURE_NODE_OS_DISK_STORAGE_ACCOUNT_TYPE</code></td>
      <td>&#10004;</td>
      <td>&#10004;</td>
      <td><strong>Optional</strong>. Set this variable if <code>AZURE_ENABLE_NODE_DATA_DISK</code> is <code>true</code>. Type of Azure storage account for worker VM disks. Example: <code>Premium_LRS</code>.</td>
    </tr>
  </tbody>
</table>
