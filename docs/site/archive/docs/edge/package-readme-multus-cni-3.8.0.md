# Multus CNI

This package enables you to attach multiple network interfaces to pods in Kubernetes using [multus-cni](https://github.com/k8snetworkplumbingwg/multus-cni).

This documentation provides information about the specific TCE package. Please visit the [TCE package management page](https://tanzucommunityedition.io/docs/v0.11/package-management/) for general information about installation, removal, troubleshooting, and other topics.

## Installation

### Installation of dependencies

No dependencies required to install the Multus package.

### Installation of package

Install the Multus CNI through tanzu command:

```bash
tanzu package install multus-cni --package-name multus-cni.community.tanzu.vmware.com --version ${MULTUS_PACKAGE_VERSION}
```

> You can get the `${MULTUS_PACKAGE_VERSION}` from running `tanzu package
> available list multus-cni.community.tanzu.vmware.com`. Specifying a
> namespace may be required depending on where your package repository was
> installed.

## Uninstallation of Multus CNI package

The following steps are used to uninstall the Multus CNI package.

1. Delete the Multus CNI resources through the following command:

    ```bash
    tanzu package installed delete <multus-cni-pkg-install-name> <-y>
    ```

1. Remove leftover Multus CNI's network configuration files on the cluster nodes. To remove such resources, one possible way is to use a daemonset to clean up the leftover Multus CNI related network configurations. One example is located at Multus CNI tests folder.

    ```bash
    kubectl create -f test/e2e/multihomed-testfiles/cleanup.yaml
    ```

## Options

You can set following configuration values to customize the Multus CNI installation.

### Package configuration values

#### Global

| Value       | Required/Optional | Default | Description                                            |
| ----------- | ----------------- | ------- | ------------------------------------------------------ |
| `namespace` | Optional          | `kube-system` | The namespace in which to deploy Multus CNI DaemonSet. |

#### Multus CNI configuration

| Value  | Required/Optional | Default | Description                                 |
| ------ | ----------------- | ------- | ------------------------------------------- |
| `args` | Optional          | `- "--multus-conf-file=auto"`, `- "--cni-version=0.3.1"` | The args for Multus CNI DaemonSet container. |

### Application configuration values

No available options to configure.

#### Multi-cloud configuration steps

There are currently no configuration steps necessary for installation of the Multus package to any provider.

## Components

* Multus CNI Custom Resources
* Multus CNI DaemonSet
* Multus CNI ConfigMap

## What This Package Does

Multus CNI enables attaching multiple network interfaces to pods in Kubernetes.

## Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Files

Here is an example of the package configuration file [values.yaml](bundle/config/values.yaml).

## Package Limitations

Uninstallation of Multus CNI package requires extra cleanup process to remove network configuration files on the cluster nodes. See the uninstallation guide above.

## Usage Example

This example guides you through attaching another network interface scenario that leverages the Multus CNI package. You must install the package before attempting this walkthrough.

1. After the Multus CNI DaemonSet is running, you can define your network-attachment-defs to tell Multus CNI which CNI will be used for other network interfaces:

   ```bash
   cat <<EOF | kubectl create -f -
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
    EOF
    ```

1. Deploy a sample pod using the network-attachment-defs defined above. Refer to the following lines to the pod spec:

    ```bash
    metadata:
      annotations:
      k8s.v1.cni.cncf.io/networks: macvlan-conf
    ```

1. After the pod is running, run the following command to check if the second network interface is also running:

    ```bash
    kubectl exec <your-pod> -- ip a
    ```

## Troubleshooting

Not applicable.

## Additional Documentation

See the [Multus documentation](https://github.com/k8snetworkplumbingwg/multus-cni) for more information.
