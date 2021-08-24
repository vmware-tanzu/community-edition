# Tanzu CLI Command Reference

The table below lists all of the commands and options of the Tanzu CLI, and provides links to the section in which they are documented.


<table class="table"> 
 <thead> 
  <tr> 
   <th>Command</th> 
   <th>Options</th> 
   <th>Description</th> 
  </tr> 
 </thead> 
 <tbody> 
  <tr> 
   <td colspan="3"><code>tanzu *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Common Tanzu Kubernetes Grid Options</a></td> 
   <tr> 
   <td colspan="3"><code>tanzu completion *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Output shell completion code for the specified shell</a></td> 
   <tr> 
   <td colspan="3"><code>tanzu cluster *</code></td> 
  </tr>  
  <tr>
   <td></td>
   <td><code>--log-file</code><br /><code>-v,</code>, <code>--verbose</code></td>
   <td> </td>
  </tr>
  <tr>
   <td colspan="3"><code>tanzu cluster create</code></td>
  </tr>
  <tr>
   <td></td>
   <td><code>-d</code>, <code>--dry-run</code></td>
   <td><a href="tanzu-k8s-clusters/deploy.md#manifest">Create Tanzu Kubernetes Cluster Manifest Files</a></td>
  </tr>
  <tr>
   <td></td>
   <td><code>-f</code>, <code>--file</code></td>
   <td><a href="tanzu-k8s-clusters/deploy.md">
   Deploy Tanzu Kubernetes Clusters</a>
   </td>
  </tr>
  <tr>
   <td></td>
   <td><code>--tkr</code></td>
   <td><a href="tanzu-k8s-clusters/k8s-versions.md">Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions</a><br />
   <a href="tanzu-k8s-clusters/connect-vsphere7.md">Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster</a></td>
  </tr>
  <tr>
   <td colspan="3"><code>tanzu cluster credentials update</code></td>
  </tr>
  <tr>
   <td></td>
   <td><code>-n</code>, <code>--namespace</code><br />
       <code>--vsphere-password</code><br />
       <code>--vsphere-user</td>
   <td><a href="cluster-lifecycle/secrets.md">Tanzu Kubernetes Cluster Secrets</a></td>
  </tr>
  <tr>
   <td colspan="3"><code>tanzu cluster delete</code></td>
  </tr>
  <tr>
    <td> </td>
    <td><code>-n</code>, <code>--namespace</code><br />
        <code>-y</code>, <code>--yes</code></td>  
   <td><a href="cluster-lifecycle/delete-cluster.md">Delete Tanzu Kubernetes Clusters</a></td>
  </tr>
  <tr>
   <td colspan="3"><code>tanzu cluster get</code></td>
  </tr>
  <tr>
    <td></td>
    <td><code>--disable-grouping</code><br />
        <code>--disable-no-echo</code><br />
        <code>-n</code>, <code>--namespace</code>
        <code>--show-all-conditions</code></td>
    <td><a href="tanzu-k8s-clusters/index.md">Deploying Tanzu Kubernetes Clusters</a>
   </td>
  </tr>
  <tr>
   <td colspan="2"><code>tanzu cluster kubeconfig get</code></td>
   <td><a href="tanzu-k8s-clusters/index.md">Deploying Tanzu Kubernetes Clusters</a><br />
   <a href="cluster-lifecycle/connect.md">Connect to and Examine Tanzu Kubernetes Clusters</a></td>
  </tr>
  <tr>
   <td></td>
   <td><code>--admin</code><br /><code>--export-file</code><br /><code>-n</code>, <code>--namespace</code></td>
   <td><a href="cluster-lifecycle/connect.md">Connect to and Examine Tanzu Kubernetes Clusters</a><br />
   <a href="upgrade-tkg/workload-clusters.md">Upgrade Tanzu Kubernetes Clusters</a> </td>
  </tr>
  <tr>
   <td colspan="3"><code>tanzu cluster list</code></td>
  </tr>
  <tr>
   <td></td>
   <td><code>--include-management-cluster</code></td>
   <td><a href="cluster-lifecycle/connect.md">Connect to and Examine Tanzu Kubernetes Clusters</a><br /><a href="upgrade-tkg/management-cluster.md">Upgrade Management Clusters</a><br /><a href="upgrade-tkg/workload-clusters.md">Upgrade Tanzu Kubernetes Clusters</a></td>
  </tr>
  <tr>
   <td></td>
   <td><code>-n</code>, <code>--namespace</code><br />
   <code>-o</code>, <code>--output</code></td>
   <td><a href="cluster-lifecycle/connect.md">Connect to and Examine Tanzu Kubernetes Clusters</a></td>
  </tr>
  <tr>
    <td colspan="3"><code>tanzu cluster machinehealthcheck delete</code>
  </tr>
  <tr>
  <td></td>
  <td><code>-m</code>, <code>--mhc-name</code><br />
  <code>-n</code>, <code>--namespace</code><br />
  <code>-y</code>, <code>--yes</code>   </td>
  <td><a href="cluster-lifecycle/configure-health-checks.md">Configure Machine Health Checks for Tanzu Kubernetes Clusters</a></td>
  </tr>
   <tr>
    <td colspan="3">
    <code>tanzu cluster machinehealthcheck get</code></td>
  </tr>
    <tr>
  <td></td>
  <td><code>-m</code>, <code>--mhc-name</code><br />
  <code>-n</code>, <code>--namespace</code><br /></td>
  <td><a href="cluster-lifecycle/configure-health-checks.md">Configure Machine Health Checks for Tanzu Kubernetes Clusters</a></td>
  </tr>
    <tr>
    <td colspan="3">
    <code>tanzu cluster machinehealthcheck set</code></td>
  </tr>
    <tr>
  <td></td>
  <td>
  <code>--match-labels</code><br />
  <code>-m</code>, <code>--mhc-name</code><br />
  <code>-n</code>, <code>--namespace</code><br />
  <code>--node-startup-timeout</code>
  <code>--unhealthy-conditions</code></td>
  <td><a href="cluster-lifecycle/configure-health-checks.md">Configure Machine Health Checks for Tanzu Kubernetes Clusters</a></td>
  </tr>
  <tr>
   <td colspan="3"><code>tanzu cluster scale</code></td>
  </tr>
  <tr>
   <td></td>
   <td><code>-c</code>, <code>--controlplane-machine-count</code><br />
       <code>-n</code>, <code>--namespace</code><br />
       <code>-w</code>, <code>--worker-machine-count</code></td>
   <td><a href="cluster-lifecycle/multiple-management-clusters.md">Manage Your Management Clusters</a><br /><a href="cluster-lifecycle/scale-cluster.md">Scale Tanzu Kubernetes Clusters</a></td>
  </tr>
  <tr>
   <td colspan="3"><code>tanzu cluster upgrade</code></td>
  </tr>
  <tr>
   <td></td>
   <td><code>-n</code>, <code>--namespace</code><br />
   <code>--os-arch</code><br />
   <code>--os-name</code><br />
   <code>--os-version</code><br />
   <code>-t</code>, <code>--timeout</code><br />
   <code>--tkr</code><br />
   <code>-y</code>, <code>--yes</code></td>
   <td><a href="upgrade-tkg/workload-clusters.md">Upgrade Tanzu Kubernetes Clusters</a><br />
   <a href="tanzu-k8s-clusters/k8s-versions.md">Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions</a><br /></td>
  </tr>
    <tr> 
   <td colspan="3"><code>tanzu config init *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Initializes the configuration with defaults
