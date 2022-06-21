# Whereabouts Package

This package provides the ability to assign IP addresses dynamically across your Kubernetes cluster using a CNI IPAM plugin named [whereabouts](https://github.com/k8snetworkplumbingwg/whereabouts).

## Installation

To support mutiple network, the Multus CNI package should be installed along with the Whereabouts package.

### Installation of dependencies

Install the Multus CNI through tanzu command:

```bash
tanzu package install multus-cni --package-name multus-cni.community.tanzu.vmware.com --version ${MULTUS_PACKAGE_VERSION}
```

> You can get the `${MULTUS_PACKAGE_VERSION}` from running `tanzu package
> available list multus-cni.community.tanzu.vmware.com`. Specifying a
> namespace may be required depending on where your package repository was
> installed.

### Installation of package

Install TCE Whereabouts package through tanzu command:

```bash
tanzu package install whereabouts --package-name whereabouts.community.tanzu.vmware.com --version ${WHEREABOUTS_PACKAGE_VERSION}
```

> You can get the `${WHEREABOUTS_PACKAGE_VERSION}` from running `tanzu package
> available list whereabouts.community.tanzu.vmware.com`. Specifying a
> namespace may be required depending on where your package repository was
> installed.

## Options

The following configuration values can be set to customize the Whereabouts installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy Whereabouts components. Default: kube-system |

### Whereabouts Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `whereabouts.config.resources.limits.cpu` | Optional | The limits for CPU resources of Whereabouts DaemonSet  |
| `whereabouts.config.resources.limits.memory` | Optional | The limits for memory resources of Whereabouts DaemonSet  |
| `whereabouts.config.resources.requests.cpu` | Optional | The requests for CPU resources of Whereabouts DaemonSet  |
| `whereabouts.config.resources.requests.memory` | Optional | The requests for memory resources of Whereabouts DaemonSet  |
| `ip_reconciler.config.schedule` | Optional | The schedule of ip-reconciler CronJob. Default: \*/5 \* \* \* \*  |
| `ip_reconciler.config.resources.requests.cpu` | Optional | The requests for memory resources of ip-reconciler CronJob  |
| `ip_reconciler.config.resources.requests.memory` | Optional | The requests for memory resources of ip-reconciler CronJob  |

## What This Package Does

If you need a way to assign IP addresses dynamically across your cluster -- Whereabouts is the tool for you. If you've found that you like how the host-local CNI plugin works, but, you need something that works across all the nodes in your cluster (host-local only knows how to assign IPs to pods on the same node) -- Whereabouts is just what you're looking for.

Whereabouts can be used for both IPv4 & IPv6 addressing.

## Components

* Whereabouts Custom Resources
* Whereabouts DaemonSet
* Whereabouts ServiceAccount
* Whereabouts ClusterRoleBinding
* Whereabouts ip-reconciler Cronjob

## Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Files

Here is an example of the package configuration file [values.yaml](bundle/config/values.yaml).

## Package Limitations

To use the Whereabouts Package, the Multus Package should be installed so that they can work together to configure the secondary network interface of the pods.

The primary network interface of the pods is managed by the CNI Packages like Antrea or Calico. Their IPAM tools are already defined in their packages and not exposed as a configurable option so far.

## Usage Example

The follow is a basic guide for getting started with Whereabouts.

It shows you how to attach a second network interface to a pod with an IP address assigned in the range you specified using Whereabouts.

1. Install the TCE Multus CNI package along with the TCE Whereabouts package by following the documentation to support multiple network.

1. After the Multus CNI and Whereabouts DaemonSets are running, you can define your NetworkAttachmentDefinition to tell
   * which CNI plugin will be used for the second network interface, this example uses the `ipvlan` CNI plugin
   * which IP address will be assigned for the second network interface, this example uses the `whereabouts` CNI IPAM plugin

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

1. Deploy a pod using the NetworkAttachmentDefinition named `ipvlan-conf-1` above, by adding the following lines to the pod `metadata.annotations`:

    ```bash
    metadata:
      annotations:
        k8s.v1.cni.cncf.io/networks: ipvlan-conf-1
    ```

1. After the pod is running, run the following command to describe your pod. There will be an event for adding the second network interface within the IP range we specified with Whereabouts:

   ```bash
   kubectl describe {your-pod}
   ... ...
   Events:
    Type    Reason          Age   From               Message
    ----    ------          ----  ----               -------
    Normal  AddedInterface  2m1s  multus             Add eth0 [100.96.1.6/24]
    Normal  AddedInterface  2m1s  multus             Add net1 [192.168.20.10/24] from default/ipvlan-conf-1
   ```

    You can also run the following command to check more details about the second network interface:

    ```bash
    kubectl exec {your-pod} -- ip a
    ```
