# Create Persistent Volumes with Storage Classes

This topic explains how to use dynamic storage in Tanzu Kubernetes (workload) clusters in Tanzu Kubernetes Grid.

## <a id="overview"></a>Overview: PersistentVolume, PersistentVolumeClaim, and StorageClass

Within a Kubernetes cluster, `PersistentVolume` (PV) objects provide shared storage for cluster pods that is unaffected by pod lifecycles.
Storage is provisioned to the PV through a `PersistentVolumeClaim` (PVC) object, which defines how much and how the pod accesses the underlying storage.
For more information, see [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) in the Kubernetes documentation.

Cluster administrators can define `StorageClass` objects that let cluster users dynamically create PVC and PV objects with different storage types and rules.
Tanzu Kubernetes Grid also provides default `StorageClass` objects that let users provision persistent storage in a turnkey environment.

`StorageClass` objects include a `provisioner` field identifying the internal or external service plug-in that provisions PVs, and a `parameters` field that associates the Kubernetes storage class with storage options defined at the infrastructure level, such as VM Storage Policies in vSphere.
For more information, see [Storage Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/) in the Kubernetes documentation.

## <a id="types"></a> Supported Storage Types

Tanzu Kubernetes Grid supports `StorageClass` objects for different storage types, provisioned by Kubernetes internal ("in-tree") or external ("out-of-tree") plug-ins.

**Storage Types**

- vSphere Cloud Native Storage (CNS)
- Amazon EBS
- Azure Disk
- iSCSI
- NFS

