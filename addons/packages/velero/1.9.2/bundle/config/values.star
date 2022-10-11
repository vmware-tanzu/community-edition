load("@ytt:data", "data")
load("@ytt:assert", "assert")
load("@ytt:overlay", "overlay")

#! Use this library for functions that operate on the values data.

# helpers
def secret_name():
  if values.credential.useDefaultSecret:
    return values.credential.name
  end
  return values.backupStorageLocation.spec.existingSecret.name
end

def resource(kind, name):
  return {"kind": kind,"metadata":{"name": name}}
end

def labels():
  return {"component": "velero"}
end

# validations
def validate_storage():
  values.backupStorageLocation.spec.provider or assert.fail("backupStorageLocation needs a provider, velero needs at least one backup storage location")
  values.backupStorageLocation.spec.objectStorage.bucket or assert.fail("backupStorageLocation needs a bucket")
end

# Note: aws and azure are the only object storage providers that the TCE Velero package supports.
# Neither vSphere or Docker provides storage.
def validate_storage_provider():
  providers = ["aws", "azure"]
  if values.backupStorageLocation.spec.provider:
    values.backupStorageLocation.spec.provider in providers or assert.fail("storage provider should be aws or azure")
  end
end

# Note: Docker does not provide snapshotting capabilities.
def validate_snapshot_provider():
  if values.volumeSnapshotLocation.snapshotsEnabled:
    providers = ["aws", "azure", "vsphere"]
    if values.volumeSnapshotLocation.spec.provider:
      values.volumeSnapshotLocation.spec.provider in providers or assert.fail("a snapshot provider should be either aws, vsphere  or azure")
    end
  end
end

def validate_provider_config():
  if values.backupStorageLocation.create:
    if values.backupStorageLocation.spec.provider == "aws":
      values.backupStorageLocation.spec.configAWS.region or assert.fail("a region must be set for the AWS backup storage location")
    end
    if values.backupStorageLocation.spec.provider == "azure":
      values.backupStorageLocation.spec.configAzure.resourceGroup or assert.fail("a resourceGroup must be set for the Azure backup storage location")
      values.backupStorageLocation.spec.configAzure.storageAccount or assert.fail("a storageAccount must be set for the Azure backup storage location")
    end
  end
  if values.volumeSnapshotLocation.snapshotsEnabled:
    if values.volumeSnapshotLocation.spec.provider == "aws":
      values.volumeSnapshotLocation.spec.configAWS.region or assert.fail("a region must be set for the AWS volume snapshot location")
    end
    if values.volumeSnapshotLocation.spec.provider == "vsphere":
      values.volumeSnapshotLocation.spec.configvSphere.region or assert.fail("a region must be set for the vSphere volume snapshot location")
      values.vsphere.create or assert.fail("vSphere configuration is not set when volume snapshot provider is vSphere")
    end
  end
  if values.vsphere.create:
    if values.volumeSnapshotLocation.snapshotsEnabled == False or values.volumeSnapshotLocation.spec.provider != "vsphere":
      assert.fail("volume snapshot provider is not vSphere or not enable when vsphere is supposed to created")
    end
    if values.backupStorageLocation.spec.provider != "aws":
      assert.fail("BackupStorageLocation provider is not set to aws for vsphere environment")
    end
  end
end

def validate_secret():
  if values.credential.useDefaultSecret:
    values.credential.name or assert.fail("must specify a name for the default secret to be used by Velero")
    values.credential.secretContents or assert.fail("must specify the content for the default secret to be used by Velero")
    values.credential.secretContents.cloud != None or assert.fail("the default secret must have a key named `cloud`")
    values.credential.secretContents.cloud or assert.fail("the default secret key `cloud` must contain the raw credentials")
  else:
    values.backupStorageLocation.spec.existingSecret.name or assert.fail("must specify the name of the existing secret to be used by Velero")
    values.backupStorageLocation.spec.existingSecret.key or assert.fail("must specify the key of the existing secret to be used by Velero")
  end
end

def validate_velero():
  validate_funcs = [
    validate_storage_provider,
    validate_snapshot_provider,
    validate_provider_config,
    validate_storage,
    validate_secret,
  ]
   for validate_func in validate_funcs:
     validate_func()
   end
end

# export
values = data.values
velero_app = overlay.subset({"metadata": {"labels": labels()}})

validate_velero()