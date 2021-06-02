# Package and Repository Operations

The process for creating a package is outlined in the [Tanzu Packaging Process](./tanzu-packaging-process.md) document. It thoroughly covers the steps and commands necessary for package creation. While developing packages, the need to iterate on changes will arise, and to aid in that process, tasks have been added to the Makefile. This document will describe the use and function of these tasks.

## Package Development/Maintenance Process

Below is an annotated version of the process flow from the [Tanzu Packaging Process](./tanzu-packaging-process.md) document. It maps tasks to steps in the process.

![Package Workflow](/docs/img/tanzu-packaging-workflow-with-commands.png)

## Initial Package Creation

When creating a new package, manifest files and the proper directory structure must be made. For more information on the needed files and directory structure, see the [Create Directory Structure](./tanzu-packaging-process.md#1-create-directory-structure) section of the Tanzu Packaging Addon documentation. You can run the `create-package` task to stub out the directory structure and required manifest files.

```shell
make create-package NAME=foobar
```

## Updating Upstream Assets

When changes happen to upstream manifests, you can trigger a `vendir sync` to bring down the newest manifests. For more details, see the [Tanzu Addon Packaging](./tanzu-packaging-process.md#2-add-manifests) document.

To update a specific package, run:

```shell
make vendir-sync-package PACKAGE=foobar
```

To update all packages, run:

```shell
make vendir-sync-all
```

## Locking Container Images

To ensure the integrity of your packages, it is important to reference image digests. The `kbld` command will create an image lock file containing the SHA of the images referenced in the package. For more details, see [Resolve and reference image digests](./tanzu-packaging-process.md#5-resolve-and-reference-image-digests)

You can lock images for a specific package:

```shell
make lock-package-images PACKAGE=foobar
```

Or lock all packages in the repo:

```shell
make lock-images-all
```

## Pushing Packages to an OCI Repository

When the package is ready to be pushed to your OCI repository, use the `push-package` tasks. As part of pushing a package, you'll need to supply the repository and tag. The repository is the URL and path to where you want the package stored, such as `projects.registry.vmware.com/tce`. Tag your package image appropriately, with a SHA, semantic version or latest.

```shell
make push-package PACKAGE=foobar TAG=baz
```

To push all packages, use:

```shell
make push-package-all OCI_REPOSITORY=repo.example.com/tce TAG=SHA
```

In the course of developing packages and releasing new versions of TCE, it will be necessary to support multiple repositories. The first iteration of an approach to this is documented here.

### Stages

Stages repesent the various steps that a package can take throughout its development. Typical channels can include `alpha`, `beta`, or `dev`. They can really be named anything. Having different channels will allow you to start development of a package in `dev`. You can `imgpkg` up your new packages and deploy to a development repository without fear of impacting the main production repo. Once your package is ready for more testing, or a wider audience, you would promote it to say, a `beta` channel, where it will be packaged alongside other packages.  
The latest, current version of TCE will use a repository with the name `main` and a tag of `stable`. `main` represents the production channel, or end of the pipeline. Only packages that have been thoroughly tested and are ready for production consumption should be promoted to the `main:stable` tag.

To create a channel, simply create a new `channel.yaml` file in the `addons/repos` directory using the make task:

```shell
    make create-channel NAME=foobar
```

The channel file contains the following:

```yaml
#@data/values
---

package_repository:
  #! The name of the Package Repository.
  #! example: delta-foo.example.com
  name:

  #! The imgpkgBundle for the repo image.
  #! Note: this value is not known until a imgpkgBundle has been pushed to an OCI registry
  #! example: registry.example.com/foo/delta:v1
  imgpkgBundle:

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
      #! example: registry.example.com/foo/foo@sha256:abcd1234...
      image:

      #! A short description of the package.
      #! example: The foo package provides f, o and more o functionality.
      description:
```

In the values file for the channel, you can list multiple packages. However, be aware that if another Package Repository already defines a package with the same name or image reference, `kapp-controller` will error and not load the Package Repository.

>The path to the imgpkgBundle for the PackageRepository is not known at this time. The next step must be completed and the imgpkgBundle must be pushed up to the OCI repository before you can get this value.

### Creating the Package Repository

You are now ready to generate the manifests for the package(s) and your package repository. Once generated, the Package Repository can be pushed to your OCI registry and then made available to your cluster.

Once a values yaml file has been defined and properly filled out, the [Package](https://carvel.dev/kapp-controller/docs/latest/package-authoring/#creating-the-package-cr) and [Package Repository](https://carvel.dev/kapp-controller/docs/latest/package-consumption/#adding-package-repository) CRs can be generated. This is done by executing the `generate-package-repository-metadata` Makefile task. This task takes 2 arguments, the `CHANNEL` and a `REPO_TAG`. `CHANNEL` is name that was used in the previous step. The `REPO_TAG` represents the tag for the image in the OCI repository. It can be whatever you like, but it might make sense to use a [SemVer](https://semver.org) to track the changes.

```shell
make generate-package-metadata CHANNEL=delta REPO_TAG=latest
```

The output of running this task is a directory structure for an [imgpkgBundle](https://carvel.dev/imgpkg/docs/latest/resources/#bundle). It contains the Package CRs and an image lock file. All the Packge CR's are concatenated together in the same file.

```txt
./addons/repos/generated/delta
├── .imgpkg
├── ├── images.yml
├── packages
├── ├── packages.yaml
```

At this time, you must manually push the bundle to your repo. This is not automated at this time simply because we don't want a repository change to happen by accident. Push the imgpkgBundle to your repo:

```shell
imgpkg push -b projects.registry.vmware.com/tce/delta:latest -f addons/repos/generated/delta

dir: .
dir: .imgpkg
file: .imgpkg/images.yml
dir: packages
file: packages/packages.yaml
Pushed 'registry.example.com/tce/delta@sha256:3d56fade82206187d4077b346a68a1a68cd5cdc74f64a3fb401753f68dd9f062'
Succeeded
```

This command outputs the URL/path to the imgpkgBundle with a SHA. The SHA needs to be placed into the values file for the current channel in the `package_repositry.imgpkgBundle` field. After updating that value, you can generate the CR for package repository. This command will create the `addons/repos/generated/$CHANNEL-package-repository.yaml` file. To make the new Package Repository available to your cluster, apply the YAML file.

```shell
make generate-package-repository-metadata CHANNEL=delta
===> Generating package repository metadata for delta
To push this repository to your cluster, run the following command:
        tanzu package repository install -f addons/repos/generated/delta-package-repository.yaml
```

Once again, you must manually update your cluster with the new repository.

```shell
tanzu package repository install -f addons/repos/generated/delta-package-repository.yaml
```

Verify that your cluster has the new package repository.

```shell
tanzu package repository list
NAME                     AGE
delta-foo.example.com    19s
```
