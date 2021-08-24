# Register Core Add-ons

This topic describes how to register the CNI, vSphere CPI, vSphere CSI, Pinniped, and Metrics Server add-ons with `tanzu-addons-manager`, the component that manages the lifecycle of add-ons. Skip this step if:

* Your management and Tanzu Kubernetes (workload) clusters were created using Tanzu Kubernetes Grid v1.3.0 or later.
* You already registered the add-ons by following the instructions in this topic.

## About Add-on Lifecycle Management in Tanzu Kubernetes Grid

When you create a management or a workload cluster using Tanzu Kubernetes Grid v1.3.0 or later, it automatically installs the following core add-ons in the cluster:

* CNI: `cni/calico` or `cni/antrea`
* (**vSphere only**) vSphere CPI: `cloud-provider/vsphere-cpi`
* (**vSphere only**) vSphere CSI: `csi/vsphere-csi`
* Authentication: `authentication/pinniped`
* Metrics Server: `metrics/metrics-server`

Tanzu Kubernetes Grid manages the lifecycle of these add-ons. For example, it automatically upgrades the add-ons when you upgrade your management and workload clusters using the `tanzu management-cluster upgrade` and `tanzu cluster upgrade` commands. This ensures that your Tanzu Kubernetes Grid version and add-on versions are compatible.

### Upgrades from Tanzu Kubernetes Grid v1.2.x to v1.3.x

Upgrading your management and workload clusters from Tanzu Kubernetes Grid v1.2.x to v1.3.x does not automatically upgrade the CNI, vSphere CPI, and vSphere CSI add-ons. To enable automatic lifecycle management for these add-ons, you must manually register them with the `tanzu-addons-manager`. The Pinniped and Metrics Server add-ons are new components in Tanzu Kubernetes Grid v1.3.0. You must enable them if you want to use identity management with Pinniped and Metrics Server in your upgraded clusters.

**Important:** The Dex and Gangway extensions are deprecated in Tanzu Kubernetes Grid v1.3.0 and will be removed in a future release. It is strongly recommended to migrate any existing clusters that implement the Dex and Gangway extensions to the integrated Pinniped authentication service.

## Prerequisites

Before following the instructions in this topic, confirm that:

* Your management and workload clusters have been upgraded to Tanzu Kubernetes Grid v1.3.x.
* `tanzu-addons-controller-manager` and `kapp-controller` are running in your management cluster, by using `kubectl get pods -n tkg-system`.

## Register the Core Add-ons

After you upgrade your management and workload clusters to Tanzu Kubernetes Grid v1.3.x, follow the instructions below to register the core add-ons with `tanzu-addons-manager`:

