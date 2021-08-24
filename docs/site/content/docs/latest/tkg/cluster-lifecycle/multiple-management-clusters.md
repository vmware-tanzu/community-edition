# Manage Your Management Clusters

This topic explains how to manage multiple management clusters from the same bootstrap machine,
including management clusters deployed by Tanzu Kubernetes Grid to vSphere, Azure, or Amazon EC2 and vSphere with Tanzu Supervisor Clusters designated as Tanzu Kubernetes Grid management clusters.

## <a id="login"></a> List Management Clusters and Change Context

To list available management clusters and see which one you are currently logged in to, run `tanzu login` on your bootstrap machine:

   - To change your current login context, use your up- and down-arrow keys to highlight the new management cluster and then press **Enter**.
   - To retain your current context, press **Enter** without changing the highlighting.

For example, if you have two management clusters, `my-vsphere-mgmt-cluster` and `my-aws-mgmt-cluster`, you are currently logged in to `my-vsphere-mgmt-cluster`:

   ```
   $ tanzu login
   ? Select a server  [Use arrows to move, type to filter]
   > my-vsphere-mgmt-cluster  ()
     my-aws-mgmt-cluster      ()
     + new server
   ```

## <a id="list-mc"></a> See Management Cluster Details

To see the details of a management cluster:

