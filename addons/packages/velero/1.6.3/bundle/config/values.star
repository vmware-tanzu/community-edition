load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("@ytt:overlay", "overlay")

def labels():
  return {"component": "velero"}
end

# export
values = data.values
velero_app = overlay.subset({"metadata": {"labels": labels()}})

def validate_configs():
  data.values.namespace or assert.fail("Velero namespace should be provided")
end

def validate_storage():
  data.values.backupStorageLocation.spec.common.provider or assert.fail("backupStorageLocation needs a provider")
  data.values.backupStorageLocation.spec.objectStorage.bucket or assert.fail("backupStorageLocation needs a bucket")
end


# validate
def validate_velero():
  validate_funcs = [
    validate_configs,
    validate_storage,
  ]
   for validate_func in validate_funcs:
     validate_func()
   end
end

validate_velero()