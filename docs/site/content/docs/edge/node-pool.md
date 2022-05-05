# Manage Node Pools for Different VM Types

This topic explains how to create, update and delete node pools in a workload cluster.
Node pools enable a single workload cluster to contain and manage different types of nodes, to support the diverse needs of different applications.

For example, a workload cluster can use nodes with high storage capacity to run a datastore, and thinner nodes to process application requests.

## About Node Pools

Node pools define properties for the sets of worker nodes used by a workload cluster.

Some node pool properties depend on the VM options that are available in the underlying infrastructure, but all node pools on all cloud infrastructures share the following properties:

- `name`: a unique identifier for the node pool, used for operations like updates and deletion.
- `replicas`: the number of nodes in the pool, all of which share the same properties.
- `labels`: key/value pairs set as metadata on the nodes, to match workloads to nodes in the pool. For more information, and example labels, see [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#motivation) in the Kubernetes documentation.

All workload clusters are created with a first, original node pool.
When you create additional node pools for a cluster, as described below, the first node pool provides default values for properties not set in the new node pool definitions.

## List Node Pools

To inspect the node pools currently available in a cluster, run:

```shell
tanzu cluster node-pool list CLUSTER-NAME
```

This returns a lists of all of the node pools in the cluster, and the state of the replicas in each node pool.

## Create a Node Pool

To create a node pool in a cluster:

1. Create a configuration file for the node pool.  
    See below for example configuration files for each infrastructure provider.

2. Create the node pool defined by the configuration file:

    ```shell
    tanzu cluster node-pool set CLUSTER-NAME -f <CONFIG-FILE>
    ```

    Options:

    - `--namespace` specifies the namespace of the cluster. The default value is `default`.
    - `--machine-deployment-base` specifies the base `MachineDeployment` object from which to create the new node pool.
      - Set this value as a `MachineDeployment` identifier as listed in the output of `tanzu cluster get` under `Details`.
      - The default value is the first in the cluster's array of worker node `MachineDeployment` objects, represented internally as `workerMDs[0]`.

### AWS Configuration

In addition to the required `name`, `replicas`, and `labels` properties above, configuration files for node pools on AWS support the following optional properties:

- `az`: Availability Zone
- `nodeMachineType`: Instance type

   These settings may be omitted, in which case their values inherit from the cluster's first node pool.

Example node pool definition for an AWS cluster:

```yaml
name: tkg-aws-wc-np-1
replicas: 2
az: us-west-2b
nodeMachineType: t3.large
labels:
key1: value1
key2: value2
```

### Azure Configuration

In addition to the required `name`, `replicas`, and `labels` properties above, configuration files for node pools on Microsoft Azure support the following optional properties:

- `az`: Availability Zone
- `nodeMachineType`: Instance type

These settings may be omitted, in which case their values inherit from the cluster's first node pool.

Example node pool definition for Azure cluster:

```yaml
name: tkg-azure-wc-np-1
replicas: 2
az: 2
nodeMachineType: Standard_D2s_v3
labels:
key1: value1
key2: value2
```

### vSphere Configuration

In addition to the required `name`, `replicas`, and `labels` properties above, configuration files for node pools on vSphere can include a `vsphere` block, to define optional properties specific to configuring VMs on vSphere.

Example node pool definition for vSphere cluster:

```sh
name: tkg-wc-oidc-md-1
replicas: 4
labels:
key1: value1
key2: value2
vsphere:
memoryMiB: 8192
diskGiB: 64
numCPUs: 4
datacenter: dc0
datastore: iscsi-ds-0
storagePolicyName: name
folder: vmFolder
resourcePool: rp-1
vcIP: 10.0.0.1
template: templateName
cloneMode: clone-mode
network: network-name
```

Any values not set in the `vsphere` block inherit from the values in the cluster's first node pool.

## Update Node Pools

If you only need to change the number of nodes in a node pool, use the Tanzu CLI command in [Scale Nodes Only](#scale-nodes) below. If you want to add labels as well, follow the procedure in [Add Labels and Scale Nodes](#update-labels).

**Caution:** With these procedures, do not change existing labels, the availability zone, node instance type (on AWS or Azure), or virtual machine properties (on vSphere) of the node pool. This can have severe negative impacts on running workloads. To change these properties, create a new node pool with these properties and migrate workloads to the new node pool before deleting the original.

### Scale Nodes Only

To change the number of nodes in a node pool, run:

```shell
tanzu cluster scale CLUSTER-NAME -p NODE-POOL-NAME -w NODE-COUNT
```

Where:

- `CLUSTER-NAME` is the name of the workload cluster.
- `NODE-POOL-NAME` is the name of the node pool.
- `NODE-COUNT` is the number of nodes, as an integer, that belong in this node pool.

### Add Labels and Scale Nodes

You can add labels to a node pool and scale its nodes at the same time through the node pool configuration file.

1. Open the configuration file for the node pool you want to update.

1. If you are increasing or decreasing the number of nodes in this node pool, update the number after `replicas`.

1. If you are adding labels, indent them below `labels`. For example:

    ```yaml
    labels:
      key1: value1
      key2: value2
    ```

1. Save the node pool configuration file.

1. In a terminal, run:

    ```shell
    tanzu cluster node-pool set CLUSTER-NAME -f <CONFIG-FILE>
    ```

    If the `CLUSTER-NAME` in the command and `name` in the configuration file match a node pool in the cluster, this command updates the existing node pool instead of creating a new one.

## Delete Node Pools

To delete a node pool run:

```shell
tanzu cluster node-pool delete <CLUSTER-NAME>
```

Optionally, use `--namespace` to specify the namespace of the cluster. The default value is `default`.

**Caution:**  Migrate any workloads on these nodes to other nodes before performing this operation.

`tanzu cluster node-pool delete` does not migrate workloads off of nodes before deleting them.
