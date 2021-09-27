### Managed Clusters

Managed clusters is a deployment model that features `1` management cluster and `N` workload cluster(s). The management cluster provides management and operations for Tanzu. It runs [Cluster-API](https://cluster-api.sigs.k8s.io/) which is used to manage workload clusters and multi-cluster services. The workload cluster(s) are where developer's workloads run.

When you create a management cluster, a bootstrap cluster is created on your local machine. This is a [Kind](https://kind.sigs.k8s.io/) based cluster, which runs via Docker. The bootstrap cluster creates a management cluster on your specified provider. The information for how to manage clusters in the target environment is then pivoted into the management cluster. At this point, the local bootstrap cluster is deleted. The management cluster can now create workload clusters.
