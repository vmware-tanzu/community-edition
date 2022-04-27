load("@ytt:data", "data")
load("/globals.star", "globals", "validate_infrastructure_provider")
load("@ytt:assert", "assert")
load("/libs/constants.lib.yaml", "is_dex_required")

SERVICE_TYPE_NODEPORT = "NodePort"
SERVICE_TYPE_LOADBALANCER = "LoadBalancer"

def validate_pinniped():
  data.values.tkg_cluster_role in ("management", "workload") or assert.fail("tkg_cluster_role must be provided to be either 'management' or 'workload'")
  data.values.infrastructure_provider in ("vsphere", "azure", "aws") or assert.fail("infrastructure_provider must be provided to be either 'vsphere', 'azure' or 'aws'")
  if data.values.identity_management_type:
    if is_mgmt_cluster():
      data.values.identity_management_type in ("oidc", "ldap") or assert.fail("identity_management_type for management clusters must be either 'oidc' or 'ldap'")
    end
    if is_workload_cluster():
      data.values.identity_management_type in ("none", "oidc", "ldap") or assert.fail("identity_management_type for workload clusters be either 'none', 'oidc' or 'ldap'")
    end
  end
  if render_on_workload_cluster():
    data.values.pinniped.supervisor_svc_endpoint or assert.fail("the pinniped.supervisor_svc_endpoint must be provided")
    data.values.pinniped.supervisor_ca_bundle_data or assert.fail("the pinniped.supervisor_ca_bundle_data must be provided")
  end
end

def validate_dex():
  if is_workload_cluster():
    return
  end
  validate_funcs = [validate_infrastructure_provider,
                    validate_dex_namespace,
                    validate_dex_config,
                    validate_dex_certificate,
                    validate_dex_deployment,
                    validate_dex_service]
  for validate_func in validate_funcs:
    validate_func()
  end
end

def validate_dex_namespace():
  data.values.dex.namespace or assert.fail("Dex namespace should be provided")
end

def validate_dex_config():
  globals.infrastructure_provider in ("aws", "vsphere", "azure") or assert.fail("Dex supports provider aws, vsphere or azure")
  if globals.infrastructure_provider == "vsphere":
    data.values.dex.dns.vsphere.ipAddresses or assert.fail("Dex MGMT_CLUSTER_IP should be provided for vsphere provider")
    data.values.dex.config.issuerPort or assert.fail("Dex config issuerPort should be provided for vsphere provider")
  end
  if globals.infrastructure_provider == "aws":
    data.values.dex.dns.aws.DEX_SVC_LB_HOSTNAME or assert.fail("Dex oidc issuer DEX_SVC_LB_HOSTNAME should be provided for aws provider")
  end
  if globals.infrastructure_provider == "azure":
    data.values.dex.dns.azure.DEX_SVC_LB_HOSTNAME or assert.fail("Dex DEX_SVC_LB_HOSTNAME should be provided for azure provider")
  end

  validate_ldap_config()

  data.values.dex.config.oauth2 or assert.fail("Dex oauth2 should be provided")
  data.values.dex.config.storage or assert.fail("Dex storage should be provided")
end

def validate_ldap_config():
  data.values.dex.config.ldap.host or assert.fail("Dex ldap <LDAP_HOST> should be provided")
  data.values.dex.config.ldap.insecureSkipVerify in (True, False)
  if data.values.dex.config.ldap.userSearch :
    data.values.dex.config.ldap.userSearch.baseDN or assert.fail("Dex ldap userSearch enabled. baseDN should be provided")
  end
  if data.values.dex.config.ldap.groupSearch :
    data.values.dex.config.ldap.groupSearch.baseDN or assert.fail("Dex ldap groupSearch enabled. baseDN should be provided")
  end
end

def validate_dex_certificate():
  data.values.dex.certificate.duration or assert.fail("Dex certificate duration should be provided")
  data.values.dex.certificate.renewBefore or assert.fail("Dex certificate renewBefore should be provided")
end

def validate_dex_deployment():
  data.values.dex.deployment.replicas or assert.fail("Dex deployment replicas should be provided")
end

def validate_dex_service():
  if data.values.dex.service.type:
    data.values.dex.service.type in ("LoadBalancer", "NodePort") or assert.fail("Dex service type should be LoadBalancer or NodePort")
  end
  if globals.infrastructure_provider == "aws":
    data.values.dex.dns.aws.DEX_SVC_LB_HOSTNAME or assert.fail("Dex aws dnsname DEX_SVC_LB_HOSTNAME should be provided")
  end
  if globals.infrastructure_provider == "vsphere":
    data.values.dex.dns.vsphere.ipAddresses[0] or assert.fail("Dex vsphere dns at least one ipaddress should be provided")
  end
  if globals.infrastructure_provider == "azure":
    data.values.dex.dns.azure.DEX_SVC_LB_HOSTNAME or assert.fail("Dex azure dnsname DEX_SVC_LB_HOSTNAME should be provided")
  end
end

# vsphere, aws, azure currently supported.
def get_default_service_type():
  if globals.infrastructure_provider == "vsphere":
    return SERVICE_TYPE_NODEPORT
  else:
    return SERVICE_TYPE_LOADBALANCER
  end
end

def get_pinniped_supervisor_service_type():
  if hasattr(data.values.pinniped, "supervisor") and hasattr(data.values.pinniped.supervisor, "service") and hasattr(data.values.pinniped.supervisor.service, "type") and data.values.pinniped.supervisor.service.type != None:
    return data.values.pinniped.supervisor.service.type
  else:
    return get_default_service_type()
  end
end

def is_pinniped_supervisor_service_type_NodePort():
  return get_pinniped_supervisor_service_type() == SERVICE_TYPE_NODEPORT
end

def get_dex_service_type():
  if hasattr(data.values.dex, "service") and hasattr(data.values.dex.service, "type") and data.values.dex.service.type != None:
    return data.values.dex.service.type
  else:
    return get_default_service_type()
  end
end

def is_dex_service_type_LB():
  return get_dex_service_type() == SERVICE_TYPE_LOADBALANCER
end

def is_dex_service_NodePort():
  return get_dex_service_type() == SERVICE_TYPE_NODEPORT
end

def get_dex_service_annotations():
  annotations = {}
  if hasattr(data.values.dex, "service") and hasattr(data.values.dex.service, "annotations") and data.values.dex.service.annotations:
     annotations = data.values.dex.service.annotations
  end
  if globals.infrastructure_provider == "aws":
    annotations["service.beta.kubernetes.io/aws-load-balancer-backend-protocol"] = "ssl"
  end
  return annotations
end

def validate_static_client() :
  if data.values.dex.config.staticClients and len(data.values.dex.config.staticClients) > 0:
    for client in data.values.dex.config.staticClients :
      getattr(client, "id") or assert.fail("Dex staticClients should have id")
      getattr(client, "redirectURIs") or assert.fail("Dex staticClients should have redirectURIs")
      getattr(client, "name") or assert.fail("Dex staticClients should have name")
      getattr(client, "secret") or assert.fail("Dex staticClients should have secret")
    end
  end
end

def is_mgmt_cluster():
  return data.values.tkg_cluster_role == "management"
end

def is_workload_cluster():
  return data.values.tkg_cluster_role == "workload"
end

# at present we ignore identity_management_type as a render gate for mgmt clusters
def render_on_mgmt_cluster():
  return is_mgmt_cluster()
end

def render_on_workload_cluster():
  return is_workload_cluster() and data.values.identity_management_type != "none"
end

#export
values = data.values


# validate dex
if is_dex_required():
  validate_dex()
end

# validate pinniped
validate_pinniped()
