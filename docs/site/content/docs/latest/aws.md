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
To configure your AWS account credentials and SSH key pair, perform the following steps.

1. Create an access key and access key secret for your active AWS account. For more information, see
[AWS Account and Access Keys](https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html) in the AWS documentation. 

2. Configure AWS credentials using one of the following methods:
    a. One option for configuring AWS credentials is to set local environment variables on your bootstrap machine. To use local environment variables, you specify your AWS account credentials statically in local environment variables. Set the following environment variables for your AWS account:

    ``export AWS_ACCESS_KEY_ID=aws_access_key``

    ``export AWS_SECRET_ACCESS_KEY=aws_access_key_secret``

    ``export AWS_REGION=aws_region``

    or

    b. Configure a credentials profile using the ``AWS configure`` command. Run ``AWS configure`` and enter your access key, access key secret, and region. 
    
    For more information,  [Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

3. For each region that you plan to use with Tanzu Community Edition, create a named key pair, and output a `.pem` file that includes the name. For example, the following command uses `default` and saves the file as `default.pem`.

   ```
   aws ec2 create-key-pair --key-name default --output json | jq .KeyMaterial -r > default.pem
   ```
   To create a key pair for a region that is not the default in your profile, or set locally as `AWS_DEFAULT_REGION`, include the `--region` option.

4. Log in to your Amazon EC2 dashboard and go to **Network & Security** > **Key Pairs** to verify that the created key pair is registered with your account.


