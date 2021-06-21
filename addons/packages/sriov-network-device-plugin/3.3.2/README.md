# Sriov-Network-Device-Plugin Package

This package contains Kubernetes device plugin for discovering and advertising SR-IOV virtual functions (VFs) available on a Kubernetes host. Using [sriov-network-device-plugin](https://github.com/k8snetworkplumbingwg/sriov-netowrk-device-plugin/).

## Components

* SRIOV Network Device Plugin DaemonSet

## Configuration

The following configuration values can be set to customize the sriov-network-device-plugin installation.

### Global

| Value       | Required/Optional | Description                                                   |
| ----------- | ----------------- | ------------------------------------------------------------- |
| `namespace` | Required          | The namespace in which to deploy sriov-network-device-plugin. |

### sriov-network-device-plugin DaemonSet Configuration

| Value                   | Required/Optional | Description                                                                                                     |
| ----------------------- | ----------------- | --------------------------------------------------------------------------------------------------------------- |
| `image`                 | Required          | The image used to start daemonset. Must provided with repository, name and tag.                                 |
| `imagePullPolicy`       | Required          | The image pull policy in the daemonset. Must provided with image field.                                         |
| `daemonset`             | Required          | The resources daemonset settings including resources. |
| `sriov_nodes_resources` | Required          | The SR-IOV resource list sriov-network-device-plugin daemonset refers to.  |

## Usage Example

The follow is a basic guide for getting started with sriov-network-device-plugin.

1. Firstly, install the SRIOV Network Device Plugin through tanzu command with a values file containing the SR-IOV devices information:

    ```bash
    tanzu package install sriov-network-device-plugin.community.tanzu.vmware.com --version <3.3.2> -f <your-values-yaml-file>
    ```

1. Check the config lists used for sriov-network-device-plugin DaemonSet by:

    ```bash
    tanzu package available get sriov-network-device-plugin.community.tanzu.vmware.com --values-schema -o json
    ```
