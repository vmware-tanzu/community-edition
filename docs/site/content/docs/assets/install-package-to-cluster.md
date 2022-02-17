1. List the available package repositories.

    ```sh
    tanzu package repository list --all-namespaces
    ```

1. Verify the `STATUS` is `Reconcile succeeded` in the output of the above
   command.

    ```txt
    NAME                 REPOSITORY                                           TAG                     STATUS               DETAILS  NAMESPACE
    tkg-core-repository  projects.registry.vmware.com/tce/main                0.9.1                   Reconcile succeeded           tanzu-package-repo-global
    tkg-core-repository  projects.registry.vmware.com/tkg/packages/core/repo  v1.21.2_vmware.1-tkg.1  Reconcile succeeded           tkg-system
    ```

1. List the available packages.

    ```sh
    tanzu package available list
    ```

    ```txt
    NAME                                           DISPLAY-NAME        SHORT-DESCRIPTION
    cert-manager.community.tanzu.vmware.com        cert-manager        Certificate management
    contour.community.tanzu.vmware.com             Contour             An ingress controller
    external-dns.community.tanzu.vmware.com        external-dns        This package provides DNS synchronization functionality.
    fluent-bit.community.tanzu.vmware.com          fluent-bit          Fluent Bit is a fast Log Processor and Forwarder
    gatekeeper.community.tanzu.vmware.com          gatekeeper          policy management
    grafana.community.tanzu.vmware.com             grafana             Visualization and analytics software
    harbor.community.tanzu.vmware.com              Harbor              OCI Registry
    knative-serving.community.tanzu.vmware.com     knative-serving     Knative Serving builds on Kubernetes to support deploying and serving of applications and functions as serverless containers
    local-path-storage.community.tanzu.vmware.com  local-path-storage  This package provides local path node storage and primarily supports RWO AccessMode.
    multus-cni.community.tanzu.vmware.com          multus-cni          This package provides the ability for enabling attaching multiple network interfaces to pods in Kubernetes
    prometheus.community.tanzu.vmware.com          prometheus          A time series database for your metrics
    velero.community.tanzu.vmware.com              velero              Disaster recovery capabilities
    ```

1. List the available versions of cert-manager.

    ```sh
    tanzu package available list cert-manager.community.tanzu.vmware.com
    ```

    ```txt
    NAME                                     VERSION  RELEASED-AT
    cert-manager.community.tanzu.vmware.com  1.3.3    2021-08-06 05:31:21 -0700 MST
    cert-manager.community.tanzu.vmware.com  1.4.4    2021-08-23 09:47:51 -0700 MST
    cert-manager.community.tanzu.vmware.com  1.5.3    2021-08-23 10:22:51 -0700 MST
    ```

1. Install the `1.5.3` version of cert-manager.

    ```sh
    tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version 1.5.3
    ```

    > If you have `kubectl` installed, you can now view and interact with
    > cert-manager.

1. Verify the package is now installed.

    ```sh
    tanzu package installed list
    ```

    ```txt
    NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
    cert-manager  cert-manager.community.tanzu.vmware.com  1.5.3            Reconcile succeeded
    ```

1. If desired, you can also delete the package using:

    ```sh
    tanzu package installed delete cert-manager
    ```
