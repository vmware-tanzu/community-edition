# Package Maintainership

* Proposal: [https://github.com/vmware-tanzu/community-edition/issues/2736](https://github.com/vmware-tanzu/community-edition/issues/2736)

This document applies to packages found in the official Tanzu Community Edition package repositories.

* TODO list core
* TODO list user-managed
* TODO list cli plugin repo

## Definition

A package maintainer is an individual that owns the package configuration for software deployed in Tanzu Community Editon. A package may have one to many package maintainers associated with it. This document covers key sections on:

* How to contribute a package
* Expectations of package maintainers

## Contributing a Package

This section covers what is required to introduce a new package into community-edition.

### Determine the location of your source configuration

Tanzu Community Edition requires all package source is:

* Available in a public git repository
* Licensed under (TODO)

Many packages used in TCE store their source in [community-edition/addons/packages](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages).

### Opening a Proposal (GitHub Issue)

Provide details about:

* Functionality it provides
* Software involved
* TODO(josrhosso)

**highly recommended:** Wait for `status/approved` on proposal before doing work. You're welcome to begin work immediately, but if the proposal is `status/declined`, the work may go to waste.

### Install Required Tools

Locally installed

* vendir
* kbld
* ytt
* imgpkg

Server/Cluster apps

* kapp-controller

### Creating the Package Skeleton

Directory structure

```shell
├── 1.2.3
│   ├── README.md
│   ├── bundle
│   │   ├── .imgpkg
│   │   │   └── images.yml
│   │   ├── config
│   │   │   ├── overlays
│   │   │   │   └── overlay-a.yaml
│   │   │   ├── upstream
│   │   │   │   └── upstream-crd.yaml
│   │   │   └── values.yaml
│   │   ├── vendir.lock.yml
│   │   └── vendir.yml
│   └── package.yaml
├── metadata.yaml
└── test
```

### Import upstream dependencies (via vendir)

When relevant, point at upstream dependencies

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

kbld

### Add overlays and configurable values (ytt)

* Produce schema (used in the readme.md and package.yaml)

### Create Package Metadata

* categories
* displayName
* iconSVGBase64
* longDescription
* maintainers
* providerName
* shortDescription
* supportDescription

### Create documentation

This should include a brief overview of the software components contained in the package, a description of configuration parameters, and general usage information.

Interactions and dependencies between packages

Compatible Kubernetes distributions

Considerations for underlying infrastructure (e.g. AWS, GCP, Docker, vSphere, etc)

### Run linting checks

TODO: make tasks here that can run various validations

### Open Pull Request in community-edition repo

Reference proposal issue

### Push Package Bundle into `projects.registry.vmware.com/tce/packages`

Initial software components that are to be packaged
Configuration parameters that are to be exposed
Basic tests should be provided along with execution instructions.

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

### CVE Remdiation

Security issues and/or CVE’s shall be handled within an acceptable time frame

### Exposing Schemas

### Testing (and Coverage)

The package should also include end to end tests to verify basic functionality.

### Publishing Container Images

### Publishing Package Images

### Documentation

Providing updated documentation on package usage and configuration

Maintainers should provide documentation on the usage of the package, including details on installation and configuration.

### Supporting Package End-Users

Community support is expected of Package Maintainers. Maintainers shall monitor and respond in appropriate Slack workspaces or other online messaging channels relevant to the usage of the package. Issues/Bugs/Enhancements that are filed against the package shall be addressed and responded to by maintainers.
