# Deploying Grafana + Prometheus + Contour + Cert Manager on TCE

The purpose of this guide is to guide the reader through a deployment of the monitoring packages that are available with TCE, the Tanzu Community Edition. These packages are Contour, Cert Manager, Prometheus and Grafana. Cert Manager provides secure communication between Contour and Envoy.  Contour [projectcontour.io] is a control plane for an Envoy Ingress controller. Prometheus records real-time metrics in a time series database, and Grafana, an analytics and interactive visualization web application which provides charts, graphs, and alerts when connected to a supported data source, such as Prometheus.

From a dependency perspective, Prometheus and Grafana have a dependency on an Ingress, or a HTTPProxy to be more precise, which is included by the Contour package. The ingress controller is Envoy, with Contour acting as the control plane to provide dynamic configuration updates and delegation control. Lastly, in this deployment, Contour will have a dependency on a Certificate Manager, which is also provided by the Cert Manager package. Thus, the order of package deployment will be, Certificate Manager, followed by Contour, followed by Prometheus and then finally Grafana.

We will make the assumption that a TCE workload cluster is already provisioned, and that is has been integrated with a load balancer. In this scenario, the deployment is to vSphere, and the Load Balancer services are being provided by the NSX Advanced Load Balancer (NSX ALB). Deployment of the TCE clusters and NSX ALB are beyond the scope of this document, but details on how to do these deployment operations can be found elsewhere in the official documentation.

It is also recommend that reader familiarise themselves with the [working with packages](/docs/latest/package-management.md) documention as we will be using packages extensively in this procedure.

## Examining the TCE environment

For the purposes of illustration, this is the environment that we will be using to deploy the monitoring stack. Your environment may of course be different. This environment has a TCE management cluster and a single TCE workload cluster. Context has been set to that of "admin" on the workload cluster. If Identity Management has been configured on the workload cluster, an LDAP or OID user with appropriate privileges may also be used.

```sh
% tanzu cluster list --include-management-cluster
NAME     NAMESPACE   STATUS   CONTROLPLANE  WORKERS  KUBERNETES        ROLES       PLAN
workload default     running  1/1           1/1      v1.21.2+vmware.1  <none>      dev
mgmt     tkg-system  running  1/1           1/1      v1.21.2+vmware.1  management  dev


% kubectl config get-contexts
CURRENT NAME                     CLUSTER   AUTHINFO        NAMESPACE
        mgmt-admin@mgmt          mgmt      mgmt-admin
        tanzu-cli-mgmt@mgmt      mgmt      tanzu-cli-mgmt
*       workload-admin@workload  workload  workload-admin


% kubectl get nodes -o wide
NAME                           STATUS ROLES                AGE   VERSION          INTERNAL-IP  EXTERNAL-IP  OS-IMAGE                KERNEL-VERSION  CONTAINER-RUNTIME
workload-control-plane-sjswp   Ready  control-plane,master 5d1h  v1.21.2+vmware.1 xx.xx.51.50  xx.xx.51.50  VMware Photon OS/Linux  4.19.198-1.ph3  containerd://1.4.6
workload-md-0-6555d876c9-qp6ft Ready  <none>               5d1h  v1.21.2+vmware.1 xx.xx.51.51  xx.xx.51.51  VMware Photon OS/Linux  4.19.198-1.ph3  containerd://1.4.6
```

## Add the Tanzu Community Edition Package Repository

By default, only the `tanzu core` packages are available on the workload cluster. I use the `-A` option to check all namespaces.

