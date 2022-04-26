# Getting Started with Unmanaged Clusters

<!-- markdownlint-disable MD036 -->
<!-- markdownlint-disable MD024 -->

This guide walks you through creating a Tanzu workload from source code in a git repository using the Tanzu CLI and with unmanaged clusters.

*Note:* This guide is an introduction to using the Application Toolkit package. For a complete usage guide, refer to the [Application Toolkit readme](package-readme-app-toolkit-0.2.0)

## Before You Begin

1. Ensure you have a Dockerhub account

## Install Tanzu CLI

The Tanzu CLI is used for interacting with Tanzu Community Edition including creating an application and creating a cluster.

Choose your operating system below for guidance on installation.

{{< tabs tabTotal="3" tabID="1" tabName1="Linux" tabName2="Mac" tabName3="Windows">}}
{{< tab tabNum="1" >}}

### Linux System Requirements

{{% include "/docs/assets/prereq-unmanaged-linux.md" %}}

### Package Manager

**Homebrew**

{{% include "/docs/assets/install-homebrew.md" %}}

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< tab tabNum="2" >}}

### Mac System Requirements

{{% include "/docs/assets/prereq-unmanaged-mac.md" %}}

### Package Manager

**Homebrew**

{{% include "/docs/assets/install-homebrew.md" %}}

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< tab tabNum="3" >}}

### Windows System Requirements

{{% include "/docs/assets/prereq-unmanaged-windows.md" %}}

### Package Manager

**Chocolatey**

