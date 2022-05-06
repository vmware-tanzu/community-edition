# Deploying Grafana + Prometheus + Contour + Local Path Storage + Cert Manager on Tanzu Community Edition

The purpose of this document is to guide the reader through the deployment of a monitoring stack using the community packages that are available with Tanzu Community Edition. These packages are Contour, Cert Manager, local-path-storage, Prometheus and Grafana. Cert Manager provides secure communication between Contour and Envoy.  Contour [projectcontour.io](https://projectcontour.io) is a control plane for an Envoy Ingress controller. Local-path-storage allows Kubernetes to consume local storage for persistent volumes. Prometheus records real-time metrics in a time series database, and Grafana, an analytics and interactive visualization web application which provides charts, graphs, and alerts when connected to a supported data source, such as Prometheus.

From a dependency perspective, Prometheus and Grafana need an Ingress, or a HTTPProxy to be more precise, which is included in the Contour package. The ingress controller is Envoy, with Contour acting as the control plane to provide dynamic configuration updates and delegation control. In this deployment, Contour will have a dependency on a Certificate Manager, which is also provided by the Cert Manager package.

Both Prometheus and Grafana also have a requirement to use persistent volumes (PV). To facilitate the creation of persistent volumes on local storage, the local-path-storage package is installed. Thus, the order of package deployment will be, Certificate Manager, followed by Contour, then local-path-storage, followed by Prometheus and then finally Grafana.

We will assume that a Tanzu Community Edition workload cluster is already provisioned. For more information about workload clusters, see [Deploying workload clusters](workload-clusters).
The Load Balancer services are being provided by [metallb](https://metallb.universe.tf/).
The metallb deployment can be considered a three step process if deploying via manifests, and the procedure is well documented on the [metallb](https://metallb.universe.tf/installation/#installation-by-manifest) web site.

- Create a namespace for MetalLB
- Deploy the components of MetalLB
- Create a ConfigMap which contains the list of addresses to use for VIPs / Load Balancer IP Addresses.

Deployment of the Tanzu Community Edition cluster and metallb are beyond the scope of this document.

It is also recommended that readers familiarise themselves with the [working with packages](/docs/package-management) documentation as we will be using packages extensively in this procedure.

## Add the Tanzu Community Edition Package Repository

By default, only the `tanzu core` packages are available on the Tanzu cluster. The package install resources for core packages are typically found in the `tkg-system` namespace. Add the `-A` option to this command to check all namespaces.

```sh
% tanzu package repository list -A
/ Retrieving repositories...
  NAME        REPOSITORY                                                                                 STATUS               DETAILS  NAMESPACE
  tanzu-core  projects-stg.registry.vmware.com/tkg/packages/core/repo:v1.21.2_vmware.1-tkg.1-zshippable  Reconcile succeeded           tkg-system


% tanzu package available list
/ Retrieving available packages...
  NAME  DISPLAY-NAME  SHORT-DESCRIPTION


% tanzu package available list -A
/ Retrieving available packages...
  NAME                                                DISPLAY-NAME                       SHORT-DESCRIPTION                                                                                                                   NAMESPACE
  addons-manager.tanzu.vmware.com                     tanzu-addons-manager               This package provides TKG addons lifecycle management capabilities.                                                                 tkg-system
  ako-operator.tanzu.vmware.com                       ako-operator                       NSX Advanced Load Balancer using ako-operator                                                                                       tkg-system
  antrea.tanzu.vmware.com                             antrea                             networking and network security solution for containers                                                                             tkg-system
  calico.tanzu.vmware.com                             calico                             Networking and network security solution for containers.                                                                            tkg-system
  kapp-controller.tanzu.vmware.com                    kapp-controller                    Kubernetes package manager                                                                                                          tkg-system
  load-balancer-and-ingress-service.tanzu.vmware.com  load-balancer-and-ingress-service  Provides L4+L7 load balancing for TKG clusters running on vSphere                                                                   tkg-system
  metrics-server.tanzu.vmware.com                     metrics-server                     Metrics Server is a scalable, efficient source of container resource metrics for Kubernetes built-in autoscaling pipelines.         tkg-system
  pinniped.tanzu.vmware.com                           pinniped                           Pinniped provides identity services to Kubernetes.                                                                                  tkg-system
  vsphere-cpi.tanzu.vmware.com                        vsphere-cpi                        vSphere CPI provider                                                                                                                tkg-system
  vsphere-csi.tanzu.vmware.com                        vsphere-csi                        vSphere CSI provider
```

To access the community packages, you will first need to add the `tce` repository.

```sh
% tanzu package repository add tce-repo --url projects.registry.vmware.com/tce/main:{{< pkg_repo_latest >}}
/ Adding package repository 'tce-repo'...
Added package repository 'tce-repo'
 ```

Monitor the repo until the STATUS changes to `Reconcile succeeded`. The community packages are now available to the cluster.

```sh
% tanzu package repository list -A
/ Retrieving repositories...
  NAME        REPOSITORY                                                                                 STATUS               DETAILS  NAMESPACE
  tce-repo    projects.registry.vmware.com/tce/main:{{< pkg_repo_latest >}}                                               Reconciling                   default
  tanzu-core  projects-stg.registry.vmware.com/tkg/packages/core/repo:v1.21.2_vmware.1-tkg.1-zshippable  Reconcile succeeded           tkg-system

% tanzu package repository list -A
/ Retrieving repositories...
  NAME        REPOSITORY                                                                                 STATUS               DETAILS  NAMESPACE
  tce-repo    projects.registry.vmware.com/tce/main:{{< pkg_repo_latest >}}                                               Reconcile succeeded           default
  tanzu-core  projects-stg.registry.vmware.com/tkg/packages/core/repo:v1.21.2_vmware.1-tkg.1-zshippable  Reconcile succeeded           tkg-system
  ```

Additional packages from the newly added repository should now be available.

```sh
% tanzu package available list -A
\ Retrieving available packages...
  NAME                                                DISPLAY-NAME                       SHORT-DESCRIPTION                                                                                                                  NAMESPACE
  cert-manager.community.tanzu.vmware.com             cert-manager                       Certificate management                                                                                                             default
  contour.community.tanzu.vmware.com                  Contour                            An ingress controller                                                                                                              default
  external-dns.community.tanzu.vmware.com             external-dns                       This package provides DNS synchronization functionality.                                                                           default
  fluent-bit.community.tanzu.vmware.com               fluent-bit                         Fluent Bit is a fast Log Processor and Forwarder                                                                                   default
  gatekeeper.community.tanzu.vmware.com               gatekeeper                         policy management                                                                                                                  default
  grafana.community.tanzu.vmware.com                  grafana                            Visualization and analytics software                                                                                               default
  harbor.community.tanzu.vmware.com                   Harbor                             OCI Registry                                                                                                                       default
  knative-serving.community.tanzu.vmware.com          knative-serving                    Knative Serving builds on Kubernetes to support deploying and serving of applications and functions as serverless containers       default
  local-path-storage.community.tanzu.vmware.com       local-path-storage                 This package provides local path node storage and primarily supports RWO AccessMode.                                               default
  multus-cni.community.tanzu.vmware.com               multus-cni                         This package provides the ability for enabling attaching multiple network interfaces to pods in Kubernetes                         default
  prometheus.community.tanzu.vmware.com               prometheus                         A time series database for your metrics                                                                                            default
  velero.community.tanzu.vmware.com                   velero                             Disaster recovery capabilities                                                                                                     default
  addons-manager.tanzu.vmware.com                     tanzu-addons-manager               This package provides TKG addons lifecycle management capabilities.                                                                tkg-system
  ako-operator.tanzu.vmware.com                       ako-operator                       NSX Advanced Load Balancer using ako-operator                                                                                      tkg-system
  antrea.tanzu.vmware.com                             antrea                             networking and network security solution for containers                                                                            tkg-system
  calico.tanzu.vmware.com                             calico                             Networking and network security solution for containers.                                                                           tkg-system
  kapp-controller.tanzu.vmware.com                    kapp-controller                    Kubernetes package manager                                                                                                         tkg-system
  load-balancer-and-ingress-service.tanzu.vmware.com  load-balancer-and-ingress-service  Provides L4+L7 load balancing for TKG clusters running on vSphere                                                                  tkg-system
  metrics-server.tanzu.vmware.com                     metrics-server                     Metrics Server is a scalable, efficient source of container resource metrics for Kubernetes built-in autoscaling pipelines.        tkg-system
  pinniped.tanzu.vmware.com                           pinniped                           Pinniped provides identity services to Kubernetes.                                                                                 tkg-system
  vsphere-cpi.tanzu.vmware.com                        vsphere-cpi                        vSphere CPI provider                                                                                                               tkg-system
  vsphere-csi.tanzu.vmware.com                        vsphere-csi                        vSphere CSI provider                                                                                                               tkg-system
```

## Deploy Certificate Manager

Cert-manager [cert-manager.io](http://cert-manager.io) is an optional package, but shall be installed to make the monitoring app stack more secure. We will use it to secure communications between Contour and the Envoy Ingress. Thus, Contour has a dependency on Certificate Manager, so we will need to install this package first.

Cert-manager automates certificate management in cloud native environments. It provides certificates-as-a-service capabilities. You can install the cert-manager package on your cluster through a community package.

In this example, version 1.5.1 of the Cert Manager is being deployed with its default configuration values. Other versions of the package are available and can also be used should there be a need to do so. To check which versions of a package are available, the `tanzu package available list` command is used:

```sh
% tanzu package available list cert-manager.community.tanzu.vmware.com -n default
\ Retrieving package versions for cert-manager.community.tanzu.vmware.com...
  NAME                                     VERSION  RELEASED-AT
  cert-manager.community.tanzu.vmware.com  1.3.1    2021-04-14T18:00:00Z
  cert-manager.community.tanzu.vmware.com  1.4.0    2021-06-15T18:00:00Z
  cert-manager.community.tanzu.vmware.com  1.5.1    2021-08-13T19:52:11Z
```

For some packages, bespoke changes to the configuration may be required. There is no requirement to supply any bespoke data values for the Cert Manager packages, unless you would like to install the certificate manager components to a different target namespace (the target namespace is set to *tanzu-certificates* by default). These configuration values can be queried using the `--values-schema` option with the `tanzu package available get` command, as shown below.

```sh
% tanzu package available get cert-manager.community.tanzu.vmware.com/1.5.1 -n default --values-schema
| Retrieving package details for cert-manager.community.tanzu.vmware.com/1.5.1...
  KEY        DEFAULT             TYPE    DESCRIPTION
  namespace  tanzu-certificates  string  The namespace in which to deploy cert-manager.
```

Once the version has been identified, it can be installed using the following command:

```sh
% tanzu package install cert-manager --package-name cert-manager.community.tanzu.vmware.com --version 1.5.1
/ Installing package 'cert-manager.community.tanzu.vmware.com'
| Getting namespace 'default'
- Getting package metadata for 'cert-manager.community.tanzu.vmware.com'
| Creating service account 'cert-manager-default-sa'
| Creating cluster admin role 'cert-manager-default-cluster-role'
| Creating cluster role binding 'cert-manager-default-cluster-rolebinding'
- Creating package resource
/ Package install status: Reconciling

Added installed package 'cert-manager' in namespace 'default'

 %
 ```

These commands verify that the package has been installed successfully. As can be seen, the objects for cert-manager have been placed in the `tanzu-certificates` namespace by default.

 ```sh
% tanzu package installed list
- Retrieving installed packages...
  NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
  cert-manager  cert-manager.community.tanzu.vmware.com  1.5.1            Reconcile succeeded


% kubectl get pods -A | grep cert-manager
tanzu-certificates   cert-manager-66655465c7-ngq2h                            1/1     Running   0          111s
tanzu-certificates   cert-manager-cainjector-5d765746b-pfwpq                  1/1     Running   0          111s
tanzu-certificates   cert-manager-webhook-5d97fcb595-4ntwt                    1/1     Running   0          111s
```

With the Certificate Manager successfully deployed, the next step is to deploy the Contour community package.

## Deploy Contour (Ingress)

Later we shall deploy Prometheus and Grafana, which have a requirement on an Ingress/HTTPProxy. Contour [projectcontour.io](http://projectcontour.io) provides this functionality via an Envoy Ingress controller. Contour is an open source Kubernetes Ingress controller that acts as a control plane for the Envoy edge and service proxy.

For our purposes of standing up a monitoring stack, we can provide a very simple data values file, in YAML format, when deploying Contour. In this manifest, the Envoy Ingress controller is requested to use a Load Balancer service which will be provided by `metallb`. It also requests that Contour leverage the previously deployed Cert-Manager to provision TLS certificates rather than using the upstream Contour cert-gen job to provision certificates. This secures communication between Contour and Envoy.

```yaml
envoy:
  service:
    type: LoadBalancer
certificates:
  useCertManager: true
```

This is only a subset of the configuration parameters available in Contour. To display all configuration parameters, use the `--values-schema` option to display the configuration settings against the appropriate version of the package:

```sh
% tanzu package available list contour.community.tanzu.vmware.com
- Retrieving package versions for contour.community.tanzu.vmware.com...
  NAME                                VERSION  RELEASED-AT
  contour.community.tanzu.vmware.com  1.15.1   2021-06-01T18:00:00Z
  contour.community.tanzu.vmware.com  1.17.1   2021-07-23T18:00:00Z


% tanzu package available get contour.community.tanzu.vmware.com/1.17.1 --values-schema
| Retrieving package details for contour.community.tanzu.vmware.com/1.17.1...
  KEY                                  DEFAULT         TYPE     DESCRIPTION
  envoy.hostNetwork                    false           boolean  Whether to enable host networking for the Envoy pods.
  envoy.hostPorts.enable               false           boolean  Whether to enable host ports. If false, http and https are ignored.
  envoy.hostPorts.http                 80              integer  If enable == true, the host port number to expose Envoys HTTP listener on.
  envoy.hostPorts.https                443             integer  If enable == true, the host port number to expose Envoys HTTPS listener on.
  envoy.logLevel                       info            string   The Envoy log level.
  envoy.service.type                   LoadBalancer    string   The type of Kubernetes service to provision for Envoy.
  envoy.service.annotations            <nil>           object   Annotations to set on the Envoy service.
  envoy.service.externalTrafficPolicy  Local           string   The external traffic policy for the Envoy service.
  envoy.service.loadBalancerIP         <nil>           string   If type == LoadBalancer, the desired load balancer IP for the Envoy service.
  envoy.service.nodePorts.http         <nil>           integer  If type == NodePort, the node port number to expose Envoys HTTP listener on. If not specified, a node port will be auto-assigned by Kubernetes.
  envoy.service.nodePorts.https        <nil>           integer  If type == NodePort, the node port number to expose Envoys HTTPS listener on. If not specified, a node port will be auto-assigned by Kubernetes.
  envoy.terminationGracePeriodSeconds  300             integer  The termination grace period, in seconds, for the Envoy pods.
  namespace                            projectcontour  string   The namespace in which to deploy Contour and Envoy.
  certificates.renewBefore             360h            string   If using cert-manager, how long before expiration the certificates should be renewed. If useCertManager is false, this field is ignored.
  certificates.useCertManager          false           boolean  Whether to use cert-manager to provision TLS certificates for securing communication between Contour and Envoy. If false, the upstream Contour certgen job will be used to provision certificates. If true, the cert-manager addon must be installed in the cluster.
  certificates.duration                8760h           string   If using cert-manager, how long the certificates should be valid for. If useCertManager is false, this field is ignored.
  contour.configFileContents           <nil>           object   The YAML contents of the Contour config file. See https://projectcontour.io/docs/v1.17.1/configuration/#configuration-file for more information.
  contour.logLevel                     info            string   The Contour log level. Valid options are info and debug.
  contour.replicas                     2               integer  How many Contour pod replicas to have.
  contour.useProxyProtocol             false           boolean  Whether to enable PROXY protocol for all Envoy listeners.
```

Note that there is currently no mechanism at present to display the configuration parameters in `yaml` format, but YAML examples can be found in the official package documentation.

With the above YAML manifest stored in `contour-data-values.yaml`, the Contour/Envoy Ingress can now be deployed:

```sh
% tanzu package install contour -p contour.community.tanzu.vmware.com --version 1.17.1 --values-file contour-data-values.yaml
/ Installing package 'contour.community.tanzu.vmware.com'
| Getting namespace 'default'
/ Getting package metadata for 'contour.community.tanzu.vmware.com'
| Creating service account 'contour-default-sa'
| Creating cluster admin role 'contour-default-cluster-role'
| Creating cluster role binding 'contour-default-cluster-rolebinding'
| Creating secret 'contour-default-values'
\ Creating package resource
\ Package install status: Reconciling

Added installed package 'contour' in namespace 'default'

%
```

### Check Contour configuration values have taken effect

The following commands can be used to verify that the Contour configuration values provided at deployment time have been implemented.

```sh
% tanzu package installed list
| Retrieving installed packages...
  NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
  cert-manager  cert-manager.community.tanzu.vmware.com  1.5.1            Reconcile succeeded
  contour       contour.community.tanzu.vmware.com       1.17.1           Reconcile succeeded


% tanzu package installed get contour -f /tmp/yyy
\ Retrieving installation details for contour... %


% cat /tmp/yyy
---
envoy:
 service:
   type: LoadBalancer
certificates:
  useCertManager: true
%
```

### Validating Contour functionality

A good step at this point is to verify that Envoy is working as expected. To do that, we can locate the Envoy Pod, setup port-forwarding from the container, and then connect a browser to the localhost port as shown below:

```sh
% kubectl get pods -A | grep contour
projectcontour       contour-77f5ddff5d-9vkx2                                 1/1     Running   0          11m
projectcontour       contour-77f5ddff5d-kvhsd                                 1/1     Running   0          11m
projectcontour       envoy-9kmtl                                              2/2     Running   0          11m


% kubectl get svc envoy -n projectcontour
NAMESPACE              NAME                   TYPE           CLUSTER-IP       EXTERNAL-IP    PORT(S)                      AGE
projectcontour         envoy                  LoadBalancer   100.71.239.129   172.18.255.1   80:31459/TCP,443:30549/TCP   11m

% ENVOY_POD=$(kubectl -n projectcontour get pod -l app=envoy -o name | head -1)
% echo $ENVOY_POD
pod/envoy-mfjcp


% kubectl -n projectcontour port-forward $ENVOY_POD 9001:9001
Forwarding from 127.0.0.1:9001 -> 9001
Forwarding from [::1]:9001 -> 9001
Handling connection for 9001
```

Now if you point a browser to the `localhost:9001`, the following Envoy landing page should be displayed:

![Envoy Listing](../img/envoy-listings.png?raw=true)

Note the stats/prometheus link. This will be useful to reference when testing Prometheus in a later step.

## Deploy local-path-storage

Both Prometheus and Grafana have a requirement for persistent storage, so both
have Persistent Volume Claims. By default, there is no storage provider
available in Tanzu Community Edition, and thus no default Storage
Class. To accommodate this request for persistent storage, another community
package called `local-path-storage` must be deployed in the Tanzu Community
Edition cluster. Once the package has been successfully installed and
reconciled, there should be a new default StorageClass called `local-path`
created on the cluster.

Begin by installing the required version of the package. In this guide, we are installed version 0.0.20.

```sh
% tanzu package available list
- Retrieving available packages...
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


% tanzu package available list local-path-storage.community.tanzu.vmware.com
- Retrieving package versions for local-path-storage.community.tanzu.vmware.com...
  NAME                                           VERSION  RELEASED-AT
  local-path-storage.community.tanzu.vmware.com  0.0.19
  local-path-storage.community.tanzu.vmware.com  0.0.20


% tanzu package install local-path-storage -p local-path-storage.community.tanzu.vmware.com -v 0.0.20
- Installing package 'local-path-storage.community.tanzu.vmware.com'
| Getting namespace 'default'
| Getting package metadata for 'local-path-storage.community.tanzu.vmware.com'
| Creating service account 'local-path-storage-default-sa'
| Creating cluster admin role 'local-path-storage-default-cluster-role'
| Creating cluster role binding 'local-path-storage-default-cluster-rolebinding'
- Creating package resource
| Package install status: Reconciling

Added installed package 'local-path-storage' in namespace 'default'

%
```

Once the package is installed, we can check that the pod is running, and that the StorageClass has been successfully created.

```sh
% kubectl get pods -n tanzu-local-path-storage
NAME                                      READY   STATUS    RESTARTS   AGE
local-path-provisioner-6d6784d644-lxpp6   1/1     Running   0          55s


% kubectl get sc
NAME                   PROVISIONER             RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
local-path (default)   rancher.io/local-path   Delete          WaitForFirstConsumer   false                  113s
```

Everything is now in place to proceed with the deployment of the Prometheus package.

## Deploy Prometheus

Prometheus [prometheus.io](http://prometheus.io) records real-time metrics and provides alerting capabilities. It has a requirement for an Ingress (or HTTPProxy) and that the requirement has been met by the Contour package. It also has a requirement for persistent storage and that requirement has been met by local-path-storage package. We can now proceed with the installation of the Prometheus package. Prometheus has several configuration options, which can be displayed using the `--values-schema` option with the `tanzu package available get` command. In the command outputs below,  the version is first retrieved, then the configuration options are displayed for that version. At present, there is only a single version of the Prometheus community package available:

```sh
% tanzu package available list prometheus.community.tanzu.vmware.com
- Retrieving package versions for prometheus.community.tanzu.vmware.com...
  NAME                                   VERSION  RELEASED-AT
  prometheus.community.tanzu.vmware.com  2.27.0   2021-05-12T18:00:00Z


% tanzu package available get prometheus.community.tanzu.vmware.com/2.27.0 --values-schema
| Retrieving package details for prometheus.community.tanzu.vmware.com/2.27.0...
  KEY                                                         DEFAULT                                     TYPE     DESCRIPTION
  pushgateway.deployment.containers.resources                 <nil>                                       object   pushgateway containers resource requirements (See Kubernetes OpenAPI Specification io.k8s.api.core.v1.ResourceRequirements)
  pushgateway.deployment.podAnnotations                       <nil>                                       object   pushgateway deployments pod annotations
  pushgateway.deployment.podLabels                            <nil>                                       object   pushgateway deployments pod labels
  pushgateway.deployment.replicas                             1                                           integer  Number of pushgateway replicas.
  pushgateway.service.annotations                             <nil>                                       object   pushgateway service annotations
  pushgateway.service.labels                                  <nil>                                       object   pushgateway service pod labels
  pushgateway.service.port                                    9091                                        integer  The ports that are exposed by pushgateway service.
  pushgateway.service.targetPort                              9091                                        integer  Target Port to access on the pushgateway pods.
  pushgateway.service.type                                    ClusterIP                                   string   The type of Kubernetes service to provision for pushgateway.
  alertmanager.service.annotations                            <nil>                                       object   Alertmanager service annotations
  alertmanager.service.labels                                 <nil>                                       object   Alertmanager service pod labels
  alertmanager.service.port                                   80                                          integer  The ports that are exposed by Alertmanager service.
  alertmanager.service.targetPort                             9093                                        integer  Target Port to access on the Alertmanager pods.
  alertmanager.service.type                                   ClusterIP                                   string   The type of Kubernetes service to provision for Alertmanager.
  alertmanager.config.alertmanager_yml                        See default values file                     object   The contents of the Alertmanager config file. See https://prometheus.io/docs/alerting/latest/configuration/ for more information.
  alertmanager.deployment.podLabels                           <nil>                                       object   Alertmanager deployments pod labels
  alertmanager.deployment.replicas                            1                                           integer  Number of alertmanager replicas.
  alertmanager.deployment.containers.resources                <nil>                                       object   Alertmanager containers resource requirements (See Kubernetes OpenAPI Specification io.k8s.api.core.v1.ResourceRequirements)
  alertmanager.deployment.podAnnotations                      <nil>                                       object   Alertmanager deployments pod annotations
  alertmanager.pvc.accessMode                                 ReadWriteOnce                               string   The name of the AccessModes to use for persistent volume claim. By default this is null and default provisioner is used
  alertmanager.pvc.annotations                                <nil>                                       object   Alertmanagers persistent volume claim annotations
  alertmanager.pvc.storage                                    2Gi                                         string   The storage size for Alertmanager server persistent volume claim.
  alertmanager.pvc.storageClassName                           <nil>                                       string   The name of the StorageClass to use for persistent volume claim. By default this is null and default provisioner is used
  cadvisor.daemonset.podLabels                                <nil>                                       object   cadvisor deployments pod labels
  cadvisor.daemonset.updatestrategy                           RollingUpdate                               string   The type of DaemonSet update.
  cadvisor.daemonset.containers.resources                     <nil>                                       object   cadvisor containers resource requirements (See Kubernetes OpenAPI Specification io.k8s.api.core.v1.ResourceRequirements)
  cadvisor.daemonset.podAnnotations                           <nil>                                       object   cadvisor deployments pod annotations
  ingress.tlsCertificate.tls.crt                              <nil>                                       string   Optional cert for ingress if you want to use your own TLS cert. A self signed cert is generated by default. Note that tls.crt is a key and not nested.
  ingress.tlsCertificate.tls.key                              <nil>                                       string   Optional cert private key for ingress if you want to use your own TLS cert. Note that tls.key is a key and not nested.
  ingress.tlsCertificate.ca.crt                               <nil>                                       string   Optional CA certificate. Note that ca.crt is a key and not nested.
  ingress.virtual_host_fqdn                                   prometheus.system.tanzu                     string   Hostname for accessing prometheus and alertmanager.
  ingress.alertmanagerServicePort                             80                                          integer  Alertmanager service port to proxy traffic to.
  ingress.alertmanager_prefix                                 /alertmanager/                              string   Path prefix for Alertmanager.
  ingress.enabled                                             false                                       boolean  Whether to enable Prometheus and Alertmanager Ingress. Note that this requires contour.
  ingress.prometheusServicePort                               80                                          integer  Prometheus service port to proxy traffic to.
  ingress.prometheus_prefix                                   /                                           string   Path prefix for Prometheus.
  kube_state_metrics.deployment.replicas                      1                                           integer  Number of kube-state-metrics replicas.
  kube_state_metrics.deployment.containers.resources          <nil>                                       object   kube-state-metrics containers resource requirements (See Kubernetes OpenAPI Specification io.k8s.api.core.v1.ResourceRequirements)
  kube_state_metrics.deployment.podAnnotations                <nil>                                       object   kube-state-metrics deployments pod annotations
  kube_state_metrics.deployment.podLabels                     <nil>                                       object   kube-state-metrics deployments pod labels
  kube_state_metrics.service.annotations                      <nil>                                       object   kube-state-metrics service annotations
  kube_state_metrics.service.labels                           <nil>                                       object   kube-state-metrics service pod labels
  kube_state_metrics.service.port                             80                                          integer  The ports that are exposed by kube-state-metrics service.
  kube_state_metrics.service.targetPort                       8080                                        integer  Target Port to access on the kube-state-metrics pods.
  kube_state_metrics.service.telemetryPort                    81                                          integer  The ports that are exposed by kube-state-metrics service.
  kube_state_metrics.service.telemetryTargetPort              8081                                        integer  Target Port to access on the kube-state-metrics pods.
  kube_state_metrics.service.type                             ClusterIP                                   string   The type of Kubernetes service to provision for kube-state-metrics.
  namespace                                                   prometheus                                  string   The namespace in which to deploy Prometheus.
  node_exporter.daemonset.podLabels                           <nil>                                       object   node-exporter deployments pod labels
  node_exporter.daemonset.updatestrategy                      RollingUpdate                               string   The type of DaemonSet update.
  node_exporter.daemonset.containers.resources                <nil>                                       object   node-exporter containers resource requirements (See Kubernetes OpenAPI Specification io.k8s.api.core.v1.ResourceRequirements)
  node_exporter.daemonset.hostNetwork                         false                                       boolean  The Host networking requested for this pod
  node_exporter.daemonset.podAnnotations                      <nil>                                       object   node-exporter deployments pod annotations
  node_exporter.service.labels                                <nil>                                       object   node-exporter service pod labels
  node_exporter.service.port                                  9100                                        integer  The ports that are exposed by node-exporter service.
  node_exporter.service.targetPort                            9100                                        integer  Target Port to access on the node-exporter pods.
  node_exporter.service.type                                  ClusterIP                                   string   The type of Kubernetes service to provision for node-exporter.
  node_exporter.service.annotations                           <nil>                                       object   node-exporter service annotations
  prometheus.config.alerting_rules_yml                        See default values file                     object   The YAML contents of the Prometheus alerting rules file.
  prometheus.config.alerts_yml                                <nil>                                       object   Additional prometheus alerts can be configured in this YAML file.
  prometheus.config.prometheus_yml                            See default values file                     object   The contents of the Prometheus config file. See https://prometheus.io/docs/prometheus/latest/configuration/configuration/ for more information.
  prometheus.config.recording_rules_yml                       See default values file                     object   The YAML contents of the Prometheus recording rules file.
  prometheus.config.rules_yml                                 <nil>                                       object   Additional prometheus rules can be configured in this YAML file.
  prometheus.deployment.configmapReload.containers.args       webhook-url=http://127.0.0.1:9090/-/reload  array    List of arguments passed via command-line to configmap reload container. For more guidance on configuration options consult the configmap-reload docs at https://github.com/jimmidyson/configmap-reload#usage
  prometheus.deployment.configmapReload.containers.resources  <nil>                                       object   configmap-reload containers resource requirements (io.k8s.api.core.v1.ResourceRequirements)
  prometheus.deployment.containers.args                       prometheus storage retention time = 42d     array    List of arguments passed via command-line to prometheus server. For more guidance on configuration options consult the Prometheus docs at https://prometheus.io/
  prometheus.deployment.containers.resources                  <nil>                                       object   Prometheus containers resource requirements (See Kubernetes OpenAPI Specification io.k8s.api.core.v1.ResourceRequirements)
  prometheus.deployment.podAnnotations                        <nil>                                       object   Prometheus deployments pod annotations
  prometheus.deployment.podLabels                             <nil>                                       object   Prometheus deployments pod labels
  prometheus.deployment.replicas                              1                                           integer  Number of prometheus replicas.
  prometheus.pvc.annotations                                  <nil>                                       object   Prometheuss persistent volume claim annotations
  prometheus.pvc.storage                                      150Gi                                       string   The storage size for Prometheus server persistent volume claim.
  prometheus.pvc.storageClassName                             <nil>                                       string   The name of the StorageClass to use for persistent volume claim. By default this is null and default provisioner is used
  prometheus.pvc.accessMode                                   ReadWriteOnce                               string   The name of the AccessModes to use for persistent volume claim. By default this is null and default provisioner is used
  prometheus.service.annotations                              <nil>                                       object   Prometheus service annotations
  prometheus.service.labels                                   <nil>                                       object   Prometheus service pod labels
  prometheus.service.port                                     80                                          integer  The ports that are exposed by Prometheus service.
  prometheus.service.targetPort                               9090                                        integer  Target Port to access on the Prometheus pods.
  prometheus.service.type                                     ClusterIP                                   string   The type of Kubernetes service to provision for Prometheus.
```

This displays all of the configuration settings for Prometheus. For the purposes of this deployment, a simple manifest which enables ingress and provides a virtual host fully qualified domain name are all that is needed. This is the sample manifest to modify the default Prometheus deployment, primarily to enable Ingress/HTTPProxy usage, and secondly to set a fully qualified domain name - fqdn - which would be used to access the Prometheus dashboard:

```yaml
ingress:
  enabled: true
  virtual_host_fqdn: "prometheus.local"
  prometheus_prefix: "/"
  alertmanager_prefix: "/alertmanager/"
  prometheusServicePort: 80
  alertmanagerServicePort: 80
```

While it is interesting to see how an Ingress can be configured from Prometheus, it is only for demonstration purposes in the current version of Tanzu Community Edition clusters. Due to the limitations of docker networking not being accessible directly to the host on macOS and Windows, we will not be able to use the Ingress to reach the Prometheus UI in this guide. Instead, we will be relying on port-forwarding, as we will see shortly.

We can now proceed with deploying the Prometheus package, using the `--values-file` to point to the simple manifest created previously.

```sh
% tanzu package install prometheus -p prometheus.community.tanzu.vmware.com -v 2.27.0 --values-file prometheus-data-values.yaml
- Installing package 'prometheus.community.tanzu.vmware.com'
| Getting namespace 'default'
/ Getting package metadata for 'prometheus.community.tanzu.vmware.com'
| Creating service account 'prometheus-default-sa'
| Creating cluster admin role 'prometheus-default-cluster-role'
| Creating cluster role binding 'prometheus-default-cluster-rolebinding'
| Creating secret 'prometheus-default-values'
\ Creating package resource
- Package install status: Reconciling

Added installed package 'prometheus' in namespace 'default'

%
```

During the deployment, the following helper Pods appear for the creation of the Persistent Volumes:

```sh
tanzu-local-path-storage   helper-pod-create-pvc-df5b72ae-7cee-47fb-be8d-c4e06d9c85f8   0/1     Completed           0          7s
tanzu-local-path-storage   helper-pod-create-pvc-df5b72ae-7cee-47fb-be8d-c4e06d9c85f8   0/1     Terminating         0          8s
tanzu-local-path-storage   helper-pod-create-pvc-df5b72ae-7cee-47fb-be8d-c4e06d9c85f8   0/1     Terminating         0          8s
tanzu-local-path-storage   helper-pod-create-pvc-78192f97-d7eb-440f-8d36-004a0a6a094c   1/1     Running             0          6s
tanzu-local-path-storage   helper-pod-create-pvc-78192f97-d7eb-440f-8d36-004a0a6a094c   0/1     Completed           0          8s
```

The Persistent Volume Claims (PVCs) and Persistent Volumes (PVs) then get created:

```sh
% kubectl get pvc -A
NAMESPACE    NAME                STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
prometheus   alertmanager        Bound    pvc-df5b72ae-7cee-47fb-be8d-c4e06d9c85f8   2Gi        RWO            local-path     74s
prometheus   prometheus-server   Bound    pvc-78192f97-d7eb-440f-8d36-004a0a6a094c   150Gi      RWO            local-path     74s


% kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                          STORAGECLASS   REASON   AGE
pvc-78192f97-d7eb-440f-8d36-004a0a6a094c   150Gi      RWO            Delete           Bound    prometheus/prometheus-server   local-path              68s
pvc-df5b72ae-7cee-47fb-be8d-c4e06d9c85f8   2Gi        RWO            Delete           Bound    prometheus/alertmanager        local-path              70s
```

### Check Prometheus configuration values have taken effect

The following command can be used to verify that the data values provided at deployment time have been implemented.

```sh
% tanzu package installed get prometheus -f /tmp/xxx
\ Retrieving installation details for prometheus... %

% more /tmp/xxx
---
ingress:
  enabled: true
  virtual_host_fqdn: "prometheus.local"
  prometheus_prefix: "/"
  alertmanager_prefix: "/alertmanager/"
  prometheusServicePort: 80
  alertmanagerServicePort: 80
```

The configuration values look correct. The following Kubernetes pods and services should have been created successfully.

```sh
% kubectl get pods,svc -n prometheus
NAMESPACE                  NAME                                                     READY   STATUS    RESTARTS   AGE
prometheus                 alertmanager-5965c496b4-6264w                            1/1     Running   0          105s
prometheus                 prometheus-cadvisor-xtq75                                1/1     Running   0          104s
prometheus                 prometheus-kube-state-metrics-7ff778d6cf-gmgn2           1/1     Running   0          104s
prometheus                 prometheus-node-exporter-5tzgv                           1/1     Running   0          104s
prometheus                 prometheus-node-exporter-ht4z7                           1/1     Running   0          104s
prometheus                 prometheus-pushgateway-8699dd8f74-rz7tl                  1/1     Running   0          103s
prometheus                 prometheus-server-7b9c56b66-975gj                        2/2     Running   0          103s

NAMESPACE            NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP    PORT(S)                      AGE
prometheus           alertmanager                    ClusterIP      100.68.9.137     <none>         80/TCP                       2m44s
prometheus           prometheus-kube-state-metrics   ClusterIP      None             <none>         80/TCP,81/TCP                2m44s
prometheus           prometheus-node-exporter        ClusterIP      100.71.35.129    <none>         9100/TCP                     2m43s
prometheus           prometheus-pushgateway          ClusterIP      100.64.223.3     <none>         9091/TCP                     2m46s
prometheus           prometheus-server               ClusterIP      100.71.138.73    <none>         80/TCP                       2m43s
```

Contour provides an advanced resource type called [HttpProxy](https://projectcontour.io/docs/v1.18.1/config/fundamentals/) that provides some benefits over Ingress resources. We can also examine that this resource was created successfully:

```sh
% kubectl get HTTPProxy -n prometheus
NAME                   FQDN               TLS SECRET       STATUS   STATUS DESCRIPTION
prometheus-httpproxy   prometheus.local   prometheus-tls   valid    Valid HTTPProxy
```

### Validate Prometheus functionality

Everything looks good. In most environments, you would now be able to point to the Prometheus FQDN / Ingress Load Balancer (e.g. prometheus.local), assuming it was addded to DNS. However, since docker does not expose the docker network to the host on macOS or Windows, we are going to rely on the port-forward functionality of Kubernetes where a pod port can be accessed from the host. The following command can be used to locate the port used by the prometheus-server container:

```sh
% kubectl get pods prometheus-server-7595fcbd6c-bs5t2 -n prometheus -o jsonpath='{.spec.containers[*].name}{.spec.containers[*].ports}'
prometheus-server-configmap-reload prometheus-server[{"containerPort":9090,"protocol":"TCP"}]%
```

The container port should be 9090 by default. Now simply `port-forward` this container port to your host with the following command, letting `kubectl` choose and allocate a host port, such as 52119 in this case.

```sh
% kubectl port-forward prometheus-server-7595fcbd6c-bs5t2 -n prometheus :9090
Forwarding from 127.0.0.1:52119 -> 9090
Forwarding from [::1]:52119 -> 9090
```

Now open a browser on your host, and connect to `localhost:52119` (or whatever port was chosen on your host). The Prometheus dashboard should be visible.

![Prometheus Dashboard Landing Page](../img/prometheus-standalone1.png?raw=true)

To do a very simple test, add a simple query, e.g. `prometheus_http_requests_total` and click Execute:

![Prometheus Simple Query](../img/prometheus-standalone2.png?raw=true)

To check integration between Prometheus and Envoy, another query can be executed. When the Envoy landing page was displayed earlier, there was a section called `prometheus/stats`. These can now be queried as well, since these are the metrics that Envoy is sending to Prometheus. If we return to the Envoy landing page in the browser, and click on the prometheus/stats link and examine one of these metrics, such as the `envoy_cluster_default_total_match`, and use it as a query in Prometheus (selecting Graph instead of Table this time):

![Envoy Prometheus Metric Query](../img/prometheus-standalone3.png?raw=true)

If this metric is also visible, then Prometheus is working successfully. Now let's complete the monitoring stack by provisioning Grafana, and connecting it to our Prometheus data source.

## Deploy Grafana

[Grafana](https://grafana.com/) is an analytics and interactive visualisation web application. Let's begin by displaying all of the configuring values that are available in Grafana. Once again, the package version is required to do this.

```sh
% tanzu package available list grafana.community.tanzu.vmware.com -A
\ Retrieving package versions for grafana.community.tanzu.vmware.com...
  NAME                                VERSION  RELEASED-AT           NAMESPACE
  grafana.community.tanzu.vmware.com  7.5.7    2021-05-19T18:00:00Z  default


% tanzu package available get grafana.community.tanzu.vmware.com/7.5.7 --values-schema
| Retrieving package details for grafana.community.tanzu.vmware.com/7.5.7...
  KEY                                      DEFAULT                                             TYPE     DESCRIPTION
  grafana.config.dashboardProvider_yaml    See default values file                             object   The YAML contents of the Grafana dashboard provider file. See https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards for more information.
  grafana.config.datasource_yaml           Includes default prometheus package as datasource.  object   The YAML contents of the Grafana datasource config file. See https://grafana.com/docs/grafana/latest/administration/provisioning/#example-data-source-config-file for more information.
  grafana.config.grafana_ini               See default values file                             object   The contents of the Grafana config file. See https://grafana.com/docs/grafana/latest/administration/configuration/ for more information.
  grafana.deployment.k8sSidecar            <nil>                                               object   k8s-sidecar related configuration.
  grafana.deployment.podAnnotations        <nil>                                               object   Grafana deployments pod annotations
  grafana.deployment.podLabels             <nil>                                               object   Grafana deployments pod labels
  grafana.deployment.replicas              1                                                   integer  Number of grafana replicas.
  grafana.deployment.containers.resources  <nil>                                               object   Grafana containers resource requirements (See Kubernetes OpenAPI Specification io.k8s.api.core.v1.ResourceRequirements)
  grafana.pvc.storage                      2Gi                                                 string   The storage size for persistent volume claim.
  grafana.pvc.storageClassName             <nil>                                               string   The name of the StorageClass to use for persistent volume claim. By default this is null and default provisioner is used
  grafana.pvc.accessMode                   ReadWriteOnce                                       string   The name of the AccessModes to use for persistent volume claim. By default this is null and default provisioner is used
  grafana.pvc.annotations                  <nil>                                               object   Grafanas persistent volume claim annotations
  grafana.secret.admin_user                admin                                               string   Username to access Grafana.
  grafana.secret.type                      Opaque                                              string   The Secret Type (io.k8s.api.core.v1.Secret.type)
  grafana.secret.admin_password            admin                                               string   Password to access Grafana. By default is null and grafana defaults this to "admin"
  grafana.service.labels                   <nil>                                               object   Grafana service pod labels
  grafana.service.port                     80                                                  integer  The ports that are exposed by Grafana service.
  grafana.service.targetPort               3000                                                integer  Target Port to access on the Grafana pods.
  grafana.service.type                     LoadBalancer                                        string   The type of Kubernetes service to provision for Grafana. (For vSphere set this to NodePort, For others set this to LoadBalancer)
  grafana.service.annotations              <nil>                                               object   Grafana service annotations
  ingress.servicePort                      80                                                  integer  Grafana service port to proxy traffic to.
  ingress.tlsCertificate.ca.crt            <nil>                                               string   Optional CA certificate. Note that ca.crt is a key and not nested.
  ingress.tlsCertificate.tls.crt           <nil>                                               string   Optional cert for ingress if you want to use your own TLS cert. A self signed cert is generated by default. Note that tls.crt is a key and not nested.
  ingress.tlsCertificate.tls.key           <nil>                                               string   Optional cert private key for ingress if you want to use your own TLS cert. Note that tls.key is a key and not nested.
  ingress.virtual_host_fqdn                grafana.system.tanzu                                string   Hostname for accessing grafana.
  ingress.enabled                          true                                                boolean  Whether to enable Grafana Ingress. Note that this requires contour.
  ingress.prefix                           /                                                   string   Path prefix for Grafana.
  namespace                                grafana                                             string   The namespace in which to deploy Grafana.
  ```

We will again try to keep the Grafana configuration quite simple. The Grafana service type is set to Load Balancer by default, and is preconfigured to use Prometheus as a data source. Thus, the only additional configuration required is to add a virtual host FQDN. Note that once again, this is not really relevant in this deployment since we are using port-forwarding to access the Grafana dashbaords. We are not using the FQDN due to the docker networking limitations. Here is my very simple values file for Grafana. You may want to specify a different fqdn.

```yaml
ingress:
  virtual_host_fqdn: "grafana.local"
```

We can now proceed to deploy the Grafana package with the above parameters:

```sh
% tanzu package install grafana --package-name grafana.community.tanzu.vmware.com --version 7.5.7 --values-file grafana-data-values.yaml
- Installing package 'grafana.community.tanzu.vmware.com'
| Getting namespace 'default'
/ Getting package metadata for 'grafana.community.tanzu.vmware.com'
| Creating service account 'grafana-default-sa'
| Creating cluster admin role 'grafana-default-cluster-role'
| Creating cluster role binding 'grafana-default-cluster-rolebinding'
| Creating secret 'grafana-default-values'
\ Creating package resource
| Package install status: Reconciling

Added installed package 'grafana' in namespace 'default'
%
```

### Check Grafana data values have taken effect

The following command can be used to verify that the data values provided at deployment time have been implemented.

```sh
% tanzu package installed get grafana -f /tmp/zzz
/ Retrieving installation details for grafana... %


% cat /tmp/zzz
---
ingress:
  virtual_host_fqdn: "grafana.local"
```

### Validate Grafana functionality

The following Pods, Services and HTTPProxy should have been created.

```sh
% kubectl get pods,pvc,pv,svc -n grafana
NAME                           READY   STATUS    RESTARTS   AGE
pod/grafana-594574468f-4fhfd   2/2     Running   0          8m20s

NAME                                STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistentvolumeclaim/grafana-pvc   Bound    pvc-232bb1cc-9e16-4ef2-8b43-d4bce46c79e1   2Gi        RWO            local-path     8m20s

NAME                                                        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                          STORAGECLASS   REASON   AGE
persistentvolume/pvc-232bb1cc-9e16-4ef2-8b43-d4bce46c79e1   2Gi        RWO            Delete           Bound    grafana/grafana-pvc            local-path              8m20s
persistentvolume/pvc-34f86f3a-2850-420e-97fe-d5596639dc20   2Gi        RWO            Delete           Bound    prometheus/alertmanager        local-path              17m
persistentvolume/pvc-be93c9ab-dde3-4141-afa9-803d060773a8   150Gi      RWO            Delete           Bound    prometheus/prometheus-server   local-path              17m

NAME              TYPE           CLUSTER-IP       EXTERNAL-IP    PORT(S)        AGE
service/grafana   LoadBalancer   100.67.216.246   172.18.255.2   80:31790/TCP   8m20s


% kubectl get httpproxy -A
NAMESPACE                 NAME                   FQDN                         TLS SECRET       STATUS   STATUS DESCRIPTION
tanzu-system-dashboard    grafana-httpproxy      grafana.corinternal.com      grafana-tls      valid    Valid HTTPProxy
tanzu-system-monitoring   prometheus-httpproxy   prometheus.corinternal.com   prometheus-tls   valid    Valid HTTPProxy
```

As mentioned, Grafana uses a Load Balancer service type by default, so it has been provided with its own Load Balancer IP addess by metallb. However, since this network is not accessible from the host, locate the port on the Grafana pod, port-forward it via kubectl and then connect to the Grafana dashboard. Let's find the port first, which should be 3000 by default, and then forward it from the Pod to the host:

```sh
% kubectl get pods grafana-594574468f-4fhfd -n grafana -o jsonpath='{.spec.containers[*].name}{.spec.containers[*].ports}'
grafana-sc-dashboard grafana[{"containerPort":80,"name":"service","protocol":"TCP"},{"containerPort":3000,"name":"grafana","protocol":"TCP"}]%


% kubectl port-forward grafana-594574468f-4fhfd -n grafana :3000
Forwarding from 127.0.0.1:52533 -> 3000
Forwarding from [::1]:52533 -> 3000
```

Now connect to `localhost:52533`, or whichever port was chosen by kubectl on your system, you should see the Grafana UI. The login credentials are admin/admin initially, but you will need to change the password on first login. This is the landing page:

![Grafana Landing Page](../img/grafana-landing-standalone.png?raw=true)

There is no need to add a datasource or create a dashboard - these have already been done for you.

To examine the data source, click on the icon representing datasources on the left-hand side (which looks like a cog). Here you can see the Prometheus data source is already in place:

![Grafana Data Source Prometheus](../img/grafana-data-source-standalone.png?raw=true)

Now click on the dashboards icon on the left hand side (it looks like a square of 4 smaller squares), and select `Manage` from the drop-down list. This will show the existing dashboards. There are 2 existing dashboards that have been provided; one is Kubernetes monitoring, and the other is TKG monitoring. These dashboards are based on the Kubernetes Grafana dashboards found on [GitHub](https://github.com/kubernetes-monitoring/kubernetes-mixin).

![Grafana Dashboards Manager](../img/grafana-manage-dashboards-standalone.png?raw=true)

Finally, select the TKG dashboard which is being sent metrics via the Prometheus data source. This provides an overview of the TKG cluster:

![TKG Dashboard](../img/grafana-dashboard-standalone.png?raw=true)

The full monitoring stack of Contour/Envoy Ingress, with secure communication via Cert-Manager, local storage provided by local-path-storage, alongside the Prometheus data scraper and Grafana visualization are now deployed through Tanzu Community Edition community packages onto a cluster running in Docker. Happy monitoring/analyzing.