```sh
% tanzu package installed list -A
| Retrieving installed packages...
  NAME                               PACKAGE-NAME                                        PACKAGE-VERSION  STATUS               NAMESPACE
  antrea                             antrea.tanzu.vmware.com                                              Reconcile succeeded  tkg-system
  load-balancer-and-ingress-service  load-balancer-and-ingress-service.tanzu.vmware.com                   Reconcile succeeded  tkg-system
  metrics-server                     metrics-server.tanzu.vmware.com                                      Reconcile succeeded  tkg-system
  pinniped                           pinniped.tanzu.vmware.com                                            Reconcile succeeded  tkg-system
  vsphere-cpi                        vsphere-cpi.tanzu.vmware.com                                         Reconcile succeeded  tkg-system
  vsphere-csi                        vsphere-csi.tanzu.vmware.com                                         Reconcile succeeded  tkg-system
```

To access the community packages, you will first need to add the `tce` repository.

```sh
% tanzu package repository add tce-repo \
  --url projects.registry.vmware.com/tce/main:stable
/ Adding package repository 'tce-repo'...
 Added package repository 'tce-repo'
 ```

Monitor the repo until the STATUS changes to `Reconcile succeeded`. The community packages are now available to the cluster.

```sh
% tanzu package repository list -A
| Retrieving repositories...
  NAME        REPOSITORY                                                                                 STATUS               DETAILS  NAMESPACE
  tce-repo    projects.registry.vmware.com/tce/main:stable                                               Reconcile succeeded           default
  tanzu-core  projects-stg.registry.vmware.com/tkg/packages/core/repo:v1.21.2_vmware.1-tkg.1-zshippable  Reconcile succeeded           tkg-system
  ```

Additional packages from the TCE repository should now be available.

```sh
% tanzu package available list -A
/ Retrieving available packages...
  NAME                                                DISPLAY-NAME                       SHORT-DESCRIPTION                                                                                                                                                                                       NAMESPACE
  cert-manager.community.tanzu.vmware.com             cert-manager                       Certificate management                                                                                                                                                                                  default
  contour.community.tanzu.vmware.com                  Contour                            An ingress controller                                                                                                                                                                                   default
  external-dns.community.tanzu.vmware.com             external-dns                       This package provides DNS synchronization functionality.                                                                                                                                                default
  fluent-bit.community.tanzu.vmware.com               fluent-bit                         Fluent Bit is a fast Log Processor and Forwarder                                                                                                                                                        default
  gatekeeper.community.tanzu.vmware.com               gatekeeper                         policy management                                                                                                                                                                                       default
  grafana.community.tanzu.vmware.com                  grafana                            Visualization and analytics software                                                                                                                                                                    default
  harbor.community.tanzu.vmware.com                   Harbor                             OCI Registry                                                                                                                                                                                            default
  knative-serving.community.tanzu.vmware.com          knative-serving                    Knative Serving builds on Kubernetes to support deploying and serving of applications and functions as serverless containers                                                                            default
  local-path-storage.community.tanzu.vmware.com       local-path-storage                 This package provides local path node storage and primarily supports RWO AccessMode.                                                                                                                    default
  multus-cni.community.tanzu.vmware.com               multus-cni                         This package provides the ability for enabling attaching multiple network interfaces to pods in Kubernetes                                                                                              default
  prometheus.community.tanzu.vmware.com               prometheus                         A time series database for your metrics                                                                                                                                                                 default
  velero.community.tanzu.vmware.com                   velero                             Disaster recovery capabilities                                                                                                                                                                          default
  addons-manager.tanzu.vmware.com                     tanzu-addons-manager               This package provides TKG addons lifecycle management capabilities.                                                                                                                                     tkg-system
  ako-operator.tanzu.vmware.com                       ako-operator                       NSX Advanced Load Balancer using ako-operator                                                                                                                                                           tkg-system
  antrea.tanzu.vmware.com                             antrea                             networking and network security solution for containers                                                                                                                                                 tkg-system
  calico.tanzu.vmware.com                             calico                             Networking and network security solution for containers.                                                                                                                                                tkg-system
  kapp-controller.tanzu.vmware.com                    kapp-controller                    Kubernetes package manager                                                                                                                                                                              tkg-system
  load-balancer-and-ingress-service.tanzu.vmware.com  load-balancer-and-ingress-service  Provides L4+L7 load balancing for TKG clusters running on vSphere                                                                                                                                       tkg-system
  metrics-server.tanzu.vmware.com                     metrics-server                     Metrics Server is a scalable, efficient source of container resource metrics for Kubernetes built-in autoscaling pipelines.                                                                             tkg-system
  pinniped.tanzu.vmware.com                           pinniped                           Pinniped provides identity services to Kubernetes.                                                                                                                                                      tkg-system
  vsphere-cpi.tanzu.vmware.com                        vsphere-cpi                        The Cluster API brings declarative, Kubernetes-style APIs to cluster creation, configuration and management. Cluster API Provider for vSphere is a concrete implementation of Cluster API for vSphere.  tkg-system
  vsphere-csi.tanzu.vmware.com                        vsphere-csi                        vSphere CSI provider                                                                                                                                                                                    tkg-system
```