</a></td> 
   <tr> 
  <tr>
    <td colspan="3"><code>tanzu config server delete</code></td>
  </tr>
  <tr>
   <td></td>
   <td><code>-y</code>, <code>--yes</code></td>
   <td><a href="cluster-lifecycle/multiple-management-clusters.md#delete-mc-config">Delete Management Clusters from Your Tanzu CLI Configuration</a></td>
  </tr>
  <tr>
   <td colspan="2"><code>tanzu config server list</code></td>
   <td><a href="cluster-lifecycle/multiple-management-clusters.md#delete-mc-config">Delete Management Clusters from Your Tanzu CLI Configuration</a></td>
  </tr>
     <tr> 
   <td colspan="3"><code>tanzu config show *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Shows the current configuration
</a></td> 
   <tr> 
  <tr>
    <td colspan="2"><code>tanzu init</code></td>
    <td> Not available in this version of the Tanzu CLI</td>
  </tr>
  <tr>
   <td colspan="2"><code>tanzu kubernetes-release get</code><br />
   <code>tanzu kubernetes-release available-upgrades get</code><br /></td>
   <td><a href="tanzu-k8s-clusters/k8s-versions.md">Deploy Tanzu Kubernetes Clusters with Different Kubernetes Versions</a><br />
   <a href="upgrade-tkg/workload-clusters.md">Upgrade Tanzu Kubernetes Clusters</a><br />
  <a href="tanzu-k8s-clusters/connect-vsphere7.md">Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster</a>
   </td>
  </tr>
   <tr>
   <td colspan="3">
   <code>tanzu kubernetes-release os get</code></td>
  </tr>
  <tr>
  <td></td>
  <td><code>--region</code>
  <td></td>
  </tr>

  <tr>
   <td colspan="3"><code>tanzu login</code></td>
  </tr>
  <tr>
  <td> </td>
  <td>
    <code>--apiToken</code><br />
    <code>--context</code><br />
    <code>--endpoint</code><br />
    <code>--kubeconfig</code><br />
    <code>--name</code><br />
    <code>--server</code>
    </td>
     <td><a href="cluster-lifecycle/connect.md">Connect to and Examine Tanzu Kubernetes Clusters</a><br />
   <a href="cluster-lifecycle/multiple-management-clusters.md">Manage Your Management Clusters</a><br />
   <a href="tanzu-k8s-clusters/connect-vsphere7.md">Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster</a><br /></td>
  </tr>
  <tr>
    <td colspan="2"><code>tanzu management-cluster ceip-participation get</code><br />
    <code>tanzu management-cluster ceip-participation set</code></td>
    <td><a href="cluster-lifecycle/multiple-management-clusters.md#ceip">Opt in or Out of the VMware CEIP</a></td>
  </tr>
 <tr>
   <td colspan="3"><code>tanzu management-cluster create</code></td>
  </tr>
 <tr>
   <td> </td>
   <td><code>-b</code>, <code>--bind</code><br />
      <code>--browser</code><br />
      <code>-u</code>, <code>--ui</code>
   </td>
   <td>
     <a href="mgmt-clusters/deploy-ui.md">Deploy Management Clusters with the Installer Interface</a></td>
  </tr>
 <tr>
   <td> </td>
   <td><code>-f</code>, <code>--file</code><br />
   <code>-t</code>, <code>--timeout</code><br />
   <code>-y</code>, <code>--yes</code><br /></td>
   <td><a href="mgmt-clusters/deploy-cli.md">Deploy Management Clusters from a Configuration File</a></td>
 </tr>
  <tr>
   <td> </td>
   <td><code>-e</code>, <code>--use-existing-bootstrap-cluster</code><br /></td>
   <td><a href="troubleshooting-tkg/use-existing.md">Use an Existing Bootstrap Cluster to Deploy Management Clusters</a></td>
 </tr>  
  <tr>
   <td colspan="2"><code>tanzu management-cluster credentials update</code></td>
   <td><a href="cluster-lifecycle/secrets.md">Tanzu Kubernetes Cluster Secrets</a></td>
  </tr>
  <tr>
   <td colspan="2"><code>tanzu management-cluster delete</code></td>
   <td><a href="cluster-lifecycle/multiple-management-clusters.md#delete">Delete Management Clusters</a></td>
  </tr>
  <tr>
    <td colspan="2"><code>tanzu management-cluster get</code></td>
    <td><a href="cluster-lifecycle/connect.md">Connect to and Examine Tanzu Kubernetes Clusters</a><br />
      <a href="mgmt-clusters/verify-deployment.md">Examine the Management Cluster Deployment</a><br />
      <a href="cluster-lifecycle/multiple-management-clusters.md">Manage Your Management Clusters</a><br />
      <a href="upgrade-tkg/management-cluster.md">Upgrade Management Clusters</a>
      </p></td>
  </tr>
    <tr>
    <td colspan="2"><code>tanzu management-cluster import</code></td>
    <td><a href="upgrade-tkg/index.md">Upgrading Tanzu Kubernetes Grid</a>
    </td>
  </tr>
   <tr>
    <td colspan="2"><code>tanzu management-cluster kubeconfig get</code></td>
    <td>
    </td>
  </tr>
  <tr>
   <td></td>
   <td><code>--admin</code><br /><code>--export-file</code></td>
   <td><a href="mgmt-clusters/verify-deployment.md">Examine the Management Cluster Deployment</a><br />
   <a href="mgmt-clusters/configure-id-mgmt.md">Configure Identity Management After Management Cluster Deployment</a></td>
  </tr>
  <tr>
    <td colspan="2"><code>tanzu management-cluster register</code></td>
    <td>
      <a href="mgmt-clusters/register_tmc.md">Register the Management Cluster with Tanzu Mission Control</a><br />
    </td>
  </tr>
  <tr>
   <td colspan="2"><code>tanzu management-cluster permissions aws get</code><br />
   <code>tanzu management-cluster permissions aws set</code></td>
   <td><a href="mgmt-clusters/config-aws.md">Create a Cluster Configuration File for Amazon EC2</a></td>
  </tr>
  <tr>
   <td colspan="3"><code>tanzu management-cluster upgrade</code></td>
  </tr>
  <tr>
   <td></td>
   <td><code>--os-arch</code><br />
   <code>--os-name</code><br />
   <code>--os-version</code><br />
   <code>-t</code>, <code>--timeout</code><br />
   <code>-y</code>, <code>--yes</code></td>
   <td><a href="upgrade-tkg/management-cluster.md">Upgrade Management Clusters</a></td>
  </tr>
     <tr> 
   <td colspan="3"><code>tanzu pinniped-auth *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Pinniped authentication operations
