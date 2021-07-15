# Delete Management Clusters
This topic describes how to delete management clusters from the Tanzu CLI and also how to fully delete management clusters.
## Delete Management Clusters

To delete a management cluster, run the `tanzu management-cluster delete` command.

When you run `tanzu management-cluster delete`, Tanzu Community Edition creates a temporary `kind` cleanup cluster on your bootstrap machine to manage the deletion process. The `kind` cluster is removed when the deletion process completes.

1. To see all your management clusters, run `tanzu login`.
1. If there are management clusters that you no longer require, run `tanzu management-cluster delete`. To skip the `yes/no` verification step, specify the `--yes` option.   You must be logged in to the management cluster that you want to delete.
    ```sh
    tanzu management-cluster delete my-aws-mgmt-cluster
    ```
2. If there are workload clusters running in the management cluster, the delete operation is not performed. In this case, you can delete the management cluster in two ways:
   - Run `tanzu cluster delete` to delete all of the running clusters and then run `tanzu management-cluster delete` again.
   - Run `tanzu management-cluster delete` with the `--force` option. For example,
   ```sh
   tanzu management-cluster delete my-aws-mgmt-cluster --force
   ```
**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu Community Edition operations are running.

##  Delete orphaned management clusters configurations from your local machine
Under the following conditions, you might need to remove a management cluster from the Tanzu CLI:
- You added a management cluster that someone else created to your instance of the Tanzu CLI, and now want to remove it.
- If a management cluster is deleted directly on your infrastructure provider without running `tanzu management-cluster delete`, then the management cluster continues to appear in the list of management clusters that the CLI tracks when you run `tanzu login`.

1. Run `tanzu config server list`, to see the list of management clusters that the Tanzu CLI tracks. You should see all of the management clusters that you have deployed or added to the Tanzu CLI, the location of their kubeconfig files, and their contexts.
1. Run the `tanzu config server delete` command to remove a management cluster.
   ```sh
   tanzu config server delete <MGMT-CLUSTER>
   ```
   This removes the cluster details from the `~/.tanzu/config.yaml` and `~/.kube-tkg/config.yaml` files. It does not delete the management cluster itself, if it still exists.
