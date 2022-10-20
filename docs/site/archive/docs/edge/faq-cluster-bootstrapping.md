# FAQ on Cluster Bootstrapping

## x509: certificate signed by unknown authority when deploying Management Cluster from Windows

Upon completion of [Install the Tanzu CLI](getting-started) steps for Windows, you may encounter the following when performing `tanzu management-cluster create --ui`:

```powershell
PS > tanzu management-cluster create --ui
Downloading TKG compatibility file from 'projects.registry.vmware.com/tkg/framework-zshippable/tkg-compatibility'
Error: unable to ensure prerequisites: unable to ensure tkg BOM file: failed to download TKG compatibility file from the registry: failed to list TKG compatibility image tags: Get "https://projects.registry.vmware.com/v2/": x509: certificate signed by unknown authority
Usage:
  tanzu management-cluster create [flags]
  ...
```

To workaround this, edit `%USERPROFILE%\.config\tanzu\tkg\config.yaml` with the following, then retry management cluster setup:

```yaml
release:
    version: ""
TKG_CUSTOM_IMAGE_REPOSITORY_SKIP_TLS_VERIFY: true
```

|||
|:---------------- | --- |
|Bootstrap platform| Windows |
|Target platform   | Any |
|Affected versions | v0.2.1 |

## Failed to create inotify errors

Occasionally there are failures during management cluster deployment. In the
case where deployment fails at the step indicating **"Install providers on
management cluster"**, you may see the following error in the
cap*x*-controller-manager logs when following the
[Troubleshoot Clusters with Tanzu Diagnostics](tanzu-diagnostics) instructions.

```txt
Failed to create inotify object: Too many open files
```

This can be more common when the system being used for deployment already has
several containers running.

On Linux, this issue can be resolved by increasing the inotify watch limits
using these commands:

```sh
sysctl fs.inotify.max_user_watches=1048576
sysctl fs.inotify.max_user_instances=8192
```

After running those commands you should be able to retry deploying a managment
cluster.
