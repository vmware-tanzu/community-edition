# Work with Packages

Tanzu Community Edition employs the `tanzu package` command to discover,
configure, and manage packages running in a cluster. A package is a bundle of
software that translates to one or many Kubernetes primitives such as
`Deployment`s, `Service`s, and more. This document covers how to use `tanzu
package` to interact with packages.

## Package Repositories

A package repository holds references to package(s). By installing a package
repository into a cluster, packages become available for installation. A look
at this relationship is as follows.

![tanzu packaging flow](/docs/img/pkg-mgmt-repo.png)

Packages are available in the same namespace where the PackageRepository has been installed.
There is however, a namespace that provides packages globally to the cluster if the PackageRepository
providing them is installed in it. This namespace is `tanzu-package-repo-global`.

The `tanzu-core` package repository will pre-exist on every cluster in the
`tkg-system` namespace. Packages in this repository are exclusively for cluster
bootstrapping. They should **not** be reinstalled by users.

### Adding a Package Repository

To add a package repository to a cluster, run:

```sh
tanzu package repository \
  add ${REPO_NAME} \
  --url ${URL} \
  --namespace ${NS}
```

* `${REPO_NAME}` is the friendly name that will show up in the cluster.
* `${URL}` is the location of the package repository bundle.
  * This must point to a [package repository OCI
    bundle](https://carvel.dev/kapp-controller/docs/latest/package-authoring/#creating-a-package-repository).
* `${NS}` is the Kubernetes namespace to deploy the repository into.
  * This is the namespace packages will be discoverable within. **It does not
    define the target namespace software is eventually run within**. This is typically
    defined as configuration to the package.
  * To make packages available globally in the cluster, use `tanzu-package-repo-global`
    as namespace.
  * If you'd like to create the namespace as part of the command, append
    `--create-namespace`.

#### Example(s)

1. Installing the Tanzu Community Edition package repository globally.
   This will make all the packages in the repository available in every namespace:

    ```sh
    tanzu package repository add tce-repo \
      --url projects.registry.vmware.com/tce/main:0.9.1 \
      --namespace tanzu-package-repo-global
    ```

1. Installing the Tanzu Community Edition package repository to the `default` namespace.
   This will make all the packages in the repository available only in the `default` namespace:

    ```sh
    tanzu package repository add tce-repo \
      --url projects.registry.vmware.com/tce/main:0.9.1
    ```

### Discovering Package Repositories

To discover package repositories installed in a cluster, run:

```sh
tanzu package repository list --all-namespaces
```

Package repositories are namespace scoped, thus without using
`--all-namespaces`, you'll only receive a list of repositories deployed to the
`default` namespace.

### Deleting a Package Repository

To remove an installed package repository, run:

```sh
tanzu package repository delete ${REPO_NAME} --namespace ${NS}
```

* `${REPO_NAME}` is the friendly name used when adding the repository.
* `${NS}` is Kubernetes namespace the repository was installed in.
    > If this is not set, the deletion command assumes the repository is in the
    > `default` namespace.

#### Example(s)

1. Removing the Tanzu Community Edition package repository from the global (`tanzu-package-repo-global`) namespace:

    ```sh
    tanzu package repository delete tce-repo --namespace tanzu-package-repo-global
    ```

1. Removing the Tanzu Community Edition package repository from the `default` namespace:

    ```sh
    tanzu package repository delete tce-repo
    ```

## Packages

Packages hold reference to a configuration bundle. The configuration bundle
provides instructions for how to run the software in a cluster. Source code for
Tanzu Community Edition's configuration bundles can be found in
[GitHub](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages).

![tanzu package install](/docs/img/pkg-mgmt-pkg.png)

When running `tanzu package` commands, there are two types of resources:

* Packages: Definitions of **available** software to install. Represented as a
  `Package` object in the cluster.
* Installed Packages: Declared desire to install an instance of a package.
  Represented as a `PackageInstall` object in the cluster.

You will see the above surfaced when running `tanzu package available` vs `tanzu
package installed`, respectively.

### Discovering Available Packages

To list available packages, run:

```sh
tanzu package available list --namespace ${NS}
```

* `${NS}` is the namespace the package repository was installed in.
    > If this is not set, the list command assumes the repository is in the
    > `default` namespace.

**NOTE**: If the PackageRepository has been installed in the global namespace (`tanzu-package-repo-global`) the
previous command will show the list of available packages on any namespace.

This will return a table featuring a display name, name (fully qualified), and
short description.

To resolve the available versions of a package, run:

```sh
tanzu package available list ${PACKAGE_FQN}
```

* `${PACKAGE_FQN}` is the full name of the package.
  * Typically follows the format of `${PACKAGE_NAME}.community.tanzu.vmware.com`

To install a package, the `${PACKAGE_FQN}` and `${PACKAGE_VERSION}` is required.

#### Example(s)

1. Getting a list of all packages when PackageRepository has been installed in the global
   namespace (`tanzu-package-repo-global`). You can provide a `--namespace` argument to this
   command but results should be the same.

    ```sh
    tanzu package available list
    ```

1. Getting a list of all packages from package repositories in the `my-apps`
   namespace.

    ```sh
    tanzu package available list --namespace my-apps
    ```

1. Getting the versions of `external-dns` available from the `external-dns` package in the
   `default` namespace.

    ```sh
    tanzu package available list external-dns.community.tanzu.vmware.com
    ```

### Installing a Package

A package may have its own unique installation steps or requirements, and may have dependencies on
other software, for example, Contour has a dependency on Cert Manager. Before installing a package, be sure to review its documentation. Documentation for each package can be found in the left navigation (`Packages > ${PACKAGE_NAME}`) of this site.

To install a package, run:

```sh
tanzu package install ${NAME} \
  --package-name ${PACKAGE_FQN} \
  --version ${PACKAGE_VERSION} \
  --namespace ${NS}
```

* `${NAME}` is the friendly name of the installed software.
* `${PACKAGE_FQN}` is the full name of the package.
  * Typically follows the format of `${PACKAGE_NAME}.community.tanzu.vmware.com`
* `${PACKAGE_VERSION}` is the semantic version of the package to deploy.
  * Available versions can be retrieved via `tanzu package available list
    ${PACKAGE_FQN}`.
* `${NS}` is the namespace the package can be located in. It also will determine
  where the `PackageInstall` Kubernetes object is placed.
  * When Packages are available globally, they will be available in every namespace.
  * **This does not determine which namespace(s) the software will run in.**

At the point of install, there are multiple objects that may exist in different
Kubernetes namespaces. The breakdown of how objects end up in different namespaces is as
follows.

![tanzu package namespace](/docs/img/pkg-mgmt-ns.png)

The `Package`s, or software available for install, are always available in the
**same** namespace as the `PackageRepository`. When you run a `tanzu package
install`, you must specify the same namespace as the `Package`, unless the
`PackageRepository` is setup to be a global package repository. [To understand
global package repositories and package namespacing in general, visit the
kapp-controller
documentation](https://carvel.dev/kapp-controller/docs/latest/package-consumer-concepts/#namespacing).
When it comes time for `kapp-controller` to install the software representing
the package's configuration bundle, the namespace declaration in the rendered
manifests is used to determine where the software is installed. See the
Configuring a Package section below for an example on how to customize this.

### Configuring a Package

Packages usually offer configuration to customize how their software is deployed in
a cluster. To understand the configuration for a package, visit its
documentation. Documentation for each package can be found in the left navigation (`Packages >
${PACKAGE_NAME}`) of this site.

Along with a package's documentation, a values file is stored along with the
package bundle for ease of access. Package configurations are available at
`github.com/vmware-tanzu/community-edition/tree/main/addons/packages/${PACKAGE_NAME}/${PACKAGE_VERSION}/bundle/config/values.yaml`. You can download the `values.yaml` file and customize it. For example, see this [Prometheus values file](https://github.com/vmware-tanzu/community-edition/blob/4b1a206e44588cf097e388d2ce2a354433389cb3/addons/packages/prometheus/2.27.0/bundle/config/values.yaml).

You can also get a description of the possible configuration for a package in a human readable format:

```sh
tanzu package available get ${PACKAGE_FQN}/${PACKAGE_VERSION} \
  --namespace ${NS} \
  --values-schema
```

* `${PACKAGE_FQN}` is the full name of the package.
  * Typically follows the format of `${PACKAGE_NAME}.community.tanzu.vmware.com`
* `${PACKAGE_VERSION}` is the semantic version of the package to deploy.
  * Available versions can be retrieved via `tanzu package available list
    ${PACKAGE_FQN}`.
* `${NS}` is the namespace the package can be located in.

#### Example(s)

1. Getting a configuration for package `contour.community.tanzu.vmware.com` version `1.17.1`.

    ```sh
    tanzu package available get contour.community.tanzu.vmware.com/1.17.1 --values-schema
    ```

To install a package with a customized configuration file, run:

```sh
tanzu package install ${NAME} \
  --package-name ${PACKAGE_FQN} \
  --values-file ${VALUES_FILE_PATH} \
  --version ${PACKAGE_VERSION} \
  --namespace ${NS}
```

* `${NAME}` is the friendly name of the installed software.
* `${PACKAGE_FQN}` is the full name of the package.
  * Typically follows the format of `${PACKAGE_NAME}.community.tanzu.vmware.com`
* `${PACKAGE_VERSION}` is the semantic version of the package to deploy.
  * Available versions can be retrieved via `tanzu package available list
    ${PACKAGE_FQN}`.
* `${VALUES_FILE_PATH}` is the location of the values (`.yaml`) file for
  customizing the package.
* `${NS}` is the namespace the package can be located in. It also will determine
  where the `PackageInstall` Kubernetes object is placed.
  * **This does not determine which namespace(s) the software will run in**. This is typically
    defined as configuration to the package.

#### Example

1. Customizing the namespace a package's (`contour`) software will run in.

    a. Retrieve the list of possible values using the `--values-schema` flag.

    ```text
    tanzu pacakage available get contour.community.tanzu.vmware.com/1.17.1 --values-schema

    | Retrieving package details for contour.community.tanzu.vmware.com/1.17.1...
      KEY                                  DEFAULT         TYPE     DESCRIPTION
      contour.logLevel                     info            string   The Contour log level. Valid options are info and debug.
      contour.replicas                     2               integer  How many Contour pod replicas to have.
      contour.useProxyProtocol             false           boolean  Whether to enable PROXY protocol for all Envoy listeners.
      contour.configFileContents           <nil>           object   The YAML contents of the Contour config file. See https://projectcontour.io/docs/v1.17.1/configuration/#configuration-file for more information.
      envoy.logLevel                       info            string   The Envoy log level.
      envoy.service.annotations            <nil>           object   Annotations to set on the Envoy service.
      envoy.service.externalTrafficPolicy  Local           string   The external traffic policy for the Envoy service.
      envoy.service.loadBalancerIP         <nil>           string   If type == LoadBalancer, the desired load balancer IP for the Envoy service.
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

    > Note: Value files are expected to be plain YAML and not YTT Data Values files.

### Listing Installed Packages

To list installed packages, run:

```sh
tanzu package installed list --namespace ${NS}
```

* `${NS}` is the namespace the `PackageInstall` was added to.
  * **This is not the namespace(s) the software is running in.**

### Deleting a Package

To delete an installed package, run

```sh
tanzu package installed delete ${NAME} --namespace ${NS}
```

* `${NAME}` is the name of the package installed.
* `${NS}` is the namespace the package can be located in. Specifically, it where
  the `PackageInstall` object was created.
  * This is not always the namespace where the software is running.

When deleting a package, this removes the `PackageInstall` object. From there,
the components that make up the software will be terminated by
`kapp-controller`.

#### Example(s)

1. Deleting the installed package `contour-teamb` from the `default` namespace.

    ```sh
    tanzu package installed delete contour-teamb
    ```

### Updating a Package

To update an installed package, run

```sh
tanzu package installed update ${NAME} \
  --package-name ${PACKAGE_FQN} \
  --values-file ${VALUES_FILE_PATH} \
  --version ${PACKAGE_VERSION} \
  --namespace ${NS}
```

* `${NAME}` is the friendly name of the installed software.
* `${PACKAGE_FQN}` is the full name of the package.
  * Typically follows the format of `${PACKAGE_NAME}.community.tanzu.vmware.com`
* `${PACKAGE_VERSION}` is the semantic version of the package to deploy.
  * Available versions can be retrieved via `tanzu package available list
    ${PACKAGE_FQN}`.
* `${VALUES_FILE_PATH}` is the location of the values (`.yaml`) file for
  customizing the package.
* `${NS}` is the namespace the package can be located in. It also will determine
  where the `PackageInstall` Kubernetes object is placed.
  * **This does not determine which namespace(s) the software will run in**. This is typically
    defined as configuration to the package.

The package will be updated to satisfy the new arguments provided. Arguments are optional and only those
that will need to change need to be provided.

#### Example(s)

1. Updating the installed package `contour-teamb` from the `default` namespace to a newer version
   and updated configuration.

    ```sh
    tanzu package installed update contour-teamb --version 1.17.2 --values-file contour-config.yaml
    ```

## Troubleshooting Packages

This section covers common ways to troubleshoot packages. Before reading, review
the following diagram that shows what is created when a package is installed.

![tanzu packaging troubleshooting](/docs/img/pkg-mgt-trbl.png)

### Installation Troubleshooting

If packages are failing to install, you can view the following resources:

```sh
kubectl get packageinstall,app --namespace ${NS}
```

* `packageinstall` is the object created by `tanzu package install`, it declares
  intent to install a package.
* `app` is [a CRD](https://carvel.dev/kapp-controller/docs/latest/walkthrough/)
  the defines installation details around software `kapp-controller` should
install.
  * This is an object you'll typically not mutate, it's a lower-level detail
    helpful for troubleshooting.

Running a `kubectl describe` against the objects returned from the above may
also provide some helpful troubleshooting details.

### Configuration Troubleshooting

Package configurations are rendered by `kapp-controller` to produce manifests
that are installed in the cluster. If you're having trouble understanding why
manifests are being deployed as they are, you can render manifests locally.

To do this, you can:

1. Clone the community-edition repository.

    ```sh
    git clone github.com/vmware-tanzu/community-edition
    ```

1. Open the directory of a specific package version, for example, contour.

    ```sh
    cd community-edition/addons/packages/contour/1.17.1
    ```

1. You can now render the manifest using the [ytt tool](https://carvel.dev/ytt/docs/latest).

    ```sh
    ytt -f bundle/config
    ```

    > Additionally, you can provide a values file (a sample `values.yaml` is in
    > `bundle/config/values.yaml`) by appending the -f flag.