1. Run `tanzu login` to log in to the management cluster, as described in [List Management Clusters and Change Context](#login).

1. Run `tanzu management-cluster get`.  For example:

   ```
   $ tanzu management-cluster get
   NAME         NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES
   mc-test-cli  tkg-system  running  1/1           1/1      v1.20.1+vmware.2  management

   Details:

   NAME                                                            READY  SEVERITY  REASON  SINCE  MESSAGE
   /mc-test-cli                                                    True                     29m
   ├─ClusterInfrastructure - AzureCluster/mc-test-cli              True                     30m
   ├─ControlPlane - KubeadmControlPlane/mc-test-cli-control-plane  True                     29m
   │ └─Machine/mc-test-cli-control-plane-htlc4                     True                     30m
   └─Workers
     └─MachineDeployment/mc-test-cli-md-0
       └─Machine/mc-test-cli-md-0-699df4dc76-9kgmw                 True                     30m

   Providers:

     NAMESPACE                          NAME                   TYPE                    PROVIDERNAME  VERSION  WATCHNAMESPACE
     capi-kubeadm-bootstrap-system      bootstrap-kubeadm      BootstrapProvider       kubeadm       v0.3.14
     capi-kubeadm-control-plane-system  control-plane-kubeadm  ControlPlaneProvider    kubeadm       v0.3.14
     capi-system                        cluster-api            CoreProvider            cluster-api   v0.3.14
     capz-system                        infrastructure-azure   InfrastructureProvider  azure         v0.4.8
   ```

To see more options, run `tanzu management-cluster get --help`.

## <a id="kubectl"></a> Management Clusters, `kubectl`, and `kubeconfig`

Tanzu Kubernetes Grid does not automatically change the `kubectl` context when you run `tanzu login` to change the `tanzu` CLI context. Also, Tanzu Kubernetes Grid does not set the `kubectl` context to a workload cluster when you create it. To change the `kubectl` context, use the `kubectl config use-context` command.

By default, Tanzu Kubernetes Grid saves cluster context information in the following files on your bootstrap machine:

- **Management cluster contexts**: `~/.kube-tkg/config`
- **Workload cluster contexts**: `~/.kube/config`

## <a id="config"></a> Management Clusters and Their Configuration Files

When you run `tanzu management-cluster create` for the first time, it creates the `~/.tanzu/tkg` subfolder that contains the Tanzu Kubernetes Grid configuration files. To deploy your first management cluster, you must specify the `--ui` or `--file` option with `tanzu management-cluster create`:

* `tanzu management-cluster create --ui` creates a management cluster with the installer interface and saves the settings from your installer input into a cluster configuration file `~/.tanzu/tkg/clusterconfigs/UNIQUE-ID.yaml`, where `UNIQUE-ID` is a generated filename.

* `tanzu management-cluster create --file` creates a management cluster using an existing cluster configuration file. The `--file` option applies to cluster configuration files only and does not change where the `tanzu` CLI references other files under `~/.tanzu/tkg`.

* `tanzu management-cluster create` with neither the `--ui` nor `--file` option creates a management cluster using the default cluster configuration file `~/.tanzu/tkg/cluster-config.yaml`.

The recommended practice is to use a dedicated configuration file for every management cluster that you deploy.

For more information about configuration files in Tanzu Kubernetes Grid, see [What Happens When You Create a Management Cluster](../mgmt-clusters/deploy-management-clusters.md#what-happens).

## <a id="add-mc"></a> Add Existing Management Clusters to Your Tanzu CLI

To log in to a management cluster that someone else created, run the `tanzu login` command, select **+ new server**, and then select **Server endpoint** or **Local kubeconfig** as your login type.

For example, to log in to an existing management cluster using a local kubeconfig:

1. Run `tanzu login`, use your down-arrow key to highlight **+ new server**, and press **Enter**.

   ```
   tanzu login
   ? Select a server + new server
   ```

1. When prompted, select **Local kubeconfig** as your login type and enter the path to your local kubeconfig file, context, and the name of your server. For example:

   ```
   tanzu login
   ? Select a server + new server
   ? Select login type Local kubeconfig
   ? Enter path to kubeconfig (if any) /Users/exampleuser/examples/kubeconfig
   ? Enter kube context to use new-mgmt-cluster-admin@new-mgmt-cluster
   ? Give the server a name new-mgmt-cluster
   ✔  successfully logged in to management cluster using the kubeconfig new-mgmt-cluster
   ```

   Alternatively, you can run `tanzu login` with the `--server`, `--kubeconfig`, and `--context` options.

## <a id="delete-mc-config"></a> Delete Management Clusters from Your Tanzu CLI Configuration

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

## <a id="scale"></a> Scale Management Clusters

After you deploy a management cluster, you can scale it up or down by increasing or reducing the number of node VMs that it contains. To scale a management cluster, use the `tanzu cluster scale` command with one or both of the following options:

* `--controlplane-machine-count` changes the number of management cluster control plane nodes.
* `--worker-machine-count` changes the number of management cluster worker nodes.

Because management clusters run in the `tkg-system` namespace rather than the `default` namespace, you must also specify the `--namespace` option when you scale a management cluster.

1. Run `tanzu login` before you run `tanzu cluster scale` to make sure that the management cluster to scale is the current context of the Tanzu CLI.
1. To scale a production management cluster that you originally deployed with 3 control plane nodes and 5 worker nodes to 5 and 10 nodes respectively, run the following command:

    <pre>tanzu cluster scale MANAGEMENT-CLUSTER-NAME --controlplane-machine-count 5 --worker-machine-count 10 --namespace tkg-system</pre>

If you initially deployed a development management cluster with one control plane node and you scale it up to 3 control plane nodes, Tanzu Kubernetes Grid automatically enables stacked HA on the control plane.

**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu Kubernetes Grid operations are running.

## <a id="creds"></a> Update Management Cluster Credentials (vSphere)

To update the vSphere credentials used by a management cluster, and optionally all of the workload clusters that it manages, see [Update Management and Workload Cluster Credentials](../cluster-lifecycle/secrets.md#mgmt-creds-update).

## <a id="ceip"></a> Opt in or Out of the VMware CEIP

When you deploy a management cluster by using either the installer interface or the CLI, participation in the VMware Customer Experience Improvement Program (CEIP) is enabled by default, unless you specify the option to opt out. If you remain opted in to the program, the management cluster sends information about how you use Tanzu Kubernetes Grid back to VMware at regular intervals, so that we can make improvements in future versions. Management clusters send the following information to VMware:

- The number of Tanzu Kubernetes clusters that you deploy.
- The infrastructure, network, and storage providers that you use.
- The time that it takes for the `tanzu` CLI to perform basic operations, such as `cluster create`, `cluster delete`, `cluster scale`, and `cluster upgrade`.
- The Tanzu Kubernetes Grid extensions that you implement.
- The plans that you use to deploy clusters, as well as the number and configuration of the control plane and worker nodes.
- The versions of Tanzu Kubernetes Grid and Kubernetes that you use.
- The type and size of the workloads that your clusters run, as well as their lifespan.
- Whether or not you integrate Tanzu Kubernetes Grid with Tanzu Kubernetes Grid Service for vSphere, Tanzu Mission Control, or Tanzu Observability by Wavefront.
- The nature of any problems, errors, and failures that you encounter when using Tanzu Kubernetes Grid, so that we can identify which areas of Tanzu Kubernetes Grid need to be made more robust.

If you opted out of the CEIP when you deployed a management cluster and want to opt in, or if you opted in and want to opt out, you can change your CEIP participation setting after deployment.

CEIP runs as a `cronjob` on the management cluster. It does not run on workload clusters.

1. Run `tanzu login`, as described in [List Management Clusters and Change Context](#login) to log in to the management cluster that you want to see or set CEIP status for.

1. Run the `tanzu management-cluster ceip-participation get` command to see the CEIP status of the current management cluster.

   ```
   tanzu management-cluster ceip-participation get
   ```

   The status `Opt-in` means that CEIP participation is enabled on a management cluster. `Opt-out` means that CEIP participation is disabled.

   ```
   MANAGEMENT-CLUSTER-NAME        CEIP-STATUS
   my-aws-mgmt-cluster            Opt-out
   ```

1. To enable CEIP participation on a management cluster on which it is currently disabled, run the `tanzu management-cluster ceip-participation set` command with the value `true`.

   ```
   tanzu management-cluster ceip-participation set true
   ```

1. To verify that the CEIP participation is now active, run `tanzu management-cluster ceip-participation get` again.

   The status should now be `Opt-in`.

   ```
   MANAGEMENT-CLUSTER-NAME        CEIP-STATUS
   my-aws-mgmt-cluster            Opt-in
   ```

   You can also check that the CEIP `cronjob` is running by setting the `kubectl` context to the management cluster and running `kubectl get cronjobs -A`. For example:

   ```
   kubectl config use-context my-aws-mgmt-cluster-admin@my-aws-mgmt-cluster
   ```

   ```
   kubectl get cronjobs -A
   ```

   The output shows that the `tkg-telemetry` job is running:

   ```
   NAMESPACE              NAME            SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
   tkg-system-telemetry   tkg-telemetry   0 */6 * * *   False     0        <none>          18s
   ```

1. To disable CEIP participation on a management cluster on which it is currently enabled, run the `tanzu management-cluster ceip-participation set` command with the value `false`.

   ```
   tanzu management-cluster ceip-participation set false
   ```

1. To verify that the CEIP participation is disabled, run `tanzu management-cluster ceip-participation get` again.

   The status should now be `Opt-out`.

   ```
   MANAGEMENT-CLUSTER-NAME        CEIP-STATUS
   my-aws-mgmt-cluster            Opt-out
   ```

   If you run `kubectl get cronjobs -A` again, the output shows that no job is running:

   ```
   No resources found
   ```

## <a id="namespaces"></a> Create Namespaces in the Management Cluster

To help you to organize and manage your development projects, you can optionally divide the management cluster into Kubernetes namespaces. You can then use Tanzu CLI to deploy Tanzu Kubernetes clusters to specific namespaces in your management cluster. For example, you might want to create different types of clusters in dedicated namespaces. If you do not create additional namespaces, Tanzu Kubernetes Grid creates all Tanzu Kubernetes clusters in the `default` namespace. For information about Kubernetes namespaces, see the [Kubernetes documentation](https://kubernetes.io/docs/tasks/administer-cluster/namespaces-walkthrough/).

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

1. Use `kubectl create -f` to create new namespaces, for example for development and production.

   These examples use the `production` and `development` namespaces from the Kubernetes documentation.

   <pre>kubectl create -f https://k8s.io/examples/admin/namespace-dev.json</pre>

   <pre>kubectl create -f https://k8s.io/examples/admin/namespace-prod.json</pre>

1. Run `kubectl get namespaces --show-labels` to see the new namespaces.

    ```
    development                         Active   22m   name=development
    production                          Active   22m   name=production
    ```

## <a id="delete"></a> Delete Management Clusters

To delete a management cluster, run the `tanzu management-cluster delete` command.

When you run `tanzu management-cluster delete`, Tanzu Kubernetes Grid creates a temporary `kind` cleanup cluster on your bootstrap machine to manage the deletion process. The `kind` cluster is removed when the deletion process completes.

1. To see all your management clusters, run `tanzu login` as described in [List Management Clusters and Change Context](#login).

1. If there are management clusters that you no longer require, run `tanzu management-cluster delete`.

    You must be logged in to the management cluster that you want to delete.

    ```
    tanzu management-cluster delete my-aws-mgmt-cluster
    ```

    To skip the `yes/no` verification step when you run `tanzu management-cluster delete`, specify the `--yes` option.

    ```
    tanzu management-cluster delete my-aws-mgmt-cluster --yes
    ```

1. If there are Tanzu Kubernetes clusters running in the management cluster, the delete operation is not performed.

   In this case, you can delete the management cluster in two ways:

   - Run `tanzu cluster delete` to delete all of the running clusters and then run `tanzu management-cluster delete` again.
   - Run `tanzu management-cluster delete` with the `--force` option.

   ```
   tanzu management-cluster delete my-aws-mgmt-cluster --force
   ```

**IMPORTANT**: Do not change context or edit the `.kube-tkg/config` file while Tanzu Kubernetes Grid operations are running.

## <a id="what-next"></a> What to Do Next

You can use Tanzu Kubernetes Grid to start deploying Tanzu Kubernetes clusters to different Tanzu Kubernetes Grid instances. For information, see [Deploying Tanzu Kubernetes Clusters](../tanzu-k8s-clusters/index.md).

If you have vSphere 7, you can also deploy and manage Tanzu Kubernetes clusters in vSphere with Tanzu. For information, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).
