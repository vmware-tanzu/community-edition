# Unmanaged Clusters Reference

This is the reference documentation for the `unmanaged-cluster` plugin. If
this is your first time deploying an unmanaged cluster, see
[Getting Started with Unmanaged Clusters](getting-started-unmanaged/).

## Command aliases

As a Tanzu CLI plugin, unmanaged clusters can be interacted with using `tanzu
unmanaged-cluster`. `uc`, `um`, and `unmanaged` are all valid aliases for the
command. This means the following commands equate to the same:

* `tanzu unmanaged-cluster create hello`
* `tanzu uc create hello`
* `tanzu um create hello`
* `tanzu unmanaged create hello`

## Create clusters

`create` is used to create a new cluster. By default, it:

1. Installs a cluster using the `kind` provider.
1. Installs `kapp-controller`.
1. Installs a core package repository.
1. Installs a user-managed package repository.
1. Installs a CNI package.
    * defaults to `calico`.
1. Sets your kubeconfig context to the newly created cluster.

To create a cluster, run:

```sh
tanzu unmanaged-cluster create ${CLUSTER_NAME}
```

## Use a different cluster provider

`create` supports the `--provider` flag or `Provider` configuration option (if you are using a custom configuration yaml file)
which sets the cluster bootstrapping provider.

The following providers are supported:

