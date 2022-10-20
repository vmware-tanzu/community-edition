# AWS Workload Cluster Template

When you deploy workload clusters to Amazon Web Services (AWS), you must specify options in the cluster configuration file to connect to your AWS account and identify the resources that the cluster will use. You can also specify the sizes for the control plane and worker node VMs, distribute nodes across availability zones, and share VPCs between clusters.

The template below includes all of the options that are relevant to deploying workload clusters on AWS. You can copy this template and update it to deploy workload clusters to AWS.

Mandatory options are uncommented. Optional settings are commented out.  Default values are included where applicable.

```yaml
#! ---------------------------------------------------------------------
#! Cluster creation basic configuration
#! ---------------------------------------------------------------------

#! CLUSTER_NAME:
CLUSTER_PLAN: dev
NAMESPACE: default
CNI: antrea

#! ---------------------------------------------------------------------
#! Node configuration
#! AWS-only MACHINE_TYPE settings override cloud-agnostic SIZE settings.
#! ---------------------------------------------------------------------

# SIZE:
# CONTROLPLANE_SIZE:
# WORKER_SIZE:
CONTROL_PLANE_MACHINE_TYPE: t3.small
NODE_MACHINE_TYPE: m5.large
# CONTROL_PLANE_MACHINE_COUNT: 1
# WORKER_MACHINE_COUNT: 1
# WORKER_MACHINE_COUNT_0:
# WORKER_MACHINE_COUNT_1:
# WORKER_MACHINE_COUNT_2:

#! ---------------------------------------------------------------------
#! AWS Configuration
#! ---------------------------------------------------------------------

AWS_REGION:
AWS_NODE_AZ: ""
# AWS_NODE_AZ_1: ""
# AWS_NODE_AZ_2: ""
# AWS_VPC_ID: ""
# AWS_PRIVATE_SUBNET_ID: ""
# AWS_PUBLIC_SUBNET_ID: ""
# AWS_PUBLIC_SUBNET_ID_1: ""
# AWS_PRIVATE_SUBNET_ID_1: ""
# AWS_PUBLIC_SUBNET_ID_2: ""
# AWS_PRIVATE_SUBNET_ID_2: ""
# AWS_VPC_CIDR: 10.0.0.0/16
# AWS_PRIVATE_NODE_CIDR: 10.0.0.0/24
# AWS_PUBLIC_NODE_CIDR: 10.0.1.0/24
# AWS_PRIVATE_NODE_CIDR_1: 10.0.2.0/24
# AWS_PUBLIC_NODE_CIDR_1: 10.0.3.0/24
# AWS_PRIVATE_NODE_CIDR_2: 10.0.4.0/24
# AWS_PUBLIC_NODE_CIDR_2: 10.0.5.0/24
AWS_SSH_KEY_NAME:
BASTION_HOST_ENABLED: true

#! ---------------------------------------------------------------------
#! Machine Health Check configuration
#! ---------------------------------------------------------------------

ENABLE_MHC: true
MHC_UNKNOWN_STATUS_TIMEOUT: 5m
MHC_FALSE_STATUS_TIMEOUT: 12m

#! ---------------------------------------------------------------------
#! Common configuration
#! ---------------------------------------------------------------------

# TKG_CUSTOM_IMAGE_REPOSITORY: ""
# TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE: ""

# TKG_HTTP_PROXY: ""
# TKG_HTTPS_PROXY: ""
# TKG_NO_PROXY: ""

ENABLE_AUDIT_LOGGING: true
ENABLE_DEFAULT_STORAGE_CLASS: true

CLUSTER_CIDR: 100.96.0.0/11
SERVICE_CIDR: 100.64.0.0/13

# OS_NAME: ""
# OS_VERSION: ""
# OS_ARCH: ""

#! ---------------------------------------------------------------------
#! Autoscaler configuration
#! ---------------------------------------------------------------------

ENABLE_AUTOSCALER: false
# AUTOSCALER_MAX_NODES_TOTAL: "0"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_ADD: "10m"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_DELETE: "10s"
# AUTOSCALER_SCALE_DOWN_DELAY_AFTER_FAILURE: "3m"
# AUTOSCALER_SCALE_DOWN_UNNEEDED_TIME: "10m"
# AUTOSCALER_MAX_NODE_PROVISION_TIME: "15m"
# AUTOSCALER_MIN_SIZE_0:
# AUTOSCALER_MAX_SIZE_0:
# AUTOSCALER_MIN_SIZE_1:
# AUTOSCALER_MAX_SIZE_1:
# AUTOSCALER_MIN_SIZE_2:
# AUTOSCALER_MAX_SIZE_2:

#! ---------------------------------------------------------------------
#! Antrea CNI configuration
#! ---------------------------------------------------------------------

# ANTREA_NO_SNAT: false
# ANTREA_TRAFFIC_ENCAP_MODE: "encap"
# ANTREA_PROXY: false
# ANTREA_POLICY: true
# ANTREA_TRACEFLOW: false
```
