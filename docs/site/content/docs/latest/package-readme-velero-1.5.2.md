# Velero

This package provides disaster recovery capabilities using [velero](https://velero.io/). At the moment, it leverages [minio](https://github.com/minio/minio) for object storage.

## Components

* velero Namespace
* velero Custom Resources
* velero Deployment
* cloud-credentials Secret (contains the credentials for Velero to authenticate with minio)
* minio Deployment
* minio-setup Job (configures/initializes minio)
* minio Service

## Configuration

The following configuration values can be set to customize the Velero installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy Velero. |

### Velero Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `provider` | Required | The cloud provider in use. One of: `aws`, `azure`, `vsphere`, `docker`. |
| `csi.enabled` | Optional | Whether to enable Velero's CSI features. Defaults to `false`. |
| `backupStorageLocation.name` | Optional | The name of the Backup Storage Location. |
| `backupStorageLocation.bucket` | Required | The storage bucket where backups are to be uploaded. |
| `backupStorageLocation.prefix` | Optional | The directory inside a storage bucket where backups are to be uploaded. |
| `backupStorageLocation.aws.region` | Required when the provider is `aws` | The AWS region where the S3 bucket is located. |
| `volumeSnapshotLocation.name` | Optional | The name of the Volume Snapshot Location. |
| `volumeSnapshotLocation.aws.region` | Required when the provider is `aws` | The AWS region where the Volumes and Snapshots are located. |

## Usage Example

This walkthrough guides you through an example disaster recovery scenario that leverages the Velero package. You must deploy the package before attempting this walkthrough.

⚠️ Note: For more advanced use cases and documentation, see the official Velero [documentation](https://velero.io/docs/latest/).

In the following steps, you will simulate a disaster scenario. Specifically, you will deploy a stateless workload, create a backup, delete the workload, and restore it from the backup.

1. Download the Velero CLI from the GitHub [releases](https://github.com/vmware-tanzu/velero/releases/latest) page. The following steps assume you have installed Velero into your PATH.

1. Create a new namespace for this example:

    ```bash
    kubectl create ns velero-example
    ```

1. Deploy a sample workload into the new namespace:

    ```bash
    kubectl create deploy -n velero-example nginx --image=nginx
    ```

1. Verify the workload is up and running:

    ```bash
    kubectl get pods -n velero-example
    ```

    The output should be similar to the following:

    ```bash
    NAME                     READY   STATUS    RESTARTS   AGE
    nginx-6799fc88d8-mm47k   1/1     Running   0          7s
    ```

1. Create a backup of the `velero-example` namespace:

    ```bash
    velero create backup velero-example --include-namespaces velero-example
    ```

1. Verify the backup completed successfully:

    ```bash
    velero describe backup velero-example
    ```

    The output shows the "Phase" of the backup, which should be `Completed`.

1. Delete the `velero-example` namespace to simulate a disaster scenario:

    ```bash
    kubectl delete ns velero-example
    ```

1. Verify that the namespace has been deleted:

    ```bash
    kubectl get ns
    ```

1. Restore the namespace from the velero backup:

    ```bash
    velero create restore --from-backup velero-example
    ```

1. Validate that the `velero-example` namespace has been restored:

    ```bash
    kubectl get ns velero-example
    ```

1. Validate that the workload has been restored:

    ```bash
    kubectl get pods -n velero-example
    ```
