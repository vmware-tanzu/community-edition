# Harbor Package

[Harbor](https://github.com/goharbor/harbor) is an open source trusted cloud native registry project that stores, signs, and scans content. Harbor extends the open source Docker Distribution by adding the functionalities usually required by users such as security, identity and management.

## Components

This Harbor Package integrates open source Harbor 2.3.3. See [docs for Harbor 2.3.3](https://goharbor.io/docs/2.3.0/install-config/#harbor-components).

## Configuration

The following configuration values can be set to customize the harbor installation.

### Global

| Value | Required/Optional | Default | Description |
|-------|-------------------|---------|-------------|
| `namespace` | Optional | harbor | The namespace in which to deploy Harbor.|

### Harbor Configuration

To get the configuration for this package, you will need to refer to the [TCE GitHub
repository](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages).
Select the package/version and navigate into the `bundle/config` directory.
Download or copy/paste the `values.yaml` file and rename it to `harbor-values.yaml`.

or get the template configuration file by using script below:

   ```shell
   image_url=$(kubectl get packages harbor.community.tanzu.vmware.com.2.3.3 -o jsonpath='{.spec.template.spec.fetch[0].imgpkgBundle.image}')
   imgpkg pull -b $image_url -o /tmp/harbor-package-PACKAGE-VERSION
   cp /tmp/harbor-package-PACKAGE-VERSION/config/values.yaml harbor-values.yaml
   ```

> When you are using `imgpkg` to get the configuratuion file, specifying a namespace may be required
> depending on where your package repository was installed.

|Values | Required/Optional | Default | Type | Description |
|-------|-------------------|---------|------|--------------|
|` logLevel `|   Required   | info | string | The log level of core, exporter, jobservice, registry. |
|` metrics.core.path `|   Required   | /metrics | string | The path of the metrics. |
|` metrics.core.port `|   Required   | 8001 | integer | The port of the metrics. |
|` metrics.enabled `|   Required   | false | boolean | Enable the metrics when it's true |
|` metrics.exporter.path `|   Required   | /metrics | string | The path of the metrics. |
|` metrics.exporter.port `|   Required   | 8001 | integer | The port of the metrics. |
|` metrics.jobservice.path `|   Required   | /metrics | string | The path of the metrics. |
|` metrics.jobservice.port `|   Required   | 8001 | integer | The port of the metrics. |
|` metrics.registry.path `|   Required   | /metrics | string | The path of the metrics. |
|` metrics.registry.port `|   Required   | 8001 | integer | The port of the metrics. |
|` namespace `|   Required   | harbor | string | The namespace to install Harbor. |
|` proxy.httpProxy `|   Required   |  | string | HTTP proxy URL. |
|` proxy.httpsProxy `|   Required   |  | string | HTTPS proxy URL. |
|` proxy.noProxy `|   Required   | 127.0.0.1,localhost,.local,.internal | string | Ignore proxy for the domains. |
|` pspNames `|   Required   | null | string | The PSP names used by Harbor pods. The names are separated by ','. 'null' means all PSP can be used. |
|` enableContourHttpProxy `|   Required   | true | boolean | Use contour http proxy instead of the ingress when it's true. |
|` hostname `|   Required   | harbor.yourdomain.com | string | The FQDN for accessing Harbor admin UI and Registry service. |
|` secretKey `|   Required   |  | string | The secret key used for encryption. Must be a string of 16 chars. |
|` tlsCertificate.ca.crt `|   Required   |  | string | The certificate of CA, this enables the download, link on portal to download the certificate of CA. Note that ca.crt is a key and not nested. |
|` tlsCertificate.tls.crt `|   Required   |  | string | The certificate. Note that tls.crt is a key and not nested. |
|` tlsCertificate.tls.key `|   Required   |  | string | The private key. Note that tls.key is a key and not nested. |
|` network.ipFamilies `|   Required   | [IPv4 IPv6] | array | THe array of network ipFamilies. |
|` port.https `|   Required   | 443 | integer | The network port of the Envoy service in Contour or other Ingress Controller. |
|` notary.enabled `|   Required   | true | boolean | Whether to install Notary |
|` trivy.skipUpdate `|   Required   | false | boolean | The flag to disable Trivy DB downloads from GitHub. |
|` trivy.enabled `|   Required   | true | boolean | Whether to install Trivy scanner. |
|` trivy.gitHubToken `|   Required   |  | string | the GitHub access token to download Trivy DB. |
|` trivy.replicas `|   Required   | 1 | integer | The replicas for the trivy component. |
|` database.shmSizeLimit `|   Required   |  | integer | The initial value of shmSizeLimit |
|` database.maxIdleConns `|   Required   |  | integer | The initial value of maxIdleConns |
|` database.maxOpenConns `|   Required   |  | integer | The initial value of maxOpenConns |
|` database.password `|   Required   |  | string | The initial password of the postgres database. |
|` jobservice.secret `|   Required   |  | string | Secret is used when job service communicates with other components. |
|` jobservice.replicas `|   Required   | 1 | integer | The replicas for the jobservice component. |
|` harborAdminPassword `|   Required   |  | string | The initial password of Harbor admin. |
|` persistence.imageChartStorage.caBundleSecretName `|   Optional   |  | string | Specify the "caBundleSecretName" if the storage service uses a self-signed certificate. The secret must contain keys named "ca.crt" which will be injected into the trust store of registry's and chartmuseum's containers. |
|` persistence.imageChartStorage.disableredirect `|   Optional   | false | boolean | Specify whether to disable `redirect` for images and chart storage, for backends which not supported it (such as using minio for `s3` storage type), please disable it. To disable redirects, simply set `disableredirect` to `true` instead. Refer to <https://github.com/docker/distribution/blob/master/docs/configuration.md#redirect> for the detail. |
|` persistence.imageChartStorage.oss.rootdirectory `|   Optional   |  | string | The rootdirectory in oss. |
|` persistence.imageChartStorage.oss.accesskeysecret `|   Optional   | accesskeysecret | string | Access key secert of oss. |
|` persistence.imageChartStorage.oss.bucket `|   Optional   | bucketname | string | Bucket name of oss. |
|` persistence.imageChartStorage.oss.chunksize `|   Optional   |  | string | Chunk size for the oss, eg 10M. |
|` persistence.imageChartStorage.oss.endpoint `|   Optional   |  | string | Endpoint of oss. |
|` persistence.imageChartStorage.oss.internal `|   Optional   |  | boolean | Use the internal endpoint when it's true. |
|` persistence.imageChartStorage.oss.accesskeyid `|   Optional   | accesskeyid | string | Access key id of oss. |
|` persistence.imageChartStorage.oss.encrypt `|   Optional   |  | boolean | Encrypt of oss. |
|` persistence.imageChartStorage.oss.region `|   Optional   | regionname | string | Region of oss. |
|` persistence.imageChartStorage.oss.secure `|   Optional   |  | boolean | Secure of oss. |
|` persistence.imageChartStorage.azure.accountkey `|   Optional   | base64encodedaccountkey | string | Account key of azure storage. |
|` persistence.imageChartStorage.azure.accountname `|   Optional   | accountname | string | Account name of azure storage. |
|` persistence.imageChartStorage.azure.container `|   Optional   | containername | string | Container name of azure storage. |
|` persistence.imageChartStorage.azure.realm `|   Optional   | core.windows.net | string | Realm for azure storage. |
|` persistence.imageChartStorage.filesystem.maxthreads `|   Optional   | 100 | integer | Max threads for filesystem. |
|` persistence.imageChartStorage.filesystem.rootdirectory `|   Optional   | /storage | string | The rootdirectory in filesystem. |
|` persistence.imageChartStorage.gcs.chunksize `|   Optional   | 5.24288e+06 | integer | Check size for gcs. |
|` persistence.imageChartStorage.gcs.encodedkey `|   Optional   | base64-encoded-json-key-file | string | The base64 encoded json file which contains the key |
|` persistence.imageChartStorage.gcs.rootdirectory `|   Optional   |  | string | The rootdirectory in gcs. |
|` persistence.imageChartStorage.gcs.bucket `|   Optional   | bucketname | string | Bucket name of gcs. |
|` persistence.imageChartStorage.s3.accesskey `|   Optional   |  | string | Access key of s3. |
|` persistence.imageChartStorage.s3.multipartcopychunksize `|   Optional   |  | integer | multi part copy chunk size of s3. |
|` persistence.imageChartStorage.s3.region `|   Optional   | us-west-1 | string | Region of s3. |
|` persistence.imageChartStorage.s3.skipverify `|   Optional   | false | boolean | skipverify for s3. |
|` persistence.imageChartStorage.s3.multipartcopythresholdsize `|   Optional   |  | integer | multi part copy threshold size of s3. |
|` persistence.imageChartStorage.s3.regionendpoint `|   Optional   |  | string | Region endpoint of s3, eg <http://myobjects.local> |
|` persistence.imageChartStorage.s3.v4auth `|   Optional   | true | boolean | Use v4auth for s3 when it's true. |
|` persistence.imageChartStorage.s3.bucket `|   Optional   | bucketname | string | Bucket name of s3. |
|` persistence.imageChartStorage.s3.encrypt `|   Optional   | false | boolean | Encrypt for s3. |
|` persistence.imageChartStorage.s3.keyid `|   Optional   |  | string | Keyid of s3. |
|` persistence.imageChartStorage.s3.secretkey `|   Optional   |  | string | Secret key of s3. |
|` persistence.imageChartStorage.s3.secure `|   Optional   | true | boolean | Secure for s3. |
|` persistence.imageChartStorage.s3.chunksize `|   Optional   |  | integer | Check size for s3. |
|` persistence.imageChartStorage.s3.multipartcopymaxconcurrency `|   Optional   |  | integer | multi part copy max concurrency of s3. |
|` persistence.imageChartStorage.s3.rootdirectory `|   Optional   |  | string | The rootdirectory in s3. |
|` persistence.imageChartStorage.s3.storageclass `|   Optional   | STANDARD | string | Storage class of s3. |
|` persistence.imageChartStorage.swift.domain `|   Optional   |  | string | Domain of swift. |
|` persistence.imageChartStorage.swift.prefix `|   Optional   |  | string | Prefix path of swift. |
|` persistence.imageChartStorage.swift.trustid `|   Optional   |  | string | Trust id of swift. |
|` persistence.imageChartStorage.swift.username `|   Optional   | username | string | Username of swift. |
|` persistence.imageChartStorage.swift.accesskey `|   Optional   |  | string | Access key of swift. |
|` persistence.imageChartStorage.swift.authurl `|   Optional   | <https://storage.myprovider.com/v3/auth> | string | Auth url of swift. |
|` persistence.imageChartStorage.swift.domainid `|   Optional   |  | string | Domain id of swift. |
|` persistence.imageChartStorage.swift.region `|   Optional   |  | string | Region of swift. |
|` persistence.imageChartStorage.swift.tempurlmethods `|   Optional   |  | string | Temp url methods of swift. |
|` persistence.imageChartStorage.swift.authversion `|   Optional   |  | string | Auth version of swift. |
|` persistence.imageChartStorage.swift.secretkey `|   Optional   |  | string | Secret key of swift. |
|` persistence.imageChartStorage.swift.tenantid `|   Optional   |  | string | Tenant id of swift. |
|` persistence.imageChartStorage.swift.tempurlcontainerkey `|   Optional   |  | boolean | Use temp url container key of swift when it's true. |
|` persistence.imageChartStorage.swift.tenant `|   Optional   |  | string | Tenant of swift. |
|` persistence.imageChartStorage.swift.chunksize `|   Optional   |  | string | Check size of swift, eg 5M. |
|` persistence.imageChartStorage.swift.container `|   Optional   | containername | string | Container of swift. |
|` persistence.imageChartStorage.swift.endpointtype `|   Optional   |  | string | Endpoint type of swift, eg public. |
|` persistence.imageChartStorage.swift.insecureskipverify `|   Optional   |  | boolean | Ignore the cert verify when it's true. |
|` persistence.imageChartStorage.swift.password `|   Optional   | password | string | Password of swift. |
|` persistence.imageChartStorage.type `|   Optional   | filesystem | string | Specify the type of storage: "filesystem", "azure", "gcs", "s3", "swift", "oss" and fill the information needed in the corresponding section. The type must be "filesystem" if you want to use persistent volumes for registry and chartmuseum |
|` persistence.persistentVolumeClaim.redis.size `|   Optional   | 1Gi | string | Size of the PVC. |
|` persistence.persistentVolumeClaim.redis.storageClass `|   Optional   |  | string | Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |
|` persistence.persistentVolumeClaim.redis.subPath `|   Optional   |  | string | The "subPath" if the PVC is shared with other components. |
|` persistence.persistentVolumeClaim.redis.accessMode `|   Optional   | ReadWriteOnce | string | Access mode of the PVC. |
|` persistence.persistentVolumeClaim.redis.existingClaim `|   Optional   |  | string | Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |
|` persistence.persistentVolumeClaim.registry.existingClaim `|   Optional   |  | string | Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |
|` persistence.persistentVolumeClaim.registry.size `|   Optional   | 10Gi | string | Size of the PVC. |
|` persistence.persistentVolumeClaim.registry.storageClass `|   Optional   |  | string | Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |
|` persistence.persistentVolumeClaim.registry.subPath `|   Optional   |  | string | The "subPath" if the PVC is shared with other components. |
|` persistence.persistentVolumeClaim.registry.accessMode `|   Optional   | ReadWriteOnce | string | Access mode of the PVC. |
|` persistence.persistentVolumeClaim.trivy.size `|   Optional   | 5Gi | string | Size of the PVC. |
|` persistence.persistentVolumeClaim.trivy.storageClass `|   Optional   |  | string | Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |
|` persistence.persistentVolumeClaim.trivy.subPath `|   Optional   |  | string | The "subPath" if the PVC is shared with other components. |
|` persistence.persistentVolumeClaim.trivy.accessMode `|   Optional   | ReadWriteOnce | string | Access mode of the PVC. |
|` persistence.persistentVolumeClaim.trivy.existingClaim `|   Optional   |  | string | Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |
|` persistence.persistentVolumeClaim.database.accessMode `|   Optional   | ReadWriteOnce | string | Access mode of the PVC. |
|` persistence.persistentVolumeClaim.database.existingClaim `|   Optional   |  | string | Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |
|` persistence.persistentVolumeClaim.database.size `|   Optional   | 1Gi | string | Size of the PVC. |
|` persistence.persistentVolumeClaim.database.storageClass `|   Optional   |  | string | Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |
|` persistence.persistentVolumeClaim.database.subPath `|   Optional   |  | string | The "subPath" if the PVC is shared with other components. |
|` persistence.persistentVolumeClaim.jobservice.accessMode `|   Optional   | ReadWriteOnce | string | Access mode of the PVC. |
|` persistence.persistentVolumeClaim.jobservice.existingClaim `|   Optional   |  | string | Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |
|` persistence.persistentVolumeClaim.jobservice.size `|   Optional   | 1Gi | string | Size of the PVC. |
|` persistence.persistentVolumeClaim.jobservice.storageClass `|   Optional   |  | string | Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |
|` persistence.persistentVolumeClaim.jobservice.subPath `|   Optional   |  | string | The "subPath" if the PVC is shared with other components. |
|` registry.replicas `|   Required   | 1 | integer | The replicas for the registry component. |
|` registry.secret `|   Required   |  | string | Secret is used to secure the upload state from client and registry storage backend. |
|` core.secret `|   Required   |  | string | Secret is used when core server communicates with other components. |
|` core.xsrfKey `|   Required   |  | string | The XSRF key. Must be a string of 32 chars. |
|` core.replicas `|   Required   | 1 | integer | The replicas for the core component. |
|` exporter.cacheDuration `|   Required   |  | integer | The initial value of cacheDuration. |

Please refer the following step to configure Harbor.

## Installation

The Harbor package requires use of Contour for ingress and cert-manager for certificate generation.

1. Install cert-manager Package

   ```shell
   tanzu package install cert-manager \
      --package-name cert-manager.community.tanzu.vmware.com \
      --version ${CERT_MANAGER_PACKAGE_VERSION}
   ```

   > You can get the `${CERT_MANAGER_PACKAGE_VERSION}` from running `tanzu package
   > available list cert-manager.community.tanzu.vmware.com`. Specifying a
   > namespace may be required depending on where your package repository was
   > installed.

1. Install Contour Package

   If your workload cluster supports Service type LoadBalancer, simply execute this command:

   ```shell
   tanzu package install contour \
      --package-name contour.community.tanzu.vmware.com \
      --version ${CONTOUR_PACKAGE_VERSION}
   ```

   > You can get the `${CONTOUR_PACKAGE_VERSION}` from running `tanzu package
   > available list contour.community.tanzu.vmware.com`. Specifying a
   > namespace may be required depending on where your package repository was
   > installed.

   If your workload cluster doesn't support Service type LoadBalancer, use NodePort with hostPorts enabled instead by following these steps:

   * Get the configuration for this package, by heading to [TCE GitHub repository](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages).  Select the package/version and navigate into the `bundle/config` directory. Download or copy/paste the `values.yaml` file. Rename it `contour-values.yaml`.
   * Set `envoy.service.type: NodePort` and `envoy.hostPorts.enable: true` in `contour-values.yaml`
   * Run `tanzu package install contour --package-name contour.community.tanzu.vmware.com --version ${CONTOUR_PACKAGE_VERSION} --values-file contour-values.yaml`

1. Configure Harbor Package

   Configure with the `harbor-values.yaml` file you obtained before.

   Optionally get the helper script for configuring Harbor:

   ```shell
   image_url=$(kubectl get package harbor.community.tanzu.vmware.com.2.3.3 -o jsonpath='{.spec.template.spec.fetch[0].imgpkgBundle.image}')
   imgpkg pull -b $image_url -o /tmp/harbor-package
   cp /tmp/harbor-package/config/scripts/generate-passwords.sh .
   ```

   Specify the mandatory passwords and secrets in `harbor-values.yaml`,

   or

   to Generate them automatically. run

   ```shell
   bash generate-passwords.sh harbor-values.yaml
   ```

   This step is needed only once.

   Specify other Harbor configuration (e.g. admin password, hostname, persistence setting, etc.) in `harbor-values.yaml`.

   **NOTE**: If the default storageClass in the Workload Cluster, or the specified storageClass in `harbor-values.yaml` supports the accessMode [ReadWriteMany](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes), make sure to update the accessMode from `ReadWriteOnce` to `ReadWriteMany` in `harbor-values.yaml`. [VMware vSphere 7 with vSAN 7 File Service enabled supports accessMode ReadWriteMany](https://blogs.vmware.com/virtualblocks/2020/03/12/cloud-native-storage-and-vsan-file-services-integration/) but vSphere 6.7u3 does not. If you are using vSphere 7 without vSAN File Service enabled, or you are using vSphere 6.7u3, use the default accessMode `ReadWriteOnce`.

1. Remove all the comments in the `harbor-values.yaml` file using tool [yq](https://mikefarah.gitbook.io/yq/) before installation. run

   ```shell
   yq -i eval '... comments=""' harbor-values.yaml
   ```

1. Install Harbor Package

   ```shell
   tanzu package install harbor \
      --package-name harbor.community.tanzu.vmware.com \
      --version 2.3.3 \
      --values-file harbor-values.yaml
   ```

   > You can get the `${HARBOR_PACKAGE_VERSION}` from running `tanzu package
   > available list harbor.community.tanzu.vmware.com`. Specifying a namespace may be required
   > depending on where your package repository was installed.

## Usage

### Connect to Harbor User Interface

The Harbor UI is exposed via the Envoy service load balancer that is running in the Contour Package. To allow users to connect to the Harbor UI, you must map the address of the Envoy service load balancer to the hostname of the Harbor service, for example `harbor.yourdomain.com`.

1. Obtain the address of the Envoy service load balancer.

   ```shell
   kubectl get svc envoy -n tanzu-system-ingress -o jsonpath='{.status.loadBalancer.ingress[0]}'
   ```

   On **vSphere without NSX Advanced Load Balancer (ALB)**, the Envoy service is exposed via NodePort instead of LoadBalancer, so the above output will be empty, and you can use the IP address of any worker node in the workload cluster instead. On **Amazon EC2**, it has a FQDN similar to `a82ebae93a6fe42cd66d9e145e4fb292-1299077984.us-west-2.elb.amazonaws.com`.
   On **vSphere with NSX ALB** and **Azure**, the Envoy service has a Load Balancer IP address similar to `20.54.226.44`.

1. Map the address of the Envoy service load balancer to the hostname of the Harbor service.

   * **vSphere**: If you deployed Harbor on a workload cluster that is running on vSphere, you must add an IP to hostname mapping in `/etc/hosts` or add corresponding `A` records in your DNS server. For example, if the IP address is `10.93.9.100`, add the following to `/etc/hosts`:

       ```shell
       10.93.9.100 harbor.yourdomain.com notary.harbor.yourdomain.com
       ```

     On Windows machines, the equivalent to `/etc/hosts/` is `C:\Windows\System32\Drivers\etc\hosts`.

   * **Amazon EC2 or Azure**: If you deployed Harbor on a workload cluster that is running on Amazon EC2 or Azure, you must create two DNS `CNAME` records (on Amazon EC2) or two DNS `A` records (on Azure) for the Harbor hostnames on a DNS server on the Internet.
      * One record for the Harbor hostname, for example, `harbor.yourdomain.com`, that you configured in `harbor-values.yaml`, that points to the FQDN or IP of the Envoy service load balancer.
      * Another record for the Notary service that is running in Harbor, for example, `notary.harbor.yourdomain.com`, that points to the FQDN or IP of the Envoy service load balancer.

Users can now connect to the Harbor UI by navigating to `https://harbor.yourdomain.com` in a Web browser and log in as user `admin` with the `harborAdminPassword` that you configured in `harbor-values.yaml`.

### Push and Pull Images to and from Harbor

1. If Harbor uses a self-signed certificate, download the Harbor CA certificate from `https://harbor.yourdomain.com/api/v2.0/systeminfo/getcert`, and install it on your local machine, so Docker can trust this CA certificate.

   * On Linux, save the certificate as `/etc/docker/certs.d/harbor.yourdomain.com/ca.crt`.
   * On macOS, follow [this procedure](https://blog.container-solutions.com/adding-self-signed-registry-certs-docker-mac).
   * On Windows, right-click the certificate file and select **Install Certificate**.

1. Log in to the Harbor registry with the user `admin`. When prompted, enter the `harborAdminPassword` that you set when you deployed the Harbor Extension on the workload cluster.

   ```shell
   docker login harbor.yourdomain.com -u admin
   ```

1. Tag an existing image that you have already pulled locally, for example `nginx:1.7.9`.

   ```shell
   docker tag nginx:1.7.9 harbor.yourdomain.com/library/nginx:1.7.9
   ```

1. Push the image to the Harbor registry.

   ```shell
   docker push harbor.yourdomain.com/library/nginx:1.7.9
   ```

1. Now you can pull the image from the Harbor registry on any machine where the Harbor CA certificate is installed.

   ```shell
   docker pull harbor.yourdomain.com/library/nginx:1.7.9
   ```
