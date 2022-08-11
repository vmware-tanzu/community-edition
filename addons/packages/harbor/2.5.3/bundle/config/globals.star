load("@ytt:data", "data")
load("@ytt:assert", "assert")

def get_image_location(repository, name, tag):
  return '{0}/{1}:{2}'.format(repository, name, tag)
end

def validate_infrastructure_provider():
  data.values.infrastructure_provider in ("aws", "vsphere", "azure") or assert.fail("infrastructure provider should be either aws or vsphere or azure")
end

def get_kapp_versioned_annotations(kapp):
  annotations = {}
  if kapp.config and kapp.config.versioned:
    annotations["kapp.k14s.io/versioned"] = ""
  end
  return annotations
end

def get_kapp_disable_wait_annotations():
  annotations = {}
  annotations["kapp.k14s.io/disable-wait"] = ""
  return annotations
end

def get_kapp_configmap_annotations(kapp):
  return get_kapp_versioned_annotations(kapp)
end

def get_kapp_namespace_annotations(kapp):
  annotations = {}
  if kapp.config and kapp.config.orphanNamespace:
    annotations["kapp.k14s.io/delete-strategy"] = "orphan"
  end
  return annotations
end

def get_kapp_secret_annotations(kapp):
  return get_kapp_versioned_annotations(kapp)
end

# Get kapp annotations for volumeClaimTemplates
# See https://github.com/k14s/kapp/issues/36
def get_kapp_vct_annotations():
  annotations = {}
  annotations["kapp.k14s.io/owned-for-deletion"] = ""
  return annotations
end

def get_kapp_annotations(kind):
  annotations = {}
  for deployer in data.values.app.deploy:
    if deployer.kapp:
      if kind == "ConfigMap":
        return get_kapp_configmap_annotations(deployer.kapp)
      elif kind == "Namespace":
        return get_kapp_namespace_annotations(deployer.kapp)
      elif kind == "Secret":
        return get_kapp_versioned_annotations(deployer.kapp)
      end
    end
  end
  return annotations
end

def get_resource_names_for_role(names):
  if names == None:
    return []
  else:
    return names.split(',')
  end
end

def getClusterRoleName():
   return "psp:" + data.values.namespace + ":harbor"
end

def getClusterRoleBindingName():
   return "psp:" + data.values.namespace + ":harbor"
end

#export
globals = data.values
