# TCE Headless Installation

Though the deployment of Tanzu Community Edition clusters can be done from any
machine, due to user preference or other constraints, performing the deployment
operations on a machine without a desktop environment may be required.

Our current (and recommended) installation method is to use the web based user
interface to perform new management cluster deployments. This is initiated with
the command:

```sh
tanzu management-cluster create --ui
```

This command will attempt to open a web browser pointed to the installation
wizard that steps through the various input needed to deploy a new management
cluster. If there is no desktop environment present, this command will give an
error message that it was unable to open a browser.

There are several ways to approach these "headless" installations. This
document provides a few suggestions for ways that are known to work well in the
past.

## Deployment Overview

Before going in to the various options for deploying in a headless environment,
it may help to have a basic understanding of what steps are necessary when
preparing a new management cluster deployment.

Ultimately, there are two main steps the user interface performs:

1. Collecting input in order to generate a cluster configuration YAML file
1. Passing that YAML file to the cluster create call to perform the deployment

In order to deploy from a headless machine, you will need to create this
configuration YAML file and call the `tanzu management-cluster create` command
using the file.

Here are a few different options to accomplish this task.

### Copy Generated File

The easist option may be to run the cluster creation UI on a machine that does
have a desktop environment. You can then copy the generated YAML file from this
machine over to the headless machine to perform that actual deployment.

On the machine with a desktop environment, run the command line:

```sh
tanzu management-cluster create --ui
```

As in normal operations, that will launch a browser-based UI that will prompt
for the various configuration settings needed. Follow the wizard all the way
through to final "Confirm Settings" page if the wizard.

On this page you can review your settings and go back to make any additional
changes. If you scroll all the way down to the bottom of this page, there is a
box labeled "CLI Command Equivalent" that shows the command line that can be
used to perform the second step of kicking off the deployment.

![cli command equivalent](images/wizard-cli-command.png)

Note the file path of the generated YAML file (the path is listed after the
`--file` argument). In this case, the file path is:

```sh
/home/smcginnis/.tanzu/tkg/clusterconfigs/bs83endsfl.yaml
```

Once you have the file path, you can copy that file over to your headless
machine. This can be done via a file share, thumb drive, or whatever means you
have at your disposal. Here is an example of using the SCP command to copy it
into the home directory on the headless machine:

```sh
scp /home/smcginnis/.tanzu/tkg/clusterconfigs/bs83endsfl.yaml ubuntu@deployment-vm:~/
```

Now the file is available on the headless machine to actually perform the
deployment. You can now run the command given at the end of the deployment UI
from the terminal on this machine, pointing to the file in your home directory:

```sh
tanzu management-cluster create --file ~/bs83endsfl.yaml -v 6
```

### Remote UI

Another option is to use the graphical UI to perform the installation, but
launched from the headless machine and run on your local machine with a desktop
environment.

In order to operate the UI in this scenario, we need to forward the local web
server port that the browser would normally connect to through to our desktop
machine using SSH port tunneling. To do so, first connect to your machine using
the `-L` argument for local port forwarding with a command similar to:

```sh
ssh -L 127.0.0.1:8080:127.0.0.1:8080 ubuntu@deployment-vm
```

This tells SSH to forward your localhost port 8080 to localhost:8080 on the
headless machine. With the port forwarding enabled, run this variation on the
management cluster creation command to launch the UI web server:

```sh
tanzu management-cluster create --ui --browser=none
```

The `--browser=none` argument tells the command to serve up the web interface,
but do not attempt to launch a browser instance to connect to it.

