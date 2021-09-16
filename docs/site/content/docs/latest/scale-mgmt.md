# Scale Management Clusters

After you deploy a management cluster, you can scale it up or down by increasing or reducing the number of node VMs that it contains. To scale a management cluster, use the `tanzu cluster scale` command with one or both of the following options:

* `--controlplane-machine-count` changes the number of management cluster control plane nodes.
* `--worker-machine-count` changes the number of management cluster worker nodes.

Because management clusters run in the `tkg-system` namespace rather than the `default` namespace, you must also specify the `--namespace` option when you scale a management cluster.

1. Run `tanzu login` before you run `tanzu cluster scale` to check that the management cluster you wish to scale is the current context of the Tanzu CLI.
1. To scale a production management cluster that you originally deployed with 3 control plane nodes and 5 worker nodes to 5 and 10 nodes respectively, run the following command:

   ```sh
   tanzu cluster scale <MANAGEMENT-CLUSTER-NAME> --controlplane-machine-count 5 --worker-machine-count 10 --namespace tkg-system
   ```

If you initially deployed a development management cluster with one control plane node and you scale it up to 3 control plane nodes, Tanzu Community Edition automatically enables stacked HA on the control plane.

**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu Community Edition operations are running.
