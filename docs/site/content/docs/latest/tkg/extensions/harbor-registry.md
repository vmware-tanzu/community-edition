# Deploy Harbor Registry as a Shared Service

[Harbor](https://goharbor.io) is an open source, trusted, cloud native container registry that stores, signs, and scans content. Harbor extends the open source Docker Distribution by adding the functionalities usually required by users such as security, identity control and management.

Tanzu Kubernetes Grid includes signed binaries for Harbor, that you can deploy on a shared services cluster to provide container registry services for other Tanzu Kubernetes clusters. Unlike Tanzu Kubernetes Grid extensions, which you use to deploy services on individual clusters, you deploy Harbor as a shared service. In this way, Harbor is available to all of the Tanzu Kubernetes clusters in a given Tanzu Kubernetes Grid instance. To implement Harbor as a shared service, you deploy it on a special cluster that is reserved for running shared services in a Tanzu Kubernetes Grid instance.

You can use the Harbor shared service as a private registry for images that you want to make available to all of the Tanzu Kubernetes clusters that you deploy from a given management cluster. An advantage to using the Harbor shared service is that it is managed by Kubernetes, so it provides greater reliability than a standalone registry. Also, the Harbor implementation that Tanzu Kubernetes Grid provides as a shared service has been tested for use with Tanzu Kubernetes Grid and is fully supported.

The procedures in this topic all apply to vSphere, Amazon EC2, and Azure deployments.

## <a id="internet-restricted"></a> Using the Harbor Shared Service in Internet-Restricted Environments

Another use-case for deploying Harbor as a shared service is for Tanzu Kubernetes Grid deployments in Internet-restricted environments. For more information, see [Using the Harbor Shared Service in Internet-Restricted Environments](../mgmt-clusters/airgapped-environments.md#internet-restricted-harbor).

## <a id="external-dns"></a> Harbor Registry and External DNS

VMware recommends installing the External DNS service alongside the Harbor Registry on infrastructures with load balancing (AWS, Azure, and vSphere with NSX Advanced Load Balancer), especially in production or other environments in which Harbor availability is important.

If the IP address to the shared services ingress load balancer changes, External DNS automatically picks up the change and re-maps the new address to the Harbor hostname.
This precludes the need to manually re-map the address as described in [Connect to the Harbor User Interface](#connect-ui).

## <a id="prereqs"></a> Prerequisites

- You have deployed a management cluster on vSphere, Amazon EC2, or Azure, in either an Internet-connected or Internet-restricted environment.

   If you are using Tanzu Kubernetes Grid in an Internet-restricted environment, you performed the procedure in [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](../mgmt-clusters/airgapped-environments.md) before you deployed the management cluster.
- You have downloaded and unpacked the bundle of Tanzu Kubernetes Grid extensions. For information about where to obtain the bundle, see [Download and Unpack the Tanzu Kubernetes Grid Extensions Bundle](index.md#unpack-bundle).  
- You have installed the Carvel tools. For information about installing the Carvel tools, see [Install the Carvel Tools](../install-cli.md#install-carvel).
- You have installed [yq](https://github.com/mikefarah/yq/releases):
   - For Tanzu Kubernetes Grid v1.3.0, install `yq` v3.
   - For Tanzu Kubernetes Grid v1.3.1 and later, install `yq` v4.5 or later.

**IMPORTANT**: The extensions  folder `tkg-extensions-v1.4.0+vmware.1
` contains subfolders for each type of extension, for example, `authentication`, `ingress`, `registry`, and so on. At the top level of the folder there is an additional subfolder named `extensions`. The `extensions` folder also contains subfolders for `authentication`, `ingress`, `registry`, and so on. Take care to run commands from the location provided in the instructions. Commands are usually run from within the `extensions` folder.

## <a id="prepare-tkc"></a> Prepare a Shared Services Cluster for Harbor Deployment

Each Tanzu Kubernetes Grid instance can only have one shared services cluster. You must deploy Harbor on a cluster that will only be used for shared services.

To prepare a shared services cluster for running Harbor Extension on:

1. Create a shared services cluster, if it is not already created, by following the procedure [Create a Shared Services Cluster](index.md#shared).

1. Deploy Contour Extension on the shared services cluster.

   Harbor Extension requires Contour Extension to be present on the cluster, to provide ingress control. For how to deploy Contour Extension, see [Deploy Contour on the Tanzu Kubernetes Cluster](ingress-contour.md#deploy).

1. (Optional) Deploy External DNS Extension on the shared services cluster.
  External DNS Extension is recommended for using Harbor Extension in environments with load balancing, as described in [Harbor Registry and External DNS](#external-dns), above.

Your shared services cluster is now ready for you to deploy the Harbor Extension on it.

## <a id="deploy"></a> Deploy Harbor Extension on the Shared Services Cluster

After you have deployed a shared services cluster that includes the Contour Extension, you can deploy the Harbor Extension.

1. Set the context of `kubectl` to the shared services cluster.

   ```
   kubectl config use-context tkg-services-admin@tkg-services
   ```

1. Create a namespace for the Harbor Extension on the shared services cluster.

    ```sh
    kubectl apply -f registry/harbor/namespace-role.yaml
    ```

    You should see confirmation that a `tanzu-system-registry` namespace, service account, and RBAC role bindings are created.

    ```
    namespace/tanzu-system-registry created
    serviceaccount/harbor-extension-sa created
    role.rbac.authorization.k8s.io/harbor-extension-role created
    rolebinding.rbac.authorization.k8s.io/harbor-extension-rolebinding created
    clusterrole.rbac.authorization.k8s.io/harbor-extension-cluster-role created
    clusterrolebinding.rbac.authorization.k8s.io/harbor-extension-cluster-rolebinding created
    ```

1. Make a copy of the `harbor-data-values.yaml.example` file and name it `harbor-data-values.yaml`.

    ```sh
    cp registry/harbor/harbor-data-values.yaml.example registry/harbor/harbor-data-values.yaml
    ```

    The `harbor-data-values.yaml` configures the Harbor extension.
    You can also customize your Harbor setup using `ytt` overlays.
    See [`ytt` Overlays and Example: Clean Up S3 and Trust Let's Encrypt](#ytt) below.

1. Set the mandatory passwords and secrets in `harbor-data-values.yaml`.

   You can do this in one of two ways:

   - To automatically generate random passwords and secrets, run the following command:

      ```
      bash registry/harbor/generate-passwords.sh registry/harbor/harbor-data-values.yaml
      ```  

   - To set your own passwords and secrets, update the following entries in  `harbor-data-values.yaml`:
      - `harborAdminPassword`
      - `secretKey`
      - `database.password`
      - `core.secret`
      - `core.xsrfKey`
      - `jobservice.secret`
      - `registry.secret`

1. Specify other settings in `harbor-data-values.yaml`.

   - Set the `hostname` setting to the hostname you want to use to access Harbor. For example: `harbor.yourdomain.com`.
   - To use your own certificates, update the `tls.crt`, `tls.key`, and `ca.crt` settings with the contents of your certificate, key, and CA certificate. The certificate can be signed by a trusted authority or be self-signed. If you leave these blank, Tanzu Kubernetes Grid automatically generates a self-signed certificate.
   - If you used the `generate-passwords.sh` script, optionally update the `harborAdminPassword` with something that is easier to remember.
   - Optionally update the `persistence` settings to specify how Harbor stores data.

      If you need to store a large quantity of container images in Harbor, set `persistence.persistentVolumeClaim.registry.size` to a larger number.

      If you do not update the `storageClass` under `persistence` settings, Harbor uses the shared services cluster's default `storageClass`. If the default `storageClass` or a `storageClass` that you specify in `harbor-data-values.yaml` supports the `accessMode` [`ReadWriteMany`](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes), you must update the `persistence.persistentVolumeClaim` `accessMode` settings for `registry`, `jobservice`, `database`, `redis`, and `trivy` from `ReadWriteOnce` to `ReadWriteMany`. [vSphere 7 with VMware vSAN 7 supports `accessMode: ReadWriteMany`](https://blogs.vmware.com/virtualblocks/2020/03/12/cloud-native-storage-and-vsan-file-services-integration/) but vSphere 6.7u3 does not. If you are using vSphere 7 without vSAN, or you are using vSphere 6.7u3, use the default value `ReadWriteOnce`.

   - Optionally update the other Harbor settings. The settings that are available in `harbor-data-values.yaml` are a subset of the settings that you set when deploying open source Harbor with Helm. For information about the other settings that you can configure, see [Deploying Harbor with High Availability via Helm](https://goharbor.io/docs/2.0.0/install-config/harbor-ha-helm/) in the Harbor documentation.
1. Create a Kubernetes secret named `harbor-data-values` with the values that you set in `harbor-data-values.yaml`.

    ```sh
    kubectl create secret generic harbor-data-values --from-file=values.yaml=registry/harbor/harbor-data-values.yaml -n tanzu-system-registry
    ```

1. Deploy the Harbor extension.

    ```sh
    kubectl apply -f registry/harbor/harbor-extension.yaml
    ```

   You should see the confirmation `extension.clusters.tmc.cloud.vmware.com/harbor created`.

1. View the status of the Harbor extension.

   ```
   kubectl get extension harbor -n tanzu-system-registry
   ```

   You should see information about the Harbor extension.

    ```
    NAME     STATE   HEALTH   VERSION
    harbor   3
    
    ```

1. View the status of the Harbor service itself.

    ```
    kubectl get app harbor -n tanzu-system-registry
    ```

    The status of the Harbor app should show `Reconcile Succeeded` when Harbor has deployed successfully.

    ```
    NAME     DESCRIPTION           SINCE-DEPLOY   AGE
    harbor   Reconcile succeeded   3m11s          23m
    ```

1. If the status is not `Reconcile Succeeded`, view the full status details of the Harbor service.

   Viewing the full status can help you to troubleshoot the problem.

    ```
    kubectl get app harbor -n tanzu-system-registry -o yaml
    ```

1. Check that the new services are running by listing all of the pods that are running in the cluster.

    ```
    kubectl get pods -A
    ```

   In the `tanzu-system-regisry` namespace, you should see the `harbor` `core`, `clair`, `database`, `jobservice`, `notary`, `portal`, `redis`, `registry`, and `trivy` services running in a pod with names similar to `harbor-registry-76b6ccbc75-vj4jv`.

    ```
    NAMESPACE               NAME                                    READY   STATUS    RESTARTS   AGE
    [...]
    tanzu-system-ingress    contour-6b568c9b88-h5s2r                1/1     Running   0          26m
    tanzu-system-ingress    contour-6b568c9b88-mlg2r                1/1     Running   0          26m
    tanzu-system-ingress    envoy-wfqdp                             2/2     Running   0          26m
    tanzu-system-registry   harbor-clair-9ff9b98d-6vlk4             2/2     Running   1          23m
    tanzu-system-registry   harbor-core-557b58b65c-4kzhn            1/1     Running   0          23m
    tanzu-system-registry   harbor-database-0                       1/1     Running   0          23m
    tanzu-system-registry   harbor-jobservice-847b5c8756-t6kfs      1/1     Running   0          23m
    tanzu-system-registry   harbor-notary-server-6b74b8dd56-d7swb   1/1     Running   2          23m
    tanzu-system-registry   harbor-notary-signer-69d4669884-dglzm   1/1     Running   2          23m
    tanzu-system-registry   harbor-portal-8f677757c-t4cbj           1/1     Running   0          23m
    tanzu-system-registry   harbor-redis-0                          1/1     Running   0          23m
    tanzu-system-registry   harbor-registry-85b96c7777-wsdnj        2/2     Running   0          23m
    tanzu-system-registry   harbor-trivy-0                          1/1     Running   0          23m
    tkg-system              kapp-controller-778b5f484c-fkbvg        1/1     Running   0          59m
    vmware-system-tmc       extension-manager-6c64cdd984-s99gc      1/1     Running   0          27m
    ```

1. Obtain the Harbor CA certificate from the `harbor-tls` secret in the `tanzu-system-registry namespace`.

    ```
    kubectl -n tanzu-system-registry get secret harbor-tls -o=jsonpath="{.data.ca\.crt}" | base64 -d
    ```

    Make a copy of the output.

## <a id="connect-ui"></a> Connect to the Harbor User Interface

The Harbor UI is exposed via the Envoy service load balancer that is running in the Contour extension on the shared services cluster. To allow users to connect to the Harbor UI, you must map the address of the Envoy service load balancer to the hostname of the Harbor service, for example `harbor.yourdomain.com`. How you map the address of the Envoy service load balancer to the hostname depends on whether your Tanzu Kubernetes Grid instance is running on vSphere, on Amazon EC2 or on Azure.

1. Obtain the address of the Envoy service load balancer.

   ```
   kubectl get svc envoy -n tanzu-system-ingress -o jsonpath='{.status.loadBalancer.ingress[0]}'
   ```

   On **vSphere without NSX Advanced Load Balancer (ALB)**, the Envoy service is exposed via NodePort instead of LoadBalancer, so the above output will be empty, and you can use the IP address of any worker node in the shared services cluster instead. On **Amazon EC2**, it has a FQDN similar to `a82ebae93a6fe42cd66d9e145e4fb292-1299077984.us-west-2.elb.amazonaws.com`.
   On **vSphere with NSX ALB** and **Azure**, the Envoy service has a Load Balancer IP address similar to `20.54.226.44`.

1. Map the address of the Envoy service load balancer to the hostname of the Harbor service.

   - **vSphere**: If you deployed Harbor on a shared services cluster that is running on vSphere, you must add an IP to hostname mapping in `/etc/hosts` or add corresponding `A` records in your DNS server. For example, if the IP address is `10.93.9.100`, add the following to `/etc/hosts`:

       ```
       10.93.9.100 harbor.yourdomain.com notary.harbor.yourdomain.com
       ```

      On Windows machines, the equivalent to `/etc/hosts/` is `C:\Windows\System32\Drivers\etc\hosts`.

   - **Amazon EC2 or Azure**: If you deployed Harbor on a shared services cluster that is running on Amazon EC2 or Azure, you must create two DNS `CNAME` records (on Amazon EC2) or two DNS `A` records (on Azure) for the Harbor hostnames on a DNS server on the Internet.
       - One record for the Harbor hostname, for example, `harbor.yourdomain.com`, that you configured in `harbor-data-values.yaml`, that points to the FQDN or IP of the Envoy service load balancer.
       - Another record for the Notary service that is running in Harbor, for example, `notary.harbor.yourdomain.com`, that points to the FQDN or IP of the Envoy service load balancer.

Users can now connect to the Harbor UI by navigating to `https://harbor.yourdomain.com` in a Web browser and log in as user `admin` with the `harborAdminPassword` that you configured in `harbor-data-values.yaml`.

## <a id="push-pull"></a> Push and Pull Images to and from the Harbor Extension

Now that Harbor is set up as a shared service, you can push images to it to make them available for your Tanzu Kubernetes clusters to pull.

1. If Harbor uses a self-signed certificate, download the Harbor CA certificate from `https://harbor.yourdomain.com/api/v2.0/systeminfo/getcert`, and install it on your local machine, so Docker can trust this CA certificate.

   - On Linux, save the certificate as `/etc/docker/certs.d/harbor.yourdomain.com/ca.crt`.
   - On macOS, follow [this procedure](https://blog.container-solutions.com/adding-self-signed-registry-certs-docker-mac).
   - On Windows, right-click the certificate file and select **Install Certificate**.

1. Log in to the Harbor registry with the user `admin`. When prompted, enter the `harborAdminPassword` that you set when you deployed the Harbor Extension on the shared services cluster.

   ```
   docker login harbor.yourdomain.com -u admin
   ```

1. Tag an existing image that you have already pulled locally, for example `nginx:1.7.9`.

   ```
   docker tag nginx:1.7.9 harbor.yourdomain.com/library/nginx:1.7.9
   ```

1. Push the image to the Harbor registry.

   ```
   docker push harbor.yourdomain.com/library/nginx:1.7.9
   ```

1. Now you can pull the image from the Harbor registry on any machine where the Harbor CA certificate is installed.

   ```
   docker pull harbor.yourdomain.com/library/nginx:1.7.9
   ```

## <a id="populate"></a> Push the Tanzu Kubernetes Grid Images into the Harbor Registry

The Tanzu Kubernetes Grid images are published in a public container registry and used by Tanzu Kubernetes Grid to deploy Tanzu Kubernetes clusters and extensions. When creating a Tanzu Kubernetes cluster, in order for Tanzu Kubernetes cluster nodes to pull Tanzu Kubernetes Grid images from the Harbor shared service rather than over the Internet, you must first push those images to the Harbor shared service.

This procedure is optional if your Tanzu Kubernetes clusters have internet connectivity to pull external images.

If you only want to store your application images rather than the Tanzu Kubernetes Grid images in the Harbor shared service, follow the procedure in [Trust Custom CA Certificates on Cluster Nodes](../cluster-lifecycle/secrets.md#custom-ca) to enable the Tanzu Kubernetes cluster nodes to pull images from the Harbor shared service, and skip the rest of this procedure.

**NOTE**: If your Tanzu Kubernetes Grid instance is running in an Internet-restricted environment, you must perform these steps on a machine that has an Internet connection, that can also access the Harbor registry that you have just deployed as a shared service.

1. Create a public project named `tkg` from Harbor UI. Or, you can use another project name.

1. Set the FQDN of the Harbor registry that is running as a shared service as an environment variable.

   On Windows platforms, use the `SET` command instead of `export`. Include the name of the default project in the variable. For example, if you set the Harbor hostname to `harbor.yourdomain.com`, set the following:

   ```
   export TKG_CUSTOM_IMAGE_REPOSITORY=harbor.yourdomain.com/tkg
   ```

1. Follow the step 2 and step 3 in [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](../mgmt-clusters/airgapped-environments.md) to generate and run the `publish-images.sh` script.

1. When the script finishes, add or update the following rows in the global configuration file, `~/.tanzu/tkg/config.yaml`.

   These variables ensure that when creating a Management Cluster or Tanzu Kubernetes Cluster, Tanzu Kubernetes Grid always pulls Tanzu Kubernetes Grid images from the Harbor Extension that is running as a shared service, rather than from the external internet. If your Harbor Extension uses self-signed certificates, also add `TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE` to the configuration file. Provide the CA certificate in base64 encoded format.

    ```
    TKG_CUSTOM_IMAGE_REPOSITORY: harbor.yourdomain.com/tkg
    TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE: LS0t[...]tLS0tLQ==
    ```

   If your Tanzu Kubernetes Grid instance is running in an Internet-restricted environment, you can disconnect the Internet connection now.

You can now use the `tanzu cluster create` command to deploy Tanzu Kubernetes clusters, and the images will be pulled from the Harbor Extension that is running in the shared services cluster. You can push images to the Harbor registry to make them available to all clusters that are running in the Tanzu Kubernetes Grid instance.

Connections between Tanzu Kubernetes cluster nodes and Harbor are secure, regardless of whether you use a trusted or a self-signed certificate for the Harbor shared service.

## <a id="ytt"></a> `ytt` Overlays and Example: Clean Up S3 and Trust Let's Encrypt

In addition to modifying `harbor-data-values.yaml`, you can use `ytt` overlays to configure your Harbor setup, as described in [Extensions and Shared Services](../ytt.md#extensions) in _Customizing Clusters, Plans, and Extensions with ytt Overlays_ and in the extensions mods examples in the [TKG Lab repository](https://github.com/Tanzu-Solutions-Engineering/tkg-lab).

One TKG Lab example, in the step [Prepare Manifests and Deploy Harbor Extension](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/docs/shared-services-cluster/10_harbor.md#prepare-manifests-and-deploy-harbor-extension), cleans PersistentVolumeClaim (PVC) info from Harbor's S3 storage and lets the Harbor extension trust a [Let's Encrypt](https://letsencrypt.org/certificates/) certificate authority.

The example procedure does this by running a script [`generate-and-apply-harbor-yaml.sh`](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/scripts/generate-and-apply-harbor-yaml.sh) that sets up the configuration files used to deploy the Harbor extension.
To customize the Harbor extension, the script applies three `ytt` overlay files:

- [`overlay-s3-pvc-fix.yaml`](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/tkg-extensions-mods-examples/registry/harbor/overlay-s3-pvc-fix.yaml) clears PVC data to allow the Harbor registry to use S3 for data storage
- [`trust-certificate/overlay.yaml`](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/overlay/trust-certificate/overlay.yaml) lets any extension (Harbor in this case) trust a Let's Encrypt CA; useful for OIDC providers with Let's Encrypt-based certs
- [`harbor-extension-overlay.yaml`](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/tkg-extensions-mods-examples/registry/harbor/harbor-extension-overlay.yaml) directs the Harbor extension to always deploy the Harbor app with the previous two overlays

See the TKG Lab repository and its [Step by Step setup guide](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/docs/baseline-lab-setup/step-by-step.md) for more examples.

## <a id="update"></a> Update a Running Harbor Deployment

If you need to make changes to the configuration of the Harbor extension after deployment, follow these steps to update your deployed Harbor extension.

1. Update the Harbor configuration in `registry/harbor/harbor-data-values.yaml`.

   For example, increase the amount of registry storage by updating the `persistence.persistentVolumeClaim.registry.size` value.

1. Update the Kubernetes secret, which contains the Harbor configuration.

   This command assumes that you are running it from `tkg-extensions-v1.4.0+vmware.1
/extensions`.

   ```sh
    kubectl create secret generic harbor-data-values --from-file=values.yaml=registry/harbor/harbor-data-values.yaml -n tanzu-system-registry -o yaml --dry-run | kubectl replace -f-
    ```

   Note that the final `-` on the `kubectl replace` command above is necessary to instruct `kubectl` to accept the input being piped to it from the `kubectl create secret` command.

   The Harbor extension will be reconciled using the new values you just added. The changes should show up in five minutes or less. This is handled by the Kapp controller, which synchronizes every five minutes.
