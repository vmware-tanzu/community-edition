<!-- markdownlint-disable MD031 -->
<!-- TODO @randomvariable: Fix spacing to re-enable MD031 -->

# Delete Tanzu Kubernetes Grid Extensions

If you have deployed extensions that you no longer require, you can delete them from management clusters and Tanzu Kubernetes clusters.

## <a id="prepare"></a> Prepare to Delete Extensions

1. In a terminal, navigate to the folder that contains the unpacked Tanzu Kubernetes Grid extension manifest files, `tkg-extensions-v1.4.0+vmware.1
/extensions`.

   ```
   cd <path>/tkg-extensions-v1.4.0+vmware.1
/extensions
   ```

   Run all of the commands in these procedures from this location.
1. Set the context of `kubectl` to the management cluster or Tanzu Kubernetes cluster on which the extension is deployed.

    ```
    kubectl config use-context contour-test-admin@contour-test
    ```

**IMPORTANT**: For all of the extensions, do not delete `namespace-role.yaml` before the application has been fully deleted. This leads to errors due to the service account that is used by `kapp-controller` being deleted.

## <a id="contour"></a> Delete the Contour Extension

1. Delete the Contour extension.

    ```
    kubectl delete -f ingress/contour/contour-extension.yaml
    ```
1. Delete the Contour application.

    ```
    kubectl delete app contour -n tanzu-system-ingress
    ```
1. Delete the Contour namespace.

    ```
    kubectl delete -f ingress/contour/namespace-role.yaml
    ```

## <a id="fluentbit"></a> Delete the Fluent Bit Extension

1. Delete the Fluent Bit extension.

    ```
    kubectl delete -f logging/fluent-bit/fluent-bit-extension.yaml
    ```
1. Delete the Fluent Bit application.

    ```
    kubectl delete app fluent-bit -n tanzu-system-logging
    ```
1. Delete the Fluent Bit namespace.

   ```
   kubectl delete -f logging/fluent-bit/namespace-role.yaml
   ```

## <a id="observability"></a> Delete the Prometheus and Grafana Extensions

1. Delete the Prometheus extension.

    ```
    kubectl delete -f monitoring/prometheus/prometheus-extension.yaml
    ```
1. Delete the Prometheus application.
    ```
    kubectl delete app prometheus -n tanzu-system-monitoring
    ```

1. Delete the Prometheus namespace.

   ```
   kubectl delete -f monitoring/prometheus/namespace-role.yaml
   ```
1. Delete the Grafana extension.

    ```
    kubectl delete -f monitoring/grafana/grafana-extension.yaml
    ```
1. Delete the Grafana application.
    ```
    kubectl delete app grafana -n tanzu-system-monitoring
    ```
1. Delete the Grafana namespace.

   ```
   kubectl delete -f monitoring/grafana/namespace-role.yaml
   ```

## <a id="external-dns"></a> Delete the External DNS Extension

1. Delete the External DNS extension.

    ```
    kubectl delete -f registry/dns/external-dns-extension.yaml
    ```
1. Delete the External DNS application.
    ```
    kubectl delete app external-dns -n tanzu-system-registry
    ```
1. Delete the External DNS namespace.

   ```
   kubectl delete -f dns/external-dns/namespace-role.yaml
   ```

## <a id="harbor"></a> Delete the Harbor Extension

1. Delete the Harbor extension.

    ```
    kubectl delete -f registry/harbor/harbor-extension.yaml
    ```
1. Delete the Harbor application.
    ```
    kubectl delete app harbor -n tanzu-system-registry
    ```
1. Delete the Harbor namespace.

   ```
   kubectl delete -f registry/harbor/namespace-role.yaml
   ```

## <a id="authentication"></a> Delete the Dex and Gangway Extensions

1. Delete the Dex extension.

    ```
    kubectl delete -f authentication/dex/dex-extension.yaml
    ```
1. Delete the Dex application.
    ```
    kubectl delete app dex -n tanzu-system-auth
    ```
1. Delete the Dex namespace.

   ```
   kubectl delete -f authentication/dex/namespace-role.yaml
   ```
1. Delete the Gangway extension.

    ```
    kubectl delete -f authentication/gangway/gangway-extension.yaml
    ```
1. Delete the Gangway application.
    ```
    kubectl delete app gangway -n tanzu-system-auth
    ```
1. Delete the Gangway namespace.

   ```
   kubectl delete -f authentication/gangway/namespace-role.yaml
   ```

## <a id="utils"></a> Delete the Extensions Utilities

If you delete all extensions from a cluster, you can remove common extensions utilities.

If the extensions are deployed on a Tanzu Kubernetes cluster, optionally delete the `cert-manager`.

 ```
 kubectl delete -f ../cert-manager/
 ```

Do not delete `cert-manager` from management clusters.
