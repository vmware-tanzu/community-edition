# Velero

This package provides safe operations to backup, restore, perform disaster recovery, and migrate TCE cluster resources and persistent volumes using the open source tool [Velero](https://velero.io/).

## Components

### Custom Resources

Each Velero operation – on-demand backup, scheduled backup, restore – is a custom resource, defined with a Kubernetes [Custom Resource Definition (CRD)](https://kubernetes.io/docs/concepts/api-extension/custom-resources/#customresourcedefinitions) and stored in [etcd](https://github.com/coreos/etcd).  Each time these operations are run, an equivalent Kubernetes object is created and saved to storage.

Because of this Kubernetes native way that Velero operates, you are not restricted to backing up the entire etcd. You can back up or restore all objects in a TCE cluster, or you can also filter what object to operate on by type, namespace, and/or label.

### CLI

Once the Velero package is installed, you will need to have the Velero CLI to run operations on the command line. Please see this documentation for how to install it: [Velero CLI install](https://velero.io/docs/v1.6/basic-install/#install-the-cli).

### Storage

Velero needs an object storage where to save all resource backups and the information about snapshot backups.

Velero treats object storage as the source of truth. It continuously checks to see that the correct backup resources are always present. If there is a properly formatted backup file in the storage bucket, but no corresponding backup resource in the Kubernetes API, Velero synchronizes the information from object storage to Kubernetes. This allows restore functionality to work in a cluster migration scenario, where the original backup objects do not exist in the new cluster. Likewise, if a backup object exists in Kubernetes but not in object storage, it will be deleted from Kubernetes since the backup tarball no longer exists. To learn more, please see the documentation: [How Velero works](https://velero.io/docs/v1.6/how-velero-works/).

### Server

Velero runs on the cluster as a deployment alongside installed plugins that are specific to storage providers for backup and snapshot operations. It also includes controllers that process the custom resources to perform backups, restores, and all related operations.

## Supported Providers

The TCE Velero package provides support for these providers out of the box with minimum configuration.

| Provider                          | Object Store        | Volume Snapshotter           | Plugin Provider Repo                    | Setup Instructions            |
|-----------------------------------|---------------------|------------------------------|-----------------------------------------|-------------------------------|
| [Amazon Web Services (AWS)](https://aws.amazon.com)    | AWS S3              | AWS EBS                      | [Velero plugin for AWS](https://github.com/vmware-tanzu/velero-plugin-for-aws)              | [AWS Plugin Setup](https://github.com/vmware-tanzu/velero-plugin-for-aws#setup)        |
| [Microsoft Azure](https://azure.com)                                       | Azure Blob Storage  | Azure Managed Disks          | [Velero plugin for Microsoft Azure](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure) | [Azure Plugin Setup](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure#setup)      |
| [VMware vSphere](https://github.com/vmware-tanzu/velero-plugin-for-vsphere) | N/A                | 🚫 vSphere Volumes  (on the roadmap)            | [VMware vSphere](https://github.com/vmware-tanzu/velero-plugin-for-vsphere)                    | [vSphere Plugin Setup](https://github.com/vmware-tanzu/velero-plugin-for-vsphere#velero-plugin-for-vsphere-installation-and-configuration-details)

Some other third-party storage providers, like MinIO, DigitalOcean, and others, support the same S3 API that the **AWS Velero plugin** uses.  For more information please see: [S3-Compatible object store providers for Velero](https://velero.io/docs/v1.6/supported-providers/#s3-compatible-object-store-providers).

## Configuration

The following configuration values can be set to customize the Velero installation for the different components.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy Velero. |

### Storage settings

#### Global configurations for storage

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `backupStorageLocation.name` | Required | The name of the Backup Storage Location. |
| `backupStorageLocation.provider` | Required | The cloud provider to use. One of: `aws`, `azure`. |
| `backupStorageLocation.default` | Optional | Indicates if this location is the default backup storage location. |
| `backupStorageLocation.objectStorage.bucket` | Required | The storage bucket where backups are to be uploaded. |
| `backupStorageLocation.objectStorage.prefix` | Optional | The directory inside a storage bucket where backups are to be uploaded. |

#### AWS storage

| `backupStorageLocation.configAWS.region` | Required | The AWS region where the S3 bucket is located. |

#### Azure storage

| `backupStorageLocation.configAzure.resourceGroup` | Required | The name of the resource group containing the storage account for this backup storage location. |
| `backupStorageLocation.configAzure.storageAccount` | Required | The name of the storage account for this backup storage location. |
| `backupStorageLocation.configAzure.storageAccountKeyEnvVar` | Required | Required if using a storage account access key to authenticate rather than a service principal. |
| `backupStorageLocation.configAzure.subscriptionId` | Optional | The the ID of the subscription for this backup storage location. |

### Volume snapshot settings

#### Global configurations for snapshotting

#### AWS volumes

#### Azure volumes

### Advanced

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `rbac.create` | Optional | Whether to create the Velero Role and RoleBinding to give all permissions to the namespace to Velero.|
| `rbac.name` | Optional |  A new name for the cluster RolBinding. Default is `velero`. |
| `rbac.clusterAdministrator` | Optional | Whether to create the ClusterRoleBinding to give administrator permissions to Velero. `rbac.create` must also be set to `true`.|
| `rbac.roleRefName` | Optional | Name of the cluster role to reference. Default is `cluster-admin`.|
| `rbac.clusterRoleAPIGroups` | Optional |  The name of the API groups that contain the resources for the cluster role.|
| `rbac.clusterRoleVerbs` | Optional |  The set of verbs that apply to the secret resources contained in this rule. |
| `serviceAccount.name` | Optional |  The name of the ServiceAccount the RoleBinding should reference. |
| `serviceAccount.annotations` | Optional |  Annotations for the ServiceAccount the RoleBinding should reference. |
| `serviceAccount.labels` | Optional |  Labels for the ServiceAccount the RoleBinding should reference. |

## Install and update

Steps to install and configure:

- Install the package:

Set `pkgname` to any name to identify the package by.
Set `namespace` to any existing namespace.

```sh
tanzu package install $pkgname --package-name velero.community.tanzu.vmware.com --version 1.6.3 --namespace $namespace
```

- Verify it was properly installed:

```sh
tanzu package installed list -A
```

It should be installed in the namespace specified in the configuration. The default is `velero`.

- Download and alter the ``values.yaml`` file from the TCE repository. Delete the first line the file starting with  `#@`, as this will cause an error when updating the package.

- Update the package:

```sh
tanzu package installed update $pkgname --version 1.6.3 --namespace $namespace --values-file $path/values.yaml
```

## Usage Example

This walkthrough guides you through an example disaster recovery scenario that leverages the Velero package. You must deploy the package before attempting this walkthrough.

⚠️ Note: For more advanced use cases and documentation, see the official Velero [documentation](https://velero.io/docs/latest/).

In the following steps, you will simulate a disaster scenario. Specifically, you will deploy a stateless workload, create a backup, delete the workload, and restore it from the backup.

- Download the Velero CLI from the GitHub [releases](https://github.com/vmware-tanzu/velero/releases/latest) page. The following steps assume you have installed Velero into your PATH.

- Create a new namespace for this example:

```sh
kubectl create ns velero-example
```

- Deploy a sample workload into the new namespace:

```sh
kubectl create deploy -n velero-example nginx --image=nginx
```

- Verify the workload is up and running:

```sh
kubectl get pods -n velero-example
```

The output should be similar to the following:

```sh
NAME                     READY   STATUS    RESTARTS   AGE
nginx-6799fc88d8-mm47k   1/1     Running   0          7s
```

- Create a backup of the `velero-example` namespace:

```sh
velero create backup velero-example --include-namespaces velero-example
```

- Verify the backup completed successfully:

```sh
velero describe backup velero-example
```

The output shows the "Phase" of the backup, which should be `Completed`.

- Delete the `velero-example` namespace to simulate a disaster scenario:

```sh
kubectl delete ns velero-example
```

- Verify that the namespace has been deleted:

```sh
kubectl get ns
```

- Restore the namespace from the velero backup:

```sh
velero create restore --from-backup velero-example
```

- Validate that the `velero-example` namespace has been restored:

```sh
kubectl get ns velero-example
```

- Validate that the workload has been restored:

```sh
kubectl get pods -n velero-example
```
