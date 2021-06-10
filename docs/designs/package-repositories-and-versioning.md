# Package Repositories, Ownership, and Packages

This document defines TCE's approach to publishing package repositories and
versioning. Considerations and questions that this document will address include
the following:

## Package Repositories

A [package
repository](https://carvel.dev/kapp-controller/docs/latest/packaging/#packagerepository-cr)
is a collection of
[packages](https://carvel.dev/kapp-controller/docs/latest/packaging/#packagerepository-cr).
There are 2 classes of package repositories.

a. `core`: Contains packages managed as part of a Tanzu Kubernetes
Release. These are not installed or managed by users after cluster bootstrap.

b. `user-managed`: Contains packages that can be installed, by users, on top of
a running Tanzu cluster.

The following characteristics are true for packages and package repositories.

a. All **packages** live in the `tce` git repository.

b. The `core` **package repository** is maintained in Tanzu's `core` git repository.

c. The `user-managed` **package repositories** for TCE are
maintained in the `tce` git repository.

TCE will provide 2 **user-managed** package repositories:

* `main`: contains stable packages.
  * tagged `:v${MOST-RECENT-TCE-TAG}`.
  * tagged `:stable`, representing the latest version.
    * tag is re-written every time the repository is pushed
* `beta`: contains beta packages, only meant for testing.
  * tagged `:v${MOST-RECENT-TCE-TAG}`.

## Ownership

This section details the ownership and responsibilities for packages and package
repositories.

### Package Ownership

A package has an owner or set of owners. Owner(s) are responsible
for:

* Maintaining package source
* Pushing package bundles to a registry
* Creating PRs to update the TCE package repository

To find package owners, see [CODEOWNERS](/CODEOWNERS).

### Package Repository Ownership

The community edition package repository is owned by the TCE team. The TCE team
is responsible for:

* Merging PRs requesting an update to the community edition package repository
* Pushing the community edition package repository to a registry

## Packages

The source for each package is found in `addons/packages`. The process for
packaging is defined
[here](https://github.com/vmware-tanzu/tce/blob/main/docs/designs/tanzu-packaging-process.md).

Each package should be versioned using [Semantic
Versioning](https://semver.org/). The package versioning should **not** be bound
to the packaged software. It should represent the semantic version of the
**package**, which users and tools can look at to determine whether there has
been a breaking change, minor release, or (no functionality changing) patch.
Package owners are responsible for versioning their package with the assumption
that they will **not** break semantic versioning guarantees. Unlike other
software at VMware, packages do not need to be appended with `-vmwareX`.

For each new instance of a major/minor package, a **new** directory should be
introduced representing that version. TCE recommends maintaining the source of
at least `N-2` packages, although package authors are empowered to retain more or
less. Consider the following example for the Prometheus package.

```txt
$ tree -L 2 prometheus/
prometheus/
├── v1.0.5
│   ├── bundle
│   ├── installedpackage.yaml
│   ├── prometheus-addon-clusterrolebinding.yaml
│   ├── prometheus-addon-serviceaccount.yaml
│   └── README.md
├── v1.1.0
│   ├── bundle
│   ├── installedpackage.yaml
│   ├── prometheus-addon-clusterrolebinding.yaml
│   ├── prometheus-addon-serviceaccount.yaml
│   └── README.md 
└── v2.0.0
    ├── bundle
    ├── installedpackage.yaml
    ├── prometheus-addon-clusterrolebinding.yaml
    ├── prometheus-addon-serviceaccount.yaml
    └── README.md
```

In the above, `v2.0.0` is the newest package. However, `v1.1.X` and `v1.0.X` are still
expected to be available. Creation a `v2.1.0` package would require a new
directory, while creation of a `v2.0.1` package would not (since it's only a
patch change). If `v2.1.0` is created, the package owner may consider deleting
`v1.0.5`.

🛑: We acknowledge the flat directory structure is less-than-ideal. However,
this decision was made for 2 reasons. 1) it ensured compatibility with existing
downstream packaging models. 2) it is a stop gap solution until we move to
individual repositories for each package.

### Package Updates

TCE encourages packages are constantly updated and stay aligned with the latest
stable version(s) of software they package. However it is up to package owners
as to the cadence of updates. Package owners should remember that updates to a
package does not guarantee inclusion in the TCE community repository. This gate
is currently controlled by the TCE team.

As a community project, TCE provides no guarantee of package updates in response
to CVEs or other critical issues.

### Package Update Automation

When an upstream source releases new software, ideally, automation can detect
this change and make a PR for the package. Package owners are encouraged to
setup automation around this flow. However, at this time it is not implemented,
when it is, we'll try to use the first case as a reference implementation.

### Package Promotion in Repositories

For a package to be promoted in the `main:stable` community package repository,
the following must be true.

1. A pull request incrementing the available package versioning the package
   repository must be made by the package owner(s)
1. End-to-end tests prescribed by the TCE team must pass
1. The TCE team must merge the change
