# OCI-Based Plugins

## Pushing Plugins

This demonstrates the setup our our plugin repo. It should be used to build
automation.

1. Push all plugins (1 for each OS:ARCH) to `projects.registry.vmware.com/tce/cli-plugins`.

1. Set the versions for each plugin in the `plugins/config/version.yaml`

1. Update the lock file config at `lock-config.yaml`.

1. Lock the images.

    ```
    kbld -f lock-config.yaml --file plugins --imgpkg-lock-output plugins/.imgpkg/images.yml 1>> /dev/null
    ```

1. Test the rendering

    ```sh
    ytt -f . | kbld -f .
    ```

1. Push the bundle

    ```sh
    imgpkg push -f plugins/ -b projects.registry.vmware.com/tce/cli-plugins/release:v0.10.0-dev.5
    ```

## Resolving Plugins

1. Setup your client config

    ```yaml
    apiVersion: config.tanzu.vmware.com/v1alpha1
    clientOptions:
      cli:
        compatibilityFilePath: projects-stg.registry.vmware.com/tkg
        discoverySources:
        - local:
            name: default-local
            path: standalone
        - oci:
            image: projects-stg.registry.vmware.com/tkg/packages/standalone/standalone-plugins:v0.17.0-dev-24-gbc750f62_vmware.1
            name: standalone-oci
        - oci:
            image: projects.registry.vmware.com/tce/cli-plugins/release:v0.10.0-dev.5
            name: standalone-oci-tce
        edition: tkg
        repositories:
        - gcpPluginRepository:
            bucketName: tanzu-cli-framework
            name: core
        unstableVersionSelector: experimental
      features:
        cluster:
          custom-nameservers: "false"
          dual-stack-ipv4-primary: "false"
          dual-stack-ipv6-primary: "false"
        global:
          context-aware-cli-for-plugins: "true"
          context-aware-discovery: "false"
          use-context-aware-discovery: "false"
        management-cluster:
          custom-nameservers: "false"
          dual-stack-ipv4-primary: "false"
          dual-stack-ipv6-primary: "false"
          export-from-confirm: "true"
          import: "false"
          network-separation-beta: "false"
          standalone-cluster-mode: "false"
    ```
