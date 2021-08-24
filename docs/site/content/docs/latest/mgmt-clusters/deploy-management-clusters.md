# Deploying Management Clusters

This topic summarizes how to deploy a Tanzu Kubernetes Grid management cluster or designate one from vSphere with Tanzu.
Deploying or designating a management cluster completes the Tanzu Kubernetes Grid installation process and makes Tanzu Kubernetes Grid operational.

## <a id="overview"></a> Overview

After you have performed the steps described in [Install the Tanzu CLI and Other Tools](../install-cli.md), you can deploy management clusters to the platforms of your choice.

**NOTE**: On vSphere with Tanzu, available on vSphere 7 and later, VMware recommends configuring the built-in supervisor cluster as a management cluster instead of using the `tanzu` CLI to deploy a new management cluster. Deploying a Tanzu Kubernetes Grid management cluster to vSphere 7 when vSphere with Tanzu is not enabled is supported, but the preferred option is to enable vSphere with Tanzu and use the Supervisor Cluster. For details, see [vSphere with Tanzu Provides Management Cluster](vsphere.md#mc-vsphere7).

The management cluster is a Kubernetes cluster that runs Cluster API operations on a specific cloud provider to create and manage workload clusters on that provider.
The management cluster is also where you configure the shared and in-cluster services that the workload clusters use.

## <a id="ui-cli"></a> Installer UI vs. CLI

You can deploy management clusters in two ways:

- Run the Tanzu Kubernetes Grid installer, a wizard interface that guides you through the process of deploying a management cluster. This is the recommended method.
- Create and edit YAML configuration files, and use them to deploy a management cluster with the CLI commands.

## <a id="platforms"></a> Platforms

You can deploy and manage Tanzu Kubernetes Grid management clusters on:

- vSphere 6.7u3
- vSphere 7, if vSphere with Tanzu is not enabled. For more information, see [vSphere with Tanzu Provides Management Cluster](vsphere.md#mc-vsphere7).
- Amazon Elastic Compute Cloud (Amazon EC2)
- Microsoft Azure

You can deploy the management cluster as either a single control plane, for development, or as a highly-available multi-node control plane, for production environments.

- For information about the required setup for management cluster deployment to your infrastructure of choice, see [Prepare to Deploy Management Clusters](prepare-deployment.md).
- When you have set up your infrastructure, see [Deploy Management Clusters with the Installer Interface](deploy-ui.md) or [Deploy Management Clusters from a Configuration File](deploy-cli.md).
- After you have deployed a management cluster to the platform of your choice, [Examine the Management Cluster Deployment](verify-deployment.md).
- If you want to register your management cluster with Tanzu Mission Control, follow the procedure in [Register Your Management Cluster with Tanzu Mission Control](register_tmc.md).
- After you have deployed one or more management clusters to one or more platforms, use the Tanzu CLI to [Manage Your Management Clusters](../cluster-lifecycle/multiple-management-clusters.md).
- If necessary, you can [Back Up and Restore Clusters](../cluster-lifecycle/backup-restore-mgmt-cluster.md).

## <a id="configuring"></a> Configuring the Management Cluster

You deploy your management cluster by running the `tanzu management-cluster create` command on the bootstrap machine.
You configure the management cluster in different ways, depending on whether you specify `--ui` to launch the installer interface:

* **Installer Interface**: UI input
* **CLI**: Set configuration parameters, like `AZURE_NODE_MACHINE_TYPE`:
    - As local environment variables
    - In the cluster configuration file passed to the `--file` option

The `tanzu management-cluster create` command uses these sources and inputs in the following order of increasing precedence:

1. `~/.tanzu/tkg/providers/config_default.yaml`: This file contains system defaults, and should not be changed.
1. **With the `--file` option**: The cluster configuration file, which defaults to `~/.tanzu/tkg/cluster-config.yaml`. This file configures specific invocations of `tanzu management-cluster create` and other `tanzu` commands. Use different `--file` files to save multiple configurations.
1. Local environment variables: Parameter settings in your local environment override settings from config files. Use them to make quick config choices without having to search and edit a config file.
1. **With the `--ui` option**: Installer UI input. When you run `tanzu management-cluster create --ui`, the installer sets all management cluster configuration values from user input and ignores all other CLI options.

## <a id="what-happens"></a> What Happens When You Create a Management Cluster

Running `tanzu management-cluster create` creates a temporary management cluster using a [Kubernetes in Docker](https://kind.sigs.k8s.io/) (`kind`) cluster on the bootstrap machine. After creating the temporary management cluster locally, Tanzu Kubernetes Grid uses it to provision the final management cluster in the platform of your choice.

In the process, `tanzu management-cluster create` creates or modifies CLI configuration and state files in the user's home directory on the local bootstrap machine:

<table id="config-files" width="100%" border="1" class="nice">
  <tr>
    <th width="25%" scope="col">Location</th>
    <th width="50%" scope="col">Content</th>
    <th width="25%" scope="col">Change</th>
  </tr><tr>
    <td><code>~/.tanzu/tkg/bom/</code></td>
    <td>Bill of Materials (BoM) files that list specific versions of all of the packages that Tanzu Kubernetes Grid requires when it creates a cluster with a specific OS and Kubernetes version.  Tanzu Kubernetes Grid adds to this directory as new Tanzu Kubernetes release versions are published.</td>
    <td>Add if not already present</td>
  </tr><tr>
    <td><code>~/.tanzu/tkg/providers/</code></td>
    <td>Configuration template files for Cluster API, cloud providers, and other dependencies, organized with <code>ytt</code> overlays for non-destructive modification.</td>
    <td>Add if not already present</td>
  </tr><tr>
    <td><code>~/.tanzu/tkg/providers-TIMESTAMP-HASH/</code></td>
    <td>Backups of <code>/providers</code> directories from previous installations.</td>
    <td>Add if not first installation</td>
  </tr><tr>
    <td><code>~/.tanzu/config.yaml</code></td>
    <td>Names, contexts, and certificate file locations for the management clusters that the <code>tanzu</code> CLI knows about, and which is the current one.</td>
    <td>Add new management cluster information and set it as <code>current</code>.</td>
  </tr><tr>
    <td><code>~/.tanzu/tkg/cluster-config.yaml</code></td>
    <td>Default cluster configuration file that the <code>tanzu cluster create</code> and <code>tanzu management-cluster create</code> commands use if you do not specify one with <code>--file</code>. <br />
    Best practice is to use a configuration file unique to each cluster.</td>
    <td>Add empty file if not already present.</td>
  </tr><tr>
    <td><code>~/.tanzu/tkg/clusterconfigs/IDENTIFIER.yaml</code></td>
    <td>Cluster configuration file that <code>tanzu management-cluster create --ui</code> writes out with values input from the installer interface.<br />
    <code>IDENTIFIER</code> is an unique identifier generated by the installer.</td>
    <td>Create file</td>
  </tr><tr>
    <td><code>~/.tanzu/tkg/config.yaml</code></td>
    <td>List of configurations and locations for the Tanzu Kubernetes Grid core and all of its providers. </td>
    <td>Add if not already present</td>
  </tr><tr>
    <td><code>~/.tanzu/tkg/providers/config.yaml</code></td>
    <td>Similar to <code>~/.tanzu/tkg/config.yaml</code>, but only lists providers and configurations in the <code>~/.tanzu/tkg/providers</code> directory, not configuration files used by core Tanzu Kubernetes Grid.</td>
    <td>Add if not already present</td>
  </tr><tr>
    <td><code>~/.tanzu/tkg/providers/config_default.yaml</code></td>
    <td>System-wide default configurations for providers.<br />
    Best practice is not to edit this file, but to change provider configs through <code>ytt</code> overlay files.</td>
    <td>Add if not already present</td>
  </tr><tr>
    <td><code>~/.kube-tkg/config</code></td>
    <td>Management cluster <code>kubeconfig</code> file containing names and certificates for the management clusters that the <code>tanzu</code> CLI knows about. Location overridden by the <code>KUBECONFIG</code> environment variable.</td>
    <td>Add new management cluster info and set the cluster as the <code>current-context</code>.</td>
  </tr><tr>
    <td><code>~/.kube/config</code></td>
    <td>Configuration and state for the <code>kubectl</code> CLI, including all management and workload clusters, and which is the current context.</td>
    <td>Add new management cluster name, context, and certificate info. Do not change current <code>kubectl</code> context to new cluster.</td>
  </tr>
</table>

## <a id="core-add-ons"></a> Core Add-ons

When you deploy a management or a workload cluster, Tanzu Kubernetes Grid installs the following core add-ons in the cluster:

* CNI: `cni/calico` or `cni/antrea`
* (**vSphere only**) vSphere CPI: `cloud-provider/vsphere-cpi`
* (**vSphere only**) vSphere CSI: `csi/vsphere-csi`
* Authentication: `authentication/pinniped`
* Metrics Server: `metrics/metrics-server`

Tanzu Kubernetes Grid manages the lifecycle of the core add-ons. For example, it automatically upgrades the add-ons when you upgrade your management and workload clusters.

For more information about the core add-ons, see [Update and Troubleshoot Core Add-On Configuration](../cluster-lifecycle/update-addons.md).
