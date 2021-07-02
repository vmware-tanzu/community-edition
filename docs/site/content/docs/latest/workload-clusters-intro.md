# Deploying Workload Clusters

This section describes how to use the Tanzu CLI to deploy and manage workload clusters.

## Before You Begin

Before you can create workload clusters, you must install the Tanzu CLI and deploy a management cluster.

You can use the Tanzu CLI to deploy workload clusters to the following platforms:

- vSphere
- Amazon EC2
- Docker

## About Workload Clusters

In Tanzu Community Edition, your application workloads run on workload clusters.

Tanzu Community Edition automatically deploys workload clusters to the platform on which you deployed the management cluster. For example, you cannot deploy workload clusters to Amazon EC2 from a management cluster that is running in vSphere. It is not possible to use shared services between the different providers because each provider uses different systems.

Tanzu Community Edition automatically deploys clusters from whichever management cluster you have set as the context for the CLI by using the `tanzu login` command. <!--For information about `tanzu login`, see [Manage Your Management Clusters](../cluster-lifecycle/multiple-management-clusters.md).-->

<!--- For information about how to use the Tanzu CLI to deploy Workload clusters, see [Deploy Workload Clusters](deploy.md) and its subtopics.-->
<!--- After you have deployed Workload clusters, the Tanzu CLI provides commands and options to perform the following cluster lifecycle management operations. See [Managing Cluster Lifecycles](../cluster-lifecycle/index.md).-->
<!--note to self - check back if this is need - upgrading-->
<!--For information about how to upgrade existing clusters to a new version of Kubernetes, see [Upgrade Workload Clusters](../upgrade-tkg/workload-clusters.md).-->

## <a id="kubectl"></a> Workload Clusters, `kubectl`, and `kubeconfig`

When you create a management cluster, the Tanzu CLI and `kubectl` contexts are automatically set to that management cluster. However, when you create a workload cluster, the `kubectl` context is not automatically reset to the new workload cluster. You must set the `kubectl` context to a workload cluster manually by using the `kubectl config use-context` command. <!--need to describe how -->

By default, unless you specify the `KUBECONFIG` option to save the `kubeconfig` for a cluster to a specific file, all workload clusters that you deploy are added to a shared `.kube/config` file. If you delete the shared `.kube/config` file and you still have the `.kube-tkg/config` file for the management cluster, you can recover the `.kube/config` of the workload clusters with the `tanzu cluster kubeconfig get <my-cluster>` command.

By default, information about management clusters is stored in a separate `.kube-tkg/config` file.

Do not change context or edit the `.kube-tkg/config` or `.kube/config` files while Workload Grid operations are running.

<!--## <a id="vsphere-with-tanzu"></a> Using the Tanzu CLI to Create and Manage Clusters in vSphere with Tanzu

If you have vSphere 7 and you have enabled the vSphere with Tanzu feature, you can use the Tanzu CLI to interact with the vSphere with Tanzu Supervisor Cluster, to deploy Workload clusters in vSphere with Tanzu. For more information, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](connect-vsphere7.md).-->