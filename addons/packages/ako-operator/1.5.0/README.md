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
| `avi_data_network` | Required | describes the Data Networks the AKO will be deployed with. |
| `avi_data_network_cidr` | Required | describes the Data Networks the AKO will be deployed with. |
| `avi_ca_data_b64` | Required | Avi controller credential. |
| `avi_labels` | Optional | Label used to select Clusters. The Clusters that are selected by this will be the ones affected by this AKODeploymentConfig. |
| `avi_cni_plugin` | Optional | describes which cni plugin cluster is using. AKO supported CNI: antrea,calico,canal,flannel,openshift and ncp |
| `avi_disable_static_route_sync` | Optional | describes ako should sync static routing or not. |
| `avi_control_plane_ha_provider` | Required | describes whether Avi provides control plane HA service or not. |
| `avi_management_cluster_vip_network_name` | Required | describes the data network name of the management cluster AKO will be deployed with. |
| `avi_management_cluster_vip_network_cidr` | Required | describes the data network cidr of the management cluster AKO will be deployed with. |
| `avi_control_plane_endpoint_port`| Optional | describe the port of AVI control plane endpoint |
| `avi_enable_evh` | Optional | describes should the Enhanced Virtual Hosting Model in Avi controller for Virtual Services be enabled |
| `avi_l7_only` | Optional | describes should AKO only to do layer 7 load balancing |
| `avi_services_api` | Optional |  describes should enable AKO in [services API mode](https://kubernetes-sigs.github.io/service-apis/) |
| `avi_namespace_selector_label_key` | Optional | describes the label key used for namespace migration |
| `avi_namespace_selector_label_value` | Optional | describes the label value used for namespace migration |
| `avi_enable_rhi` |  Optional | describes should the Route Health Injection be enabled for BGP |
| `avi_bgp_peer_labels` |  Optional | describe BGP peers, this is used for selective VsVip advertisement |
| `avi_no_pg_for_sni` | Optional |  describes if you want to get rid of pool groups from SNI VSes |
| `avi_advanced_l4` | Optional | describes the settings for the services API usage for L4 |
| `avi_auto_fqdn` | Optional |  describes the FQDN generation for L4: default,flat,disabled |

## Usage Example

This walkthrough guides you through using ako-operator...