</a></td> 
   <tr> 
       <tr> 
   <td colspan="3"><code>tanzu pinniped-auth login *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code>, <code>--ca-bundle strings</code>, <code>--client-id string</code>, <code>--concierge-authenticator-name string</code>, <code>--concierge-authenticator-type</code>, <code>--concierge-ca-bundle-data string</code>, <code>--concierge-authenticator-name string</code>, <code>--concierge-authenticator-type</code>, <code>--concierge-ca-bundle-data string</code>, <code>--concierge-endpoint string</code>, <code>--concierge-namespace string</code>, <code>--enable-concierge</code>, <code>--issuer string</code>, <code>--listen-port uint16</code>, <code>--request-audience string</code>, <code>--scopes strings</code>, <code>--session-cache string</code>, <code>--skip-browser</code><br /></td> 
   <td><a href="install-cli.md#common-options">Log in using an OpenID Connect provider
</a></td> 
   <tr> 
  <tr> 
   <td colspan="2">
   <code>tanzu plugin clean</code> <sup>*</sup><br />
   <code>*tanzu plugin install*</code> <sup>*</sup><br />
   <code>tanzu plugin list</code> <sup>*</sup>
 </td>
   <td><a href="install-cli.md">Install the Tanzu CLI</a></td>
  </tr> 
     <tr> 
   <td colspan="3"><code>tanzu plugin delete *</code></td> 
  </tr>
 </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Deletes a Tanzu plugin  
   <tr> 
   <td colspan="3"><code>tanzu plugin describe *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Describes a Tanzu plugin
