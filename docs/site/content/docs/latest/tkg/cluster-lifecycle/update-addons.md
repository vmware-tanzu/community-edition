
# Update and Troubleshoot Core Add-On Configuration

This topic describes how to update and troubleshoot the default configuration of core add-ons in Tanzu Kubernetes Grid.

## <a id="default-configuration"></a> Default Core Add-On Configuration

Tanzu Kubernetes Grid automatically manages the lifecycle of its core add-ons, which includes the CNI, Metrics Server, Pinniped, vSphere CPI, and vSphere CSI add-ons. For more information, see [Core Add-ons](../mgmt-clusters/deploy-management-clusters.md#core-add-ons) in _Deploying Management Clusters_.

To review the default configuration of these add-ons, you can:

* Download the following templates from `projects.registry.vmware.com/tkg/tanzu_core/addons`:

   * `antrea-templates`
   * `calico-templates`
   * `metrics-server-templates`
   * `pinniped-templates`
   * `vsphere-cpi-templates`
   * `vsphere-csi-templates`

* Examine the Kubernetes secret for your target add-on by running the `kubectl get secret CLUSTER-NAME-ADD-ON-NAME-addon -n CLUSTER-NAMESPACE` command against the management cluster.

For example, to review the default configuration of the Antrea add-on:

   * Review the Antrea templates:

      1. Locate the version tag for `antrea-templates` in the Tanzu Kubernetes release (TKr) that you used to deploy your cluster. You can retrieve the TKr by running the `kubectl get tkr` command against the management cluster:

         1. Run `kubectl get clusters CLUSTER-NAME -n CLUSTER-NAMESPACE  --show-labels`.

         1. In the output, locate the value of `tanzuKubernetesRelease`. For example, `tanzuKubernetesRelease=v1.20.5---vmware.1-tkg.1`.

         1. Run `kubectl get tkr TKR-VERSION`, where `TKR-VERSION` is the value that you retrieved above. For example:

            ```
            kubectl get tkr v1.20.5---vmware.1-tkg.1 -o yaml
            ```

         1. In the output, locate the version tag under `tanzu_core/addons/antrea-templates`.

          Alternatively, you can review the TKr in `~/tanzu/tkg/bom/YOUR-TKR-BOM-FILE`.

      1. Download the templates. For example:

         ```
         imgpkg pull -i projects.registry.vmware.com/tkg/tanzu_core/addons/antrea-templates:v1.3.1 -o antrea-templates
         ```

      1. Navigate to `antrea-templates` and review the templates.

   * Retrieve and review the Antrea add-on secret. To retrieve the secret, run the following command against the management cluster:

      ```
      kubectl get secret CLUSTER-NAME-antrea-addon -n CLUSTER-NAMESPACE
      ```

      Where:

      * `CLUSTER-NAME` is the name of your target cluster. If you want to review the Antrea add-on secret for a workload cluster, `CLUSTER-NAME` is the name of your workload cluster.
      * `CLUSTER-NAMESPACE` is the namespace of your target cluster.

## <a id="resources"></a> Updating and Troubleshooting Core Add-on Configuration

You can update and troubleshoot the default configuration of a core add-on by modifying the following resources:

<table class="table">
  <thead>
    <tr>
      <th>Type</th>
      <th>Resources</th>
      <th width="60%">Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Configuration updates</td>
      <td>Add-on secret</td>
      <td><p>To update the default configuration of a core add-on, you can:</p>
         <ul>
           <li>Modify the <code>values.yaml</code> section of the add-on secret. For more information, see <a href="#values-yaml">Update the values.yaml section</a> below.</li>
           <li>Add an overlay to the add-on secret. For more information, see <a href="#overlay">Add an Overlay</a> below.</li>
         </ul>
      </td>
    </tr>
    <tr>
      <td>Troubleshooting</td>
      <td>App custom resource (CR) and add-on secret</td>
      <td><p>Same as above. Additionally, if you need to apply temporary changes to your add-on configuration, you can:</p>
        <ul>
          <li>Pause secret reconciliation.</li>
          <li>Pause app CR reconciliation.</li>
        </ul>
        <p><strong>This disables lifecycle management for the add-on. Use with caution.</strong> For more information, see <a href="#pause">Pause Core Add-on Lifecycle Management</a> below.</p>
      </td>
    </tr>
  </tbody>
</table>

For more information about add-on secrets and app CRs, see [Key Components and Objects](#components) below.

### <a id="updating"></a> Updating Core Add-on Configuration

You can update the default configuration of a core add-on by modifying the `values.yaml` section of the add-on secret or by adding an overlay to the add-on secret. These changes are persistent.

#### <a id="values-yaml"></a> Update the values.yaml section

In the `values.yaml` section, you can update the following configuration settings:

<table class="table">
 <thead>
    <tr>
      <th>Add-on</th>
      <th width="40%">Setting</th>
      <th width="45%">Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Antrea</td>
      <td><code>antrea.config.defaultMTU</code></td>
      <td>By default, this parameter is set to <code>null</code>.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.oidc.CLIENT_ID</code><sup>&#42;</sup> (v1.3.0) or <code>pinniped.upstream_oidc_client_id</code> (v1.3.1+)</td>
      <td>The client ID of your OIDC provider.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.oidc.CLIENT_SECRET</code> (v1.3.0) or <code>pinniped.upstream_oidc_client_secret</code> (v1.3.1+)</td>
      <td>The client secret of your OIDC provider.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.oidc.issuer</code> (v1.3.0) or <code>pinniped.upstream_oidc_issuer_url</code> (v1.3.1+)</td>
      <td>The URL of your OIDC provider.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.oidc.scopes</code> (v1.3.0) or <code>pinniped.upstream_oidc_additional_scopes</code> (v1.3.1+)</td>
      <td>A list of additional scopes to request in the token response.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.oidc.claimMapping</code> (v1.3.0) or <code>pinniped.upstream_oidc_claims</code> (v1.3.1+)</td>
      <td>OIDC claim mapping.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.ldap.host</code></td>
      <td>The IP or DNS address of your LDAP server. If you want to change the default port 636 to a different port, specify <code>"host:port"</code>.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.ldap.bindDN</code> and <code>dex.config.ldap.bindPW</code></td>
      <td>The DN and password for an application service account.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.ldap.userSearch</code></td>
      <td>Search attributes for users.</td>
    </tr>
    <tr>
      <td>Pinniped</td>
      <td><code>dex.config.ldap.groupSearch</code></td>
      <td>Search attributes for groups.</td>
    </tr>
    <tr>
      <td>vSphere CSI</td>
      <td><code>vsphereCSI.provisionTimeout</code></td>
      <td>By default, this parameter is set to <code>300s</code>.</td>
    </tr>
    <tr>
      <td>vSphere CSI</td>
      <td><code>vsphereCSI.attachTimeout</code></td>
      <td>By default, this parameter is set to <code>300s</code>.</td>
    </tr>
 </tbody>
</table>

<sup>&#42;</sup> If you want to update a Pinniped setting that starts with `dex.`, you must restart `dex` in the management cluster after you update the add-on secret.

To modify the `values.yaml` section of an add-on secret:

1. Retrieve the add-on secret by running the `kubectl get secret CLUSTER-NAME-ADD-ON-NAME-addon -n CLUSTER-NAMESPACE` command against the management cluster. For example:

   ```
   kubectl get secret example-mgmt-cluster-antrea-addon -n tkg-system -o jsonpath={.data.values\\.yaml} | base64 -D > values.yaml
   ```

1. Update the `values.yaml` section. You can update any of the values listed in the table above.

1. Apply your update by running the `kubectl apply` command. Alternatively, you can use the `kubectl edit` command to update the add-on secret.

1. After updating the secret, check the status of the add-on by running the `kubectl get app` command. For example:

   ```
   $ kubectl get app antrea -n tkg-system
   NAME           DESCRIPTION             SINCE-DEPLOY    AGE
   antrea         Reconcile succeeded     3m23s           7h50m
   ```

   If the returned status is `Reconcile failed`, run the following command to get details on the failure:

   ```
   kubectl get app antrea -n tkg-system -o yaml
   ```

The example below updates the default MTU for the Antrea add-on.

```
...
stringData:
  values.yaml: |
    #@data/values
    #@overlay/match-child-defaults missing_ok=True
    ---
    infraProvider: vsphere
    antrea:
      config:
        defaultMTU: 8900
```

#### <a id="overlay"></a> Add an Overlay

If you want to update a configuration setting that is not supported by the default add-on templates, you can add an overlay to the add-on secret. The example below instructs Pinniped to use `LoadBalancer` instead of the default `NodePort` on vSphere:

```
...
stringData:
 overlays.yaml: |
   #@ load("@ytt:overlay", "overlay")
   #@overlay/match by=overlay.subset({"kind": "Service", "metadata": {"name": "pinniped-supervisor", "namespace": "pinniped-supervisor"}})
   ---
   #@overlay/replace
   spec:
     type: LoadBalancer
     selector:
       app: pinniped-supervisor
     ports:
       - name: https
         protocol: TCP
         port: 443
         targetPort: 8443
 values.yaml: |
   #@data/values
   #@overlay/match-child-defaults missing_ok=True
   ---
   infrastructure_provider: vsphere
   tkg_cluster_role: management
```

To add an overlay to an add-on secret:

1. Retrieve the add-on secret by running the `kubectl get secret CLUSTER-NAME-ADD-ON-NAME-addon -n CLUSTER-NAMESPACE` command against the management cluster. For example:

   ```
   kubectl get secret example-mgmt-cluster-pinniped-addon -n tkg-system -o jsonpath={.data.values\\.yaml} | base64 -D > values.yaml
   ```

1. Add your `overlay.yaml` section under `stringData`.

1. Apply the update by running the `kubectl apply` command. Alternatively, you can use the `kubectl edit` command to update the add-on secret.

1. After updating the secret, check the status of the add-on by running the `kubectl get app` command. For example:

   ```
   $ kubectl get app pinniped -n tkg-system
   NAME           DESCRIPTION             SINCE-DEPLOY    AGE
   pinniped       Reconcile succeeded     3m23s           7h50m
   ```

   If the returned status is `Reconcile failed`, run the following command to get details on the failure:

   ```
   kubectl get app pinniped -n tkg-system -o yaml
   ```

### <a id="troubleshooting"></a> Troubleshooting Core Add-on Configuration

Before troubleshooting the core add-ons, review the following sections:

* [Key Components and Objects](#components) below.
* [Updating Core Add-on Configuration](#updating) above.
* [Pause Core Add-on Lifecycle Management](#pause) below.

#### <a id="components"></a> Key Components and Objects

Tanzu Kubernetes Grid uses the following components and objects for core add-on management.

**Components in the management cluster:**

* `kapp-controller`, a local package manager: When you deploy a management cluster, the Tanzu CLI installs `kapp-controller` in the cluster. `kapp-controller` deploys `tanzu-addons-manager` and the core add-ons. It also deploys and manages `kapp-controller` in each Tanzu Kubernetes (workload) cluster that you deploy from that management cluster.
* `tanzu-addons-manager`: Manages the lifecycle of the core add-ons in the management cluster and workload clusters that you deploy from your management cluster.
* `tkr-controller`: Creates Tanzu Kubernetes releases (TKr) and BoM ConfigMaps in the management cluster.

**Component in workload clusters:**

`kapp-controller` deploys the core add-ons in the workload cluster in which it runs.

**Objects:**

* **Secret:** The Tanzu CLI creates a secret for each core add-on, per cluster. These secrets define the configuration of the core add-ons. All add-on secrets are created in the management cluster. `tanzu-addons-manager` reads the secrets and uses the configuration information they contain to create app CRs.
* **App CR:** For each add-on, `tanzu-addons-manager` creates an app CR in the target cluster. Then, `kapp-controller` reconciles the CR and deploys the add-on.
* **BoM ConfigMap:** Provides metadata information about the core add-ons, such as image location, to `tanzu-addons-manager`.

You can use the following commands to monitor the status of these components and objects:

<table class="table">
  <thead>
    <tr>
     <th>Command</th>
     <th width="50%">Description</th>
    </tr>
  </thead>
  <tbody>
  <tr>
   <td><code>kubectl get app ADD-ON -n tkg-system -o yaml</code></td>
   <td>Check the app CR in your target cluster. For example, <code>kubectl get app antrea -n tkg-system -o yaml</code>.</td>
  </tr>
  <tr>
   <td><code>kubectl get cluster CLUSTER-NAME -n CLUSTER-NAMESPACE -o jsonpath={.metadata.labels.tanzuKubernetesRelease}</code></td>
   <td>In the management cluster, check if the TKr label of your target cluster points to the correct TKr.</td>
  </tr>
  <tr>
   <td><code>kubectl get tanzukubernetesrelease TKR-NAME</code></td>
   <td>Check if the TKr is present in the management cluster.</td>
  </tr>
  <tr>
   <td><code>kubectl get configmaps -n tkr-system -l 'tanzuKubernetesRelease=TKR-NAME'</code></td>
   <td>Check if the BoM ConfigMap corresponding to your TKr is present in the management cluster.</td>
  </tr>
  <tr>
   <td><code>kubectl get app CLUSTER-NAME-kapp-controller -n CLUSTER-NAMESPACE</code></td>
   <td>For workload clusters, check if the <code>kapp-controller</code> app CR is present in the management cluster.</td>
  </tr>
  <tr>
   <td><code>kubectl logs deployment/tanzu-addons-controller-manager -n tkg-system</code></td>
   <td>Check <code>tanzu-addons-manager</code> logs in the management cluster.</td>
  </tr>
  <tr>
   <td><code>kubectl get configmap -n tkg-system | grep ADD-ON-ctrl</code></td>
   <td>Check if your updates to the add-on secret have been applied. The sync period is 5 minutes.</td>
  </tr>
  </tbody>
  </table>

#### <a id="pause"></a> Pause Core Add-on Lifecycle Management

**IMPORTANT:** The commands in this section disable add-on lifecycle management. Whenever possible, use the procedures described in [Updating Add-on Configuration](#updating)  above instead.

If you need to temporary pause add-on lifecycle management for a core add-on, you can use the commands below:

   * To pause secret reconciliation, run the following command against the management cluster:

      ```
      kubectl patch secret/CLUSTER-NAME-ADD-ON-NAME-addon -n CLUSTER-NAMESPACE -p '{"metadata":{"annotations":{"tkg.tanzu.vmware.com/addon-paused": ""}}}' --type=merge
      ```

      After you run this command, `tanzu-addons-manager` stops reconciling the secret.

   * To pause app CR reconciliation, run the following command against your target cluster:

      ```
      kubectl patch app/ADD-ON-NAME -n tkg-system -p '{"spec":{"paused":true}}' --type=merge
      ```

      After you run this command, `kapp-controller` stops reconciling the app CR.

If you want to temporary modify a core add-on app, pause secret reconciliation first and then pause app CR reconciliation. After you unpause add-on lifecycle management, `tanzu-addons-manager` and `kapp-controller` resume secret and app CR reconciliation:

* To unpause secret reconciliation, remove `tkg.tanzu.vmware.com/addon-paused` from the secret annotations.

* To unpause app CR reconciliation, update the app CR with `{"spec":{"paused":false}}` or remove the variable.
