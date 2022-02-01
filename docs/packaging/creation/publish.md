# Publishing a Package

The final step in creating a package is to publish it. This can be done with the
[Carvel imgpkg tool](../tooling/). imgpkg is a tool that allows users to store a
set of arbitrary files as an OCI image and push it to a container registry.

## Repository

Packages are pushed as OCI bundles to a container registry.

**All** official TCE packages must be hosted in `projects.registry.vmware.com`.

## Example Usage

To publish your package, you first need to login to your OCI Registry. This can
typically be done with the `docker login` command. Use the `imgpkg` command to
create a bundle from the `bundle` directory of your package. This example uses
the cert-manager package and `ttl.sh` for the OCI Registry.

 ```shell
 imgpkg push --bundle projects.registry.vmware.com/tce/packages/cert-manager:6h --file bundle/

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
 Pushed 'projects.registry.vmware.com/tce/packages/cert-manager@sha256:7335d2f20d000695e7880731ad24406c3d98dff5008edb04fa30f16e31abbb1a'
 Succeeded
 ```

> For successfully pushed images, `imgpkg` will report the digest of the
package. In this example,
`projects.registry.vmware.com/tce/packages/cert-manager@sha256:7335d2f20d000695e7880731ad24406c3d98dff5008edb04fa30f16e31abbb1a`.
Take note of this value as it is needed in the Package CR.
