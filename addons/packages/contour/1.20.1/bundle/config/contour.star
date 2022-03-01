load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("@ytt:struct", "struct")

# ##########
# DEFAULTING
# ##########

def get_envoy_service_type():
  if data.values.envoy.service.type:
    return data.values.envoy.service.type
  elif data.values.infrastructureProvider == "docker":
    return "NodePort"
  elif data.values.infrastructureProvider == "vsphere":
    return "NodePort"
  else:
    return "LoadBalancer"
  end
end

def get_envoy_service_external_traffic_policy():
  if data.values.envoy.service.externalTrafficPolicy:
    return data.values.envoy.service.externalTrafficPolicy
  elif data.values.infrastructureProvider == "vsphere":
    return "Cluster"
  else:
    return "Local"
  end
end

def get_envoy_service_annotations():
  annotations = {}

  if data.values.infrastructureProvider == "aws":
    if data.values.envoy.service.aws.loadBalancerType == "nlb":
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
  if data.values.infrastructureProvider:
    data.values.infrastructureProvider in ("docker", "aws", "vsphere", "azure") or assert.fail("infrastructureProvider must be either docker or aws or vsphere or azure")
  end
end

def validate_contour_namespace():
  data.values.namespace or assert.fail("namespace must be provided")
end

def validate_contour_deployment():
  data.values.contour.replicas or assert.fail("contour.replicas must be provided")
end

def validate_contour_certificate():
  if data.values.certificates.useCertManager:
    data.values.certificates.duration or assert.fail("certificates.duration must be provided when certificates.useCertManager is true")
    data.values.certificates.renewBefore or assert.fail("certificates.renewBefore must be provided when certificates.useCertManager is true")
  end
end

def validate_envoy_deployment():
  if data.values.envoy.hostPorts.enable:
    data.values.envoy.hostPorts.http or assert.fail("envoy.hostPorts.http must be provided when envoy.hostPorts.enable is true")
    data.values.envoy.hostPorts.https or assert.fail("envoy.hostPorts.https must be provided when envoy.hostPorts.enable is true")
  end
  
  data.values.envoy.logLevel in ("trace", "debug", "info", "warning", "warn", "error", "critical", "off") or assert.fail("envoy.logLevel must be one of trace|debug|info|warning/warn|error|critical|off")
  
  data.values.envoy.terminationGracePeriodSeconds or assert.fail("envoy.terminationGracePeriodSeconds must be provided")
end

def validate_envoy_service():
  if data.values.envoy.service.type:
    data.values.envoy.service.type in ("LoadBalancer", "NodePort", "ClusterIP") or assert.fail("envoy.service.type must be either LoadBalancer or NodePort or ClusterIP")
  end

  if data.values.envoy.service.externalTrafficPolicy:
    data.values.envoy.service.externalTrafficPolicy in ("Cluster", "Local") or assert.fail("envoy.service.externalTrafficPolicy must be either Cluster or Local")
  end

  if data.values.infrastructureProvider == "aws":
    data.values.envoy.service.aws.loadBalancerType in ("classic", "nlb") or assert.fail("envoy.service.aws.loadBalancerType must be either classic or nlb when infrastructureProvider is aws")
  end
end

validate_contour()
