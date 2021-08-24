# Upgrading Tanzu Kubernetes Grid

To upgrade Tanzu Kubernetes Grid, you download and install the new version of the Tanzu CLI on the machine that you use as the bootstrap machine. You must also download and install base image templates and VMs, depending on whether you are upgrading clusters that you previously deployed to vSphere, Amazon EC2, or Azure.

After you have installed the new versions of the components, you use the `tanzu management-cluster upgrade` and `tanzu cluster upgrade` CLI commands to upgrade management clusters and Tanzu Kubernetes clusters.

## <a id="prereq"></a> Prerequisites

Before you begin the upgrade to Tanzu Kubernetes Grid v1.3.1, you must ensure the following prequisites are met.

* Your current deployment is Tanzu Kubernetes Grid v1.2.x or v1.3.0. To upgrade to v1.3.x from Tanzu Kubernetes Grid versions earlier than v1.2, you must first upgrade to v1.2.x with the `tkg` v1.2.x CLI.
* If your deployment is running on vSphere, you have migrated your clusters from an HA Proxy Load Balancer to Kube-VIP. You must complete this migration before upgrading to Tanzu Kubernetes Grid v1.3.x. For migration instructions, see [Migrate Clusters from an HA Proxy Load Balancer to Kube-VIP](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-upgrade-tkg-migrate-haproxy.html) in the v1.2 documentation.

## <a id="process"></a> Procedure

The next sections are the overall steps required to upgrade Tanzu Kubernetes Grid. This procedure assumes that you are upgrading to Tanzu Kubernetes Grid v1.3.1.

Some steps are only required if you are performing a major upgrade from Tanzu Kubernetes Grid v1.2.x to v1.3.1 and are not required if you are performing a minor upgrade from Tanzu Kubernetes Grid v1.3.0 to v1.3.1.

If you deployed the previous version of Tanzu Kubernetes Grid on vSphere in an Internet-restricted environment, see [Upgrading vSphere Deployments in an Internet Restricted Environment](internet-restrict-upgrade.md).

## <a id="download-cli"></a> Download and Install the Tanzu CLI

This step is required for both major v1.2.x to v1.3.1 and minor v1.3.0 to v1.3.1 upgrades.

To download the Tanzu CLI, perform the following steps.

1. Follow the instructions in [Install the Tanzu CLI and Other Tools](../install-cli.md) to download and install the Tanzu CLI and `kubectl` on the machine where you currently run your `tkg` commands for v1.2.x or `tanzu` commands for v1.3.0.  
1. After you install `tanzu`, run `tanzu version` to check that the correct version of the Tanzu CLI is properly installed.
1. After you install `kubectl`, run `kubectl version` to check that the correct version of `kubectl` is properly installed.

For information about Tanzu CLI commands and options that are available, see the [Tanzu CLI Command Reference](../tanzu-cli-reference.md).

## <a id="import"></a> Import Configuration Files from Existing v1.2 Management Clusters

This step is only required for major upgrades from Tanzu Kubernetes Grid v1.2.x to v1.3.x.

On the machine where you currently run `tkg` commands to manage your Tanzu Kubernetes Grid v1.2 clusters, import your existing management cluster configurations into the configuration file used by the `tanzu` CLI.

If you store your management cluster configuration files in the default location, `~/.tkg/config.yaml`,  run the following command:

```
tanzu management-cluster import
```

If the configuration files for your management clusters are stored in a non-default location or if you have multiple management cluster config.yaml files, specify the location of the file in the command.

```
tanzu management-cluster import -f /path/to/config.yaml
```

If the import is successful, the command returns a message similar to the following:

```
the old providers folder /Users/username/.tkg/providers is backed up to /Users/username/.tkg/providers-20210309174325-0aw0ckqa
successfully imported server: tkg-cluster-mc

Management cluster configuration imported successfully
```

By importing your existing configuration files, your Tanzu Kubernetes Grid v1.2 management clusters are now accessible by the Tanzu CLI.

In `~/.tanzu/config.yaml`, you should see new entries for the management clusters you imported in the `servers` section. For example:

```
servers:
- managementClusterOpts:
    context: tkg-cluster-mc-admin@tkg-cluster-mc
    path: /Users/username/.kube-tkg/tkg-cluster-mc_config
  name: tkg-cluster-mc
  type: managementcluster
```

## <a id="update-scripts"></a> Replace v1.2 <code>tkg</code> Commands with <code>tanzu</code> Commands

This step is only required for major upgrades from Tanzu Kubernetes Grid v1.2.x to v1.3.x.

If you use any `tkg` commands in automation scripts for your Tanzu Kubernetes clusters, make sure to replace any `tkg` command invocations with equivlent `tanzu` CLI commands before you upgrade to v1.3.x.

Tanzu Kubernetes Grid v1.3.x clusters cannot be managed with the `tkg` CLI.