See [Default Storage Classes](#defaults) below for vSphere CNS, Azure EBS, and Azure Disk default storage classes.

**Plug-in Locations**

- Kubernetes internal ("in-tree") storage.
    - Ships with core Kubernetes; `provider` values are prefixed with `kubernetes.io`, e.g. `kubernetes.io/aws-ebs`.
- External ("out-of-tree") storage.
    - Can be anywhere defined by `provider` value, e.g. `csi.vsphere.vmware.com`.
    - Follow the Container Storage Interface (CSI) standard for external storage.

## <a id="defaults"></a> Default Storage Classes

Tanzu Kubernetes Grid provides default `StorageClass` objects that let workload cluster users provision persistent storage on their infrastructure in a turnkey environment, without needing `StorageClass` objects created by a cluster administrator.

The `ENABLE_DEFAULT_STORAGE_CLASS` variable is set to `true` by default in the cluster configuration file passed to `--file` option of `tanzu cluster create`, to enable the default storage class for a workload cluster.

The Tanzu Kubernetes Grid default storage class definitions are:

**vSphere CNS**

```
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

**Amazon EBS**

```
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: default
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: kubernetes.io/aws-ebs
```

See the [Amazon EBS](https://kubernetes.io/docs/concepts/storage/storage-classes/#aws-ebs) storage class parameters in the Kubernetes documentation.

**Azure Disk**

```
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

## <a id="create-policy"></a> Set Up CNS and Create a Storage Policy (vSphere)

vSphere administrators can set up vSphere CNS and create storage policies for virtual disk (VMDK) storage, based on the needs of Tanzu Kubernetes Grid cluster users.

You can use either vSAN or local VMFS (Virtual Machine File System) for persistent storage in a Kubernetes cluster, as follows:

**vSAN Storage**:

To create a storage policy for vSAN storage in the vSphere Client, browse to **Home** > **Policies and Profiles** > **VM Storage Policies** and click **Create** to launch the **Create VM Storage Policy** wizard.

Follow the instructions in [Create a Storage Policy](https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.storage.doc/GUID-8D51CECC-ED3B-424E-BFE2-43379729A653.html) in the vSphere documentation. Make sure to:

  - In the **Policy structure** pane, under **Datastore specific rules**, select **Enable rules for "vSAN" storage**.
  - Configure other panes or accept defaults as needed.
  - Record the storage policy name for reference as the `storagePolicyName` value in `StorageClass` objects.

**Local VMFS Storage**:

To create a storage policy for local storage, apply a tag to the storage and create a storage policy based on the tag as follows:

1. From the top-level vSphere menu, select **Tags &amp; Custom Attributes**

1. In the **Tags** pane, select **Categories** and click **New**.

1. Enter a category name, such as `tkg-storage`.
Use the checkboxes to associate it with **Datacenter** and the storage objects, **Folder** and **Datastore**.
Click **Create**.

1. From the top-level **Storage** view, select your VMFS volume, and in its **Summary** pane, click **Tags** > **Assign...**.

1. From the **Assign Tag** popup, click **Add Tag**.

1. From the **Create Tag** popup, give the tag a name, such as `tkg-storage-ds1` and assign it the **Category** you created.  Click **OK**.

1. From **Assign Tag**, select the tag and click **Assign**.

1. From top-level vSphere, select **VM Storage Policies** > **Create a Storage Policy**.  A configuration wizard starts.

1. In the **Name and description** pane, enter a name for your storage policy.
Record the storage policy name for reference as the `storagePolicyName` value in `StorageClass` objects.

1. In the **Policy structure** pane, under **Datastore specific rules**, select **Enable tag-based placement rules**.

1. In the **Tag based placement** pane, click **Add Tag Rule** and configure:

    - **Tag category**: Select your category name
    - **Usage option**: `Use storage tagged with`
    - **Tags**: Browse and select your tag name

1. Confirm and configure other panes or accept defaults as needed, then click **Review and finish**. **Finish** to create the storage policy.

## <a id="create-class"></a> Create a Custom Storage Class

Cluster administrators can create a new storage class as follows:

1. On vSphere, select or create the VM storage policy to use as the basis for the Kubernetes `StorageClass`.
    - vSphere administrators can create a storage policy by following [Create a Storage Policy (vSphere)](#create-policy), above.
1. Create a `StorageClass` configuration `.yaml` with `provisioner`, `parameters`, and other options.
    - On vSphere, associate a Kubernetes storage class with a vSphere storage policy by setting its `storagePolicyName` parameter to the vSphere storage policy name, as a double-quoted string.
1. Pass the file to `kubectl create -f`
1. Verify the storage class by running `kubectl describe storageclass <storageclass metadata.name>`.

Examples:

- [Enabling Dynamic Provisioning](https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/#enabling-dynamic-provisioning) in the Kubernetes documentation
- [CSI - Container Storage Interface](https://cloud-provider-vsphere.sigs.k8s.io/container_storage_interface.html) in the Kubernetes vSphere Cloud Provider documentation

## <a id="use"></a> Use a Custom Storage Class in a Cluster

To provision persistent storage for their cluster nodes that does not use one of the [Default Storage Classes](#defaults) described above, cluster users include a custom storage class in a pod configuration as follows:

1. Set the context of `kubectl` to the cluster. For example:

  ```
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
    For an example, see [Dynamic Provisioning and StorageClass API](https://vmware.github.io/vsphere-storage-for-kubernetes/documentation/storageclass.html) in the vSphere Storage for Kubernetes documentation.
    1. Pass the file to `kubectl create -f`
    1. Run `kubectl get pod <pod metadata.name>` to verify the pod.

## <a id="expansion"></a> Enable Offline Volume Expansion for vSphere CSI (vSphere 7)

To enable [offline volume expansion](https://vsphere-csi-driver.sigs.k8s.io/features/volume_expansion.html) for vSphere CSI storage used by workload clusters, you need to add a `csi-resizer` sidecar pod to the cluster's CSI processes.

The CSI configuration for workload clusters is encoded as a Kubernetes secret.
This procedure adds the `csi-resizer` process by revising the CSI configuration secret.
It adds to the secret a `stringData` definition that combines two encoded configuration data strings: a `values.yaml` string containing the secret's prior CSI configuration data, and a new `overlays.yaml` string that deploys the `csi-resizer` pod.

1. Log into the management cluster for the workload cluster you are changing, and run `tanzu cluster list` if you need to retrieve the name of the workload cluster.

1. Retrieve the name of the CSI secret for the workload cluster, using label selectors `vsphere-csi` and the cluster name:

  ```
  $ kubectl get secret \
      -l tkg.tanzu.vmware.com/cluster-name=NAME_OF_WORKLOAD_CLUSTER \
      -l tkg.tanzu.vmware.com/addon-name=vsphere-csi
  my-wc-vsphere-csi-secret
  ```
  
1. Save a backup of the secret's content, in YAML format, to `vsphere-csi-secret.yaml`:

  ```
  kubectl get secret my-wc-vsphere-csi-secret -o yaml > vsphere-csi-secret.yaml
  ```

1. Output the secret's current content again, with the `data.values` values `base64`-decoded into plain YAML.

  ```
  $ kubectl get secret my-wc-vsphere-csi-secret -o jsonpath={.data.values\\.yaml} | base64 -d

  #@data/values
  #@overlay/match-child-defaults missing_ok=True
  ---
  vsphereCSI:
    CSIAttacherImage:
      repository: projects.registry.vmware.com/tkg
      path: csi/csi-attacher
      tag: v3.0.0_vmware.1
      pullPolicy: IfNotPresent
    vsphereCSIControllerImage:
      repository: projects.registry.vmware.com/tkg
      path: csi/vsphere-block-csi-driver
      tag: v2.1.0_vmware.1
      pullPolicy: IfNotPresent
    livenessProbeImage:
      repository: projects.registry.vmware.com/tkg
      path: csi/csi-livenessprobe
      tag: v2.1.0_vmware.1
      pullPolicy: IfNotPresent
    vsphereSyncerImage:
      repository: projects.registry.vmware.com/tkg
      path: csi/volume-metadata-syncer
      tag: v2.1.0_vmware.1
      pullPolicy: IfNotPresent
    CSIProvisionerImage:
      repository: projects.registry.vmware.com/tkg
      path: csi/csi-provisioner
      tag: v2.0.0_vmware.1
      pullPolicy: IfNotPresent
    CSINodeDriverRegistrarImage:
      repository: projects.registry.vmware.com/tkg
      path: csi/csi-node-driver-registrar
      tag: v2.0.1_vmware.1
      pullPolicy: IfNotPresent
    namespace: kube-system
    clusterName: wc-1
    server: 10.170.104.114
    datacenter: /dc0
    publicNetwork: VM Network
    username: <MY-VSPHERE-USERNAME>
    password: <MY-VSPHERE-PASSWORD>
 
  ```

1. Open `vsphere-csi-secret.yaml` in an editor and do the following to make it look like the code below:
  1. After the first line, add two lines that define `stringData`, and `values.yaml` as its first element.
  1. Copy the `data.values` output from the previous step.
  1. After the third line, paste in the `data.values` output and indent it as the value of `values.yaml`.
  1. Immediately below the `values.yaml` definition, add another `stringData` definition for `overlays.yaml` as shown below. Do not modify other definitions in the file.

      ```
      apiVersion: v1
      stringData:
        values.yaml: |
          #@data/values
          #@overlay/match-child-defaults missing_ok=True
          ---
          vsphereCSI:
            CSIAttacherImage:
              repository: projects.registry.vmware.com/tkg
              path: csi/csi-attacher
              tag: v3.0.0_vmware.1
              pullPolicy: IfNotPresent
            vsphereCSIControllerImage:
              repository: projects.registry.vmware.com/tkg
              path: csi/vsphere-block-csi-driver
              tag: v2.1.0_vmware.1
              pullPolicy: IfNotPresent
            livenessProbeImage:
              repository: projects.registry.vmware.com/tkg
              path: csi/csi-livenessprobe
              tag: v2.1.0_vmware.1
              pullPolicy: IfNotPresent
            vsphereSyncerImage:
              repository: projects.registry.vmware.com/tkg
              path: csi/volume-metadata-syncer
              tag: v2.1.0_vmware.1
              pullPolicy: IfNotPresent
            CSIProvisionerImage:
              repository: projects.registry.vmware.com/tkg
              path: csi/csi-provisioner
              tag: v2.0.0_vmware.1
              pullPolicy: IfNotPresent
            CSINodeDriverRegistrarImage:
              repository: projects.registry.vmware.com/tkg
              path: csi/csi-node-driver-registrar
              tag: v2.0.1_vmware.1
              pullPolicy: IfNotPresent
            namespace: kube-system
            clusterName: wc-1
            server: 10.170.104.114
            datacenter: /dc0
            publicNetwork: VM Network
            username: <MY-VSPHERE-USERNAME>
            password: <MY-VSPHERE-PASSWORD>
        overlays.yaml: |
          #@ load("@ytt:overlay", "overlay")    
          #@overlay/match by=overlay.subset({"kind": "Deployment", "metadata": {"name": "vsphere-csi-controller"}})
          ---
          spec:
            template:
              spec:
                containers:
                #@overlay/append
                  - name: csi-resizer
                    image: projects.registry.vmware.com/tkg/kubernetes-csi_external-resizer:v1.0.0_vmware.1
                    args:
                      - "--v=4"
                      - "--timeout=300s"
                      - "--csi-address=$(ADDRESS)"
                      - "--leader-election"
                    env:
                      - name: ADDRESS
                        value: /csi/csi.sock
                    volumeMounts:
                      - mountPath: /csi
                        name: socket-dir
      kind: Secret
      ...
      ```

1. Run `kubectl apply` to update the cluster's secret with the revised definitions and re-create the `csi-controller` pod:

  ```
  kubectl apply -f vsphere-csi-secret.yaml 
  ```

1. To verify that the `vsphere-csi-controller` and external resizer are working on the cluster:
  1. Confirm that `vsphere-csi-controller` is running on the workload cluster with six healthy pods:

      ```
      $ kubectl get pods -n kube-system -l app=vsphere-csi-controller
      NAME                                     READY   STATUS    RESTARTS   AGE
      vsphere-csi-controller-<ID-HASH>   6/6     Running   0          6m49s
      ```

  1. Check the logs of the `vsphere-csi-controller` to see that the external resizer started.

      ```
      $ kubectl logs vsphere-csi-controller-<ID-HASH> -n kube-system -c csi-resizer
        I0308 23:44:45.035254       1 main.go:79] Version : v1.0.0-0-gb22717d
        I0308 23:44:45.037856       1 connection.go:153] Connecting to unix:///csi/csi.sock
        I0308 23:44:45.038572       1 common.go:111] Probing CSI driver for readiness
        I0308 23:44:45.040579       1 csi_resizer.go:77] CSI driver name: "csi.vsphere.vmware.com"
        W0308 23:44:45.040612       1 metrics.go:142] metrics endpoint will not be started because `metrics-address` was not specified.
        I0308 23:44:45.042184       1 controller.go:117] Register Pod informer for resizer csi.vsphere.vmware.com
        I0308 23:44:45.043182       1 leaderelection.go:243] attempting to acquire leader lease  kube-system/external-resizer-csi-vsphere-vmware-com...
        I0308 23:44:45.073383       1 leaderelection.go:253] successfully acquired lease kube-system/external-resizer-csi-vsphere-vmware-com
        I0308 23:44:45.076028       1 leader_election.go:172] new leader detected, current leader: vsphere-csi-controller-87d7dcf48-jcht2
        I0308 23:44:45.079332       1 leader_election.go:165] became leader, starting
        I0308 23:44:45.079638       1 controller.go:241] Starting external resizer csi.vsphere.vmware.com 
      ```