* To register the CNI, vSphere CPI, and vSphere CSI add-ons, see [Register the CNI, vSphere CPI, and vSphere CSI Add-ons](#existing-components).
* To register the Pinniped and Metrics Server add-ons, see [Enable the Pinniped and Metrics Server Add-ons](#new-components).

### <a id="existing-components"></a> Register the CNI, vSphere CPI, and vSphere CSI Add-ons

When registering the CNI, vSphere CPI, or vSphere CSI add-on with `tanzu-addons-manager`, register the add-on that is running in the management cluster first and then in each workload cluster.

To register the CNI, vSphere CPI, or vSphere CSI add-on that is running in a management or a workload cluster:

1. Create a configuration file for the cluster.

   1. Set the following variables. For example, for vSphere:

      ```
      # This is the name of your target management or workload cluster.
      CLUSTER_NAME=YOUR-CLUSTER-NAME
      # For the management cluster, the namespace must be "tkg-system". For workload clusters, the default namespace is "default".
      NAMESPACE=YOUR-CLUSTER-NAMESPACE
      # CLUSTER_PLAN can be "dev", "prod", and so on.
      CLUSTER_PLAN=YOUR-CLUSTER-PLAN

      # If you are creating the configuration file for a management cluster, you can use the following commands to retrieve the values of CLUSTER_CIDR and SERVICE_CIDR from your management cluster.
      CLUSTER_CIDR=$(kubectl get kubeadmconfig -n tkg-system -l cluster.x-k8s.io/cluster-name=$CLUSTER_NAME,cluster.x-k8s.io/control-plane=  -o jsonpath='{.items[0].spec.clusterConfiguration.networking.podSubnet}')
      SERVICE_CIDR=$(kubectl get kubeadmconfig -n tkg-system -l cluster.x-k8s.io/cluster-name=$CLUSTER_NAME,cluster.x-k8s.io/control-plane=  -o jsonpath='{.items[0].spec.clusterConfiguration.networking.serviceSubnet}')

      # If you are creating the configuration file for a workload cluster, you can use the following commands to retrieve the values of CLUSTER_CIDR and SERVICE_CIDR from your workload cluster.
      CLUSTER_CIDR=$(kubectl get cluster "${CLUSTER_NAME}" -n "${NAMESPACE}" -o jsonpath='{.spec.clusterNetwork.pods.cidrBlocks[0]}')
      SERVICE_CIDR=$(kubectl get cluster "${CLUSTER_NAME}" -n "${NAMESPACE}" -o jsonpath='{.spec.clusterNetwork.services.cidrBlocks[0]}')

      # Set these variables if your cluster is running on vSphere. You can use the commands bellow to retrieve their values from the cluster. If your cluster is running on Amazon EC2 or Azure, set the AWS_ or AZURE_ variables that you configured when you deployed the cluster.
      VSPHERE_SERVER=$(kubectl get VsphereCluster "${CLUSTER_NAME}" -n "${NAMESPACE}" -o jsonpath='{.spec.server}')
      VSPHERE_DATACENTER=$(kubectl get VsphereMachineTemplate "${CLUSTER_NAME}-control-plane" -n "${NAMESPACE}" -o jsonpath='{.spec.template.spec.datacenter}')
      VSPHERE_RESOURCE_POOL=$(kubectl get VsphereMachineTemplate "${CLUSTER_NAME}-control-plane" -n "${NAMESPACE}" -o jsonpath='{.spec.template.spec.resourcePool}')
      VSPHERE_DATASTORE=$(kubectl get VsphereMachineTemplate "${CLUSTER_NAME}-control-plane" -n "${NAMESPACE}" -o jsonpath='{.spec.template.spec.datastore}')
      VSPHERE_FOLDER=$(kubectl get VsphereMachineTemplate "${CLUSTER_NAME}-control-plane" -n "${NAMESPACE}" -o jsonpath='{.spec.template.spec.folder}')
      VSPHERE_NETWORK=$(kubectl get VsphereMachineTemplate "${CLUSTER_NAME}-control-plane" -n "${NAMESPACE}" -o jsonpath='{.spec.template.spec.network.devices[0].networkName}')
      VSPHERE_SSH_AUTHORIZED_KEY=$(kubectl get KubeadmControlPlane "${CLUSTER_NAME}-control-plane" -n "${NAMESPACE}" -o jsonpath='{.spec.kubeadmConfigSpec.users[0].sshAuthorizedKeys[0]}')
      VSPHERE_TLS_THUMBPRINT=$(kubectl get VsphereCluster "${CLUSTER_NAME}" -n "${NAMESPACE}" -o jsonpath='{.spec.thumbprint}')
      VSPHERE_INSECURE=TRUE-OR-FALSE
      VSPHERE_USERNAME='YOUR-VSPHERE-USERNAME'
      VSPHERE_PASSWORD='YOUR-VSPHERE-PASSWORD'
      VSPHERE_CONTROL_PLANE_ENDPOINT=FQDN-OR-IP
      ```

   1. If your cluster uses Calico as the CNI provider, add or include the following line in the cluster configuration file:

      ```
      CNI: calico
      ```

   1. Create the configuration file. You can use the `echo` command to write the variables that you set above and their values to the file. For example:

      ```
      echo "CLUSTER_CIDR: ${CLUSTER_CIDR}" >> config.yaml
      ```

1. Set the `_TKG_CLUSTER_FORCE_ROLE` environment variable to `management` or `workload`. For example:

   ```
   export _TKG_CLUSTER_FORCE_ROLE="management"
   ```

   ```
   export _TKG_CLUSTER_FORCE_ROLE="workload"
   ```

   On Windows, use the `SET` command.

1. Register the add-ons.

   * **CNI add-on**:

      1. Set the `FILTER_BY_ADDON_TYPE` and `REMOVE_CRS_FOR_ADDON_TYPE` environment variables to the values below:

         ```
         export FILTER_BY_ADDON_TYPE="cni/antrea"
         ```

         ```
         export REMOVE_CRS_FOR_ADDON_TYPE="cni/antrea"
         ```

         If your cluster uses Calico, replace `cni/antrea` with `cni/calico`.

      1. Generate a manifest file for the add-on. For example:

         ```
         tanzu cluster create ${CLUSTER_NAME} --dry-run -f config.yaml > ${CLUSTER_NAME}-addon-manifest.yaml
         ```

         Where `CLUSTER_NAME` is the name of your target management or workload cluster.

      1. Review the manifest and then apply it to the management cluster.

         ```
         kubectl apply –f ${CLUSTER_NAME}-addon-manifest.yaml
         ```

   * **vSphere CPI add-on**:

      1. Set the `FILTER_BY_ADDON_TYPE` and `REMOVE_CRS_FOR_ADDON_TYPE` environment variables to the values below.

         ```
         export FILTER_BY_ADDON_TYPE="cloud-provider/vsphere-cpi"
         ```

         ```
         export REMOVE_CRS_FOR_ADDON_TYPE="cloud-provider/vsphere-cpi"
         ```

      1. Generate a manifest file for the add-on. For example:

         ```
         tanzu cluster create ${CLUSTER_NAME} --dry-run -f config.yaml > ${CLUSTER_NAME}-addon-manifest.yaml
         ```

         Where `CLUSTER_NAME` is the name of your target management or workload cluster.

      1. Review the manifest and then apply it to the management cluster.

         ```
         kubectl apply –f ${CLUSTER_NAME}-addon-manifest.yaml
         ```

   * **vSphere CSI add-on**:

      1. Set the `FILTER_BY_ADDON_TYPE` and `REMOVE_CRS_FOR_ADDON_TYPE` environment variables to the values below.

         ```
         export FILTER_BY_ADDON_TYPE="csi/vsphere-csi"
         ```

         ```
         export REMOVE_CRS_FOR_ADDON_TYPE="csi/vsphere-csi"
         ```

      1. Generate a manifest file for the add-on. For example:

         ```
         tanzu cluster create ${CLUSTER_NAME} --dry-run -f config.yaml > ${CLUSTER_NAME}-addon-manifest.yaml
         ```

         Where `CLUSTER_NAME` is the name of your target management or workload cluster.

      1. Review the manifest and then apply it to the management cluster.

         ```
         kubectl apply –f ${CLUSTER_NAME}-addon-manifest.yaml
         ```

### <a id="new-components"></a> Enable the Pinniped and Metrics Server Add-ons

To enable the Pinniped and Metrics Server add-ons, follow the instructions below:

* [Pinniped Add-on](#pinniped): Follow these instructions if you want to enable identity management with Pinniped in your upgraded clusters.
* [Metrics Server Add-on](#metrics-server): Follow these instructions if you want to enable Metrics Server in your upgraded clusters.

#### <a id="pinniped"></a> Pinniped Add-on

To enable identity management with Pinniped, you must enable the Pinniped add-on in your upgraded management and workload clusters. For information about identity management in Tanzu Kubernetes Grid v1.3.x, see [Enabling Identity Management in Tanzu Kubernetes Grid](../mgmt-clusters/enabling-id-mgmt.md).

**Important:** The Dex and Gangway extensions are deprecated in Tanzu Kubernetes Grid v1.3.0 and will be removed in a future release. It is strongly recommended to migrate any existing clusters that implement the Dex and Gangway extensions to the integrated Pinniped authentication service. If you implemented the Dex and Gangway extensions in Tanzu Kubernetes Grid v1.2.x, delete Dex and Gangway from both the management cluster and workload clusters before enabling the Pinniped add-on. Back up your Dex and Gangway configuration settings, such as `ConfigMap`.

To enable identity management with Pinniped:

1. Create a configuration file for your management cluster as described in step 1 in [Register the CNI, vSphere CPI, and vSphere CSI Add-ons](#existing-components) above.

1. Obtain your OIDC or LDAP identity provider details and add the following settings to the configuration file.

   ```
   # Identity management type. This must be "oidc" or "ldap".

   IDENTITY_MANAGEMENT_TYPE:

   # Set these variables if you want to configure OIDC.

   CERT_DURATION: 2160h
   CERT_RENEW_BEFORE: 360h
   OIDC_IDENTITY_PROVIDER_ISSUER_URL:
   OIDC_IDENTITY_PROVIDER_CLIENT_ID:
   OIDC_IDENTITY_PROVIDER_CLIENT_SECRET:
   OIDC_IDENTITY_PROVIDER_SCOPES: "email,profile,groups"
   OIDC_IDENTITY_PROVIDER_USERNAME_CLAIM:
   OIDC_IDENTITY_PROVIDER_GROUPS_CLAIM:

   # Set these variables if you want to configure LDAP.

   LDAP_BIND_DN:
   LDAP_BIND_PASSWORD:
   LDAP_HOST:
   LDAP_USER_SEARCH_BASE_DN:
   LDAP_USER_SEARCH_FILTER:
   LDAP_USER_SEARCH_USERNAME: userPrincipalName
   LDAP_USER_SEARCH_ID_ATTRIBUTE: DN
   LDAP_USER_SEARCH_EMAIL_ATTRIBUTE: DN
   LDAP_USER_SEARCH_NAME_ATTRIBUTE:
   LDAP_GROUP_SEARCH_BASE_DN:
   LDAP_GROUP_SEARCH_FILTER:
   LDAP_GROUP_SEARCH_USER_ATTRIBUTE: DN
   LDAP_GROUP_SEARCH_GROUP_ATTRIBUTE:
   LDAP_GROUP_SEARCH_NAME_ATTRIBUTE: cn
   LDAP_ROOT_CA_DATA_B64:
   ```

   For more information about these variables, see [Variables for Configuring Identity Providers - OIDC](../tanzu-config-reference.md#identity-management-oidc) and [Variables for Configuring Identity Providers - LDAP](../tanzu-config-reference.md#identity-management-ldap).

1. Set the `_TKG_CLUSTER_FORCE_ROLE` environment variable to `management`.

   ```
   export _TKG_CLUSTER_FORCE_ROLE="management"
   ```

   On Windows, use the `SET` command.

1. Set the `FILTER_BY_ADDON_TYPE` environment variable to `authentication/pinniped`.

   ```
   export FILTER_BY_ADDON_TYPE="authentication/pinniped"
   ```

1. Generate a manifest file for the add-on. For example:

   ```
   tanzu cluster create ${CLUSTER_NAME} --dry-run -f config.yaml > ${CLUSTER_NAME}-addon-manifest.yaml
   ```

   Where `CLUSTER_NAME` is the name of your target management cluster.

1. Review the manifest and then apply it to the cluster. For example:

   ```
   kubectl apply –f ${CLUSTER_NAME}-addon-manifest.yaml
   ```

   If you configured your management cluster to use OIDC authentication above, you must provide the callback URI for the management cluster to your OIDC identity provider. For more information, see [Configure Identity Management After Management Cluster Deployment](../mgmt-clusters/configure-id-mgmt.md).

1. Enable the Pinniped add-on in each workload cluster that is managed by your management cluster. For each cluster, follow these steps:

   1. Create a configuration file for the cluster as described in step 1 in [Register the CNI, vSphere CPI, and vSphere CSI Add-ons](#existing-components) above.

   1. Add the following variables to the configuration file.

      ```
      # This is the Pinniped supervisor service endpoint in the management cluster.
      SUPERVISOR_ISSUER_URL:

      # Pinniped uses this b64-encoded CA bundle data for communication between the management cluster and the workload cluster.
      SUPERVISOR_ISSUER_CA_BUNDLE_DATA_B64:
      ```

      You can retrieve these values by running `kubectl get configmap pinniped-info -n kube-public -o yaml` against the management cluster.

   1. Set the `_TKG_CLUSTER_FORCE_ROLE` environment variable to `workload`.

      ```
      export _TKG_CLUSTER_FORCE_ROLE="workload"
      ```

   1. Set the `FILTER_BY_ADDON_TYPE` environment variable to `authentication/pinniped`.

      ```
      export FILTER_BY_ADDON_TYPE="authentication/pinniped"
      ```

   1. Generate a manifest file for the add-on. For example:

      ```
      tanzu cluster create ${CLUSTER_NAME} --dry-run -f config.yaml > ${CLUSTER_NAME}-addon-manifest.yaml
      ```

      Where `CLUSTER_NAME` is the name of your target workload cluster.

   1. Review the manifest and then apply it to the management cluster:

      ```
      kubectl apply –f ${CLUSTER_NAME}-addon-manifest.yaml
      ```

      For information about how to grant users access to workload clusters on which you have implemented identity management, see [Authenticate Connections to a Workload Cluster](../cluster-lifecycle/connect.html#id-mgmt).

#### <a id="metrics-server"></a> Metrics Server Add-on

When enabling the Metrics Server add-on, enable the add-on in the management cluster first and then in each workload cluster.
To enable the Metrics Server add-on in a management or a workload cluster:

1. Create a configuration file for the cluster as described in step 1 in [Register the CNI, vSphere CPI, and vSphere CSI Add-ons](#existing-components) above.

1. Set the `_TKG_CLUSTER_FORCE_ROLE` environment variable to `management` or `workload`. For example:

   ```
   export _TKG_CLUSTER_FORCE_ROLE="management"
   ```

   ```
   export _TKG_CLUSTER_FORCE_ROLE="workload"
   ```

   On Windows, use the `SET` command.

1. Set the `FILTER_BY_ADDON_TYPE` environment variable to `metrics/metrics-server`.

   ```
   export FILTER_BY_ADDON_TYPE="metrics/metrics-server"
   ```

1. Generate a manifest file for the add-on. For example:

   ```
   tanzu cluster create ${CLUSTER_NAME} --dry-run -f config.yaml > ${CLUSTER_NAME}-addon-manifest.yaml
   ```

   Where `CLUSTER_NAME` is the name of your target management or workload cluster.

1. Review the manifest and then apply it to the management cluster. For example:

   ```
   kubectl apply –f ${CLUSTER_NAME}-addon-manifest.yaml
   ```
