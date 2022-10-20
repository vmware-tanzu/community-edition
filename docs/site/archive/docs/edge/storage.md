# Create Persistent Volumes with Storage Classes

This topic explains how to use dynamic storage in workload clusters.

## Overview: PersistentVolume, PersistentVolumeClaim, and StorageClass

Within a Kubernetes cluster, `PersistentVolume` (PV) objects provide shared storage for cluster pods that is unaffected by pod lifecycles.
Storage is provisioned to the PV through a `PersistentVolumeClaim` (PVC) object, which defines how much and how the pod accesses the underlying storage.
For more information, see [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) in the Kubernetes documentation.

Cluster administrators can define `StorageClass` objects that let cluster users dynamically create PVC and PV objects with different storage types and rules.
Tanzu also provides default `StorageClass` objects that let users provision persistent storage in a turnkey environment.

`StorageClass` objects include a `provisioner` field identifying the internal or external service plug-in that provisions PVs, and a `parameters` field that associates the Kubernetes storage class with storage options defined at the infrastructure level, such as VM Storage Policies in vSphere.
For more information, see [Storage Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/) in the Kubernetes documentation.

## Supported Storage Types

Tanzu Community Edition supports `StorageClass` objects for different storage types, provisioned by Kubernetes internal ("in-tree") or external ("out-of-tree") plug-ins.

### Storage Types

- vSphere Cloud Native Storage (CNS)
- Amazon EBS
- Azure Disk
- iSCSI
- NFS

See [Default Storage Classes](#defaults) below for vSphere CNS, Azure EBS, and Azure Disk default storage classes.

### Plug-in Locations

- Kubernetes internal ("in-tree") storage.
      - Ships with core Kubernetes; `provider` values are prefixed with `kubernetes.io`, e.g. `kubernetes.io/aws-ebs`.
- External ("out-of-tree") storage.
      - Can be anywhere defined by `provider` value, e.g. `csi.vsphere.vmware.com`.
      - Follow the Container Storage Interface (CSI) standard for external storage.

## Default Storage Classes

Tanzu provides default `StorageClass` objects that let workload cluster users provision persistent storage on their infrastructure in a turnkey environment, without needing `StorageClass` objects created by a cluster administrator.

The `ENABLE_DEFAULT_STORAGE_CLASS` variable is set to `true` by default in the cluster configuration file passed to `--file` option of `tanzu cluster create`, to enable the default storage class for a workload cluster.

The Tanzu default storage class definitions are:

### vSphere CNS

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: default
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: csi.vsphere.vmware.com
parameters:
  storagePolicyName: optional
```

See the [vSphere](https://kubernetes.io/docs/concepts/storage/storage-classes/#vsphere) CSI storage class parameters in the Kubernetes documentation.

### Amazon EBS

```sh
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: default
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: kubernetes.io/aws-ebs
```

See the [Amazon EBS](https://kubernetes.io/docs/concepts/storage/storage-classes/#aws-ebs) storage class parameters in the Kubernetes documentation.

### Azure Disk

```sh
apiVersion: storage.k8s.io/v1beta1
kind: StorageClass
metadata:
  name: default
  annotations:
    storageclass.beta.kubernetes.io/is-default-class: "true"
  labels:
    kubernetes.io/cluster-service: "true"
provisioner: kubernetes.io/azure-disk
parameters:
  kind: Managed
  storageaccounttype: Standard_LRS
  cachingmode: ReadOnly
volumeBindingMode: WaitForFirstConsumer
```

See the [Azure Disk](https://kubernetes.io/docs/concepts/storage/storage-classes/#aws-ebs) storage class parameters in the Kubernetes documentation.

## Create a Custom Storage Class

Cluster administrators can create a new storage class as follows:

Before You Begin:

(vSphere only) Select or create the VM storage policy to use as the basis for the Kubernetes `StorageClass`. vSphere administrators can create a storage policy by following the steps in [Set Up vSphere CNS and Create a Storage Policy in vSphere](vsphere-cns).

1. Create a `StorageClass` configuration `.yaml` with `provisioner`, `parameters`, and other options.
2. (vSphere only) Associate a Kubernetes storage class with a vSphere storage policy by setting its `storagePolicyName` parameter to the vSphere storage policy name, as a double-quoted string.
3. Pass the file to `kubectl create -f`
4. Verify the storage class by running `kubectl describe storageclass <storageclass metadata.name>`.

For example, see [Enabling Dynamic Provisioning](https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/#enabling-dynamic-provisioning) in the Kubernetes documentation.

For vSphere CSI information and resources, see [VMware vSphere Container Storage Plug-in Documentation](https://docs.vmware.com/en/VMware-vSphere-Container-Storage-Plug-in/index.html).

## Use a Custom Storage Class in a Cluster

To provision persistent storage for their cluster nodes that does not use one of the [Default Storage Classes](#defaults) described above, cluster users include a custom storage class in a pod configuration as follows:

1. Set the context of `kubectl` to the cluster. For example:

   ```sh
   kubectl config use-context my-cluster-admin@my-cluster
   ```

1. Select or create a storage class.

    - **Select**:
        - To list available storage classes, run `kubectl get storageclass`.
    - **Create**
        - Cluster admins can create storage classes by following [Create a Custom Storage Class](#create), above.

1. Create a PVC and its PV:

    1. Create a `PersistentVolumeClaim` configuration `.yaml` with `spec.storageClassName` set to the `metadata.name` value of your `StorageClass` object.
    For an example, see [Enabling Dynamic Provisioning](https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/#enabling-dynamic-provisioning) in the Kubernetes documentation.
    1. Pass the file to `kubectl create -f`
    1. Run `kubectl describe pvc <pvc metadata.name>` to verify the PVC.
    1. A PV is automatically created with the PVC. Record its name, listed in the `kubectl describe pvd` output after `Successfully provisioned volume`.
    1. Run `kubectl describe pv <pv unique name>` to verify the PV.

1. Create a pod using the PVC:

    1. Create a `Pod` configuration `.yaml` that sets `spec.volumes` to include your PVC under `persistentVolumeClaim.claimName`.
    For an example, see [Dynamic Provisioning and StorageClass API](https://docs.vmware.com/en/VMware-vSphere-Container-Storage-Plug-in/2.0/vmware-vsphere-csp-getting-started/GUID-606E179E-4856-484C-8619-773848175396.html) in the vSphere Storage for Kubernetes documentation.
    1. Pass the file to `kubectl create -f`
    1. Run `kubectl get pod <pod metadata.name>` to verify the pod.
