# Development

## Updating ExternalDNS version in the package

These are the necessary steps for updating ExternalDNS.

We use [vendir](https://carvel.dev/vendir/) to vendor in the upstream
ExternalDNS manifests, and [ytt](https://carvel.dev/ytt/) overlays to change
them.

To update the upstream manifests, use in the `./bundle` directory:

```shell
vendir sync
```

Check to see if anything changed upstream after the sync, and update
any overlays accordingly.

To update the ExternalDNS image that is substituted into the manifests by
[kbld](https://carvel.dev/kbld/):

1. Update the `kbld-config.yml` with the new ExternalDNS image.

2. Run in the `./bundle` directory:

```shell
kbld --imgpkg-lock-output .imgpkg/images.yml -f .
```

This will update the image in the imgpkg lock.

To rebuild the package image—which consists of all the yaml files within the `./bundle`
directory for local testing—run:

```shell
imgpkg push -b <some-image-registry>/external-dns-package:dev -f .
```

**Replace `<some-image-registry>` with a registry you can push to for testing.**

You can then change the `imgpkgBundle.image` in `package.yaml`. Note that this is
not something you will commit in your final PR. Once the bundle changes are
merged, an official build of the imgpkgBundle will be created and updated in the
`package.yaml`.

For local testing you can apply the `./metadata.yaml` and `./package.yaml` and
run the `e2e` tests.

```shell
kubectl apply -n tanzu-package-repo-global -f ../metadata.yaml
kubectl apply -n tanzu-package-repo-global -f ./package.yaml
```

To run the e2e tests:

```shell
cd test
make e2e-test
```

## Official Build instructions

Like above, but push the bundle to the TCE registry:

```shell
imgpkg push -b projects.registry.vmware.com/tce/external-dns:${version_semver}
```

NOTE: in the version semver `+` should be substituted with `_`.

The resulting sha256 bundle url should be copied into the `package.yaml`
`imgpkgBundle.image` field and committed.

## Generate package `valuesSchema` from ytt schema

1. Generate OpenAPIv3 schema from `schema.yaml` using `ytt`.

```bash
ytt -f bundle/config/schema.yaml --data-values-schema-inspect --output openapi-v3
```

1. Copy contents of `components.schemas.dataValues` into
`spec.valuesSchema.openAPIv3` of `package.yaml`.
