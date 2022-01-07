# Package Maintainership

* Proposal: [https://github.com/vmware-tanzu/community-edition/issues/2736](https://github.com/vmware-tanzu/community-edition/issues/2736)

This document applies to packages found in the official Tanzu Community Edition package repositories. There are currently 2 types of packages: Core and User Managed.

Core packages are those that are necessary to bootstrap clusters and that are required to actually run the management and workload clusters. They are defined by the upstream Tanzu Framework package and cannot be changed by end users. The current list of core packages is as follows:

* ako-operator
* antrea
* calico
* kapp-controller
* load-balancer-and-ingress-service
* local-path-storage
* pinniped
* secretgen-controller
* sriov-network-device-plugin
* vsphere-cpi
* vsphere-csi

User managed packages are those that are available to be installed by an end user to a running cluster. These packages are optional installs, not required by TCE to manage or run clusters. The current list of user managed packages is as follows:

* cert-manager
* contour
* external-dns
* fluent-bit
* gatekeeper
* grafana
* harbor
* knative-serving
* local-path-storage
* multus-cni
* prometheus
* velero
* whereabouts

TODO list cli plugin repo - Why is this relevant to a package maintainer doc?

## Definition

A package maintainer is an individual or organization that owns the package configuration for software deployed in Tanzu Community Editon. A package may have one to many package maintainers associated with it. This document covers key sections on:

* How to contribute a package
* Expectations of package maintainers

## Contributing a Package

This section covers what is required to introduce a new package into Tanzu Community Edition.

### Determine the location of your source configuration

Tanzu Community Edition requires all package source is:

* Available in a public git repository
* Licensed under (TODO)

Many packages used in TCE store their source in [community-edition/addons/packages](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages). Other options are to incorporate package code directly into the main source repository or to create a separate source repository that is specifically for package related files.

### Opening a Proposal (GitHub Issue)

To submit a package for consideration into TCE's package repository, you will need to create a GitHub Issue. The issue should provide details about:

* Functionality it provides
* Software involved
* TODO (joshrosso)

**highly recommended:** Wait for `status/approved` on proposal before doing work. You're welcome to begin work immediately, but if the proposal is `status/declined`, the work may go to waste.

### Required Tools

To develop packages for Tanzu Community Edition, the following [Carvel](https://carvel.dev) tools are necessary to be installed on the local development machine:

* Carvel Toolsuite
  * vendir
  * kbld
  * ytt
  * imgpkg

### Creating the Package Skeleton

Packages should conform to the following directory structure.

```shell
├── 1.2.3
│   ├── README.md
│   ├── bundle
│   │   ├── .imgpkg
│   │   │   └── images.yml
│   │   ├── config
│   │   │   ├── overlays
│   │   │   │   ├── overlay-a.yaml
│   │   │   │   └── overlay-b.yaml
│   │   │   ├── upstream
│   │   │   │   ├── upstream-a.yaml
│   │   │   │   └── upstream-b.yaml
│   │   │   └── values.yaml
│   │   ├── schema.yml
│   │   ├── vendir.lock.yml
│   │   └── vendir.yml
│   └── package.yaml
├── metadata.yaml
└── test
```

### Import upstream dependencies (via vendir)

When relevant, point at upstream dependencies. For example, if the software being packaged is available as a release on GitHub, reference the release and version in the `vendir.yml` file. By running `vendir`, the upstream resources will be downloaded to your local package. Doing this ensures that you have the proper upstream resources.

### Push upstream container images into `projects.registry.vmware.com/tce/images` (via imgpkg)

```shell
wget -q -O -
https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml | grep -i image:
          image: "quay.io/jetstack/cert-manager-cainjector:v1.6.1"
          image: "quay.io/jetstack/cert-manager-controller:v1.6.1"
          image: "quay.io/jetstack/cert-manager-webhook:v1.6.1"
```

```shell
$ crane copy quay.io/jetstack/cert-manager-controller:v1.6.1 projects.registry.vmware.com/tce/images/jetstack/cert-manager-controller:v1.6.1

2021/12/13 16:49:13 Copying from quay.io/jetstack/cert-manager-controller:v1.6.1 to projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller:v1.6.1
2021/12/13 16:49:15 existing manifest: sha256:41917b5d23b4abe3f5c34a156b1554e49e41185431361af46640580e4d6258fc
2021/12/13 16:49:16 existing blob: sha256:ec52731e927332d44613a9b1d70e396792d20a50bccfa06332a371e1c68d7785
2021/12/13 16:49:17 existing blob: sha256:dc34538f67ce001ae34667e7a528f5d7f1b7373b4c897cec96b54920a46cde65
2021/12/13 16:49:17 pushed blob: sha256:a6dbf7b27db03dd5a6e8d423d831a2574a72cc170d47fbae95318d3eeae32149
2021/12/13 16:49:57 pushed blob: sha256:29e5180199b812b0af5fe3d7cbe11787ba3234935537ec14ad0adf56847f005d
2021/12/13 16:49:58 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller@sha256:e2be0d9dfa684e1abf5ef9b24b601b1ca6b9dd6d725342b13c18b44156518b49: digest: sha256:e2be0d9dfa684e1abf5ef9b24b601b1ca6b9dd6d725342b13c18b44156518b49 size: 947
2021/12/13 16:49:59 existing blob: sha256:ec52731e927332d44613a9b1d70e396792d20a50bccfa06332a371e1c68d7785
2021/12/13 16:49:59 existing blob: sha256:dc34538f67ce001ae34667e7a528f5d7f1b7373b4c897cec96b54920a46cde65
2021/12/13 16:50:00 pushed blob: sha256:24882da6a70629e1639eb5bff873474c56a8c794a4adeca7cde9ed3fcda12102
2021/12/13 16:50:42 pushed blob: sha256:313817109359e805c69c3824ca6bc0a4a491e8b418399f0beea479d140541973
2021/12/13 16:50:43 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller@sha256:8898cc51a41a7848076cd7735e5a86feee734f13e802c563ef1deaafe6685040: digest: sha256:8898cc51a41a7848076cd7735e5a86feee734f13e802c563ef1deaafe6685040 size: 947
2021/12/13 16:50:44 existing blob: sha256:ec52731e927332d44613a9b1d70e396792d20a50bccfa06332a371e1c68d7785
2021/12/13 16:50:44 existing blob: sha256:dc34538f67ce001ae34667e7a528f5d7f1b7373b4c897cec96b54920a46cde65
2021/12/13 16:50:45 pushed blob: sha256:0714e6c1a7c35f6ea4fa848f83b7a8f341e3dcf44b5a5721fc569132d151a40c
2021/12/13 16:51:23 pushed blob: sha256:b68f7fa8b507c96446c17634e98eadacfac7b0473da27558ea4c9df64edd0fb6
2021/12/13 16:51:24 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller@sha256:7a60aca7f3c33e58f722229a139514b24cee45881b4c39428ae3cc252ef3190d: digest: sha256:7a60aca7f3c33e58f722229a139514b24cee45881b4c39428ae3cc252ef3190d size: 947
2021/12/13 16:51:25 existing blob: sha256:ec52731e927332d44613a9b1d70e396792d20a50bccfa06332a371e1c68d7785
2021/12/13 16:51:25 existing blob: sha256:dc34538f67ce001ae34667e7a528f5d7f1b7373b4c897cec96b54920a46cde65
2021/12/13 16:51:26 pushed blob: sha256:19542d9fe421c98aa84668010a0842501e30f6a99007846962ec1f2bcf6f6b37
2021/12/13 16:52:14 pushed blob: sha256:2a38dfa462ca3cb493a46809d9f587c3df314c96c62697a9a23aad9782f00990
2021/12/13 16:52:14 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller@sha256:1faa4c99e61db1e2227ca074de4e40c4e9008335f009fd6fd139c07ac4d5024b: digest: sha256:1faa4c99e61db1e2227ca074de4e40c4e9008335f009fd6fd139c07ac4d5024b size: 947
2021/12/13 16:52:15 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller:v1.6.1: digest: sha256:fef465f62524ed89c27451752385ab69e5c35ea4bc48b62bf61f733916ea674c size: 1723
```

### Resolve digest values for all images in `projects.registry.vmware.com/tce/images`

[kbld](https://carvel.dev/kbld/) is a Carvel tool that enables you to ensure that you're using the correct versions of software when you are creating a package. It allows you to build your package configuration with immutable image references. kbld scans a package configuration for image references and resolves those references to digests. For example, it allows you to  specify image `cert-manager:1.5.3` which is actually `quay.io/jetstack/cert-manager-controller@sha256:7b039d469ed739a652f3bb8a1ddc122942b66cceeb85bac315449724ee64287f`.

kbld scans a package configuration for any references to images and creates a mapping of image tags to a URL with a `sha256` digest. As images with the same name and tag on different registries are not necessarily the same images, by referring to an image with a digest, you're guaranteed to get the image that you're expecting. This is similar to providing a checksum file alongside an executable on a download site.

### Add overlays and configurable values (ytt)

Overlays provide a means for the package maintainer to modify or configure the behavior of the underlying software in the package. Overlays are processed by [ytt](https://carvel.dev/ytt/).

The package maintainer will create a `schema.yaml` file that contains the configuration values available in the package. For each configuration value there should be a template or overlay that modifies the underlying software's configuration.

### Create Package Metadata

The package maintainer will need to provide general metadata about the package. This metadata will be exposed in the package repository that makes a package available to a cluster.

* Category(s) of the type of functionality provided
* Display name
* Icon image in SVG base64 format
* Short description of the package
* Long Description of the package
* List of Package Maintainers
* Provider Name
* Where/how to find support for the package

### Create documentation

This should include a brief overview of the software components contained in the package, a description of configuration parameters, and general usage information. This documentation is not inteneded to replace, or be as extensive as the official documentation for the software.

The package documentation should highlight dependencies or considerations on other packages, software, Kubernetes distributions, or underlying infrastructure (e.g. AWS, GCP, Docker, vSphere, etc).

### Run linting checks

* Validate Markdown for documentation
* Validate schema is present
* Validate overlays and templates
* TODO - more checks

### Push Package Bundle to an OCI Registry

Initial software components that are to be packaged
Configuration parameters that are to be exposed
Basic tests should be provided along with execution instructions.

### Open Pull Request in community-edition repo

When your package is complete and ready for acceptance, create a Pull Request in the [community-edition](https://github.com/vmware-tanzu/community-edition/pulls) GitHub repository. The PR should reference the GitHub issue created for introducing the package.

## Expectations (Role of Package Maintainer)

A package maintainer is a person or team responsible for the creation, maintenance and support of a package. Ideally, a maintainer is associated with the group or organization that is responsible for the development of the software contained within the package. This closeness will ensure that knowledgeable persons are able to develop and support the package. However, being involved in the project is not required to be a package maintainer.

Packages are contributed [according to the guides and documentation provided](TODO:Link) in Tanzu Community Edition.

### Who Can be a Maintainer

// TODO(joshrosso): we need to figure this out with PM and such

### Package Validation (Linting)

Base image checks (Alpine not allowed)
Markdown linting
Directory structure, required files (package.yaml, metadata.yaml, lock files, readme.md etc)

### Software Updates

Creating new packages to track the latest version of the underlying software

### Package Updates

Providing package updates containing bug fixes and/or enhancements to the package itself

Creating new packages to track the latest version of the underlying software

### CVE Remediation

Security concerns and CVE's can occur within a package in 2 different manners. The first is in the underlying, packaged software. The other is in the package itself.

If a concern is present in the underlying software, the package maintainer is expected to update the package with patched or updated versions of the effected software when available. The package maintainer is not expected to patch the underlying software.

If the concern is in the package itself, the package maintainer is expected to provide a patch as soon as possible. For example, if an overlay introduces a vulnerability, it is expected that the package maintainer fix the overlay to remedy the concern.

### Exposing Schemas

The package maintainer is responsible for creating a [schema](https://carvel.dev/ytt/docs/latest/how-to-write-schema/) that details the available configuration parameters. This schema can be used to generate package documentation and a default values.yaml file for use when installing the package.

### Testing (and Coverage)

The package should also include end to end tests to verify basic functionality.

### Publishing Container Images

### Publishing Package Images

### Documentation

Providing updated documentation on package usage and configuration

Maintainers should provide documentation on the usage of the package, including details on installation and configuration.

### Supporting Package End-Users

Community support is expected of Package Maintainers. Maintainers shall monitor and respond in appropriate Slack workspaces or other online messaging channels relevant to the usage of the package. Issues/Bugs/Enhancements that are filed against the package shall be addressed and responded to by maintainers.
