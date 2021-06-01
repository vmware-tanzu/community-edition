### Prepare to deploy packages 

When you deploy a management cluster, the `kubectl` context is not automatically set to the context of the management cluster, and this context needs to be set. Before you can run `kubectl` operations on a management cluster, you must obtain its `kubeconfig`. <!--Tanzu Kubernetes Grid provides two contexts for every management cluster and Tanzu Kubernetes cluster:  -->

<!--- The `admin` context of a cluster gives you full access to that cluster. If you implemented identity management on the cluster, using the `admin` context allows you to run `kubectl` operations without requiring authentication with your identity provider (IDP). 
- If you implemented identity management on the cluster, using the regular context requires you to authenticate with your IDP before you can run `kubectl` operations on the cluster.-->

## Procedure

   
1.On the bootstrap machine, run the `tanzu login` command to see the available management clusters and which one is the current login context for the CLI. 
1. To retrieve a `kubeconfig` for the management cluster, run the `tanzu management-cluster kubeconfig get` command 
   
1. Set the context of `kubectl` to the management cluster.

   ```sh
   kubectl config use-context my-mgmnt-cluster-admin@my-mgmnt-cluster
   ```
1. Use `kubectl` commands to examine the resources of the management cluster. For example, run `kubectl get nodes`, `kubectl get pods`, or `kubectl get namespaces` to see the nodes, pods, and namespaces running in the management cluster.
1. (Optional) The `tanzu management-cluster kubeconfig get` command retrieves `kubeconfig` configuration information for the current management cluster, with options as follows:

- `--export-file FILE`
  - **Without option**: Add the retrieved cluster configuration information to the `kubectl` CLI's current `kubeconfig` file, whether it is the default `~/.kube/config` or set by the `KUBECONFIG` environment variable.
  - **With option**: Write the cluster configuration to a standalone `kubeconfig` file `FILE` that you can share with others.

- `--admin`
  - **Without option**: Generate a _regular `kubeconfig`_ that requires the user to authenticate with an external identity provider, and grants them access to cluster resources based on their assigned roles.
      - The context name for this `kubeconfig` includes a `tanzu-cli-` prefix. For example, `tanzu-cli-id-mgmt-test@id-mgmt-test`.
  - **With option**: Generate an _administrator `kubeconfig`_ containing embedded credentials that lets the user access the cluster without logging in to an identity provider, and grants full access to the cluster's resources.
      - The context name for this `kubeconfig` includes an `-admin` suffix. For example, `id-mgmt-test-admin@id-mgmt-test`.

For example, to generate a standalone `kubeconfig` file to share with someone to grant them full access to your current management cluster:

   ```sh
   tanzu management-cluster kubeconfig get --admin --export-file MC-ADMIN-KUBECONFIG
   ```

<!--To retrieve a `kubeconfig` for a workload cluster, run `tanzu cluster kubeconfig get` as described in [Retrieve Tanzu Kubernetes Cluster `kubeconfig`](../cluster-lifecycle/connect.md#kubeconfig).-->
<!-- add this after the cluster content>

**IMPORTANT**: By default, unless you set the `KUBECONFIG` environment variable to save the `kubeconfig` for a cluster to a specific file, all clusters that you deploy from the Tanzu CLI are added to a shared `.kube-tkg/config` file. If you delete the shared `.kube-tkg/config` file, all management clusters become orphaned and thus unusable.
