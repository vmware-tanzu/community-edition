# Multus CNI

This package provides the ability for enabling attaching multiple network interfaces to pods in Kubernetes using [multus-cni](https://github.com/k8snetworkplumbingwg/multus-cni).

## Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Components

* Multus CNI Custom Resources
* Multus CNI DaemonSet
* Multus CNI ConfigMap

## Configuration

The following configuration values can be set to customize the Multus CNI installation.

### Global

| Value       | Required/Optional | Description                                            |
| ----------- | ----------------- | ------------------------------------------------------ |
| `namespace` | Optional          | The namespace in which to deploy Multus CNI DaemonSet. |

### Multus CNI configuration

| Value  | Required/Optional | Description                                  |
| ------ | ----------------- | -------------------------------------------- |
| `args` | Optional          | The args for Multus CNI DaemonSet container. |

## Usage Example

This example guides you about attaching another network interface scenario that leverages the Multus CNI package. You must deploy the package before attempting this walkthrough.

1. Install the Multus CNI through tanzu command:

    ```bash
    tanzu package install multus-cni --package-name multus-cni.community.tanzu.vmware.com --version ${MULTUS_PACKAGE_VERSION}
    ```

    > You can get the `${MULTUS_PACKAGE_VERSION}` from running `tanzu package
    > available list multus-cni.community.tanzu.vmware.com`. Specifying a
    > namespace may be required depending on where your package repository was
    > installed.

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

1. Deploy a sample pod using the network-attachment-defs defined above by the following  lines to the pod spec:

    ```bash
    metadata:
      annotations:
      k8s.v1.cni.cncf.io/networks: macvlan-conf
    ```

1. After the pod is running, run the following command to check if the second network interface is up and running:

    ```bash
    kubectl exec <your-pod> -- ip a
    ```

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
