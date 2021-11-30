## Installing a Package

This section walks you through installing the cert-manager package in your cluster as an example of package installation.

{{% include "/docs/assets/cert-manager-desc.md" %}}

For detailed instruction on package management, see [Work with Packages](../package-management).

### Prerequisites

- Before you install packages, you should have one of the following cluster configurations running:

  - A [management cluster](https://tanzucommunityedition.io/docs/latest/glossary/#management-cluster) and a [workload cluster](https://tanzucommunityedition.io/docs/latest/glossary/#workload-cluster).

    **or**

  - A [standalone cluster](https://tanzucommunityedition.io/docs/latest/glossary/#standalone-cluster)
- If you deployed a management/workload cluster, you will install cert-manager in the workload cluster. If you deployed a standalone cluster, you will install cert-manager in the standalone cluster.

For more information, see [Planning Your Installation](https://tanzucommunityedition.io/docs/latest/installation-planning/).

### Procedure

1. Make sure your `kubectl` context is set to either the workload cluster or standalone cluster. See Prerequisites above.

    ```sh
    kubectl config use-context <CLUSTER-NAME>-admin@<CLUSTER-NAME>
    ```

    Where ``<CLUSTER-NAME>`` is the name of workload or standalone cluster where you want to install a package.

1. Install the Tanzu Community Edition package repository into the `tanzu-package-repo-global` namespace.

    ```sh
    tanzu package repository add tce-repo --url projects.registry.vmware.com/tce/main:0.9.1 --namespace tanzu-package-repo-global
    ```

    > Package repositories installed into the `tanzu-package-repo-global` namespace are available to the entire cluster.  
    > Use the `--namespace` argument in the `tanzu package repository add` command to install a package repository into a specific namespace. If you install a package repository into another namespace, you must specify that namespace as an argument in the `tanzu package install` command  when you install a package from that repository.  
    > A `tanzu-core` repository is also installed in the `tkg-system` namespace
    > clusters. This repository holds lower-level components that are **not**
    > meant to be installed by the user! These packages are used during cluster
    > boostrapping.

1. Verify the package repository has reconciled.

    ```sh
    tanzu package repository list --namespace tanzu-package-repo-global
    ```

    The output will look similar to the following:

    ```sh
    / Retrieving repositories...
      NAME      REPOSITORY                                    STATUS
    DETAILS
      tce-repo  projects.registry.vmware.com/tce/main:0.9.1  Reconcile succeeded
    ```

    > It may take some time to see `Reconcile succeeded`. Until then, packages
    > won't show up in the available list described in the next step.

1. List the available packages.

    ```sh
    tanzu package available list
    ```

    The output will look similar to the following:

    ```sh
    - Retrieving available packages...
     NAME                                           DISPLAY-NAME        SHORT-DESCRIPTION
     cert-manager.community.tanzu.vmware.com        cert-manager        Certificate management
     contour.community.tanzu.vmware.com             Contour             An ingress controller
     external-dns.community.tanzu.vmware.com        external-dns        This package provides DNS...
     fluent-bit.community.tanzu.vmware.com          fluent-bit          Fluent Bit is a fast Log Processor and...
     gatekeeper.community.tanzu.vmware.com          gatekeeper          policy management
     grafana.community.tanzu.vmware.com             grafana             Visualization and analytics software
     harbor.community.tanzu.vmware.com              Harbor              OCI Registry
     knative-serving.community.tanzu.vmware.com     knative-serving     Knative Serving builds on Kubernetes to...
     local-path-storage.community.tanzu.vmware.com  local-path-storage  This package provides local path node...
     multus-cni.community.tanzu.vmware.com          multus-cni          This package provides the ability for...
     prometheus.community.tanzu.vmware.com          prometheus          A time series database for your metrics
     velero.community.tanzu.vmware.com              velero              Disaster recovery capabilities
    ```

1. List the available versions for the `cert-manager` package.

    ```shell
    tanzu package available list cert-manager.community.tanzu.vmware.com
    ```

    The output will look similar to the following:

    ```sh
    / Retrieving package versions for cert-manager.community.tanzu.vmware.com...
    NAME                                     VERSION  RELEASED-AT
    cert-manager.community.tanzu.vmware.com  1.3.3    2021-08-06T12:31:21Z
    cert-manager.community.tanzu.vmware.com  1.4.4    2021-08-23T16:47:51Z
    cert-manager.community.tanzu.vmware.com  1.5.3    2021-08-23T17:22:51Z
    ```

    **NOTE**: The available versions of a package may have changed since this guide was written.

1. Install the package to the cluster.

    ```sh
    tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version 1.5.3
    ```

    The output will look similar to the following:

    ```sh
    | Installing package 'cert-manager.community.tanzu.vmware.com'
    / Getting package metadata for cert-manager.community.tanzu.vmware.com
    - Creating service account 'cert-manager-default-sa'
    \ Creating cluster admin role 'cert-manager-default-cluster-role'
  
    Creating package resource
    / Package install status: Reconciling

    Added installed package 'cert-manager' in namespace 'default'

    ```

    **NOTE**: Use one of the available package versions, since the one described
    in this guide might no longer be available.

1. Verify cert-manager is installed in the cluster.

     ```sh
     tanzu package installed list
     ```

     The output will look similar to the following:

     ```sh
     | Retrieving installed packages...
     NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
     cert-manager  cert-manager.community.tanzu.vmware.com  1.5.3            Reconcile succeeded
     ```

1. To remove a package from the cluster, run the following command:

     ```shell
     tanzu package installed delete cert-manager
     ```

     The output will look similar to the following:

     ```sh
     | Uninstalling package 'cert-manager' from namespace 'default'
     | Getting package install for 'cert-manager'
     \ Deleting package install 'cert-manager' from namespace 'default'
     \ Package uninstall status: ReconcileSucceeded
     \ Package uninstall status: Reconciling
     \ Package uninstall status: Deleting
     | Deleting admin role 'cert-manager-default-cluster-role'

     / Deleting service account 'cert-manager-default-sa'
     Uninstalled package 'cert-manager' from namespace 'default'
     ```

For more information about package management, see [Work with Packages](../package-management). For details on installing a specific package,
see the package's documentations in the left navigation bar (`Packages >
${PACKAGE_NAME}`).
