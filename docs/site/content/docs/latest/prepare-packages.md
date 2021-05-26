### View Management Cluster Details With Tanzu CLI and `kubectl`

Tanzu CLI provides commands that facilitate many of the operations that you can perform with your management cluster. However, for certain operations, you still need to use `kubectl`. 

When you deploy a management cluster, the `kubectl` context is not automatically set to context of the management cluster. Tanzu Kubernetes Grid provides two contexts for every management cluster and Tanzu Kubernetes cluster:  

- The `admin` context of a cluster gives you full access to that cluster. If you implemented identity management on the cluster, using the `admin` context allows you to run `kubectl` operations without requiring authentication with your identity provider (IDP). 
- If you implemented identity management on the cluster, using the regular context requires you to authenticate with your IDP before you can run `kubectl` operations on the cluster.

Before you can run `kubectl` operations on a management cluster, you must obtain its `kubeconfig`.
   
1. On the bootstrap machine, run the `tanzu login` command to see the available management clusters and which one is the current login context for the CLI. 

   For more information, see [List Management Clusters and Change Context](../cluster-lifecycle/multiple-management-clusters.md#login).

1. To see the details of the management cluster, run `tanzu management-cluster get`.  

   For more information, see [See Management Cluster Details](../cluster-lifecycle/multiple-management-clusters.md#list-mc).

1. To retrieve a `kubeconfig` for the management cluster, run the `tanzu management-cluster kubeconfig get` command as described in [Retrieve Management Cluster `kubeconfig`](#kubeconfig).
   
1. Set the context of `kubectl` to the management cluster.

   ```
   kubectl config use-context my-mgmnt-cluster-admin@my-mgmnt-cluster
   ```

1. Use `kubectl` commands to examine the resources of the management cluster.

   For example, run `kubectl get nodes`, `kubectl get pods`, or `kubectl get namespaces` to see the nodes, pods, and namespaces running in the management cluster.


## Retrieve Management Cluster `kubeconfig`

The `tanzu management-cluster kubeconfig get` command retrieves `kubeconfig` configuration information for the current management cluster, with options as follows:

- `--export-file FILE`
  - **Without option**: Add the retrieved cluster configuration information to the `kubectl` CLI's current `kubeconfig` file, whether it is the default `~/.kube/config` or set by the `KUBECONFIG` environment variable.
  - **With option**: Write the cluster configuration to a standalone `kubeconfig` file `FILE` that you can share with others.

- `--admin`
  - **Without option**: Generate a _regular `kubeconfig`_ that requires the user to authenticate with an external identity provider, and grants them access to cluster resources based on their assigned roles.
      - The context name for this `kubeconfig` includes a `tanzu-cli-` prefix. For example, `tanzu-cli-id-mgmt-test@id-mgmt-test`.
  - **With option**: Generate an _administrator `kubeconfig`_ containing embedded credentials that lets the user access the cluster without logging in to an identity provider, and grants full access to the cluster's resources.
      - The context name for this `kubeconfig` includes an `-admin` suffix. For example, `id-mgmt-test-admin@id-mgmt-test`.

For example, to generate a standalone `kubeconfig` file to share with someone to grant them full access to your current management cluster:

   ```
   tanzu management-cluster kubeconfig get --admin --export-file MC-ADMIN-KUBECONFIG
   ```

To retrieve a `kubeconfig` for a workload cluster, run `tanzu cluster kubeconfig get` as described in [Retrieve Tanzu Kubernetes Cluster `kubeconfig`](../cluster-lifecycle/connect.md#kubeconfig).


## <a id="what-next"></a> What to Do Next

You can now use Tanzu Kubernetes Grid to start deploying Tanzu Kubernetes clusters. For information, see [Deploying Tanzu Kubernetes Clusters](../tanzu-k8s-clusters/index.md).

If you need to deploy more than one management cluster, on any or all of vSphere, Azure, and Amazon EC2, see [Manage Your Management Clusters](../cluster-lifecycle/multiple-management-clusters.md). This topic also provides information about how to add existing management clusters to your CLI instance, obtain credentials, scale and delete management clusters, and how to opt in or out of the CEIP.