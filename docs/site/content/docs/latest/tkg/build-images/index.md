# Building Machine Images

You can build custom machine images for Tanzu Kubernetes Grid to use as a VM template for the management and Tanzu Kubernetes (workload) cluster nodes that it creates.
Each custom machine image packages a base OS version and a Kubernetes version, along with any additional customizations, into an image that runs on vSphere, Amazon EC2, or Microsoft Azure infrastructure.
The base OS can be an OS that VMware supports but does not distribute, such as Red Hat Enterprise Linux (RHEL) v7.

This topic provides background on custom images for Tanzu Kubernetes Grid, and explains how to build them.

## <a id="overview"></a> Overview: Kubernetes Image Builder

To build custom machine images for Tanzu Kubernetes Grid cluster nodes, you use the container image from the upstream [Kubernetes Image Builder](https://github.com/kubernetes-sigs/image-builder) project.
Kubernetes Image Builder runs on your local workstation and uses the following:

* [Packer](https://www.packer.io/downloads) automates and standardizes the image-building process for current and future CAPI providers, and packages the images for their target infrastructure once they are built.
* [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html) standardizes the process of configuring and provisioning machines across multiple target distribution families, such as Ubuntu and CentOS.
* To build the images, Image Builder uses native infrastructure for each provider:
  - **Amazon EC2**
      - You build your custom images from base AMIs that are published on Amazon EC2, such as official Ubuntu AMIs.
      - The custom image is built inside AWS and then stored in your AWS account in one or more regions.
      - See [Building Images for AWS](https://image-builder.sigs.k8s.io/capi/providers/aws.html) in the Image Builder documentation.
  - **Azure**:
      - You can store your custom image in an Azure Shared Image Gallery.
      - See [Building Images for Azure](https://image-builder.sigs.k8s.io/capi/providers/azure.html) in the Image Builder documentation.
  - **vSphere**:
      - Image Builder builds Open Virtualization Archive (OVA) images
      - You build the machine images from the Linux distribution's original installation `ISO`.
      - You import the resulting OVA into a vSphere cluster, take a snapshot for fast cloning, and then mark the machine image as a `vm template`.
      - See [Building Images for vSphere](https://image-builder.sigs.k8s.io/capi/providers/vsphere.html) in the Image Builder documentation.

See [Customization](https://image-builder.sigs.k8s.io/capi/capi.html#customization) in the Image Builder documentation for how you can customize your image.
Before making any modifications, consult with VMware Customer Reliability Engineering (CRE) for best practices and recommendations.

After you have created a custom image, you enable the `tanzu` CLI to use it by creating a custom Tanzu Kubernetes release (TKr) based on the image.

### <a id="replace"></a> Custom Images Replace Default Images

For common combinations of OS version, Kubernetes version, and target infrastructure, Tanzu Kubernetes Grid provides default machine images.
For example, one `ova-ubuntu-2004-v1.20.5+vmware.2-tkg` image serves as the OVA image for Ubuntu v20.04 and Kubernetes v1.20.5 on vSphere.

For other combinations of OS version, Kubernetes version, and infrastructure, such as with the RHEL v7 OS, there are no default machine images, but you can build them.

If you build and use a custom image with the same OS version, Kubernetes version, and infrastructure that a default image already has, your custom image replaces the default.
The `tanzu` CLI then creates new clusters using your custom image, and no longer uses the default image,
for that combination of OS version, Kubernetes version, and target infrastructure.

### <a id="cluster-api"></a> Cluster API

[Cluster API (CAPI)](https://github.com/kubernetes-sigs/cluster-api) is built on the principles of immutable infrastructure. All nodes
that make up a cluster are derived from a common template or machine image.

When CAPI creates a cluster from a machine image, it expects several
things to be configured, installed, and accessible or running, including:

* The versions of `kubeadm`, `kubelet` and `kubectl` specified in the cluster manifest.
* A container runtime, most often `containerd`.
* All required images for `kubeadm init` and `kubeadm join`. You must include any images that are not published and must be pulled locally, as with VMware-signed images.
* `cloud-init` configured to accept bootstrap instructions.

## <a id="build"></a> Build a Custom Machine Image

This procedure builds a custom machine image for Tanzu Kubernetes Grid to use when creating management or workload cluster nodes.
It works by:

1. Collecting parameter strings that give a Kubernetes Image Builder command the context and inputs that it needs to create the custom image.
1. Passing the parameter strings to a long `docker run` command that runs the Kubernetes Image Builder command in a container.

### <a id="prerequisites"></a> Prerequisites

To create a custom machine image, you need:

* An account on your target infrastructure, AWS, Azure, or vSphere
* A macOS or Linux workstation with the following installed:
  * [Docker Desktop](https://www.docker.com/products/docker-desktop)
  * For **AWS**: The `aws` command-line interface (CLI)
  * For **Azure**: The `az` CLI
  * For **vSphere**: A local copy of the `OVFTool` linux [installer](https://code.vmware.com/web/tool/4.4.0/ovf)

### <a id="procedure"></a> Procedure

1. On **AWS** and **Azure**, log in to your infrastructure CLI. Authenticate and specify your region, if prompted:
  * **AWS**: Run `aws configure`.
  * **Azure**: Run `az login`.

1. On **vSphere**, do the following:

    1. Download Open Virtualization Format (OVF) tool from [VMware {code}](https://code.vmware.com/web/tool/4.4.0/ovf). You will need the installer for x86_64 Linux, as you will not be installing this locally, rather installing into a Linux container. In the following step, this file is referred to as `YOUR-OVFTOOL-INSTALLER-FILE`, and should be in the same directory as your new `Dockerfile`

    1. Create a `Dockerfile` and fill in values as shown:

      ```
      FROM k8s.gcr.io/scl-image-builder/cluster-node-image-builder-amd64:v0.1.9
      USER root
      ENV LC_CTYPE=POSIX
      ENV OVFTOOL_FILENAME=YOUR-OVFTOOL-INSTALLER-FILE
      ADD $OVFTOOL_FILENAME /tmp/
      RUN /bin/sh /tmp/$OVFTOOL_FILENAME --console --required --eulas-agreed && \
          rm -f /tmp/$OVFTOOL_FILENAME
      USER imagebuilder
      ENV IB_OVFTOOL=1
      ```

    1. Build a new container image from the `Dockerfile`. It is recommended to give a custom name that will be meaningful to you:

      ```bash
      docker build . -t projects.registry.vmware.com/tkg/imagebuilder-byoi:v0.1.9
      ```

    1. Create a vSphere credentials JSON file and fill in its values:

      ```
      {
      "cluster": "",
      "convert_to_template": "false",
      "create_snapshot": "true",
      "datacenter": "",
      "datastore": "",
      "folder": "",
      "insecure_connection": "false",
      "linked_clone": "true",
      "network": "",
      "password": "",
      "resource_pool": "",
      "template": "",
      "username": "",
      "vcenter_server": ""
      }
      ```

1. Determine the Image Builder configuration version that you want to build from.

   - Search the VMware {code} [Sample Exchange](https://code.vmware.com/samples) for `TKG Image Builder` to list the available versions.
   - Each version corresponds to the Kubernetes version that Image Builder uses. For example, `TKG-Image-Builder-for-Kubernetes-v1.20.5-master.zip` builds a Kubernetes v1.20.5 image.
   - If you need to create a management cluster, which you must do when you first install Tanzu Kubernetes Grid, choose the default Kubernetes version of your Tanzu Kubernetes Grid version. For example, in Tanzu Kubernetes Grid v1.3.1, the default Kubernetes version is v1.20.5.

1. The Image Builder configurations have two different architectures and build instructions, based on their Kubernetes versions:
  - For v1.20.5, v1.20.4, v1.19.9, v1.19.8, v1.18.17, v1.18,16, or v1.17.16, continue with the procedure below.
  - For v1.19.3, v1.19.1, v1.18.10, v1.18.8, v1.17.13, and v1.17.77, Follow the **Build an Image with Kubernetes Image Builder** instructions in the Tanzu Kubernetes Grid v1.2 documentation:
      - [Build and Use Custom AMI images on Amazon EC2](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-build-images-capa-images.html)
      - [Build and Use Custom VM images on Azure](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-build-images-capz-images.html)
      - [Build and Use Custom OVA Images on vSphere](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-build-images-capv-images.html)

  After creating a custom image file following the v1.2 procedure, continue with [Use a Custom Machine Image](#use) below.
  Do not follow the Tanzu Kubernetes Grid v1.2 procedure to add a reference to the custom image to a Bill of Materials (BoM) file.

1. Download the configuration code zip file, and unpack its contents.

1. `cd` into the `TKG-Image-Builder-` directory, so that the `tkg.json` file is in your current directory.

1. Collect the following parameter strings to plug into the command in the next step. Many of these specify `docker run -v` parameters that copy your current working directories into the `/home/imagebuilder` directory of the container used to build the image.

  - `AUTHENTICATION`: Copies your local CLI directory:
      - **AWS**: Use `~/.aws:/home/imagebuilder/.aws`
      - **Azure**: Use `~/.azure:/home/imagebuilder/.azure`
      - **vSphere**: `/PATH/TO/CREDENTIALS.json:/home/imagebuilder/vsphere.json`
  - `SOURCES`: Copies the repo's `tkg.json` file, which lists download sources for versioned OS, Kubernetes, container network interface (CNI). images:
      - Use `/PATH/TO/tkg.json:/home/imagebuilder/tkg.json`
  - `ROLES`: Copies the repo's `tkg` directory, which contains Ansible roles required by Image Builder.
      - Use `/PATH/TO/tkg:/home/imagebuilder/tkg`
      - To add custom Ansible roles, edit the `tkg.json` file to reformat the `custom_role_names` setting with escaped quotes (`\"`), in order to make it a list with multiple roles. For example: <br />
      `  "custom_role_names": "\"/home/imagebuilder/tkg /home/imagebuilder/mycustomrole\"",`
  - `TESTS`: Copies a `goss` test directory designed for the image's target infrastructure, OS, and Kubernetes verson:
      - Use the filename of a file in the repo's `goss` directory, for example `amazon-ubuntu-1.20.5+vmware.2-goss-spec.yaml`.
  - `CUSTOMIZATIONS`: Copies a customizations file in JSON format.
  - (Azure) `AZURE-CREDS`: Path to an Azure credentials file, as described in the [Image Builder documentation](https://image-builder.sigs.k8s.io/capi/container-image.html#examples).
  - `CONTAINER`: A container hosted on Google Cloud Platform:
      - **AWS**: Use `k8s.gcr.io/scl-image-builder/cluster-node-image-builder-amd64:v0.1.9`
      - **Azure**: Use `k8s.gcr.io/scl-image-builder/cluster-node-image-builder-amd64:v0.1.9`
      - **vSphere**: Use `k8s.gcr.io/scl-image-builder/cluster-node-image-builder-amd64:v0.1.9`
  - `COMMAND`: Use a command like one of the following, based on the custom image OS.  For vSphere and Azure images, the commands start with `build-node-ova-` and `build-azure-sig-`:
      - `build-ami-ubuntu-2004`: Ubuntu v20.04
      - `build-ami-ubuntu-1804`: Ubuntu v18.04
      - `build-ami-amazon-2`: Amazon Linux 2

1. Using the strings above, run the Image Builder in a Docker container:

    ```sh
    docker run -it --rm \
        -v AUTHENTICATION \
        -v SOURCES \
        -v ROLES \
        -v /PATH/TO/goss/TESTS.yaml:/home/imagebuilder/goss/goss.yaml \
        -v /PATH/TO/CUSTOMIZATIONS.json:/home/imagebuilder/CUSTOMIZATIONS.json \
        --env PACKER_VAR_FILES="tkg.json CUSTOMIZATIONS.json" \
        --env-file AZURE-CREDS \
        CONTAINER \
        COMMAND \
    ```

    Notes:

    * Omit `env-file` if you are not building an image for Azure.
    * This command may take several minutes to complete.

    For example, to create a custom image with Ubuntu v20.04 and Kubernetes v1.20.5 to run on AWS, running from the directory that contains `tkg.json`:

    ```sh
    docker run -it --rm \
        -v ~/.aws:/home/imagebuilder/.aws \
        -v $(pwd)/tkg.json:/home/imagebuilder/tkg.json \
        -v $(pwd)/tkg:/home/imagebuilder/tkg \
        -v $(pwd)/goss/amazon-ubuntu-1.20.5+vmware.2-goss-spec.yaml:/home/imagebuilder/goss/goss.yaml \
        -v /PATH/TO/CUSTOMIZATIONS.json /home/imagebuilder/aws.json \
        --env PACKER_VAR_FILES="tkg.json aws.json" \
        k8s.gcr.io/scl-image-builder/cluster-node-image-builder-amd64:v0.1.9 \
        build-ami-ubuntu-2004
    ```

    For vSphere, you must use the custom container image created above. You must also set a version string that will match what you pass in your custom TKr in the later steps. While VMware published OVAs will have a version string like `v1.20.5+vmware.2-tkg.1`, it is recommended that the `-tkg.1` be replaced with a string meaningful to your organization. To set this version string, define it in a `metadata.json` file like the following:

    ```json
    {
      "VERSION": "v1.20.5+vmware.2-myorg.0"
    }
    ```

    When building OVAs, the `.ova` file is saved to the local filesystem of your workstation. Whatever folder you want those OVAs to be saved in should be mounted to `/home/imagebuilder/output` within the container.
    Then, create the OVA using the container image:

    ```sh
    docker run -it --rm \
      -v /PATH/TO/CREDENTIALS.json:/home/imagebuilder/vsphere.json \
      -v $(pwd)/tkg.json:/home/imagebuilder/tkg.json \
      -v $(pwd)/tkg:/home/imagebuilder/tkg \
      -v $(pwd)/goss/vsphere-ubuntu-1.20.5+vmware.2-goss-spec.yaml:/home/imagebuilder/goss/goss.yaml \
      -v $(pwd)/metadata.json:/home/imagebuilder/metadata.json \
      -v /PATH/TO/OVA/DIR:/home/imagebuilder/output \
      --env PACKER_VAR_FILES="tkg.json vsphere.json" \
        --env OVF_CUSTOM_PROPERTIES=/home/imagebuilder/metadata.json \
      projects.registry.vmware.com/tkg/imagebuilder-byoi:v0.1.9 \
      build-node-ova-vsphere-ubuntu-2004
    ```

## <a id="use"></a> Use a Custom Machine Image

After you have created a custom image, you enable the `tanzu` CLI to use it by creating a custom Tanzu Kubernetes release (TKr) based on the image.

To create a custom TKr, you add it to the Bill of Materials (BoM) of the TKr for the image's Kubernetes version.
For example, to add a custom image that you built with Kubernetes v1.20.5, you modify the current `~/.tanzu/tkg/bom/tkr-bom-v1.20.5.yaml` file.

1. From your `~/.tanzu/tkg/bom/` directory, open the TKr BoM corresponding to your custom image's Kubernetes version. For example with a filename like `tkr-bom-v1.20.5+vmware.2-tkg.1.yaml` for Kubernetes v1.20.5.
  - If the directory lacks the TKr BoM file that you need, you can bring it in by deploying a cluster with the desired Kubernetes version, as described in [Deploy a Cluster with a Non-Default Kubernetes Version](../tanzu-k8s-clusters/k8s-versions.md#non-default).

1. In the BoM file, find the image definition blocks for your infrastructure: `ova` for vSphere, `ami` for AWS, and `azure` for Azure.

1. Determine whether an existing definition block applies to your image's OS, as listed by `osinfo.name`, `.version`, and `.arch`.

1. If no existing block applies to your image's `osinfo`, add a new block as follows.
If an existing block does apply, replace its values as follows:

  - **vSphere**:
      - `name:` a unique name for your OVA that includes the OS version, like `my-ubuntu-2004`
      - `version:` follow existing `version` value format, but use the unique `VERSION` assigned in `metadata.json` when you created the OVA, for example `v1.20.5+vmware.2-myorg.0`.
  - **AWS** - for each region that you plan to use the custom image in:
      - `id:` follow existing `id` value format, but use a unique hex string at the end, for example `ami-693a5e2348b25e428`
  - **Azure**:
      - `sku:` a unique SKU for your image that includes the OS version, like `my-k8s-1dot20dot4-ubuntu-2004`

    If the BoM file defines images under regions, your new or modified custom image definition block must be listed first in its region.
    Within each region, the cluster creation process picks the first suitable image listed.

1. Save the BoM file.  If its filename includes a plus (`+`) character, save the modified file under a new filename that replaces the `+` with a triple dash (`---`).  For example, `tkr-bom-v1.20.5---vmware.2-tkg.1.yaml`.

1. `base64`-encode the file contents into a binary string, for example:

    ```
    cat tkr-bom-v1.20.5---vmware.2-tkg.1.yaml | base64 -w 0
    ```

1. Create a `ConfigMap` YAML file in the `tkr-system` namespace, also without a `+` in its filename, and fill in values as shown:

    ```
    apiVersion: v1
    kind: ConfigMap
    metadata:
     name: CUSTOM-TKG-BOM
     labels:
       tanzuKubernetesRelease: CUSTOM-TKR
    binaryData:
     bomContent: BOM-BINARY-CONTENT
    ```

    Where:

    - `CUSTOM-TKG-BOM` is the name of the `ConfigMap` YAML file, without the `.yaml` extension, such as `my-custom-tkr-bom-v1.20.5---vmware.2-tkg.1`
    - `CUSTOM-TKR` is a name for your TKr, such as `my-custom-tkr-v1.20.5---vmware.2-tkg.1`

    - `BOM-BINARY-CONTENT` is the `base64`-encoded content of your customized BoM file.

1. Save the `ConfigMap` file, set the `kubectl` context to a management cluster you want to add TKr to, and apply the file to the cluster, for example:

    ```
    kubectl apply -f my-custom-tkr-bom-v1.20.5---vmware.2-tkg.1.yaml

    ```

    * Once the `ConfigMap` is created, the TKr Controller reconciles the new object by creating a `TanzuKubernetesRelease`.<br />
    The default reconciliation period is??600 seconds.
    You can avoid this delay by deleting the TKr Controller pod, which makes the pod restore and reconcile immediately:
      1. List pods in the `tkr-system` namespace:

          ```
          kubectl get pod -n tkr-system
          ```

      1. Retrieve the name of the TKr Controller pod, which looks like `tkr-controller-manager-f7bbb4bd4-d5lfd`
      1. Delete the pod:

          ```
          kubectl delete pod -n tkr-system TKG-CONTROLLER
          ```

        Where `TKG-CONTROLLER` is the name of the TKr Controller pod.

1. To check that the custom TKr was added, run `tanzu kubernetes-release get` or `kubectl get tkr` or and look for the `CUSTOM-TKR` value set above in the output.

Once your custom TKr is listed by the `kubectl` and `tanzu` CLIs, you can use it to create management or workload clusters as described below.

### <a id="use-mc"></a> Use a Custom Machine Image in a Management Cluster

To create a management cluster that uses your custom image as the base OS for its nodes:

1. Upload the image to your cloud provider.

1. When you run the installer interface, select the custom image in the **OS Image** pane, as described in [Select the Base OS Image(../mgmt-clusters/deploy-ui.md#base-os).

### <a id="use-mc"></a> Use a Custom Machine Image in a Workload Cluster

To create a workload cluster that uses your custom image as the base OS for its nodes,
pass its TKr name as listed by `tanzu kubernetes-release get` to the `--tkr` option of `tanzu cluster create`.
