### Management cluster description
The management cluster provides management and operations for your instance. It runs [Cluster-API](https://cluster-api.sigs.k8s.io/) which is used to create workload clusters, as well as creating shared services for all the clusters within the instance.  The management cluster is not intended to be used for application workloads. A management cluster is deployed using the Tanzu Community Edition Installer.

When you create a management cluster, a bootstrap cluster is created on your local machine. This is a [Kind](https://kind.sigs.k8s.io/)  based cluster -  a cluster in a container.  This bootstrap cluster then creates a cluster on your specified provider. The Cluster APIs then pivots this cluster into a management cluster.
At this point, the local bootstrap cluster is deleted.  The management cluster can now instantiate more workload clusters.
