# Deploy Management Clusters from a Configuration File

You can use the Tanzu CLI to deploy a management cluster to vSphere, Amazon Elastic Compute Cloud (Amazon EC2), and Microsoft Azure with a configuration that you specify in a YAML configuration file.

## <a id="prereqs"></a> Prerequisites

Before you can deploy a management cluster, you must make sure that your environment meets the requirements for the target infrastructure provider.

### General Prerequisites

- Make sure that you have met the all of the requirements and followed all of the procedures in [Install the Tanzu CLI and Other Tools](../install-cli.md).
- For production deployments, it is strongly recommended to enable identity management for your clusters. For information about the preparatory steps to perform before you deploy a management cluster, see [Enabling Identity Management in Tanzu Kubernetes Grid](enabling-id-mgmt.md).
- If you want to register your management cluster with Tanzu Mission Control, follow the procedure in [Register Your Management Cluster with Tanzu Mission Control](register_tmc.md).
- If you are deploying clusters in an internet-restricted environment to either vSphere or Amazon EC2, you must also perform the steps in  [Deploying Tanzu Kubernetes Grid in an Internet-Restricted Environment](airgapped-environments.md).
- It is **strongly recommended** to use the Tanzu Kubernetes Grid installer interface rather than the CLI to deploy your first management cluster to a given infrastructure provider. When you deploy a management cluster by using the installer interface, it populates a [cluster configuration file](#config) for the management cluster with the required parameters. You can use the created configuration file as a model for future deployments from the CLI to this infrastructure provider.
- If you plan on registering the management cluster with Tanzu Mission Control, ensure that your Tanzu Kubernetes clusters meet the requirements listed in [Requirements for Registering a Tanzu Kubernetes Cluster with Tanzu Mission Control](https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-concepts/GUID-3AE5F733-7FA7-4B34-8935-C25D41D15EF9.html) in the Tanzu Mission Control documentation.
- Read the [Tanzu Kubernetes Grid 1.3.1 Release Notes](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3.1/rn/VMware-Tanzu-Kubernetes-Grid-131-Release-Notes.html) for updates related to security patches.

### vSphere Prerequisites

- Make sure that you have met the all of the requirements listed in [Prepare to Deploy Management Clusters to vSphere](vsphere.md).
- **NOTE**: On vSphere with Tanzu, you do not need to deploy a management cluster. See [Use the Tanzu CLI with a vSphere with Tanzu Supervisor Cluster](../tanzu-k8s-clusters/connect-vsphere7.md).

### Amazon EC2 Prerequisites

- Make sure that you have met the all of the requirements listed [Prepare to Deploy Management Clusters to Amazon EC2](aws.md).
- For information about the configurations of the different sizes of node instances, for example, `t3.large` or `t3.xlarge`, see [Amazon EC2 Instance Types](https://aws.amazon.com/ec2/instance-types/).
- For information about when to create a Virtual Private Cloud (VPC) and when to reuse an existing VPC, see [Resource Usage in Your Amazon Web Services Account](aws.md#aws-resources).
- If this is the first time that you are deploying a management cluster to Amazon EC2, create a Cloud Formation stack for Tanzu Kubernetes Grid in your AWS account by following the instructions in [Create IAM Resources](#create) below.

#### <a id="create"></a> Create IAM Resources

Before you deploy a management cluster to Amazon EC2 for the first time,
you must create a CloudFormation stack for Tanzu Kubernetes Grid, `tkg-cloud-vmware-com`, in your AWS account.
This CloudFormation stack includes the identity and access management (IAM) resources that Tanzu Kubernetes Grid needs to create and run clusters on Amazon EC2. For more information, see [Required IAM Resources](aws.md#iam-permissions) in _Prepare to Deploy Management Clusters to Amazon EC2_.

1. If you have already created the CloudFormation stack for Tanzu Kubernetes Grid in your AWS account, skip the rest of this procedure.

1. If you have not already created the CloudFormation stack for Tanzu Kubernetes Grid in your AWS account, ensure that AWS authentication variables are set either in the local environment or in your AWS default credential provider chain. For instructions, see [Configure AWS Account Credentials and SSH Key](aws.md#register-ssh).

   If you have configured AWS credentials in multiple places, the credential settings used to create the CloudFormation stack are applied in the following order of precedence:

   - Credentials set in the local environment variables `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN` and `AWS_REGION` are applied first.
   - Credentials stored in a shared credentials file as part of the default credential provider chain. You can specify the location of the credentials file to use in the local environment variable `AWS_SHARED_CREDENTIAL_FILE`. If this environment variable in not defined, the default location of `$HOME/.aws/credentials` is used. If you use credential profiles, the command uses the profile name specified in the `AWS_PROFILE` local environment configuration variable. If you do not specify a value for this variable, the profile named `default` is used.

    For an example of how the default AWS credential provider chain is interpreted for Java apps, see [Working with AWS Credentials](https://docs.aws.amazon.com/sdk-for-java/latest/developer-guide/credentials.html) in the AWS documentation.

1. Run the following command:

   ```
   tanzu management-cluster permissions aws set
   ```

  For more information about this command, run `tanzu management-cluster permissions aws set --help`.

**IMPORTANT:** The `tanzu management-cluster permissions aws set` command replaces the `clusterawsadm` command line utility that existed in Tanzu Kubernetes Grid v1.1.x and earlier. For existing management and Tanzu Kubernetes clusters initially deployed with v1.1.x or earlier, continue to use the CloudFormation stack that was created by running the `clusterawsadm alpha bootstrap create-stack` command. For Tanzu Kubernetes Grid v1.2 and later clusters, use the `tkg-cloud-vmware-com` stack.

### Microsoft Azure Prerequisites

- Make sure that you have met the requirements listed in [Prepare to Deploy Management Clusters to Microsoft Azure](azure.md).
- For information about the configurations of the different sizes of node instances for Azure, for example, `Standard_D2s_v3` or `Standard_D4s_v3`, see [Sizes for virtual machines in Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/sizes).

## <a id="config"></a> Create the Cluster Configuration File

Before creating a management cluster using the Tanzu CLI, you must define its configuration in a YAML configuration file that provides the base configuration for the cluster. When you deploy the management cluster from the CLI, you specify this file by using the `--file` option of the `tanzu management-cluster create` command.

Running `tanzu management-cluster create` command for the first time creates the `~/.tanzu/tkg` subdirectory that contains the Tanzu Kubernetes Grid configuration files.

If you have previously deployed a management cluster by running `tanzu management-cluster create --ui`, the `~/.tanzu/tkg/clusterconfigs` directory contains management cluster configuration files with settings saved from each invocation of the installer interface. Depending the infrastructure on which you deployed the management cluster, you can use these files as templates for  cluster configuration files for new deployments to the same infrastructure. Alternatively, you can create management cluster configuration files from the templates that are provided in this documentation.

- To use the configuration file from a previous deployment that you performed by using the installer interface, make a copy of the configuration file with a new name, open it in a text editor, and update the configuration. For information about how to update all of the settings, see the [Tanzu CLI Configuration File Variable Reference](../tanzu-config-reference.md).
- To create a new configuration file, see [Create a Management Cluster Configuration File](create-config-file.md). This section provides configuration file templates for each infrastructure provider.

VMware recommends using a dedicated configuration file for each management cluster, with configuration settings specific to a single infrastructure.

## <a id="set-variable"></a> (v1.3.1 Only) Set the <code>TKG_BOM_CUSTOM_IMAGE_TAG</code>

Before you can deploy a management cluster, you must specify the correct BOM file to use as a local environment variable. In the event of a patch release to Tanzu Kubernetes Grid, the BOM file may require an update to coincide with updated base image files.

**Note** For more information about recent security patch updates to VMware Tanzu Kubernetes Grid v1.3, see the [VMware Tanzu Kubernetes Grid v1.3.1 Release Notes](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3.1/rn/VMware-Tanzu-Kubernetes-Grid-131-Release-Notes.html) and this [Knowledgebase Article](https://kb.vmware.com/s/article/83781).

On the machine where you run the Tanzu CLI, perform the following steps:

1. Remove any existing BOM data.

   ```
   rm -rf ~/.tanzu/tkg/bom
   ```

1. Specify the updated BOM to use by setting the following variable.

   ```
   export TKG_BOM_CUSTOM_IMAGE_TAG="v1.3.1-patch1"
   ```

1. Run `tanzu management-cluster create` command with no additional parameters.

   ```
   tanzu management-cluster create
   ```

   This command produces an error but results in the BOM files being downloaded to `~/.tanzu/tkg/bom`.

## <a id="mc-create"></a> Run the `tanzu management-cluster create` Command

After you have created or updated the cluster configuration file and downloaded the most recent BOM, you can deploy a management cluster by running the `tanzu management-cluster create --file CONFIG-FILE` command, where `CONFIG-FILE` is the name of the configuration file. If your configuration file is the default `~/.tanzu/tkg/cluster-config.yaml`, you can omit the `--file` option.

<p class="note warning"><strong>Warning</strong>: The <code>tanzu management-cluster create</code> command takes time to complete.
While <code>tanzu management-cluster create</code> is running, do not run additional invocations of <code>tanzu management-cluster create</code> on the same bootstrap machine to deploy multiple management clusters, change context, or edit <code>~/.kube-tkg/config</code>.</p>

To deploy a management cluster, run the `tanzu management-cluster create` command. For example:

```
tanzu management-cluster create --file path/to/cluster-config-file.yaml
```

### Validation Checks

When you run `tanzu management-cluster create`, the command performs several validation checks before deploying the management cluster. The checks are different depending on the infrastructure to which you are deploying the management cluster.

- **vSphere**

    The command verifies that the target vSphere infrastructure meets the following requirements:

    * The vSphere credentials that you provided are valid.
    * Nodes meet the minimum size requirements.
    * Base image template exists in vSphere and is valid for the specified Kubernetes version.
    * Required resources including the resource pool, datastores, and folder exist in vSphere.

- **Amazon EC2**

    The command verifies that the target Amazon EC2 infrastructure meets the following requirements:

    * The AWS credentials that you provided are valid.
    * Cloud Formation stack exists.
    * Node Instance type is supported.
    * Region and AZ match.

- **Azure**

     The command verifies that the target Azure infrastructure meets the following requirements:

     * The Azure credentials that you provided are valid.
     * The public SSH key is encoded in base64 format.
     * The node instance type is supported.

If any of these conditions are not met, the `tanzu management-cluster create` command fails.

### Monitoring Progress

When you run `tanzu management-cluster create`, you can follow the progress of the deployment of the management cluster in the terminal. The first run of `tanzu management-cluster create` takes longer than subsequent runs because it has to pull the required Docker images into the image store on your bootstrap machine. Subsequent runs do not require this step, so are faster.

If `tanzu management-cluster create` fails before the management cluster deploys, you should clean up artifacts on your bootstrap machine before you re-run `tanzu management-cluster create`. See the [Troubleshooting Tips](../troubleshooting-tkg/tips.html) topic for details. If the machine on which you run `tanzu management-cluster create` shuts down or restarts before the local operations finish, the deployment will fail.

If the deployment succeeds, you see a confirmation message in the terminal:

```
Management cluster created! You can now create your first workload cluster by running tanzu cluster create [name] -f [file]
```

## <a id="what-next"></a> What to Do Next

- If you enabled identity management on the management cluster, you must perform post-deployment configuration steps to allow users to access the management cluster. For more information, see [Configure Identity Management After Management Cluster Deployment](configure-id-mgmt.md).
- For information about what happened during the deployment of the management cluster, how to connect `kubectl` to the management cluster, and how to create namespaces see [Examine the Management Cluster Deployment](verify-deployment.md).
- If you need to deploy more than one management cluster, on any or all of vSphere, Azure, and Amazon EC2, see [Manage Your Management Clusters](../cluster-lifecycle/multiple-management-clusters.md). This topic also provides information about how to add existing management clusters to your CLI instance, obtain credentials, scale and delete management clusters, add namespaces, and how to opt in or out of the CEIP.
