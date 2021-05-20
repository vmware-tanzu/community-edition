# BOMs and Providers

## TCE Customization Hacks

Note, the following sections are workarounds we currently use to ship
customization's to TKG's BOM/TKR. These steps are intended to be moved away from
as soon as possible.

### Management Cluster: BOM/TKR Customization

1. Clone the `tkg-cli` repository.

1. Create a new branch off the branch you wish to build atop.

1. Sync the bom and tkr to ensure you have the correct version.

   ```sh
   make sync-bom
   ```

   > This will place the bom and tkr in `bom/originals/`

1. Edit the TKR/BOM in the above directory as needed.

1. Generate the bin-data based on your changes.

  ```sh
  make generate-bindata
  ```

1. Commit and push your changes.

1. Ensure the `install-cli` target in the TCE Makefile references the proper
   branch.

  ```make
 .PHONY: install-cli
install-cli:
  TANZU_CORE_REPO_BRANCH="v1.3.0" TKG_PROVIDERS_REPO_BRANCH="tce"  TKG_CLI_REPO_BRANCH="tce" BUILD_VERSION=${CORE_BUILD_VERSION} hack/build-tanzu.sh
  ```

  > Note in the above `TKG_CLI_REPO_BRANCH` is set to `tce`.

1. Update go.mod `replace` to point to your correct revision of tkg-providers.

  ```txt
  github.com/vmware-tanzu-private/tkg-providers => github.com/vmware-tanzu-private/tkg-providers v1.3.0-rc.1.0.20210415172650-c5f3a82f00e4
  ```

### Management Cluster: Providers Customization

This customization will alter providers such that all changes are reflected in
management clusters that are created.

1. Clone the tkg-providers repo.

  ```sh
  https://github.com/vmware-tanzu-private/tkg-providers
  ```

1. Create a new branch off the branch you wish to build atop.

1. Alter the providers in any way you need to.

1. Pump the provider changes into bin-data.

  ```sh
  make generate-bindata
  ```

1. Commit and push the changes.

1. Ensure the `install-cli` target in the TCE Makefile references the proper
   branch.

  ```make
 .PHONY: install-cli
install-cli:
  TANZU_CORE_REPO_BRANCH="v1.3.0" TKG_PROVIDERS_REPO_BRANCH="tce"  TKG_CLI_REPO_BRANCH="tce" BUILD_VERSION=${CORE_BUILD_VERSION} hack/build-tanzu.sh
  ```

  > Note in the above `TKG_PROVIDERS_REPO_BRANCH` is set to `tce`.

1. Update go.mod `replace` to point to your correct revision of tkg-providers.

  ```txt
  github.com/vmware-tanzu-private/tkg-providers => github.com/vmware-tanzu-private/tkg-providers v1.3.0-rc.1.0.20210415172650-c5f3a82f00e4
  ```

### Guest Cluster: TKR BOM Customization

This customization drives the BOM used in the workload clusters.

Currently cannot figure out how this is generated. So hacking around it as
follows.

1. Pull down the existing TKG TKR BOM.

    ```sh
    crane pull projects.registry.vmware.com/tce/tkr-bom:v1.20.4_vmware.1-tkg.1 test.tar
    ```

1. Untar the output.

    ```sh
    tar xvf test.tar
    tar zxvf ${LAYER_HASH}.tar.gz
    ```

1. Edit the BOM file.

    ```sh
    vim tkr-bom-v1.20.4+vmware.1-tkg.1.yaml
    ```

    > Version could vary

1. Push the BOM file to the TCE repository.

    ```sh
    imgpkg push -i projects.registry.vmware.com/tce/tkr-bom:v1.20.4_vmware.1-tkg.1 -f tkr-bom-v1.20.4+vmware.1-tkg.1.yaml
    ```

    > At the time of this writing image (`-i`) is being used, not bundle (`-b`).

## Future Plan

Long-term, we need to figure out how to wire this up with the TKG Bolt process.
This **cannot** be our longterm solution and should be considered intentional
technical debt.
