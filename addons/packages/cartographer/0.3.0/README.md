# Cartographer

Cartographer allows you to create secure and reusable supply chains that define
all of your application CI and CD in one place, in cluster.

## Components

* cartographer

## Supported Providers

The following table shows the providers this package can work with.

| AWS  | Azure | vSphere | Docker |
|------|-------|---------|--------|
| ✅   | ✅    | ✅      | ✅     |

## Configuration

The Cartographer package has no configurable properties.

## Installation

The Cartographer package requires use of cert-manager for certificate
generation.

1. Install cert-manager Package

   ```shell
   tanzu package install cert-manager \
      --package-name cert-manager.community.tanzu.vmware.com \
      --version ${CERT_MANAGER_PACKAGE_VERSION}
   ```

   > You can get the `${CERT_MANAGER_PACKAGE_VERSION}` from running `tanzu
   > package available list cert-manager.community.tanzu.vmware.com`.
   > Specifying a namespace may be required depending on where your package
   > repository was installed.

2. Install the Cartographer package

   ```shell
   tanzu package install cartographer \
      --package-name cartographer.community.tanzu.vmware.com \
      --version ${CARTOGRAPHER_PACKAGE_VERSION}
   ```

   > You can get the `${CARTOGRAPHER_PACKAGE_VERSION}` from running `tanzu
   > package available list cartographer.community.tanzu.vmware.com`.
   > Specifying a namespace may be required depending on where your package
   > repository was installed.

## Documentation

For documentation specific to Cartographer, check out
[cartographer.sh](https://cartographer.sh) and the main repository
[vmware-tanzu/cartographer](https://github.com/vmware-tanzu/cartographer).

[carvel]: https://carvel.dev/
[Cartographer]: https://cartographer.sh
[kapp-controller]: https://github.com/vmware-tanzu/carvel-kapp-controller
[Tanzu CLI]: https://github.com/vmware-tanzu/tanzu-framework
[cert-manager]: https://github.com/cert-manager/cert-manager
