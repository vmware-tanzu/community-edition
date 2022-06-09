# AWS EBS CSI Driver

This package provides cloud storage interface using [aws-ebs-csi-driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver).

## Supported Providers

The following tables shows the providers this package can work with. Other cloud provider support will be added  
in the future.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
|  ✅ |  ❌  | ❌  |  ❌  |

## Configuration

The following configuration values can be set to customize the aws-ebs-csi-driver installation.

### Global

None

### aws-ebs-csi-driver Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `awsEBSCSIDriver.namespace` | Required | The namespace of the Kubernetes cluster in cluster ID. Default value is `null`. |

## Usage Example

To learn more about how to use AWS EBS CSI Driver refer to [AWS EBS CSI Driver website](https://github.com/kubernetes-sigs/aws-ebs-csi-driver)
