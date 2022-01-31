---
title: Publishing a Package
weight: 3
---

The final step in creating a package is to publish it. This can be done with the [Carvel imgpkg tool](../tooling/). imgpkg is a tool that allows users to store a set of arbitrary files as an OCI image.

## Example Usage

To publish your package, you first need to login to your OCI Registry. This can typically be done with the `docker login` command. Use the `imgpkg` command to create a bundle from the `bundle` directory of your package. This example uses the cert-manager package and `ttl.sh` for the OCI Registry.

 ```shell
 imgpkg push --bundle ttl.sh/seemiller/cert-manager:6h --file bundle/

 dir: .
 dir: .imgpkg
 file: .imgpkg/images.yml
 dir: config
 dir: config/overlays
 file: config/overlays/annotations.yaml
 file: config/overlays/deployment.yaml
 file: config/overlays/misc.yaml
 dir: config/upstream
 file: config/upstream/cert-manager.yaml
 file: config/values.yaml
 file: vendir.lock.yml
 file: vendir.yml
 Pushed 'ttl.sh/seemiller/cert-manager@sha256:7335d2f20d000695e7880731ad24406c3d98dff5008edb04fa30f16e31abbb1a'
 Succeeded
 ```

> For successfully pushed images, `imgpkg` will report the digest of the package. In this example, `ttl.sh/seemiller/cert-manager@sha256:7335d2f2...`. Take note of this value as it is needed in the Package CR.

