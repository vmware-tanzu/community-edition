# Terminology

This document contains definitions for terms used in our documentation.

* **Packages:** software installed into a Kubernetes cluster. For example, [Project
Contour](https://projectcontour.io).
  * **User-Managed packages**: Deployed into clusters and lifecycle managed independent
  of a cluster. For example [Project
  Contour](https://projectcontour.io).
  * **Core packages**: Deployed into clusters, typically after cluster bootstrap.
  Lifecycle managed as part of cluster. For example,
  [Antrea](https://github.com/vmware-tanzu/antrea).
    * The packaging details in most
  of this document are relevant to core and user-managed packages. However, much of the details
  around discovery, repositories, and CLI interaction are only relevant for
  user-managed packagess.
* **Add-ons:** same as packages (see above)
* **Extensions:** same as packages (see above)
