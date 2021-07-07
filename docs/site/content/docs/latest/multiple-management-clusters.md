#  Delete Management Clusters from Your Tanzu CLI Configuration

It is possible that you might add a management cluster that someone else created to your instance of the Tanzu CLI, that at some point you no longer require. Similarly, if you deployed a management cluster and that management cluster has been deleted from your infrastructure provider by means other than by running `tanzu management-cluster delete`, that management cluster will continue to appear in the list of management clusters that the CLI tracks when you run `tanzu login`. In these cases, you can remove the management cluster from the list of management clusters that the Tanzu CLI tracks.

1. Run `tanzu config server list`, to see the list of management clusters that the Tanzu CLI tracks.

   ```
   tanzu config server list
   ```

   You should see all of the management clusters that you have either deployed yourself or added to the Tanzu CLI, the location of their kubeconfig files, and their contexts.

1. Run the `tanzu config server delete` command to remove a management cluster.

   ```
   tanzu config server delete my-vsphere-mc
   ```

Running the `tanzu config server delete` command removes the cluster details from the `~/.tanzu/config.yaml` and `~/.kube-tkg/config.yaml` files. It does not delete the management cluster itself, if it still exists. To delete a  management cluster rather than just remove it from the Tanzu CLI configuration, see [Delete Management Clusters](#delete).

## Delete Management Clusters

To delete a management cluster, run the `tanzu management-cluster delete` command.

When you run `tanzu management-cluster delete`, Tanzu Commuity Edition creates a temporary `kind` cleanup cluster on your bootstrap machine to manage the deletion process. The `kind` cluster is removed when the deletion process completes.

1. To see all your management clusters, run `tanzu login`.

1. If there are management clusters that you no longer require, run `tanzu management-cluster delete`.

    You must be logged in to the management cluster that you want to delete.

    ```
    tanzu management-cluster delete my-aws-mgmt-cluster
    ```

    To skip the `yes/no` verification step when you run `tanzu management-cluster delete`, specify the `--yes` option.

    ```
    tanzu management-cluster delete my-aws-mgmt-cluster --yes
    ```

1. If there are workload clusters running in the management cluster, the delete operation is not performed.

   In this case, you can delete the management cluster in two ways:

   - Run `tanzu cluster delete` to delete all of the running clusters and then run `tanzu management-cluster delete` again.
   - Run `tanzu management-cluster delete` with the `--force` option. For example,

   ```sh
   tanzu management-cluster delete my-aws-mgmt-cluster --force
   ```

**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu Commuity Edition operations are running.

## Scale Management Clusters

After you deploy a management cluster, you can scale it up or down by increasing or reducing the number of node VMs that it contains. To scale a management cluster, use the `tanzu cluster scale` command with one or both of the following options:

* `--controlplane-machine-count` changes the number of management cluster control plane nodes.
* `--worker-machine-count` changes the number of management cluster worker nodes.

Because management clusters run in the `tkg-system` namespace rather than the `default` namespace, you must also specify the `--namespace` option when you scale a management cluster.

1. Run `tanzu login` before you run `tanzu cluster scale` to make sure that the management cluster to scale is the current context of the Tanzu CLI.
1. To scale a production management cluster that you originally deployed with 3 control plane nodes and 5 worker nodes to 5 and 10 nodes respectively, run the following command:

```sh
tanzu cluster scale MANAGEMENT-CLUSTER-NAME --controlplane-machine-count 5 --worker-machine-count 10 --namespace tkg-system
```

If you initially deployed a development management cluster with one control plane node and you scale it up to 3 control plane nodes, Tanzu Commuity Edition automatically enables stacked HA on the control plane.

**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu Community Edition operations are running.

## <a id="namespaces"></a> Create Namespaces in the Management Cluster

To help you to organize and manage your development projects, you can optionally divide the management cluster into Kubernetes namespaces. You can then use Tanzu CLI to deploy workload clusters to specific namespaces in your management cluster. For example, you might want to create different types of clusters in dedicated namespaces. If you do not create additional namespaces, Tanzu Commuity Edition creates all workload clusters in the `default` namespace. For information about Kubernetes namespaces, see the [Kubernetes documentation](https://kubernetes.io/docs/tasks/administer-cluster/namespaces-walkthrough/).

1. Make sure that `kubectl` is connected to the correct management cluster context by displaying the current context.

   ```
   kubectl config current-context
   ```

1. List the namespaces that are currently present in the management cluster.

   ```
   kubectl get namespaces
   ```

   You will see that the management cluster already includes several namespaces for the different services that it provides:

    ```
    capi-kubeadm-bootstrap-system       Active   4m7s
    capi-kubeadm-control-plane-system   Active   4m5s
    capi-system                         Active   4m11s
    capi-webhook-system                 Active   4m13s
    capv-system                         Active   3m59s
    cert-manager                        Active   6m56s
    default                             Active   7m11s
    kube-node-lease                     Active   7m12s
    kube-public                         Active   7m12s
    kube-system                         Active   7m12s
    tkg-system                          Active   3m57s
    ```  

1. Use `kubectl create -f` to create new namespaces, for example for development and production. These examples use the `production` and `development` namespaces from the Kubernetes documentation.

```sh
kubectl create -f https://k8s.io/examples/admin/namespace-dev.json
kubectl create -f https://k8s.io/examples/admin/namespace-prod.json
```

2. Run `kubectl get namespaces --show-labels` to see the new namespaces.

    ```
    development                         Active   22m   name=development
    production                          Active   22m   name=production
    ```



