A package repository is a collection of packages.  A package repository defines metadata information that makes it possible to discover, install, manage, and upgrade packages on your clusters. Before a package can be deployed in a cluster, it must be made discoverable via a package repository.

A package repository is a collection of Kubernetes custom resources that are handled by the Tanzu Community Edition kapp-controller. Similar to a Linux package repository, a Tanzu package repository declaratively defines metadata information that makes it possible to discover, install, manage, and upgrade software packages on running clusters.

Tanzu Community Edition provides a package repository called `tce-repo` that provides a collection of packages necessary to start building an application platform on Kubernetes. You can create your own package repository to distribute different software.
