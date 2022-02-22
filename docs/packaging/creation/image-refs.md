# Image References

[kbld](./tooling.md#kbld) is used to create immutable image references within the package configuration. By resolving image references with tags to image references with SHA-256 digests, you're guaranteed to get exactly what you expect in your package.

When kbld runs, it parses your configuration files and finds images. It will then lookup the images on their registries and get their `sha256` digest. This mapping will then be placed into an `images.yml` lock file in the `bundle/.imgpkg` directory. The mapping file can be used for different scenarios in the future; one being the ability to copy a package to removable media for transfer to an air-gapped network, and the second being retrieval to a cluster by kapp-controller.

## Example Usage

1. Run the following command to create the `images.yml` file:

    ```shell
    kbld --file bundle --imgpkg-lock-output bundle/.imgpkg/images.yml 1>> /dev/null
    ```

   Here is what the `images.yml` file should look like, using the cert-manager package as an example.

    ```yaml
    ---
    apiVersion: imgpkg.carvel.dev/v1alpha1
    images:
    - annotations:
        kbld.carvel.dev/id: quay.io/jetstack/cert-manager-cainjector:v1.5.3
        kbld.carvel.dev/origins: |
          - resolved:
              tag: v1.5.3
              url: quay.io/jetstack/cert-manager-cainjector:v1.5.3
      image: quay.io/jetstack/cert-manager-cainjector@sha256:de02e3f445cfe7c035f2a9939b948c4d043011713389d9437311a62740f20bef
    - annotations:
        kbld.carvel.dev/id: quay.io/jetstack/cert-manager-controller:v1.5.3
        kbld.carvel.dev/origins: |
          - resolved:
              tag: v1.5.3
              url: quay.io/jetstack/cert-manager-controller:v1.5.3
      image: quay.io/jetstack/cert-manager-controller@sha256:7b039d469ed739a652f3bb8a1ddc122942b66cceeb85bac315449724ee64287f
    - annotations:
        kbld.carvel.dev/id: quay.io/jetstack/cert-manager-webhook:v1.5.3
        kbld.carvel.dev/origins: |
          - resolved:
              tag: v1.5.3
              url: quay.io/jetstack/cert-manager-webhook:v1.5.3
      image: quay.io/jetstack/cert-manager-webhook@sha256:ed6354190d259524d32ae74471f93bf46bfdcf4df6f73629eedf576cd87e10b8
    kind: ImagesLock
    ```
