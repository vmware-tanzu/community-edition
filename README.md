<!-- markdownlint-disable MD033 -->
<img src="docs/images/logos/tce-logo-only.png" width="150" align="left">

# Tanzu Community Edition

Tanzu Community Edition is a fully-featured, easy to manage, Kubernetes platform
for learners and users. It is a freely available, community supported, and open
source distribution of VMware Tanzu. It can be installed and deployed in minutes to your
local workstation or favorite infrastructure provider. Along with cluster
management, powered by [Cluster API](https://github.com/kubernetes-sigs/cluster-api),
Tanzu Community Edition enables higher-level functionality via its robust
[package management](https://tanzucommunityedition.io/docs/latest/package-management)
built on top of [Carvel's kapp-controller](https://carvel.dev/kapp-controller/),
and opinionated, yet extensible, [Carvel packages](#packages).

![overview](docs/images/overview.gif)

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/4906/badge)](https://bestpractices.coreinfrastructure.org/projects/4906)
[![Build - Create Dev/Staging](https://github.com/vmware-tanzu/community-edition/actions/workflows/build-staging.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/build-staging.yaml)
[![Check - Main (All tests)](https://github.com/vmware-tanzu/community-edition/actions/workflows/check-main.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/check-main.yaml)
[![Check - imagelint](https://github.com/vmware-tanzu/community-edition/actions/workflows/check-imagelint.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/check-imagelint.yaml)
[![E2E Test - vSphere Management and Workload Cluster](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-vsphere-management-and-workload-cluster.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-vsphere-management-and-workload-cluster.yaml)
[![E2E Test - vSphere Standalone Cluster](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-vsphere-standalone-cluster.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-vsphere-standalone-cluster.yaml)
[![E2E Test - Azure Management and Workload Cluster](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-azure-management-and-workload-cluster.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-azure-management-and-workload-cluster.yaml)
[![E2E Test - Azure Standalone Cluster](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-azure-standalone-cluster.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-azure-standalone-cluster.yaml)
[![E2E Test - AWS Management and Workload Cluster](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-aws-management-and-workload-cluster.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-aws-management-and-workload-cluster.yaml)
[![E2E Test - AWS Standalone Cluster](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-aws-standalone-cluster.yaml/badge.svg)](https://github.com/vmware-tanzu/community-edition/actions/workflows/e2e-aws-standalone-cluster.yaml)

## Getting Started

* [Getting Started Guide](https://tanzucommunityedition.io/docs/latest/getting-started)
* [Documentation](https://tanzucommunityedition.io/docs/latest/getting-started)

## Installation

We recommend installing Tanzu Community Edition using a package manager. If that
is not possible, manual steps are detailed last.

### Mac/Linux via homebrew

```sh
brew install vmware-tanzu/tanzu/tanzu-community-edition
```

After install, homebrew will prompt you with a configure script, run it.

```txt
******************************************************************************
* To initialize all plugins required by TCE, an additional step is required.
* To complete the installation, please run the following shell script:
*
* ${HOMEBREW_EXEC_DIR}/configure-tce.sh
******************************************************************************
```

### Windows via chocolatey

```sh
choco install tanzu-community-edition
```

### Manual (Mac/Linux/Windows)

1. [Download the release tarball](https://github.com/vmware-tanzu/community-edition/releases) based on your operating system.
1. Unpack the release tarball.
    * Unzip on Windows.
    * `tar zxvf <release tarball>` on Mac/Linux.
1. Enter the directory of the unpacked release.
1. Run the install script.
    * `install.bat` on Windows as Administrator.
    * `install.sh` on Mac/Linux

## Packages

Packages provide higher-level functionality to build atop your Kubernetes-based
platform. Packages included, by default, in Tanzu Community Edition are:

| Name | Description | Documentation |
|------|-------------|---------------|
| Load Balancer Operator | Provides load balancer integrations to clusters | [Load Balancer Operator package docs](./addons/packages/ako-operator) |
| Cert Manager | Provides certificate management provisioning within the cluster | [Cert Manager package docs](./addons/packages/cert-manager) |
| Contour | Provides ingress support to workloads | [Contour package docs](./addons/packages/contour) |
| ExternalDNS | Provides discoverability of services via public DNS | [ExternalDNS package docs](./addons/packages/external-dns) |
| Harbor | Provides cloud native container registry service | [Harbor package docs](./addons/packages/harbor) |
| Fluent-Bit | Log processor and forwarder | [Fluent Bit package docs](./addons/packages/fluent-bit) |
| Gatekeeper | Provides policy enforcement within clusters | [Gatekeeper package docs](./addons/packages/gatekeeper) |
| Grafana | Metrics visualization and analytics | [Grafana package docs](./addons/packages/grafana) |
| Knative Serving | Provides serving functionality to clusters | [knative serving package docs](./addons/packages/knative-serving) |
| Prometheus | Time series database for metrics. Includes AlertManager | [Prometheus package docs](./addons/packages/prometheus) |
| Velero | Provides disaster recovery capabilities | [Velero package docs](./addons/packages/velero) |
| Multus CNI | Provides ability for attaching multiple network interfaces to pods in Kubernetes | [Multus CNI package docs](./addons/packages/multus-cni) |
| Whereabouts | Provides A CNI IPAM plugin that assigns IP addresses cluster-wide | [Whereabouts package docs](./addons/packages/whereabouts) |
| Sriov Network Device Plugin | The SR-IOV Network Device Plugin is Kubernetes device plugin for discovering and advertising SR-IOV virtual functions (VFs) available on a Kubernetes host. | [Sriov Network Device Plugin package docs](./addons/packages/sriov-network-device-plugin) |

## Contributing

If you are ready to jump in and test, add code, or help with documentation,
follow the instructions on our [Contribution Guidelines](./CONTRIBUTING.md) to
get started and at all times, follow our [Code of
Conduct](./CODE_OF_CONDUCT.md).

Before opening an issue or pull request, please search for any existing issues
or existing pull requests. If an issue does not exist, please create one for
your feedback! If one exists, please feel free to comment and add any
additional context you may have!

## Repository Layout

The following describes the key directories that make up this repository.

* `addons/`: the source configuration of our packages and package repository
  available to be installed in TCE clusters
  * `packages/`: software packages installable in TCE clusters
  * `repos/`: bundles of packages that can be installed in TCE clusters
    making all packages within available
* `cli/`: plugins that add TCE-specific functionality to the `tanzu` CLI
  * `cmd/plugin/${PLUGIN_NAME}/`: individual plugin project and go module
* `docs/`: documentation and our hugo-based website
* `hack/`: scripts used for development and build processes
* `test/`: scripts, configuration, and code used for end-to-end testing

## Support

If you have any questions about Tanzu Community Edition, please join [#tanzu-community-edition](https://kubernetes.slack.com/messages/tanzu-community-edition) on [Kubernetes slack](http://slack.k8s.io/).

Please submit [bugs or enhancements requests](https://github.com/vmware-tanzu/community-edition/issues/new/choose) in GitHub.
More information about troubleshooting and our triage process is available [here](https://tanzucommunityedition.io/docs/latest/trouble-faq/).

Information about our roadmap is available [here](https://github.com/vmware-tanzu/community-edition/issues/1293).
