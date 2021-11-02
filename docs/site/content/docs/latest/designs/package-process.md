# Package Process

This document covers the creation of packages for use in Tanzu Community Edition. This is a working design doc that will evolve over time as packages are
implemented. Along with being a design doc, this asset walks you
through the packaging process.

## Terminology

For definitions of extensions, add-ons, core add-ons, user-managed add-ons and
more, see the [Glossary](../installation-planning/#glossary). The packaging details in most of this document are relevant to core and user-managed packages. However, much of the details around discovery, repositories, and CLI interaction are only relevant to user-managed packages.

## Packages

Packaging of external, third party software and functionality is done with the [Carvel](https://carvel.dev/) toolkit.
The end result is an OCI bundle stored in a container registry. For discovery,
deployment, and management operations, the `tanzu` CLI is used, as shown below.

```sh
$ tanzu package install gatekeeper --package-name gatekeeper.community.tanzu.vmware.com --version 3.2.3 --namespace default
| Installing package 'gatekeeper.community.tanzu.vmware.com'
/ Getting package metadata for gatekeeper.community.tanzu.vmware.com
- Creating service account 'gatekeeper-default-sa'
\ Creating cluster admin role 'gatekeeper-default-cluster-role'

- Creating package resource
/ Package install status: Reconciling

 Added installed package 'gatekeeper' in namespace 'default'
```

> This experience is specific to user-managed packages.

For details on how these packages are discovered, deployed, and managed, see
[Package Management](/docs/latest/package-management.md).

### Packaging Workflow

The following flow describes how we package user-managed packages. These steps
are described in detail in the subsequent sections.

![tanzu packaging flow](/docs/img/tanzu-packaging-flow.png)

### 1. Create Directory Structure

Each package lives in a separate directory, named after the package. The
create-package make target will construct the directories and default files. You
can run it by setting `NAME` and `VERSION` variables.

```sh
make create-package NAME=gatekeeper VERSION=3.2.3                                                                                                                                                                                                                                                                                                                          ─╯
mkdir: created directory 'addons/packages/gatekeeper/3.2.3/bundle/.imgpkg'
mkdir: created directory 'addons/packages/gatekeeper/3.2.3/bundle/config/overlay'
mkdir: created directory 'addons/packages/gatekeeper/3.2.3/bundle/config/upstream'

package bootstrapped at addons/packages/gatekeeper/3.2.3
```

The above script creates the following files and directory structure.

```txt
addons/packages/gatekeeper
├── 3.2.3
│   ├── README.md
│   ├── bundle
│   │   ├── .imgpkg
│   │   ├── config
│   │   │   ├── overlay
│   │   │   ├── upstream
│   │   │   └── values.yaml
│   │   └── vendir.yml
│   └── package.yaml
└── metadata.yaml
```

The files and directories are used for the following.

* **README**: Contains the package's documentation.
* **bundle**: Contains the package's imgpkg bundle.
* **bundle/.imgpkg**: Contains metadata for the bundle.
* **bundle/config/upstream**: Contains the package's deployment manifests. Typically
  sourced by upstream.
* **bundle/config/overlay**: Contains the package's overlay applied atop the
  upstream manifest.
* **bundle/config/values.yaml**: User configurable values
* **bundle/vendir.yml**: Defines the location of the upstream resources
* **package.yaml**: Descriptive metadata for the specific version of the package
* **metadata.yaml**: Descriptive metadata for the package

### 2. Add Manifest(s)

In order to stay aligned with upstream, store unmodified manifests. For example,
[gatekeeper](https://github.com/open-policy-agent/gatekeeper)'s upstream
manifest is located
[here](https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/deploy/gatekeeper.yaml).
By storing the configuration of the upstream manifest, you can easily update the
manifest and have customizations applied via
[overlays](https://carvel.dev/ytt/#example:example-overlay).

To ensure integrity of the sourced upstream manifests,
[vendir](https://carvel.dev/vendir/docs/latest/vendir-spec) is used. It will
download and create a lock file that ensures the manifest matches a specific
commit.

In the `bundle` directory, create/modify the `vendir.yml` file. The following
demonstrates the configuration for gatekeeper.

```yaml
apiVersion: vendir.k14s.io/v1alpha1
kind: Config
minimumRequiredVersion: 0.12.0
directories:
  - path: config/upstream
    contents:
      - path: .
        git:
          url: https://github.com/open-policy-agent/gatekeeper
          ref: v3.2.3
        newRootPath: deploy
```

> There are multiple sources you can use. Ideally, packages use either `git` or `githubReleases` such that we can lock in the version. Using the `http` source does not give us the same guarentee as the aforementioned sources.

This configuration means vendir will manage the `config/upstream` directory. To
download the assets and produce a lock file, run the following.

```sh
vendir sync
```

There is also a make task for this.

```sh
make vendir-sync-package PACKAGE=gatekeeper VERSION=3.2.3
```

A lock file will be created at `bundle/vendir.lock.yml`. It will contain the
following lock metadata.

```yaml
apiVersion: vendir.k14s.io/v1alpha1
directories:
  - contents:
      - git:
          commitTitle: Prepare v3.2.3 release (#1084)...
          sha: 15def468c9cbfffc79c6d8e29c484b71713303ae
          tags:
            - v3.2.3
        path: .
    path: config/upstream
kind: LockConfig
```

With the above in place, the directories and files will appear as follows.

```txt
addons/packages/gatekeeper
├── 3.2.3
│   ├── README.md
│   ├── bundle
│   │   ├── .imgpkg
│   │   ├── config
│   │   │   ├── overlays
│   │   │   ├── upstream
│   │   │   │   └── gatekeeper.yaml
│   │   │   └── values.yaml
│   │   ├── vendir.lock.yml
│   │   └── vendir.yml
│   └── package.yaml
└── metadata.yaml
```

### 3. Create Overlay(s)

For each object (e.g. `Deployment`) you need to modify from upstream, an overlay
file should be created. Overlays are used to ensure we import
unmodified-upstream manifests and apply specific configuration on top.

Consider the following `gatekeeper.yaml` added in the previous step.

```yaml
---
#! upstream.yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: gatekeeper-deployment
  labels:
    app: gatekeeper
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gatekeeper
  template:
    metadata:
      labels:
        app: gatekeeper
    spec:
      containers:
      - name: gatekeeper
        image: gatekeeper:1.14.2
        ports:
        - containerPort: 80
```

Assume you want to modify `metadata.labels` to a static value and
`spec.replicas` to a user-configurable value.

Create a file named `overlay-deployment.yaml` in the `bundle/overlay`
directory.

```yaml
---
#! overlay-deployment.yaml

#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind":"Deployment", "metadata":{"name":"gatekeeper-deployment"}})
---
metadata:
  labels:
    #@overlay/match missing_ok=True
    class: gatekeeper
    #@overlay/match missing_ok=True
    owned-by: tanzu

#@overlay/match by=overlay.subset({"kind":"Deployment", "metadata": {"name": "gatekeeper-deployment"}})
---
spec:
  #@overlay/match missing_ok=True
  replicas: #@ data.values.runtime.replicas
```

> ⚠️: Do not templatize or overlay container image fields. `kbld` will be used to
create and/or reference image digest SHAs.

_Detailed overlay documentation is available [in the Carvel
site](https://carvel.dev/ytt/#example:example-overlay)._

### 4. Default Values

For every user-configurable value defined above, a `values.yaml` file should
contain defaults and documentation for what the parameter impacts. If a value
is overriding an upstream value, prefer to use that upstream value. For example,
if the upstream default namespace is `foo-ns`, prefer to use `foo-ns` as the
default setting for the namespace in the values.yaml file.

Create/modify a `values.yaml` file in `bundle/config`.

```yaml
#@data/values
---

#! The namespace in which to deploy gatekeeper.
namespace: gatekeeper

#! The amount of replicas that should exist in gatekeeper.
runtime:
  replicas: 3
```

### [Optional]: Validate Templating

With the above in place, you can validate that overlays and templating are
working as expected. The conceptual flow is as follows.

![templating flow](/docs/img/tanzu-templating-flow.png)

To run the above, you can use `ytt` as follows. If successful, the transformed manifest will be displayed,
otherwise, an error message is displayed.

```sh
ytt -f addons/packages/gatekeeper/3.2.3/bundle/config
```

### 5. Resolve and reference image digests

To ensure integrity of packages, it is important to reference an [image
digest](https://github.com/opencontainers/image-spec/blob/master/descriptor.md#digests)
rather than a tag. A tag's underlying image can change arbitrarily. Whereas
referencing a SHA (via digest) will ensure consistency on every pull.

kbld is used to create a lock file, which we name `images.yml`. This file contains an
`ImagesLock` resource. `ImagesLock` is similar to a
[go.sum](https://golang.org/ref/mod#go). The image field in the source manifests
**are not mutated**. Instead, the SHA will be swapped out for the tag upon
deployment. The relationship is as follows.

![tanzu packaging flow](/docs/img/tanzu-kbld-flow.png)

To find all container image references, create an `ImagesLock`, and ensure the
digest's SHA is referenced, you can run `kbld` as follows.

```sh
kbld --file addons/packages/gatekeeper/3.2.3/bundle \
  --imgpkg-lock-output addons/packages/gatekeeper/3.2.3/bundle/.imgpkg/images.yml
```

There is also a make task for this.

```sh
make lock-package-images PACKAGE=gatekeeper VERSION=3.2.3
```

This will produce the following file `bundle/.imgpkg/images.yml`.

```yaml
---
apiVersion: imgpkg.carvel.dev/v1alpha1
images:
  - annotations:
      kbld.carvel.dev/id: openpolicyagent/gatekeeper:v3.2.3
    image: index.docker.io/openpolicyagent/gatekeeper@sha256:9cd6e864...
kind: ImagesLock
```

By placing this file in `bundle/.imgpkg`, it will not pollute the
`bundle/config` directory and risk being deployed into Kubernetes
clusters. At this point, the following directories and files should be in place.

```txt
addons/packages/gatekeeper
├── 3.2.3
│   ├── README.md
│   ├── bundle
│   │   ├── .imgpkg
│   │   │   └── images.yml
│   │   ├── config
│   │   │   ├── overlays
│   │   │   │   ├── overlay-deployment.yaml
│   │   │   ├── upstream
│   │   │   │   └── gatekeeper.yaml
│   │   │   └── values.yaml
│   │   ├── vendir.lock.yml
│   │   └── vendir.yml
│   └── package.yaml
└── metadata.yaml
```

### 6. Bundle configuration and deploy to registry

All the manifests and configuration are bundled in an OCI-compliant package.
This ensures immutability of configuration upon a release. The bundles are
stored in a container registry.

`imgpkg` is used to create the bundle and push it to the container registry. It
leverages your underlying container registry, so you must set up authentication
on the system you'll create the bundle from (e.g. `docker login`).

To ensure metadata about the package is captured, add the following `Bundle` file
into `bundle/.imgpkg/bundle.yaml`.

```yaml
apiVersion: imgpkg.carvel.dev/v1alpha1
kind: Bundle
metadata:
  name: gatekeeper
authors:
  - name: Joe Engineer
    email: engineerj@example.com
websites:
  - url: github.com/open-policy-agent/gatekeeper
```

The following packages and pushes the bundle.

```sh
imgpkg push \
  --bundle $(OCI_REGISTRY)/gatekeeper/3.2.3:$(BUNDLE_TAG) \
  --file addons/packages/gatekeeper/3.2.3/bundle
```

There is also a make task for this.

```sh
make push-package PACKAGE=gatekeeper VERSION=3.2.3
```

The results of this look as follows. Notice at the end of a successful push, imgpkg reports the URL and digest of
the package. This information will be used in the next step.

```sh
===> pushing gatekeeper/3.2.3
dir: .
dir: .imgpkg
file: .imgpkg/bundle.yml
file: .imgpkg/images.yml
dir: config
dir: config/overlays
file: config/overlays/overlay-deployment.yaml
dir: config/upstream
file: config/upstream/gatekeeper.yaml
file: config/values.yaml
file: vendir.lock.yml
file: vendir.yml
Pushed 'projects.registry.vmware.com/tce/gatekeeper@sha256:b7a21027...'
Succeeded
```

### 7. Create/Modify a Package CR

A `Package` is used to define metadata and templating information about a piece
of software. A `Package` CR is created for every addon and points to the OCI
registry where the `imgpkg` bundle can be found. This file also captures some version specific
information about the package, such as version, license, release notes. The `Package` CR is put into
a directory structure with other packages to eventually form a
`PackageRepository`. The `Package` CR is **not** deployed to the cluster,
instead the `PackageRepsoitory` bundle, containing many `Package`s is. Once
the `PackageRepository` is in place, `kapp-controller` will make `Package` CRs
in the cluster. This relationship can be seen as follows.

![Package and Package Repository](/docs/img/tanzu-carvel-new-apis.png)

An example `Package` for `gatekeeper` would read as follows.

```yaml
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: gatekeeper.community.tanzu.vmware.com.3.2.3
  namespace: gatekeeper
spec:
  refName: gatekeeper.community.tanzu.vmware.com
  version: 3.2.3
  releaseNotes: "gatekeeper 3.2.3 https://github.com/open-policy-agent/gatekeeper/releases/tag/v3.2.3"
  licenses:
    - "Apache 2.0"
  template:
    spec:
      fetch:
        - imgpkgBundle:
            image: projects.registry.vmware.com/tce/gatekeeper@sha256:b7a21027...
      template:
        - ytt:
            paths:
              - config/
        - kbld:
            paths:
              - "-"
              - .imgpkg/images.yml
      deploy:
        - kapp: {}
```

* `metadata.name`: Concatenation of `spec.refName` and `version` (see below).
* `spec.refName`: Name that will show up to consumers. **Must be unique across
  packages**.
* `version`: version number of this package instance, must use semver. The
  version used should reflect the version of the packaged software. For example,
  if `gatekeeper`'s main container image is version `3.2.3`, this
  package should be the same.
* `spec.template.spec.fetch[0].imgpkgBundle.image`: THe URL of the location of this package in an OCI registry.
  This value is obtained from the result of the `imgpkg push` command.

### 8. Package Metadata

The final step in creating a package is to update the `metadata.yaml` file. This file contains general
information about the package overall, not specific to a version. Here is an overview of the types of metadata captured
in the file.

* Display friendly name
* Short and long descriptions.
* Authoring organization
* Maintainers
* Descriptive categories
* SVG logo

At this point, the package has been created, pushed and documented. The package is ready to be deployed to a cluster
as part of a package repository.

### 9. Creating a Package Repository

Tanzu Community Edition maintains a `addons/repos` directory where the main repository definition file, `main.yaml`, is kept.
This file is a simple, custom yaml file defining the specific versions of packages to be included in the repository.
An example of this file is as follows:

```yaml
---
packages:
  - name: gatekeeper
    versions:
      - 3.2.3
```

There is a makefile task, `generate-package-repo`, that generates the package repository from this file. `kapp-controller`
currently expects the package repositories to be in the format of an imgpkgBundle. This task will generate that bundle.
When the task is executed, `make generate-package-repo CHANNEL=main`, the following steps are performed:

* Create `addons/repos/generated/main` directory
* Create `addons/repos/generated/main/.imgpkg` for imgpkg
* Create `addons/repos/generated/main/packages/packages.yaml`
* Iterate over packages and concatenates the package `metadata.yaml` and specific package version's `package.yaml` into the repository's `packages.yaml` file
* Create an imgpkg `images.yml` lock file
* Push the bundle to the OCI Registry.

> The package repository will be tagged `:latest`

Upon successful completion, instructions for installing the package repository to your cluster are shown.

```sh
tanzu package repository add repo-name --namespace default --url projects.registry...
```

Tanzu Community Edition will maintain a `main` repo, but a `beta` or `package-foo` repo could be created for development work or to provide
multiple versions of the `foo` software.

## Common Packaging Considerations

### Preventing kapp-controller from Mutating Resources After Deploy

At times, a resource deployed and managed by kapp-controller may be expectedly
mutated by another process. For example, a configmap may be deployed alongside
an [operator](https://operatorhub.io/what-is-an-operator). When the operator
mutates the configmap, kapp-controller will eventually trigger an update and
refresh the configmap back to its original state.

To prevent this behavior, an annotation is added named
`kapp.k14s.io/update-strategy` set to the value of `skip`. It's likely you'll do
this via an [overlay](#3-create-overlays). Below is an example of how you'd set
this up for an upstream configmap.

Upstream Configmap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-domain
  namespace: knative-serving
  labels:
    serving.knative.dev/release: "v0.20.0"
  annotations:
    knative.dev/example-checksum: "74c3fc6a"
data:
  _example: |
    ################################
    #                              #
    #    EXAMPLE CONFIGURATION     #
    #                              #
    ################################
```

Overlay

```yaml
---
#! overlay-configmap-configdomain.yaml

#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind":"ConfigMap", "metadata":{"name":"config-domain"}})
---
metadata:
  annotations:
    #@overlay/match missing_ok=True
    kapp.k14s.io/update-strategy: skip
```

With the above in place, updates will not cause the `config-domain` ConfigMap to
be mutated.

For more details on this annotation, see the [kapp Apply Ordering
documentation](https://carvel.dev/kapp/docs/latest/apply-ordering).

### Ensuring Order of Deploying Assets

It may be important that your package deploys specific components before others.
For example, you may wish for a Deployment that satisfies a [validating
webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#validatingadmissionwebhook)
to be up before applying a ValidatingWebhookConfiguration. This would ensure
the service that does validation is up and healthy before blocking API traffic
to its endpoint.

To prevent this behavior, the annotations `kapp.k14s.io/change-group` and
`kapp.k14s.io/change-rule` are used. It's likely you'll do this via an
[overlay](#3-create-overlays). Below is an example of how you'd set this up for
an upstream Deployment and ValidatingWebhookConfiguration.

Upstream

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    control-plane: controller-manager
    gatekeeper.sh/operation: audit
    gatekeeper.sh/system: "yes"
  name: gatekeeper-controller-manager
  namespace: gatekeeper-system
  annotations:
spec:
---
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  creationTimestamp: null
  labels:
    gatekeeper.sh/system: "yes"
  name: gatekeeper-validating-webhook-configuration
  annotations:
    # it is very important this resource (ValidatingWebhookConfiguration) is applied
    # last. Otherwise, it can wire up the admission request before components required
    # to satisfy it are deployed.
    kapp.k14s.io/change-group: "tce.gatekeeper/vwc"
    kapp.k14s.io/change-rule: "upsert after upserting tce.gatekeeper/deployment"
webhooks:
```

Overlays

```yaml
---
#! overlay-deployment-gatekeeperaudit.yaml

#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind":"Deployment", "metadata":{"name":"gatekeeper-controller-manager"}})
---
metadata:
  annotations:
    #@overlay/match missing_ok=True
    kapp.k14s.io/change-group: "tce.gatekeeper/deployment"
```

```yaml
---
#! overlay-validatingwebhookconfiguration-gatekeeper.yaml

#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind":"ValidatingWebhookConfiguration", "metadata":{"name":"gatekeeper-validating-webhook-configuration"}})
---
metadata:
  annotations:
    #@overlay/match missing_ok=True
    kapp.k14s.io/change-group: "tce.gatekeeper/vwc"
    #@overlay/match missing_ok=True
    kapp.k14s.io/change-rule: "upsert after upserting tce.gatekeeper/deployment"
```

With the above overlays applied, the ValidatingWebhookConfiguration will not be
applied until the Deployment is healthy.

For more details on this annotation, see the [kapp Apply Ordering
documentation](https://carvel.dev/kapp/docs/latest/apply-ordering).

## Designed Pending Details

This section covers concerns that need design work.

### Versioning of Multiple PackageRepository Instances

With the introduction of the `PackageRepository`, we need to determine how we
are going to handle the ever growing number of package instances
(package+version) that will grow over time.

* Do we maintain a `default` repo with all the latest packages?
* How to we offer older packages?
