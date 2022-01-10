# Sriov-Network-Device-Plugin Package

This package contains Kubernetes device plugin for discovering and advertising SR-IOV virtual functions (VFs) available on a Kubernetes host. Using [sriov-network-device-plugin](https://github.com/k8snetworkplumbingwg/sriov-netowrk-device-plugin/).

## Components

* SRIOV Network Device Plugin DaemonSet
* SRIOV Network Device Plugin ConfigMap
* SRIOV Network Device Plugin ServiceAccount

## Configuration

The following configuration values can be set to customize the sriov-network-device-plugin installation.

### Global

| Value       | Required/Optional | Description                                                   |
| ----------- | ----------------- | ------------------------------------------------------------- |
| `namespace` | Optional          | The namespace in which to deploy sriov-network-device-plugin. |

### sriov-network-device-plugin DaemonSet Configuration

| Value                   | Required/Optional | Description                                                                                                     |
| ----------------------- | ----------------- | --------------------------------------------------------------------------------------------------------------- |
| `daemonset.resources`             | Optional          | The requests and limits for cpu and memory resources daemonset requires. |
| `daemonset.args`             | Required          | The arguments passed into daemonset's container. |
| `sriov_nodes_resources` | Required          | The SR-IOV resources list that daemonset of each node pool refers to.  |

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

## Example for sriov_nodes_resources

The `sriov_nodes_resources` field in values file should be like:

````yaml
sriov_nodes_resources:
    key:value1: |
        {
            "resourceList": [{
                    "resourceName": "intel_sriov_netdevice",
                    "selectors": {
                        "vendors": ["8086"],
                        "devices": ["154c", "10ed"],
                        "drivers": ["i40evf", "iavf", "ixgbevf"]
                    }
                }
            ]
        }
    key:value2: |
        {
            "resourceList": [{
                    "resourceName": "mlnx_sriov_rdma",
                    "selectors": {
                        "vendors": ["15b3"],
                        "devices": ["1018"],
                        "drivers": ["mlx5_ib"],
                        "isRdma": true
                    }
                }
            ]
        }
        ...
````

Each entry in the list refers to one node pool with a same set of SR-IOV resources list.

1. `key:values` refers to the node label used to specify the SR-IOV node pools on each node.

1. The value field refers to the available SR-IOV resources on that node pool.
