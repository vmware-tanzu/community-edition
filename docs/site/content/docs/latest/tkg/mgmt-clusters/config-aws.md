# Management Cluster Configuration for Amazon EC2

To create a cluster configuration file, you can copy an existing configuration file for a previous deployment to Amazon EC2 and update it. Alternatively, you can create a file from scratch by using an empty template.

## Management Cluster Configuration Template

The template below includes all of the options that are relevant to deploying management clusters on Amazon EC2. You can copy this template and use it to deploy management clusters to Amazon EC2.

- For information about how to update the settings that are common to all infrastructure providers, see [Create a Management Cluster Configuration File](create-config-file.md)
- For information about all configuration file variables, see the [Tanzu CLI Configuration File Variable Reference](../tanzu-config-reference.md).
- For examples of how to configure the vSphere settings, see the sections below the template.

Mandatory options are uncommented. Optional settings are commented out. Default values are included where applicable.

```
#! ---------------------------------------------------------------------
#! Basic cluster creation configuration
#! ---------------------------------------------------------------------

CLUSTER_NAME:
CLUSTER_PLAN: dev
INFRASTRUCTURE_PROVIDER: aws
ENABLE_CEIP_PARTICIPATION: true
# TMC_REGISTRATION_URL:
ENABLE_AUDIT_LOGGING: true
CLUSTER_CIDR: 100.96.0.0/11
SERVICE_CIDR: 100.64.0.0/13

#! ---------------------------------------------------------------------
#! Image repository configuration
#! ---------------------------------------------------------------------

# TKG_CUSTOM_IMAGE_REPOSITORY: ""
# TKG_CUSTOM_IMAGE_REPOSITORY_CA_CERTIFICATE: ""

#! ---------------------------------------------------------------------
#! Proxy configuration
#! ---------------------------------------------------------------------

# TKG_HTTP_PROXY: ""
# TKG_HTTPS_PROXY: ""
# TKG_NO_PROXY: ""

#! ---------------------------------------------------------------------
#! Node configuration
#! AWS-only MACHINE_TYPE settings override cloud-agnostic SIZE settings.
#! ---------------------------------------------------------------------

# SIZE:
# CONTROLPLANE_SIZE:
# WORKER_SIZE:
CONTROL_PLANE_MACHINE_TYPE: t3.small
NODE_MACHINE_TYPE: m5.large
# OS_NAME: ""
# OS_VERSION: ""
# OS_ARCH: ""

#! ---------------------------------------------------------------------
#! AWS configuration
#! ---------------------------------------------------------------------

AWS_REGION:
AWS_NODE_AZ: ""
AWS_ACCESS_KEY_ID:
AWS_SECRET_ACCESS_KEY:
AWS_SSH_KEY_NAME:
BASTION_HOST_ENABLED: true
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


#! ---------------------------------------------------------------------
#! Machine Health Check configuration
#! ---------------------------------------------------------------------

ENABLE_MHC: true
MHC_UNKNOWN_STATUS_TIMEOUT: 5m
MHC_FALSE_STATUS_TIMEOUT: 12m

#! ---------------------------------------------------------------------
#! Identity management configuration
#! ---------------------------------------------------------------------

IDENTITY_MANAGEMENT_TYPE: "oidc"

#! Settings for OIDC
# CERT_DURATION: 2160h
# CERT_RENEW_BEFORE: 360h
# OIDC_IDENTITY_PROVIDER_CLIENT_ID:
# OIDC_IDENTITY_PROVIDER_CLIENT_SECRET:
# OIDC_IDENTITY_PROVIDER_GROUPS_CLAIM: groups
# OIDC_IDENTITY_PROVIDER_ISSUER_URL:
# OIDC_IDENTITY_PROVIDER_SCOPES: email
# OIDC_IDENTITY_PROVIDER_USERNAME_CLAIM: email

#! The following two variables are used to configure Pinniped JWTAuthenticator for workload clusters
# SUPERVISOR_ISSUER_URL:
# SUPERVISOR_ISSUER_CA_BUNDLE_DATA:

#! Settings for LDAP
# LDAP_BIND_DN:
# LDAP_BIND_PASSWORD:
# LDAP_HOST:
# LDAP_USER_SEARCH_BASE_DN:
# LDAP_USER_SEARCH_FILTER:
# LDAP_USER_SEARCH_USERNAME: userPrincipalName
# LDAP_USER_SEARCH_ID_ATTRIBUTE: DN
# LDAP_USER_SEARCH_EMAIL_ATTRIBUTE: DN
# LDAP_USER_SEARCH_NAME_ATTRIBUTE:
# LDAP_GROUP_SEARCH_BASE_DN:
# LDAP_GROUP_SEARCH_FILTER:
# LDAP_GROUP_SEARCH_USER_ATTRIBUTE: DN
# LDAP_GROUP_SEARCH_GROUP_ATTRIBUTE:
# LDAP_GROUP_SEARCH_NAME_ATTRIBUTE: cn
# LDAP_ROOT_CA_DATA_B64:

#! ---------------------------------------------------------------------
#! Antrea CNI configuration
#! ---------------------------------------------------------------------

# ANTREA_NO_SNAT: false
# ANTREA_TRAFFIC_ENCAP_MODE: "encap"
# ANTREA_PROXY: false
# ANTREA_POLICY: true
# ANTREA_TRACEFLOW: false
```

