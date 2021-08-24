# Back Up and Restore Clusters

To back up and restore Tanzu Kubernetes Grid clusters, you can use [Velero](https://velero.io/docs), an open source community standard tool for backing up and restoring Kubernetes cluster objects and persistent volumes.
Velero supports a variety of [storage providers](https://velero.io/docs/main/supported-providers/) to store its backups.

If a Tanzu Kubernetes Grid management or workload cluster crashes and fails to recover, the infrastructure administrator can use a Velero backup to restore its contents to a new cluster, including cluster extensions and internal API objects for the workload clusters.

The following sections explain how to set up a Velero server on Tanzu Kubernetes Grid management or workload clusters, and direct it from the Velero CLI to back up and restore the clusters.

**NOTE**: You must create a new cluster to restore to; you cannot restore a cluster backup to an existing cluster.

## <a id="setup"></a> Setup Overview

To back up and restore Tanzu Kubernetes Grid clusters, you need:

- The Velero CLI running on your local workstation; see [Install the Velero CLI](#cli).
- A storage provider with locations to save the backups to; see [Set Up a Storage Provider](#storage).
- A Velero server running on each cluster that you want to back up; see [Deploy Velero Server to Clusters](#server).

## <a id="cli"></a> Install the Velero CLI

1. Go to [the Tanzu Kubernetes Grid downloads page](https://my.vmware.com/en/web/vmware/downloads/info/slug/infrastructure_operations_management/vmware_tanzu_kubernetes_grid/1_x) and log in with your My VMware credentials.
1. Under **Product Downloads**, click **Go to Downloads**.
1. Scroll to the **Velero** entries and download the Velero CLI `.gz` file for your workstation OS. Its filename starts with `velero-linux-`, `velero-mac-`, or `velero-windows64-`.
1. Use the `gunzip` command or the extraction tool of your choice to unpack the binary:

   ```sh  
   gzip -d <RELEASE-TARBALL-NAME>.gz
   ```

1. Rename the CLI binary for your platform to `velero`, make sure that it is executable, and add it to your `PATH`.

   - macOS and Linux platforms:

       1. Move the binary into the `/usr/local/bin` folder and rename it to `velero`.
       1. Make the file executable:

       ```sh  
       chmod +x /usr/local/bin/velero
       ```

   - Windows platforms:

        1. Create a new `Program Files\velero` folder and copy the binary into it.
        1. Rename the binary to `velero.exe`.
        1. Right-click the `velero` folder, select **Properties** > **Security**, and make sure that your user account has the **Full Control** permission.
        1. Use Windows Search to search for `env`.
        1. Select **Edit the system environment variables** and click the **Environment Variables** button.
        1. Select the `Path` row under **System variables**, and click **Edit**.
        1. Click **New** to add a new row and enter the path to the `velero` binary.
1. On **vSphere with Tanzu**:
    - Install the vSphere plugin for `kubectl` by following the procedure in [Download and Install the Kubernetes CLI Tools for vSphere](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-0F6E45C4-3CB1-4562-9370-686668519FCA.html).
    - This plugin retrieves the Supervisor Cluster credentials to enable the Velero CLI to access it.

## <a id="storage"></a> Set Up a Storage Provider

To back up Tanzu Kubernetes Grid clusters, you need storage locations for:

- Cluster object storage backups for Kubernetes metadata in both management clusters and workload clusters
- Volume snapshots for data used by workload clusters

See [Backup Storage Locations and Volume Snapshot Locations](https://velero.io/docs/main/locations/) in the Velero documentation.
Velero supports a variety of [storage providers](https://velero.io/docs/main/supported-providers).

VMware recommends dedicating a unique storage bucket to each cluster.

### <a id="storage-vsphere"></a> Storage for vSphere

On vSphere, cluster object storage backups and volume snapshots save to the same storage location.
This location must be S3-compatible external storage on Amazon Web Services (AWS), or an S3 provider such as [MinIO](https://min.io).

To set up storage for Velero on vSphere, follow the installation procedures in the Velero Plugin for AWS repository, depending on what kind of cluster you are backing up:

* **Management cluster or workload cluster deployed by Tanzu Kubernetes Grid**: See [Velero Plugin for vSphere in Vanilla Kubernetes Cluster](https://github.com/vmware-tanzu/velero-plugin-for-vsphere/blob/v1.1.0/docs/vanilla.md) for the v1.1.0 plugin
* **vSphere with Tanzu Supervisor cluster**: See [Velero Plugin for vSphere in vSphere with Tanzu Supervisor Cluster](https://github.com/vmware-tanzu/velero-plugin-for-vsphere/blob/v1.1.0/docs/supervisor.md) for the v1.1.0 plugin

### <a id="storage-aws"></a> Storage for and on AWS

To set up storage for Velero on AWS, follow the procedures in the Velero Plugins for AWS repository:

1. [Create an S3 bucket](https://github.com/vmware-tanzu/velero-plugin-for-aws#Create-S3-bucket)

1. [Set permissions for Velero](https://github.com/vmware-tanzu/velero-plugin-for-aws#Set-permissions-for-Velero)

Set up S3 storage as needed for each plugin.  The object store plugin stores and retrieves cluster object backups, and the volume snapshotter stores and retrieves data volumes.

### <a id="storage-azure"></a> Storage for and on Azure

To set up storage for Velero on Azure, follow the procedures in the Velero Plugins for Azure repository:

1. [Create an Azure storage account and blob container](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure#Create-Azure-storage-account-and-blob-container)

1. [Get the resource group containing your VMs and disks](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure#Get-resource-group-containing-your-VMs-and-disks)

1. [Set permissions for Velero](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure#Set-permissions-for-Velero)

Set up S3 storage as needed for each plugin.  The object store plugin stores and retrieves cluster object backups, and the volume snapshotter stores and retrieves data volumes.

## <a id="server"></a> Deploy Velero Server to Clusters

To deploy the Velero Server to a cluster, you run the `velero install` command.
This command creates a namespace called `velero` on the cluster, and places a deployment named `velero` in it.

`velero install` installs the Velero server to the current default cluster in your `kubeconfig`, or else you can specify a different cluster with the `--kubeconfig` flag.

How you run the `velero install` command and otherwise set up Velero on a cluster depends on your infrastructure and storage provider:

### <a id="server-vsphere-no-tanzu"></a> Velero Server on vSphere without Tanzu

This procedure applies to management clusters deployed by Tanzu Kubernetes Grid and the workload clusters they create.
If a vSphere with Tanzu Supervisor cluster serves as your Tanzu Kubernetes Grid management cluster, see the **vSphere with Tanzu** instructions below.

1. Install the Velero server to the current default cluster in your `kubeconfig` by running `velero install`, as described in the [Install](https://github.com/vmware-tanzu/velero-plugin-for-vsphere/blob/v1.1.0/docs/vanilla.md#install) section for Vanilla Kubernetes clusters in the Velero Plugin for vSphere v1.1.0 repository.  Include option values as follows:
  - `--provider aws`
  - `--plugins velero/velero-plugin-for-aws:v1.1.0`
  - `--bucket $BUCKET`: name of your S3 bucket
  - `--backup-location-config region=$REGION`: region the bucket is in
  - `--snapshot-location-config region=$REGION`: region the bucket is in
  - For bucket access via username / password, include `--secret-file ./velero-creds` pointing to local file that looks like:

      ```
      [default]
      aws_access_key_id=<AWS_ACCESS_KEY_ID>
      aws_secret_access_key=<AWS_SECRET_ACCESS_KEY>
      ```

  - For bucket access via `kube2iam`:

      ```
      --pod-annotations iam.amazonaws.com/role=arn:aws:iam::<AWS_ACCOUNT_ID>:role/<VELERO_ROLE_NAME>``
      --no-secret
      ```

  - (Optional) `--kubeconfig` to install the Velero server to a cluster other than the current default.
  - For additional options, see [Install and start Velero](https://github.com/vmware-tanzu/velero-plugin-for-aws#install-and-start-velero).

  For example, to use MinIO as the object storage, following the [MinIO server setup instructions](https://velero.io/docs/main/contributions/minio/#set-up-server) in the Velero documentation:

  ```sh  
  velero install --provider aws --plugins "velero/velero-plugin-for-aws:v1.1.0" --bucket velero --secret-file ./credentials-velero --backup-location-config "region=minio,s3ForcePathStyle=true,s3Url=minio_server_url" --snapshot-location-config region="default"
  ```

  Installing the Velero server to the cluster creates a namespace in the cluster called `velero`, and places a deployment named `velero` in it.

1. For workload clusters with CSI-based volumes, add the Velero Plugin for vSphere. This plugin lets Velero use your S3 bucket to store CSI volume snapshots for workload data, in addition to storing cluster objects:
    1. Download the [Velero Plugin for vSphere](https://github.com/vmware-tanzu/velero-plugin-for-vsphere/tree/v1.1.0) v1.1.0 image.
    1. Run `velero plugin add PLUGIN-IMAGE` with the plugin image name.
        - `PLUGIN-IMAGE` is the container image name listed in the [Velero Plugin for vSphere repo](https://github.com/vmware-tanzu/velero-plugin-for-vsphere/releases/tag/v1.1.0) v1.1.0, for example, `vsphereveleroplugin/velero-plugin-for-vsphere:1.1.0`.
    1. Enable the plugin by adding the following **VirtualMachine** permissions to the role you created for the Tanzu Kubernetes Grid account, if you did not already include them when you [created the account](../mgmt-clusters/vsphere.md#vsphere-permissions):
        - **Configuration** > **Toggle disk change tracking**
        - **Provisioning** > **Allow read-only disk access**
        - **Provisioning** > **Allow virtual machine download**
        - **Snapshot management** > **Create snapshot**
        - **Snapshot management** > **Remove snapshot**

### <a id="server-vsphere-tanzu"></a> Velero Server on vSphere with Tanzu

vSphere with Tanzu Supervisor clusters do not have the Kubernetes API server permissions required to retrieve Kubernetes cluster objects, so you need to install Velero with a `Velero vSphere Operator` that elevates Velero's permissions.

To install Velero on a Supervisor cluster, follow [Installing Velero on a Supervisor Cluster](https://github.com/vmware-tanzu/velero-plugin-for-vsphere/blob/v1.1.0/docs/velero-vsphere-operator-user-manual.md#installing-velero-on-supervisor-cluster) in the Velero Plugin for vSphere v1.1.0 repository.

**NOTE:** Tanzu Kubernetes Grid does not support backing up the Kubernetes object metadata for the Supervisor cluster, which captures its state.
You can use Velero to back up data volume snapshots for user workloads running on the Supervisor cluster, as well as objects and data from workload clusters managed by the Supervisor cluster.

### <a id="server-vsphere-aws"></a> Velero Server on AWS

To install Velero on clusters on AWS, follow the [Install and start Velero](https://github.com/vmware-tanzu/velero-plugin-for-aws#Install-and-start-Velero) procedure in the Velero Plugins for AWS repository.

### <a id="server-vsphere-azure"></a> Velero Server on Azure

To install Velero on clusters on Azure, follow the [Install and start Velero](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure#Install-and-start-Velero) procedure in the Velero Plugins for Azure repository.

## <a id="vsphere"></a> vSphere Backup and Restore

These sections describe how to back up and restore Tanzu Kubernetes Grid clusters on vSphere.

### <a id="vsphere-backup"></a> Back Up a Cluster on vSphere

1. Follow the [Deploy Velero Server to Clusters](#server) instructions above to deploy a Velero server onto the cluster, along with the Velero Plugin for vSphere if needed.

1. If you are backing up a management cluster, set `Cluster.Spec.Paused` to `true` for all of its workload clusters:

   ```sh  
   kubectl patch cluster workload_cluster_name --type='merge' -p '{"spec":{"paused": true}}'
   ```

1. Back up the cluster:

   ```sh  
   velero backup create your_backup_name --exclude-namespaces=tkg-system
   ```

   Excluding `tkg-system` objects avoids creating duplicate cluster API objects when restoring a cluster.

1. If you backed up a management cluster, set `Cluster.Spec.Paused` back to `false` for all of its workload clusters.

### <a id="vsphere-restore"></a> Restore a Cluster on vSphere

1. Create a new cluster. You cannot restore a cluster backup to an existing cluster.

1. Follow the [Deploy Velero Server to Clusters](#server) instructions above to deploy a Velero server onto the new cluster, along with the Velero Plugin for vSphere if needed.

1. Restore the cluster:

   ```sh  
   velero restore create your_restore_name --from-backup your_backup_name
   ```

1. Set `Cluster.Spec.Paused` field to `false` for all workload clusters:

   ```sh
   kubectl patch cluster cluster_name --type='merge' -p '{"spec":{"paused": false}}'
   ```

## <a id="aws"></a> AWS Backup and Restore

These sections describe how to back up and restore clusters on AWS.

### <a id="aws-backup"></a> Back Up a Cluster on AWS

1. Follow the [Velero Plugin for AWS setup instructions](https://github.com/vmware-tanzu/velero-plugin-for-aws#setup) to install Velero server on the cluster.

1. If you are backing up a management cluster, set `Cluster.Spec.Paused` to `true` for all of its workload clusters:

   ```sh
   kubectl patch cluster workload_cluster_name --type='merge' -p '{"spec":{"paused": true}}'
   ```

1. Back up the cluster:

   ```sh  
   velero backup create your_backup_name --exclude-namespaces=tkg-system
   ```

   Excluding `tkg-system` objects avoids creating duplicate cluster API objects when restoring a cluster.

1. If you backed up a management cluster, set `Cluster.Spec.Paused` back to `false` for all of its workload clusters.

### <a id="aws-restore"></a> Restore a Cluster on AWS

1. Create a new cluster. You cannot restore a cluster backup to an existing cluster.

1. Follow the [Velero Plugin for AWS setup instructions](https://github.com/vmware-tanzu/velero-plugin-for-aws#setup) to install Velero server on the new cluster.

1. Restore the cluster:

   ```sh  
   velero backup get
   velero restore create your_restore_name --from-backup your_backup_name
   ```

1. Set `Cluster.Spec.Paused` to `false` for all workload clusters:

   ```sh
   kubectl patch cluster cluster_name --type='merge' -p '{"spec":{"paused": false}}'
   ```

## <a id="azure"></a> Azure Backup and Restore

These sections describe how to back up and restore clusters on Azure.

### <a id="azure-backup"></a> Back Up a Cluster on Azure

1. Follow the [Velero Plugin for Azure setup instructions](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure#setup) to install Velero server on the cluster.

1. If you are backing up a management cluster, set `Cluster.Spec.Paused` to `true` for all of its workload clusters:

   ```sh
   kubectl patch cluster workload_cluster_name --type='merge' -p '{"spec":{"paused": true}}'
   ```

1. Back up the cluster:

   ```sh  
   velero backup create your_backup_name --exclude-namespaces=tkg-system
   ```

   Excluding `tkg-system` objects avoids creating duplicate cluster API objects when restoring a cluster.

1. If `velero backup` returns a `transport is closing` error, try again after increasing the memory limit, as described in [Update resource requests and limits after install](https://velero.io/docs/v1.4/customize-installation/#update-resource-requests-and-limits-after-install) in the Velero documentation.

1. If you backed up a management cluster, set `Cluster.Spec.Paused` back to `false` for all of its workload clusters.

### <a id="azure-restore"></a> Restore a Cluster on Azure

1. Create a new cluster. You cannot restore a cluster backup to an existing cluster.

1. Follow the [Velero Plugin for Azure setup instructions](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure#setup) to install Velero server on the new cluster.

1. Restore the cluster:

   ```sh  
   velero backup get
   velero restore create your_restore_name --from-backup your_backup_name
   ```

1. Set `Cluster.Spec.Paused` to `false` for all workload clusters:

   ```sh
   kubectl patch cluster cluster_name --type='merge' -p '{"spec":{"paused": false}}'
   ```
