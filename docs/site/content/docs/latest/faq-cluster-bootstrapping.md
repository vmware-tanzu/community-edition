# FAQ on Cluster Bootstrapping

## x509: certificate signed by unknown authority when deploying Management Cluster from Windows

Upon completion of [Install Tanzu Community Edition](getting-started) steps for Windows, you may encounter the following when performing `tanzu management-cluster create --ui`:

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
