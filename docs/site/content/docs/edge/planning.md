# Planning your Deployment

A Tanzu Community Edition deployment consists of the following high level steps:

1. **Decision time:** You must decide, which of the following suits your desired outcome:

    Do you want a single node, local workstation cluster suitable for a development/test environment? If so deploy an [unmanaged cluster](planning/#unmanaged-cluster). Start here: [Deploy Unmanaged Clusters](getting-started-unmanaged).

    **or**

    Do you want a full-featured, scalable Kubernetes implementation suitable for a development or production environment? If so deploy a [managed cluster](planning/#managed-cluster). Start here: [Deploy Managed Clusters](getting-started).

1. **Install Tanzu CLI:** Regardless of which route you choose, you will start by installing the Tanzu CLI. Installation steps are available in both the [managed](getting-started/#install-tanzu-cli) and [unmanaged](getting-started-unmanaged/#install-tanzu-cli) cluster Getting Started Guides.

1. **Deploy your chosen cluster type:** Now that you have chosen your cluster deployment type and installed the Tanzu CLI, you are ready to deploy your cluster. Deployment steps are available in both the [managed](getting-started/#deploy-clusters) and [unmanaged](getting-started-unmanaged/#deploy-a-cluster) cluster Getting Started Guides.

1. **Deploy packages:** If you want to deploy applications on any cluster you will need some [packages](planning/#packages).  
       - For unmanaged clusters, the [package repository](planning/#package-repository) is already installed, so you can start to install individual packages. Steps for a sample package installation are available in [Create Unmanaged Clusters](getting-started-unmanaged/#deploy-a-package).  
       - For a managed cluster deployment, you must first install the package repository, and then you are ready to install individual packages. Steps are available for both in [Create Managed Clusters](getting-started/#deploy-a-package).

## Component Descriptions

### Managed Cluster

{{% include "/docs/assets/mgmt-desc.md" %}}

### Unmanaged Cluster

{{% include "/docs/assets/unmanaged-desc.md" %}}

### Packages

{{% include "/docs/assets/package-description.md" %}}

### Package Repository

{{% include "/docs/assets/package-repository.md" %}}
