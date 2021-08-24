# Enable Identity Management After Management Cluster Deployment

This topic describes how to enable identity management in Tanzu Kubernetes Grid as a post-deployment step.

## <a id="overview"></a> Overview

If you did not configure identity management when you deployed your management cluster,
you can enable it as a post-deployment step by following the instructions below:

1. [Obtain your identity provider details.](#idp)
1. [Generate a Kubernetes secret for the Pinniped add-on.](#secret)
1. [Check the status of the identity management service.](#status)
1. [(OIDC only) Provide the callback URI to your OIDC provider.](#callback)
1. [Generate a non-admin `kubeconfig`.](#kubeconfig)
1. [Create role bindings for your management cluster users.](#role-binding)
1. [Enable identity management in workload clusters.](#workload)

You must enable identity management in the management cluster first and then in each Tanzu Kubernetes (workload) cluster that it manages.

## <a id="idp"></a> Obtain Your Identity Provider Details

Before you can enable identity management, you must have an identity provider.
Tanzu Kubernetes Grid supports LDAPS and OIDC identity providers. For more information,
see [Obtain Your Identity Provider Details](../mgmt-clusters/enabling-id-mgmt.md#idp) in _Enabling Identity Management in Tanzu Kubernetes Grid_ and then return to this topic.

## <a id="secret"></a> Generate a Kubernetes Secret for the Pinniped Add-on

This procedure configures the Pinniped add-on in your management cluster.

To generate a Kubernetes secret for the Pinniped add-on:

1. Create a cluster configuration file using the configuration settings that you defined when you deployed your management cluster. Include the following configuration variables in the file:

   * Basic cluster variables.

      ```
      # This is the name of your target management cluster.
      CLUSTER_NAME:
      # For the management cluster, the default namespace is "tkg-system".
      NAMESPACE:
      CLUSTER_PLAN:
      CLUSTER_CIDR:
      SERVICE_CIDR:
      ```

   * vSphere-, AWS-, or Azure-specific variables that you set when you deployed your management cluster. For information about these variables, see [Tanzu CLI Configuration File Variable Reference](../tanzu-config-reference.md).

      For example:

      ```
      VSPHERE_SERVER:
      VSPHERE_DATACENTER:
      VSPHERE_RESOURCE_POOL:
      VSPHERE_DATASTORE:
      VSPHERE_FOLDER:
      VSPHERE_NETWORK:
      VSPHERE_SSH_AUTHORIZED_KEY:
      VSPHERE_TLS_THUMBPRINT:
      VSPHERE_INSECURE:
      VSPHERE_USERNAME:
      VSPHERE_PASSWORD:
      VSPHERE_CONTROL_PLANE_ENDPOINT:
      ```

   * OIDC or LDAP identity provider details.

      **Note:** Set these variables only for the management cluster. You do not need to set them for your workload clusters.

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

      For more information, see [Variables for Configuring Identity Providers - OIDC](../tanzu-config-reference.md#identity-management-oidc) and [Variables for Configuring Identity Providers - LDAP](../tanzu-config-reference.md#identity-management-ldap).

1. Set the `_TKG_CLUSTER_FORCE_ROLE` environment variable to `management`.

   ```
   export _TKG_CLUSTER_FORCE_ROLE="management"
   ```

   On Windows, use the `SET` command.

1. Set the `FILTER_BY_ADDON_TYPE` environment variable to `authentication/pinniped`.

   ```
   export FILTER_BY_ADDON_TYPE="authentication/pinniped"
   ```

1. Generate a manifest file for the cluster.

   ```
   tanzu cluster create CLUSTER-NAME --dry-run -f CLUSTER-CONFIG-FILE > CLUSTER-NAME-example-secret.yaml
   ```

   Where:

   * `CLUSTER-NAME` is the name of your target management cluster.
   * `CLUSTER-CONFIG-FILE` is the configuration file that you created above.

   The resulting manifest contains only the secret for the Pinniped add-on.

1. Review the secret and then apply it to the cluster. For example:

   ```
   kubectl apply -f CLUSTER-NAME-example-secret.yaml
   ```

1. After applying the secret, check the status of the Pinniped add-on by running the `kubectl get app` command.

   ```
   $ kubectl get app pinniped -n tkg-system
   NAME           DESCRIPTION             SINCE-DEPLOY    AGE
   pinniped       Reconcile succeeded     3m23s           7h50m
   ```

   If the returned status is `Reconcile failed`, run the following command to get details on the failure.

   ```
   kubectl get app pinniped -n tkg-system -o yaml
   ```

   For more information about troubleshooting the Pinniped add-on, see [Troubleshooting Core Add-on Configuration](./update-addons.md#troubleshooting) in _Update and Troubleshoot Core Add-On Configuration_.

## <a id="status"></a> Check the Status of the Identity Management Service

Confirm that the Pinniped service is running correctly. To check the status of the Pinniped service:

* If you are configuring an OIDC identity provider, follow the instructions in [Check the Status of an OIDC Identity Management Service](../mgmt-clusters/configure-id-mgmt.md#oidc) and then return to this topic.
* If you are configuring a LDAP identity provider, follow the instructions in [Check the Status of an LDAP Identity Management Service](../mgmt-clusters/configure-id-mgmt.md#ldap) and then return to this topic.

## <a id="callback"></a> (OIDC Only) Provide the Callback URI to the OIDC Provider

If you are configuring your management cluster to use OIDC authentication, you must provide the callback URI for the management cluster to your OIDC identity provider. To configure the callback URI, follow the instructions in [Provide the Callback URI to the OIDC Provider](../mgmt-clusters/configure-id-mgmt.md#callback) and then return to this topic.

## <a id="kubeconfig"></a> Generate a Non-Admin `kubeconfig`

To allow authenticated users to connect to the management cluster, generate a non-admin `kubeconfig`. To generate the non-admin `kubeconfig` file, follow the instructions in [Generate a `kubeconfig` to Allow Authenticated Users to Connect to the Management Cluster](../mgmt-clusters/configure-id-mgmt.md#gen-kubeconfig) and then return to this topic.

## <a id="role-binding"></a> Create Role Bindings for Your Management Cluster Users

To complete the identity management configuration of the management cluster, you must create role bindings for the users who use the `kubeconfig` that you generated in the above step. To create a role binding, follow the instructions in [Create a Role Binding on the Management Cluster](../mgmt-clusters/configure-id-mgmt.md#create-rolebinding) and then return to this topic.

## <a id="workload"></a> Enable Identity Management in Workload Clusters

Any workload clusters that you create after you enable identity management in the management cluster are automatically configured to use the same identity management service. If a workload cluster was created before you enabled identity management in your management cluster, you must enable it manually.

To enable identity management in a workload cluster:

1. Generate a Kubernetes secret for the Pinniped add-on.

   1. Create a cluster configuration file using the configuration settings that you defined when you deployed your workload cluster. Include the following variables:

      * Basic cluster variables.

         ```
         # This is the name of your target workload cluster.
         CLUSTER_NAME:
         # For workload clusters, the default namespace is "default".
         NAMESPACE:
         CLUSTER_PLAN:
         CLUSTER_CIDR:
         SERVICE_CIDR:
         ```

      * vSphere-, AWS-, or Azure-specific variables that you set when you deployed your workload cluster. For information about these variables, see [Tanzu CLI Configuration File Variable Reference](../tanzu-config-reference.md).

      * Supervisor issuer URL and CA bundle data.

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

   1. Generate the secret for the Pinniped add-on.

      ```
      tanzu cluster create CLUSTER-NAME --dry-run -f CLUSTER-CONFIG-FILE > CLUSTER-NAME-example-secret.yaml
      ```

1. Review the secret and apply it to the management cluster.

1. After applying the secret, check the status of the Pinniped add-on by running the `kubectl get app` command against the workload cluster.

   ```
   $ kubectl get app pinniped -n tkg-system
   NAME           DESCRIPTION             SINCE-DEPLOY    AGE
   pinniped       Reconcile succeeded     3m23s           7h50m
   ```

1. Configure role-based access control on the workload cluster by following the instructions in [Authenticate Connections to a Workload Cluster](./connect.md#id-mgmt).
