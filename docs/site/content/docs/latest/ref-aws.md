## Reference information for AWS Account

If you encounter issues deploying a cluster to AWS EC2, review the following reference content:

1. [Resource quotas and ports](#resource-quotas)

2. [Virtual Private Clouds and NAT Gateway Limits](#vpc)

3. [Required Permissions for the AWS Account](#permissions)

## Resource quotas and ports {#resource-quotas}

- Ensure your AWS account has sufficient resource quotas for the following:

   - Virtual Private Cloud (VPC) instances. By default, each management cluster that you deploy creates one VPC and one or three NAT gateways. The default NAT gateway quota is 5 instances per availability zone, per account.
   - Elastic IP (EIP) addresses. The default EIP quota is 5 EIP addresses per region, per account. For more information, see [Amazon VPC Quotas](https://docs.aws.amazon.com/vpc/latest/userguide/amazon-vpc-limits.html) in the AWS documentation and [Resource Usage in Your Amazon Web Services Account](#aws-resources) below:

- Ensure traffic is allowed between your local bootstrap machine and port 6443 of all VMs in the clusters you create. Port 6443 is where the Kubernetes API is exposed.

- For development management clusters that are not configured for high availability, Tanzu Community Edition provisions the following resources:
    - 3 VMs, including a control plane node, a worker node (to run the cluster agent extensions) and, by default, a bastion host. If you specify additional VMs in your node pool, those are provisioned as well.
    - 4 security groups, one for the load balancer and one for each of the initial VMs.
    - 1 private subnet and 1 public subnet in the specified availability zone.
    - 1 public and 1 private route table in the specified availability zone.
    - 1 classic load balancer.
    - 1 internet gateway.
    - 1 NAT gateway in the specified availability zone.
    - By default, 1 EIP, for the NAT gateway, when clusters are deployed in their own VPC. You can optionally share VPCs rather than creating new ones, such as a workload cluster sharing a VPC with its management cluster.

- For production management clusters, which are configured for high availability, Tanzu Community Edition provisions the following resources to support distribution across three availability zones:

    - 3  control plane VMs
    - 3  private and public subnets
    - 3  private and public route tables
    - 3  NAT gateways
    - By default, 3 EIPs, one for each NAT gateway, for clusters deployed in their own VPC. You can optionally share VPCs rather than creating new ones, such as a workload cluster sharing a VPC with its management cluster.

- AWS implements a set of default limits or quotas on these types of resources and allows you to modify the limits. Typically, the default limits are sufficient to get started creating clusters from Tanzu Installer. However, as you increase the number of clusters you are running or the workloads on your clusters, you will encroach on these limits. When you reach the limits imposed by AWS, any attempts to provision that type of resource fail. As a result, Tanzu is unable to create a new cluster, or you might be unable to create additional deployments on your existing clusters. Therefore, regularly assess the limits you have specified in AWS account and adjust them as necessary to fit your business needs.

For information about the sizes of cluster node instances, see
[Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/) in the AWS documentation.

## Virtual Private Clouds and NAT Gateway Limits {#vpc}

- If you create a new Virtual Private Cloud (VPC) when you deploy a management cluster, Tanzu also creates a dedicated NAT gateway for the management cluster or if you deploy a production management cluster, three NAT gateways, one in each of the availability zones.

    If you create a new Virtual Private Cloud (VPC) when you deploy a development management cluster, Tanzu creates a dedicated NAT gateway for the management cluster.  If you deploy a production management cluster,  Tanzu creates three NAT gateways, one in each of the availability zones. In this case, by default, Tanzu creates a new VPC and one or three NAT gateways for each Tanzu cluster that you deploy from that management cluster. By default, AWS allows five NAT gateways per availability zone per account. Consequently, if you always create a new VPC for each cluster, you can create only five development clusters in a single availability zone. If you already have five NAT gateways in use, Tanzu is unable to provision the necessary resources when you attempt to create a new cluster. If you do not want to change the default quotas, to create more than five development clusters in a given availability zone, you must share existing VPCs, and therefore their NAT gateways, between multiple clusters.

- There are three possible scenarios regarding VPCs and NAT gateway usage when you deploy management clusters and workload clusters:

### Create a new VPC and NAT gateway(s) for every management cluster and Tanzu Kubernetes cluster**

   If you deploy a management cluster and use the option to create a new VPC and if you make no modifications to the configuration when you deploy workload clusters from that management cluster, the deployment of each of the workload clusters also creates a VPC and one or three NAT gateways. In this scenario, you can deploy one development management cluster and up to 4 development workload clusters, due to the default limit of 5 NAT gateways per availability zone.

 ### Reuse a VPC and NAT gateway(s) that already exist in your availability zone(s)

   If a VPC already exists in the availability zone(s) in which you are deploying a management cluster, for example a VPC that you created manually or by using tools such as CloudFormation or Terraform, you can specify that the management cluster should use this VPC. In this case, all of the workload clusters that you deploy from that management cluster also use the specified VPC and its NAT gateway(s).

   An existing VPC must be configured with the following networking:

   - Two subnets for development clusters or six subnets for production clusters
   - One NAT gateway for development clusters or three NAT gateways for production clusters
   - One internet gateway and corresponding routing tables

 ### Create a new VPC and NAT gateway(s) for the management cluster and deploy workload clusters that share that VPC and NAT gateway(s)

   If you are starting with an empty availability zone(s), you can deploy a management cluster and use the option to create a new VPC. If you want the workload clusters to share a VPC that Tanzu created, you must modify the cluster configuration when you deploy workload clusters from this management cluster.
## Required Permissions for the AWS Account {#permissions}

Your AWS account must have at least the following permissions:

* [Required IAM Resources](#iam-permissions): Tanzu Kubernetes Grid creates these resources when you deploy a management cluster to your AWS account for the first time.
* [Required Permissions for `tanzu management-cluster create`](#user-permissions): Tanzu Kubernetes
Grid uses these permissions when you run `tanzu management-cluster create` or deploy your management clusters from the installer interface.

### Required IAM Resources {#iam-permissions}

When you deploy your first management cluster to Amazon EC2, you instruct Tanzu to create a CloudFormation stack, `tkg-cloud-vmware-com`, in your AWS account. This CloudFormation stack defines the identity and access management (IAM) resources that Tanzu  uses to deploy and run clusters on Amazon EC2, which includes the following IAM policies, roles, and profiles:

**AWS::IAM::InstanceProfile:** <br>
    - control-plane.tkg.cloud.vmware.com <br>
    - controllers.tkg.cloud.vmware.com <br>
    - nodes.tkg.cloud.vmware.com <br>

**AWS::IAM::ManagedPolicy:** <br>
    - arn:aws:iam::YOUR-ACCOUNT-ID:policy/control-plane.tkg.cloud.vmware.com<br>
      This policy is attached to the control-plane.tkg.cloud.vmware.com IAM role.<br>
    - arn:aws:iam::YOUR-ACCOUNT-ID:policy/nodes.tkg.cloud.vmware.com <br>
  This policy is attached to the control-plane.tkg.cloud.vmware.com and nodes.tkg.cloud.vmware.com IAM roles.<br>
    - arn:aws:iam::YOUR-ACCOUNT-ID:policy/controllers.tkg.cloud.vmware.com <br>
This policy is attached to the controllers.tkg.cloud.vmware.com and control-plane.tkg.cloud.vmware.com IAM roles.

**AWS::IAM::Role:** <br>
    - control-plane.tkg.cloud.vmware.com <br>
    - controllers.tkg.cloud.vmware.com <br>
    - nodes.tkg.cloud.vmware.com <br>

The AWS user that you provide to Tanzu when you create the CloudFormation stack must have permissions to manage IAM resources, such as IAM policies, roles, and instance profiles. You need to create only one CloudFormation stack per AWS account, regardless of whether you use a single or multiple AWS regions for your Tanzu Kubernetes Grid environment.

After Tanzu Kubernetes Grid creates the CloudFormation stack, AWS stores its template as part of the stack. To retrieve the template from CloudFormation, you can navigate to **CloudFormation** > **Stacks** in the AWS console or use the `aws cloudformation get-template` CLI command. For more information about CloudFormation stacks, see [Working with Stacks](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacks.html) in the AWS documentation.

### Required AWS Permissions for `tanzu management-cluster create` {#user-permissions}

The AWS user that you provide to Tanzu when you deploy a management cluster must have at least the following permissions:

* The permissions that are defined in the `control-plane.tkg.cloud.vmware.com`, `nodes.tkg.cloud.vmware.com`, and `controllers.tkg.cloud.vmware.com` IAM polices of the `tkg-cloud-vmware-com` CloudFormation stack. To retrieve these policies from CloudFormation, you can navigate to **CloudFormation** > **Stacks** in the AWS console. For more information, see [Required IAM Resources](#iam-permissions) above.
* If you intend to deploy the management cluster from the installer interface, your AWS user must also have the `"ec2:DescribeInstanceTypeOfferings"` and `"ec2:DescribeInstanceTypes"` permissions. If your AWS user does not currently have these permissions, you can create a custom policy that includes the permissions and attach it to your AWS user.

For example, the `control-plane.tkg.cloud.vmware.com`, `nodes.tkg.cloud.vmware.com`, and `controllers.tkg.cloud.vmware.com` IAM polices include the following permissions:

The `control-plane.tkg.cloud.vmware.com` IAM policy:

```sh
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

```sh
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

```sh
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

## Tag AWS Resources

If both of the following are true, you must add the `kubernetes.io/cluster/YOUR-CLUSTER-NAME=shared` tag to the public subnet or subnets that you intend to use for the management cluster:

* You deploy the management cluster to an existing VPC that was not created by Tanzu.
* You want to create services of type `LoadBalancer` in the management cluster.

Adding the `kubernetes.io/cluster/YOUR-CLUSTER-NAME=shared` tag to the public subnet or subnets enables you to create services of type `LoadBalancer` in the management cluster.
To add this tag, follow the steps below:

1. Gather the ID or IDs of the public subnet or subnets within your existing VPC that you want to use for the management cluster. To deploy a `prod` management cluster, you must provide three subnets.

2. Create the required tag by running the following command:

    ```sh
    aws ec2 create-tags --resources YOUR-PUBLIC-SUBNET-ID-OR-IDS --tags Key=kubernetes.io/cluster/YOUR-CLUSTER-NAME,Value=shared
    ```

    Where:

    * `YOUR-PUBLIC-SUBNET-ID-OR-IDS` is the ID or IDs of the public subnet or subnets that you gathered in the previous step.
    * `YOUR-CLUSTER-NAME` is the name of the management cluster that you want to deploy.

    For example:

    ```sh
    aws ec2 create-tags --resources subnet-00bd5d8c88a5305c6 subnet-0b93f0fdbae3436e8 subnet-06b29d20291797698 --tags Key=kubernetes.io/cluster/my-management-cluster,Value=shared
    ```
