</a></td> 
   <tr>   
   <tr> 
</a></td> 
   <tr> 
   <td colspan="3"><code>tanzu plugin upgrade *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Upgrades a Tanzu plugin
</a></td> 
   <tr> 
   <tr> 
   <td colspan="3"><code>tanzu update *</code></td> 
  </tr>
   <tr> 
   <td></td> 
   <td><code>-h</code>, <code>--help</code><br /></td> 
   <td><a href="install-cli.md#common-options">Updates the Tanzu CLI
</a></td> 
   <tr> 
   <tr> 
   <td colspan="2"><code>tanzu version</code></td>
   <td><a href="install-cli.md">Install the Tanzu CLI</a></td>
  </tr>
 </tbody>
</table>

<sup>*</sup> Some <code>tanzu plugin</code> commands such as <code>tanzu plugin repo</code> and <code>tanzu plugin update</code> are not functional in the current release.

## <a id="equivalents"></a> Table of Equivalents

The `tanzu` command-line interface (CLI) works differently from the `tkg` CLI used in previous versions of Tanzu Kubernetes Grid.
Because of these differences, many `tkg` commands do not have direct `tanzu` equivalents, and vice-versa.
But other commands do have direct or close equivalencies across both CLIs:

<table class="table">
 <thead>
  <tr>
   <th>TKG CLI command</th>
   <th>Tanzu CLI command</th>
   <th>Notes</th>
  </tr>
 </thead>
 <tbody>
   <tr>
    <td><code>tkg init</code></td>
    <td><code>tanzu management-cluster create</code></td>
    <td>Some <code>tkg</code> flags not replicated in <code>tanzu</code></td>
  </tr>
   <tr>
    <td><code>tkg init --ui</code></td>
    <td><code>tanzu management-cluster create --ui</code></td>
    <td>Some <code>tkg</code> flags not replicated in <code>tanzu</code></td>
  </tr>
  <tr>
    <td><code>tkg set management-cluster MANAGEMENT_CLUSTER_NAME</code></td>
    <td><code>tanzu login --server SERVER</code>
  </tr>
  <tr>
    <td><code>tkg add management-cluster</code></td>
    <td><code>tanzu login --kubeconfig CONFIG --context OPTIONAL_CONTEXT --name SERVER-NAME</code></td>
    <td> </td>
  </tr>
  <tr>
    <td><code>tkg get management-cluster</code></td>
    <td><code>tanzu login</code></td>
    <td>The <code>tanzu login</code> command displays the currently configured management clusters</td>
  </tr>
  <tr>
    <td><code>tkg create cluster</code></td>
    <td><code>tanzu cluster create</code></td>
    <td>Some <code>tkg</code> flags not replicated in <code>tanzu</code></td>
  </tr>
  <tr>
    <td><code>tkg config cluster</code></td>
    <td><code>tanzu cluster create --dry-run</code></td>
    <td>Some <code>tkg</code> flags not replicated in <code>tanzu</code></td>
  </tr>
  <tr>
    <td><code>tkg delete cluster</code></td>
    <td><code>tanzu cluster delete</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg get cluster</code></td>
    <td><code>tanzu cluster list</code></td>
    <td></td>
  </tr>
    <tr>
    <td><code>tkg scale cluster</code></td>
    <td><code>tanzu cluster scale</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg upgrade cluster</code></td>
    <td><code>tanzu cluster upgrade</code></td>
    <td></td>
  </tr>
   <tr>
    <td><code>tkg get credentials</code></td>
    <td><code>tanzu management-cluster kubeconfig get --admin</code></td>
    <td></td>
  </tr>
     <tr>
    <td><code>tkg get credentials WORKLOAD_CLUSTER_NAME</code></td>
    <td><code>tanzu cluster kubeconfig get WORKLOAD_CLUSTER_NAME --admin</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg delete machinehealthcheck</code></td>
    <td><code>tanzu cluster machinehealthcheck delete</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg get machinehealthcheck</code></td>
    <td><code>tanzu cluster machinehealthcheck get</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg set machinehealthcheck</code></td>
    <td><code>tanzu cluster machinehealthcheck set</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg get ceip-participation</code></td>
    <td><code>tanzu management-cluster ceip-participation get</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg set ceip-participation</code></td>
    <td><code>tanzu management-cluster ceip-participation set</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg delete management-cluster</code></td>
    <td><code>tanzu management-cluster delete</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg upgrade management-cluster</code></td>
    <td><code>tanzu management-cluster upgrade</code></td>
    <td></td>
  </tr>
  <tr>
    <td><code>tkg version</code></td>
    <td><code>tanzu version</code></td>
    <td></td>
  </tr>
   <tr>
    <td><code>tkg get kubernetesversions</code></td>
    <td><code>tanzu kubernetes-release get</code></td>
    <td></td>
  </tr>
 </tbody>
</table>
