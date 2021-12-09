# Harbor

[Harbor](https://github.com/goharbor/harbor) is an open source trusted cloud native registry project that stores, signs, and scans content. Harbor extends the open source Docker Distribution by adding the functionalities usually required by users such as security, identity and management.

## Supported Providers

The following table shows the providers this package can work with.
| AWS | Azure | vSphere | Docker |
|:---:|:-----:|:-------:|:------:|
| ✅  |    ✅  | ✅      |   ✅    |

## Components

This Harbor Package integrates open source Harbor 2.4.0. See [docs for Harbor 2.4.0](https://goharbor.io/docs/2.4.0/install-config/#harbor-components).

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

   Optionally get the `harbor-values.yaml` file to configure harbor.
   Download the values.yaml file from [addons/packages/harbor/2.4.0/bundle/config/values.yaml](https://github.com/vmware-tanzu/community-edition/blob/main/addons/packages/harbor/2.4.0/bundle/config/values.yaml) to check all configuration values for Harbor Package and rename it to `harbor-values.yaml`.

   Or get the template configuration file by using script below:

   ```shell
   image_url=$(kubectl get packages harbor.community.tanzu.vmware.com.2.4.0 -o jsonpath='{.spec.template.spec.fetch[0].imgpkgBundle.image}')
   imgpkg pull -b $image_url -o /tmp/harbor-package-PACKAGE-VERSION
   cp /tmp/harbor-package-PACKAGE-VERSION/config/values.yaml harbor-values.yaml
   ```

   > When you are using `imgpkg` to get the configuratuion file, specifying a namespace may be required
   > depending on where your package repository was installed.

   Optionally get the helper script for configuring Harbor:

   ```shell
   image_url=$(kubectl get package harbor.community.tanzu.vmware.com.2.4.0 -o jsonpath='{.spec.template.spec.fetch[0].imgpkgBundle.image}')
   imgpkg pull -b $image_url -o /tmp/harbor-package
   cp /tmp/harbor-package/config/scripts/generate-passwords.sh .
   ```

   Specify the mandatory passwords and secrets in `harbor-values.yaml`,

   or to generate them automatically. run

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
      --version 2.4.0 \
      --values-file harbor-values.yaml
   ```

   > Specifying a namespace may be required depending on where your package repository was installed.

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

## Configuration

The following lightwight pass-through values can be set to customize the Harbor installation.

### Global

| Value | Required/Optional | Default | Description |
|:-------|:-------------------|:---------|:-------------|
| `namespace` | Optional | harbor | The namespace in which to deploy Harbor.|

### General Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` harborAdminPassword `| The initial password of Harbor admin. |             | string |
|` secretKey `| The secret key used for encryption. Must be a string of 16 chars. |             | string |
|` hostname `| The FQDN for accessing Harbor admin UI and Registry service. | harbor.yourdomain.com | string |
|` logLevel `| The log level of core, exporter, jobservice, registry. | info | string |
|` port.https `| The network port of the Envoy service in Contour or other Ingress Controller. | 443 | integer |
|` pspNames `| The PSP names used by Harbor pods. The names are separated by ','. 'null' means all PSP can be used. | null | string |
|` enableContourHttpProxy `| Use contour http proxy instead of the ingress when it's true. | true | boolean |
|` network.ipFamilies `| THe array of network ipFamilies. | [IPv4 IPv6] | array |

### Proxy Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` proxy.noProxy `| Ignore proxy for the domains. | 127.0.0.1,localhost,.local,.internal | string |
|` proxy.httpProxy `| HTTP proxy URL. |  | string |
|` proxy.httpsProxy `| HTTPS proxy URL. |  | string |

### Registry Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` registry.replicas `| The replicas for the registry component. | 1 | integer |
|` registry.secret `| Secret is used to secure the upload state from client and registry storage backend. |             | string |

### Core Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` core.replicas `| The replicas for the core component. | 1 | integer |
|` core.secret `| Secret is used when core server communicates with other components. |             | string |
|` core.xsrfKey `| The XSRF key. Must be a string of 32 chars. |             | string |

### Metrics Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` metrics.jobservice.path `| The path of the metrics. | /metrics | string |
|` metrics.jobservice.port `| The port of the metrics. | 8001 | integer |
|` metrics.registry.port `| The port of the metrics. | 8001 | integer |
|` metrics.registry.path `| The path of the metrics. | /metrics | string |
|` metrics.core.path `| The path of the metrics. | /metrics | string |
|` metrics.core.port `| The port of the metrics. | 8001 | integer |
|` metrics.enabled `| Enable the metrics when it's true | false | boolean |
|` metrics.exporter.path `| The path of the metrics. | /metrics | string |
|` metrics.exporter.port `| The port of the metrics. | 8001 | integer |

### Database Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` database.password `| The initial password of the postgres database. |             | string |
|` database.shmSizeLimit `| The initial value of shmSizeLimit |             | integer |
|` database.maxIdleConns `| The initial value of maxIdleConns |             | integer |
|` database.maxOpenConns `| The initial value of maxOpenConns |             | integer |

### JobService Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` jobservice.replicas `| The replicas for the jobservice component. | 1 | integer |
|` jobservice.secret `| Secret is used when job service communicates with other components. |             | string |

### Notary Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` notary.enabled `| Whether to install Notary | true | boolean |

### Exporter Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` exporter.cacheDuration `| The initial value of cacheDuration. |             | integer |

### tlsCertificate Settings

| Value | Description | Default | Type |
|-------|-------------------|---------|-------------|
|` tlsCertificate.ca.crt `| The certificate of CA, this enables the download, link on portal to download the certificate of CA. Note that ca.crt is a key and not nested. |             | string |
|` tlsCertificate.tls.crt `| The certificate. Note that tls.crt is a key and not nested. |             | string |
|` tlsCertificate.tls.key `| The private key. Note that tls.key is a key and not nested. |             | string |

### Trivy Settings

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` trivy.enabled `| Whether to install Trivy scanner. | true | boolean |
|` trivy.gitHubToken `| the GitHub access token to download Trivy DB. |  | string |
|` trivy.replicas `| The replicas for the trivy component. | 1 | integer |
|` trivy.skipUpdate `| The flag to disable Trivy DB downloads from GitHub. | false | boolean |  

### Storage Settings

#### *General*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.imageChartStorage.type `| Specify the type of storage: "filesystem", "azure", "gcs", "s3","swift", "oss" and fill the information needed in the corresponding section. The type must be "filesystem" if you want to use persistent volumes for registry and chartmuseum | filesystem | string |
|` persistence.imageChartStorage.disableredirect `| Specify whether to disable `redirect` for images and chart storage, for backends which not supported it (such as using minio for `s3` storage type), please disable it. To disable redirects, simply set `disableredirect` to `true` instead. Refer to <https://github.com/docker/distribution/blob/master/docs/configuration.md#redirect> for the detail. | false | boolean |
|` persistence.imageChartStorage.caBundleSecretName `| Specify the "caBundleSecretName" if the storage service uses a self-signed certificate. The secret must contain keys named "ca.crt" which will be injected into the trust store of registry's and chartmuseum's containers. |  | string |

#### *FileSystem*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.imageChartStorage.filesystem.rootdirectory `| The root directory in filesystem. | /storage | string |
|` persistence.imageChartStorage.filesystem.maxthreads `| Max threads for filesystem. | 100 | integer |

#### *Azure*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.imageChartStorage.azure.accountkey `| Account key of azure storage. | base64encodedaccountkey | string |
|` persistence.imageChartStorage.azure.accountname `| Account name of azure storage. | accountname | string |
|` persistence.imageChartStorage.azure.container `| Container name of azure storage. | containername | string |
|` persistence.imageChartStorage.azure.realm `| Realm for azure storage. | core.windows.net | string |

#### *OSS*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.imageChartStorage.oss.chunksize `| Chunk size for the oss, eg 10M. |             | string |
|` persistence.imageChartStorage.oss.endpoint `| Endpoint of oss. |             | string |
|` persistence.imageChartStorage.oss.internal `| Use the internal endpoint when it's true. |             | boolean |
|` persistence.imageChartStorage.oss.rootdirectory `| The rootdirectory in oss. |             | string |
|` persistence.imageChartStorage.oss.bucket `| Bucket name of oss. | bucketname | string |
|` persistence.imageChartStorage.oss.accesskeysecret `| Access key secert of oss. | accesskeysecret | string |
|` persistence.imageChartStorage.oss.encrypt `| Encrypt of oss. |             | boolean |
|` persistence.imageChartStorage.oss.region `| Region of oss. | regionname | string |
|` persistence.imageChartStorage.oss.secure `| Secure of oss. |             | boolean |
|` persistence.imageChartStorage.oss.accesskeyid `| Access key id of oss. | accesskeyid | string |

#### *S3*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.imageChartStorage.s3.encrypt `| Encrypt for s3. | false | boolean |
|` persistence.imageChartStorage.s3.regionendpoint `| Region endpoint of s3, eg <http://myobjects.local> |             | string |
|` persistence.imageChartStorage.s3.secretkey `| Secret key of s3. |             | string |
|` persistence.imageChartStorage.s3.skipverify `| skipverify for s3. | false | boolean |
|` persistence.imageChartStorage.s3.v4auth `| Use v4auth for s3 when it's true. | true | boolean |
|` persistence.imageChartStorage.s3.chunksize `| Check size for s3. |             | integer |
|` persistence.imageChartStorage.s3.multipartcopychunksize `| multi part copy chunk size of s3. |             | integer |
|` persistence.imageChartStorage.s3.multipartcopythresholdsize `| multi part copy threshold size of s3. |             | integer |
|` persistence.imageChartStorage.s3.secure `| Secure for s3. | true | boolean |
|` persistence.imageChartStorage.s3.bucket `| Bucket name of s3. | bucketname | string |
|` persistence.imageChartStorage.s3.multipartcopymaxconcurrency `| multi part copy max concurrency of s3. |             | integer |
|` persistence.imageChartStorage.s3.rootdirectory `| The rootdirectory in s3. |             | string |
|` persistence.imageChartStorage.s3.storageclass `| Storage class of s3. | STANDARD | string |
|` persistence.imageChartStorage.s3.accesskey `| Access key of s3. |             | string |
|` persistence.imageChartStorage.s3.region `| Region of s3. | us-west-1 | string |
|` persistence.imageChartStorage.s3.keyid `| Keyid of s3. |             | string |

#### *Swift*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.imageChartStorage.swift.container `| Container of swift. | containername | string |
|` persistence.imageChartStorage.swift.domain `| Domain of swift. |             | string |
|` persistence.imageChartStorage.swift.endpointtype `| Endpoint type of swift, eg public. |             | string |
|` persistence.imageChartStorage.swift.insecureskipverify `| Ignore the cert verify when it's true. |             | boolean |
|` persistence.imageChartStorage.swift.region `| Region of swift. |             | string |
|` persistence.imageChartStorage.swift.tenant `| Tenant of swift. |             | string |
|` persistence.imageChartStorage.swift.authversion `| Auth version of swift. |             | string |
|` persistence.imageChartStorage.swift.chunksize `| Check size of swift, eg 5M. |             | string |
|` persistence.imageChartStorage.swift.tenantid `| Tenant id of swift. |             | string |
|` persistence.imageChartStorage.swift.accesskey `| Access key of swift. |             | string |
|` persistence.imageChartStorage.swift.domainid `| Domain id of swift. |             | string |
|` persistence.imageChartStorage.swift.tempurlcontainerkey `| Use temp url container key of swift when it's true. |             | boolean |
|` persistence.imageChartStorage.swift.prefix `| Prefix path of swift. |             | string |
|` persistence.imageChartStorage.swift.secretkey `| Secret key of swift. |             | string |
|` persistence.imageChartStorage.swift.tempurlmethods `| Temp url methods of swift. |             | string |
|` persistence.imageChartStorage.swift.trustid `| Trust id of swift. |             | string |
|` persistence.imageChartStorage.swift.username `| Username of swift. | username | string |
|` persistence.imageChartStorage.swift.authurl `| Auth url of swift. | <https://storage.myprovider.com/v3/auth> | string |
|` persistence.imageChartStorage.swift.password `| Password of swift. | password | string |

#### *GCS*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.imageChartStorage.gcs.bucket `| Bucket name of gcs. | bucketname | string |
|` persistence.imageChartStorage.gcs.chunksize `| Check size for gcs. | 5.24288e+06 | integer |
|` persistence.imageChartStorage.gcs.encodedkey `| The base64 encoded json file which contains the key | base64-encoded-json-key-file | string |
|` persistence.imageChartStorage.gcs.rootdirectory `| The rootdirectory in gcs. |             | string |

### Persistent Volume Claim Settings

#### *Database*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.persistentVolumeClaim.database.storageClass `| Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |  | string |
|` persistence.persistentVolumeClaim.database.subPath `| The "subPath" if the PVC is shared with other components. |  | string |
|` persistence.persistentVolumeClaim.database.accessMode `| Access mode of the PVC. | ReadWriteOnce | string |
|` persistence.persistentVolumeClaim.database.existingClaim `| Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |  | string |
|` persistence.persistentVolumeClaim.database.size `| Size of the PVC. | 1Gi | string |

#### *JobService*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.persistentVolumeClaim.jobservice.subPath `| The "subPath" if the PVC is shared with other components. |  | string |
|` persistence.persistentVolumeClaim.jobservice.accessMode `| Access mode of the PVC. | ReadWriteOnce | string |
|` persistence.persistentVolumeClaim.jobservice.existingClaim `| Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |  | string |
|` persistence.persistentVolumeClaim.jobservice.size `| Size of the PVC. | 1Gi | string |
|` persistence.persistentVolumeClaim.jobservice.storageClass `| Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |  | string |

#### *Redis*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.persistentVolumeClaim.redis.accessMode `| Access mode of the PVC. | ReadWriteOnce | string |
|` persistence.persistentVolumeClaim.redis.existingClaim `| Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |  | string |
|` persistence.persistentVolumeClaim.redis.size `| Size of the PVC. | 1Gi | string |
|` persistence.persistentVolumeClaim.redis.storageClass `| Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |  | string |
|` persistence.persistentVolumeClaim.redis.subPath `| The "subPath" if the PVC is shared with other components. |  | string |

#### *Registry*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.persistentVolumeClaim.registry.accessMode `| Access mode of the PVC. | ReadWriteOnce | string |
|` persistence.persistentVolumeClaim.registry.existingClaim `| Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |  | string |
|` persistence.persistentVolumeClaim.registry.size `| Size of the PVC. | 10Gi | string |
|` persistence.persistentVolumeClaim.registry.storageClass `| Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |  | string |
|` persistence.persistentVolumeClaim.registry.subPath `| The "subPath" if the PVC is shared with other components. |  | string |

#### *Trivy*

|Values | Description | Default | Type |
|-------|-------------|---------|------|
|` persistence.persistentVolumeClaim.trivy.subPath `| The "subPath" if the PVC is shared with other components. |  | string |
|` persistence.persistentVolumeClaim.trivy.accessMode `| Access mode of the PVC. | ReadWriteOnce | string |
|` persistence.persistentVolumeClaim.trivy.existingClaim `| Use the existing PVC which must be created manually before bound, and specify the "subPath" if the PVC is shared with other components |  | string |
|` persistence.persistentVolumeClaim.trivy.size `| Size of the PVC. | 5Gi | string |
|` persistence.persistentVolumeClaim.trivy.storageClass `| Specify the "storageClass" used to provision the volume. Or the default StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning |  | string |
