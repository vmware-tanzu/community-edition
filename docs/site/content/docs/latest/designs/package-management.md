# Tanzu Package Management

This document covers the management of packages from the server-side and client-side
perspective in Tanzu Community Edition. This is a working design doc that will evolve over time as our
package management is implemented and enhanced.

## Server Side

This section describes the server-side management of extensions. This
specifically focuses on `kapp-controller` and the related Packaging APIs.

### Overview and APIs

Tanzu Community Edition offers package management using the [Carvel Packaging
APIs](https://carvel.dev/kapp-controller/docs/latest/packaging). The primary
APIs are as follows.

* [Package](https://carvel.dev/kapp-controller/docs/latest/packaging/#package):
  Contains metadata about a package and the OCI bundle that satisfies it.
  Typically the OCI bundle referenced is an
  [imgpkg](https://carvel.dev/imgpkg/docs/latest) bundle of configuration. A
  Package is eventually bundled in a `PackageRepository`.
* [PackageRepository](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-repositories):
  A bundle of `Package`s. The bundle is created using `imgpkg` and pushed up to an
  OCI repository. `kapp-controller` watches the object and makes packages
  available in the cluster.
* [IntalledPackage](https://carvel.dev/kapp-controller/docs/latest/packaging/#installedpackage-cr):
  Intent in the cluster to install a package. The is applied by a client-side tool
  such as `tanzu` or `kubectl` CLIs. An `InstalledPackage` references a `Package`.

### Package Creation

`Package`s are created by `kapp-controller` based on the contents of a
`PackageRepository`. A `PackageRepository` points to an OCI bundle of `Package`
manifests. Inclusion of a `PackageRepository` in a cluster inherently makes the
`Package`s available. This can be visually represented as follows.

![Package and Package Repository](/docs/img/tanzu-carvel-new-apis.png)

With the above in place, a user can see all the packages using `kubectl`.

```sh
$ kubectl get packages
NAME                              PUBLIC-NAME            VERSION      AGE
pkg.test.carvel.dev.1.0.0         pkg.test.carvel.dev   1.0.0        7s
pkg.test.carvel.dev.2.0.0         pkg.test.carvel.dev   2.0.0        7s
pkg.test.carvel.dev.3.0.0-rc.1    pkg.test.carvel.dev   3.0.0-rc.1   7s
```

### Package Installation (Server Side)

To install a `Package`, an `InstalledPackage` resource is applied to the
cluster. This instructs `kapp-controller` to lookup the configuration bundle
referenced in the corresponding `Package` manifest. It then downloads those
assets, renders them (e.g. `ytt`) and applies them to the cluster. This can be
visually represented as follows.

![InstalledPackage Flow](/docs/img/tanzu-carvel-installed-package.png)

## Client Side

This section describes the client-side management of extensions. This
specifically focuses on our usage of `tanzu` CLI to discover, configure, deploy,
and manage packages.

### Package Discovery

The `tanzu` CLI is able to discover packages known to the cluster. It discovers
these packages by viewing all [Package
CRs](https://carvel.dev/kapp-controller/docs/latest/packaging/#package) that are available to the
cluster. These Packages can be sourced from 1 or many
[PackageRepository](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-repositories)
CRs. This model is inspired by popular Linux package managers where inclusion of
a repository (e.g. `/etc/apt/sources.list`)  will make new packages available to
the manager (e.g. `apt`). With this, a command such as the following is
possible.

```sh
tanzu package available list
| Retrieving available packages...
  NAME                                           DISPLAY-NAME        SHORT-DESCRIPTION
  cert-manager.community.tanzu.vmware.com        cert-manager        Certificate management
  contour-operator.community.tanzu.vmware.com    contour-operator    Layer 7 Ingress
```

In the above, the `tanzu` CLI is aggregating and listing metadata from
already-existent objects. Namely the following from each `Package` instance:

* `NAME`: `spec.publicName`
* `DISPLAY-NAME`: `spec.version`
* `SHORT-DESCRIPTION`: `spec.`

This is visually represented as follows.

![tanzu package list](/docs/img/tanzu-package-list.png)

### Package Configuration

The `tanzu` CLI is able to configure packages before installation.
This is achieved by inspecting the `--values-schema` for each package
and by providing a `--values-file` of YAML values to be configured.
Alternatively, one may look up the `config/value.yaml` file embedded in each package on [GitHub](https://github.com/vmware-tanzu/community-edition).
Consider the following means of capturing the configuration.

a. Retrieve the list of possible values using the `--values-schema` flag.

```text
tanzu package available get contour.community.tanzu.vmware.com/1.17.1 --values-schema

| Retrieving package details for contour.community.tanzu.vmware.com/1.17.1...
  KEY                                  DEFAULT         TYPE     DESCRIPTION
  contour.logLevel                     info            string   The Contour log level. Valid options are info and debug.
  contour.replicas                     2               integer  How many Contour pod replicas to have.
  contour.useProxyProtocol             false           boolean  Whether to enable PROXY protocol for all Envoy listeners.
  contour.configFileContents           <nil>           object   The YAML contents of the Contour config file. See https://projectcontour.io/docs/v1.17.1/configuration/#configuration-file for more information.
  envoy.logLevel                       info            string   The Envoy log level.
  envoy.service.annotations            <nil>           object   Annotations to set on the Envoy service.
  envoy.service.externalTrafficPolicy  Local           string   The external traffic policy for the Envoy service.
  envoy.service.nodePorts.http         <nil>           integer  If type == NodePort, the node port number to expose Envoy's HTTP listener on. If not specified, a node port will be auto-assigned by Kubernetes.
  envoy.service.nodePorts.https        <nil>           integer  If type == NodePort, the node port number to expose Envoy's HTTPS listener on. If not specified, a node port will be auto-assigned by Kubernetes.
  envoy.service.type                   LoadBalancer    string   The type of Kubernetes service to provision for Envoy.
  envoy.terminationGracePeriodSeconds  300             integer  The termination grace period, in seconds, for the Envoy pods.
  envoy.hostNetwork                    false           boolean  Whether to enable host networking for the Envoy pods.
  envoy.hostPorts.http                 80              integer  If enable == true, the host port number to expose Envoy's HTTP listener on.
  envoy.hostPorts.https                443             integer  If enable == true, the host port number to expose Envoy's HTTPS listener on.
  envoy.hostPorts.enable               false           boolean  Whether to enable host ports. If false, http and https are ignored.
  namespace                            projectcontour  string   The namespace in which to deploy Contour and Envoy.
  certificates.renewBefore             360h            string   If using cert-manager, how long before expiration the certificates should be renewed. If useCertManager is false, this field is ignored.
  certificates.useCertManager          false           boolean  Whether to use cert-manager to provision TLS certificates for securing communication between Contour and Envoy. If false, the upstream Contour certgen job will be used to provision certificates. If true, the cert-manager addon must be installed in the cluster.
  certificates.duration                8760h           string   If using cert-manager, how long the certificates should be valid for. If useCertManager is false, this field is ignored.
```

* `KEY` denotes the yaml access key that can be populated in a `values.yaml` file. Nested keys are denoted by a `.`
* `DEFAULT` denotes the default value if not configured in a provided `values.yaml` file
* `TYPE` tells us what `yaml` type is expected from the key/value pair in a `values.yaml` file
* `DESCRIPTION` is a short hand description of what the value configures for the package

b. Create a `values.yaml` file and define the `namespace` and `logLevel` values based on the `--values-schema`.

```yaml
namespace: custom-namespace
contour:
  logLevel: debug
```

c. Apply the `value.yaml` file during installation.

```sh
tanzu package install contour \
  --package-name contour.community.tanzu.vmware.com \
  --version 1.17.1 \
  --values-file values.yaml
```

> Note: Value files are expected to use `ytt` syntax. Please refer to the [Carvel `ytt` documentation for further details](https://carvel.dev/ytt/docs/latest/)

### Package Installation

The `tanzu` CLI is able to install packages. Where installation is the
declaration of intent to deploy workloads into the cluster. This is achieved
using the [InstalledPackage
CR](https://carvel.dev/kapp-controller/docs/latest/packaging/#installedpackage-cr).
Once present in the cluster, kapp-controller (server side) is able to resolve
what resources need to be created and reconciled to create them.

The primary work on the `tanzu` CLI side is to translate a `Package` CR into an
`InstalledPackage` based on the user's desire. Consider the following available
packages in a cluster.

```sh
tanzu package available list
| Retrieving available packages...
  NAME                                           DISPLAY-NAME        SHORT-DESCRIPTION
  cert-manager.community.tanzu.vmware.com        cert-manager        Certificate management
  contour-operator.community.tanzu.vmware.com    contour-operator    Layer 7 Ingress
  contour.community.tanzu.vmware.com             Contour             An ingress controller
  external-dns.community.tanzu.vmware.com        external-dns        This package provides DNS synchronization functionality.
  fluent-bit.community.tanzu.vmware.com          fluent-bit          Fluent Bit is a fast Log Processor and Forwarder
  gatekeeper.community.tanzu.vmware.com          gatekeeper          policy management
  grafana.community.tanzu.vmware.com             grafana             Visualization and analytics software
  harbor.community.tanzu.vmware.com              Harbor              OCI Registry
  knative-serving.community.tanzu.vmware.com     knative-serving     Knative Serving builds on Kubernetes to support deploying and serving of applications and functions as serverless containers
  local-path-storage.community.tanzu.vmware.com  local-path-storage  This package provides local path node storage and primarily supports RWO AccessMode.
  multus-cni.community.tanzu.vmware.com          multus-cni          This package provides the ability for enabling attaching multiple network interfaces to pods in Kubernetes
  prometheus.community.tanzu.vmware.com          prometheus          A time series database for your metrics
  velero.community.tanzu.vmware.com              velero              Disaster recovery capabilities
```

Should the user want to install the `knative-serving.tce.vmware.com:0.21.0-vmware0` package,
they can run the following command.

```sh
tanzu package install knative-serving.tce.vmware.com

Looking up package to install: knative-serving.tce.vmware.com:
Installed package in default/knative-serving.tce.vmware.com:0.21.0-vmware0
```

In running this command, the `tanzu` client will have done the following.

1. Read the contents of the `knative-serving.tce.vmware.com:0.21.0-vmware0` `Package`.

   The following is an example of what the Package contents might be.

    ```yaml
    apiVersion: package.carvel.dev/v1alpha1
    kind: Package
    metadata:
      # Resource name. Should not be referenced by InstalledPackage.
      # Should only be populated to comply with Kubernetes resource schema.
      # spec.publicName/spec.version fields are primary identifiers
      # used in references from InstalledPackage
      name: knative.tce.vmware.com.0.21.0-vmware0
      # Package is a cluster scoped resource, so no namespace
    spec:
      # Name of the package; Referenced by InstalledPackage (required)
      publicName: knative.tce.vmware.com
      # Package version; Referenced by InstalledPackage;
      # Must be valid semver (required)
      version: 0.21.0-vmware0
      # App template used to create the underlying App CR.
      # See 'App CR Spec' docs for more info
      template:
        spec:
          fetch:
          - imgpkgBundle:
              image: registry.tkg.vmware.run/tkg-knative@sha256:...
          template:
          - ytt:
              paths:
              - config/
          - kbld:
              paths:
              # - must be quoted when included with paths
              - "-"
              - .imgpkg/images.yml
          deploy:
          - kapp: {}
    ```

1. Created a `knative-serving-0-21` `InstalledPackage`.

   The following is an example of what the `InstalledPackage` contents might be.

    ```yaml
    apiVersion: install.package.carvel.dev/v1alpha1
    kind: InstalledPackage
    metadata:
      name: knative-serving-0-21
      namespace: my-ns
    spec:
      # specifies service account that will be used to install underlying package contents
      serviceAccountName: knative-sa
      packageRef:
        # Public name of the package to install. (required)
        publicName: knative.tce.vmware.com
        # Specifies a specific version of a package to install (optional)
        # Either version or versionSelection is required.
        version: 0.21.0-vmware0
        # Selects version of a package based on constraints provided (optional)
        # Either version or versionSelection is required.
        versionSelection:
          # Constraint to limit acceptable versions of a package;
          # Latest version satisying the contraint is chosen;
          # Newly available, acceptable later versions are picked up and installed automatically. (optional)
          constraint: ">0.20"
          # Include prereleases when selecting version. (optional)
          prereleases: {}
    # Populated by the controller
    status:
      packageRef:
        # Kubernetes resource name of the package chosen against the constraints
        name: knative.tce.vmware.com.0.21.0-vmware0
      # Derived from the underlying App's Status
      conditions:
      - type: ValuesSchemaCheckFailed
      - type: ReconcileSucceeded
      - type: ReconcileFailed
      - type: Reconciling
    ```

1. Applied the `InstalledPackage` to the cluster.

This is visually represented as follows.

![tanzu package install](/docs/img/tanzu-package-install.png)

#### Including Package Configuration

A user may wish to bring additional configuration as described in [package
configuration](package-management/#package-configuration)
into their installation. Configuration may be included as a flag during
install `--values-file/-f`. The above example could be run again with the
following.

```sh
tanzu package install knative-serving.tce.vmware.com --values-file knative-serving-config.yaml

Looking up package to install: knative-serving.tce.vmware.com:
Installed package in default/knative-serving.tce.vmware.com:0.21.0-vmware0
```

The implication of including this configuration would do the following.

1. Apply `knative-serving-confg.yaml` as a Kubernetes secret.

    * For example:

    ```yaml
    ---
    apiVersion: v1
    kind: Secret
    metadata:
      name: knative-serving-config
    stringData:
      values.yml: |
        #@data/values
        ---
        hello_msg: "hi"
    ```

1. Add `spec.values` into the `InstalledPackage` CR before applying.

    ```yaml
      # Values to be included in package's templating step
      # (currently only included in the first templating step) (optional)
      values:
      - secretRef:
          name: knative-serving-config
    ```

### Package Repository Discovery

The `tanzu` CLI is able to list all package repositories known to the cluster.
This is essentially a list of all `PackageRespository` objects. The CLI
interaction would look as follows.

```sh
tanzu package repository list

NAME               VERSION
tce-main           1.12
```

In the above, the `tanzu` CLI is aggregating and listing metadata from
already-existent objects. Namely the following from each `PackageRepository` instance:

* `NAME`: `metadata.name`
* `VERSION`: Unknown at this point, will packagerepos be versioned?

### Package Repository Creation

The `tanzu` CLI is able to install package repositories. In turn, all packages
referenced in that repo will be created and made available in the cluster by
`kapp-controller`. A package repository points to an OCI bundle that contains
multiple `Package` manifests. `tanzu` CLI only needs to apply the
`PackageRepository` CR. Theses repository manifests can be found in the Tanzu Community Edition
GitHub repo. The flow could look as follows.

```sh
tanzu package repository install -f ${REPO_MANIFEST_LOCATION}

installed package repo
```

The flow would look as follows.

![tanzu package repo install](/docs/img/tanzu-package-repo-install.png)

### Package Repository Deletion

The `tanzu` CLI is able to delete package repositories. In turn, all packages
referenced in that repo will be deleted  by `kapp-controller`. The flow could
look as follows.

```sh
tanzu package repository delete ${REPO_NAME}

deleted package repo
```

The flow would look as follows.

![tanzu package repo delete](/docs/img/tanzu-package-repo-delete.png)

## Designed Pending Details

This section covers concerns that need design work.

### Upgrading Packages and PackageRepositories

We need a design around how `Package` and `PackageRepsitory` upgrades will work
from a client-side perspective.
