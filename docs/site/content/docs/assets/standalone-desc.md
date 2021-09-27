### Standalone Clusters

A standalone cluster is a faster way to get a functioning workload cluster with less resources than managed clusters. These clusters do not require a long-running management cluster. A standalone cluster is created using a bootstrap cluster on your local machine with [Kind](https://kind.sigs.k8s.io/). After the standalone cluster is created, the bootstrap cluster is destroyed. Any operations against the standalone cluster, e.g. deletion, will re-invoke the bootstrap cluster.
