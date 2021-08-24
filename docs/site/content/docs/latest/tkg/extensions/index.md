<!-- markdownlint-disable MD031 -->
<!-- TODO @randomvariable: Fix spacing to re-enable MD031 -->

# Deploying and Managing Extensions and Shared Services

Tanzu Kubernetes Grid includes binaries for tools that provide in-cluster and shared services to the clusters running in your Tanzu Kubernetes Grid instance. All of the provided binaries and container images are built and signed by VMware.

## <a id="locations"></a> Locations and Dependencies

You can add functionalities to Tanzu Kubernetes clusters by installing extensions to different cluster locations as follows:

<table id='TK-NAME' border="1" class="nice" >
  <tr>
    <th>Function</th>
    <th>Extension</th>
    <th>Location</th>
    <th>Procedure</th>
  </tr>
  <tr>
    <td>Ingress Control</td>
    <td>Contour</td>
    <td>Tanzu Kubernetes or Shared Service cluster</td>
    <td><a href="ingress-contour.md">Implementing Ingress Control with Contour</a></td>
  </tr><tr>
    <td>Service Discovery</td>
    <td>External DNS</td>
    <td>Tanzu Kubernetes or Shared Service cluster</td>
    <td><a href="external-dns.md">Implementing Service Discovery with External DNS</a></td>
  </tr><tr>
    <td>Log Forwarding</td>
    <td>Fluent Bit</td>
    <td>Tanzu Kubernetes cluster</td>
    <td><a href="logging-fluentbit.md">Implementing Log Forwarding with Fluentbit</a></td>
  </tr><tr>
    <td>Container Registry</td>
    <td>Harbor</td>
    <td>Shared Services cluster</td>
    <td><a href="harbor-registry.md">Deploy Harbor Registry as a Shared Service</a></td>
  </tr><tr>
    <td rowspan="2">Monitoring</td>
    <td>Prometheus</td>
    <td>Tanzu Kubernetes cluster</td>
    <td rowspan="2"><a href="monitoring.md">Implementing Monitoring with Prometheus and Grafana</a></td>
  </tr><tr>
    <td>Grafana</td>
    <td>Tanzu Kubernetes cluster</td>
  </tr>
</table>

Some extensions require or are enhanced by other extensions deployed to the same cluster:

* **Contour** is required by Harbor, External DNS, and Grafana
* **Prometheus** is required by Grafana
* **External DNS** is recommended for Harbor on infrastructures with load balancing (AWS, Azure, and vSphere with NSX Advanced Load Balancer), especially in production or other environments in which Harbor availability is important.

## <a id="preparing"></a> Preparing to Deploy the Extensions

Before you can deploy the Tanzu Kubernetes Grid extensions, you must prepare your bootstrap environment.

