DRAFT WIP DRAFT WIP

<!-- Taken from: https://github.com/vmware-tanzu-private/tkg-docs/tree/main/tkg-docs.vmware.com/aws -->
# Prepare to Deploy a Management or Stand-alone Cluster to Amazon EC2

This topic explains how to prepare Amazon EC2 before you deploy a management or stand-alone cluster.

Before you can use the Tanzu CLI or installer interface to deploy a management or stand-alone cluster, you must prepare the bootstrap machine on which you run the Tanzu CLI and set up your Amazon Web Services Account (AWS) account. To enable Tanzu Kubernetes Grid VMs to launch on Amazon EC2, you must configure your AWS account credentials and then provide the public key part of an SSH key pair to Amazon EC2 for every region in which you plan to deploy management clusters.

**Before you begin**

- Ensure the Tanzu CLI is installed locally. See [Install the Tanzu CLI](../install-cli.md).
- Install [`jq`]( https://stedolan.github.io/jq/download/) is installed locally. The AWS CLI uses `jq` to process JSON when creating SSH key pairs. 
- Install the [AWS CLI]( https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- You have the access key and access key secret for an active AWS account. For more information, see [AWS Account and Access Keys](https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html) in the AWS documentation. 


**Procedure to Configure AWS Account Credentials and SSH Key**


1. Create an access key and access key secret for your active AWS account. For more information, see [AWS Account and Access Keys](https://docs.aws.amazon.com/powershell/latest/userguide/).
2. Configure AWS account credentials using one of the following methods:
    - Specify your AWS account credentials statically in local environment variables.
    - Configure a credentials profile using the AWS configure command



To configure your AWS account credentials and SSH key pair, perform the following steps.

### <a id="account-setup"></a> Configure AWS Credentials

Tanzu Kubernetes Grid uses the default AWS credentials provider chain.
You must set your account credentials to create an SSH key pair for the region where you plan to deploy Tanzu Kubernetes Grid clusters.

To deploy your management cluster on AWS, you have several options for configuring the AWS account used to access EC2.

 - You can specify your AWS account credentials statically in local environment variables. Set the following environment variables for your AWS account:

export AWS_ACCESS_KEY_ID=<em>aws_access_key, where
aws_access_key is your AWS access key.

- <code>export AWS_SECRET_ACCESS_KEY=<em>aws_access_key_secret</em></code>, where <code><em>aws_access_key_secret</em></code> is your AWS access key secret.

- <code>export AWS_SESSION_TOKEN=<em>aws_session_token</em></code>, where
    <code><em>aws_session_token</em></code> is the AWS session token granted to your account. You only need to specify this variable if you are required to use a temporary access key. For more information about using temporary access keys, see [Understanding and getting your AWS credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#temporary-access-keys).

- <code>export AWS_REGION=<em>aws_region</em></code>, where
    <code><em>aws_region</em></code> is the AWS region in which you intend to deploy the cluster. For example, `us-west-2`.

    For the full list of AWS regions, see [AWS Service Endpoints](https://docs.aws.amazon.com/general/latest/gr/rande.html). In addition to the regular AWS regions, you can also specify the `us-gov-east` and
      `us-gov-west` regions in AWS GovCloud.
 - You can use a credentials profile, which you can store in a shared credentials file, such as `~/.aws/credentials`, or a shared config file, such as `~/.aws/config`. You can manage profiles by using the `aws configure` command.

### <a id="aws-account-env-vars"></a> Local Environment Variables

One option for configuring AWS credentials is to set local environment variables on your bootstrap machine. To use local environment variables, 

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
