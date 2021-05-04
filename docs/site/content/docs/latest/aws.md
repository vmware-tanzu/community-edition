DRAFT WIP DRAFT WIP

<!-- Taken from: https://github.com/vmware-tanzu-private/tkg-docs/tree/main/tkg-docs.vmware.com/aws -->
# Prepare to Deploy a Management or Stand-alone Cluster to Amazon EC2

This topic explains how to prepare Amazon EC2 before you deploy a management or stand-alone cluster.

Before you can use the Tanzu CLI or installer interface to deploy a management cluster, you must prepare the bootstrap machine on which you run the Tanzu CLI and set up your Amazon Web Services Account (AWS) account.


- The Tanzu CLI installed locally. See [Install the Tanzu CLI](../install-cli.md).
- You have the access key and access key secret for an active AWS account. For more information, see [AWS Account and Access Keys](https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html) in the AWS documentation. 
- Your AWS account must have at least the permissions described in [Required Permissions for the AWS Account](#permissions).
- Your AWS account has sufficient resource quotas for the following.
For more information, see [Amazon VPC Quotas](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html) in the AWS documentation and [Resource Usage in Your Amazon Web Services Account](#aws-resources) below:
   - Virtual Private Cloud (VPC) instances. By default, each management cluster that you deploy creates one VPC and one or three NAT gateways. The default NAT gateway quota is 5 instances per availability zone, per account.
   - Elastic IP (EIP) addresses. The default EIP quota is 5 EIP addresses per region, per account.
- Traffic is allowed between your local bootstrap machine and port 6443 of all VMs in the clusters you create. Port 6443 is where the Kubernetes API is exposed.
- Traffic is allowed between your local bootstrap machine and the image repositories listed in the management cluster Bill of Materials (BOM) file, over port 443, for TCP.&#42;
   - The BOM file is under `~/.tanzu/tkg/bom/` and its name includes the Tanzu Kubernetes Grid version, for example `tkg-bom-1.3.0+vmware.1.yaml` for v1.3.0.
   - Run a DNS lookup on all `imageRepository` values to find their IPs, for example `projects.registry.vmware.com/tkg` requires network access to `208.91.0.233`.

- The [AWS CLI]( https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html) installed locally.
- [`jq`]( https://stedolan.github.io/jq/download/) installed locally.

   The AWS CLI uses `jq` to process JSON when creating SSH key pairs. It is also used to prepare the environment or configuration variables when you deploy Tanzu Kubernetes Grid by using the CLI.

&#42;Or see [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md) for installing without external network access.

## <a id="aws-resources"></a> Resource Usage in Your AWS Account

For each cluster that you create, Tanzu Kubernetes Grid provisions a set of resources in your AWS account.

For development management clusters that are not configured for high availability, Tanzu Kubernetes Grid provisions the following resources:

- 3 VMs, including a control plane node, a worker node (to run the cluster agent extensions) and, by default, a bastion host. If you specify additional VMs in your node pool, those are provisioned as well.
- 4 security groups, one for the load balancer and one for each of the initial VMs.
- 1 private subnet and 1 public subnet in the specified availability zone.
- 1 public and 1 private route table in the specified availability zone.
- 1 classic load balancer.
- 1 internet gateway.
- 1 NAT gateway in the specified availability zone.
- By default, 1 EIP, for the NAT gateway, when clusters are deployed in their own VPC. You can optionally share VPCs rather than creating new ones, such as a workload cluster sharing a VPC with its management cluster.

For production management clusters, which are configured for high availability, Tanzu Kubernetes Grid provisions the following resources to support distribution across three availability zones:

- 3  control plane VMs
- 3  private and public subnets
- 3  private and public route tables
- 3  NAT gateways
- By default, 3 EIPs, one for each NAT gateway, for clusters deployed in their own VPC. You can optionally share VPCs rather than creating new ones, such as a workload cluster sharing a VPC with its management cluster.

AWS implements a set of default limits or quotas on these types of resources and allows you to modify the limits. Typically, the default limits are sufficient to get started creating clusters from Tanzu Kubernetes Grid. However, as you increase the number of clusters you are running or the workloads on your clusters, you will encroach on these limits. When you reach the limits imposed by AWS, any attempts to provision that type of resource fail. As a result, Tanzu Kubernetes Grid will be unable to create a new cluster, or you might be unable to create additional deployments on your existing clusters. Therefore, regularly assess the limits you have specified in AWS account and adjust them as necessary to fit your business needs.

For information about the sizes of cluster node instances, see
[Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/) in the AWS documentation.

### <a id="aws-vpc"></a> Virtual Private Clouds and NAT Gateway Limits

If you create a new Virtual Private Cloud (VPC) when you deploy a management cluster, Tanzu Kubernetes Grid also creates a dedicated NAT gateway for the management cluster or if you deploy a production management cluster, three NAT gateways, one in each of the availability zones. In this case, by default, Tanzu Kubernetes Grid creates a new VPC and one or three NAT gateways for each Tanzu Kubernetes cluster that you deploy from that management cluster. By default, AWS allows five NAT gateways per availability zone per account. Consequently, if you always create a new VPC for each cluster, you can create only five development clusters in a single availability zone. If you already have five NAT gateways in use, Tanzu Kubernetes Grid is unable to provision the necessary resources when you attempt to create a new cluster. If you do not want to change the default quotas, to create more than five development clusters in a given availability zone, you must share existing VPCs, and therefore their NAT gateways, between multiple clusters.

There are three possible scenarios regarding VPCs and NAT gateway usage when you deploy management clusters and Tanzu Kubernetes clusters.

- **Create a new VPC and NAT gateway(s) for every management cluster and Tanzu Kubernetes cluster**

   If you deploy a management cluster and use the option to create a new VPC and if you make no modifications to the configuration when you deploy Tanzu Kubernetes clusters from that management cluster, the deployment of each of the Tanzu Kubernetes clusters also creates a VPC and one or three NAT gateways. In this scenario, you can deploy one development management cluster and up to 4 development Tanzu Kubernetes clusters, due to the default limit of 5 NAT gateways per availability zone.

- **Reuse a VPC and NAT gateway(s) that already exist in your availability zone(s)**

   If a VPC already exists in the availability zone(s) in which you are deploying a management cluster, for example a VPC that you created manually or by using tools such as CloudFormation or Terraform, you can specify that the management cluster should use this VPC. In this case, all of the Tanzu Kubernetes clusters that you deploy from that management cluster also use the specified VPC and its NAT gateway(s).

   An existing VPC must be configured with the following networking:

   - Two subnets for development clusters or six subnets for production clusters
   - One NAT gateway for development clusters or three NAT gateways for production clusters
   - One internet gateway and corresponding routing tables

- **Create a new VPC and NAT gateway(s) for the management cluster and deploy Tanzu Kubernetes clusters that share that VPC and NAT gateway(s)**

   If you are starting with an empty availability zone(s), you can deploy a management cluster and use the option to create a new VPC. If you want the Tanzu Kubernetes clusters to share a VPC that Tanzu Kubernetes Grid created, you must modify the cluster configuration when you deploy Tanzu Kubernetes clusters from this management cluster.

For information about how to deploy management clusters that either create or reuse a VPC, see [Deploy Management Clusters with the Installer Interface](deploy-ui.md) and [Deploy Management Clusters to Amazon EC2  CLI](deploy-cli.md).

For information about how to deploy Tanzu Kubernetes clusters that share a VPC that Tanzu Kubernetes Grid created when you deployed the management cluster, see [Deploy a Cluster that Shares a VPC with the Management Cluster](../tanzu-k8s-clusters/aws.md#aws-vpc).

## <a id="permissions"></a> Required Permissions for the AWS Account

Your AWS account must have at least the following permissions:

* [Required IAM Resources](#iam-permissions): Tanzu Kubernetes Grid creates these resources when you deploy a management cluster to your AWS account for the first time.
* [Required Permissions for `tanzu management-cluster create`](#user-permissions): Tanzu Kubernetes
Grid uses these permissions when you run `tanzu management-cluster create` or deploy your management clusters from the installer interface.

### <a id="iam-permissions"></a> Required IAM Resources

When you deploy your first management cluster to Amazon EC2, you instruct Tanzu Kubernetes Grid to create a CloudFormation stack, `tkg-cloud-vmware-com`, in your AWS account. This CloudFormation stack defines the identity and access management (IAM) resources that Tanzu Kubernetes Grid uses to deploy and run clusters on Amazon EC2, which includes the following IAM policies, roles, and profiles:

* `AWS::IAM::InstanceProfile`:
   * `control-plane.tkg.cloud.vmware.com`
   * `controllers.tkg.cloud.vmware.com`
   * `nodes.tkg.cloud.vmware.com`

* `AWS::IAM::ManagedPolicy`:
  * `arn:aws:iam::YOUR-ACCOUNT-ID:policy/control-plane.tkg.cloud.vmware.com`. This policy is attached to the `control-plane.tkg.cloud.vmware.com` IAM role.
  * `arn:aws:iam::YOUR-ACCOUNT-ID:policy/nodes.tkg.cloud.vmware.com`. This policy is attached to the `control-plane.tkg.cloud.vmware.com` and `nodes.tkg.cloud.vmware.com` IAM roles.
  * `arn:aws:iam::YOUR-ACCOUNT-ID:policy/controllers.tkg.cloud.vmware.com`. This policy is attached to the `controllers.tkg.cloud.vmware.com` and `control-plane.tkg.cloud.vmware.com` IAM roles.

* `AWS::IAM::Role`:

   * `control-plane.tkg.cloud.vmware.com`
   * `controllers.tkg.cloud.vmware.com`
   * `nodes.tkg.cloud.vmware.com`

The AWS user that you provide to Tanzu Kubernetes Grid when you create the CloudFormation stack must have permissions to manage IAM resources, such as IAM policies, roles, and instance profiles. You need to create only one CloudFormation stack per AWS account, regardless of whether you use a single or multiple AWS regions for your Tanzu Kubernetes Grid environment.

After Tanzu Kubernetes Grid creates the CloudFormation stack, AWS stores its template as part of the stack. To retrieve the template from CloudFormation, you can navigate to **CloudFormation** > **Stacks** in the AWS console or use the `aws cloudformation get-template` CLI command. For more information about CloudFormation stacks, see [Working with Stacks](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacks.html) in the AWS documentation.

### <a id="user-permissions"></a> Required AWS Permissions for `tanzu management-cluster create`

The AWS user that you provide to Tanzu Kubernetes Grid when you deploy a management cluster must have at least the following permissions:

* The permissions that are defined in the `control-plane.tkg.cloud.vmware.com`, `nodes.tkg.cloud.vmware.com`, and `controllers.tkg.cloud.vmware.com` IAM polices of the `tkg-cloud-vmware-com` CloudFormation stack. To retrieve these policies from CloudFormation, you can navigate to **CloudFormation** > **Stacks** in the AWS console. For more information, see [Required IAM Resources](#iam-permissions) above.
* If you intend to deploy the management cluster from the installer interface, your AWS user must also have the `"ec2:DescribeInstanceTypeOfferings"` and `"ec2:DescribeInstanceTypes"` permissions. If your AWS user does not currently have these permissions, you can create a custom policy that includes the permissions and attach it to your AWS user.

For example, in Tanzu Kubernetes Grid v1.3.0, the `control-plane.tkg.cloud.vmware.com`, `nodes.tkg.cloud.vmware.com`, and `controllers.tkg.cloud.vmware.com` IAM polices include the following permissions:

The `control-plane.tkg.cloud.vmware.com` IAM policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "autoscaling:DescribeAutoScalingGroups",
                "autoscaling:DescribeLaunchConfigurations",
                "autoscaling:DescribeTags",
                "ec2:DescribeInstances",
                "ec2:DescribeImages",
                "ec2:DescribeRegions",
                "ec2:DescribeRouteTables",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeSubnets",
                "ec2:DescribeVolumes",
                "ec2:CreateSecurityGroup",
                "ec2:CreateTags",
                "ec2:CreateVolume",
                "ec2:ModifyInstanceAttribute",
                "ec2:ModifyVolume",
                "ec2:AttachVolume",
                "ec2:AuthorizeSecurityGroupIngress",
                "ec2:CreateRoute",
                "ec2:DeleteRoute",
                "ec2:DeleteSecurityGroup",
                "ec2:DeleteVolume",
                "ec2:DetachVolume",
                "ec2:RevokeSecurityGroupIngress",
                "ec2:DescribeVpcs",
                "elasticloadbalancing:AddTags",
                "elasticloadbalancing:AttachLoadBalancerToSubnets",
                "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
                "elasticloadbalancing:CreateLoadBalancer",
                "elasticloadbalancing:CreateLoadBalancerPolicy",
                "elasticloadbalancing:CreateLoadBalancerListeners",
                "elasticloadbalancing:ConfigureHealthCheck",
                "elasticloadbalancing:DeleteLoadBalancer",
                "elasticloadbalancing:DeleteLoadBalancerListeners",
                "elasticloadbalancing:DescribeLoadBalancers",
                "elasticloadbalancing:DescribeLoadBalancerAttributes",
                "elasticloadbalancing:DetachLoadBalancerFromSubnets",
                "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
                "elasticloadbalancing:ModifyLoadBalancerAttributes",
                "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
                "elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
                "elasticloadbalancing:AddTags",
                "elasticloadbalancing:CreateListener",
                "elasticloadbalancing:CreateTargetGroup",
                "elasticloadbalancing:DeleteListener",
                "elasticloadbalancing:DeleteTargetGroup",
                "elasticloadbalancing:DescribeListeners",
                "elasticloadbalancing:DescribeLoadBalancerPolicies",
                "elasticloadbalancing:DescribeTargetGroups",
                "elasticloadbalancing:DescribeTargetHealth",
                "elasticloadbalancing:ModifyListener",
                "elasticloadbalancing:ModifyTargetGroup",
                "elasticloadbalancing:RegisterTargets",
                "elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
                "iam:CreateServiceLinkedRole",
                "kms:DescribeKey"
            ],
            "Resource": [
                "*"
            ],
            "Effect": "Allow"
        }
    ]
}
```

The `nodes.tkg.cloud.vmware.com` IAM policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "ec2:DescribeInstances",
                "ec2:DescribeRegions",
                "ecr:GetAuthorizationToken",
                "ecr:BatchCheckLayerAvailability",
                "ecr:GetDownloadUrlForLayer",
                "ecr:GetRepositoryPolicy",
                "ecr:DescribeRepositories",
                "ecr:ListImages",
                "ecr:BatchGetImage"
            ],
            "Resource": [
                "*"
            ],
            "Effect": "Allow"
        },
        {
            "Action": [
                "secretsmanager:DeleteSecret",
                "secretsmanager:GetSecretValue"
            ],
            "Resource": [
                "arn:*:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*"
            ],
            "Effect": "Allow"
        },
        {
            "Action": [
                "ssm:UpdateInstanceInformation",
                "ssmmessages:CreateControlChannel",
                "ssmmessages:CreateDataChannel",
                "ssmmessages:OpenControlChannel",
                "ssmmessages:OpenDataChannel",
                "s3:GetEncryptionConfiguration"
            ],
            "Resource": [
                "*"
            ],
            "Effect": "Allow"
        }
    ]
}
```

The `controllers.tkg.cloud.vmware.com` IAM policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "ec2:AllocateAddress",
                "ec2:AssociateRouteTable",
                "ec2:AttachInternetGateway",
                "ec2:AuthorizeSecurityGroupIngress",
                "ec2:CreateInternetGateway",
                "ec2:CreateNatGateway",
                "ec2:CreateRoute",
                "ec2:CreateRouteTable",
                "ec2:CreateSecurityGroup",
                "ec2:CreateSubnet",
                "ec2:CreateTags",
                "ec2:CreateVpc",
                "ec2:ModifyVpcAttribute",
                "ec2:DeleteInternetGateway",
                "ec2:DeleteNatGateway",
                "ec2:DeleteRouteTable",
                "ec2:DeleteSecurityGroup",
                "ec2:DeleteSubnet",
                "ec2:DeleteTags",
                "ec2:DeleteVpc",
                "ec2:DescribeAccountAttributes",
                "ec2:DescribeAddresses",
                "ec2:DescribeAvailabilityZones",
                "ec2:DescribeInstances",
                "ec2:DescribeInternetGateways",
                "ec2:DescribeImages",
                "ec2:DescribeNatGateways",
                "ec2:DescribeNetworkInterfaces",
                "ec2:DescribeNetworkInterfaceAttribute",
                "ec2:DescribeRouteTables",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeSubnets",
                "ec2:DescribeVpcs",
                "ec2:DescribeVpcAttribute",
                "ec2:DescribeVolumes",
                "ec2:DetachInternetGateway",
                "ec2:DisassociateRouteTable",
                "ec2:DisassociateAddress",
                "ec2:ModifyInstanceAttribute",
                "ec2:ModifyNetworkInterfaceAttribute",
                "ec2:ModifySubnetAttribute",
                "ec2:ReleaseAddress",
                "ec2:RevokeSecurityGroupIngress",
                "ec2:RunInstances",
                "ec2:TerminateInstances",
                "tag:GetResources",
                "elasticloadbalancing:AddTags",
                "elasticloadbalancing:CreateLoadBalancer",
                "elasticloadbalancing:ConfigureHealthCheck",
                "elasticloadbalancing:DeleteLoadBalancer",
                "elasticloadbalancing:DescribeLoadBalancers",
                "elasticloadbalancing:DescribeLoadBalancerAttributes",
                "elasticloadbalancing:DescribeTags",
                "elasticloadbalancing:ModifyLoadBalancerAttributes",
                "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
                "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
                "elasticloadbalancing:RemoveTags",
                "autoscaling:DescribeAutoScalingGroups",
                "autoscaling:DescribeInstanceRefreshes",
                "ec2:CreateLaunchTemplate",
                "ec2:CreateLaunchTemplateVersion",
                "ec2:DescribeLaunchTemplates",
                "ec2:DescribeLaunchTemplateVersions",
                "ec2:DeleteLaunchTemplate",
                "ec2:DeleteLaunchTemplateVersions"
            ],
            "Resource": [
                "*"
            ],
            "Effect": "Allow"
        },
        {
            "Action": [
                "autoscaling:CreateAutoScalingGroup",
                "autoscaling:UpdateAutoScalingGroup",
                "autoscaling:CreateOrUpdateTags",
                "autoscaling:StartInstanceRefresh",
                "autoscaling:DeleteAutoScalingGroup",
                "autoscaling:DeleteTags"
            ],
            "Resource": [
                "arn:*:autoscaling:*:*:autoScalingGroup:*:autoScalingGroupName/*"
            ],
            "Effect": "Allow"
        },
        {
            "Condition": {
                "StringLike": {
                    "iam:AWSServiceName": "autoscaling.amazonaws.com"
                }
            },
            "Action": [
                "iam:CreateServiceLinkedRole"
            ],
            "Resource": [
                "arn:*:iam::*:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling"
            ],
            "Effect": "Allow"
        },
        {
            "Condition": {
                "StringLike": {
                    "iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"
                }
            },
            "Action": [
                "iam:CreateServiceLinkedRole"
            ],
            "Resource": [
                "arn:*:iam::*:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing"
            ],
            "Effect": "Allow"
        },
        {
            "Condition": {
                "StringLike": {
                    "iam:AWSServiceName": "spot.amazonaws.com"
                }
            },
            "Action": [
                "iam:CreateServiceLinkedRole"
            ],
            "Resource": [
                "arn:*:iam::*:role/aws-service-role/spot.amazonaws.com/AWSServiceRoleForEC2Spot"
            ],
            "Effect": "Allow"
        },
        {
            "Action": [
                "iam:PassRole"
            ],
            "Resource": [
                "arn:*:iam::*:role/*.tkg.cloud.vmware.com"
            ],
            "Effect": "Allow"
        },
        {
            "Action": [
                "secretsmanager:CreateSecret",
                "secretsmanager:DeleteSecret",
                "secretsmanager:TagResource"
            ],
            "Resource": [
                "arn:*:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*"
            ],
            "Effect": "Allow"
        }
    ]
}
```

## <a id="register-ssh"></a> Configure AWS Account Credentials and SSH Key

To enable Tanzu Kubernetes Grid VMs to launch on Amazon EC2, you must configure your AWS account credentials and then provide the public key part of an SSH key pair to Amazon EC2 for every region in which you plan to deploy management clusters.

To configure your AWS account credentials and SSH key pair, perform the following steps.

### <a id="account-setup"></a> Configure AWS Credentials

Tanzu Kubernetes Grid uses the default AWS credentials provider chain.
You must set your account credentials to create an SSH key pair for the region where you plan to deploy Tanzu Kubernetes Grid clusters.

To deploy your management cluster on AWS, you have several options for configuring the AWS account used to access EC2.

 - You can specify your AWS account credentials statically in local environment variables.
 - You can use a credentials profile, which you can store in a shared credentials file, such as `~/.aws/credentials`, or a shared config file, such as `~/.aws/config`. You can manage profiles by using the `aws configure` command.

### <a id="aws-account-env-vars"></a> Local Environment Variables

One option for configuring AWS credentials is to set local environment variables on your bootstrap machine. To use local environment variables, set the following environment variables for your AWS account:

- <code>export AWS_ACCESS_KEY_ID=<em>aws_access_key</em></code>, where
    <code><em>aws_access_key</em></code> is your AWS access key.

- <code>export AWS_SECRET_ACCESS_KEY=<em>aws_access_key_secret</em></code>, where <code><em>aws_access_key_secret</em></code> is your AWS access key secret.

- <code>export AWS_SESSION_TOKEN=<em>aws_session_token</em></code>, where
    <code><em>aws_session_token</em></code> is the AWS session token granted to your account. You only need to specify this variable if you are required to use a temporary access key. For more information about using temporary access keys, see [Understanding and getting your AWS credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#temporary-access-keys).

- <code>export AWS_REGION=<em>aws_region</em></code>, where
    <code><em>aws_region</em></code> is the AWS region in which you intend to deploy the cluster. For example, `us-west-2`.

    For the full list of AWS regions, see [AWS Service Endpoints](https://docs.aws.amazon.com/general/latest/gr/rande.html). In addition to the regular AWS regions, you can also specify the `us-gov-east` and
      `us-gov-west` regions in AWS GovCloud.

### <a id="profiles"></a> Credential Files and Profiles

As an alternative to using local environment variables, you can store AWS credentials in a shared or local credentials file. An AWS credential file can store multiple accounts as named profiles. The credential files and profiles are applied after local environment variables as part of the AWS default credential provider chain.

To set up credentials files and profiles for your AWS account on the bootstrap machine, you can use the [aws configure](https://docs.aws.amazon.com/cli/latest/reference/configure/index.html) CLI command.

To customize which AWS credential files and profiles to use, you can set the following environment variables:

- <code>export AWS_SHARED_CREDENTIAL_FILE=<em>path_to_credentials_file</em></code> where <code><em>path_to_credentials_file</em></code> is the location and name of the credentials file that contains your AWS access key information. If you do not define this environment variable, the default location and filename is `$HOME/.aws/credentials`.

- <code>export AWS_PROFILE=<em>profile_name</em></code> where <code><em>profile_name</em></code> is the profile name that contains the AWS access key you want to use. If you do not specify a value for this variable, the profile name `default` is used. For more information about using named profiles, see [Named profiles](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html) in AWS documentation.

**NOTE:** Any named profiles that you create in your AWS credentials file appear as selectable options in the **AWS Credential Profile** drop-down in the Tanzu Kubernetes Grid Installer UI for Amazon EC2.

For more information about working AWS credentials and the default AWS credential provider chain, see [Best practices for managing AWS access keys](https://docs.aws.amazon.com/general/latest/gr/aws-access-keys-best-practices.html#more-resources) in the AWS documentation.

### <a id="register-ssh"></a> Register an SSH Public Key with Your AWS Account

After you have set your AWS account credentials using either local environment variables or in a credentials file and profile, you can generate an SSH key pair for your AWS account.

**NOTE**: AWS supports only RSA keys. The keys required by AWS are of a different format to those required by vSphere. You cannot use the same key pair for both vSphere and AWS deployments.

If you do not already have an SSH key pair for the account and region you are using to deploy the management cluster, create one by performing the steps below:

1. For each region that you plan to use with Tanzu Kubernetes Grid, create a named key pair, and output a `.pem` file that includes the name. For example, the following command uses `default` and saves the file as `default.pem`.

   ```
   aws ec2 create-key-pair --key-name default --output json | jq .KeyMaterial -r > default.pem
   ```
   To create a key pair for a region that is not the default in your profile, or set locally as `AWS_DEFAULT_REGION`, include the `--region` option.

1. Log in to your Amazon EC2 dashboard and go to **Network & Security** > **Key Pairs** to verify that the created key pair is registered with your account.

## <a id="aws-tags"></a>Tag AWS Resources

If both of the following are true, you must add the `kubernetes.io/cluster/YOUR-CLUSTER-NAME=shared` tag to the public subnet or subnets that you intend to use for the management cluster:

* You deploy the management cluster to an existing VPC that was not created by Tanzu Kubernetes Grid.
* You want to create services of type `LoadBalancer` in the management cluster.

Adding the `kubernetes.io/cluster/YOUR-CLUSTER-NAME=shared` tag to the public subnet or subnets enables you to create services of type `LoadBalancer` in the management cluster.
To add this tag, follow the steps below:

1. Gather the ID or IDs of the public subnet or subnets within your existing VPC that you want to use for the management cluster. To deploy a `prod` management cluster, you must provide three subnets.

1. Create the required tag by running the following command:

    ```
    aws ec2 create-tags --resources YOUR-PUBLIC-SUBNET-ID-OR-IDS --tags Key=kubernetes.io/cluster/YOUR-CLUSTER-NAME,Value=shared
    ```

    Where:

    * `YOUR-PUBLIC-SUBNET-ID-OR-IDS` is the ID or IDs of the public subnet or subnets that you gathered in the previous step.
    * `YOUR-CLUSTER-NAME` is the name of the management cluster that you want to deploy.

    For example:

    ```
    aws ec2 create-tags --resources subnet-00bd5d8c88a5305c6 subnet-0b93f0fdbae3436e8 subnet-06b29d20291797698 --tags Key=kubernetes.io/cluster/my-management-cluster,Value=shared
    ```

If you want to use services of type `LoadBalancer` in
a Tanzu Kubernetes cluster after you deploy the cluster to a VPC that was not
created by Tanzu Kubernetes Grid, follow the tagging instructions in [Deploy a Cluster to an Existing VPC and Add Subnet Tags (Amazon EC2)](../tanzu-k8s-clusters/aws.md#own-vpc).

## <a id="what-next"></a> What to Do Next

For production deployments, it is strongly recommended to enable identity management for your clusters. For information about the preparatory steps to perform before you deploy a management cluster, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).

If you are using Tanzu Kubernetes Grid in an environment with an external internet connection, once you have set up identity management, you are  ready to deploy management clusters to Amazon EC2.

- [Deploy Management Clusters with the Installer Interface](deploy-ui.md). This is the preferred option for first deployments.
- [Deploy Management Clusters from a Configuration File](deploy-cli.md). This is the more complicated method that allows greater flexibility of configuration.
- If you want to deploy clusters to vSphere and Azure as well as to Amazon EC2, see [Prepare to Deploy Management Clusters to vSphere](vsphere.html) and [Prepare to Deploy Management Clusters to Microsoft Azure](azure.md) for the required setup for those platforms.
