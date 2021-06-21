# sriov-cni Package

This package enables the configuration and usage of SR-IOV VF networks in containers and orchestrators like Kubernetes. using [sriov-cni](https://github.com/k8snetworkplumbingwg/sriov-cni/).

## Components

* SRIOV CNI DaemonSet

## Configuration

The following configuration values can be set to customize the sriov-cni installation.

### Global

| Value       | Required/Optional | Description                                 |
| ----------- | ----------------- | ------------------------------------------- |
| `namespace` | Optional          | The namespace in which to deploy sriov-cni. |

### sriov-cni Configuration

| Value                 | Required/Optional | Description                                                                                                       |
| --------------------- | ----------------- | ----------------------------------------------------------------------------------------------------------------- |
| `daemonset.resources` | Optional          | The resources requirement for SR-IOV CNI image, such as limits and requests, should be provided in complete form. |
| `daemonset.args`      | Optional          | The args passed to image entrypoint.                                                                              |

## Usage Example

The follow is a basic guide for getting started with sriov-cni. To enable SRIOV CNI, Multus CNI and SRIOV Network Device Plugin are needed. So please install Multus CNI package first before install sriov-cni package. Otherwise sriov-cni won't work as expected.

1. Firstly, install the Multus CNI through tanzu command:

    ```bash
    tanzu package install sriov-cni.community.tanzu.vmware.com
    ```

1. Or you can specify the args used for sriov-cni DaemonSet by:

    ```bash
    tanzu package available get sriov-cni.community.tanzu.vmware.com --vaules-schema -o json
    ```

1. Add a values.yaml file contains the above values schema in your current directory. You can change the image in it and then pass it to package install process by:

    ```bash
    tanzu package install sriov-cni.community.tanzu.vmware.com --version <2.6.1> -f <your-values-yaml-file>
    ```
