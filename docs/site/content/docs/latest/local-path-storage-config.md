# Configuring the Local Path Storage Package

This package provides local path node storage and primarily supports RWO access mode.
It utilizes the Kubernetes [Local Persistent Volume feature](https://kubernetes.io/blog/2018/04/13/local-persistent-volumes-beta/)
and in Tanzu Community Edition, it is primarily intended for use with Docker.

## Limitations

The local-path-storage binds to a single host node
and is not intended to dynamically change hosts.
Therefore, a PVC can _only_ be used by the node that creates it.
This can lead to unintended data loss when scaling or when pods roll from one node to another.
This can make scheduling difficult since applications are "tied" to the node that creates it's PV.

Furthermore, local-path-storage does _not_ enforce capacity limitations
and may possibly overwhelm the local node's disc capacity.

See the [Local Path Storage Provisioner readme](https://github.com/rancher/local-path-provisioner)
for further configuration options.

## Configuration

| Value       | Required/Optional | Description                                         |
|:-------------|:-------------------|:-----------------------------------------------------|
| `namespace` | Required          | The namespace to deploy the local-path-storage pods |

*Note:* The local-path-storage provides a config map that may be modified _after_ installation.
This includes a `config.json` that can be used to further configure the storage provider.
Additionally, `setup` and `teardown` scripts are defined in the config map and are used in the lifecycle of persistent volumes.
The local-path-storage pods will dynamically reload the config map upon configuration without need to reapply the deployment.

## Usage Examples

A StorageClass is required in order to use PVCs and store data (which is necessary for services
like Prometheus). The local-path-storage provider enables local Docker clusters to store data locally.
Using a local PVC with Docker lets a developer work quickly on their own workstation with Docker.

A local storage provider may also be used in special cases for caching, sharding data in distributed datastores,
and other node failure tolerant storage models.
Note that local storage providers are generally not suitable for most production use cases.