- To deploy the extensions, you update configuration files with information about your environment. You then use `kubectl` to apply preconfigured YAML files that pull data from the updated configuration files to create and update clusters that implement the extensions. The YAML files include calls to `ytt`, `kapp`, and `kbld` commands, so these tools must be present on your bootstrap environment when you deploy the extensions. For information about installing `ytt`, `kapp`, and `kbld`, see [Install the Carvel Tools](../install-cli#install-carvel).
- If you are using Tanzu Kubernetes Grid in an Internet-restricted environment, see [Deploying the Tanzu Kubernetes Grid Extensions in an Internet Restricted Environment](../mgmt-clusters/airgapped-environments.md#deploying-the-tanzu-kubernetes-grid-extensions-in-an-internet-restricted-environment-10).

## <a id="unpack-bundle"></a> Download and Unpack the Tanzu Kubernetes Grid Extensions Bundle

The Tanzu Kubernetes Grid extension manifests are provided in a separate bundle to the Tanzu CLI and other binaries.

1. On the system that you use as the bootstrap machine, go to [the Tanzu Kubernetes Grid downloads page](https://my.vmware.com/en/web/vmware/downloads/info/slug/infrastructure_operations_management/vmware_tanzu_kubernetes_grid/1_x) and log in with your My VMware credentials.
1. Under **Product Downloads**, click **Go to Downloads**.
1. Scroll to **VMware Tanzu Kubernetes Grid Extensions Manifest 1.4.0** and click **Download Now**.
1. Use either the `tar` command or the extraction tool of your choice to unpack the bundle of YAML manifest files for the Tanzu Kubernetes Grid extensions.

   <pre>tar -xzf tkg-extensions-manifests-v1.4.0-vmware.1.tar.gz</pre>

   For convenience, unpack the bundle in the same location as the one from which you run `tanzu` and `kubectl` commands.

**IMPORTANT**:

- After you unpack the bundle, the extensions files are contained in a folder named `tkg-extensions-v1.4.0+vmware.1
`. This folder contains subfolders for each type of extension, for example, `authentication`, `ingress`, `registry`, and so on. At the top level of the folder there is an additional subfolder named `extensions`. The `extensions` folder also contains subfolders for `authentication`, `ingress`, `registry`, and so on. In the procedures to deploy the extensions, take care to run commands from the location provided in the instructions. Commands are usually run from within the `extensions` folder.
- For historical reasons, the extensions bundle includes the manifests for the Dex and Gangway extensions. Tanzu Kubernetes Grid v1.3 introduces user authentication with Pinniped, that run automatically in management clusters if you enable identity management during deployment. For new deployments, enable Pinniped and Dex in your management clusters. Do not use the Dex and Gangway extensions. For information about identity management with Pinniped, see [Enabling Identity Management in Tanzu Kubernetes Grid](../mgmt-clusters/enabling-id-mgmt.md). For information about migrating existing Dex and Gangway deployments to Pinniped, see [Register Core Add-ons](../upgrade-tkg/addons.md).

## <a id="cert-mgr"></a> Install Cert Manager on Workload Clusters

Before you can deploy Tanzu Kubernetes Grid extensions, you must install `cert-manager`, which provides automated certificate management, on workload clusters. The `cert-manager` service already runs by default in management clusters.

All extensions other than Fluent Bit require `cert-manager` to be running on workload clusters. Fluent Bit does not use `cert-manager`.

To install the `cert-manager` service on a workload cluster, specify the cluster with `kubectl config use-context` and then do the following:

1. Deploy `cert-manager` on the cluster.

   ```sh
   kubectl apply -f cert-manager/
   ```

1. Check that the Kapp controller and cert-manager services are running as pods in the cluster.

   ```sh
   kubectl get pods -A
   ```

   The command output should show:

   - A `kapp-controller` pod with a name like `kapp-controller-cd55bbd6b-vt2c4` running in the namespace `tkg-system`.
   - For extensions other than Fluent Bit, pods with names like `cert-manager-69877b5f94-8kwx9`, `cert-manager-cainjector-7594d76f5f-8tstw`, and `cert-manager-webhook-5fc8c6dc54-nlvzp` running in the namespace `cert-manager`.
   - A `Ready` status of `1/1` for all of these pods. If this status is not displayed, stop and troubleshoot the pods before proceeding.

## <a id="shared"></a> Create a Shared Services Cluster

The Harbor service runs on a shared services cluster, to serve all the other clusters in an installation.
The Harbor service requires the Contour service to also run on the shared services cluster.
In many environments, the Harbor service also benefits from External DNS running on its cluster, as described in [Harbor Registry and External DNS](harbor-registry.md#external-dns).

Each Tanzu Kubernetes Grid instance can only have one shared services cluster.

To deploy a shared services cluster:

1. Create a cluster configuration YAML file for the target cluster.
To deploy to a shared services cluster, for example named `tkg-services`, it is recommended to use the `prod` cluster plan rather than the `dev` plan.
For example:

   ```
   INFRASTRUCTURE_PROVIDER: vsphere
   CLUSTER_NAME: tkg-services
   CLUSTER_PLAN: prod
   ```

1. **vSphere**: To deploy the cluster to vSphere, add a line to the configuration file that sets `VSPHERE_CONTROL_PLANE_ENDPOINT` to a static virtual IP (VIP) address for the control plane of the shared services cluster. Ensure that this IP address is not in the DHCP range, but is in the same subnet as the DHCP range. If you mapped a fully qualified domain name (FQDN) to the VIP address, you can specify the FQDN instead of the VIP address. For example:

   ```
   VSPHERE_CONTROL_PLANE_ENDPOINT: 10.10.10.10
   ```

1. Deploy the cluster by passing the cluster configuration file to the `tanzu cluster create`:

   ```
   tanzu cluster create tkg-services --file tkg-services-config.yaml
   ```

   Throughout the rest of these procedures, the cluster that you just deployed is referred to as the shared services cluster.

1. Set the context of `kubectl` to the context of your management cluster.
   For example, if your cluster is named `mgmt-cluster`, run the following command.

   ```
   kubectl config use-context mgmt-cluster-admin@mgmt-cluster
   ```

1. Add the label `tanzu-services` to the shared services cluster, as its cluster role. This label identifies the shared services cluster to the management cluster and workload clusters.

   ```
   kubectl label cluster.cluster.x-k8s.io/tkg-services cluster-role.tkg.tanzu.vmware.com/tanzu-services="" --overwrite=true
   ```

   You should see the confirmation `cluster.cluster.x-k8s.io/tkg-services labeled`.

1. Check that the label has been correctly applied by running the following command.

   ```
   tanzu cluster list --include-management-cluster
   ```
   You should see that the `tkg-services` cluster has the `tanzu-services` role.

   ```
     NAME              NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES           PLAN
     another-cluster   default     running  1/1           1/1      v1.20.5+vmware.1  <none>          dev
     tkg-services      default     running  3/3           3/3      v1.20.5+vmware.1  tanzu-services  prod
     mgmt-cluster      tkg-system  running  1/1           1/1      v1.20.5+vmware.1  management      dev
   ```

1. In a terminal, navigate to the folder that contains the unpacked Tanzu Kubernetes Grid extension manifest files, `tkg-extensions-v1.4.0+vmware.1
/extensions`.

   ```
   cd <path>/tkg-extensions-v1.4.0+vmware.1
/extensions
   ```

   You should see folders for `authentication`, `ingress`, `logging`, `monitoring`, `registry`, and some YAML files. Run all of the commands in this procedure from this location.

1. Get the `admin` credentials of the shared services cluster on which to deploy Harbor.

   ```
   tanzu cluster kubeconfig get tkg-services --admin
   ```

1. Set the context of `kubectl` to the shared services cluster.

   ```
   kubectl config use-context tkg-services-admin@tkg-services
   ```

## <a id="add-certs"></a> Add Certificates to the Kapp Controller

Previous versions of Tanzu Kubernetes Grid required the user to install the `kapp-controller` service to any extension cluster.
As of v1.3, all management and workload clusters are created with the `kapp-controller` service pre-installed.
If the cluster configuration file specifies a private registry with `TKG_CUSTOM_IMAGE_REPOSITORY` and `TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE` variables, the `kapp-controller` is configured to trust the private registry.

To enable a cluster's Kapp Controller to trust additional private registries, add their certificates to its configuration:

1. If needed, set the current `kubectl` context to the cluster with the Kapp Controller you are changing:

  ```
  kubectl config use-context CLUSTER-CONTEXT
  ```

1. Open the Kapp Controller's ConfigMap file in an editor:

  ```
  kubectl edit configmap -n tkg-system kapp-controller-config
  ```

1. Edit the ConfigMap file to add new certificates to the `data.caCerts` property:

  ```
  apiVersion: v1
  kind: ConfigMap
  metadata:
    # Name must be `kapp-controller-config` for kapp controller to pick it up
    name: kapp-controller-config
    # Namespace must match the namespace kapp-controller is deployed to
    namespace: tkg-system
  data:
    # A cert chain of trusted ca certs. These will be added to the system-wide
    # cert pool of trusted ca's (optional)
    caCerts: |
      -----BEGIN CERTIFICATE-----
      <Existing Certificate>
      -----END CERTIFICATE-----
      -----BEGIN CERTIFICATE-----
      <New Certificate>
      -----END CERTIFICATE-----
    # The url/ip of a proxy for kapp controller to use when making network requests (optional)
    httpProxy: ""
    # The url/ip of a tls capable proxy for kapp controller to use when making network requests (optional)
    httpsProxy: ""
    # A comma delimited list of domain names which kapp controller should bypass the proxy for when making requests (optional)
    noProxy: ""
    ```

1. Save the ConfigMap and exit the editor.

1. Delete the `kapp-controller` pod, so that it regenerates with the new configuration:

  ```
  kubectl delete pod -n tkg-system -l app=kapp-controller
  ```

## <a id="upgrading"></a>  Upgrading the Tanzu Kubernetes Grid Extensions

For information about how to upgrade the Tanzu Kubernetes Grid extensions from a previous release, see [Upgrade Tanzu Kubernetes Grid Extensions](../upgrade-tkg/extensions.md).
