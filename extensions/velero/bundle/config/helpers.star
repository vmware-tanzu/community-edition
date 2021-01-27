# Helper functions for ytt templates.

# backupStorageLocationConfig builds the config object for a Velero
# BackupStorageLocation given a provider (e.g., AWS, Azure, minio)
def backupStorageLocationConfig(provider, backupStorageLocationValues):
    if provider == 'minio':
        return {
            'region': 'minio',
            's3ForcePathStyle': "true",
            's3Url': 'http://minio:9000'
        }
    elif provider == 'aws':
        return {
            'region': backupStorageLocationValues.aws.region
        }
    elif provider == 'azure':
        azureValues = backupStorageLocationValues.azure
        return {
            'resourceGroup':  azureValues.resourceGroup,
            'storageAccount':  azureValues.storageAccount,
            'storageAccountKeyEnvVar':  azureValues.storageAccountKeyEnvVar,
            'subscriptionId':  azureValues.subscriptionId,
            'blockSizeInBytes':  azureValues.blockSizeInBytes
        }
    end
end

# volumeSnapshotLocationConfig builds the config object for a Velero
# VolumeSnapshotLocation given a provider (e.g., AWS, Azure, vsphere) 
def volumeSnapshotLocationConfig(provider, volumeSnapshotLocationValues):
    if provider == 'aws':
        return {
            'region': volumeSnapshotLocationValues.aws.region
        }
    elif provider == 'azure':
        return {
            'apiTimeout': volumeSnapshotLocationValues.azure.apiTimeout,
            'resourceGroup': volumeSnapshotLocationValues.azure.resourceGroup,
            'subscriptionId': volumeSnapshotLocationValues.azure.subscriptionId,
            'incremental': volumeSnapshotLocationValues.azure.incremental
        }
    elif provider == 'vsphere':
        return {}
    else:
        return {}
    end
end
