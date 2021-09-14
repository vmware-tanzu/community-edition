## Installing and Managing Packages
You can discover and deploy packages through the Tanzu CLI. Packages extend the functionality of Tanzu Community Edition.

### Before you begin
Ensure you have deployed either a management/workload cluster or a standalone cluster.

### Procedure

1. Make sure your `kubectl` context is set to either the workload cluster or standalone cluster.

    ```sh
    kubectl config use-context <CLUSTER-NAME>-admin@<CLUSTER-NAME>
    ```
    > Where ``<CLUSTER-NAME>`` is the name of workload or standalone cluster where you want to install package.


1. Install the Tanzu Community Edition package repository.

    ```sh
    tanzu package repository add tce-repo --url projects.registry.vmware.com/tce/main:stable
    ```

   > By installing the Tanzu Community Edition package repository, [kapp-controller](https://carvel.dev/kapp-controller/) will make multiple packages available in the cluster.

1. [Optional] Verify the package repository has reconciled.

    ```sh
    tanzu package repository list
    ```

    The output will look similar to the following:

    ```sh
    / Retrieving repositories...
      NAME      REPOSITORY                                    STATUS
    DETAILS
      tce-repo  projects.registry.vmware.com/tce/main:stable  Reconcile succeeded
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
     contour-operator.community.tanzu.vmware.com    contour-operator    Layer 7 Ingress
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

1. [Optional]: Get additional information about a package

    ```shell
    tanzu package available get cert-manager.community.tanzu.vmware.com
    ```
    The output will look similar to the following:
    ```sh
    / Retrieving package details for cert-manager.community.tanzu.vmware.com...
    NAME:               cert-manager.community.tanzu.vmware.com
    DISPLAY-NAME:       cert-manager
    SHORT-DESCRIPTION:  Certificate management
    PACKAGE-PROVIDER:   VMware
    LONG-DESCRIPTION:   Provides certificate management provisioning within the cluster
    MAINTAINERS:        [{Nicholas Seemiller}]
    ```

1. [Optional]: See available package versions

    ```shell
    tanzu package available list cert-manager.community.tanzu.vmware.com
    ```
    The output will look similar to the following:
    ```sh
    / Retrieving package versions for cert-manager.community.tanzu.vmware.com...
    NAME                                     VERSION  RELEASED-AT
    cert-manager.community.tanzu.vmware.com  1.3.1    2021-04-14T18:00:00Z
    cert-manager.community.tanzu.vmware.com  1.4.0    2021-06-15T18:00:00Z
    ```

1. [Optional]: Download the configuration for a package. For the moment, you will need to refer to the
   [TCE GitHub repository](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages). Select the package/version
   and navigate into the `bundle/config` directory. Download or copy/paste the `values.yaml` file.

1. [Optional]: Alter the ``values.yaml`` file using a text editor of your choice. Ensure that there are no lines in the file starting with `#!` or `#@`, as these cause an error when installing to the cluster.

1. Install the package to the cluster.

    ```sh
    tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version 1.4.0
    ```
    > If using a custom configuration values file, append `--values-file values.yaml` to the installation command. <br>

    The output will look similar to the following:

    ```sh
    | Installing package 'cert-manager.community.tanzu.vmware.com'
    / Getting package metadata for cert-manager.community.tanzu.vmware.com
    - Creating service account 'cert-manager-default-sa'
    \ Creating cluster admin role 'cert-manager-default-cluster-role'

    - Creating package resource
    / Package install status: Reconciling

    Added installed package 'cert-manager' in namespace 'default'
    ```



1. Verify cert-manager is installed in the cluster.

     ```sh
     tanzu package installed list
     | Retrieving installed packages...
     NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
     cert-manager  cert-manager.community.tanzu.vmware.com  1.4.0            Reconcile succeeded
     ```

1. For troubleshooting, you can view `PackageInstall` and `App` objects in the cluster using the following kubectl command.

     ```sh
     kubectl get packageInstall,apps --all-namespaces
     ```
     The output will look similar to the following:
     ```sh
     NAMESPACE    NAME                                                 PACKAGE NAME                              PACKAGE VERSION                    DESCRIPTION           AGE
     default      packageinstall.packaging.carvel.dev/cert-manager     cert-manager.community.tanzu.vmware.com   1.4.0                              Reconcile succeeded   18m
     tkg-system   packageinstall.packaging.carvel.dev/antrea           antrea.tanzu.vmware.com                   0.13.3+vmware.1-tkg.1-zshippable   Reconcile succeeded   17d
     tkg-system   packageinstall.packaging.carvel.dev/metrics-server   metrics-server.tanzu.vmware.com           0.4.0+vmware.1-tkg.1-zshippable    Reconcile succeeded   17d

     NAMESPACE    NAME                                  DESCRIPTION           SINCE-DEPLOY   AGE
     default      app.kappctrl.k14s.io/cert-manager     Reconcile succeeded   12s            18m
     tkg-system   app.kappctrl.k14s.io/antrea           Reconcile succeeded   24s            17d
     tkg-system   app.kappctrl.k14s.io/metrics-server   Reconcile succeeded   28s            17d
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

For more information about how this package model works from a server-side and client-side perspective, see the
[Package Management design doc](./designs/package-management.md).

_Note:_ For installation of packages that require persistent storage
(like Prometheus or Grafana), you must _first_ have a [`StorageClass`](https://kubernetes.io/docs/concepts/storage/storage-classes/) installed.
TCE provides the [`local-path-storage`](../package-readme-local-path-storage-0.0.19) package for quick and easy local node storage.
It is primarily intended to be used with Docker, but will work on any infrastructure.
This installs a default storage class if none is available
and enables persistent storage to the local node's filesystem.
But it has limitations in multi-node deployments and is not intended to dynamically change hosts.
More information can be found in the [`local-path-storage` package documentation.](../package-readme-local-path-storage-0.0.19)
