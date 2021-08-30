# Prepare to Deploy a Management or Standalone Cluster to Amazon EC2

This topic explains how to prepare your environment before you deploy a management or standalone cluster on Amazon EC2.

To enable Tanzu Community Edition VMs to launch on Amazon EC2, you must configure your AWS account credentials and then provide the public key part of an SSH key pair to Amazon EC2 for every region in which you plan to deploy a management or standalone cluster.

## Before you begin

- Ensure the Tanzu CLI is installed locally on the bootstrap machine. See [Install the Tanzu CLI](installation-cli.md).
- Install [`jq`]( https://stedolan.github.io/jq/download/) locally on the bootstrap machine. The AWS CLI uses `jq` to process JSON when creating SSH key pairs.
- Install the [AWS CLI]( https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- Ensure you have an active AWS account.
- Review the reference information for your AWS account here: [Reference information for AWS Account](ref-aws.md).


## Procedure

To configure your AWS account credentials and SSH key pair, perform the following steps:

1. Create an access key and access key secret for your active AWS account. For more information, see
[AWS Account and Access Keys](https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html) in the AWS documentation.

2. Configure AWS credentials using one of the following methods:<br>
    a. Set local environment variables on your local bootstrap machine. To use local environment variables, you specify your AWS account credentials statically in local environment variables. Set the following environment variables for your AWS account:

    ```sh
    export AWS_ACCESS_KEY_ID=aws_access_key

    export AWS_SECRET_ACCESS_KEY=aws_access_key_secret

    export AWS_REGION=aws_region
    ```

    or

    b. Configure a credentials profile using the ``AWS configure`` command. Run ``AWS configure`` and enter your access key, access key secret, and region. For more information, see [Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

3. For each region that you plan to use with Tanzu Community Edition, create a named key pair, and output a `.pem` file that includes the name. For example, the following command uses `default` and saves the file as `default.pem`:

   ```sh
   aws ec2 create-key-pair --key-name default --output json | jq .KeyMaterial -r > default.pem
   ```
   To create a key pair for a region that is not the default in your profile, or set locally as `AWS_DEFAULT_REGION`, include the `--region` option.

4. Log in to your Amazon EC2 dashboard and go to **Network & Security** > **Key Pairs** to verify that the created key pair is registered with your account.



