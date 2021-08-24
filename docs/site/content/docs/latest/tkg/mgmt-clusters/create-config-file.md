# Create a Management Cluster Configuration File

This documentation includes configuration file templates that you can use to deploy management clusters to each of vSphere, Amazon EC2, and Azure. The templates include all of the options that are relevant to deploying management clusters on a given infrastructure provider. You can copy the templates and follow the instructions in this section to update them.

Consult the [Tanzu CLI Configuration File Variable Reference](../tanzu-config-reference.md) for details about each setting. The sections below also contain links to other sections of this documentation to provide additional information.

**IMPORTANT**:

- As described in [Configuring the Management Cluster](deploy-management-clusters.md#configuring), environment variables override values from a cluster configuration file. To use all settings from a cluster configuration file, unset any conflicting environment variables before you deploy the management cluster from the CLI.
- Tanzu Kubernetes Grid does not support IPv6 addresses. This is because upstream Kubernetes only provides alpha support for IPv6. Always provide IPv4 addresses in settings in the configuration file.
- Some parameters configure identical properties. For example, the `SIZE` property configures the same infrastructure settings as all of the control plane and worker node size and type properties for the different infrastructure providers, but at a more general level. In such cases, avoid setting conflicting or redundant properties.  

## Create the Configuration File

1. Copy and paste the contents of the template for your infrastructure provider into a text editor.

   Copy a template from one of the following locations:

   - [Management Cluster Configuration for vSphere](config-vsphere.md)
   - [Management Cluster Configuration for Amazon EC2](config-aws.md)
   - [Management Cluster Configuration for Microsoft Azure](config-azure.md)

   For example, if you have already deployed a management cluster from the installer interface, you can save the file in the default location for cluster configurations,  `~/.tanzu/tkg/clusterconfigs`.

1. Save the file with a `.yaml` extension and an appropriate name, for example `aws-mgmt-cluster-config.yaml`.

The subsequent sections describe how to update the settings that are common to all infrastructure providers as well as the settings that are specific to each of vSphere, Amazon EC2, and Azure.

## <a id="basic"></a> Configure Basic Management Cluster Creation Information

The basic management cluster creation settings define the infrastructure on which to deploy the management cluster and other basic settings. They are common to all infrastructure providers.

- For `CLUSTER_PLAN` specify whether you want to deploy a development cluster, which provides a single control plane node, or a production cluster, which provides a highly available management cluster with three control plane nodes. Specify `dev` or `prod`.
- For `INFRASTRUCTURE_PROVIDER`, specify `aws`, `azure`, or `vsphere`.

   ```
   INFRASTRUCTURE_PROVIDER: aws
   ```

   ```
   INFRASTRUCTURE_PROVIDER: azure
   ```

   ```
   INFRASTRUCTURE_PROVIDER: vsphere
   ```

- Optionally disable participation in the VMware Customer Experience Improvement Program (CEIP) by setting `ENABLE_CEIP_PARTICIPATION` to `false`. For information about the CEIP, see [Opt in or Out of the VMware CEIP](../cluster-lifecycle/multiple-management-clusters.md#ceip) and [https://www.vmware.com/solutions/trustvmware/ceip.html](https://www.vmware.com/solutions/trustvmware/ceip.html).
- Optionally uncomment and update `TMC_REGISTRATION_URL` to register the management cluster with Tanzu Mission Control. For information about Tanzu Mission Control, see [Register Your Management Cluster with Tanzu Mission Control](register_tmc.md).
- Optionally disable audit logging by setting `ENABLE_AUDIT_LOGGING` to `false`. For information about  audit logging, see [Audit Logging](../troubleshooting-tkg/audit-logging.md).
- If the recommended CIDR ranges of 100.64.0.0/13 and 100.96.0.0/11 are unavailable, update `CLUSTER_CIDR` for the cluster pod network and `SERVICE_CIDR` for the cluster service network.

For example:

```
#! ---------------------------------------------------------------------
#! Basic cluster creation configuration
#! ---------------------------------------------------------------------

CLUSTER_NAME: aws-mgmt-cluster
CLUSTER_PLAN: dev
INFRASTRUCTURE_PROVIDER: aws
ENABLE_CEIP_PARTICIPATION: true 
TMC_REGISTRATION_URL: https://tmc-org.cloud.vmware.com/installer?id=[...]&source=registration
ENABLE_AUDIT_LOGGING: true
CLUSTER_CIDR: 100.96.0.0/11
SERVICE_CIDR: 100.64.0.0/13      
```

## <a id="identity-mgmt"></a> Configure Identity Management

Set `IDENTITY_MANAGEMENT_TYPE` to <code>ldap</code> or <code>oidc</code>. Set <code>none</code> to disable identity management. It is strongly recommended to enable identity management for production deployments.

For information identity management in Tanzu Kubernetes Grid, and the pre-deployment steps to perform, see [Configure Identity Management After Management Cluster Deployment](configure-id-mgmt.md).

```
IDENTITY_MANAGEMENT_TYPE: oidc
```

```
IDENTITY_MANAGEMENT_TYPE: ldap
```

### OIDC

To configure OIDC, update the variables below. For information about how to configure the variables, see [Identity Providers - OIDC](../tanzu-config-reference.md#identity-management-oidc) in the *Tanzu CLI Configuration File Variable Reference*.

For example:

```
OIDC_IDENTITY_PROVIDER_CLIENT_ID: 0oa2i[...]NKst4x7
OIDC_IDENTITY_PROVIDER_CLIENT_SECRET: <encoded:LVVnMFNsZFIy[...]TMTV3WUdPZDJ2Xw==>
OIDC_IDENTITY_PROVIDER_GROUPS_CLAIM: groups
OIDC_IDENTITY_PROVIDER_ISSUER_URL: https://dev-[...].okta.com
OIDC_IDENTITY_PROVIDER_SCOPES: openid,groups,email
OIDC_IDENTITY_PROVIDER_USERNAME_CLAIM: email
```

### LDAP

To configure LDAP, uncomment and update the `LDAP_*` variables with information about your LDAPS server. For information about how to configure the variables, see [Identity Providers - LDAP](../tanzu-config-reference.md#identity-management-ldap) in the *Tanzu CLI Configuration File Variable Reference*.

For example:

```
LDAP_BIND_DN: ""
LDAP_BIND_PASSWORD: ""
LDAP_GROUP_SEARCH_BASE_DN: dc=example,dc=com
LDAP_GROUP_SEARCH_FILTER: (objectClass=posixGroup)
LDAP_GROUP_SEARCH_GROUP_ATTRIBUTE: memberUid
LDAP_GROUP_SEARCH_NAME_ATTRIBUTE: cn
LDAP_GROUP_SEARCH_USER_ATTRIBUTE: uid
LDAP_HOST: ldaps.example.com:636
LDAP_ROOT_CA_DATA_B64: ""
LDAP_USER_SEARCH_BASE_DN: ou=people,dc=example,dc=com
LDAP_USER_SEARCH_FILTER: (objectClass=posixAccount)
LDAP_USER_SEARCH_NAME_ATTRIBUTE: uid
LDAP_USER_SEARCH_USERNAME: uid
```

## <a id="proxies"></a> Configure Proxies

To optionally send outgoing HTTP(S) traffic from the management cluster to a proxy, for example in an internet-restricted environment, uncomment and set the `*_PROXY` settings.
The proxy settings are common to all infrastructure providers. You can choose to use one proxy for HTTP requests and another proxy for HTTPS requests or to use the same proxy for both HTTP and HTTPS requests.

* (**Required**) `TKG_HTTP_PROXY`: This is the URL of the proxy that handles HTTP requests. To set the URL, use the format below:

   ```
   PROTOCOL://USERNAME:PASSWORD@FQDN-OR-IP:PORT
   ```

   Where:

   * (**Required**) `PROTOCOL`: This must be `http`.
   * (**Optional**) `USERNAME` and `PASSWORD`: This is your HTTP proxy username and password. You must set `USERNAME` and `PASSWORD` if the proxy requires authentication.
   * (**Required**) `FQDN-OR-IP`: This is the FQDN or IP address of your HTTP proxy.
   * (**Required**) `PORT`: This is the port number that your HTTP proxy uses.

   For example, `http://user:password@myproxy.com:1234`.

* (**Required**) `TKG_HTTPS_PROXY`: This is the URL of the proxy that handles HTTPS requests. You can set `TKG_HTTPS_PROXY` to the same value as `TKG_HTTP_PROXY` or provide a different value. To set the value, use the URL format from the previous step, where:

   * (**Required**) `PROTOCOL`: This must be `http`.
   * (**Optional**) `USERNAME` and `PASSWORD`: This is your HTTPS proxy username and password. You must set `USERNAME` and `PASSWORD` if the proxy requires authentication.
   * (**Required**) `FQDN-OR-IP`: This is the FQDN or IP address of your HTTPS proxy.
   * (**Required**) `PORT`: This is the port number that your HTTPS proxy uses.

   For example, `http://user:password@myproxy.com:1234`.

* (**Optional**) `TKG_NO_PROXY`: This sets one or more comma-separated network CIDRs or hostnames that must bypass the HTTP(S) proxy, for example to enable the management cluster to communicate directly with infrastructure that runs on the same network, behind the same proxy.
Do not use spaces. For example, `noproxy.yourdomain.com,192.168.0.0/24`.

   Internally, Tanzu Kubernetes Grid appends `localhost`, `127.0.0.1`, the values of `CLUSTER_CIDR` and `SERVICE_CIDR`, `.svc`, and `.svc.cluster.local` to the value that you set in `TKG_NO_PROXY`. It also appends your AWS VPC CIDR and `169.254.0.0/16` for deployments to Amazon EC2 and your Azure VNET CIDR, `169.254.0.0/16`, and `168.63.129.16` for deployments to Azure. For vSphere, you must manually add the CIDR of `VSPHERE_NETWORK`, which includes the IP address of your control plane endpoint, to `TKG_NO_PROXY`. If you set `VSPHERE_CONTROL_PLANE_ENDPOINT` to an FQDN, add both the FQDN and `VSPHERE_NETWORK` to `TKG_NO_PROXY`.

   **Important:** If the cluster VMs need to communicate with external services and infrastructure endpoints in your Tanzu Kubernetes Grid environment, ensure that those endpoints are reachable by the proxies that you set above or add them to `TKG_NO_PROXY`. Depending on your environment configuration, this may include, but is not limited to:

   * Your OIDC or LDAP server
   * Harbor
   * NSX-T
   * NSX Advanced Load Balancer
   * AWS VPC CIDRs that are external to the cluster

For example:

```
#! ---------------------------------------------------------------------
#! Proxy configuration
#! ---------------------------------------------------------------------

TKG_HTTP_PROXY: "http://myproxy.com:1234"
TKG_HTTPS_PROXY: "http://myproxy.com:1234"
TKG_NO_PROXY: "noproxy.yourdomain.com,192.168.0.0/24"
```

## <a id="nodes"></a> Configure Node Settings

By default, all cluster nodes run Ubuntu v20.04, for all infrastructure providers. On vSphere you can optionally deploy clusters that run Photon OS on their nodes. On Amazon EC2, nodes can optionally run Amazon Linux 2. For the architecture, the default and only current choice is `amd64`. For the OS and version settings, see see [Node Configuration](../tanzu-config-reference.md#size) in the *Tanzu CLI Configuration File Variable Reference*.

For example:

```
#! ---------------------------------------------------------------------
#! Node configuration 
#! ---------------------------------------------------------------------

OS_NAME: "photon"
OS_VERSION: "3"
OS_ARCH: "amd64"
```

How you set node compute configuration and sizes depends on the infrastructure provider. For information, see [Management Cluster Configuration for vSphere](config-vsphere.md), [Management Cluster Configuration for Amazon EC2](config-aws.md), or [Management Cluster Configuration for Microsoft Azure](config-azure.md).

## <a id="mhc"></a> Configure Machine Health Checks

Optionally update variables based on your deployment preferences and using the guidelines described in the Configuration Parameter Reference. Alternatively, disable Machine Health Checks by setting `ENABLE_MHC: "false"`.

For information about how to configure the Machine Health Check settings, see [Machine Health Checks](../tanzu-config-reference.md#mhc) in the *Tanzu CLI Configuration File Variable Reference* and [Configure Machine Health Checks for Tanzu Kubernetes Clusters](../cluster-lifecycle/configure-health-checks.md).

For example:

```
ENABLE_MHC: "true"
MHC_UNKNOWN_STATUS_TIMEOUT: 10m
MHC_FALSE_STATUS_TIMEOUT: 20m
```

## <a id="registry"></a> Configure a Private Image Registry

If you are deploying the management cluster in an Internet-restricted environment, uncomment and update the `TKG_CUSTOM_IMAGE_REPOSITORY_*` settings. If you are deploying the management cluster in an environment that has access to the external internet, you do not need to configure these settings.

For information about deployments in Internet-restricted environments, see [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md). The private image registry settings are common to all infrastructure providers.

For example:

```
#! ---------------------------------------------------------------------
#! Image repository configuration
#! ---------------------------------------------------------------------

TKG_CUSTOM_IMAGE_REPOSITORY: "custom-image-repository.io/yourproject"
TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE: "LS0t[...]tLS0tLQ=="
```

## <a id="antrea"></a>  Configure Antrea CNI

By default, clusters that you deploy with the Tanzu CLI provide in-cluster container networking with the Antrea container network interface (CNI).

You can optionally disable Source Network Address Translation (SNAT) for pod traffic, implement `hybrid`, `noEncap`, `NetworkPolicyOnly` traffic encapsulation modes, use proxies and network policies, and implement Traceflow.

For more information about Antrea, see the following resources:

- [VMware Container Networking with Antrea](https://www.vmware.com/products/antrea-container-networking.html) product page on vmware.com.
- [Antrea open source project page](https://antrea.io/)
- [Antrea documentation](https://antrea.io/docs/v0.11.3/)
- [Deploying Antrea for Kubernetes Networking](https://www.vmware.com/content/dam/digitalmarketing/vmware/en/pdf/products/vmware-deploy-antrea-k8s-networking.pdf) whitepaper on vmware.com.

To optionally configure these features on Antrea, uncomment and update the `ANTREA_*` variables. For example:

```
#! ---------------------------------------------------------------------
#! Antrea CNI configuration
#! ---------------------------------------------------------------------

ANTREA_NO_SNAT: true
ANTREA_TRAFFIC_ENCAP_MODE: "hybrid"
ANTREA_PROXY: true
ANTREA_POLICY: true 
ANTREA_TRACEFLOW: false
```

## What to Do Next

Continue to update the configuration file settings for vSphere, Amazon EC2, or Azure. For the configuration file settings that are specific to each infrastructure provider, see the corresponding topic:

- [Management Cluster Configuration for vSphere](config-vsphere.md)
- [Management Cluster Configuration for Amazon EC2](config-aws.md)
- [Management Cluster Configuration for Microsoft Azure](config-azure.md)
