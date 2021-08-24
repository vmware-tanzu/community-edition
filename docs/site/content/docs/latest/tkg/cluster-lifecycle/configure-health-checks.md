# Configure Machine Health Checks for Tanzu Kubernetes Clusters

This topic describes how to use the Tanzu CLI to
create, update, retrieve, and delete `MachineHealthCheck` objects for
Tanzu Kubernetes clusters.

## <a id="about"></a> About `MachineHealthCheck`

[`MachineHealthCheck`](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/machine-health-check.html#machinehealthcheck) is a controller that provides node health monitoring and node auto-repair for Tanzu Kubernetes clusters.

This controller is enabled in the global Tanzu Kubernetes Grid configuration by default, for all Tanzu Kubernetes clusters. You can override your global Tanzu Kubernetes Grid configuration for individual Tanzu Kubernetes clusters in two ways:

- When deploying the management cluster. You can enable or disable
the default `MachineHealthCheck` in either the Tanzu Kubernetes Grid installer
interface or the cluster configuration file. Each Tanzu Kubernetes cluster that you deploy with your management cluster inherits this configuration by default. For more information, see
[Deploying Management Clusters](../mgmt-clusters/deploy-management-clusters.md).
- After creating a Tanzu Kubernetes cluster.
You can use the Tanzu CLI to create, update, retrieve, and
delete `MachineHealthCheck` objects for individual Tanzu Kubernetes clusters.
See the sections below.

When `MachineHealthCheck` is enabled in a Tanzu Kubernetes cluster, it runs in the same namespace as the cluster.

```
#! ---------------------------------------------------------------------
#! Machine Health Check configuration
#! ---------------------------------------------------------------------
ENABLE_MHC: true
MHC_UNKNOWN_STATUS_TIMEOUT: 5m
MHC_FALSE_STATUS_TIMEOUT: 12m
```

## <a id="create"></a> Create or Update a `MachineHealthCheck`

To create a `MachineHealthCheck` with the default configuration, run the following command:

```
tanzu cluster machinehealthcheck set CLUSTER-NAME
```

Where `CLUSTER-NAME` is the name of the Tanzu Kubernetes cluster you want to monitor.

You can also use this command to create `MachineHealthCheck` objects with custom configuration options or update existing `MachineHealthCheck` objects. To set custom configuration options for a `MachineHealthCheck`, run the
`tanzu cluster machinehealthcheck set` command with one or more
of the following:

   * `--mhc-name`: By default, when you run `tanzu cluster machinehealthcheck set CLUSTER-NAME`, the command sets the name of the `MachineHealthCheck` to `CLUSTER-NAME`. Specify the `--mhc-name` option if you want to set a different name. For example:

       ```
       tanzu cluster machinehealthcheck set my-cluster --mhc-name my-mhc
       ```

   * `--match-labels`: This option filters machines by label keys and
   values. You can specify one or more label constraints.
   The `MachineHealthCheck` is applied to all machines that satisfy these constraints.
   Use the syntax below:

       ```
       tanzu cluster machinehealthcheck set my-cluster --match-labels "key1:value1,key2:value2"
       ```

       For example:

       ```
       tanzu cluster machinehealthcheck set my-cluster --match-labels "node-pool:my-cluster-worker-pool"
       ```

   * `--node-startup-timeout`: This option controls the amount of time that the
   `MachineHealthCheck` waits for a machine to join the cluster before considering the machine unhealthy. For example, the command below sets the `--node-startup-timeout` option to `10m`:

       ```
       tanzu cluster machinehealthcheck set my-cluster --node-startup-timeout 10m
       ```

       If a machine fails to join the cluster within this amount of time, the
       `MachineHealthCheck` recreates the machine.

   * `--unhealthy-conditions`: This option can set the
   `Ready`, `MemoryPressure`, `DiskPressure`, `PIDPressure`, and
   `NetworkUnavailable` conditions. The `MachineHealthCheck` uses the conditions that you set to determine whether a node is unhealthy. To set the status of a condition, use
   `True`, `False`, or `Unknown`. For example:

       ```
       tanzu cluster machinehealthcheck set my-cluster --unhealthy-conditions "Ready:False:5m,Ready:Unknown:5m"
       ```

       In the example above, if the status of the `Ready` node condition remains
       `Unknown` or `False` for longer than `5m`, the `MachineHealthCheck` considers the machine unhealthy and recreates it.

## <a id="get"></a> Retrieve a `MachineHealthCheck`

To retrieve a `MachineHealthCheck` object, run the following command:

```
tanzu cluster machinehealthcheck get CLUSTER-NAME
```

If you assigned a non-default name to the object, specify the `--mhc-name`
flag.

## <a id="delete"></a> Delete a `MachineHealthCheck`

To delete a `MachineHealthCheck` object, run the following command:

```
tanzu cluster machinehealthcheck delete CLUSTER-NAME
```

If you assigned a non-default name to the object, specify the `--mhc-name`
flag.