## Deploy Certificate Manager

Cert-manager [cert-manager.io](http://cert-manager.io) is an optional package, but we shall install it anyway to make the monitoring app stack more secure. We will use it to secure communications between Contour and the Envoy Ingress. Thus, Contour has a dependency on Certificate Manager, so we will need to install this package first.

Cert-manager automates certificate management in cloud native environments. It provides certificates-as-a-service capabilities. You can install the cert-manager package on your cluster through a community package.

For some packages, bespoke changes to the configuration may be required. There is no requirement to supply any bespoke data values for the Cert Manager. Thus, the package may be deployed with its default configuration values. In this example, version 1.5.1 of the Cert Manager is being deployed. Other versions may be available and can also be used. To check which versions of a package are available, use the `list` option:

```sh
% tanzu package available list -A
| Retrieving available packages...
  NAME                                                DISPLAY-NAME                       SHORT-DESCRIPTION                                                                                                                                                                                       NAMESPACE
  cert-manager.community.tanzu.vmware.com             cert-manager                       Certificate management                                                                                                                                                                                  default
  contour.community.tanzu.vmware.com                  Contour                            An ingress controller                                                                                                                                                                                   default
  external-dns.community.tanzu.vmware.com             external-dns                       This package provides DNS synchronization functionality.                                                                                                                                                default
  fluent-bit.community.tanzu.vmware.com               fluent-bit                         Fluent Bit is a fast Log Processor and Forwarder                                                                                                                                                        default
  gatekeeper.community.tanzu.vmware.com               gatekeeper                         policy management                                                                                                                                                                                       default
  grafana.community.tanzu.vmware.com                  grafana                            Visualization and analytics software                                                                                                                                                                    default
  harbor.community.tanzu.vmware.com                   Harbor                             OCI Registry                                                                                                                                                                                            default
  knative-serving.community.tanzu.vmware.com          knative-serving                    Knative Serving builds on Kubernetes to support deploying and serving of applications and functions as serverless containers                                                                            default
  local-path-storage.community.tanzu.vmware.com       local-path-storage                 This package provides local path node storage and primarily supports RWO AccessMode.                                                                                                                    default
  multus-cni.community.tanzu.vmware.com               multus-cni                         This package provides the ability for enabling attaching multiple network interfaces to pods in Kubernetes                                                                                              default
  prometheus.community.tanzu.vmware.com               prometheus                         A time series database for your metrics                                                                                                                                                                 default
  velero.community.tanzu.vmware.com                   velero                             Disaster recovery capabilities                                                                                                                                                                          default
  addons-manager.tanzu.vmware.com                     tanzu-addons-manager               This package provides TKG addons lifecycle management capabilities.                                                                                                                                     tkg-system
  ako-operator.tanzu.vmware.com                       ako-operator                       NSX Advanced Load Balancer using ako-operator                                                                                                                                                           tkg-system
  antrea.tanzu.vmware.com                             antrea                             networking and network security solution for containers                                                                                                                                                 tkg-system
  calico.tanzu.vmware.com                             calico                             Networking and network security solution for containers.                                                                                                                                                tkg-system
  kapp-controller.tanzu.vmware.com                    kapp-controller                    Kubernetes package manager                                                                                                                                                                              tkg-system
  load-balancer-and-ingress-service.tanzu.vmware.com  load-balancer-and-ingress-service  Provides L4+L7 load balancing for TKG clusters running on vSphere                                                                                                                                       tkg-system
  metrics-server.tanzu.vmware.com                     metrics-server                     Metrics Server is a scalable, efficient source of container resource metrics for Kubernetes built-in autoscaling pipelines.                                                                             tkg-system
  pinniped.tanzu.vmware.com                           pinniped                           Pinniped provides identity services to Kubernetes.                                                                                                                                                      tkg-system
  vsphere-cpi.tanzu.vmware.com                        vsphere-cpi                        The Cluster API brings declarative, Kubernetes-style APIs to cluster creation, configuration and management. Cluster API Provider for vSphere is a concrete implementation of Cluster API for vSphere.  tkg-system
  vsphere-csi.tanzu.vmware.com                        vsphere-csi                        vSphere CSI provider                                                                                                                                                                                    tkg-system


% tanzu package available get cert-manager.community.tanzu.vmware.com -n default
| Retrieving package details for cert-manager.community.tanzu.vmware.com...
NAME:                 cert-manager.community.tanzu.vmware.com
DISPLAY-NAME:         cert-manager
SHORT-DESCRIPTION:    Certificate management
PACKAGE-PROVIDER:     VMware
LONG-DESCRIPTION:     Provides certificate management provisioning within the cluster
MAINTAINERS:          [{Nicholas Seemiller}]
SUPPORT:              Go to https://cert-manager.io/ for documentation or the #cert-manager channel on Kubernetes slack
CATEGORY:             [certificate management]

% tanzu package available list cert-manager.community.tanzu.vmware.com -n default
\ Retrieving package versions for cert-manager.community.tanzu.vmware.com...
  NAME                                     VERSION  RELEASED-AT
  cert-manager.community.tanzu.vmware.com  1.3.1    2021-04-14T18:00:00Z
  cert-manager.community.tanzu.vmware.com  1.4.0    2021-06-15T18:00:00Z
  cert-manager.community.tanzu.vmware.com  1.5.1    2021-08-13T19:52:11Z
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

The following commands will verify that the package has been installed.

 ```sh
% tanzu package installed list
- Retrieving installed packages...
  NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
  cert-manager  cert-manager.community.tanzu.vmware.com  1.5.1            Reconcile succeeded


% kubectl get pods -A | grep cert-manager
tanzu-certificates      cert-manager-6476798c86-phqjh                               1/1     Running     0          20m
tanzu-certificates      cert-manager-cainjector-766549fd55-292j4                    1/1     Running     0          20m
tanzu-certificates      cert-manager-webhook-79878cbcbb-kttq9
```

With the Certificate Manager successfully deployed, the next step is to deploy an Ingress. Envoy, managed by Contour, is also available as a package with TCE.

## Deploy Contour (Ingress)

Later we shall deploy Prometheus and Grafana, which have a requirement on an Ingress/HTTPProxy. Contour [projectcontour.io](http://projectcontour.io) provides this functionality via an Envoy Ingress controller. Contour is an open source Kubernetes Ingress controller that acts as a control plane for the Envoy edge and service proxy.​

Prometheus has a requirement on an Ingress. Contour provides this functionality. Contour is an open source Kubernetes Ingress controller that acts as a control plane for the Envoy edge and service proxy.​

For our purposes of standing up a monitoring stack, we can provide a very simple data values file, in YAML format, when deploying Contour. In this manifest, we are requesting that the Envoy Ingress controller use a Load Balancer service which will be provided by NSX ALB, and that Contour leverages the previously deployed Cert-Manager to provision TLS certificates rather than using the upstream Contour cert-gen job to provision certificates. This secures communication between Contour and Envoy. You can optionally set more the number of Contour replicas as well:

```yaml
envoy:
  service:
    type: LoadBalancer
certificates:
  useCertManager: true
```

This is only a subset of the configuration parameters available in Contour. To display all configuation parameters, use the `--values-schema` option to display the configuration settings against the appropriate version of the package:

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

Note that there is currently no mechanism at present to display the configuration parameters in `yaml` format, but that further YAML examples can be found in the official package documentation.

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

### Check Contour data values have taken effect

The following command can be used to verify that the data values provided at deployment time have been implemented.

```sh
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

A good step at this point is to verify that Envoy is working as expected. To do that, we can locate the Envoy Pod, setup port-forwarding, and connect a browser to it once it has been as shown below:

```sh
% kubectl get pods -A | grep contour
projectcontour          contour-df5cc8689-7h9kh                                     1/1     Running     0          17m
projectcontour          contour-df5cc8689-q497w                                     1/1     Running     0          17m
projectcontour          envoy-mfjcp                                                 2/2     Running     0          17m


% kubectl get svc envoy -n projectcontour
NAME    TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
envoy   LoadBalancer   100.67.34.204   xx.xx.62.22   80:30639/TCP,443:32539/TCP   18h


% ENVOY_POD=$(kubectl -n projectcontour get pod -l app=envoy -o name | head -1)
% echo $ENVOY_POD
pod/envoy-mfjcp


% kubectl -n projectcontour port-forward $ENVOY_POD 9001
Forwarding from 127.0.0.1:9001 -> 9001
Forwarding from [::1]:9001 -> 9001
Handling connection for 9001
```

Note that I have deliberately obfuscated the first two octets of the IP address allocated to Envoy above. Now if you point a browser to the localhost:9001, the following Envoy landing page should be displayed:

![Envoy Listing](/docs/img/envoy-listings.png?raw=true)

Everything is now in place to deploy Prometheus.

## Deploy Prometheus

Prometheus [prometheus.io](http://prometheus.io) records real-time metrics and provides alerting capabilities. It has a requirement for an Ingress (or HTTPProxy) and that the requirement has been met by Contour. We can now proceed with the installation of the Prometheus community package. Prometheus has quite a number of configuration options, which can once again be displayed using the following commands. First determine the version, and then display the configuration options for that version. At present, there is only a single version of the Prometheus community package available:

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

This displays all of the configuration settings for Prometheus. For the purposes of this deployment, a simple manifest which enabled ingress and provides a virtual host fully qualified domain name are all that is needed. This is the sample manifest to modify the default Prometheus deployment, primarily to enable Ingress/HTTPProxy usage, and secondly to set a fully qualified domain name - fqdn - which would be used to access the Prometheus dashboard:

```yaml
ingress:
  enabled: true
  virtual_host_fqdn: "prometheus.rainpole.com"
  prometheus_prefix: "/"
  alertmanager_prefix: "/alertmanager/"
  prometheusServicePort: 80
  alertmanagerServicePort: 80
```

The virtual host fully qualified domain name may be added to the DNS using the same IP as that assigned to the Envoy Ingress by the NSX ALB. In the case of this example, this should map to IP Address xx.xx.62.22 as seen in `kubectl get svc` output after deploying Contour. If you have admin access to your DNS server for this environment, then you could add it manually. Another alternative is to integrate your deployment with an [ExternalDNS](https://github.com/kubernetes-sigs/external-dns). ExternalDNS synchronises exposed Kubernetes Services and Ingresses with DNS providers. What this means is that the ExternalDNS controller will interact with your infrastructure provider and will register the DNS name in the DNS service of the infrastructure provider. A final alternative would be to simply add the FQDN to the /etc/hosts file of the desktop where you are running the browser. Either way, assume that this step has now been done for this procedure. We can now deploy Prometheus:

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

### Check Prometheus data values have taken effect

The following command can be used to verify that the data values provided at deployment time have been implemented.

```sh
% tanzu package installed get prometheus -f /tmp/xxx
\ Retrieving installation details for prometheus... %

% more /tmp/xxx
---
ingress:
  enabled: true
  virtual_host_fqdn: "prometheus.rainpole.com"
  prometheus_prefix: "/"
  alertmanager_prefix: "/alertmanager/"
  prometheusServicePort: 80
  alertmanagerServicePort: 80
```

### Validate Prometheus functionality

The following Pods and Services should have been created successfully.

```sh
% kubectl get pods,svc -n prometheus
NAME                                                READY STATUS  RESTARTS AGE
pod/alertmanager-c45d9bf8c-86p5g                    1/1   Running 0        104s
pod/prometheus-cadvisor-bsbw4                       1/1   Running 0        106s
pod/prometheus-kube-state-metrics-7454948844-fxbwd  1/1   Running 0        104s
pod/prometheus-node-exporter-l6j42                  1/1   Running 0        104s
pod/prometheus-node-exporter-r8qcg                  1/1   Running 0        104s
pod/prometheus-pushgateway-6c69cb4d9c-6sttd         1/1   Running 0        104s
pod/prometheus-server-6587f4456c-xqxj6              2/2   Running 0        104s

NAME                                   TYPE       CLUSTER-IP      EXTERNAL-IP  PORT(S)        AGE
service/alertmanager                   ClusterIP  100.69.242.136  <none>       80/TCP         106s
service/prometheus-kube-state-metrics  ClusterIP  None            <none>       80/TCP,81/TCP  104s
service/prometheus-node-exporter       ClusterIP  100.71.136.28   <none>       9100/TCP       104s
service/prometheus-pushgateway         ClusterIP  100.70.127.19   <none>       9091/TCP       106s
service/prometheus-server              ClusterIP  100.66.19.135   <none>       80/TCP         104s
```

Contour provides an advanced resource type called [HttpProxy](https://projectcontour.io/docs/v1.18.1/config/fundamentals/) that provides some benefits over Ingress resources. We can also examine that this resource was created successfully:

```sh
% kubectl get HTTPProxy -A
NAMESPACE  NAME                  FQDN                    TLS SECRET     STATUS STATUS DESCRIPTION
prometheus prometheus-httpproxy  prometheus.rainpole.com prometheus-tls valid  Valid HTTPProxy
```

To verify that Prometheus is working correctly, point to the Prometheus FQDN (e.g. http:// prometheus.rainpole.com). If everything has worked correctly, you should be able to see a Prometheus dashboard:

![Envoy Dashboard Landing Page](/docs/img/envoy-db1.png?raw=true)

To do a very simple test, add a simple query, e.g. `prometheus_http_requests_total` and click Execute:

![Envoy Simple Query](/docs/img/envoy-db2.png?raw=true)

To check integration between Prometheus and Envoy, another query can be executed. When the Envoy landing page was displayed earlier, there was a section called `prometheus/stats`. These can now be queried as well, since these are the metrics that Envoy is sending to Prometheus. If we return to the Envoy landing page in the browser, and click on the prometheus/stats link and examine the metrics. one of these metrics, such as the `envoy_cluster_default_total_match`, and use it as a query in Prometheus (selecting Graph instead of Table this time):

![Envoy Prometheus Metric Query](/docs/img/envoy-db3.png?raw=true)

If you see something similar to this, then it would appear that Prometheus is working successfully. Now let's complete the monitoring stack by provisioning Grafana, and connecting it to our Prometheus data source.

## Deploy Grafana

[Grafana] (<https://grafana.com/>) is an analytics and interactive visualisation web application. Let's begin by displaying all of the configuring values that are available in Grafana. Once again, the package version is required to do this.

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

We will again try to keep it quite simple. Through the data values file, we can provide a data source (to Prometheus). The prometheus URL is an internal Kubernetes URL, made up of Pod Name and Namespace of the Prometheus Server. Since Grafana is also running in the cluster, they are able to communicate using internal K8s networking.

You will probably want to use a different virtual host fdqn, and you can add that it to your DNS once the Grafana Load Balancer service has allocated it with an IP address after deployment. As mentioned other options are to use an ExternalDNS, or add the entry to the local /etc/hosts file of the desktop which will launch the browser to Grafana. Since I have admin access to my DNS server, I can simply add this manually to my DNS.

The Grafana service type is set to Load Balancer by default.

```yaml
grafana:
  config:
    datasource_yaml: |-
      apiVersion: 1
      datasources:
        - name: Prometheus
          type: prometheus
          url: prometheus-server.prometheus.svc.cluster.local
          access: proxy
          isDefault: true
ingress:
  virtual_host_fqdn: "grafana.rainpole.com"
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
grafana:
  config:
    datasource_yaml: |-
      apiVersion: 1
      datasources:
        - name: Prometheus
          type: prometheus
          url: prometheus-server.prometheus.svc.cluster.local
          access: proxy
          isDefault: true
ingress:
  virtual_host_fqdn: "grafana.rainpole.com"
```

The following Pods, Services and HTTPProxy should have been created.

```sh
% kubectl get pods -A | grep grafana
grafana                 grafana-74ccf5fd4-27wm2                                     2/2     Running     0          109s

% kubectl get svc -A | grep grafana
grafana              grafana                         LoadBalancer   100.70.202.170   xx.xx.62.23   80:31227/TCP                 118s

$ kubectl get httpproxy -A
NAMESPACE                 NAME                   FQDN                         TLS SECRET       STATUS   STATUS DESCRIPTION
tanzu-system-dashboard    grafana-httpproxy      grafana.corinternal.com      grafana-tls      valid    Valid HTTPProxy
tanzu-system-monitoring   prometheus-httpproxy   prometheus.corinternal.com   prometheus-tls   valid    Valid HTTPProxy
```

As mentioned, Grafana uses a Load Balancer service type by default, so it has been provided with its own Load Balancer IP addess by NSX ALB. I have once more intentionally obfuscated the first two octets of the address. You can now add this to your DNS, like you did with Prometheus.

### Validate Grafana functionality

After adding your virtual host FQDN to your DNS, you can now connect to the Grafana dashboard using the FDQN. You connect directly to the Load Balancer IP address allocated to the Service. The login credentials are `admin/admin` initially, but you will need to change the password on first login. This is the landing page:

![Grafana Landing Page](/docs/img/grafana-landing-page.png?raw=true)

There is no need to add a datasource or create a dashboard - these have already been done for you.

To examine the data source, click on the icon representing datas sources on the left hand side (which looks like a cog). Here you can see the Prometheus data source that we placed in the data values manifest file when we deployed Grafana is already in place:

![Grafana Data Source Prometheus](/docs/img/grafana-data-source.png?raw=true)

Now click on the dashboards icon on the left hand side (it looks like a square of 4 smaller squares), and select `Manage` from the drop-down list. This will show the existing dashboards. There are 2 existing dashboards that have been provided; one is Kubernetes monitoring and the other is TKG monitoring. These dashboards are based on the Kubernetes Grafana dashboards found on [GitHub](https://github.com/kubernetes-monitoring/kubernetes-mixin).

![Grafana Dashboards Manager](/docs/img/grafana-manage-dashboards.png?raw=true)

Finally, select the TKG dashboard which is being sent metrics via the Prometheus data source. This provides an overview of the TKG cluster:

![TKG Dashboard](/docs/img/grafana-tkg-dashboard.png?raw=true)

The full monitoring stack of Contour/Envoy Ingress, with secure communication via Cert-Manager, alongside the Prometheus data scraper and Grafana visualization are now deployed through TCE community packages. Happy monitoring/analyzing.
