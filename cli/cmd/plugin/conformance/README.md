# Conformance Plugin for Tanzu Community Edition

The `tanzu conformance` plugin is a wrapper around [Sonobuoy](https://sonobuoy.io) that is bundled with Tanzu Community Edition.

Sonobuoy's basic functions and commands will still work with the `conformance` plugin - the biggest difference is substituting `sonobuoy` for `tanzu conformance`.

To do a basic, non-destructive conformance test on the Tanzu cluster currently active in your Kubeconfig, run the following:

```sh
tanzu conformance run --wait
```

Once this command is finished, you can then inspect the results by downloading and extracting them.

```sh
results=$(tanzu conformance retrieve)
tanzu conformance results $results
```

You may entirely remove Sonobuoy from the cluster once you are done. This will remove the `sonobuoy` namespace and everything within it.

```sh
tanzu conformance delete --wait
```

You can find out more about Sonobuoy's capabilities in the [Sonobuoy documentation](https://sonobuoy.io/docs/), and more about conformance testing from the [Cloud Native Computing Foundation's conformance repo](https://github.com/cncf/k8s-conformance#certified-kubernetes) and [the Kubernetes conformance tests](https://github.com/kubernetes/kubernetes/tree/master/test/conformance).

## Special Considerations

When running conformance tests against a cluster with a `DEV` plan, or other single-node clusters, conformance tests related to scheduling will likely fail.
This is because for Kubernetes cluster to be considered conformant, there must be at least two (2) nodes where pods may be scheduled to run.
For testing software on a local laptop, these conformance errors can be safely ignored.
