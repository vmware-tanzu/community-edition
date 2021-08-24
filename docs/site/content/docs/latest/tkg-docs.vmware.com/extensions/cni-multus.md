# Implementing Multiple CNIs with Multus

[Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni) is a Container Network Interface ([CNI](https://www.cni.dev)) plugin for Kubernetes that lets pods run multiple other CNI plug-ins, each of which is associated with a different address range.
Multus acts as a "meta-plugin" to let different workloads use different network interfaces.

This topic explains how to install the Multus package onto a Tanzu Kubernetes (workload) cluster and use it to create pods with multiple network interfaces.
For example, Antrea or Calico as the primary CNI, and a secondary interface such as [macvlan](https://www.cni.dev/plugins/current/main/macvlan/) or [ipvlan](https://www.cni.dev/plugins/current/main/ipvlan/), or [SR-IOV](https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.networking.doc/GUID-CC021803-30EA-444D-BCBE-618E0D836B9F.html) or [DPDK](https://www.dpdk.org/) devices for hardware or accelerated interfaces.

Binaries for macvlan and ipvlan are already installed in the workload cluster node template.

## <a id="prereqs"></a> Prerequisites

- A bootstrap machine with the Tanzu CLI, Tanzu CLI plugins, and `kubectl` installed as described in [Install the Tanzu CLI and Other Tools](../install-cli.md) and plugins the `packages` plugin
- A Tanzu Kubernetes Grid management cluster and workload cluster running on vSphere, Amazon EC2, or Azure.

## <a id="install"></a> Install the Multus CNI Package

**NOTE:** Once the Multus CNI is installed in a cluster, it should not be deleted.
See [Deleting Multus Unsupported](#delete) below.

To install the Multus CNI package on a workload cluster and configure the cluster to use it:

1. Install the Multus package:

   * **Default Configuration**: Run `tanzu package install multus-cni`.

   * **Custom Configuration**:

        1. Create a configuration file `multus-cni-values.yaml` that retrieves the Multus image from Docker and deploys it as a Daemonset.
        See the Multus [`entrypoint.sh` script](https://github.com/k8snetworkplumbingwg/multus-cni/blob/master/images/entrypoint.sh#L50) for more information.<br />
        For example:

          ```
          #@data/values
          #@overlay/match-child-defaults missing_ok=True
          ---

          namespace: kube-system

          image:
            repository: docker.io/nfvpe
            name: multus
            tag: stable

          #! DaemonSet related configuration
          #@overlay/replace
          daemonset:
            #! please refer to https://github.com/k8snetworkplumbingwg/multus-cni/blob/master/images/entrypoint.sh#L50
            args:
              - "--multus-conf-file=auto"
              - "--cni-version=0.3.1"
          ```

        1. Run `tanzu package install multus-cni --config multus-cni-values.yaml`

1. Create a custom resource definition (CRD) for `NetworkAttachmentDefinition` that defines the CNI configuration for network interfaces to be used by Multus CNI.

   1. Create a CRD specification. For example, this `multus-cni-crd.yaml` specifies a `NetworkAttachmentDefinition` named `macvlan-conf` that configures a `macvlan` CNI:

      ```
      apiVersion: "k8s.cni.cncf.io/v1"
      kind: NetworkAttachmentDefinition
      metadata:
        name: macvlan-conf
      spec:
        config: '{
            "cniVersion": "0.3.0",
            "type": "macvlan",
            "master": "eth0",
            "mode": "bridge",
            "ipam": {
              "type": "host-local",
              "subnet": "192.168.1.0/24",
              "rangeStart": "192.168.1.200",
              "rangeEnd": "192.168.1.216",
              "routes": [
                { "dst": "0.0.0.0/0" }
              ],
              "gateway": "192.168.1.1"
            }
          }'
      ```

    1. Create the resource; for example `kubectl create -f multus-cni-crd.yaml`

1. Create a pod with the annotation `k8s.v1.cni.cncf.io/networks`, which takes a comma-delimited list of the names of `NetworkAttachmentDefinition` custom resource.

   1. Create the pod specification, for example `my-multi-cni-pod.yaml`:

      ```
      apiVersion: v1
      kind: Pod
      metadata:
        name: sample-pod
        annotations:
          k8s.v1.cni.cncf.io/networks: macvlan-conf
      spec:
        containers:
        - name: sample-pod
          command: ["/bin/ash", "-c", "trap : TERM INT; sleep infinity & wait"]
          image: harbor-repo.vmware.com/dockerhub-proxy-cache/library/alpine

      ```

    1. Create the pod; for example `kubectl create -f my-multi-cni-crd.yaml` creates the pod `sample-pod`.

Once the pod is created, it will have three network interfaces:

- `lo` the loopback interface
- `eth0` the default pod network managed by Antrea or Calico CNI
- `net1` the new interface created via the annotation `k8s.v1.cni.cncf.io/networks: macvlan-conf`.

Note: The default network gets the name `eth0` and additional network pod interfaces get the name as `net1`, `net2`, and so on.

## <a id="validate"></a>Validating Multus

Run `kubectl describe pod` on the pod, and confirm that the annotation `k8s.v1.cni.cncf.io/network-status` lists all network interfaces.
For example:

```
$ kubectl describe pod sample-pod

Name:         sample-pod
Namespace:    default
Priority:     0
Node:         tcecluster-md-0-6476897f75-rl9vt/10.170.109.225
Start Time:   Thu, 27 May 2021 15:31:20 +0000
Labels:       <none>
Annotations:  k8s.v1.cni.cncf.io/network-status:
                [{
                    "name": "",
                    "interface": "eth0",
                    "ips": [
                        "100.96.1.80"
                    ],
                    "mac": "66:39:dc:63:50:a3",
                    "default": true,
                    "dns": {}
                },{
                    "name": "default/macvlan-conf",
                    "interface": "net1",
                    "ips": [
                        "192.168.1.201"
                    ],
                    "mac": "02:77:cb:a0:60:e3",
                    "dns": {}
                }]
              k8s.v1.cni.cncf.io/networks: macvlan-conf

```

## <a id="delete"></a>Deleting Multus Unsupported

Once the Multus CNI is installed in a cluster, it should not be deleted.

Deleting Multus does not uninstall the the Multus configuration file `/etc/cni/net.d/00-multus.conf` from the CNI scripts directory, which prevents the cluster from creating new pods.

This is a known issue; see [Issue #461](https://github.com/k8snetworkplumbingwg/multus-cni/issues/461) in the Multus repository.
