# Tanzu Community Edition Package Versioning

This document defines TCE's approach to versioning packages and publishing package repositories. Considerations and questions that this document will address include the following:

* What happens when an application/functionality releases a new version?
* What is the prescribed process that TCE will use to make available that new version?
* How many package versions will TCE support?
  If N=$current_version, will we support n-2?
  Does N represent major? minor? patch?
* Can multiple versions of a package be available simultaneously?
* How will TCE respond to CVE's in existing packages?
* What channels should TCE provide for its repository (e.g. stable, beta. alpha)
* How we are going to think about "true" community-owned repositories
  Users will not always want their software to become part of TCE.
* We need to provide guidance on how folks can bring their own packages + bring their own package repositories.
* How will TCE handle `core` packages? These are packages required by Kubernetes itself.
* What can be automated?
* How will packages be tested?



## Directory Structure

Packages and package repositories are contained in the `packages` directory. Within this directory there are directories for the 2 types of packages (`core` and `user-managed`), package repositories (`repositories`) and miscellaneous support files (`misc`).

```txt
./packages
├── core
├── misc
├── repositories
├── user-managed
```

### core

The `core` directory contains packages that are required for Kubernetes itself. These packages are not intended for end user use, modification or installation.

### misc

The `misc` directory contains miscellaneous files used to support the development and testing of packages and package repositories.

### repostiories

The `repositories` directory contains the files used to generate the package repository manifests.

### user-managed

The `user-managed` directory contains packages that are intended for end user consumption. They represent additional functionality that a user can optionally add to extend and enhance their Kubernetes workload clusters.

### `core` and `user-managed` Directory Structure

The `core` and `user-managed` directories will house the individual packages that require versioning. The proposed structure of these directories to support versioning is to place version specific files into a sub-directory named after the version.

For example, consider the cert-manager package. If TCE was to support the previous 3 versions, there would be 