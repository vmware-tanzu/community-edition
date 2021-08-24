# Upgrade Tanzu Kubernetes Grid Extensions

This topic describes how to upgrade Tanzu Kubernetes Grid extensions. You upgrade the extensions after you upgrade to v1.3.x. Tanzu Kubernetes Grid extensions are deployed and managed by `kapp-controller` from the [Carvel Tools](https://carvel.dev/).

You should upgrade the following extensions on clusters that have been upgraded to Tanzu Kubernetes Grid v1.3.x:

- Contour
- Harbor
- Fluent Bit
- Prometheus
- Grafana
- External DNS

## <a id="upgrade-consider"></a> Considerations for Upgrading Extensions from v1.2.x to v1.3.x

This section lists some changes in Tanzu v1.3.x that impact the upgrade of extensions from v1.2.x to v1.3.x.

### <a id="dex-gang"></a> Dex and Gangway Extensions Upgrade

Tanzu Kubernetes Grid v1.3.x allows you to manage authentication with Pinniped instead of using Dex and Gangway.

Instead of upgrading the Dex and Gangway extensions, you should migrate your clusters to use Dex and Pinniped. After you migrate, you use the add-on manager to perform upgrades of Dex and Pinniped. For instructions on how to migrate your clusters, see [Register Core Add-ons](addons.md).

### <a id="registery"></a> Registry Update

Tanzu Kubernetes Grid v1.3.x switches the registry from `registry.tkg.vmware.run` to `projects.registry.vmware.com/tkg`.

To implement this registry change, you must apply the `cert-manager` included in Tanzu Kubernetes Grid v1.3.x on each cluster where you are upgrading Contour, Prometheus and Grafana extensions.

### <a id="tmc-ext-mgmr-removal"></a> Extension Manager Removal

Tanzu Kubernetes Grid v1.3.1 removes the Tanzu Mission Control extension manager from the extensions bundle. Extensions are no longer wrapped inside the extension resource. Instead of deploying extensions with Tanzu Mission Control extension manager, you deploy them by using the Kapp controller. As part of the upgrade, you must remove the extension resource for each extension being upgraded on the cluster.

## <a id="prereq"></a> Prerequisites

This procedure assumes that you are upgrading to Tanzu Kubernetes Grid v1.3.1 from either v1.2.x or v1.3.0.

To upgrade Tanzu Kubernetes Grid extensions to v1.3.1:

- You previously deployed one or more of the extensions on clusters running Tanzu Kubernetes Grid v1.2.x or v1.3.0.
- You have upgraded the management clusters to Tanzu Kubernetes Grid v1.3.1 or later.
- You have upgraded the clusters on which the extensions are running to Tanzu Kubernetes Grid v1.3.1 or later.
- You have installed the Carvel tools. For information about installing the Carvel tools, see [Install the Carvel Tools](../install-cli.md#install-carvel).
- You have downloaded and unpacked the bundle of Tanzu Kubernetes Grid extensions for v1.3.1 to the `tkg-extensions-v1.4.0+vmware.1
/extensions` folder. For information about where to obtain the bundle, see [Download and Unpack the Tanzu Kubernetes Grid Extensions Bundle](../extensions/index.md#unpack-bundle).
- Read the [Tanzu Kubernetes Grid 1.3.1 Release Notes](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3.1/rn/VMware-Tanzu-Kubernetes-Grid-131-Release-Notes.html) for updates related to security patches.

## <a id="contour-upgrade"></a> Upgrade the Contour Extension

Follow these steps to upgrade the Contour extension.
  
1. In a terminal, set the context of `kubectl` to the Tanzu Kubernetes cluster where you want to upgrade Contour.

1. Delete the Contour extension resource.

   ```sh
   kubectl delete extension contour -n tanzu-system-ingress
   ```

1. Delete the extension manager unless the cluster is attached in Tanzu Mission Control. If the cluster is attached in Tanzu Mission Control, you can skip this entire step.

    - Change directory to the extension bundle for the current version of the extension.
       - **v1.3.0** `tkg-extensions-v1.3.0+vmware.1/extensions`
       - **v1.2.1** `tkg-extensions-v1.2.1+vmware.1/extensions`
       - **v1.2.0** `tkg-extensions-v1.2.0+vmware.1/extensions`

    - Remove the extension manager.

      ```sh
      kubectl delete -f tmc-extension-manager.yaml
      ```

1. Navigate to the `tkg-extensions-v1.4.0+vmware.1
` folder where you downloaded the new extensions bundle.

1. Change directories to the `extensions` subfolder.

   ```sh
   cd extensions
   ```

1. Update the Kapp controller.

   ```sh
   kubectl apply -f kapp-controller.yaml
   ```

1. Update `cert-manager` to switch to the new registry.

   ```sh
   kubectl apply -f ../cert-manager/
   ```
  
1. Change directories to `ingress/contour`.

1. Obtain the current `contour-data-values.yaml` and secret used by the current Contour Extension.

   ```sh
   kubectl get secret contour-data-values -n tanzu-system-ingress -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > current-contour-data-values.yaml
   ```

1. Copy the `contour-data-values.yaml.example` file for your provider and name it `contour-data-values.yaml`.

   **vSphere:**

   ```
   cp vsphere/contour-data-values.yaml.example vsphere/contour-data-values.yaml
   ```

   **Amazon EC2**:

   ```
   cp aws/contour-data-values.yaml.example aws/contour-data-values.yaml
   ```

   **Azure**:

   ```
   cp azure/contour-data-values.yaml.example azure/contour-data-values.yaml
   ```

1. Manually copy over any customizations in `current-contour-data-values.yaml` into the `contour-data-values.yaml` file for your provider. For example, you may have customized `NodePort`.

   Do not change this configuration data. Otherwise, the upgrade will fail.

1. Update the `contour-data-values` secret with the new `contour-data-values.yaml`.

   **vSphere:**

   ```
   kubectl create secret generic contour-data-values --from-file=values.yaml=vsphere/contour-data-values.yaml -n tanzu-system-ingress -o yaml --dry-run | kubectl replace -f-
   ```

   **Amazon EC2**:

   ```
   kubectl create secret generic contour-data-values --from-file=values.yaml=aws/contour-data-values.yaml -n tanzu-system-ingress -o yaml --dry-run | kubectl replace -f-
   ```

   **Azure**:

   ```
   kubectl create secret generic contour-data-values --from-file=values.yaml=azure/contour-data-values.yaml -n tanzu-system-ingress -o yaml --dry-run | kubectl replace -f-
   ```
  
1. Deploy the new Contour extension.

   ```
   kubectl apply -f contour-extension.yaml
   ```

   You should see the confirmation
   `extension.clusters.tmc.cloud.vmware.com/contour created`.

1. View the status of the Contour service itself.

    ```
    kubectl get app contour -n tanzu-system-ingress
    ```

    The status of the Contour app should show `Reconcile Succeeded` when Contour has deployed successfully.

    ```
    NAME      DESCRIPTION           SINCE-DEPLOY   AGE
    contour   Reconcile succeeded   15s            72s
    ```

1. Check that the ingress and HTTP proxy resources are valid and that the ingress traffic is working.

    ```sh
    kubectl get ingress -A
    kubectl get httpproxy -A
    ```

## <a id="harbor-upgrade"></a> Upgrade the Harbor Extension

Follow these steps to upgrade the Harbor extension.

1. In a terminal, set the context of `kubectl` to the Tanzu Kubernetes cluster where you want to upgrade Harbor.

1. Delete the Harbor extension resource.

   ```sh
   kubectl delete extension harbor -n tanzu-system-registry
   ```

1. Delete the extension manager unless the cluster is attached in Tanzu Mission Control. If the cluster is attached in Tanzu Mission Control, you can skip this entire step.

    - Change directory to the extension bundle for the current version of the extension.
       - **v1.3.0** `tkg-extensions-v1.3.0+vmware.1/extensions`
       - **v1.2.1** `tkg-extensions-v1.2.1+vmware.1/extensions`
       - **v1.2.0** `tkg-extensions-v1.2.0+vmware.1/extensions`

    - Remove the extension manager.

      ```sh
      kubectl delete -f tmc-extension-manager.yaml
      ```

1. Navigate to the `tkg-extensions-v1.4.0+vmware.1
` folder where you downloaded the new extensions bundle.

1. Switch `kubectl` config context to the shared services cluster.

1. If you have not done so already, upgrade the Contour Extension on the shared services cluster as described in the previous procedure, [Upgrade the Contour Extension](#contour-upgrade). Contour is a dependency for Harbor.

1. Change directories to the `extensions` subfolder.

   ```sh
   cd extensions
   ```

1. Update the Kapp controller.

   ```sh
   kubectl apply -f kapp-controller.yaml
   ```

1. Change directories to `registry/harbor`.

1. Obtain the `harbor-data-values.yaml` and secret used by the current Harbor Extension.

   ```sh
   kubectl get secret harbor-data-values -n tanzu-system-registry -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > current-harbor-data-values.yaml
   ```

1. Copy `harbor-data-values.yaml.example` to `harbor-data-values.yaml`.

   ```sh
   cp harbor-data-values.yaml.example harbor-data-values.yaml
   ```

1. Manually copy all custom configuration data in `current-harbor-data-values.yaml` into the `harbor-data-values.yaml` file. You may need to copy in individual values such as `hostname`, `harborAdminPassword`, `secretKey`, `core.secret`, `persistence`, and so on.

   Do not change this configuration data. Otherwise, the upgrade will fail.

1. Update the `harbor-data-values` secret with the new `harbor-data-values.yaml`.

   ```sh
   kubectl create secret generic harbor-data-values --from-file=values.yaml=harbor-data-values.yaml -n tanzu-system-registry -o yaml --dry-run | kubectl replace -f-
   ```

1. Deploy the new Harbor Extension.

   ```sh
   kubectl apply -f harbor-extension.yaml
   ```

1. Retrieve the status of Harbor Extension.

   ```sh
   kubectl get app harbor -n tanzu-system-registry
   ```

   The Harbor App status should change to `Reconcile succeeded` after the Harbor Extension is deployed successfully.

   View detailed status:

   ```sh
   kubectl get app harbor -n tanzu-system-registry -o yaml
   ```

1. In a web browser, navigate to the current Harbor portal URL.

1. Ensure that you can sign in with the current username and password and that all your projects and images are present.

After the upgrade completes, all commands such as `docker login`, `docker pull`, `docker push` should be fully functional since the Harbor CA certificate does not change during the upgrade.

## <a id="fluent-bit-upgrade"></a> Upgrade the Fluent Bit Extension

Follow these steps to upgrade the Fluent Bit Extension.

**Note** The Fluent Bit extension in Tanzu Kubernetes Grid v1.3.x adds support for **Syslog**. To install and configure **Syslog** in Fluent Bit, see [Prepare the Fluent Bit Configuration File for a Syslog Output Plugin](../extensions/logging-fluentbit.md#syslog).

1. In a terminal, set the context of `kubectl` to the Tanzu Kubernetes cluster where you want to upgrade Fluent Bit.

1. Delete the Fluent Bit extension resource.

   ```sh
   kubectl delete extension fluent-bit -n tanzu-system-logging
   ```

1. Delete the extension manager unless the cluster is attached in Tanzu Mission Control. If the cluster is attached in Tanzu Mission Control, you can skip this entire step.

    - Change directory to the extension bundle for the current version of the extension.
       - **v1.3.0** `tkg-extensions-v1.3.0+vmware.1/extensions`
       - **v1.2.1** `tkg-extensions-v1.2.1+vmware.1/extensions`
       - **v1.2.0** `tkg-extensions-v1.2.0+vmware.1/extensions`

    - Remove the extension manager.

      ```sh
      kubectl delete -f tmc-extension-manager.yaml
      ```

1. Navigate to the `tkg-extensions-v1.4.0+vmware.1
` folder where you downloaded the new extensions bundle.

1. Change directories to the `extensions` subfolder.

   ```sh
   cd extensions
   ```

1. Update the Kapp controller.

   ```sh
   kubectl apply -f kapp-controller.yaml
   ```

1. Set the context of `kubectl` to the Tanzu Kubernetes cluster where you want to upgrade Fluent Bit.

1. Change directories to `logging/fluent-bit`.

1. Obtain the current `fluent-bit-data-values.yaml` files and secrets used by the current Fluent Bit Extension.

   **Elastic Search**

   ```sh
   kubectl get secret fluent-bit-data-values -n tanzu-system-logging -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > elasticsearch/current-fluent-bit-data-values.yaml
   ```

   **Kafka**

   ```sh
   kubectl get secret fluent-bit-data-values -n tanzu-system-logging -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > kafka/current-fluent-bit-data-values.yaml
   ```

   **Splunk**

   ```sh
   kubectl get secret fluent-bit-data-values -n tanzu-system-logging -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > splunk/current-fluent-bit-data-values.yaml
   ```

   **HTTP**

   ```sh
   kubectl get secret fluent-bit-data-values -n tanzu-system-logging -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > http/current-fluent-bit-data-values.yaml
   ```

   **syslog** (Only applicable if you are upgrading from v1.3.0 or later)

   ```sh
   kubectl get secret fluent-bit-data-values -n tanzu-system-logging -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > syslog/current-fluent-bit-data-values.yaml
   ```

1. Copy the `fluent-bit-data-values.yaml.example` file for each backend logging component and name it `fluent-bit-data-values.yaml`.

   **Elastic Search**

   ```sh
   cp elasticsearch/fluent-bit-data-values.yaml.example elasticsearch/fluent-bit-data-values.yaml
   ```

   **Kafka**

   ```sh
   cp kafka/fluent-bit-data-values.yaml.example kafka/fluent-bit-data-values.yaml
   ```

   **Splunk**

   ```sh
   cp splunk/fluent-bit-data-values.yaml.example splunk/fluent-bit-data-values.yaml
   ```

   **HTTP**

   ```sh
   cp http/fluent-bit-data-values.yaml.example http/fluent-bit-data-values.yaml
   ```

   **syslog** (Only applicable if you are upgrading from v1.3.0 or later)

   ```sh
   cp syslog/fluent-bit-data-values.yaml.example syslog/fluent-bit-data-values.yaml
   ```

1. For each backend component, manually copy over any customizations in your `current-fluent-bit-data-values.yaml` files into the `fluent-bit-data-values.yaml` file for the component.

1. Update the `fluent-bit-data-values` secrets with the new `fluent-bit-data-values.yaml` file.

   **ElasticSearch:**

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=elasticsearch/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```

   **Kafka**:

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=kafka/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```

   **Splunk**:

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=splunk/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```

   **HTTP**:

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=http/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```
  
1. Deploy the Fluent Bit extension.

   ```sh
   kubectl apply -f fluent-bit-extension.yaml
   ```

   You should see the confirmation `extension.clusters.tmc.cloud.vmware.com/fluent-bit created`.

1. View the status of the Fluent Bit service itself.

    ```sh
    kubectl get app fluent-bit -n tanzu-system-logging
    ```

    The status of the Fluent Bit app should show `Reconcile Succeeded` when Fluent Bit has deployed successfully.

    ```sh
    NAME         DESCRIPTION           SINCE-DEPLOY   AGE
    fluent-bit   Reconcile succeeded   54s            14m
    ```

1. Check that the Fluent Bit daemon set is running and that log collection and forwarding is functioning.

   ```sh
   kubectl get ds -n tanzu-system-logging
   ```

## <a id="prom_upgrade"></a> Upgrade the Prometheus Extension

Follow these steps to upgrade the Prometheus Extension.

1. In a terminal, set the context of `kubectl` to the Tanzu Kubernetes cluster where you want to upgrade Prometheus.

1. Delete the Prometheus extension resource.

   ```sh
   kubectl delete extension prometheus -n tanzu-system-monitoring

   ```

1. Delete the extension manager unless the cluster is attached in Tanzu Mission Control. If the cluster is attached in Tanzu Mission Control, you can skip this entire step.

    - Change directory to the extension bundle for the current version of the extension.
       - **v1.3.0** `tkg-extensions-v1.3.0+vmware.1/extensions`
       - **v1.2.1** `tkg-extensions-v1.2.1+vmware.1/extensions`
       - **v1.2.0** `tkg-extensions-v1.2.0+vmware.1/extensions`

    - Remove the extension manager.

      ```sh
      kubectl delete -f tmc-extension-manager.yaml
      ```

1. Navigate to the `tkg-extensions-v1.4.0+vmware.1
` folder where you downloaded the new extensions bundle.

1. Change directories to the `extensions` subfolder.

   ```sh
   cd extensions
   ```

1. Update the Kapp controller.

   ```sh
   kubectl apply -f kapp-controller.yaml
   ```

1. Update `cert-manager` to switch to the new registry.

   ```sh
   kubectl apply -f ../cert-manager/
   ```

1. Change directories to `monitoring/prometheus`.

1. Obtain the current `prometheus-data-values.yaml` files and secrets used by the current Prometheus Extension.

   ```sh
   kubectl get secret prometheus-data-values -n tanzu-system-monitoring -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > current-prometheus-data-values.yaml
   ```

1. Copy the `prometheus-data-values.yaml.example` file to `prometheus-data-values.yaml`.

   ```sh
   cp prometheus-data-values.yaml.example prometheus-data-values.yaml
   ```

1. Manually copy over any customizations in your `current-prometheus-data-values.yaml` files into the `prometheus-data-values.yaml` file.

1. Update the `prometheus-data-values` secret with the new `prometheus-data-values.yaml`.

   ```sh
   kubectl create secret generic prometheus-data-values --from-file=values.yaml=prometheus-data-values.yaml -n tanzu-system-monitoring -o yaml --dry-run | kubectl replace -f-
   ```

1. Deploy the Prometheus extension.

    ```sh
    kubectl apply -f prometheus-extension.yaml
    ```

    You should see a confirmation that `extensions.clusters.tmc.cloud.vmware.com/prometheus` has been created.

2. The extension takes several minutes to deploy. To check the status of the deployment, use the `kubectl get extension` and the `kubectl get app` commands.

    ```sh
    kubectl get app prometheus -n tanzu-system-monitoring
    ```

    While the extension is being deployed, the "Description" field from the `kubectl get app` command shows a status of `Reconciling`. After Prometheus is deployed successfully, the status of the Prometheus app shown by the `kubectl get app` command changes to `Reconcile succeeded`.

    You can view detailed status information with this command:

    ```sh
    kubectl get app prometheus -n tanzu-system-monitoring -o yaml
    ```

## <a id="upgrade-grafana"></a> Upgrade the Grafana Extension

Follow these steps to upgrade the Grafana extension.

1. In a terminal, set the context of `kubectl` to the Tanzu Kubernetes cluster where you want to upgrade Grafana.

1. Delete the Grafana extension resource.

   ```sh
   kubectl delete extension grafana -n tanzu-system-monitoring

   ```

1. Delete the extension manager unless the cluster is attached in Tanzu Mission Control. If the cluster is attached in Tanzu Mission Control, you can skip this entire step.

    - Change directory to the extension bundle for the current version of the extension.
       - **v1.3.0** `tkg-extensions-v1.3.0+vmware.1/extensions`
       - **v1.2.1** `tkg-extensions-v1.2.1+vmware.1/extensions`
       - **v1.2.0** `tkg-extensions-v1.2.0+vmware.1/extensions`

    - Remove the extension manager.

      ```sh
      kubectl delete -f tmc-extension-manager.yaml
      ```

1. Navigate to the `tkg-extensions-v1.4.0+vmware.1
` folder where you downloaded the new extensions bundle.

1. Change directories to the `extensions` subfolder.

   ```sh
   cd extensions
   ```

1. Update the Kapp controller.

   ```sh
   kubectl apply -f kapp-controller.yaml
   ```

1. Update `cert-manager` to switch to the new registry.

   ```sh
   kubectl apply -f ../cert-manager/
   ```

1. If you have not done so already, upgrade the Contour Extension as described in the previous procedure, [Upgrade the Contour Extension](#contour-upgrade). Contour is a dependency for Grafana.

1. Change directories to `monitoring/grafana`.

1. Obtain the current `grafana-data-values.yaml` files and secrets used by the current Grafana Extension.

   ```sh
   kubectl get secret grafana-data-values -n tanzu-system-monitoring -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > current-grafana-data-values.yaml
   ```

1. Copy the `grafana-data-values.yaml.example` file to `grafana-data-values.yaml`.

   ```sh
   cp grafana-data-values.yaml.example grafana-data-values.yaml
   ```

1. Manually copy over any customizations in your `current-grafana-data-values.yaml` files into the `grafana-data-values.yaml` file.

1. Update the `grafana-data-values` secret with the new `grafana-data-values.yaml` file.

   ```sh
   kubectl create secret generic grafana-data-values --from-file=values.yaml=grafana-data-values.yaml -n tanzu-system-monitoring -o yaml --dry-run | kubectl replace -f-
   ```

1. Deploy the new Grafana extension.

    ```
    kubectl apply -f grafana-extension.yaml
    ```

    You should see a confirmation that `extensions.clusters.tmc.cloud.vmware.com/grafana` was created.

1. The extension takes several minutes to deploy. To check the status of the deployment, use the `kubectl get app -n tanzu-system-monitoring` command:

    ```
    kubectl get app grafana -n tanzu-system-monitoring
    ```

    While the extension is being deployed, the "Description" field from the `kubectl get app` command shows a status of `Reconciling`. After Grafana is deployed successfully, the status of the Grafana app as shown by the `kubectl get app` command changes to `Reconcile succeeded`.

    You can view detailed status information with this command:

    ```
    kubectl get app grafana -n tanzu-system-monitoring -o yaml
    ```

## <a id="upgrade-external-dns"></a> Upgrade the External DNS Extension

External DNS is available as an extension starting in Tanzu Kubernentes Grid v1.3.0 and later. This upgrade procedure only applies if you already installed External DNS in v1.3.0 and are upgrading the extension to v1.3.1 or later.

If you want to install the External DNS extension, see [Implementing Service Discovery with External DNS](../extensions/external-dns.md).

Follow these steps to upgrade the External DNS extension.

1. In a terminal, set the context of `kubectl` to the Tanzu Kubernetes cluster where you want to upgrade External DNS.

1. Delete the External DNS extension resource.

   ```sh
   kubectl delete extension external-dns -n tanzu-system-service-discovery

   ```

1. Delete the extension manager unless the cluster is attached in Tanzu Mission Control. If the cluster is attached in Tanzu Mission Control, you can skip this entire step.

    - Change directory to the extension bundle for the current version of the extension.
       - **v1.3.0** `tkg-extensions-v1.3.0+vmware.1/extensions`

    - Remove the extension manager.

      ```sh
      kubectl delete -f tmc-extension-manager.yaml
      ```

1. Navigate to the `tkg-extensions-v1.4.0+vmware.1
` folder where you downloaded the new extensions bundle.

1. Switch `kubectl` config context to the shared services cluster.

1. If you have not done so already, upgrade the Contour and Harbor extensions on the shared services cluster as described in the previous procedure, [Upgrade the Contour Extension](#contour-upgrade) and [Upgrade the Harbor Extension](#harbor-upgrade).

1. Change directories to the `extensions` subfolder.

   ```sh
   cd extensions
   ```

1. Update the Kapp controller.

   ```sh
   kubectl apply -f kapp-controller.yaml
   ```

1. Obtain the current `external-dns-data-values.yaml` files and secrets used by the current External DNS extension.

   ```sh
   kubectl get secret external-dns-data-values -n tanzu-system-service-discovery -o 'go-template={{ index .data "values.yaml" }}' | base64 -d > current-external-dns-data-values.yaml
   ```

1. Copy the `external-dns-data-values-PROVIDER.yaml.example` file where `PROVIDER` matches your DNS service provider to `external-dns-data-values.yaml`. For example:

   ```sh
   cp external-dns-values-azure.yaml.example external-dns-data-values.yaml
   ```

1. Manually copy over any customizations in your `current-external-dns-data-values.yaml` files into the `external-dns-values.yaml` file.

1. Update the `external-dns-values` secret with the new `external-dns-values.yaml` file.

   ```sh
   kubectl create secret generic external-dns-values --from-file=values.yaml=external-dns-values.yaml -n tanzu-system-service-discovery -o yaml --dry-run | kubectl replace -f-
   ```

1. Deploy the new External DNS extension.

    ```
    kubectl apply -f external-dns-extension.yaml
    ```

1. The extension takes several minutes to deploy. To check the status of the deployment, use the `kubectl get app -n tanzu-system-service-discovery` command:

    ```
    kubectl get app external-dns -n tanzu-system-service-discovery
    ```

    While the extension is being deployed, the "Description" field from the `kubectl get app` command shows a status of `Reconciling`. After External DNS is deployed successfully, the status of the External DNS app as shown by the `kubectl get app` command changes to `Reconcile succeeded`.

    You can view detailed status information with this command:

    ```
    kubectl get app external-dns -n tanzu-system-service-discovery -o yaml
    ```
