# ako-operator

This package provides NSX Advanced Load Balancer using ako-operator.

## Components

* ako-operator

## Configuration

The following configuration values can be set to customize the ako-operator installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `avi_enable` | Required | describes whether Avi is used or not. |
| `namespace` | Required | the namespace in which to deploy ako-operator. |
| `cluster_name` | Required | speficies the AVI Cloud AKO will be deployed with. |
| `avi_disable_ingress_class` | Optional | DisableIngressClass will prevent AKO Operator to install AKO IngressClass into workload clusters for old version of K8s. |
| `avi_ingress_default_ingress_controller` | Optional | describes ako is the default ingress controller to use. |
| `avi_ingress_shard_vs_size` | Optional | describes ingress shared virtual service size. |
| `avi_ingress_service_type` | Optional | describes ingress methods for a service. |
| `avi_ingress_node_network_list` | Optional | describes the details of network and CIDRs are used in pool placement network for vcenter cloud. |
| `avi_admin_credential_name` | Required | the name of a Secret resource which includes the username and password to access and configure the Avi Controller. |
| `avi_ca_name` | Required | Avi controller credential name. |
| `avi_controller` | Required | Avi controller ip. |
| `avi_username` | Required | Avi controller username. |
| `avi_password` | Required | Avi controller password. |
| `avi_cloud_name` | Required | the configured cloud name on the Avi controller. |
| `avi_service_engine_group` | Required | the group name of Service Engine that's to be used by the set of AKO Deployments. |
| `avi_data_network` | Required | describes the Data Networks the AKO will be deployed with. |
| `avi_data_network_cidr` | Required | describes the Data Networks the AKO will be deployed with. |
| `avi_ca_data_b64` | Required | Avi controller credential. |
| `avi_labels` | Optional | Label used to select Clusters. The Clusters that are selected by this will be the ones affected by this AKODeploymentConfig. |
| `avi_cni_plugin` | Optional | describes which cni plugin cluster is using. |
| `avi_disable_static_route_sync` | Optional | describes ako should sync static routing or not. |
| `avi_control_plane_ha_provider` | Required | describes whether Avi provides control plane HA service or not. |
| `avi_management_cluster_vip_network_name` | Required | describes the data network name of the management cluster AKO will be deployed with. |
| `avi_management_cluster_vip_network_cidr` | Required | describes the data network cidr of the management cluster AKO will be deployed with. |

## Usage Example

This walkthrough guides you through using ako-operator...