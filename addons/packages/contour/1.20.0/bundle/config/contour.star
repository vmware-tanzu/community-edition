load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("@ytt:struct", "struct")

def get_resource_names_for_role(names):
  if names == None:
    return []
  else:
    return names.split(',')
  end
end

SERVICE_TYPE_NODEPORT = "NodePort"
SERVICE_TYPE_LOADBALANCER = "LoadBalancer"

def get_envoy_service_type():
  if data.values.envoy.service.type:
    return data.values.envoy.service.type
  elif data.values.infrastructure_provider == "vsphere":
    return SERVICE_TYPE_NODEPORT
  else:
    return SERVICE_TYPE_LOADBALANCER
  end
end

def get_envoy_service_annotations():
  annotations = {}
  # This annotation tells kapp to disable wait for getting service type loadbalancer public IP.
  # This annotation is expected to be only used for local kind testing where service type loadbalancer will not get a public IP.
  if data.values.envoy.service.disableWait:
    annotations.update(get_kapp_disable_wait_annotations())
  end
  if data.values.infrastructure_provider == "aws":
    if data.values.envoy.service.aws.LBType == "nlb":
      annotations["service.beta.kubernetes.io/aws-load-balancer-type"] = "nlb"
    else:
      annotations["service.beta.kubernetes.io/aws-load-balancer-backend-protocol"] = "tcp"
      if data.values.contour.useProxyProtocol:
        annotations["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"] = "*"
      end
    end
  end
  if data.values.envoy.service.annotations:
    annotations_kvs = struct.decode(data.values.envoy.service.annotations)
    annotations.update(annotations_kvs)
  end
  return annotations
end

def get_kapp_disable_wait_annotations():
  annotations = {}
  annotations["kapp.k14s.io/disable-wait"] = ""
  return annotations
end

# ##########
# VALIDATION
# ##########

def validate_contour():
  validate_funcs = [validate_infrastructure_provider,
                    validate_contour_namespace,
                    validate_contour_deployment,
                    validate_contour_certificate,
                    validate_envoy_deployment,
                    validate_envoy_service]
   for validate_func in validate_funcs:
     validate_func()
   end
end

def validate_infrastructure_provider():
  data.values.infrastructure_provider in ("aws", "vsphere", "azure") or assert.fail("infrastructure provider should be either aws or vsphere or azure")
end

def validate_contour_namespace():
  data.values.namespace or assert.fail("Contour namespace should be provided")
end

def validate_contour_deployment():
  data.values.contour.replicas or assert.fail("Contour replicas should be provided")
end

def validate_contour_certificate():
  data.values.certificates.duration or assert.fail("Contour certificate duration should be provided")
  data.values.certificates.renewBefore or assert.fail("Contour certificate renewBefore should be provided")
end

def validate_envoy_deployment():
  if data.values.envoy.hostPorts.enable:
    data.values.envoy.hostPorts.http or assert.fail("Envoy http hostPort should be provided if hostPorts.enable is true")
    data.values.envoy.hostPorts.https or assert.fail("Envoy https hostPort should be provided if hostPorts.enable is true")
  end
  data.values.envoy.logLevel or assert.fail("Envoy log level should be provided")
  data.values.envoy.logLevel in ("trace", "debug", "info", "warning", "warn", "error", "critical", "off") or assert.fail("Envoy log level should be trace|debug|info|warning/warn|error|critical|off")
  data.values.envoy.terminationGracePeriodSeconds or assert.fail("Envoy terminationGracePeriodSeconds should be provided")
end

def validate_envoy_service():
  if data.values.envoy.service.type:
    if data.values.infrastructure_provider == "vsphere":
      data.values.envoy.service.type in ("LoadBalancer", "NodePort") or assert.fail("For vsphere infrastructure provider the envoy.service.type should be LoadBalancer or NodePort")
    else:
      data.values.envoy.service.type in ("LoadBalancer", "NodePort", "ClusterIP") or assert.fail("Envoy service type should be LoadBalancer or NodePort or ClusterIP")
    end
  end
  if data.values.envoy.service.externalTrafficPolicy:
    data.values.envoy.service.externalTrafficPolicy in ("Cluster", "Local") or assert.fail("Envoy service externalTrafficPolicy should be Cluster or Local")
  end
  if data.values.infrastructure_provider == "aws":
    data.values.envoy.service.aws.LBType in ("classic", "nlb") or assert.fail("Envoy service aws LB Type should be classic or nlb")
  end
end

validate_contour()
