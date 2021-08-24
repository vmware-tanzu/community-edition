<!-- markdownlint-disable MD031 -->
<!-- TODO @randomvariable: Fix spacing to re-enable MD031 -->
# Deploy Grafana on Tanzu Kubernetes Clusters

Tanzu Kubernetes Grid includes a Grafana extension that you can deploy on your Tanzu Kubernetes clusters.
Grafana allows you to visualize and analyze metrics data collected by Prometheus on your clusters.

## <a id="prereqs"></a> Prerequisites

Before you begin this procedure, you must have:

- Downloaded and unpacked the bundle of Tanzu Kubernetes Grid extensions. For information about where to obtain the bundle, see [Download and Unpack the Tanzu Kubernetes Grid Extensions Bundle](index.md#unpack-bundle).
- Installed the Carvel tools. For information about installing the Carvel tools, see [Install the Carvel Tools](../install-cli.md#install-carvel).
- Deployed a Tanzu Kubernetes Grid management cluster on vSphere, Amazon EC2, or Azure.
- Deployed a Tanzu Kubernetes cluster. The examples in this topic use a cluster named `monitoring-cluster`.
- Installed the Prometheus extension on the Tanzu Kubernetes cluster. For instructions on how to install Prometheus, see [Deploy Prometheus on Tanzu Kubernetes Clusters](prometheus.md).
- Installed Contour for ingress control on the Tanzu Kubernetes cluster. For information on installing Contour, see [Implementing Ingress Control with Contour](ingress-contour.md).

**IMPORTANT**:

- Tanzu Kubernetes Grid does not support IPv6 addresses because upstream Kubernetes only provides alpha support for IPv6. In the following procedures, you must always provide IPv4 addresses.
- The extensions  folder `tkg-extensions-v1.4.0+vmware.1
` contains subfolders for each type of extension, for example, `authentication`, `ingress`, `registry`, and so on. At the top level of the folder there is an additional subfolder named `extensions`. The `extensions` folder also contains subfolders for `authentication`, `ingress`, `registry`, and so on. Take care to run commands from the location provided in the instructions. Commands are usually run from within the `extensions` folder.

## <a id="prepare-tkc"></a> Prepare the Tanzu Kubernetes Cluster for the Grafana Extension

To deploy the Grafana extension, you must prepare the specific Tanzu Kubernetes cluster where you plan to deploy the extension. First you must install a few supporting applications on the Tanzu Kubernetes cluster.

If you have already installed another extension such as the Prometheus extension onto the Tanzu Kubernetes cluster, then you can skip this section and proceed directly to [Prepare the Configuration File for the Grafana Extension](#config).

This procedure applies to Tanzu Kubernetes clusters running on vSphere, Amazon EC2, and Azure.

1. In a terminal, navigate to the folder that contains the unpacked Tanzu Kubernetes Grid extension manifest files,
`tkg-extensions-v1.4.0+vmware.1
/extensions`.

   ```
   cd <path>/tkg-extensions-v1.4.0+vmware.1
/extensions
   ```

   You should see subfolders named `authentication`, `ingress`, `logging`, `monitoring`, `registry`, `service-discovery` and some YAML files.

1. Retrieve the `admin` credentials of the Tanzu Kubernetes cluster.

    ```sh
    tanzu cluster kubeconfig get monitoring-cluster --admin
    ```

1. Set the context of `kubectl` to the Tanzu Kubernetes cluster.

   ```sh
   kubectl config use-context monitoring-cluster-admin@monitoring-cluster
   ```

1. If you haven't already, install `cert-manager` on the Tanzu Kubernetes workload cluster by following the procedure in [Install Cert Manager on Workload Clusters](./index.md#cert-mgr
).

When all pods are ready, the Tanzu Kubernetes cluster is ready for you to deploy the Prometheus extension. To do so, follow the procedure in [Prepare the Configuration Files for the Grafana Extension](#config).

## <a id="config"></a> Prepare the Grafana Extension Configuration File

This procedure describes how to prepare the Grafana extenson configuration file for Tanzu Kubernetes clusters.  This configuration file applies to Tanzu Kubernetes clusters running on vSphere, Amazon EC2, and Azure and is required to deploy the Grafana extension.

For additional configuration options, you can use `ytt` overlays as described in [Extensions and Shared Services](../ytt.md#extensions) in _Customizing Clusters, Plans, and Extensions with ytt Overlays_ and in the extensions mods examples in the [TKG Lab repository](https://github.com/Tanzu-Solutions-Engineering/tkg-lab).

1. Make a copy of the `grafana-data-values.yaml.example` file for your infrastructure platform, and name the file `grafana-data-values.yaml`.

    ```sh
    cp monitoring/grafana/grafana-data-values.yaml.example monitoring/grafana/grafana-data-values.yaml
    ```

1. Edit the `grafana-data-values.yaml` file and replace `<ADMIN_PASSWORD>` with a Base64 encoded password.

   To generate a Base64 encoded password, run the following command:

   ```
   echo -n 'mypassword' | base64
   ```

   You can also use the Base64 encoding tool at [https://www.base64encode.org/](https://www.base64encode.org/) to encode your password. For example, by using either method, a password of `mypassword` results in the encoded password `bXlwYXNzd29yZA==`.

1. Save `grafana-data-values.yaml` when you are finished.

1. You can now either deploy the Grafana extension with default values, or you can you can customize the deployment.

    - To deploy Grafana using default configuration values, proceed directly to [Deploy Grafana on a Tanzu Kubernetes Cluster](#deploy). By default, `grafana-data-values.yaml` only contains the configuration of the infrastructure provider and a default administrative password.

    - To customize your Grafana deployment, see [Customize Your Grafana Deployment](#customize). For example, you can specify LDAP authentication or configure storage for Grafana.

## <a id="customize"></a> Customize the Configuration of the Grafana Extension

You can customize the configuration of the Grafana extension by editing the `tkg-extensions-v1.4.0+vmware.1
/extensions/monitoring/grafana/grafana-data-values.yaml` file.

If you modify this file _before_ you deploy the Grafana extension, then the custom settings take effect immediately upon deployment. For instructions, see [Deploy Grafana on a Tanzu Kubernetes Cluster](#deploy).

If you modify this file _after_ you deploy the Grafana extension, then you must update your running deployment. For instructions, see [Update a Running Grafana Extension](#update).

<a id="config-table"></a> **Grafana Extension Configuration Parameters**

The following table describes configuration parameters of the Grafana extension and their default values. To customize Grafana, specify the parameters and their custom values in the `grafana-data-values.yaml` file of your Tanzu Kubernetes cluster.

| Parameter                                            | Type and Description                                                                                                                     | Default                                        |
| ---------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| `monitoring.namespace`                               | String. Namespace in which to deploy Grafana.                                                                                            | `tanzu-system-monitoring`                      |
| `monitoring.create_namespace`                        | Boolean. The flag indicates whether to create the namespace specified by `monitoring.namespace`.                                         | `false`                                        |
| `monitoring.grafana.cluster_role.apiGroups`          | List. API group defined for Grafana `ClusterRole`.                                                                                       | `[""]`                                         |
| `monitoring.grafana.cluster_role.resources`          | List. Resources defined for Grafana `ClusterRole`.                                                                                       | `["configmaps", "secrets"]`                    |
| `monitoring.grafana.cluster_role.verbs`              | List. Access permission defined for `ClusterRole`.                                                                                       | `["get", "watch", "list"]`                     |
| `monitoring.grafana.config.grafana_ini`              | Config file. Grafana configuration file details.                                                                                         | `grafana.ini`                                  |
| `monitoring.grafana.config.datasource.type`          | String. Grafana datasource type.                                                                                                         | `prometheus`                                   |
| `monitoring.grafana.config.datasource.access`        | String. Access mode, proxy or direct (Server or Browser in the UI).                                                                      | `proxy`                                        |
| `monitoring.grafana.config.datasource.isDefault`     | Boolean. Flag to mark the default Grafana datasource.                                                                                    | `true`                                         |
| `monitoring.grafana.config.provider_yaml`            | YAML file. Config file to define Grafana dashboard provider.                                                                             | `provider.yaml`                                |
| `monitoring.grafana.service.type`                    | String. Type of Kubernetes Service to expose Grafana: `ClusterIP`, `NodePort`, `LoadBalancer`.                                           | vSphere: `NodePort`, AWS/Azure: `LoadBalancer` |
| `monitoring.grafana.pvc.storage_class`               | String. `StorageClass` to use for Persistent Volume Claim. By default this is null and default provisioner is used.                      | `null`                                         |
| `monitoring.grafana.pvc.accessMode`                  | String. Define access mode for Persistent Volume Claim: `ReadWriteOnce`, `ReadOnlyMany`, `ReadWriteMany`.                                | `ReadWriteOnce`                                |
| `monitoring.grafana.pvc.storage`                     | String. Define storage size for Persistent Volume Claim.                                                                                 | 2Gi                                            |
| `monitoring.grafana.deployment.replicas`             | Integer. Number of Grafana replicas.                                                                                                     | 1                                              |
| `monitoring.grafana.image.repository`                | String. Repository containing the Grafana image.                                                                                         | `projects.registry.vmware.com/tkg/grafana`     |
| `monitoring.grafana.image.name`                      | String. Name of the Grafana image.                                                                                                       | `grafana`                                      |
| `monitoring.grafana.image.tag`                       | String. Image tag of the Grafana image.                                                                                                  | `v7.0.3_vmware.1`                              |
| `monitoring.grafana.image.pullPolicy`                | String. Image pull policy for the Grafana image.                                                                                         | `IfNotPresent`                                 |
| `monitoring.grafana.secret.type`                     | String. Secret type defined for Grafana dashboard.                                                                                       | `Opaque`                                       |
| `monitoring.grafana.secret.admin_user`               | String. Username to access the Grafana dashboard.                                                                                        | `YWRtaW4=` (`admin` in Base64 encoding)        |
| `monitoring.grafana.secret.admin_password`           | String. Password to access Grafana dashboard.                                                                                            | null                                           |
| `monitoring.grafana.secret.ldap_toml`                | String. If using LDAP authentication, LDAP configuration file path.                                                                      | `""`                                           |
| `monitoring.grafana_init_container.image.repository` | String. Repository containing the Grafana init container image.                                                                          | `projects.registry.vmware.com/tkg/grafana`     |
| `monitoring.grafana_init_container.image.name`       | String. Name of the Grafana `init` container image.                                                                                      | `k8s-sidecar`                                  |
| `monitoring.grafana_init_container.image.tag`        | String. Image tag of the Grafana `init` container image.                                                                                 | `0.1.99`                                       |
| `monitoring.grafana_init_container.image.pullPolicy` | String. Image pull policy for the Grafana init container image.                                                                          | `IfNotPresent`                                 |
| `monitoring.grafana_sc_dashboard.image.repository`   | String. Repository containing the Grafana dashboard image.                                                                               | `projects.registry.vmware.com/tkg/grafana`     |
| `monitoring.grafana_sc_dashboard.image.name`         | String. Name of the Grafana dashboard image.                                                                                             | `k8s-sidecar`                                  |
| `monitoring.grafana_sc_dashboard.image.tag`          | String. Image tag of the Grafana dashboard image.                                                                                        | `0.1.99`                                       |
| `monitoring.grafana_sc_dashboard.image.pullPolicy`   | String. Image pull policy for the Grafana dashboard image.                                                                               | `IfNotPresent`                                 |
| `monitoring.grafana.ingress.enabled`                 | Boolean. Enable/disable ingress for Grafana.                                                                                             | `true`                                         |
| `monitoring.grafana.ingress.virtual_host_fqdn`       | String. Hostname for accessing Grafana.                                                                                                  | `grafana.system.tanzu`                         |
| `monitoring.grafana.ingress.prefix`                  | String. Path prefix for grafana.                                                                                                         | `/`                                            |
| `monitoring.grafana.ingress.tlsCertificate.tls.crt`  | String. Optional certificate for ingress if you want to use your own TLS certificate; a self-signed certificate is generated by default. | Generated certificate                          |
| `monitoring.grafana.ingress.tlsCertificate.tls.key`  | String. Optional certificate private key for ingress if you want to use your own TLS certificate.                                        | Generated certificate private key              |

## <a id="deploy"></a> Deploy Grafana on a Tanzu Kubernetes Cluster

After you have prepared a Tanzu Kubernetes cluster, you can deploy the Grafana extension on the cluster. As part of the preparation, you have updated the appropriate configuration file for your platform and optionally customized your deployment.

This procedure applies to Tanzu Kubernetes clusters running on vSphere, Amazon EC2, and Azure.

1. Create the namespace and RBAC roles for Grafana.

    ```sh
    kubectl apply -f extensions/monitoring/grafana/namespace-role.yaml
    ```

    You should see confirmation that a `tanzu-system-monitoring` namespace, service account, and RBAC role bindings are created for Grafana.

    ```
    namespace/tanzu-system-monitoring unchanged
    serviceaccount/grafana-extension-sa created
    role.rbac.authorization.k8s.io/grafana-extension-role created
    rolebinding.rbac.authorization.k8s.io/grafana-extension-rolebinding created
    clusterrole.rbac.authorization.k8s.io/grafana-extension-cluster-role created
    clusterrolebinding.rbac.authorization.k8s.io/grafana-extension-cluster-rolebinding created
    ```

    In this case, you may notice that the output states `namespace/tanzu-system-monitoring unchanged`. This output is an example of what you would see if you have already deployed the Prometheus extension, which is likely the case  since Grafana uses Prometheus as its datasource. If you installed Grafana first, then the output shows `namespace/tanzu-system-monitoring created` instead.

1. Create a Kubernetes secret that encodes the values stored in the `grafana-data-values.yaml` configuration file.

    ```sh
    kubectl -n tanzu-system-monitoring create secret generic grafana-data-values --from-file=values.yaml=extensions/monitoring/grafana/grafana-data-values.yaml
    ```

1. Deploy the Grafana extension.

    ```
    kubectl apply -f extensions/monitoring/grafana/grafana-extension.yaml
    ```

    You should see a confirmation that `extensions.clusters.tmc.cloud.vmware.com/grafana` was created.

1. The extension takes several minutes to deploy. To check the status of the deployment, use the `kubectl -n tanzu-system-monitoring get app` command:

    ```
    kubectl -n tanzu-system-monitoring get app grafana
    ```

    While the extension is being deployed, the "Description" field from the `kubectl get app` command shows a status of `Reconciling`. After Grafana is deployed successfully, the status of the Grafana app as shown by the `kubectl get app` command changes to `Reconcile succeeded`.

    You can view detailed status information with this command:

    ```
    kubectl get app grafana -n tanzu-system-monitoring -o yaml
    ```
1. After you have deployed Grafana, configure Contour. The Grafana extension requires Contour to be present and creates a Contour HTTPProxy object with an FQDN of `grafana.system.tanzu`.

1. To use this FQDN to access the Grafana dashboard, create an entry in your local `/etc/hosts` file that points an IP address this FQDN.

    **Amazon EC2** or **Azure**:
    Use the IP address of the LoadBalancer for the Envoy service
    in the `tanzu-system-ingress` namespace.

    **vSphere**: Use the IP address of a worker node.

1. Use a browser to navigate to `https://grafana.system.tanzu`.

   Since the site uses self-signed certificates, you may need to navigate through a browser-specific security warning before you are able to access the dashboard.

## <a id="update"></a> Update a Running Grafana Extension

If you need to make changes to the Grafana extension after it has been deployed, you must update the Kubernetes secret that the extension uses for its configuration. The steps below describe how to update the Kubernetes secret, and how to then update the configuration of the Grafana extension.

This procedure applies to Tanzu Kubernetes clusters running on vSphere, Amazon EC2, and Azure.

1. Locate the `grafana-data-values.yaml` file you created in the [Prepare the Configuration File for the Grafana Extension](#config). You must make your changes to this file. If you no longer have this file, you can recreate it with the following `kubectl` command:

    <pre>
    kubectl -n tanzu-system-monitoring get secret grafana-data-values -o 'go-template=&lbrace;&lbrace; index .data "values.yaml" &rbrace;&rbrace;' | base64 -d > grafana-data-values.yaml
    </pre>

    Note that macOS users will need to use the `-D` parameter to `base64`, instead of the lowercase `-d` shown above.

1. Using the information in [Customize the Configuration of the Grafana Extension](#customize) as a reference, make the necessary changes to the `grafana-data-values.yaml` file.

1. After you have made all applicable changes, save the file.

1. Update the Kubernetes secret.

   These command assume that you are running them from `tkg-extensions-v1.4.0+vmware.1
/extensions`.

   ```sh
   kubectl -n tanzu-system-monitoring create secret generic grafana-data-values --from-file=values.yaml=extensions/monitoring/grafana/grafana-data-values.yaml -o yaml --dry-run=client | kubectl replace -f -
    ```

   Note that the final `-` on the `kubectl replace` command above is necessary to instruct `kubectl` to accept the input being piped to it from the `kubectl create secret` command.

1. The Grafana extension is reconciled using the new values that you just added. The changes should show up in five minutes or less. Updates are handled by the Kapp controller, which synchronizes every five minutes.

    If you need the changes to the configuration reflected sooner, then you can delete the Grafana pod using `kubectl delete pod`. The pod is recreated via the Kapp controller with the new settings. You can also change `syncPeriod` in `grafana-extension.yaml` to a lower value and re-apply the configuration with `kubectl apply -f grafana-extension.yaml`.

## <a id="remove"></a> Remove the Grafana Extension

For information on how to remove the Grafana extension from a Tanzu Kubernetes cluster, see [Delete Tanzu Kubernetes Grid Extensions](delete-extensions#observability).
