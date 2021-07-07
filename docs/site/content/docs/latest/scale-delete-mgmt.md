# Scale or Delete Management Clusters
##  Delete Management Clusters from Your Tanzu CLI Configuration

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



