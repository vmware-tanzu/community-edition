# Upgrade Tanzu Kubernetes Clusters

After you have upgraded a management cluster, you can [upgrade the Tanzu Kubernetes clusters](workload-clusters.md) that the management cluster manages.

**IMPORTANT**: Management clusters and Tanzu Kubernetes clusters use client certificates to authenticate clients. These certificates are valid for one year. To renew them, upgrade your clusters at least once a year.

## <a id="prereqs"></a> Prerequisites

- You performed the steps in [Upgrading Tanzu Kubernetes Grid](index.md) that occur before the step for upgrading Tanzu Kubernetes Clusters.
- You performed the steps in [Upgrade Management Clusters](management-cluster.md) to upgrade the management cluster that manages the Tanzu Kubernetes clusters that you want to upgrade.
- If you are upgrading clusters that run on vSphere, before you can upgrade clusters to a non-default version of Kubernetes for your version of Tanzu Kubernetes Grid, the appropriate base image template OVAs must be available in vSphere as VM templates. For information about importing OVA files into vSphere, see [Prepare to Upgrade Clusters on vSphere](index.md#vsphere).
- If you are upgrading clusters that run on Amazon EC2, the Amazon Linux 2 Amazon Machine Images (AMI) that include the supported Kubernetes versions are publicly available to all Amazon EC2 users, in all supported AWS regions. Tanzu Kubernetes Grid automatically uses the appropriate AMI for the Kubernetes version that you specify during upgrade.
- If you are upgrading clusters that run on Azure, ensure that you completed the
steps in [Prepare to Upgrade Clusters on Azure](index.md#azure).

## <a id="procedure"></a> Procedure

The upgrade process upgrades the version of Kubernetes in all of the control plane and worker nodes of your Tanzu Kubernetes clusters.

1. Run the `tanzu login` command to see an interactive list of available management clusters.

   ```
   tanzu login
   ```

1. Select a management cluster to switch the context of the Tanzu CLI. You should select the management cluster that manages the clusters you want to upgrade. See [List Management Clusters and Change Context](../cluster-lifecycle/multiple-management-clusters.md#login) for more information.

1. Run the `tanzu cluster list` command with the `--include-management-cluster` option.

   ```
   tanzu cluster list --include-management-cluster
   ```

   The `tanzu cluster list` command shows the version of Kubernetes that is running in the management cluster and all of the clusters that it manages. In this example, you can see that the management cluster has already been upgraded to v1.20.5, but the Tanzu Kubernetes clusters are running older versions of Kubernetes.

   ```
     NAME                 NAMESPACE   STATUS    CONTROLPLANE  WORKERS  KUBERNETES         ROLES       PLAN
     k8s-1-17-13-cluster  default     running   1/1           1/1      v1.17.13+vmware.1  <none>      dev
     k8s-1-18-10-cluster  default     running   1/1           1/1      v1.18.10+vmware.1  <none>      dev
     k8s-1-19-3-cluster   default     running   1/1           1/1      v1.19.3+vmware.1   <none>      dev
     mgmt-cluster         tkg-system  running   1/1           1/1      v1.20.5+vmware.1   management  dev
   ```

1. If your v1.2 management cluster used the connectivity API to support Harbor access, you need to remove `tanzu-system-connectivity` artifacts from each workload cluster as follows:

   1. Set the context of `kubectl` to the context of your workload cluster:

      ```
      kubectl config use-context WORKLOAD-CLUSTER-admin@WORKLOAD-CLUSTER
      ```

      Where `WORKLOAD-CLUSTER` is the name of your workload cluster.

   1. Delete the `tanzu-system-connectivity` namespace:

      ```
      kubectl delete ns tanzu-system-connectivity
      ```

  **Note**: If your Harbor registry in v1.2 used a fictitious domain name such as `harbor.system.tanzu` instead of a FQDN, you cannot upgrade workload clusters automatically, but must instead create a new v1.3 workload cluster and migrate the workloads to the new cluster manually.

1. Before you upgrade a Tanzu Kubernetes cluster, remove all unmanaged `kapp-controller` deployment artifacts from the Tanzu Kubernetes cluster. An unmanaged `kapp-controller` deployment is a deployment that exists outside of the `vmware-system-tmc` namespace. You can assume it is in the `kapp-controller` namespace.

    1. Get the credentials of the cluster.

       ```
       tanzu cluster kubeconfig get CLUSTER-NAME --admin
       ```

       For example, using cluster k8s-1-19-3-cluster:

       ```
       tanzu cluster kubeconfig get k8s-1-19-3-cluster --admin
       ```

    1. Set the context of kubectl to the cluster:

       ```
       kubectl config use-context k8s-1-19-3-cluster-admin@k8s-1-19-3-cluster
       ```

    1. Delete the `kapp-controller` deployment on the cluster.

       ```
       kubectl delete deployment kapp-controller -n kapp-controller
       ```

       **Note:** If you receive a `NotFound` error message, ignore the error. You should continue with the following deletion steps in case you have any orphaned objects related to a pre-existing `kapp-controller` deployment.

       ```
       Error from server (NotFound): deployments.apps "kapp-controller" not found
       ```

    1. Delete all `kapp-controller` objects.

       ```
       kubectl delete clusterrole kapp-controller-cluster-role
       kubectl delete clusterrolebinding kapp-controller-cluster-role-binding
       kubectl delete serviceaccount kapp-controller-sa -n kapp-controller
       ```

1. To discover which versions of Kubernetes are made available by a management cluster, run the `tanzu kubernetes-release get` command.

   ```
   tanzu kubernetes-release get
   ```

   The output lists all of the versions of Kubernetes that you can use to deploy clusters, with the following notes:
      - `COMPATIBLE`: The current management cluster can deploy workload clusters with this Tanzu Kubernetes release (`tkr`).
      - `UPGRADEAVAILABLE`: This `tkr` is not the most current in its Kubernetes version line. Any workload clusters running this `tkr` version can be upgraded to newer versions.

   For example:

   ```
     NAME                       VERSION                  COMPATIBLE  UPGRADEAVAILABLE
     v1.17.16---vmware.2-tkg.2  v1.17.16+vmware.2-tkg.2  True        True
     v1.18.16---vmware.1-tkg.2  v1.18.16+vmware.1-tkg.2  True        True
     v1.18.17---vmware.1-tkg.1  v1.18.17+vmware.1-tkg.1  True        True
     v1.19.8---vmware.1-tkg.2   v1.19.8+vmware.1-tkg.2   True        True
     v1.19.9---vmware.1-tkg.1   v1.19.9+vmware.1-tkg.1   True        True
     v1.20.4---vmware.1-tkg.2   v1.20.4+vmware.1-tkg.2   True        True
     v1.20.5---vmware.1-tkg.1   v1.20.5+vmware.1-tkg.1   True        False
   ```

1. To discover the newer `tkr` versions to which you can upgrade a workload cluster running an older `tkr` version, run the `tanzu kubernetes-release available-upgrades get` command.

   ```
   tanzu kubernetes-release available-upgrades get v1.19.8---vmware.1-tkg.1
     NAME                      VERSION
     v1.19.9---vmware.1-tkg.1  v1.19.9+vmware.1-tkg.1 
     v1.20.4---vmware.1-tkg.2  v1.20.4+vmware.1-tkg.2
     v1.20.5---vmware.1-tkg.1  v1.20.5+vmware.1-tkg.1     
   ```

   You cannot skip minor versions when upgrading your `tkr` version. For example, you cannot upgrade a cluster directly from v1.18.x to v1.20.x. You must upgrade a v1.18.x cluster to v1.19.x before upgrading the cluster to v1.20.x.

1. Run the `tanzu cluster upgrade CLUSTER-NAME` command and enter `y` to confirm.

   To upgrade the cluster to the default version of Kubernetes for this release of Tanzu Kubernetes Grid, run the `tanzu cluster upgrade` command without any options. For example, the following command upgrades the cluster `k8s-1-19-3-cluster` from v1.19.3 to v1.20.4.

   ```
   tanzu cluster upgrade k8s-1-19-3-cluster
   ```

   If the cluster is not running in the `default` namespace, specify the `--namespace` option.

   ```
   tanzu cluster upgrade CLUSTER-NAME --namespace NAMESPACE-NAME
   ```

   To skip the confirmation step when you upgrade a cluster, specify the `--yes` option.

   ```
   tanzu cluster upgrade CLUSTER-NAME --yes
   ```

   If an upgrade times out before it completes, run `tanzu cluster upgrade` again and specify the `--timeout` option with a value greater than the default of 30 minutes.

   ```
   tanzu cluster upgrade CLUSTER-NAME --timeout 45m0s
   ```

   If multiple base VM images in your IaaS account have the same version of Kubernetes that you are upgrading to, use the `--os-name` option to specify the OS you want.
   See [Selecting an OS During Cluster Upgrade](cluster-os-upgrade.md) for more information.

   For example, on vSphere if you have uploaded both Photon and Ubuntu OVA templates with Kubernetes v1.20.5, specify `--os-name ubuntu` to upgrade your workload cluster to run on an Ubuntu VM.

   ```
   tanzu cluster upgrade CLUSTER-NAME --os-name ubuntu
   ```

   Since you cannot skip minor versions of `tkr`, the upgrade command fails if you try to upgrade a cluster that is more than one minor version behind the default version. For example, you cannot upgrade directly from v1.18.x to v1.20.x. To upgrade a cluster to a version of Kubernetes that is not the default version for this release of Tanzu Kubernetes Grid, specify the `--tkr` option with the `NAME` of the chosen version, as listed by `tanzu kubernetes-release get` above. For example, to upgrade the cluster `k8s-1-18-10-cluster` from v1.18.10 to v1.19.8.

   ```
   tanzu cluster upgrade k8s-1-18-10-cluster --tkr v1.19.9---vmware.1-tkg.1 --yes
   ```

1. When the upgrade finishes, run the `tanzu cluster list` command with the `--include-management-cluster` option again, to check that the Tanzu Kubernetes cluster has been upgraded.

   ```
   tanzu cluster list --include-management-cluster
   ```

   You see that the `k8s-1-17-13-cluster` and `k8s-1-19-3-cluster` Tanzu Kubernetes clusters are now running Kubernetes v1.18.17 and v1.20.5 respectively.

   ```
     NAME                 NAMESPACE   STATUS    CONTROLPLANE  WORKERS  KUBERNETES         ROLES       PLAN
     k8s-1-17-13-cluster  default     running   1/1           1/1      v1.18.17+vmware.1  <none>      dev
     k8s-1-18-10-cluster  default     running   1/1           1/1      v1.19.9+vmware.1   <none>      dev
     k8s-1-19-3-cluster   default     running   1/1           1/1      v1.20.5+vmware.1   <none>      dev
     mgmt-cluster         tkg-system  running   1/1           1/1      v1.20.5+vmware.1   management  dev
   ```

## <a id="what-next"></a> What to Do Next

You can now continue to use the Tanzu CLI to manage your clusters, and run your applications with the new version of Kubernetes.

To complete the upgrade, you should upgrade any extensions you have deployed such as Contour, Fluent Bit or Prometheus that are running on your Tanzu Kubernetes clusters.

You must also register any add-ons such as CNI, vSphere CPI, Pinniped or Metrics Server that you will be using in your Tanzu Kubernetes Grid deployment.

For more information on upgrading extensions, see [Upgrade Tanzu Kubernetes Grid Extensions](extensions.md).

For instructions on how to register add-ons after upgrading your clusters from Tanzu Kubernetes Grid v1.2.x to v1.3.x, see [Register Core Add-ons](addons.md).
