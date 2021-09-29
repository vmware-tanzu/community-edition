You can discover and deploy packages through the Tanzu CLI. Packages extend the functionality of Tanzu Community Edition. 
Packages are software installed into a Kubernetes cluster. For example, [Project
  Contour](https://projectcontour.io).

- User-Managed packages: Deployed into clusters and the lifecycle of the package is managed independently of the cluster. For example [Project  Contour](https://projectcontour.io).
- Core packages: Deployed into clusters, typically after cluster is bootstrapped. The lifecycle is managed as part of a cluster. For example, [Antrea](https://github.com/vmware-tanzu/antrea).
