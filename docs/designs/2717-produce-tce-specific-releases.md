# Produce community-edition releases that are decoupled from Tanzu Kubernetes Grid

* Proposal:
  [https://github.com/vmware-tanzu/community-edition/issues/2717](https://github.com/vmware-tanzu/community-edition/issues/2717)

## Abstract

At the time of writing, Tanzu Community Edition (TCE) relies on the Bill of
Materials (BOMs) produced for VMware's paid product, Tanzu Kubernetes Grid
(TKG). This means that creation of clusters, both management and workload,
**always** use the same packages and versions as the rest of TKG. A BOM is
declares the components and their versions used in clusters. This includes
everything from Kubernetes versions to the software that runs on top of
Kubernetes.

This binding to a TKG BOM was a decision made to ensure we could release TCE.
Now we are at stage where we must unbind from TKG's BOM and begin producing
TCE-specific BOMs. Achieving this decoupling will enable TCE to release unique
components along with newer versions of existing components, such as Kubernetes.

### Why

Doing the above enables:

* Building and packaging in the open, rather than behind VMware's network.
* Packages that can be re-produced by users.
* Releasing Kubernetes versions ahead of TKG.
* Releasing core packages (e.g. antrea) ahead of TKG.
* Releasing core-packages specific to TCE.
* Patching TCE releases out-of-band from TKG.

## Proposal

To achieve a decoupled TCE, we plan to produce:

1. **[Build Metadata API](#1-build-metadata-api)**: a build BOM that describes
   what should be contained in a TKr.
2. **[Build Library](#2-build-library)**: reads the Build BOM and produces the
   required components and a TKr.
3. **[Host Images](#3-host-images)**: created with [image
   builder](https://github.com/kubernetes-sigs/image-builder) and pushed to
   cloud providers.
4. **[Container Images](#4-container-images)**: copied from upstream projects
   into `http://projects.registry.vmware.com/tce`.
5. **[Core Package Repository](#5-core-package-repository)**: created and pushed
   to our .
6. **[User-managed Package Repository](#6-user-managed-package-repository)**: A
   user-managed package repository.
7. **[Tanzu Kubernetes Release (TKR)](#7-tanzu-kubernetes-release)**: created
   and pushed to our OCI registry, used at runtime to create clusters.

### 1. Build Metadata API

The Build metadata API is a build BOM that describes what components should be
included in a Tkr. This API is described via a manifest by the TCE team. This
manifest can then be satisfied using the build library (see next section). The
following provides a high-level representation of what the build manifest will
be translated into via the build library.

![Diagram of build manifest assets being processed by
CI](imgs-2717/build-manifest-ci.png)

Packages (bundles pushed using `imgpkg`) are **not** pushed in this process.
Instead, package maintainers are expected to have already pushed a given
package. The build metadata above only contains a reference to the pushed
package. With this, the CI can construct a package repository, which is a image
holding references and metadata to the packages.

The asset produced above will power cluster creation at runtime in the TKr.
The TKr should **not** include build information such as specific packages used.
Instead, it should reference what is relevant at runtime. This includes assets
such as package repositories and host OS images.

##### Schema

The proposed schema of the build manifest is:

```go
type TanzuBuild struct {
  // the version of this release and build
  Version string
  // packages to create a core package repo of
  // created repo will be uploaded to an OCI repo
  // uri of uploaded repo will be added to the generated TKR
  CorePackages []MetaPackage `yaml:"corePackages"`
  // packages to create a user-managed package repo of
  // created repo will be uploaded to an OCI repo
  // uri of uploaded repo will be added to the generated TKR
  UserPackages []MetaPackage `yaml:"userPackages"`
  // host images to add to the generated TKR
  HostImages []HostImage `yaml:"hostImages"`
  // k8s metadata added to the generated TKR
  KubernetesMeta `yaml:"kubernetesMeta"`
}

type MetaPackage struct {
  // name of the metapackage, e.g. contour
  Name string
  // each versioned package to make available
  Packages []Package
  // Embeds contents of package metadata:
  // https://carvel.dev/kapp-controller/docs/latest/packaging/#package-metadata
  PackageMetadata string  `yaml:"packageMetadata"`
}

type Package struct {
  // fully qualified name of package, e.g. contour.dev.1.5.3
  Name string
  // location of the OCI bundle
  ImageBundleURI string `yaml:"imageBundleUri"`
  // Arbitrary set of options to add in the package CR
  Options map[string]interface{}
}

type HostImage struct {
  // name of the host image (ami-dfdal)
  Name string
  // the infrastructure provider it relates to (aws, azure, etc)
  Provider string
  // metadata relevant to the specific provider. (aws-region: us-west-2)
  Metadata map[string]string
}

type KubernetesMeta struct {
  // the version of Kubernetes this TKR represents
  Version string
  // the kubernetes components used in bootstrap
  // these should match what is inside the host images
  // otherwise kubeadm will pull down new images on each
  // host during bootstrap
  Components []ContainerImage
}

type ContainerImage struct {
  // the name of the container image
  Name string
  // the host, and url used to access the container image
  Repository string
  // the tag associated with the image
  Tag string
}
```

Rendering this against a configuration ([go
playground](https://go.dev/play/p/vFYA7stM_y8)) produces the following.

```yaml
version: 1.25+tce.1
corePackages:
    - name: antrea
      packages:
        - name: antrea.2.1.2
          imageBundleUri: projects.registry.vmware.com/tce/packages/antrea:2.1.2
          options: {}
        - name: antrea.2.5.4
          imageBundleUri: projects.registry.vmware.com/tce/packages/antrea:2.5.4
          options: {}
      packageMetadata: |-
        apiVersion: data.packaging.carvel.dev/v1alpha1
        kind: PackageMetadata
        metadata:
          # Must consist of at least three segments separated by a '.'
          # Cannot have a trailing '.'
          name: antrea.vmware.com
          # The namespace this package metadata is available in
          namespace: my-ns
        spec:
          # Human friendly name of the package (optional; string)
          displayName: "Fluent Bit"
          # Long description of the package (optional; string)
          longDescription: "Fluent bit is an open source..."
          # Short desription of the package (optional; string)
          shortDescription: "Log processing and forwarding"
          # Base64 encoded icon (optional; string)
          iconSVGBase64: YXNmZGdlcmdlcg==
          # Name of the entity distributing the package (optional; string)
          providerName: VMware
          # List of maintainer info for the package.
          # Currently only supports the name key. (optional; array of maintner info)
          maintainers:
          - name: "Person 1"
          - name: "Person 2"
          # Classifiers of the package (optional; Array of strings)
          categories:
          - "logging"
          - "daemon-set"
          # Description of the support available for the package (optional; string)
    - name: vsphere-csi
      packages:
        - name: vsphere-csi.4.1.2
          imageBundleUri: projects.registry.vmware.com/tce/packages/vsphere-csi:4.1.2
          options: {}
        - name: vsphere-csi.4.1.5
          imageBundleUri: projects.registry.vmware.com/tce/packages/vsphere-csi:4.1.5
          options: {}
      packageMetadata: |-
        apiVersion: data.packaging.carvel.dev/v1alpha1
        kind: PackageMetadata
        metadata:
          # Must consist of at least three segments separated by a '.'
          # Cannot have a trailing '.'
          name: vsphere-csi.vmware.com
          # The namespace this package metadata is available in
          namespace: my-ns
        spec:
          # Human friendly name of the package (optional; string)
          displayName: "Fluent Bit"
          # Long description of the package (optional; string)
          longDescription: "Fluent bit is an open source..."
          # Short desription of the package (optional; string)
          shortDescription: "Log processing and forwarding"
          # Base64 encoded icon (optional; string)
          iconSVGBase64: YXNmZGdlcmdlcg==
          # Name of the entity distributing the package (optional; string)
          providerName: VMware
          # List of maintainer info for the package.
          # Currently only supports the name key. (optional; array of maintner info)
          maintainers:
          - name: "Person 1"
          - name: "Person 2"
          # Classifiers of the package (optional; Array of strings)
          categories:
          - "logging"
          - "daemon-set"
          # Description of the support available for the package (optional; string)
userPackages:
    - name: fluent-bit
      packages:
        - name: fluent-bit.1.8.2
          imageBundleUri: projects.registry.vmware.com/tce/packages/fluent-bit:1.8.2
          options: {}
        - name: fluent-bit.1.9.0
          imageBundleUri: projects.registry.vmware.com/tce/packages/fluent-bit:1.9.0
          options: {}
      packageMetadata: |-
        apiVersion: data.packaging.carvel.dev/v1alpha1
        kind: PackageMetadata
        metadata:
          # Must consist of at least three segments separated by a '.'
          # Cannot have a trailing '.'
          name: fluent-bit.vmware.com
          # The namespace this package metadata is available in
          namespace: my-ns
        spec:
          # Human friendly name of the package (optional; string)
          displayName: "Fluent Bit"
          # Long description of the package (optional; string)
          longDescription: "Fluent bit is an open source..."
          # Short desription of the package (optional; string)
          shortDescription: "Log processing and forwarding"
          # Base64 encoded icon (optional; string)
          iconSVGBase64: YXNmZGdlcmdlcg==
          # Name of the entity distributing the package (optional; string)
          providerName: VMware
          # List of maintainer info for the package.
          # Currently only supports the name key. (optional; array of maintner info)
          maintainers:
          - name: "Person 1"
          - name: "Person 2"
          # Classifiers of the package (optional; Array of strings)
          categories:
          - "logging"
          - "daemon-set"
          # Description of the support available for the package (optional; string)
    - name: contour
      packages:
        - name: contour.1.1.0
          imageBundleUri: projects.registry.vmware.com/tce/packages/contour:1.1.0
          options: {}
        - name: contour.2.0.0
          imageBundleUri: projects.registry.vmware.com/tce/packages/contour:2.0.0
          options: {}
      packageMetadata: |-
        apiVersion: data.packaging.carvel.dev/v1alpha1
        kind: PackageMetadata
        metadata:
          # Must consist of at least three segments separated by a '.'
          # Cannot have a trailing '.'
          name: contour.vmware.com
          # The namespace this package metadata is available in
          namespace: my-ns
        spec:
          # Human friendly name of the package (optional; string)
          displayName: "Fluent Bit"
          # Long description of the package (optional; string)
          longDescription: "Fluent bit is an open source..."
          # Short desription of the package (optional; string)
          shortDescription: "Log processing and forwarding"
          # Base64 encoded icon (optional; string)
          iconSVGBase64: YXNmZGdlcmdlcg==
          # Name of the entity distributing the package (optional; string)
          providerName: VMware
          # List of maintainer info for the package.
          # Currently only supports the name key. (optional; array of maintner info)
          maintainers:
          - name: "Person 1"
          - name: "Person 2"
          # Classifiers of the package (optional; Array of strings)
          categories:
          - "logging"
          - "daemon-set"
          # Description of the support available for the package (optional; string)
hostImages:
    - name: ami-192385
      provider: aws
      metadata:
        public: "true"
        region: us-west-2
    - name: ami-39402
      provider: aws
      metadata:
        public: "true"
        region: us-east-2
    - name: azi-3992
      provider: azure
      metadata:
        placement: east
        public: "true"
kubernetesMeta:
    version: "1.25"
    components:
        - name: etcd
          repository: projects.registry.vmware.com/tce/etcd
          tag: v3.5.0
        - name: pause
          repository: projects.registry.vmware.com/tce/pause
          tag: v3.4.1
        - name: etcd
          repository: projects.registry.vmware.com/tce/etcd
          tag: v1.8.0
```

### 2. Build Library

### 3. Host Images

Out-of-band from the TKr creation, the TCE project will build upstream
**host** images using
[image-builder](https://github.com/kubernetes-sigs/image-builder). This will
kick-off using our own automation. Scripts and code for how we assemble images
will be made available in the `build/` directory.

#### 3.1 Building Images

#### 3.2 Publishing Images

### 4. Container Images

In order for this proposal to work, package authors **must** copy all upstream
container images to
`projects.registry.vmware.com/tce/packages/${UPSTREAM_PROJECT}/${IMAGE_NAME}:${IMAGE_VERSION}`
and update the imagelock file of the package to point to this new location.

#### 4.1 Copying Container Images

Package authors should copy using the `crane` toolset. Consider the following
example, representing the release of cert-manager `v1.6.1`. The relevant images
are:

```sh
wget -q -O -
https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml | grep -i image:
          image: "quay.io/jetstack/cert-manager-cainjector:v1.6.1"
          image: "quay.io/jetstack/cert-manager-controller:v1.6.1"
          image: "quay.io/jetstack/cert-manager-webhook:v1.6.1"
```

For each image, the container image is pushed to the same URI, prefixed with
`projects.registry.vmware.com/tce/packages`. The following is an example of the
`cert-manager-controller` image.

```sh
$ crane copy quay.io/jetstack/cert-manager-controller:v1.6.1 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller:v1.6.1

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

> `crane` is used instead of `imgpkg` because `crane` maintains the digest
> value.  it also copies **all** architectures over, so in the case of arm64,
> this is made available in the copy. The issue requesting this functionality in
> `imgpkg` is [here](https://github.com/vmware-tanzu/carvel-imgpkg/issues/310).

Using the above `copy`, the architectures are retained and digests (SHA) are
retained.

```sh
$ crane manifest projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller:v1.6.1 | grep -i digest

 "digest": "sha256:41917b5d23b4abe3f5c34a156b1554e49e41185431361af46640580e4d6258fc",
 "digest": "sha256:e2be0d9dfa684e1abf5ef9b24b601b1ca6b9dd6d725342b13c18b44156518b49",
 "digest": "sha256:8898cc51a41a7848076cd7735e5a86feee734f13e802c563ef1deaafe6685040",
 "digest": "sha256:7a60aca7f3c33e58f722229a139514b24cee45881b4c39428ae3cc252ef3190d",
 "digest": "sha256:1faa4c99e61db1e2227ca074de4e40c4e9008335f009fd6fd139c07ac4d5024b",

$ crane manifest quay.io/jetstack/cert-manager-controller:v1.6.1 | grep -i digest
 "digest": "sha256:41917b5d23b4abe3f5c34a156b1554e49e41185431361af46640580e4d6258fc",
 "digest": "sha256:e2be0d9dfa684e1abf5ef9b24b601b1ca6b9dd6d725342b13c18b44156518b49",
 "digest": "sha256:8898cc51a41a7848076cd7735e5a86feee734f13e802c563ef1deaafe6685040",
 "digest": "sha256:7a60aca7f3c33e58f722229a139514b24cee45881b4c39428ae3cc252ef3190d",
 "digest": "sha256:1faa4c99e61db1e2227ca074de4e40c4e9008335f009fd6fd139c07ac4d5024b",
```

In some special cases, such as usage of a container base image with licensing
issues, the TCE team may require package authors to custom build container
images rather than use this copy approach.

### 5. Core Package Repository

Based on the `corePackages` in the build meta, the build library should be able
to produce a compliant list of packages, bundled up such that they can be
referenced via a
[PackageRepository](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-repository)
CRD in a cluster.

To see what this end-state should look like, view the existing `0.9.1`
repository.

```
$ crane export projects.registry.vmware.com/tce/main:0.9.1 - | tar xv
.
.imgpkg
.imgpkg/images.yml
packages
packages/packages.yaml
```

In the above, `packages/packages.yaml` contains many
[PackageMetadata](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-metadata)
and [Package](https://carvel.dev/kapp-controller/docs/latest/packaging/#package)
objects. The generation of these object can be entirely sourced from the
contents of `corePackageRepo` in the build metadata.

### 6. User-managed Package Repository

Producing the user-managed package repository is the same as the core package
repository. However, it sources its packages from the `userManagedRepo` section
in the build metadata.

### 7. Tanzu Kubernetes Release

### Producing Host Images

< TODO >

### Signing

As part of this work, we will **not** sign our machine images, container images,
packages, or package repositories. CLI binaries will continue to be signed to
ensure Mac and Windows users are not prompted to approve usage.

### Release Automation

In order to facilitate creation of compliant Tanzu Kubernetes Release (TKr) files, automation is needed to translate the data from the TCE release build manifest (described in the next section) into the TKr format.

This tooling will be created as a Go library along with a command line utility.
The code will initially reside in the Tanzu Community Edition git repository for convenience, but is a strong candidate for being moved into its own repository.
A Go library will be useful to others seeking to build further tooling on top of TCE, and a command line utility will expose the functionality so that CI/CD processes may utilize it.

#### Go API

At the top level, the  API for creating a TKr from a TanzuBuild will look like the following Go functions.
The `TanzuBuild` Go struct will also be defined by the API package.

```go
import framework "github.com/tanzu-framework/apis/run/v1alpha1"
// TODO:nrb - any validation of values? Conformance to URIs for container images? SemVer enforcement is probably limiting in case upstream images don't use it.

// ReadManifest will parse a TanzuBuild struct from a given YAML file.
func ReadManifest(filePath string) (TanzuBuild, error)

// TranslateToTKR will construct a [framework.TanzuKubernetesRelease](https://github.com/vmware-tanzu/tanzu-framework/blob/main/apis/run/v1alpha1/tanzukubernetesrelease_types.go#L97) struct from a TanzuBuild struct.
// The primary work here is to map relevant fields from a TanzuBuild into a TanzuKubernetesRelease, with little to no manipulation of values.
func TranslateToTKR(manifest TanzuBuild) (framework.TanzuKubernetesRelease, error)

// WriteTKR will output a framework.TanzuKubernetesRelease struct as YAML into a provided io.Writer.
func WriteTKR(tkr framework.TKr, out io.Writer)
```

#### CLI tooling

The CLI tool will be named `tkrgen` and provide a thin wrapper around the Go library.

Usage:

```shell
  tkrgen tce-manifest.yaml -o tce-tkr.yaml
```

#### TCE Release Build Manifest

A build manifest is used to:

* Create a core package repository
* Create a user-managed package repository
* Push the core package repository to an OCI repository
* Push the user-managed package repository to an OCI repository
* Push (upstream) OS images to relevant cloud providers
* Generate Tanzu Kubernetes Release (TKr)
* Push the TKr to an OCI repository


#### Assemble the core package repository

### Ownership

This proposal covers several aspects that bring a community-edition release
together. There are 3 groups involved in this process.

* Package Maintainers
* Tanzu Community Edition Maintainers
* Tanzu Framework Maintainers

The follow diagram categorizes the units of work into the above owners.

![Grouping of tasks based on their owners](imgs-2717/task-owners.png)

## Compatibility

< TODO >

< If this change impacts compatibility of previous versions of TCE or software
integrated with TCE, please call it out here. If incompatibilities can be
mitigated, please add it here. >

## Alternatives Considered

< TODO >

< If alternatives were considered, please add details here >
