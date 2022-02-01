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
| cert_manager | | [See cert-manager documentation](https://tanzucommunityedition.io/docs/latest/package-readme-cert-manager-1.6.1/#configuration)|
| contour | | [See contour documentation](https://tanzucommunityedition.io/docs/latest/package-readme-contour-1.20.1/#configuration-reference) |
| excluded_packages  | Array of package names | Allows installers to skip deploying named packages |
| knative_serving | | [See knative documentation](https://tanzucommunityedition.io/docs/latest/package-readme-knative-serving-1.0.0/#configuration) |
| kpack | | [See kpack documentation](https://tanzucommunityedition.io/docs/latest/package-readme-kpack-0.5.1/#kpack-configuration) |

## Installation

### Prerequisites

* The following steps require the Tanzu CLI app plugin. See [here](https://github.com/vmware-tanzu/apps-cli-plugin#getting-started) on how to install.

### Steps

The App Toolkit package is a meta-packages - a collection of packages - to build and deploy workloads. It requires that TCE is installed and the Tanzu CLI is available.

If a load balancer is available in your cluster, no value definitions are required.

```shell
tanzu package install app-toolkit --package-name app-toolkit.community.tanzu.vmware.com --version 0.1.0
```

### Using custom values

You can also change the default installation values by providing a file with predetermined values. In this documentation we will refer to this file as `values.yaml`.

```shell
tanzu package install app-toolkit --package-name app-toolkit.community.tanzu.vmware.com --version 0.1.0 -f values.yaml
```

### Exclude packages

To exclude certain packages, create a `values.yaml` file and name the excluded packages.

For example:

```yaml
excluded_packages:
  - contour.community.tanzu.vmware.com
  - cert-manager.community.tanzu.vmware.com
```

### Local Docker installation without load balancer

To deploy contour with a ClusterIP configuration, and Knative to serve at localhost, add these settings to a `values.yaml` file. Also, add a Docker registry for kpack to store produced container images.

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

## Usage Example

This example illustrates creating a simple Spring Boot web app on a locally deployed TCE cluster using Docker.

1. Begin with Docker and the Tanzu CLI installed on your development machine.

  For best performance, configure Docker with at least 8GB of RAM.

1. Create a cluster named demo-local exposing 80 and 443 ports.

    ```shell
    tanzu unmanaged-cluster create demo-local -p 80:80 -p 443:443
    ```

1. Follow the steps listed above in the "Local Docker installation without load balancer" section

1. Set up kpack and required resources

    ```shell
    ytt --data-values-file values.yaml --ignore-unknown-comments -f https://raw.githubusercontent.com/cgsamp/tanzu-simple-web-app/main/config/example_sc/rbac.yaml | kubectl apply -f -
    ytt --data-values-file values.yaml --ignore-unknown-comments -f https://raw.githubusercontent.com/cgsamp/tanzu-simple-web-app/main/config/example_sc/kpack-templates.yaml | kubectl apply -f -
    ytt --data-values-file values.yaml --ignore-unknown-comments -f https://raw.githubusercontent.com/cgsamp/tanzu-simple-web-app/main/config/example_sc/builder.yaml | kubectl apply -f -
    ```

1. Deploy the example supply chain

    ```shell
    ytt --data-values-file values.yaml --ignore-unknown-comments -f https://raw.githubusercontent.com/cgsamp/tanzu-simple-web-app/main/config/example_sc/supplychain.yaml | kubectl apply -f -
    ```

1. Deploy the Tanzu workload

    ```shell
    tanzu apps workload create -f https://raw.githubusercontent.com/cgsamp/tanzu-simple-web-app/main/example/workload.yaml
    ```

1. Watch the logs of the workload to see it build and deploy

    ```shell
    tanzu apps workload tail tanzu-simple-web-app
    ```

1. Access the workload's url in your browser, or with curl

    ```shell
    curl http://tanzu-simple-web-app.default.127-0-0-1.sslip.io/
    ```
