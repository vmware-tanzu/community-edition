## Installing and Managing Packages

With a cluster bootstrapped, you're ready to configure and install packages to the cluster.

1. Make sure your `kubectl` context is set to the workload cluster.

    ```sh
    kubectl config use-context ${GUEST_CLUSTER_NAME}-admin@${GUEST_CLUSTER_NAME}
    ```

1. Install the TCE package repository.

    ```sh
    tanzu package repository install --default
    ```

   > By installing the TCE package repository, kapp-controller will make multiple packages available in the cluster.

1. List the available packages.

    ```md
    tanzu package list

    NAME                             VERSION         DESCRIPTION
    cert-manager.tce.vmware.com      1.1.0-vmware0   This package provides certificate management functionality.
    contour-operator.tce.vmware.com  1.11.0-vmware0  This package provides an ingress controller.
    external-dns.tce.vmware.com      0.7.6-vmware0   This package provides external DNS capabilities.
    fluent-bit.tce.vmware.com        1.7.2-vmware0   Fluent Bit is an open source Log Processor and Forwarder.
    gatekeeper.tce.vmware.com        3.2.3-vmware0   This package provides custom admission control.
    grafana.tce.vmware.com           7.4.3-vmware0   Grafana is open source visualization and analytics software.
    knative-serving.tce.vmware.com   0.22.0-vmware0  This package provides serverless functionality.
    prometheus.tce.vmware.com        2.25.0-vmware0  A time series database for your metrics.
    velero.tce.vmware.com            1.5.2-vmware0   This package provides disaster recovery capabilities.
    ```

1. [Optional]: Download the configuration for a package.

   ```md
   tanzu package configure fluent-bit.tce.vmware.com

   Looking up config for package: fluent-bit.tce.vmware.com:
   Values files saved to fluent-bit.tce.vmware.com-values.yaml. Configure this file before installing the package.
   ```

1. [Optional]: Alter the values files.

   ```sh
   vim fluent-bit.tce.vmware.com-values.yaml
   ```

1. Install the package to the cluster.

    ```sh
    tanzu package install fluent-bit.tce.vmware.com --config fluent-bit.tce.vmware.com-values.yaml

    Looking up package to install: fluent-bit.tce.vmware.com:
    Installed package in default/fluent-bit.tce.vmware.com:1.7.2-vmware0
   ```

   > The `--config` flag is optional based on whether you customized the configuration file from the previous steps.

1. Verify fluent-bit is installed in the cluster.

    ```sh
    kubectl -n fluent-bit get all
    pod/fluent-bit-hgtc2   1/1     Running   0          27m
    pod/fluent-bit-j6jdj   1/1     Running   0          27m

    NAME                        DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
    daemonset.apps/fluent-bit   2         2         2       2            2           <none>          27m
    ```

1. For troubleshooting, you can view `InstalledPackage` and `App` objects in the cluster.

    ```sh
    kubectl get installedpackage,apps --all-namespaces

    NAMESPACE         NAME                    DESCRIPTION           SINCE-DEPLOY   AGE
    default           gatekeeper              Reconcile succeeded   13s            16s
    kapp-controller   tce-main.tanzu.vmware   Reconcile succeeded   17s            2m
    tkg-system        antrea                  Reconcile succeeded   116s           19h
    tkg-system        metrics-server          Reconcile succeeded   2m10s          19h
    ```

If you're interested in how this package model works from a server-side and client-side perspective, please read our
[Tanzu Add-on Management design doc](./designs/tanzu-addon-management.md).