On your local desktop machine, you can now use your web browser to connect to
[http://localhost:8080](http://localhost:8080) and be forwarded through to the
tanzu command running on the headless machine.

In this scenario you can use the UI to complete the deployment completely, or
if you prefer, you can continue through to the "Confirm Settings" page and use
the provided CLI command at the end to run it in the terminal on the headless
machine.

To exit the kickstart UI once complete, press `Ctrl-C` in the terminal to exit
the command.

### Manual Creation

It is also an option to manually create the YAML file. In this case you can
just use a text editor to create a file such as `~/mycluster.yaml` and then run
the command:

```sh
tanzu management-cluster create --file ~/mycluster.yaml -v 6
```

The settings needed in the YAML file vary depending on the infrastructure being
deployed to. Refer to the [Tanzu CLI Configuration File Variable
Reference](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.3/vmware-tanzu-kubernetes-grid-13/GUID-tanzu-config-reference.html)
for a complete list of options.

Here are some minimal templates to help get started.

#### vSphere Template

The following can be used to create a YAML file. Pay particular attention to
the `VSPHERE_*` settings to match your environment and desired configuration.

```yaml
AVI_ENABLE: "false"
CLUSTER_CIDR: 100.96.0.0/11
CLUSTER_PLAN: dev
ENABLE_CEIP_PARTICIPATION: "false"
ENABLE_MHC: "true"
IDENTITY_MANAGEMENT_TYPE: none
INFRASTRUCTURE_PROVIDER: vsphere
SERVICE_CIDR: 100.64.0.0/13
VSPHERE_CONTROL_PLANE_DISK_GIB: "40"
VSPHERE_CONTROL_PLANE_ENDPOINT: 192.168.1.210
VSPHERE_CONTROL_PLANE_MEM_MIB: "8192"
VSPHERE_CONTROL_PLANE_NUM_CPUS: "2"
VSPHERE_DATACENTER: /Datacenter1
VSPHERE_DATASTORE: /Datacenter1/datastore/VMData
VSPHERE_FOLDER: /Datacenter1/vm
VSPHERE_NETWORK: VM Network
VSPHERE_PASSWORD:
VSPHERE_RESOURCE_POOL: /Datacenter1/host/192.168.1.201/Resources
VSPHERE_SERVER: 192.168.1.202
VSPHERE_SSH_AUTHORIZED_KEY:
VSPHERE_TLS_THUMBPRINT:
VSPHERE_USERNAME: administrator@lab.local
VSPHERE_WORKER_DISK_GIB: "40"
VSPHERE_WORKER_MEM_MIB: "8192"
VSPHERE_WORKER_NUM_CPUS: "2"
```

#### AWS Template

These settings may be used for an AWS deployment. The trickiest part with AWS
is the need to populate the `AWS_B64ENCODED_CREDENTIALS` value. This is a
base64 encoded value of your AWS account settings.

To get the value to supply to `AWS_B64ENCODED_CREDENTIALS`, follow these steps:

1. Create a temporary file to store your credential information containing the
   following settings:

   ```ini
   [default]
   aws_access_key_id = [AWS_ACCESS_KEY]
   aws_secret_access_key = [AWS_SECRET_ACCESS_KEY]
   region = us-east-1
   ```

1. Get the base64 encryption of this file:

   ```sh
   base64 credentials.txt
   ```

1. Copy the output from that command and use it as the value for the
   `AWS_B64ENCODED_CREDENTIALS` setting.

Make sure to update the `NODE_MACHINE_TYPE` and other values to match your
desired configuration.

```sh
AWS_ACCESS_KEY_ID:
AWS_AMI_ID: ami-0bcd9ed3ef40fad77
AWS_B64ENCODED_CREDENTIALS:
AWS_NODE_AZ: us-east-2b
AWS_PRIVATE_NODE_CIDR: 10.0.16.0/20
AWS_PUBLIC_NODE_CIDR: 10.0.0.0/20
AWS_REGION: us-east-2
AWS_SECRET_ACCESS_KEY:
AWS_SSH_KEY_NAME:
AWS_VPC_CIDR: 10.0.0.0/16
BASTION_HOST_ENABLED: "true"
CLUSTER_CIDR: 100.97.0.0/11
CLUSTER_PLAN: dev
CONTROL_PLANE_MACHINE_TYPE: t3a.large
ENABLE_CEIP_PARTICIPATION: "false"
ENABLE_MHC: "true"
IDENTITY_MANAGEMENT_TYPE: none
INFRASTRUCTURE_PROVIDER: aws
NODE_MACHINE_TYPE: t3a.large
SERVICE_CIDR: 100.64.0.0/13
```
