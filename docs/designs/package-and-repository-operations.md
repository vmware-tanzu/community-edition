# Package and Repository Operations

The process for creating a package is outlined in the [Tanzu Addon Packaging](./tanzu-addon-packaging.md) document. That document thoroughly covers in depth the steps and commands necessary for package creation. While developing packages, the need to iterate on changes will arise, and to aid in that process, tasks have been added to the Makefile. This document will describe the use and function of these tasks.

## Package Development/Maintenance Process

Here is an annotated version of the process flow from the Tanzu Addon Packaging document. It maps tasks to steps in the process. You will most likely iterate on the templating and bundling manifest steps.

>_TODO: Update this diagram_

![Package Workflow](../images/tanzu-packaging-workflow-with-commands.png)

## Initial Package Creation

When first starting off creating a new package, certain manifest files and the proper directory structure must be made.

Manifests needed for creating new packages are:

- `clusterrolebinding.yaml`
  
  > Permissions needed by`kapp-controller` to create objects. See [Create RBAC Assets](./tanzu-addon-packaging.md#7-create-rbac-assets) for more details.
  
- `serviceaccount.yaml`

  > Service Account needed by`kapp-controller` to create objects. See [Create RBAC Assets](./tanzu-addon-packaging.md#7-create-rbac-assets) for more details
  
- `installedPackage.yaml`
  
  > The declaration of intent to install a package, which kapp-controller will act on. See [Create a sample InstalledPackage](./tanzu-addon-packaging.md#10-reate-a-sample-installedpackage) for more details.

- `bundle/vendir.yaml`
  > File defining the upstream manifests for the package. See [Add Manifests](./tanzu-addon-packaging.md#2-add-manifests) for more details.
  
- `bundle/config/values.yaml`
  > File containing user configurable values. See [Create Default Values](./tanzu-addon-packaging.md#4-create-default-values) for more details.

- `addons/repos/main/packages/<<pacakge>>.yaml`
  > Package mainfest for the repository. See [Create a Package CR](./tanzu-addon-packaging.md#8-create-a-package-cr) for more details.

You can run the `create-package` task to stub out the directory structure and required manifest files listed above.

```shell
make create-package NAME=foobar
```

## Updating Upstream Assets

When changes happen to upstream manifests, you can trigger a `vendir sync` to bring down the newest manifests.

To update a specific package, run:

```shell
make vendir-sync-package PACKAGE=foobar
```

To update all packages, run:

```shell
make vendir-sync-all
```

## Locking Container Images

To ensure the integrity of your packages, it is important to reference image digests. The `kbld` command will create an image lock file containing the SHA of the images referenced in the package. For more details, see [Resolve and reference image digests](./tanzu-addon-packaging.md#5-resolve-and-reference-image-digests)

You can lock images for a specific package:

```shell
make lock-package-images PACKAGE=foobar
```

Or lock all packages in the repo:

```shell
make lock-images-all
```

## Pushing Packages to an OCI Repository

When the package is ready to be pushed to your OCI repository, use the `push-package` tasks. As part of pushing a package, you'll need to supply the repository and tag. The repository is the URL and path to where you want the package stored, such as `projects.registry.vmware.com/tce`. Tag your package image appripriately, with a SHA, semantic version or latest.

```shell
make push-package PACKAGE=foobar TAG=baz
```

To push all packages, use:

```shell
make push-package-all OCI_REPOSITORY=repo.example.com/tce TAG=SHA
```

In the course of developing packages and releasing new versions of TCE, it will be necessary to support multiple repositories. The first iteration of an approach to this is documented here.

### Stages

Stages repesent the various steps that a package can take throughout its development. Typical stages can include `alpha`, `beta`, or `dev`. They can really be named anything. Having different stages will allow you to start development of a package in `dev`. You can `imgpkg` up your new packages and deploy to a development repository without fear of impacting the main production repo. Once your package is ready for more testing, or a wider audience, you would promote it to say, a `beta` stage, where it will be packaged alongside other packages.  
The latest, current version of TCE will use a repository with the name `main` and a tag of `stable`. `main` represents the production stage, or end of the pipeline. Only packages that have been thoroughly tested and are ready for production consumption should be promoted to the `main:stable` tag.

To create a stage, simply create a new `stage.yaml` file in the `addons/repos` directory.

>_TODO: Consider making this a Makefile task or script_

```shell
export STAGE=delta
cat >> addons/repos/${STAGE}.yaml <<EOL
#@data/values
---

package_repository:
  #! The name of the Package Repository.
  #! example: delta-foo.example.com
  name:

  #! The imageBundle or URL for the repo image.
  #! Note: this value is not known until a imgpkgBundle/image has been pushed to an OCI registry
  #! example: registry.example.com/foo/delta:v1
  imageBundle:
  url:

  packages:
      #! The name of the package.
      #! example: foo
    - name:

      #! The domain that the package belongs to. This is used in conjunction with the name to create a fully qualified domain name for the package.
      #! example: example.com
      domain:

      #! The version of the package.
      #! example: 0.0.1
      version:

      #! The path to the image in the OCI repository.
      #! example: registry.example.com/foo/foo@sha256:abababababababababababababababababababababababababababababababab
      image:

      #! A short description of the package.
      #! example: The foo package provides f, o and more o functionality.
      description:
EOL
```

In the values file for the stage, you can list multiple packages. However, be aware that if another Package Repository already defines a package with the same name or image reference, `kapp-controller` will error and not load the Package Repository.

>Also note, the path to the imageBundle or the url for the PackageRepository is not known at this time. The next step must be completed and the image bundle must be pushed up to the OCI repository before you can get this value.

### CR Generation

Once a values yaml file has been defined and properly filled out, the [Package](https://carvel.dev/kapp-controller/docs/latest/package-authoring/#creating-the-package-cr) and [Package Repository](https://carvel.dev/kapp-controller/docs/latest/package-consumption/#adding-package-repository) CRs can be generated. This is done by executing the `generate-package-repository-metadata` Makefile task. This task takes 2 arguments, the `STAGE` and a `REPO_TAG`. `STAGE` is name that was used in the previous step. The `REPO_TAG` represents the tag for the image in the OCI repository. It can be whatever you like, but it might make sense to use a [SemVer](https://semver.org) to track the changes.

```shell
make generate-package-metadata STAGE=delta REPO_TAG=latest
```

#### imgpkg Bundle

The output of running this task is a directory structure for an [imgpkgBundle](https://carvel.dev/imgpkg/docs/latest/resources/#bundle). It contains the Package CRs and an image lock file. All the Packge CR's are concatenated together in the same file.

```txt
./addons/repos/stages/delta
├── .imgpkg
├── ├── images.yml
├── packages
├── ├── packages.yaml
```

At this time, you must manually push the bundle to your repo. This is not automated at this time simply because we don't want a repository change to happen by accident. Push the imgpkgBundle to your repo:

```shell
imgpkg push -b projects.registry.vmware.com/tce/delta:latest -f addons/repos/stages/delta                                                                                                                                         ─╯
dir: .
dir: .imgpkg
file: .imgpkg/images.yml
dir: packages
file: packages/packages.yaml
Pushed 'registry.example.com/tce/delta@sha256:3d56fade82206187d4077b346a68a1a68cd5cdc74f64a3fb401753f68dd9f062'
Succeeded
```

Note that this command outputs the URL/path to the imgpkgBundle with a SHA. This value needs to be placed into the values file for the current stage in the `package_repositry.imgpkgBundle` field. After updating that value, you can generate the CR for package repository. This command will create the `addons/repos/stages/$STAGE-package-repository.yaml` file. To make the new Package Repository available to your cluster, apply the YAML file.

```shell
make generate-package-repository-metadata STAGE=delta
===> Generating package repository metadata for delta
To push this repository to your cluster, run the following command:
        tanzu package repository install -f addons/repos/stages/delta-package-repository.yaml
```

Once again, you must manually update your cluster with the new repository.

```shell
tanzu package repository install -f addons/repos/stages/delta-package-repository.yaml
```

Verify that your cluster has the new package repository.

```shell
tanzu package repository list
NAME                     AGE
delta-foo.example.com    19s
```
