# Deploying Tanzu Kubernetes Clusters

This section describes how you use the Tanzu CLI to deploy and manage Tanzu Kubernetes clusters.

Before you can create Tanzu Kubernetes clusters, you must install the Tanzu CLI and deploy a management cluster. For information, see [Install the Tanzu CLI and Other Tools](../install-cli.md) and [Deploying Management Clusters](../mgmt-clusters/deploy-management-clusters.md).

You can use the Tanzu CLI to deploy Tanzu Kubernetes clusters to the following platforms:

- vSphere 6.7u3
- vSphere 7 (see [below](#vsphere-with-tanzu))
- Amazon EC2
- Microsoft Azure

## About Tanzu Kubernetes Clusters

In VMware Tanzu Kubernetes Grid, Tanzu Kubernetes clusters are the Kubernetes clusters in which your application workloads run.

Tanzu Kubernetes Grid automatically deploys clusters to the platform on which you deployed the management cluster. For example, you cannot deploy clusters to Amazon EC2 or Azure from a management cluster that is running in vSphere, or the reverse. It is not possible to use shared services between the different providers because, for example, vSphere clusters are reliant on sharing vSphere networks and storage, while Amazon EC2 and Azure use their own systems. Tanzu Kubernetes Grid automatically deploys clusters from whichever management cluster you have set as the context for the CLI by using the `tanzu login` command. For information about `tanzu login`, see [Manage Your Management Clusters](../cluster-lifecycle/multiple-management-clusters.md).

- For information about how to use the Tanzu CLI to deploy Tanzu Kubernetes clusters, see [Deploy Tanzu Kubernetes Clusters](deploy.md) and its subtopics.
- After you have deployed Tanzu Kubernetes clusters, the Tanzu CLI provides commands and options to perform the following cluster lifecycle management operations. See [Managing Cluster Lifecycles](../cluster-lifecycle/index.md).

For information about how to upgrade existing clusters to a new version of Kubernetes, see [Upgrade Tanzu Kubernetes Clusters](../upgrade-tkg/workload-clusters.md).

## <a id="kubectl"></a> Tanzu Kubernetes Clusters, `kubectl`, and `kubeconfig`

When you create a management cluster, the Tanzu CLI and `kubectl` contexts are automatically set to that management cluster. However, Tanzu Kubernetes Grid does not automatically set the `kubectl` context to a Tanzu Kubernetes cluster when you create it. You must set the `kubectl` context to a Tanzu Kubernetes cluster manually by using the `kubectl config use-context` command.

By default, unless you specify the `KUBECONFIG` option to save the `kubeconfig` for a cluster to a specific file, all Tanzu Kubernetes clusters that you deploy are added to a shared `.kube/config` file. If you delete the shared `.kube/config` file and you still have the `.kube-tkg/config` file for the management cluster, you can recover the `.kube/config` of the Tanzu Kubernetes clusters with the `tanzu cluster kubeconfig get <my-cluster>` command.

By default, information about management clusters is stored in a separate `.kube-tkg/config` file.

Do not change context or edit the `.kube-tkg/config` or `.kube/config` files while Tanzu Kubernetes Grid operations are running.

## <a id="vsphere-with-tanzu"></a> Using the Tanzu CLI to Create and Manage Clusters in vSphere with Tanzu

If you have vSphere 7 and you have enabled the vSphere with Tanzu feature, you can use the Tanzu CLI to interact with the vSphere with Tanzu Supervisor Cluster, to deploy Tanzu Kubernetes clusters in vSphere with Tanzu. For more information, see [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](connect-vsphere7.md).