# Enable Volume Expansion for vSphere CSI (vSphere 7)

To enable [volume expansion](https://vsphere-csi-driver.sigs.k8s.io/features/volume_expansion.html) for vSphere CSI storage used by workload clusters, you need to add a `csi-resizer` sidecar pod to the cluster's CSI processes.

The CSI configuration for workload clusters is encoded as a Kubernetes secret.
This procedure adds the `csi-resizer` process by revising the CSI configuration secret.
It adds to the secret a `stringData` definition that combines two encoded configuration data strings: a `values.yaml` string containing the secret's prior CSI configuration data, and a new `overlays.yaml` string that deploys the `csi-resizer` pod.

NOTE: Online volume expansion is supported in vSphere 7.0 as of Update 2; see [Volume Expansion in vSphere with Tanzu](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-8D7C8AAA-BD59-49EF-AB02-C1B2FF46F59B.html).

1. Log into the management cluster for the workload cluster you are changing, and run `tanzu cluster list` if you need to retrieve the name of the workload cluster.

1. Retrieve the name of the CSI secret for the workload cluster, using label selectors `vsphere-csi` and the cluster name:

    ```sh
    $ kubectl get secret \
        -l tkg.tanzu.vmware.com/cluster-name=NAME_OF_WORKLOAD_CLUSTER \
        -l tkg.tanzu.vmware.com/addon-name=vsphere-csi
    my-wc-vsphere-csi-secret
    ```

1. Save a backup of the secret's content, in YAML format, to `vsphere-csi-secret.yaml`:

   ```sh
   kubectl get secret my-wc-vsphere-csi-secret -o yaml > vsphere-csi-secret.yaml
   ```

1. Output the secret's current content again, with the `data.values` values `base64`-decoded into plain YAML.

    ```sh
    kubectl get secret my-wc-vsphere-csi-secret -o jsonpath={.data.values\\.yaml} | base64 -d

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
1. Delete the existing definition for `values.yaml`, which is a long string.
1. After the first line, add a line that defines `stringData` and indent `values.yaml` to make it the first element.
1. Copy the `data.values` output from the previous step.
1. After the third line, paste in the `data.values` output and indent it as the value of `values.yaml`.
1. Immediately below the `values.yaml` definition, add another `stringData` definition for `overlays.yaml` as shown below. Do not modify other definitions in the file.

    ```sh
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

    ```sh
    kubectl apply -f vsphere-csi-secret.yaml
    ```

1. To verify that the `vsphere-csi-controller` and external resizer are working on the cluster:

    1. Confirm that `vsphere-csi-controller` is running on the workload cluster with six healthy pods:

        ```sh
        $ kubectl get pods -n kube-system -l app=vsphere-csi-controller
        NAME                                     READY   STATUS    RESTARTS   AGE
        vsphere-csi-controller-<ID-HASH>   6/6     Running   0          6m49s
        ```

    1. Check the logs of the `vsphere-csi-controller` to see that the external resizer started.

        ```sh
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

For more information about expanding vSphere CSI storage volumes in online or offline mode, see [Expand a Persistent Volume in Online Mode](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-98F6BEB2-7A0A-4C05-B687-BD329C5D2E32.html) and [Expand a Persistent Volume in Offline Mode](https://docs.vmware.com/en/VMware-vSphere/7.0/vmware-vsphere-with-tanzu/GUID-90082E1C-DC01-4610-ABA2-6A4E97C18CBC.html).
