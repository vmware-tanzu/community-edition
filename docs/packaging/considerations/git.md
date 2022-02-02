# Git Repository

Your package will consist of multiple configuration files ranging from upstream
manifests to configuration overlays that provide a customized installation
experience. This configuration must be stored in Git. The TCE project does not
mandate how you manage your git repository, as long as your package conforms to
conventions discussed in this guide and produces a valid, installable, package.

## Package Repository

The TCE project recommends a package's source configuration is either:

1. Hosted in the upstream project.
1. Hosted in a `https://github.com/vmware-tanzu/package-for-${PACKAGE_NAME}`
   repository.

The latter option is the most common when:

* The upstream project does not want to host or maintain a carvel package.
* The upstream project is not (primarily) governed by VMware.

For an example of the latter option, see the [package-for-kpack
repository](https://github.com/vmware-tanzu/package-for-kpack).

## Package in the TCE repository

For legacy reasons, some packages are hosted in the TCE repository under
[addons/packages](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages).
We do not recommend this approach going forward. Exceptions can be made if
creating a new repository is not possible. If it is not possible, be sure to
include this detail in your proposal.
