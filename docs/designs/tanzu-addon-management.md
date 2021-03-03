# Tanzu Add-on Management

This document covers the management of add-ons of a server-side and client-side
perspective. This is a working design doc that will evolve over time as our
add-on management is implemented and enhanced.

## Server Side

### Overview and APIs

## Client Side

This section describes the client side management of extensions. This
specifically focuses on our usage of `tanzu` CLI to discover, configure, deploy,
and manage add-ons.

### Overview and APIs

### Package Discovery

The `tanzu` CLI is able to discover packages known to the cluster. It discovers
these packages by viewing all available [Package
CRs](https://carvel.dev/kapp-controller/docs/latest/packaging/#package) to the
cluster. These Packages can be sourced from 1 or many
[PackageRepository](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-repositories)
CRs. This model is inspired by popular Linux package managers where inclusion of
a repository (e.g. `/etc/apt/sources.list`)  will make new packages available to
the manager (e.g. `apt`). With this, a command such as the following is
possible.

```
tanzu package list

NAME               VERSION    REPO
knative-serving    1.12       tce-main 
contour            2.32       tce-main
nvidia-driver      1.11       nvidia-main
```

In the above, the `tanzu` CLI is aggregating and listing metadata from
already-existent objects. Namely the following from each `Package` instance:

* `NAME`: `spec.publicName`
* `VERSION`: `spec.version`
* `REPO`: TODO(see
[https://github.com/vmware-tanzu/carvel-kapp-controller/issues/124](https://github.com/vmware-tanzu/carvel-kapp-controller/issues/124))

This is visually represented as follows.

<img src="../images/tanzu-package-list.png">

### Package Configuration

### Package Deployment

This is visually represented as follows.

<img src="../images/tanzu-package-install.png">

### Package Management

### Package Repository Discovery

### Package Repository Creation

### Package Repository Deletion
