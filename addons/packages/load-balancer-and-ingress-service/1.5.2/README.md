# load-balancer-and-ingress-service

This package provides NSX Advanced Load Balancer.

## Components

* load-balancer-and-ingress-service

## Configuration

The following configuration values can be set to customize the load-balancer-and-ingress-service installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `name` | Required | describes the name of configuration.  |
| `namespace` | Required | describes which namespace ako will be deployed in. |
| `is_cluster_service` | Optional | describes if AKO is running in Kubernetes cluster. |
| `replica_count` | Required | describes the number of pods in AKO statefulset. |
| `ako_settings_log_level` | Optional | describes AKO log level, enum: INFO,DEBUG,WARN,ERROR. |
| `ako_settings_full_sync_frequency` | Required | describes how often AKO polls the Avi controller to update itself with cloud configurations. |
| `ako_settings_api_server_port` | Optional | describes internal port for AKO's API server for the liveness probe of the AKO pod default=8080. |
| `ako_settings_delete_config` | Required | describes if user wants to delete AKO created objects from AVI or not. |
| `ako_settings_disable_static_route_sync` | Required | describes ako should sync static routing or not. |
| `ako_settings_cluster_name` | Required | speficies the AVI Cloud AKO will be deployed with. |
| `ako_settings_cni_plugin` | Optional | describes which cni plugin cluster is using. |
| `ako_settings_enable_EVH` | Optional | describes enabling the Enhanced Virtual Hosting Model in Avi Controller for the Virtual Services or not, default value is false. |
| `ako_settings_l7_only` | Optional | describes if you want AKO only to do layer 7 load balancing. Default value is false. |
| `ako_settings_sync_namespace` | Optional | describes should AKO sync objects from this namespace   |
| `ako_settings_namespace_selector_key` | Optional | describes the namespace selector key. namespace_selector used for namespace migration, same label has to be present on namespace/s which needs migration/sync to AKO.|
| `ako_settings_namespace_selector_value` | Optional | describes the namespace selector key. namespace_selector used for namespace migration, same label has to be present on namespace/s which needs migration/sync to AKO. |
| `ako_settings_services_API` | Optional | describes if enables AKO in services API mode. Default value is false. |
| `network_settings_subnet_ip` | Required | describes the Data Networks gateway the AKO will be deployed with. |
| `network_settings_subnet_prefix` | Required | describes the Data Networks mask the AKO will be deployed with. |
| `network_settings_network_name` | Required | describes the Data Networks the AKO will be deployed with. |
| `network_settings_node_network_list` | Optional | describes the details of network and CIDRs are used in pool placement network for vcenter cloud. |
| `network_settings_bgp_peer_labels` | Optional | describes BGP peers, this is used for selective VsVip advertisement. |
| `network_settings_enable_RHI` | Optional | describes the cluster wide setting for BGP peering. Default value is false. |
| `network_settings_vip_network_list` | Required | describes network name of the VIP network |
| `l7_settings_disable_ingress_class` | Required | DisableIngressClass will prevent AKO Operator to install AKO IngressClass into workload clusters for old version of K8s. |
| `l7_settings_default_ing_controller` | Optional | describes ako is the default ingress controller to use. |
| `l7_settings_l7_sharding_scheme` | Optional | describes the hostname. |
| `l7_settings_service_type` | Optional | describes ingress methods for a service. enum NodePort,ClusterIP,NodePortLocal |
| `l7_settings_shard_vs_size` | Optional | describes ingress shared virtual service size. |
| `l7_settings_pass_through_shard_size` | Optional | describes the passthrough virtualservice numbers, ENUMs: LARGE, MEDIUM, SMALL |
| `l4_settings_default_domain` | Optional | describes default sub-domain to use for L4 VSes. |
| `l4_settings_advanced_l4` | Optional | describes the settings for the services API usage. Default value is false. |
| `l4_settings_auto_FQDN` | Optional | describes the FQDN generation. Valid value should be default(<svc>.<ns>.<subdomain>), flat (<svc>-<ns>.<subdomain>) or disabled, default values is disabled. |
| `controller_settings_service_engine_group_name` | Required | the group name of Service Engine that's to be used by the set of AKO Deployments. |
| `controller_settings_controller_version` | Required | describes The controller API version. |
| `controller_settings_cloud_name` | Required | the configured cloud name on the Avi controller. |
| `controller_settings_controller_ip` | Required | Avi controller ip. |
| `nodeport_selector_key` | Optional | Only applicable if serviceType is NodePort. |
| `nodeport_selector_value` | Optional | Only applicable if serviceType is NodePort. |
| `resources_limits_cpu` | Required | describes AKO statefulset cpu resources limitation. |
| `resources_limits_memory` | Required | describes AKO statefulset memory resources limitation. |
| `resources_request_cpu` | Required | describes AKO statefulset requests cpu resources. |
| `resources_request_memory` | Required | describes AKO statefulset requests memory resources. |
| `rbac_psp_enabled` | Required | describes if creates the pod security policy. |
| `rbac_psp_policy_api_version` | Optional | describes which api version should be use if pod secrurity policy is enabled. |
| `persistent_volume_claim` | Optional | describes which PVC using for AKO. |
| `mount_path` | Optional | describes AKO logs mount path. |
| `log_file` | Optional | describes where to store AKO logs. |
| `avi_credentials_username` | Required | describes username that addon manager will use to deploy avi secret. |
| `avi_credentials_password` | Required | describes password that addon manager will use to deploy avi secret. |
| `avi_credentials_certificate_authority_data` | Required | describes certificate_authority_data that addon manager will use to deploy avi secret. |

## Usage Example

This walkthrough guides you through using load-balancer-and-ingress-service...
