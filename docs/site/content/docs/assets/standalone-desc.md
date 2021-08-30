## Standalone cluster description

A standalone cluster is a faster way to get a functioning cluster with minimal resources. A standalone cluster functions as a workload cluster, it can run application workloads. It does not contain any of the components related to cluster management.  It is deployed using the Tanzu Kubernetes Grid installer interface.

When you create a standalone cluster, a bootstrap cluster is created on your local machine. This is a [Kind](https://kind.sigs.k8s.io/)  based cluster - a cluster in a container.  This bootstrap cluster then creates a cluster on your specified provider, but it does not pivot into a management cluster - it functions as a workload cluster.