1. Install using [chocolatey](https://chocolatey.org/install), in **Powershell, as an administrator**.

    ```sh
    choco install tanzu-community-edition
    ```

### Direct Download

{{% include "/docs/assets/direct-download.md" %}}

{{< /tab >}}
{{< /tabs >}}

## Create an Unmanaged Cluster

1. Create a cluster with 80 and 443 ports exposed.

    ```sh
    tanzu unmanaged-cluster create app-demo -p 80:80 -p 443:443
    ```

    > **Note**: Ensure you expose the ports as they will enable your applications to be accessible!

2. Wait for the cluster to initialize. You will see this at the end of the successful execution of the command.

    ```txt
    ✅ Cluster created

    🎮 kubectl context set to beepboop

    View available packages:
       tanzu package available list
    View running pods:
       kubectl get po -A
    Delete this cluster:
       tanzu unmanaged-cluster delete beepboop
    ```

    > A container image larger than 1GB is used for running clusters. This may
    > cause your first cluster to take significantly longer to bootstrap than
    > subsequent clusters.

## Set-up Credentials

We will be creating a registry secret that will be provided to the Application Toolkit package to push and pull your application images from your registry.

1. Install `secretgen-controller` package for exporting your registry secret across all namespaces

  ```shell
  tanzu package install secretgen-controller --package-name secretgen-controller.community.tanzu.vmware.com --version 0.7.1 -n tkg-system
  ```

 1. Create a registry secret `registry-credentials` using your Docker credentials

  ```shell
  tanzu secret registry add registry-credentials --server https://index.docker.io/v1/ --username REGISTRY_USERNAME --password REGISTRY_PASSWORD --export-to-all-namespaces
  ```

  where

* `REGISTRY_USER`: the username for the Dockerhub account with write access
* `REGISTRY_PASS`: the password for the same account. You can also use the `--pasword-env-var`, `--password-file`, or `--password-stdin` options to provide your password in case you have special characters in your password.

**Note**: For the purposes of this Getting Started Guide, please ensure the name of the registry secret remains `registry-credentials`.

## Install App Toolkit Package

1. Verify if the App Toolkit package is available to install.

    ```sh
    tanzu package available get app-toolkit.community.tanzu.vmware.com
    ```

    ```txt
    NAME:                 app-toolkit.community.tanzu.vmware.com
    DISPLAY-NAME:         App-Toolkit package for TCE
    SHORT-DESCRIPTION:    Kubernetes-native toolkit to support application lifecycle
    PACKAGE-PROVIDER:     VMware
    LONG-DESCRIPTION:     Meta-Package to install a set of TCE packages that help users to create, iterate and manage their applications
    MAINTAINERS:          [{Glenio Borges} {Kris Applegate} {Ryan Collins}]
    SUPPORT:              For help go to #tanzu-community-edition in the kubernetes slack workspace.
    CATEGORY:             [application lifecycle]
    ```

1. Create a `app-toolkit-values.yaml`, copy the below content file, and update the file with your Dockerhub account credentials.

```yaml
contour:
  envoy:
    service:
      type: ClusterIP
    hostPorts:
      enable: true

knative_serving:
  domain:
    type: real
    name: 127-0-0-1.sslip.io

kpack:
  kp_default_repository: index.docker.io/[YOUR_DOCKERHUB_USERNAME]/app-toolkit-install
  kp_default_repository_username: [YOUR_DOCKERHUB_USERNAME]
  kp_default_repository_password: [YOUR_DOCKERHUB_PASSWORD]

cartographer_catalog:
  registry:
    server: index.docker.io
    repository: [YOUR_DOCKERHUB_USERNAME]
```

1. Install the Application Toolkit package.

Application Toolkit package will install software required to create a running application from your source code.

Application Toolkit will also prepare the `default` namespace so that you can start creating applications immediately. For preparing another namespace, please refer to the [Prepare your Developer Namespace](package-readme-app-toolkit-0.2.0#set-up-the-developer-namespace) section.

```shell
tanzu package install app-toolkit --package-name app-toolkit.community.tanzu.vmware.com --version 0.2.0 -f app-toolkit-values.yaml -n tanzu-package-repo-global
```

1. Verify the package is now installed.

  ```sh
  tanzu package installed list
  ```

  ```txt
  NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
  app-toolkit               app-toolkit.community.tanzu.vmware.com               0.2.0            Reconcile succeeded  tanzu-package-repo-global
  cartographer              cartographer.community.tanzu.vmware.com              0.3.0            Reconcile succeeded  tanzu-package-repo-global
  cartographer-catalog      cartographer-catalog.community.tanzu.vmware.com      0.3.0            Reconcile succeeded  tanzu-package-repo-global
  cert-manager              cert-manager.community.tanzu.vmware.com              1.6.1            Reconcile succeeded  tanzu-package-repo-global
  contour                   contour.community.tanzu.vmware.com                   1.20.1           Reconcile succeeded  tanzu-package-repo-global
  fluxcd-source-controller  fluxcd-source-controller.community.tanzu.vmware.com  0.21.2           Reconcile succeeded  tanzu-package-repo-global
  knative-serving           knative-serving.community.tanzu.vmware.com           1.0.0            Reconcile succeeded  tanzu-package-repo-global
  kpack                     kpack.community.tanzu.vmware.com                     0.5.2            Reconcile succeeded  tanzu-package-repo-global
  kpack-dependencies        kpack-dependencies.community.tanzu.vmware.com        0.0.9            Reconcile succeeded  tanzu-package-repo-global
  cni                       calico.community.tanzu.vmware.com                    3.22.1           Reconcile succeeded  tkg-system
  secretgen-controller      secretgen-controller.community.tanzu.vmware.com      0.7.1            Reconcile succeeded  tkg-system
  ```

## Create a Tanzu Workload

1. We will use a sample git repo to demonstrate how to create a Tanzu workload from your source code using `tanzu apps`.
  
For using your own repository, please refer to Create a Tanzu workload section in Application Toolkit Package Docs.

  ```shell
  tanzu apps workload create hello-world \
              --git-repo  https://github.com/vmware-tanzu/application-toolkit-sample-app \
              --git-branch main \
              --type web \
              --label app.kubernetes.io/part-of=hello-world \
              --yes \
              --tail \
              --namespace default
  ```

1. Watch the logs of the workload to see it build and deploy. You'll know it's complete when you see `Workload "tanzu-simple-web-app" is ready`

1. After the workload is built and running, you can get the URL of the workload by running the command below.

  ```shell
  tanzu apps workload get hello-world
  ```
  
1. You can then curl the URL using the below command

  `curl http://hello-world.default.127-0-0-1.sslip.io`

  You will see the below output
  `Greetings from Application Toolkit in Tanzu Community Edition!`

## Next Steps

* [Application Toolkit Reference Documentation](package-readme-app-toolkit-0.2.0.md)
