# whereabouts Package

This package provides the ability to assign IP addresses dynamically across your Kubernetes cluster using a CNI IPAM plugin named [whereabouts](https://github.com/k8snetworkplumbingwg/whereabouts).

## Components

* Whereabouts Custom Resources
* Whereabouts DaemonSet
* Whereabouts ServiceAccount
* Whereabouts ClusterRoleBinding

## Configuration

The following configuration values can be set to customize the whereabouts installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy whereabouts components. Default: kube-system |

### whereabouts Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `whereabouts.config.resources.limits.cpu` | Optional | The limits for CPU resources of whereabouts DeamonSet  |
| `whereabouts.config.resources.limits.memory` | Optional | The limits for memory resources of whereabouts DeamonSet  |
| `whereabouts.config.resources.requests.cpu` | Optional | The requests for CPU resources of whereabouts DeamonSet  |
| `whereabouts.config.resources.requests.memory` | Optional | The requests for memory resources of whereabouts DeamonSet  |

## Usage

The follow is a basic guide for getting started with whereabouts.

This example guides you about attaching the second network interface to a pod with IP address assigned in the range you specified using whereabouts.

1. Install TCE Multus CNI package to support multiple network by following
   [doc](https://github.com/vmware-tanzu/community-edition/blob/main/addons/packages/multus-cni/3.7.1/README.md#usage-example):

1. Install TCE whereabouts package through Tanzu CLI

    ```bash
    tanzu package install whereabouts --package-name whereabouts.community.tanzu.vmware.com --version ${MULTUS_PACKAGE_VERSION}
    ```

    > You can get the `${MULTUS_PACKAGE_VERSION}` from running `tanzu package
    > available list whereabouts.community.tanzu.vmware.com`. Specifying a
    > namespace may be required depending on where your package repository was
    > installed.

1. After the Multus CNI and whereabouts DaemonSet are running, you can define your NetworkAttachmentDefinition to tell
   * which CNI plugin will be used for the second network interface, in particular this example uses `ipvlan` CNI plugin
   * what IP addressed will be assigned for the second network interface,  in particular this example uses `whereabouts` CNI IPAM plugin

   ```bash
   cat <<EOF | kubectl create -f -
    apiVersion: "k8s.cni.cncf.io/v1"
    kind: NetworkAttachmentDefinition
    metadata:
    name: ipvlan-conf-1
    spec:
    config: '{
        "cniVersion": "0.3.0",
        "name": "ipvlan-conf-1",
        "type": "ipvlan",
        "master": "eth0",
        "mode": "bridge",
        "ipam": {
            "type": "whereabouts",
            "range": "192.168.20.0/24",
            "gateway": "192.168.20.1",
            "range_start": "192.168.20.2",
            "range_end": "192.168.20.100"
        }
        }'
    EOF
    ```

1. Deploy a pod using the NetworkAttachmentDefinition named `ipvlan-conf-1` as above by adding following lines to the pod `metadata.annotations`:

    ```bash
    metadata:
      annotations:
        k8s.v1.cni.cncf.io/networks: ipvlan-conf-1
    ```

1. After the pod is running, run following command to describe your pod and there will be an event for adding the second network interface within IP range
   we specified with whereabouts.

   ```bash
   kubectl describe {your-pod}
   ... ...
   Events:
    Type    Reason          Age   From               Message
    ----    ------          ----  ----               -------
    Normal  AddedInterface  2m1s  multus             Add eth0 [100.96.1.6/24]
    Normal  AddedInterface  2m1s  multus             Add net1 [192.168.20.10/24] from default/ipvlan-conf-1
   ```

    You can also run following command to check more details about the second network interface:

    ```bash
    kubectl exec {your-pod} -- ip a
    ```
