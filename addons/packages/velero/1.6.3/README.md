# Velero

This package provides safe operations to backup, restore, perform disaster recovery, and migrate TCE cluster resources and persistent volumes using the open source tool [Velero](https://velero.io/).

## Components

### Custom Resources

Each Velero operation – on-demand backup, scheduled backup, restore – is a custom resource, defined with a Kubernetes [Custom Resource Definition (CRD)](https://kubernetes.io/docs/concepts/api-extension/custom-resources/#customresourcedefinitions) and stored in [etcd](https://github.com/coreos/etcd).  Each time these operations are run, an equivalent Kubernetes object is created and saved to storage.

Because of the Kubernetes native way that Velero operates, you can back up or restore all objects in a TCE cluster, or you can filter what object to operate on by type, namespace, and/or label.

### CLI

Once the Velero package is installed, you will need to have the Velero CLI to run operations on the command line. Please see this documentation for how to install it: [Velero Docs - CLI Install](https://velero.io/docs/v1.6/basic-install/#install-the-cli).

### Storage

Velero needs an object storage where to save all resource backups and the information about snapshot backups.

Velero treats this object storage as the source of truth: it continuously checks to see that the correct backups are always present. If there is a Velero backup resource in the storage bucket, but no corresponding backup object in the TCE cluster, Velero synchronizes the information from object storage to TCE. Likewise, if a backup object exists in TCE, but not in the object storage, it will be deleted from TCE since the backup resource no longer exists, likely because it was delete.

The TCE Velero package supports these storage providers: AWS, Azure, and Minio.

### Server

Velero runs on the cluster as a deployment alongside with installed plugins that are specific to storage providers for different backup and snapshot operations. It also includes controllers that process the custom resources to perform backups, restores, and all related operations.

## Configuration

WIP
