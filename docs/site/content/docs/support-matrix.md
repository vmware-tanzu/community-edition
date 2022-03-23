# Support Matrix

The following topic provides details of the supported operating systems, infrastructure providers, and Kubernetes versions.

{{% include "/docs/assets/unmanaged-cluster-note.md" %}}

## Local Client/Bootstrap Machine Prerequisites

Before you install the Tanzu CLI, **one** of the following operating system and
hardware/software configurations is required on your local machine.

{{% include "/docs/assets/prereq-linux.md" %}}

{{% include "/docs/assets/prereq-mac.md" %}}

{{% include "/docs/assets/prereq-windows.md" %}}

## Infrastructure Providers (Target Platforms)

After you install the Tanzu CLI Edition on your local machine, you can use it to deploy a cluster to **one** of the following infrastructure providers:
| Cloud Infrastructure Provider    |Local Infrastructure Provider|
|:------------------------ |:------------------------ |
|[Amazon Web Services (AWS)](https://github.com/kubernetes-sigs/cluster-api-provider-aws) |[Docker](https://github.com/kubernetes-sigs/cluster-api/tree/main/test/infrastructure/docker) |
|[Microsoft Azure](https://github.com/kubernetes-sigs/cluster-api-provider-azure)| |
|[vSphere](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/)| |

The Cloud Infrastructure providers support the following infrastructure platforms and operating systems (OSs):

| |**vSphere** | **AWS** | **Azure** |
|:------------------------ |:------------------------ |:------------------------ |:------------------------
|**Infrastructure Platform**|vSphere 6.7U3 and later, vSphere 7| Native AWS |Native Azure  |
|**Kubernetes node OS**|Photon OS 3, Ubuntu 20.04|Amazon Linux 2, Ubuntu 20.04 |Ubuntu 18.04, Ubuntu 20.04 |

If you are using Docker as your target infrastructure provider, the following additional configuration is needed:

|**Local Docker**|
|:------------------------|
|6 GB of RAM and 4 CPUs (with no other containers running).|
|15 GB of local machine disk storage for images |
|Bootstrapping a cluster to Docker from a Windows bootstrap machine is currently experimental.|

Note: Check your Docker configuration as follows:  

Linux: Run `docker system info`  
Mac and Windows: In Docker Desktop, select Preferences > Resources > Advanced

## Support Matrix Summary for Operating System and Infrastructure Provider

{{% include "/docs/assets/support-matrix.md" %}}

## Supported Kubernetes Versions

Tanzu Community Edition supports the following Kubernetes versions: `1.21.2, 1.20.8, 1.19.12`