For a reference of command equivalents, see [Table of Equivalents](../tanzu-cli-reference.md#equivalents) in the _Tanzu CLI Command Reference_.

## <a id="vsphere"></a> Prepare to Upgrade Clusters on vSphere

This step is required for both major v1.2.x to v1.3.1 and minor v1.3.0 to v1.3.1 upgrades.

Before you can upgrade a Tanzu Kubernetes Grid deployment on vSphere, you must import into vSphere new versions of the base image templates that the upgraded management and Tanzu Kubernetes clusters will run.
VMware publishes base image templates in OVA format for each supported OS and Kubernetes version.
After importing the OVAs, you must convert the resulting VMs into VM templates.

This procedure assumes that you are upgrading to Tanzu Kubernetes Grid v1.3.1.

1. Go to [the Tanzu Kubernetes Grid downloads page](https://my.vmware.com/en/web/vmware/downloads/info/slug/infrastructure_operations_management/vmware_tanzu_kubernetes_grid/1_x) and log in with your My VMware credentials.
1. Download the latest Tanzu Kubernetes Grid OVAs for the OS and Kubernetes version lines that your management and Tanzu Kubernetes clusters are running.

   For example, for Photon v3 images:

   - Kubernetes v1.20.5: **Photon v3 Kubernetes v1.20.5 OVA**
   - Kubernetes v1.19.9: **Photon v3 Kubernetes v1.19.9 OVA**
   - Kubernetes v1.18.17: **Photon v3 Kubernetes v1.18.17 OVA**

   For Ubuntu 20.04 images:

   - Kubernetes v1.20.5: **Ubuntu 2004 Kubernetes v1.20.5 OVA**
   - Kubernetes v1.19.9: **Ubuntu 2004 Kubernetes v1.19.9 OVA**
   - Kubernetes v1.18.17: **Ubuntu 2004 Kubernetes v1.18.17 OVA**

   **Important**: Make sure you download the most recent OVA base image templates in the event of security patch releases.
   You can find updated base image templates that include security patches on the Tanzu Kubernetes Grid product download page.

1. In the vSphere Client, right-click an object in the vCenter Server inventory and select **Deploy OVF template**.
1. Select **Local file**, click the button to upload files, and navigate to a downloaded OVA file on your local machine.
1. Follow the installer prompts to deploy a VM from the OVA.

    - Accept or modify the appliance name.
    - Select the destination datacenter or folder.
    - Select the destination host, cluster, or resource pool.
    - Accept the end user license agreements (EULA).
    - Select the disk format and destination datastore.
    - Select the network for the VM to connect to.
1. Click **Finish** to deploy the VM.
1. When the OVA deployment finishes, right-click the VM and select **Template** > **Convert to Template**.  
1. In the **VMs and Templates** view, right-click the new template, select **Add Permission**, and assign your Tanzu Kubernetes Grid user, for example, `tkg-user`, to the template with the Tanzu Kubernetes Grid role, for example, `TKG`. You created this user and role in [Prepare to Deploy Management Clusters to vSphere](../mgmt-clusters/vsphere.md#vsphere-permissions).

Repeat the procedure for each of the Kubernetes versions for which you have downloaded the OVA file.

### VMware Cloud on AWS SDDC Compatibility

If you are upgrading Tanzu Kubernetes clusters that are deployed on VMware Cloud on AWS, verify that the underlying Software-Defined Datacenter (SDDC) version used by your existing deployment is compatible with the version of Tanzu Kubernetes Grid you are upgrading to.  

To view the version of an SDDC, select **View Details** on the SDDC tile in VMware Cloud Console and click on the **Support** pane.

To validate compatibility with Tanzu Kubernetes Grid, refer to the [VMware Product Interoperablity Matrix](https://interopmatrix.vmware.com/#/Interoperability).

## <a id="aws"></a> Prepare to Upgrade Clusters on Amazon EC2

No specific action is required for either major v1.2.x to v1.3.1 or minor v1.3.0 to v1.3.1 upgrades.

Amazon Linux 2 Amazon Machine Images (AMI) that include the supported Kubernetes versions are publicly available to all Amazon EC2 users, in all supported AWS regions. Tanzu Kubernetes Grid automatically uses the appropriate AMI for the Kubernetes version that you specify during upgrade.

In Tanzu Kubernetes Grid v1.2 and later, you created the required IAM resources by enabling the **Automate creation of AWS CloudFormation Stack** checkbox in the installer interface or by running the `tkg config permissions aws` command from the CLI.

When upgrading your management cluster to v1.3.x, you do not need to recreate the AWS CloudFormation Stack.

For more information, see [Deploy Management Clusters with the Installer Interface](../mgmt-clusters/deploy-ui.md) or [Deploy Management Clusters to with the CLI](../mgmt-clusters/deploy-cli.md).

## <a id="azure"></a> Prepare to Upgrade Clusters on Azure

This step is required for both major v1.2.x to v1.3.1 and minor v1.3.0 to v1.3.1 upgrades.

Before upgrading a Tanzu Kubernetes Grid deployment on Azure, you must accept the terms for the new default VM image and for each non-default VM image that you plan to use for your cluster VMs. You need to accept these terms once per subscription.

To accept the terms:

1. List all available VM images for Tanzu Kubernetes Grid in the Azure Marketplace:

   ```
   az vm image list --publisher vmware-inc --offer tkg-capi --all
   ```

1. Accept the terms for the new default VM image:

   ```
   az vm image terms accept --urn publisher:offer:sku:version
   ```

   For example, to accept the terms for the default VM image in Tanzu Kubernetes Grid v1.3.1, `k8s-1dot20dot5-ubuntu-2004`, run:

   ```
   az vm image terms accept --urn vmware-inc:tkg-capi:k8s-1dot20dot5-ubuntu-2004:2021.05.17
   ```

1. If you plan to upgrade any of your Tanzu Kubernetes clusters to a non-default Kubernetes version, such as v1.19.8 or v1.18.16, accept the terms for each non-default version that you want to use for your cluster VMs.

## <a id="set-variable"></a> Set the <code>TKG_BOM_CUSTOM_IMAGE_TAG</code>

This step is required for both major v1.2.x to v1.3.1 and minor v1.3.0 to v1.3.1 upgrades.

Before you can upgrade a management cluster to v1.3, you must specify the correct BOM file to use as a local environment variable. In the event of a patch release to Tanzu Kubernetes Grid, the BOM file may require an update to coincide with updated base image files.

**Note** For information about the most recent security patch updates to VMware Tanzu Kubernetes Grid v1.3, see the [VMware Tanzu Kubernetes Grid v1.3.1 Release Notes](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3.1/rn/VMware-Tanzu-Kubernetes-Grid-131-Release-Notes.html) and this [Knowledgebase Article](https://kb.vmware.com/s/article/83781).

On the machine where you run the Tanzu CLI, perform the following steps:

1. Remove any existing BOM data.

   ```
   rm -rf ~/.tanzu/tkg/bom
   ```

1. Specify the updated BOM to use by setting the following variable.

   ```
   export TKG_BOM_CUSTOM_IMAGE_TAG="v1.3.1-patch1"
   ```

1. Run `tanzu management-cluster create` command with no additional parameters.

   ```
   tanzu management-cluster create
   ```

   This command produces an error but results in the BOM files being downloaded to `~/.tanzu/tkg/bom`.

## <a id="mgmt-cluster"></a> Upgrade Management Clusters

This step is required for both major v1.2.x to v1.3.1 and minor v1.3.0 to v1.3.1 upgrades.

To upgrade Tanzu Kubernetes Grid, you must upgrade all management clusters in your deployment. You cannot upgrade Tanzu Kubernetes clusters until you have upgraded the management clusters that manage them.

Follow the procedure in [Upgrade Management Clusters](management-cluster.md) to upgrade your management clusters.

## <a id="workload-cluster"></a> Upgrade Workload Clusters

This step is required for both major v1.2.x to v1.3.1 and minor v1.3.0 to v1.3.1 upgrades.

After you upgrade the management clusters in your deployment, you can upgrade the Tanzu Kubernetes clusters that are managed by those management clusters.

Follow the procedure in [Upgrade Tanzu Kubernetes Clusters](workload-clusters.md) to upgrade the Tanzu Kubernetes clusters that are running your workloads.

## <a id="extensions"></a> Upgrade the Tanzu Kubernetes Grid Extensions

This step is required for both major v1.2.x to v1.3.1 and minor v1.3.0 to v1.3.1 upgrades.

If you implemented any or all of the Tanzu Kubernetes Grid extensions in Tanzu Kubernetes Grid v1.2.x or v1.3.0, you must upgrade the extensions after you upgrade your management and workload clusters to Tanzu Kubernetes Grid v1.3.1. For information about how to upgrade the extensions, see [Upgrade Tanzu Kubernetes Grid Extensions](extensions.md).

## <a id="add-ons"></a> Register Core Add-ons

This step is only required for major upgrades from Tanzu Kubernetes Grid v1.2.x to v1.3.x.

After you upgrade your management and workload clusters from Tanzu Kubernetes Grid v1.2.x to v1.3.x, follow the instructions in [Register Core Add-ons](addons.md) to register the CNI, vSphere CPI, vSphere CSI, Pinniped, and Metrics Server add-ons with `tanzu-addons-manager`, the component that manages the lifecycle of add-ons.

## <a id="crashd"></a> Upgrade Crash Recovery and Diagnostics

This step is required for both major v1.2.x to v1.3.1 and minor v1.3.0 to v1.3.1 upgrades.

For information about how to upgrade Crash Recovery and Diagnostics, see [Install or Upgrade the Crash Recovery and Diagnostics Binary](../troubleshooting-tkg/crashd.md#install).

## <a id="what-next"></a> What to Do Next

- Examine your upgraded management clusters or register them in Tanzu Mission Control. See [Examine the Management Cluster Deployment](../mgmt-clusters/verify-deployment.md) and [Register Your Management Cluster with Tanzu Mission Control](../mgmt-clusters/register_tmc.md).
- If you have not done so already, enable identity management in Tanzu Kubernetes Grid. See [Enabling Identity Management in Tanzu Kubernetes Grid](../mgmt-clusters/enabling-id-mgmt.md).
