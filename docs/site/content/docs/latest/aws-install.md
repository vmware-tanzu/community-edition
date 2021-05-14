### <a id="iaas-aws"></a> Configure a Amazon EC2 Infrastructure Provider

1. In the **IaaS Provider** section, enter credentials for your Amazon EC2 account. You have two options:
    - In the **AWS Credential Profile** drop-down, you can select an already existing AWS credential profile. If you select a profile, the access key and session token information configured for your profile are passed to the Installer without displaying actual values in the UI. For information about setting up credential profiles, see [Credential Files and Profiles](aws.md#profiles).
    - Alternately, enter AWS account credentials directly in the **Access Key ID** and **Secret Access Key** fields for your Amazon EC2 account. Optionally specify an AWS session token in **Session Token** if your AWS account is configured to require temporary credentials. For more information on acquiring session tokens, see [Using temporary credentials with AWS resources](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html).
1. In **Region**, select the AWS region in which to deploy the management cluster. If you intend to deploy a production management cluster, this region must have at least three availability zones. This region must also be registered with the SSH key entered in the next field.
1. In **SSH Key Name**, specify the name of an SSH key that is already registered with your Amazon EC2 account and in the region where you are deploying the management cluster. You may have set this up in [Configure AWS Account Credentials and SSH Key](aws.md#register-ssh).
1. If this is the first time that you are deploying a management cluster, select the **Automate creation of AWS CloudFormation Stack** checkbox, and click **Connect**.

   This CloudFormation stack creates the identity and access management (IAM) resources that Tanzu Kubernetes Grid needs to deploy and run clusters on Amazon EC2. For more information, see [Required IAM Resources](aws.md#iam-permissions) in _Prepare to Deploy Management Clusters to Amazon EC2_.

   **IMPORTANT:** The **Automate creation of AWS CloudFormation Stack** checkbox replaces the `clusterawsadm` command line utility that existed in Tanzu Kubernetes Grid v1.1.x and earlier. For existing management and Tanzu Kubernetes clusters initially deployed with v1.1.x or earlier, continue to use the CloudFormation stack that was created by running the `clusterawsadm alpha bootstrap create-stack` command.

   ![Configure the connection to AWS](../images/connect-to-aws.png)
1. If the connection is successful, click **Next**.
1. In the **VPC for AWS** section, do one of the following:

    - To create a new VPC, select **Create new VPC on AWS**, check that the pre-filled CIDR block is available, and click **Next**. If the recommended CIDR block is not available, enter a new IP range in CIDR format for the management cluster to use. The recommended CIDR block for **VPC CIDR** is 10.0.0.0/16.

       ![Create a new VPC](../images/aws-new-vpc.png)

    - To use an existing VPC, select **Select an existing VPC** and select the **VPC ID** from the drop-down menu. The **VPC CIDR** block is filled in automatically when you select the VPC.

       ![Use and existing VPC](../images/aws-existing-vpc.png)

For the next steps, go to [Configure the Management Cluster Settings](#config-mgmt-cluster).