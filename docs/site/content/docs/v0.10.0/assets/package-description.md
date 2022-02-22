Packages extend the functionality of Tanzu Community Edition. You can discover and deploy packages through the Tanzu CLI. A Tanzu package is an aggregation of Kubernetes configurations, and its associated software container image, into a versioned and distributable bundle, that can be deployed as an OCI container image. Packages are installed into a Tanzu cluster.

- User-Managed packages: Deployed into clusters and the lifecycle of the package is managed independently of the cluster. For example [Project  Contour](https://projectcontour.io).
- Core packages: Deployed into clusters, typically after cluster is bootstrapped. The lifecycle is managed as part of a cluster. For example, [Antrea](https://github.com/vmware-tanzu/antrea).
