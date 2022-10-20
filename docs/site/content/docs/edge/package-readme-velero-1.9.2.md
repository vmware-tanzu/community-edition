# Velero

This package provides safe operations to backup, restore, perform disaster recovery, and migrate Tanzu Community Edition cluster resources and persistent volumes using the open source tool [Velero](https://velero.io/).

## Supported Providers

The following table shows the providers this package can work with.

| AWS | Azure | vSphere | Docker | CSI |
|:---:|:---:|:---:|:---:|:---:|
| ✅ | ✅ | ✅ | ✅ | ✅ |

## Components

* Velero Deployment
* Restic DaemonSet

### Custom Resources

Each Velero operation – on-demand backup, scheduled backup, restore – is a custom resource, defined with a Kubernetes [Custom Resource Definition (CRD)](https://kubernetes.io/docs/concepts/api-extension/custom-resources/#customresourcedefinitions) and stored in [etcd](https://github.com/coreos/etcd).  Each time these operations are run, an equivalent Kubernetes object is created and saved to storage.

Because of this Kubernetes native way that Velero operates, you are not restricted to backing up the entire etcd. You can back up or restore all objects in a Tanzu Community Edition cluster, or you can also filter what object to operate on by type, namespace, and/or label.

### CLI

Once the Velero package is installed, you will need to have the Velero CLI to run operations on the command line. Please see this documentation for how to install it: [Velero CLI install](https://velero.io/docs/v1.6/basic-install/#install-the-cli).

### Storage

Velero needs an object storage where to save all resource backups and the information about snapshot backups.

Velero treats object storage as the source of truth. It continuously checks to see that the correct backup resources are always present. If there is a properly formatted backup file in the storage bucket, but no corresponding backup resource in the Kubernetes API, Velero synchronizes the information from object storage to Kubernetes. This allows restore functionality to work in a cluster migration scenario, where the original backup objects do not exist in the new cluster. Likewise, if a backup object exists in Kubernetes but not in object storage, it will be deleted from Kubernetes since the backup tarball no longer exists. To learn more, please see the documentation: [How Velero works](https://velero.io/docs/v1.6/how-velero-works/).

### Server

Velero runs on the cluster as a deployment alongside installed plugins that are specific to storage providers for backup and snapshot operations. It also includes controllers that process the custom resources to perform backups, restores, and all related operations.

## Providers

The Tanzu Community Edition Velero package provides support for these providers out of the box with minimal configuration.

| Provider                          | Object Store        | Volume Snapshotter           | Plugin Provider Repo                    | Setup Instructions            |
|-----------------------------------|---------------------|------------------------------|-----------------------------------------|-------------------------------|
| [Amazon Web Services (AWS)](https://aws.amazon.com)    | AWS S3              | AWS EBS                      | [Velero plugin for AWS](https://github.com/vmware-tanzu/velero-plugin-for-aws)              | [AWS Plugin Setup](https://github.com/vmware-tanzu/velero-plugin-for-aws#setup)        |
| [Microsoft Azure](https://azure.com)                                       | Azure Blob Storage  | Azure Managed Disks          | [Velero plugin for Microsoft Azure](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure) | [Azure Plugin Setup](https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure#setup)      |
| [VMware vSphere](https://github.com/vmware-tanzu/velero-plugin-for-vsphere) | On-premise                |  vSphere Volumes            | [VMware vSphere](https://github.com/vmware-tanzu/velero-plugin-for-vsphere)                    | [vSphere Plugin Setup](https://github.com/vmware-tanzu/velero-plugin-for-vsphere#velero-plugin-for-vsphere-installation-and-configuration-details)

Some other third-party storage providers, like MinIO, DigitalOcean, and others, support the same S3 API that the **AWS Velero plugin** uses.  For more information please see: [S3-Compatible object store providers for Velero](https://velero.io/docs/v1.6/supported-providers/#s3-compatible-object-store-providers).

## Configuration

The following configuration values can be set to customize the Velero installation for the different components.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Required | The namespace in which to deploy Velero. |
| `features` | Required | A comma separated string to set some Velero function enable. Current supported features included: EnableCSI, EnableAPIGroupVersions, EnableUploadProgress. Default to "". |

### Storage settings

#### Global configurations for storage

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `backupStorageLocation.create` | Required | create indicate whether to create BackupStorageLocation. Required. Default to true. |
| `backupStorageLocation.name` | Optional | The name of the Backup Storage Location. |
| `backupStorageLocation.provider` | Optional | The cloud provider to use. One of: `aws`, `azure` and `vsphere`. |
| `backupStorageLocation.default` | Optional | Indicates if this location is the default backup storage location. |
| `backupStorageLocation.objectStorage.bucket` | Optional | The storage bucket where backups are to be uploaded. |
| `backupStorageLocation.objectStorage.prefix` | Optional | The directory inside a storage bucket where backups are to be uploaded. |

#### AWS storage

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `backupStorageLocation.configAWS.region` | Optional | The AWS region where the S3 bucket is located. |

#### Azure storage

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `backupStorageLocation.configAzure.resourceGroup` | Optional | The name of the resource group containing the storage account for this backup storage location. |
| `backupStorageLocation.configAzure.storageAccount` | Optional | The name of the storage account for this backup storage location. |
| `backupStorageLocation.configAzure.storageAccountKeyEnvVar` | Optional | Required if using a storage account access key to authenticate rather than a service principal. |
| `backupStorageLocation.configAzure.subscriptionId` | Optional | The the ID of the subscription for this backup storage location. |

### Volume snapshot settings

#### Global configurations for snapshotting

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `volumeSnapshotLocation.snapshotsEnabled` | Required | Indicates whether to create a volumesnapshotlocation CR. If false => disable the snapshot feature. |
| `volumeSnapshotLocation.name` | Optional | The name of the volume snapshot location where snapshots are being taken. |
| `volumeSnapshotLocation.spec.provider` | Optional | The name for the volume snapshot provider. Required if snapshots are enabled. Valid values are `aws` and `azure` |

#### AWS volumes

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `volumeSnapshotLocation.spec.configAWS.region` | Optional | The AWS region where the volumes/snapshots are located |
| `volumeSnapshotLocation.spec.configAWS.profile` | Optional | The AWS profile within the credentials file to use for the volume snapshot location. Default is "default". |

#### Azure volumes

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `volumeSnapshotLocation.spec.configAzure.apiTimeout` | Optional | Indicates how long to wait for an Azure API request to complete before timeout. Defaults to 2m0s. |
| `volumeSnapshotLocation.spec.configAzure.resourceGroup` | Optional | The name of the resource group where volume snapshots should be stored, if different from the cluster's resource group. |
| `volumeSnapshotLocation.spec.configAzure.subscriptionId` | Optional | The ID of the subscription where volume snapshots should be stored, if different from the cluster's subscription. Requires "resourceGroup" to also be set. |
| `volumeSnapshotLocation.spec.configAzure.incremental` | Optional | Azure offers the option to take full or incremental snapshots of managed disks. Set this parameter to true, to take incremental snapshots. If the parameter is omitted or set to false, full snapshots are taken (default). |

#### vSphere volumes

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `volumeSnapshotLocation.spec.configvSphere.region` | Optional | Region is the S3 region where the volumes/snapshots are located. Defaults to minio. |
| `volumeSnapshotLocation.spec.configvSphere.bucket` | Optional | Bucket is the name of the bucket to store volumes/snapshots in. Defaults to velero. |

### vSphere specific

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `vsphere.create` | Required | Whether to deploy vSphere plugin. |
| `vsphere.namespace` | Optional | The namespace where vSphere secret is located. Defaults to "kube-system". |
| `vsphere.clusterName` | Optional | The name found in k8s current-context. Defaults to "tkg-mgmt-vc" |
| `vsphere.server` | Optional | vSphere VC server IP address in CSI vSphere configuration file's [VirtualCenter x.x.x.x] section tag. |
| `vsphere.username` | Optional | user name in CSI vSphere configuration file's [VirtualCenter x.x.x.x] section. |
| `vsphere.password` | Optional | password value in [VirtualCenter x.x.x.x] section. |
| `vsphere.datacenter` | Optional | datacenters value in [VirtualCenter x.x.x.x] section. |
| `vsphere.publicNetwork` | Optional | public-network value in [Network] section. |

### MinIO

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `vsphere.deployDefaultMinio` | Required | Whehter to deploy default MinIO to provide OSS for velero. Defaults to false. |
| `vsphere.namespace` | Optional | The namespace to delploy default MinIO. |
| `vsphere.accessKey` | Optional | MinIO's access key. |
| `vsphere.secretAccessKey` | Optional | MinIO's secret access key. |

### Images

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `images.update` | Required | Whether to update images. Defaults to false. |
| `images.velero` | Optional | Value to replace velero image. |
| `images.minio` | Optional | Value to replace minio image. |
| `images.minioClient` | Optional | Value to replace minio client image. |

### Advanced

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `rbac.create` | Required | Whether to create the Velero Role and RoleBinding to give all permissions to the namespace to Velero.|
| `rbac.name` | Optional |  A new name for the cluster RolBinding. Default is `velero`. |
| `rbac.clusterAdministrator` | Optional | Whether to create the ClusterRoleBinding to give administrator permissions to Velero. `rbac.create` must also be set to `true`.|
| `rbac.roleRefName` | Optional | Name of the cluster role to reference. Default is `cluster-admin`.|
| `rbac.clusterRoleAPIGroups` | Optional |  The name of the API groups that contain the resources for the cluster role.|
| `rbac.clusterRoleVerbs` | Optional |  The set of verbs that apply to the secret resources contained in this rule. |
| `serviceAccount.name` | Optional |  The name of the ServiceAccount the RoleBinding should reference. |
| `serviceAccount.annotations` | Optional |  Annotations for the ServiceAccount the RoleBinding should reference. |
| `serviceAccount.labels` | Optional |  Labels for the ServiceAccount the RoleBinding should reference. |
| `restic.create` | Required | Whether to deploy the restic daemonset. |
| `restic.defaultVolumesToRestic` | Optional | Bool flag to configure Velero server to use restic by default to backup all pod volumes on all backups. |
| `restic.defaultResticPruneFrequency` | Optional | How often 'restic prune' is run for restic repositories by default. |
| `restic.cpuLimit` | Optional | CPU limit for restic pod. A value of "0" is treated as unbounded. (default "1000m"). |
| `restic.cpuRequest` | Optional | CPU request for restic pod. A value of "0" is treated as unbounded. (default "500m"). |
| `restic.memoryLimit` | Optional | Memory limit for restic pod. A value of "0" is treated as unbounded. (default "1Gi"). |
| `restic.memoryRequest` | Optional | Memory request for restic pod. A value of "0" is treated as unbounded. (default "512Mi"). |
| `restic.hostPath` | Optional | The path on host, where Restic DaemonSet reads pod volumes from. Optional. (default "/var/lib/kubelet/pods") |
| `restic.restoreHelperImage` | Optional | value to replace velero Restic restore helper image. Optional. (default "velero/velero-restic-restore-helper:v1.9.2") |
| `plugins` | Required | plugins list to included in Velero deployment InitContainers. User can add customized plugins too. Allow to be an empty array. |
| `environmentVariables` | Required | Environment variables for Velero Deployment and Restic Daemonset. Default to empty array. |

## Installation

The Velero package can easily be installed to a cluster. However, for the package to function, you must configure settings for your cloud provider.

### AWS Configuration

To configure Velero for AWS, you will first need to setup your AWS account. Your account will need the
following:

* An S3 bucket to store the backups
* IAM user with permissions to access EC2 and S3

There are multiple ways to configure your AWS account. Detailed instructions and other options can be found in the [velero-plugin-for-aws](https://github.com/vmware-tanzu/velero-plugin-for-aws#setup) setup instructions. This guide follows
the first option in those instructions.

> You will also need the AWS CLI to complete these steps. Follow the [user guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html) to install it if needed.

1. Create an S3 bucket to store backups.

    ```shell
    export BUCKET=<< bucket name >>
    export REGION=<< region >>
    aws s3api create-bucket \
        --bucket ${BUCKET} \
        --region ${REGION}
    ```

1. Create an IAM user. This user will be explicitly for use with Velero, and given the appropriate permissions to perform backups.

    ```shell
    aws iam create-user --user-name velero
    ```

1. Create the policy that gives the necessary S3 and EC2 permissions.

    ```shell
    cat > velero-policy.json <<EOF
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": [
                    "ec2:DescribeVolumes",
                    "ec2:DescribeSnapshots",
                    "ec2:CreateTags",
                    "ec2:CreateVolume",
                    "ec2:CreateSnapshot",
                    "ec2:DeleteSnapshot"
                ],
                "Resource": "*"
            },
            {
                "Effect": "Allow",
                "Action": [
                    "s3:GetObject",
                    "s3:DeleteObject",
                    "s3:PutObject",
                    "s3:AbortMultipartUpload",
                    "s3:ListMultipartUploadParts"
                ],
                "Resource": [
                    "arn:aws:s3:::${BUCKET}/*"
                ]
            },
            {
                "Effect": "Allow",
                "Action": [
                    "s3:ListBucket"
                ],
                "Resource": [
                    "arn:aws:s3:::${BUCKET}"
                ]
            }
        ]
    }
    EOF
    ```

    Apply the policy to the `velero` IAM user.

    ```shell
    aws iam put-user-policy \
      --user-name velero \
      --policy-name velero \
      --policy-document file://velero-policy.json
    ```

1. Create an access key for the `velero` user.

    ```shell
    aws iam create-access-key --user-name velero
    ```

    The response will look like the following:

    ```json
    {
        "AccessKey": {
            "UserName": "velero",
            "AccessKeyId": "AKIA4CRA12345EXAMPLE",
            "Status": "Active",
            "SecretAccessKey": "RbJSykwB1Z3c094BKgfJi2YzBCG12345EXAMPLE",
            "CreateDate": "2021-09-23T15:32:15+00:00"
        }
    }
    ```

    > Be sure to save the `AccessKeyId` and `SecretAccessKey` values in a safe place. This is the only time that you will
    > have access to these values.

1. Create a `values.yaml` file to configure the package to use AWS. Values that you will need to provide are:

    * Bucket name
    * Region
    * AWS Access Key
    * AWS Secret Access Key

    ```shell
    export AWS_ACCESS_KEY=<< AWS ACCESS KEY >>
    export AWS_SECRET_ACCESS_KEY=<< AWS SECRET ACCESS KEY >>

    cat > values.yaml <<EOF
    ---
    namespace: velero
    restic:
      create: false
    credential:
      useDefaultSecret: true
      name: cloud-credentials
      secretContents:
        cloud: |
          [default]
          aws_access_key_id=${AWS_ACCESS_KEY}
          aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}
    backupStorageLocation:
      create: true
      name: backup
      spec:
        provider: aws
        default: true
        objectStorage:
          bucket: ${BUCKET}
          prefix: backup
        configAWS:
          region: ${REGION}
    volumeSnapshotLocation:
      snapshotsEnabled: true
      name: default
      spec:
        provider: aws
        configAWS:
          region: ${REGION}
    vsphere:
      create: false
    minio:
      deployDefaultMinio: false
    images:
      update: false
    features: ""
    environmentVariables:
    plugins:
      - name: velero-plugin-for-aws
        image: velero/velero-plugin-for-aws:v1.5.1
    EOF
    ```

1. Install the Velero package.

    ```sh
    tanzu package install velero --package-name velero.community.tanzu.vmware.com --version 1.9.2 --values-file values.yaml
    ```

1. Verify that the Velero package was properly installed.

    ```sh
    tanzu package installed list
    | Retrieving installed packages...
      NAME    PACKAGE-NAME                       PACKAGE-VERSION  STATUS
      velero  velero.community.tanzu.vmware.com  1.9.2            Reconcile succeeded
    ```

### vSphere configuration

This is for using Velero in TCE vSphere environment, so first, need to create a TCE cluster.

For creating vSphere TCE cluster, please reference to document: [Prepare to Deploy a Management Cluster to vSphere](https://tanzucommunityedition.io/docs/v0.11/vsphere/)

Please be aware the vCenter IAAS infrastructure should be created before going through the guide.

For vSphere and vCenter guidance, please reference to: [VMware vSphere Document](https://docs.vmware.com/cn/VMware-vSphere/index.html)

1. Due to vSphere environment doesn't provide a S3 compatible object storage, the Carvel package provides a way to deploy a default MinIO for that job. Just set DEPLOY_MINIO to true and provide AWS_ACCESS_KEY and AWS_SECRET_ACCESS_KEY in the next step. Of course, you can also choose to use your own existing S3 compatible object storage.

1. Create a `values.yaml` file to configure the package to use vSphere. Values that you will need to provide are:

    * S3 storage provider's Bucket name
    * S3 storage provider's Region
    * S3 storage provider's Access Key
    * S3 storage provider's Secret Access Key
    * Namespace that vsphere is installed in the workload cluster
    * Workload cluster name. If your environment already set the vSphere workload cluster as the default cluster to kubectl, cluster name can be retrieved by this command: `kubectl config current-context | awk -F@ '{ print $1 }'`
    * vCenter server's IP address
    * vCenter login user name
    * vCenter login password
    * Datacenter created in vCenter to deploy workload cluster
    * Public network used by vCenter
    * Whether to deploy a default MinIO as S3 provider. If deployed, the default MinIO will use provided S3 access key and secret key.

    ```shell
    export AWS_ACCESS_KEY=< S3 ACCESS KEY >
    export AWS_SECRET_ACCESS_KEY=< S3 SECRET ACCESS KEY >
    export BUCKET=< S3 BUCKET >
    export REGION=< S3 REGION >
    export VSPHERE_NAMESPACE=< VSPHERE NAMESPACE IN CLUSTER >
    export WORKLOAD_CLUSTER_NAME=`kubectl config current-context | awk -F@ '{ print $1 }'`
    export VC_SERVER=< VCENTER SERVER IP >
    export VC_USERNAME=< VCENTER USERNAME >
    export VC_PASSWORD=< VCENTER PASSWORD >
    export VC_DATACENTER=< VCENTER DATACENTER >
    export VC_PUBLIC_NETWORK=< VCENTER PUBLIC NETWORK >
    export DEPLOY_MINIO=< DEPLOY MINIO >

    cat > values.yaml <<EOF
    ---
    namespace: velero
    restic:
      create: false
    credential:
      useDefaultSecret: true
      name: cloud-credentials
      secretContents:
        cloud: |
          [default]
          aws_access_key_id=${AWS_ACCESS_KEY}
          aws_secret_access_key=${AWS_SECRET_ACCESS_KEY}
    backupStorageLocation:
      create: true
      name: backup
      spec:
        provider: aws
        default: true
        objectStorage:
          bucket: ${BUCKET}
        configAWS:
          region: ${REGION}
          s3ForcePathStyle: true
          s3Url: ${S3URL}
    volumeSnapshotLocation:
      snapshotsEnabled: true
      name: default
      spec:
        provider: vsphere
        configAWS:
          region: ${REGION}
          bucket: ${BUCKET}
    vsphere:
      create: true
      namespace: ${VSPHERE_NAMESPACE}
      clusterName: ${WORKLOAD_CLUSTER_NAME}
      server: ${VC_SERVER}
      username: ${VC_USERNAME}
      password: ${VC_PASSWORD}
      datacenter: ${VC_DATACENTER}
      publiceNetwork: ${$VC_PUBLIC_NETWORK}
    minio:
      deployDefaultMinio: ${DEPLOY_MINIO}
      namespace: ${AWS_ACCESS_KEY}
      accessKey: ${AWS_ACCESS_KEY}
      secretAccessKey: minio123
    images:
      update: false
    features: ""
    environmentVariables:
    plugins:
      - name: velero-plugin-for-aws
        image: velero/velero-plugin-for-aws:v1.5.1
      - name: velero-plugin-for-vsphere
        image: vsphereveleroplugin/velero-plugin-for-vsphere:v1.4.0
    EOF
    ```

1. Install the Velero package.

    ```sh
    tanzu package install velero --package-name velero.community.tanzu.vmware.com --version 1.9.2 --values-file values.yaml
    ```

1. Verify that the Velero package was properly installed.

    ```sh
    tanzu package installed list
    | Retrieving installed packages...
      NAME    PACKAGE-NAME                       PACKAGE-VERSION  STATUS
      velero  velero.community.tanzu.vmware.com  1.9.2            Reconcile succeeded
    ```

## Usage Example

This walkthrough guides you through an example disaster recovery scenario that leverages the Velero package. You must deploy the package before attempting this walkthrough.

⚠️ Note: For more advanced use cases and documentation, see the official Velero [documentation](https://velero.io/docs/latest/).

In the following steps, you will simulate a disaster scenario. Specifically, you will deploy a stateless workload, create a backup, delete the workload, and restore it from the backup.

1. Download the Velero CLI from the GitHub [releases](https://github.com/vmware-tanzu/velero/releases/latest) page. The following steps assume you have installed Velero into your PATH.

1. Create a new namespace for this example:

    ```sh
    kubectl create ns velero-example
    ```

1. Deploy a sample workload into the new namespace:

    ```sh
    kubectl create deploy -n velero-example nginx --image=nginx
    ```

1. Verify the workload is up and running:

    ```sh
    kubectl get pods -n velero-example
    ```

    The output should be similar to the following:

    ```sh
    NAME                     READY   STATUS    RESTARTS   AGE
    nginx-6799fc88d8-mm47k   1/1     Running   0          7s
    ```

1. Create a backup of the `velero-example` namespace:

    ```sh
    velero create backup velero-example --include-namespaces velero-example
    ```

1. Verify the backup completed successfully:

    ```sh
    velero describe backup velero-example
    ```

    The output shows the "Phase" of the backup, which should be `Completed`.

1. Delete the `velero-example` namespace to simulate a disaster scenario:

    ```sh
    kubectl delete ns velero-example
    ```

1. Verify that the namespace has been deleted:

    ```sh
    kubectl get ns
    ```

1. Restore the namespace from the velero backup:

    ```sh
    velero create restore --from-backup velero-example
    ```

1. Validate that the `velero-example` namespace has been restored:

    ```sh
    kubectl get ns velero-example
    ```

1. Validate that the workload has been restored:

    ```sh
    kubectl get pods -n velero-example
    ```

## Limitation

Velero Carvel package is still not full fledged yet. The followings are features not included but supported by velero client:

* Multiple BackupStorageLocation provision.
