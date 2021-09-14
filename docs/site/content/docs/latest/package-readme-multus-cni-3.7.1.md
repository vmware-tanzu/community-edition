# Multus CNI

This package provides functionality that enables attaching multiple network interfaces to pods in Kubernetes using [multus-cni](https://github.com/k8snetworkplumbingwg/multus-cni).

## Components

* Multus CNI Custom Resources
* Multus CNI DaemonSet
* Multus CNI ConfigMap

## Configuration

The following configuration values can be set to customize the Multus CNI installation.

### Global

|Value |Required/Optional |Description |
|:-------|:-------------------|:-------------|
|`namespace` |Optional |The namespace in which to deploy Multus CNI DaemonSet. |

### Multus CNI configuration

|Value |Required/Optional |Description |
|:-------|:-------------------|:-------------|
|`image` |Optional |The image used by Multus CNI DaemonSet. |
|`args` |Optional |The args for Multus CNI DaemonSet container. |

## Usage Example

This example provides the steps for attaching a second network interface that leverages the Multus CNI package.

### Before You Begin
You must deploy the Multus CNI package.

### Procedure

1. Install the Multus CNI through the tanzu command:

    ```bash
    tanzu package install multus-cni --package-name multus-cni.community.tanzu.vmware.com --version ${MULTUS_PACKAGE_VERSION}
    ```

    > You can get the `${MULTUS_PACKAGE_VERSION}` from running `tanzu package
    > available list multus-cni.community.tanzu.vmware.com`. Specifying a
    > namespace may be required depending on where your package repository was
    > installed.

2. After the Multus CNI DaemonSet is running, you can define your network-attachment-defs to tell Multus CNI which CNI will be used for other network interfaces:

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

3. Deploy a sample pod using the network-attachment-defs defined above by adding the following lines to the pod spec:

    ```bash
    metadata:
      annotations:
      k8s.v1.cni.cncf.io/networks: macvlan-conf
    ```

4. After the pod is running, run the following command to check if the second network interface is up and running:

    ```bash
    kubectl exec <your-pod> -- ip a
    ```