* `kind`: _Default provider._ A tool for running local Kubernetes clusters using Docker container “nodes”. [Documentation site.](https://kind.sigs.k8s.io/)
* `minikube`: Local Kubernetes, focusing on making it easy to learn and develop for Kubernetes. Supports container and virtual machine managers. [Documentation site.](https://minikube.sigs.k8s.io/docs/)

_Note:_ In order to use the `minikube` provider, you first must install minikube to your system.
[Read the minikube "Start" page](https://minikube.sigs.k8s.io/docs/start/) to learn how to install and get going with minikube.

## Deploy multi-node clusters

`create` supports `--control-plane-node-count`
and `--worker-node-count` to create multi-node clusters in a supported provider.

_Note:_ The `kind` provider does _not_ support deploying multiple control planes
with no worker nodes. For this type of granular configuration, see [Customize cluster provider.](#customize-cluster-provider)

_Note:_ The `minikube` provider does _not_ support deploying multiple control planes.

* To deploy an unmanaged cluster with 5 total nodes
using the default `kind` provider:

  ```sh
  tanzu unmanaged-cluster create --control-plane-node-count 2 --worker-node-count 3
  ```

## Install additional package repository

`create` supports the `--additional-repo` flag to automatically install package repositories
during cluster bootstrapping. This flag may be specified multiple times to install multiple repositories. Values should be valid registry URLs that point to package repositories.

By default, if you do not specify the `--additional-repo` flag, the default Tanzu Community Edition package repository is installed.

If the `--additional-repo` flag is provided, the default package repository will **not** be installed.

* To deploy a cluster with an additional package repo:

  ```sh
  tanzu unmanaged-cluster create --additional-repo my-repo.registry-url.com/path
  ```

Note: Steps for  deploying a package (after you have created an unmanaged cluster ) from the default package repository are provided in the [Getting Started Guide](getting-started-unmanaged/#deploy-a-package).

## Install packages

_Warning:_ Installing packages during bootstrapping is an experimental feature. Use with caution.

`create` supports the `--install-package` flag to automatically install a package from a package repository.
The name of the package must be the fully qualified name of the package in the package repository, or a prefix of the package name in the package repository.

* To install the latest version of `fluent-bit` with default values during cluster creation:

  ```sh
  tanzu unmanaged-cluster create my-cluster --install-package fluent-bit
  ```

* To designate a package version or install a package with a customized configuration file, use a mapping. The expected format is:

    `text
    name:version:config-file-path`

    ```sh
    tanzu unmanged-cluster create --install-package fluent-bit:1.7.5:path-to-my-config.yaml
    ```

  Both `version` and `config-file-path` are optional.

* To install the most recent package, use the keyword `latest` or an empty string for the version
to select the most recent semantic version for the specified package. Example:

  ```sh
  tanzu unmanaged-cluster create my-cluster --install-package external-dns:latest:path-to-my-config.yaml
  ```

## Install multiple packages

_Warning:_ Installing packages during bootstrapping is an experimental feature. Use with caution.

Install-package mappings may be specified multiple times via multiple `--install-package` flags
or within a single `--install-package` flag, delimited by a comma:

* To install the fluent-bit package at version 1.7.5 with no values yaml file
and the external-dns package at the latest version with a values yaml file.

  ```sh
  tanzu unmanaged-cluster create my-cluster --install-package fluent-bit:1.7.5,external-dns::path-to-my-config.yaml
  ```

For the most granularity and configurability, you can configure all options via an unmanaged cluster configuration file that includes the `--install-package` flag and mappings.  

* Generate a configuration file: For more information see [Create a configuration file](#custom-configuration) below.

  ```sh
  tanzu unmanaged-cluster create my-cluster -f my-config.yaml
  ```

The following example truncated configuration file will create a cluster with three packages automatically installed:

* fluent-bit at the latest version with default values
* external-dns at version 0.10.0 with the default values
* app-toolkit at the latest version configured with the provided values

```yaml
InstallPackages:
- name: fluent-bit.community.tanzu.vmware.com
- name: external-dns.community.tanzu.vmware.com
  config: external-values.yaml
  version: 0.10.0
- name: app-toolkit.community.tanzu.vmware.com
  config: values.yaml
```

**Note:** A package may have unique installation steps or requirements, and may have dependencies on other software, for example, Contour has a dependency on Cert Manager. Before installing a package, be sure to review its documentation. Documentation for each package can be found in the left navigation (Packages > ${PACKAGE_NAME}) of this site.

## List clusters

`list` or `ls` is used to list all known clusters.

* To list known clusters, run:

  ```sh
  tanzu unmanaged-cluster list
  ```

## Delete clusters

`delete` or `rm` is used to delete a cluster. It will:

1. Attempt to delete the cluster based on the provider.
    * by default, clusters use `kind`, this will delete the `kind` cluster.
1. Attempt to remove the cluster's directory.
    * located at `~/.config/tanzu/tkg/unmanaged/${CLUSTER_NAME}/`.

* To delete a cluster, run:

  ```sh
  tanzu unmanaged-cluster delete ${CLUSTER_NAME}
  ```

## Create a configuration file

`configure`, `config`, or `conf` creates a configuration file for cluster creation:

* To create a configuration file to modify how
clusters are created. :

  ```sh
  tanzu unmanaged-cluster configure ${CLUSTER_NAME}
  ```

  The final configuration file is available here:

  `~/.config/tanzu/tkg/unmanaged/${CLUSTER_NAME}/config.yaml`

  Tip: Reviewing this file can help in troubleshooting issues during cluster
bootstrapping.

* To create a cluster with the configuration file, use the `-f` flag to
specify this configuration file:

  ```sh
  tanzu unmanaged-cluster create my-cluster -f my-config.yaml
    ```

Along with a configuration file, `unmanaged-cluster` respects settings from
other settings such as flags. The order in which settings are resolved is:

1. Defaults (lowest precedence)
1. Configuration File
1. Environment Variables
1. Flags (highest precedence)

### Customize cluster provider

Use the `ProviderConfiguration` field in the configuration file
to give provider specific and granular customizations.
Note that some other provider specific configs may be ignored
when `ProviderConfiguration` is used.

* Kind provider: Use the `rawKindConfig` field
  to enter an entire [`kind` configuration file](https://kind.sigs.k8s.io/docs/user/configuration/)
  _or_ a partial config snippet to be used when bootstrapping.
  During bootstrapping, the default kind bootstrapping options are merged with any user provided `rawKindConfig`
  but the values given via the Tanzu CLI and env variables take the highest precedence.
  Any missing values will get the default.
  Merging is done on best effort basis and honors Tanzu CLI flag values over all others.
  View the kind config file that is generated here:  
   `~/.config/tanzu/tkg/unmanaged/${CLUSTER_NAME}/kindconfig.yaml`

  For example, the following partial kind configuration file deploys a control plane with port mappings and 2 worker nodes,
  all using the default VMware hosted kind node images.

  ```yaml
  ClusterName: my-kind-cluster
  KubeconfigPath: ""
  ExistingClusterKubeconfig: ""
  NodeImage: ""
  Provider: kind
  ProviderConfiguration:
    rawKindConfig: |
      nodes:
      - role: control-plane
        extraPortMappings:
        - containerPort: 888
          hostPort: 888
          listenAddress: "127.0.0.1"
          protocol: TCP
        - role: worker
        - role: worker
  Cni: calico
  CniConfiguration: {}
  PodCidr: 10.244.0.0/16
  ServiceCidr: 10.96.0.0/16
  TkrLocation: ""
  AdditionalPackageRepos: []
  PortsToForward: []
  SkipPreflight: false
  ControlPlaneNodeCount: "1"
  WorkerNodeCount: "0"
  InstallPackages: []
  ```

* Minikube provider:
  * `driver` - Optional: Sets the driver to run Kubernetes in. [Selecting a driver depends on your operating system.](https://minikube.sigs.k8s.io/docs/drivers/)
  * `containerRuntime` - Optional: Sets the container runtime to use. Valid options: `docker`, `cri-o`, `containerd`, `auto`.
  * `rawMinikubeArgs` - Optional: The raw flags and arguments to pass to the minikube binary.  
    _Warning:_ use with caution. Flags and arguments provided through this method are not checked or validated by the unmanaged-cluster plugin.

  Example using config options:

  ```yaml
  ClusterName: my-minikube-cluster
  KubeconfigPath: ""
  ExistingClusterKubeconfig: ""
  NodeImage: ""
  Provider: minikube
  ProviderConfiguration:
    driver: vmware
    containerRuntime: docker
    rawMinikubeArgs: --disk-size=30000mb
  Cni: calico
  CniConfiguration: {}
  PodCidr: 10.244.0.0/16
  ServiceCidr: 10.96.0.0/16
  TkrLocation: ""
  AdditionalPackageRepos: []
  PortsToForward: []
  SkipPreflight: false
  ControlPlaneNodeCount: "1"
  WorkerNodeCount: "0"
  InstallPackages: []
  ```

  The above configuration file can be used with another port mapping via the `-p` CLI flag.
  This will result in the _same_ deployment, but the port mapping configuration is merged
  resulting in the first node getting the additional port mapping.

  ```sh
  tanzu unmanaged-cluster create -f my-config-file -p 123:123
  ```

  For the _most_ granular configuration of kind, enter a _complete_ kind configuration file under `rawKindConfig`
  with no additional CLI flags or environment variables given.

## Install to existing cluster

If you wish to install the Tanzu components, such as `kapp-controller` and the
package repositories into an **existing** unmanaged cluster, you can do so with the
`--existing-cluster-kubeconfig`/`e` flags or `existingClusterKubeconfig`
configuration field. The following example demonstrates installing into an
existing [minikube](https://minikube.sigs.k8s.io) cluster.

1. Create a `minikube` cluster.

    ```sh
    $ minikube start

    * minikube v1.24.0 on Arch rolling
    * Automatically selected the docker driver. Other choices: kvm2, ssh
    * Starting control plane node minikube in cluster minikube
    * Pulling base image ...

    * Preparing Kubernetes v1.22.3 on Docker 20.10.8 ...
      - Generating certificates and keys ...
      - Booting up control plane ...
      - Configuring RBAC rules ...
    * Verifying Kubernetes components...
      - Using image gcr.io/k8s-minikube/storage-provisioner:v5
    * Enabled addons: storage-provisioner, default-storageclass
    * Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
    ```

1. Install the unmanaged cluster components

    ```sh
    tanzu unmanaged-cluster create -e ~/.kube/config --cni=none
    ```

    * `~/.kube/config` is the location of the kubeconfig used to access the
      `minikube` cluster.
    * `--cni=none` is set since `minikube` already sets up a network for pods.

1. Now you can use the Tanzu CLI to interact with the cluster.

    ```sh
    tanzu package list -A
    ```

## Disable CNI installation

To create a cluster **without** a CNI installed, run:

```sh
tanzu unmanaged-cluster create --cni=none
```

This will skip CNI installation and prompt the following during the CNI
installation step.

```txt
🌐 Installing CNI
   No CNI installed: CNI was set to none.
```

The cluster creation will complete successfully. After that, you are free to
install a CNI into the cluster.

## Customize the distribution

Unmanaged clusters gather details on how to create a cluster from a Tanzu Kubernetes Release (TKr) file. For each release of unmanaged clusters, a [default TKr is
set](https://github.com/vmware-tanzu/community-edition/blob/d0a8622e164c1e345686470b7bcce0c6be9c58f5/cli/cmd/plugin/unmanaged-cluster/tkr/tkr.go#L14-L16).

When creating clusters, you can point to a different TKr using the `--tkr` flag.
TKrs for unmanaged clusters are available here: `projects.registry.vmware.com/tce/tkr`.  
 Use [imgpkg](https://carvel.dev/imgpkg), to query available TKrs:

```sh
$ imgpkg tag list -i projects.registry.vmware.com/tce/tkr
Tags

Name
sha256-2fd337282cf17357c6329f441dc970ec900145faef9e2ec6122f98fa75d529c3.imgpkg
sha256-33f63314fb72ead645715f6ac85128c0fe0fd380d14f0a79eddba3dd361b73dd.imgpkg
sha256-ac6566268e0f113a4b91bab870a34353685e886f97e248633bb2c2fcf6490dc8.imgpkg
v1.21.5
v1.21.5-1
v1.21.5-2
v1.21.5-3
v1.22.2
```

To create a cluster with an alternative TKr, you can run:

```sh
tanzu unmanaged-cluster create --tkr projects.registry.vmware.com/tce/tkr:v1.22.2
```

The `--tkr` option also supports local files.

```sh
tanzu unmanaged-cluster create --tkr path-to-my-tkr-file.yaml
```

To customize a TKr, you can download using `imgpkg`:

```sh
$ imgpkg pull -i projects.registry.vmware.com/tce/tkr:v1.22.2 -o tkr
Pulling image 'projects.registry.vmware.com/tce/tkr@sha256:7c1a241dc57fe94f02be4dd6d7e4b29f159415417164abc4b5ab6bb10cf4cbaa'
Extracting layer 'sha256:e17e901811682a2c8c91c8865f3344a21fdf8f83f012de167c15d2ab06cc494a' (1/1)

Succeeded
```

You can then edit the above TKr in the `tkr/tkr-bom-v1.22.2.yaml`. After
modifying it, you may also wish to rename the YAML file. Once you have made your
modifications, you can repush it using:

```sh
imgpkg push -f ./tkr/tkr-bom-CUSTOM.yaml -i ${YOUR_REGISTRY}:${YOUR_TAG}
```

Once pushed, you can reference this repo or local file using the `--tkr` flag.

## Exit codes

Unmanaged clusters provide meaningful exit codes.
These are useful when deploying unmanaged clusters in automation or CI/CD.
To see the exit code of a process, execute `echo $?`.

The exit codes are defined as follows:

* 0  - Success.
* 1  - Configuration is invalid.
* 2  - Could not create local cluster directories.
* 3  - Unable to get TKR BOM.
* 4  - Could not render config.
* 5  - TKR BOM not parsable.
* 6  - Could not resolve kapp controller bundle.
* 7  - Unable to create cluster.
* 8  - Unable to use existing cluster (if provided).
* 9  - Could not install kapp controller to cluster.
* 10 - Could not install core package repo to cluster.
* 11 - Could not install additional package repo
* 12 - Could not install CNI package.
* 13 - Failed to merge kubeconfig and set context
* 14 - Could not install designated packages

## Limitations

This section details known limitations of unmanaged clusters.

### Can't Upgrade Kubernetes

By design, unmanaged clusters do not lifecycle-manage Kubernetes. They are not meant to be long-running with real workloads. To change Kubernetes versions, delete the existing cluster and create a new cluster with a different configuration.

### Deploy to Windows

`kind`, the default provider, has several known limitations when deploying to Windows.
For example, deploying a [load balancer has networking considerations.](https://kind.sigs.k8s.io/docs/user/loadbalancer/)
Be sure to familiarize yourself with the [`kind` documentation](https://kind.sigs.k8s.io/) in order to [customize your unmanaged-cluster deployment](#custom-configuration) for your needs.
