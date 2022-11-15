# Load Balancer Operator

This package reconciles Cluster API objects and manages the lifecycles of [AVI Kubernetes Operator (AKO)](https://github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/tree/release-1.7.0) using [ako-operator](https://github.com/vmware-samples/load-balancer-operator-for-kubernetes).

## Components

* [Deployment](./bundle/config/upstream/akooperator/deployment.yaml) for the Load Balancer Operator (ako-operator)
* [CRD](./bundle/config/upstream/akooperator/akodeploymentconfig.yaml) for `AKODeploymentConfig`
* [Two `AKODeploymentConfig` samples](./bundle/config/upstream/akooperator/akodeploymentconfig.yaml) for management and workload clusters respectively

## Supported Providers

The following tables shows the providers this package can work with. Other cloud provider support will be added
in the future.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
|  ❌ |   ❌ | ✅  |  ❌  |

## Mailing lists

* Use [tkg-infrax-akita](mailto:tkg-infrax-akita@groups.vmware.com) to report security concerns to the AKO Team,
  who uses the list to privately discuss security issues and fixes prior to disclosure.

[comment]: <> (* Join the[AKO Distributors]&#40;mailto:tkg-infrax-akita@groups.vmware.com&#41; mailing list for early private information and vulnerability disclosure.)

[comment]: <> (  Early disclosure may include mitigating steps and additional information on security patch releases.)

[comment]: <> (* Send new membership requests to tkg-infrax-akita@groups.vmware.com.)

[comment]: <> (  In the body of your request please specify how you qualify for)

[comment]: <> (  membership and fulfill each criterion listed in the Membership Criteria section above.)

## ako-operator integration with Tanzu Community Edition (TCE)

### Prerequisites

* vSphere is the supported cloud provider for this package. Before proceeding, make sure you have a running management cluster on vSphere
following the [steps](https://tanzucommunityedition.io/docs/vsphere-intro/). Verify it is running with `tanzu management-cluster get`.
* Retrieve the context of the management cluster with `tanzu management-cluster kubeconfig get --admin`, and then switch to the context.
* Install and set up the AVI Controller on the vCenter Server, following the [documentation](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.4/vmware-tanzu-kubernetes-grid-14/GUID-mgmt-clusters-install-nsx-adv-lb.html#install-avi-controller-on-vcenter-server-7).

### Installation

Fetch a list of available versions.

```shell
tanzu package available list ako-operator.tanzu.vmware.com -n tkg-system
```

The ako-operator package with specified version and values can be installed with the Tanzu package CLI.
Refer to [next section](#configuration-values) for a detailed list of configuration values in `values.yaml`.

```shell
tanzu package install ako-operator --values-file values.yaml \
  --package-name ako-operator.community.tanzu.vmware.com \
  --version 1.7.0 --namespace tkg-system
```

Verify the ako-operator app has reconciled successfully in the management cluster, by asserting the following grep is non-empty.

```shell
kubectl get apps -A | grep ako-operator
```

### Removal

Issue a package delete command with the Tanzu package CLI to remove ako-operator. For example,

```shell
tanzu package installed delete ako-operator --namespace tkg-system
```

## Configuration

### Configuration values

A sample configuration values file `values.yaml` is provided [here](./bundle/config/values.yaml). You should modify this template
and customize configurations, especially those related to AVI controller credentials.
Make sure they are correctly configured for the fields including `avi_controller`, `avi_username/password`, `avi_ca_data_b64` etc.

A minimum configuration values file looks like

```yaml
akoOperator:
  avi_enable: true
  namespace: tkg-system-networking
  config:
    avi_admin_credential_name: avi-controller-credentials
    avi_ca_name: avi-controller-ca
    avi_controller: "10.191.164.223"
    avi_username: "<avi_username>"
    avi_password: "<avi password>"
    avi_cloud_name: "Default-Cloud"
    avi_service_engine_group: "Default-Group"
    avi_management_cluster_service_engine_group: "MC-SEG"
    avi_data_network: "VM Network"
    avi_data_network_cidr: "10.191.160.0/20"
    avi_management_cluster_vip_network_name: "VM Network"
    avi_management_cluster_vip_network_cidr: "10.191.160.0/20"
    avi_ca_data_b64: "LS0...LQo="
```

### Global

The following configuration values can be set to customize the ako-operator installation.

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `avi_enable` | Required | describes whether Avi is used or not. |
| `namespace` | Required | the namespace in which to deploy ako-operator. |
| `cluster_name` | Required | specifies the AVI Cloud AKO will be deployed with. |
| `avi_disable_ingress_class` | Optional | DisableIngressClass will prevent AKO Operator to install AKO IngressClass into workload clusters for old version of K8s. |
| `avi_ingress_default_ingress_controller` | Optional | describes ako is the default ingress controller to use. |
| `avi_ingress_shard_vs_size` | Optional | describes ingress shared virtual service size. Valid value should be SMALL, MEDIUM, LARGE or DEDICATED, default value is SMALL. |
| `avi_ingress_service_type` | Optional | describes ingress methods for a service. Valid value should be NodePort, ClusterIP or NodePortLocal. |
| `avi_ingress_node_network_list` | Optional | describes the details of network and CIDRs are used in pool placement network for vcenter cloud. |
| `avi_admin_credential_name` | Required | the name of a Secret resource which includes the username and password to access and configure the Avi Controller. |
| `avi_ca_name` | Required | Avi controller credential name. |
| `avi_controller` | Required | Avi controller ip. |
| `avi_controller_version` | optional | The version of the Avi controller you want AKO Operator and AKO talk to. |
| `avi_username` | Required | Avi controller username. |
| `avi_password` | Required | Avi controller password. |
| `avi_cloud_name` | Required | the configured cloud name on the Avi controller. |
| `avi_service_engine_group` | Required | the group name of Service Engine that's to be used by the set of AKO Deployments. |
| `avi_management_cluster_service_engine_group` | Required | the group name of Service Engine that's to be used by the AKO in management cluster. |
| `avi_data_network` | Required | describes the Data Networks the AKO will be deployed with. |
| `avi_data_network_cidr` | Required | describes the Data Networks the AKO will be deployed with. |
| `avi_control_plane_network` | Required | describes the Control Plane Networks of the workload cluster. |
| `avi_control_plane_network_cidr` | Required | describes the Control Plane Networks of the workload cluster. |
| `avi_ca_data_b64` | Required | Avi controller credential. |
| `avi_labels` | Optional | Label used to select Clusters. The Clusters that are selected by this will be the ones affected by this AKODeploymentConfig. |
| `avi_cni_plugin` | Optional | describes which cni plugin cluster is using. AKO supported CNI: antrea,calico,canal,flannel,openshift and ncp |
| `avi_disable_static_route_sync` | Optional | describes ako should sync static routing or not. |
| `avi_control_plane_ha_provider` | Required | describes whether Avi provides control plane HA service or not. |
| `avi_management_cluster_vip_network_name` | Required | describes the data network name of the management cluster AKO will be deployed with. |
| `avi_management_cluster_vip_network_cidr` | Required | describes the data network cidr of the management cluster AKO will be deployed with. |
| `avi_management_cluster_control_plane_vip_network_name` | Required | describes the name of the management cluster control plane network. |
| `avi_management_cluster_control_plane_vip_network_cidr` | Required | describes the cidr of the management cluster control plane network. |
| `avi_control_plane_endpoint_port`| Optional | describe the port of AVI control plane endpoint |

### Verify your configurations

After the package is installed and reconcile successfully, two AKODeploymentConfigs are deployed globally.
You can list and inspect them with `kubectl get akodeploymentconfig`. They are the configuration interfaces for the ako-operator.
In general,

* `install-ako-for-management-cluster` specifies the parameters to deploy AKO statefulset in the management cluster, which
  then provides L4/L7 load balancing for applications.
* `install-ako-for-all` specifies the parameters to deploy AKO statefulsets and reconcile users in each future workload clusters. The ako-operator
  passes these parameters to each workload clusters via addon secrets.

In the management cluster, ensure that there is no error level message in the ako-operator manage log.

```shell
# retrieve the pod name of the ako-operator manager
kubectl get pods -n tkg-system-networking

# check if the log is error-free.
kubectl logs ako-operator-controller-manager-xxxxxxxx-xxxxx manager -n tkg-system-networking
```

## Usage Example

### Create an L4 Service in management cluster

Deploy a `LoadBalancer` type of service with the following manifest.

```shell
$ cat <<EOF >sample-svc.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: corgi-test
  labels:
    corgi: test
spec:
  selector:
    matchLabels:
      corgi: test
  template:
    metadata:
      labels:
        corgi: test
    spec:
      containers:
        - name: nginx
          image: harbor-repo.vmware.com/dockerhub-proxy-cache/library/nginx
          ports:
            - containerPort: 80

---
apiVersion: v1
kind: Service
metadata:
  name: corgi-test
spec:
  type: LoadBalancer
  selector:
    corgi: test
  ports:
    - nodePort: 30008
      port: 80
      targetPort: 80
EOF
```

After that, verify that the service has been assigned an external IP using `kubectl`. Visit the AVI controller
dashboard at `http://<avi_controller_ip>` and ensure a virtual service has been created accordingly for this deployment.

```shell
$ kubectl get svc
NAME         TYPE           CLUSTER-IP       EXTERNAL-IP     PORT(S)        AGE
corgi-test   LoadBalancer   100.67.166.194   10.191.167.75   80:30008/TCP   103s
```

![virtual_svc_on_avi_controller](./images/virtual_svc_on_avi_controller.png)