## Amazon EC2 Connection Settings

Specify information about your AWS account and the region and availability zone in which you want to deploy the cluster. If you have set these values as environment variables on the machine on which you run Tanzu CLI commands, you can omit these settings.

The `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are not mandatory. If not provided, the CLI will find them in the AWS default credentials provider chain. If provided, the values for `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` must be base64 encoded.

For example:

```
AWS_REGION: eu-west-1
AWS_NODE_AZ: "eu-west-1a"
AWS_ACCESS_KEY_ID:  <encoded:QUtJQVQ[...]SU82TTM=>
AWS_SECRET_ACCESS_KEY: <encoded:eGN4RHJmLzZ[...]SR08yY2ticQ==>
AWS_SSH_KEY_NAME: default
BASTION_HOST_ENABLED: true
```

## <a id="node-sizes"></a> Configure Node Sizes

The Tanzu CLI creates the individual nodes of Tanzu Kubernetes clusters according to settings that you provide in the configuration file. On Amazon EC2, you can configure all node VMs to have the same predefined configurations or set different predefined configurations for control plane and worker nodes. By using these settings, you can create Tanzu Kubernetes clusters that have nodes with different configurations to the management cluster nodes. You can also create clusters in which the control plane nodes and worker nodes have different configurations.

When you created the management cluster, the instance types for the node machines are set in the `CONTROL_PLANE_MACHINE_TYPE` and `NODE_MACHINE_TYPE` options. By default, these settings are also used for Tanzu Kubernetes clusters. The minimum configuration  is 2 CPUs and 8 GB memory. The list of compatible instance types varies in different regions.

```
CONTROL_PLANE_MACHINE_TYPE: "t3.large"
NODE_MACHINE_TYPE: "m5.large"
```

You can override these settings by using the `SIZE`, `CONTROLPLANE_SIZE` and `WORKER_SIZE` options. To create a Tanzu Kubernetes cluster in which all of the control plane and worker node VMs are the same size, specify the `SIZE` variable. If you set the `SIZE` variable, all nodes will be created with the configuration that you set. For information about the configurations of the different sizes of node instances for Amazon EC2, see [Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/).

```
SIZE: "t3.large"
```

To create a Tanzu Kubernetes cluster in which the control plane and worker node VMs are different sizes, specify the `CONTROLPLANE_SIZE` and `WORKER_SIZE` options.

```
CONTROLPLANE_SIZE: "t3.large"
WORKER_SIZE: "m5.xlarge"
```

You can combine the `CONTROLPLANE_SIZE` and `WORKER_SIZE` options with the `SIZE` option. For example, if you specify `SIZE: "t3.large"` with `WORKER_SIZE: "m5.xlarge"`, the control plane nodes will be set to `t3.large` and worker nodes will be set to `m5.xlarge`.

```
SIZE: "t3.large"
WORKER_SIZE: "m5.xlarge"
```

## Use a New VPC

If you want to deploy a development management cluster with a single control plane node to a new VPC, uncomment and update the following rows related to AWS infrastructure.

```
AWS_REGION:
AWS_NODE_AZ:
AWS_VPC_CIDR:
AWS_PRIVATE_NODE_CIDR:
AWS_PUBLIC_NODE_CIDR:
CONTROL_PLANE_MACHINE_TYPE:
NODE_MACHINE_TYPE:
AWS_SSH_KEY_NAME:
BASTION_HOST_ENABLED:
SERVICE_CIDR:
CLUSTER_CIDR:
```

If you want to deploy a production management cluster with three control plane nodes to a new VPC, also uncomment and update the following variables:

```
AWS_NODE_AZ_1:
AWS_NODE_AZ_2:
AWS_PRIVATE_NODE_CIDR_1:
AWS_PRIVATE_NODE_CIDR_2:
AWS_PUBLIC_NODE_CIDR_1:
AWS_PUBLIC_NODE_CIDR_2:
```

For example, the configuration of a production management cluster on new VPC might look like this:

```
AWS_REGION: us-west-2
AWS_NODE_AZ: us-west-2a
AWS_NODE_AZ_1: us-west-2b
AWS_NODE_AZ_2: us-west-2c
AWS_PRIVATE_NODE_CIDR: 10.0.0.0/24
AWS_PRIVATE_NODE_CIDR_1: 10.0.2.0/24
AWS_PRIVATE_NODE_CIDR_2: 10.0.4.0/24
AWS_PUBLIC_NODE_CIDR: 10.0.1.0/24
AWS_PUBLIC_NODE_CIDR_1: 10.0.3.0/24
AWS_PUBLIC_NODE_CIDR_2: 10.0.5.0/24
AWS_SSH_KEY_NAME: tkg
AWS_VPC_CIDR: 10.0.0.0/16
BASTION_HOST_ENABLED: "true"
CONTROL_PLANE_MACHINE_TYPE: m5.large
NODE_MACHINE_TYPE: m5.large
SERVICE_CIDR: 100.64.0.0/13
CLUSTER_CIDR: 100.96.0.0/11
```

## Use an Existing VPC

If you want to deploy a development management cluster with a single control plane node to an existing VPC, uncomment and update the following rows related to AWS infrastructure.

```
AWS_REGION:
AWS_NODE_AZ:
AWS_PRIVATE_SUBNET_ID:
AWS_PUBLIC_SUBNET_ID:
AWS_SSH_KEY_NAME:
AWS_VPC_ID:
BASTION_HOST_ENABLED:
CONTROL_PLANE_MACHINE_TYPE:
NODE_MACHINE_TYPE:
SERVICE_CIDR:
CLUSTER_CIDR:
```

If you want to deploy a production management cluster with three control plane nodes to an existing VPC, also uncomment and update the following variables:

```
AWS_NODE_AZ_1:
AWS_NODE_AZ_2:
AWS_PRIVATE_SUBNET_ID_1:
AWS_PRIVATE_SUBNET_ID_2:
AWS_PUBLIC_SUBNET_ID_1:
AWS_PUBLIC_SUBNET_ID_2:
```

For example, the configuration of a production management cluster on an existing VPC might look like this:

```
AWS_REGION: us-west-2
AWS_NODE_AZ: us-west-2a
AWS_NODE_AZ_1: us-west-2b
AWS_NODE_AZ_2: us-west-2c
AWS_PRIVATE_SUBNET_ID: subnet-ID
AWS_PRIVATE_SUBNET_ID_1: subnet-ID
AWS_PRIVATE_SUBNET_ID_2: subnet-ID
AWS_PUBLIC_SUBNET_ID: subnet-ID
AWS_PUBLIC_SUBNET_ID_1: subnet-ID
AWS_PUBLIC_SUBNET_ID_2: subnet-ID
AWS_SSH_KEY_NAME: tkg
AWS_VPC_ID: vpc-ID
BASTION_HOST_ENABLED: "true"
CONTROL_PLANE_MACHINE_TYPE: m5.large
NODE_MACHINE_TYPE: m5.large
SERVICE_CIDR: 100.64.0.0/13
CLUSTER_CIDR: 100.96.0.0/11
```

## What to Do Next

After you have finished updating the management cluster configuration file, create the management cluster by following the instructions in [Deploy Management Clusters from a Configuration File](deploy-cli.md#mc-create).
