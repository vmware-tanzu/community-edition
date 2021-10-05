## Installing a Package

This section walks you through installing a package (cert-manager) in your cluster. For
detailed instruction on package management, see [Work with Packages](../package-management).

1. Make sure your `kubectl` context is set to either the workload cluster or standalone cluster.

    ```sh
    kubectl config use-context <CLUSTER-NAME>-admin@<CLUSTER-NAME>
    ```

    Where ``<CLUSTER-NAME>`` is the name of workload or standalone cluster where you want to install package.

1. Install the Tanzu Community Edition package repository into the `tanzu-package-repo-global` namespace.

    ```sh
    tanzu package repository add tce-repo --url projects.registry.vmware.com/tce/main:0.9.1 --namespace tanzu-package-repo-global
    ```

    > Repositories installed into the `tanzu-package-repo-global` namespace will provide their packages to the entire
    > cluster. It is possible to install package repositories into specific namespaces when using the `--namespace` argument.
    > To install a package from a repository in another namespace will require you to specify that namespace as an argument
    > to the `tanzu package install` command.
                                                                                                                                
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

    > A `tanzu-core` repository is also installed in the `tkg-system` namespace
    > clusters. This repository holds lower-level components that are **not**
    > meant to be installed by the user! These packages are used during cluster
    > boostrapping.
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

1. List the available versions for the `cert-manager` package.

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

1. Install the package to the cluster.

    ```sh
    tanzu package install cert-manager \
      --package-name cert-manager.community.tanzu.vmware.com \
      --version 1.4.0
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

1. Verify cert-manager is installed in the cluster.

     ```sh
     tanzu package installed list
     ```

     The output will look similar to the following:

     ```sh
     | Retrieving installed packages...
     NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
     cert-manager  cert-manager.community.tanzu.vmware.com  1.4.0            Reconcile succeeded
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
