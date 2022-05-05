# Upgrade Management & Workload Cluster

The management cluster upgrade process first upgrades the Cluster API providers for vSphere, Amazon EC2, or Azure that are running in the management cluster. Then, it upgrades the version of Kubernetes in all of the control planes and worker nodes of the management cluster.

## Prerequisites

- vSphere: If you are upgrading clusters that run on vSphere, the appropriate base image template OVAs must be available in vSphere as VM templates. For more information, see steps 1, and 2 in the [Prepare to Deploy a Management Cluster to vSphere](https://tanzucommunityedition.io/docs/v0.11/vsphere/#procedure).
- Amazon EC2: If you are upgrading clusters that run on Amazon EC2, the Amazon Linux 2 Amazon Machine Images (AMI) that include the supported Kubernetes versions are publicly available to all Amazon EC2 users, in all supported AWS regions. Tanzu Community Edition automatically uses the appropriate AMI for the Kubernetes version that you specify during an upgrade.
- Azure: If you are upgrading clusters that run on Azure, you must accept the terms for the new default VM image and each non-default VM image that you plan to use for your cluster VMs. To accept the terms, complete the following steps:

  1. List all available VM images in the Azure Marketplace:

        ```sh
        az vm image list --publisher vmware-inc --offer tkg-capi --all
        ```

  1. Accept the terms for the new default VM image:

        ```sh
        az vm image terms accept --urn publisher:offer:sku:version
        ```

  1. For example, to accept the terms for k8s-1dot21dot2-ubuntu-2004, run:

        ```sh
        az vm image terms accept --urn vmware-inc:tkg-capi:k8s-1dot21dot2-ubuntu-2004:2021.05.17
        ```

## Upgrade Management Cluster

1. Run the tanzu login command to see an interactive list of management clusters available for upgrade.

    ```sh
    tanzu login
    ```

    To change your current login context, use your up and down arrow keys to highlight a management cluster and then press Enter.

1. The `tanzu login` command does not automatically set the kubectl context, run the following commands to set the context.

    First, capture the management cluster’s kubeconfig and take note of the command for accessing the cluster in the output message:

    ```sh
    tanzu management-cluster kubeconfig get <MGMT-CLUSTER-NAME> --admin
    ```

    Set your kubectl context to the management cluster:

    ```sh
    kubectl config use-context <MGMT-CLUSTER-NAME>-admin@<MGMT-CLUSTER-NAME>
    ```

1. Run the `tanzu management-cluster upgrade` command:

    ```sh
    tanzu management-cluster upgrade
    ```

    ```txt
   --namespace   #The namespace where the workload cluster was created.
   --os-arch     #OS architecture to use during cluster upgrade.
   --os-name     #OS name to use during cluster upgrade.
   --os-version  #OS version to use during cluster upgrade.
   --timeout     #Time duration to wait for an operation before timeout.
   --tkr         #TanzuKubernetesRelease(TKr) to upgrade to.
    ```

    Options for the different cloud providers are:

    vSphere:

    ```sh
    --os-name ubuntu --os-version 20.04 --os-arch amd64
    --os-name photon --os-version 3 --os-arch amd64
    ```

    AWS:

    ```sh
    --os-name ubuntu --os-version 20.04 --os-arch amd64
    --os-name amazon --os-version 2 --os-arch amd64
    ```

    Azure:

    ```sh
    --os-name ubuntu --os-version 20.04 --os-arch amd64
    --os-name ubuntu --os-version 18.04 --os-arch amd64
    ```

1. When the upgrade finishes, run the `tanzu cluster list` command with the `--include-management-cluster` option again to check that the management cluster has been upgraded.

    ```sh
    tanzu cluster list --include-management-cluster
    ```

    You should see the management cluster is now running the new version of Kubernetes, but the workload clusters are still running previous versions of Kubernetes.

    <!-- need sample here -->

1. Regenerate the kubeconfig by running:

    ```sh
    tanzu cluster kubeconfig get <MANAGEMENT-CLUSTER-NAME>
    ```

## Upgrade Workload Clusters

After you have upgraded a management cluster, you can upgrade workload clusters.

1. Run the `tanzu login` command to see an interactive list of available management clusters.

   ```sh
   tanzu login
   ```

1. Select a management cluster to switch the context of the Tanzu CLI. You should select the management cluster that manages the clusters you want to upgrade.

1. Run the `tanzu cluster list` command with the `--include-management-cluster` option.

   ```sh
   tanzu cluster list --include-management-cluster
   ```

   The `tanzu cluster list` command shows the version of Kubernetes that is running in the management cluster and all of the clusters that it manages. In this example, you can see that the management cluster has already been upgraded to v1.22.3, but the workload clusters are running older versions of Kubernetes.

   ```txt
     NAME                 NAMESPACE   STATUS    CONTROLPLANE  WORKERS  KUBERNETES         ROLES       PLAN
     k8s-1-20-8-cluster   default     running   1/1           1/1      v1.20-8+vmware.1   <none>      dev
     k8s-1-21-2-cluster   default     running   1/1           1/1      v1.21-2+vmware.1   <none>      dev
     mgmt-cluster         tkg-system  running   1/1           1/1      v1.22.3+vmware.1   management  dev
   ```

1. To discover which versions of Kubernetes are made available by a management cluster, run the `tanzu kubernetes-release get` command.

   ```sh
   tanzu kubernetes-release get
   ```

   The output lists all of the versions of Kubernetes that you can use to deploy clusters, with the following notes:
      - `COMPATIBLE`: The current management cluster can deploy workload clusters with this Tanzu Kubernetes release (`tkr`).
      - `UPGRADEAVAILABLE`: This `tkr` is not the most current in its Kubernetes version line. Any workload clusters running this `tkr` version can be upgraded to newer versions.

   For example:

   ```txt
     NAME                        VERSION                   COMPATIBLE  UPGRADEAVAILABLE
     v1.19.16---vmware.1-tkg.1   v1.19.16+vmware.1-tkg.1   True        False
     v1.20.8---vmware.1-tkg.1    v1.20.8+vmware.1-tkg.1    True        True
     v1.20.12---vmware.1-tkg.1   v1.20.12+vmware.1-tkg.1   True        True
     v1.21.2---vmware.1-tkg.1    v1.21.2+vmware.1-tkg.1    True        True
     v1.21.6---vmware.1-tkg.1    v1.21.6+vmware.1-tkg.1    True        False
     v1.22.3---vmware.1-tkg.1    v1.22.3+vmware.1-tkg.1    True        False
   ```

1. To discover the newer `tkr` versions you can upgrade to, run the `tanzu kubernetes-release available-upgrades get` command, specifying the current version of the cluster.

   ```sh
   tanzu kubernetes-release available-upgrades get v1.20.8---vmware.1-tkg.1
   ```

   This command lists all of the available Kubernetes versions to which you can upgrade clusters that are running the specified version.

   ```txt
    NAME                       VERSION
    v1.20.12---vmware.1-tkg.1  v1.20.12+vmware.1-tkg.1
    v1.21.2---vmware.1-tkg.1   v1.21.2+vmware.1-tkg.1
    v1.21.6---vmware.1-tkg.1   v1.21.6+vmware.1-tkg.1
   ```

   You can also discover the `tkr` versions that are available for a specific workload cluster by specifying the cluster name in the `tanzu cluster available-upgrades get` command.

   ```sh
   tanzu cluster available-upgrades get k8s-1-20-8-cluster
   ```

   This command lists all of the Kubernetes versions that are compatible with the specified cluster.

   ```txt
   NAME                         VERSION                            COMPATIBLE
   v1.20.8---vmware.1-tkg.1     v1.20.8+vmware.1-tkg.1             True
   v1.20.12---vmware.1-tkg.2    v1.20.12+vmware.1-tkg.2            True
   v1.21.2---vmware.1-tkg.1     v1.21.2+vmware.1-tkg.1             True
   v1.21.6---vmware.1-tkg.1     v1.21.6+vmware.1-tkg.1             True
   v1.22.3---vmware.1-tkg.1     v1.22.3+vmware.1-tkg.1             True
   ```

   You cannot skip minor versions when upgrading your `tkr` version. For example, you cannot upgrade a cluster directly from v1.20.x to v1.22.x. You must upgrade a v1.20.x cluster to v1.21.x before upgrading the cluster to v1.22.x.

1. Run the `tanzu cluster upgrade` command with the following options

   ```sh
   tanzu cluster upgrade <WORKLOAD-CLUSTER-NAME>
   ```

   ```txt
   --namespace   #The namespace where the workload cluster was created.
   --os-arch     #OS arch to use during cluster upgrade.
   --os-name     #OS name to use during cluster upgrade.
   --os-version  #OS version to use during cluster upgrade.
   --timeout     #Time duration to wait for an operation before timeout.
   --tkr string  #TanzuKubernetesRelease(TKr) to upgrade to.
   ```

    Options for the different cloud providers are:

    vSphere:

    ```sh
    --os-name ubuntu --os-version 20.04 --os-arch amd64
    --os-name photon --os-version 3 --os-arch amd64
    ```

    AWS:

    ```sh
    --os-name ubuntu --os-version 20.04 --os-arch amd64
    --os-name amazon --os-version 2 --os-arch amd64
    ```

    Azure:

    ```sh
    --os-name ubuntu --os-version 20.04 --os-arch amd64
    --os-name ubuntu --os-version 18.04 --os-arch amd64
    ```

1. When the upgrade finishes, run the `tanzu cluster list` command with the `--include-management-cluster` option again, to check that the workload cluster has been upgraded.

   ```sh
   tanzu cluster list --include-management-cluster
   ```

   You see that the `k8s-1-20-8-cluster` and `k8s-1-21-2-cluster` workload clusters are now running Kubernetes v1.21.6 and v1.22.3 respectively.

   ```txt
     NAME                 NAMESPACE   STATUS    CONTROLPLANE  WORKERS  KUBERNETES         ROLES       PLAN
     k8s-1-20-8-cluster   default     running   1/1           1/1      v1.21.6+vmware.1  <none>      dev
     k8s-1-21-2-cluster   default     running   1/1           1/1      v1.22.3+vmware.1   <none>      dev
     mgmt-cluster         tkg-system  running   1/1           1/1
     v1.22.3+vmware.1   management  dev
   ```

1. Regenerate the `kubeconfig` by running:

   ```sh
   tanzu cluster kubeconfig get <WORKLOAD-CLUSTER-NAME>
   ```
