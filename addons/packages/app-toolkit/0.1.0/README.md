# (Experimental) Application Toolkit

Application Toolkit package is a meta-package that installs a collection of packages for creating, iterating and managing applications

## Supported Providers

The following table shows the providers this package can work with.

| Docker |
|:---:|
| ✅  |

## Components

* App Toolkit version: 0.1.0

| Name | Package | Version |
|---------|----------------|---|
| Cartographer | cartographer.community.tanzu.vmware.com | 0.2.2 |
| cert-manager | cert-manager.community.tanzu.vmware.com | 1.6.1 |
| Contour | contour.community.tanzu.vmware.com | 1.20.1 |
| Flux CD Source Controller | fluxcd-source-controller.community.tanzu.vmware.com | 0.21.2 |
| Knative Serving | knative-serving.community.tanzu.vmware.com | 1.0.0 |
| kpack | kpack.community.tanzu.vmware.com | 0.5.1 |

## Configuration

| Config | Values | Description |
|--------|--------|-------------|
| cert_manager | | [See cert-manager documentation](https://tanzucommunityedition.io/docs/v0.11/package-readme-cert-manager-1.6.1/#configuration)|
| contour | | [See contour documentation](https://tanzucommunityedition.io/docs/v0.11/package-readme-contour-1.20.1/#configuration-reference) |
| excluded_packages  | Array of package names | Allows installers to skip deploying named packages |
| knative_serving | | [See knative documentation](https://tanzucommunityedition.io/docs/v0.11/package-readme-knative-serving-1.0.0/#configuration) |
| kpack | | [See kpack documentation](https://tanzucommunityedition.io/docs/v0.11/package-readme-kpack-0.5.1/#kpack-configuration) |

## Installation

### Prerequisites

* The Tanzu CLI is installed, see [Getting Started](https://tanzucommunityedition.io/docs/getting-started/#install-tanzu-cli).
* The Carvel ytt tool is installed, see [Installing ytt](https://carvel.dev/ytt/docs/v0.40.0/install/)
* The following steps require the Tanzu CLI app plugin. See [here](https://github.com/vmware-tanzu/apps-cli-plugin#getting-started) on how to install.
* You will need to provide a container registry for both `kpack` and the sample supply chain to leverage. In the examples below a free DockerHub account is shown, but you can reference the [Tanzu Community Edition kpack documentation](https://tanzucommunityedition.io/docs/v0.10/package-readme-kpack-0.5.0/) to see examples utilizing other registries.
* A targeted Tanzu unmanaged cluster, see [Usage Example](#usage-example)

### Prepare a values.yaml file for local Docker installation without load balancer

In the following example, you will deploy our traffic ingress mechanism, Contour, with a ClusterIP configuration, and Knative to serve as localhost. You will add these settings to a `values.yaml` file, and add a Docker registry for kpack to store buildpacks on. This example is one method of configuring App-toolkit, feel free to adapt the configuration to reflect your environment.

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
  kp_default_repository: https://index.docker.io/v1/
  kp_default_repository_username: [your username]
  kp_default_repository_password: [your password]
```

### Install App-toolkit using the custom values file

You can also change the default installation values by providing a file with predetermined values. In this documentation, we will refer to this file as `values.yaml`.

```shell
tanzu package install app-toolkit --package-name app-toolkit.community.tanzu.vmware.com --version 0.1.0 -f values.yaml -n tanzu-package-repo-global
```

### Optional: Exclude packages

To exclude certain packages, create or add to a `values.yaml` file the name of the excluded packages.

For example:

```yaml
excluded_packages:
  - contour.community.tanzu.vmware.com
  - cert-manager.community.tanzu.vmware.com
```

## Usage Example

This example illustrates creating a simple Spring Boot web app on the locally deployed Tanzu Community Edition cluster using Docker.

1. Begin with Docker and the Tanzu CLI installed on your development machine and ensure you have the Carvel tools installed as well as we will be leveraging `ytt` in the examples below. For best performance ensure that you have plenty of free space on the local filesystem that Docker is using.

1. Create an unmanaged cluster named demo-local exposing 80 and 443 ports.

    ```shell
    tanzu unmanaged-cluster create demo-local -p 80:80 -p 443:443
    ```

1. Ensure you have installed the App-toolkit package. You can find an example of the `values.yaml` above in the [installation section](#prepare-a-valuesyaml-file-for-local-docker-installation-without-load-balancer):

    ```shell
    tanzu package install app-toolkit --package-name app-toolkit.community.tanzu.vmware.com --version 0.1.0 -f values.yaml -n tanzu-package-repo-global
    ```

    You can validate this by checking that all the packages have successfully reconciled using the command:

    ```shell
    tanzu package installed list -A
    ```

    You should see output similar to the following:

    ```shell
    tanzu package installed list -A
    | Retrieving installed packages...
      NAME                      PACKAGE-NAME                                         PACKAGE-VERSION        STATUS               NAMESPACE
      app-toolkit               app-toolkit.community.tanzu.vmware.com               0.1.0                  Reconcile succeeded  tanzu-package-repo-global
      cartographer              cartographer.community.tanzu.vmware.com              0.2.2                  Reconcile succeeded  tanzu-package-repo-global
      cert-manager              cert-manager.community.tanzu.vmware.com              1.6.1                  Reconcile succeeded  tanzu-package-repo-global
      contour                   contour.community.tanzu.vmware.com                   1.19.1                 Reconcile succeeded  tanzu-package-repo-global
      fluxcd-source-controller  fluxcd-source-controller.community.tanzu.vmware.com  0.21.2                 Reconcile succeeded  tanzu-package-repo-global
      knative-serving           knative-serving.community.tanzu.vmware.com           1.0.0                  Reconcile succeeded  tanzu-package-repo-global
      kpack                     kpack.community.tanzu.vmware.com                     0.5.1                  Reconcile succeeded  tanzu-package-repo-global
      cni                       antrea.tanzu.vmware.com                              0.13.3+vmware.1-tkg.1  Reconcile succeeded  tkg-system
      ```

1. You'll now need to prepare a second values file (apart from the one used to install App-toolkit) with the additional configuration necessary for `kpack` and `cartographer`. In the below examples, this second values file will be referenced as `supplychain-example-values.yaml`.

    ```yaml
    kpack:
      registry:
        url: <REGISTRY_URL>
        username: <REGISTRY_USERNAME>
        password: <REGISTRY_PASSWORD>
      builder:
        # path to the container repository where kpack build artifacts are stored
        tag: <REGISTRY_TAG>
      # A comma-separated list of languages e.g. [java,nodejs] that will be supported for development
      # Allowed values are:
      # - java
      # - nodejs
      # - dotnet-core
      # - go
      # - ruby
      # - php
      languages: [java]
      image_prefix: <REGISTRY_PREFIX>

    ```

    Where:
    * `<REGISTRY_URL>` is a valid OCI registry to store kpack images. For DockerHub this would be `https://index.docker.io/v1/`
    * `<REGISTRY_USERNAME>` and `REGISTRY_PASSWORD` are the credentials for the specified registry.
    * `<REGISTRY_TAG>` is the path to the container repository where kpack build artifacts are stored. For DockerHub, it is username/tag, e.g. csamp/builder
    * `<REGISTRY_PREFIX>` is the prefix for your images as they reside on the registry. For DockerHub, it is the username followed by a trailing slash, e.g. csamp/

1. Use `ytt` to insert the values from the above `supplychain-example-values.yaml` into the Kubernetes resource definitions below. You can inspect the files and see which resources each one is creating

    ```shell
    ytt --data-values-file supplychain-example-values.yaml --ignore-unknown-comments -f https://raw.githubusercontent.com/vmware-tanzu/community-edition/main/addons/packages/app-toolkit/0.1.0/test/example_sc/developer-namespace.yaml | kubectl apply -f -
    ytt --data-values-file supplychain-example-values.yaml --ignore-unknown-comments -f https://raw.githubusercontent.com/vmware-tanzu/community-edition/main/addons/packages/app-toolkit/0.1.0/test/example_sc/kpack-dependencies.yaml | kubectl apply -f -
    ytt --data-values-file supplychain-example-values.yaml --ignore-unknown-comments -f https://raw.githubusercontent.com/vmware-tanzu/community-edition/main/addons/packages/app-toolkit/0.1.0/test/example_sc/supplychain-templates.yaml | kubectl apply -f -
    ```

1. Deploy the example supply chain

    ```shell
    ytt --data-values-file supplychain-example-values.yaml --ignore-unknown-comments -f https://raw.githubusercontent.com/vmware-tanzu/community-edition/main/addons/packages/app-toolkit/0.1.0/test/example_sc/supplychain.yaml | kubectl apply -f -
    ```

1. Deploy the Tanzu workload

    ```shell
    tanzu apps workload create tanzu-simple-web-app --git-repo https://github.com/cgsamp/tanzu-simple-web-app --git-branch main --label apps.tanzu.vmware.com/workload-type=web
    ```

1. Watch the logs of the workload to see it build and deploy. You'll know it's complete when you see `Build successful`

    ```shell
    tanzu apps workload tail tanzu-simple-web-app
    ```

1. Access the workload's url in your browser, or with curl

    ```shell
    curl http://tanzu-simple-web-app.default.127-0-0-1.sslip.io/
    ```
  
## Troubleshooting

Error: invalid configuration: no configuration has been provided

* Make sure you are providing a valid `values.yaml` file to the `tanzu package install` command
* You also must have a TCE unmanaged cluster targeted to install App Toolkit

Insufficient CPU or Memory Error

* Make sure the environment you're running your cluster in has enough resources allocated. You can find TCE's unmanaged-cluster specifications [here](https://tanzucommunityedition.io/docs/v0.11/getting-started-unmanaged/)

sample-app deploy fails with MissingValueAtPath error

* Double check the formatting for the registry credentials provided in [Usage Example](#usage-example). Different registry types expect different formats for each of the fields.

Error when you curl the `tanzu-simple-web-app` url

* The service can sometimes take a minute or two to setup, even after the build shows a success with `tanzu app workload tail tanzu-simple-web-app`
* You can also double check that the knative service for the tanzu-simple-web-app was created and is running by checking `kubectl get ksvc`
