# App-toolkit 0.2.0

Application Toolkit package is a package that installs a set of packages for creating, iterating and managing applications.

For a simple tutorial on how to use Application Toolkit, please refer to the ["Getting Started with Creating a Tanzu Workload" guide](https://tanzucommunityedition.io/docs/v0.12/getting-started-tanzu-workloads/).

This documentation provides information about the specific TCE package. Please visit the [TCE package management page](https://tanzucommunityedition.io/docs/v0.12/package-management/) for general information about installation, removal, troubleshooting, and other topics.

## Installation

### Installation of dependencies

Ensure you have these in place:

* the Tanzu CLI is installed; see [Getting Started](https://tanzucommunityedition.io/docs/v0.12/getting-started/#install-tanzu-cli).
* an unmanaged cluster.
* an OCI Compliant Container registry credentials to be used in kpack, kpack-dependencies, and cartographer-catalog package configuration, and for secret creation.
* a load balancer or Ingress configuration to be used in conjunction with Contour and Knative-Serving.

### Installation of package

#### Step 1: Create the registry secret

Install `secretgen-controller` based on the version in the secretgen-controller package docs. For example, if the version is 0.7.1:

```shell
tanzu package install secretgen-controller --package-name secretgen-controller.community.tanzu.vmware.com --version 0.7.1 -n tkg-system
```

Then create a registry secret `registry-credentials` using the below command:

```shell
tanzu secret registry add registry-credentials --server REGISTRY_URL --username REGISTRY_USER --password REGISTRY_PASS --export-to-all-namespaces
```

* `REGISTRY_URL` - URL for the registry you plan to upload your builds to.
  * For Dockerhub, it would be <https://index.docker.io/v1/>
  * For GCR, it would be <gcr.io>
  * For Harbor, it would be <myharbor.example.com>
* `REGISTRY_USER`: the username for the account with write access to the registry specified with `REGISTRY_URL`
* `REGISTRY_PASS`: the password for the same account. If your password includes special characters, you'll want to double-check that the credential is populated correctly. You can also use the `--pasword-env-var`, `--password-file`, or `--password-stdin` options to provide your password, if you prefer.

### Step 2: Prepare an app-toolkit-values.yaml

As mentioned in the prerequisite, ensure you provide the image registry configuration and ingress configuration while installing the App-Toolkit package.

```yaml
contour:

knative_serving:

kpack:
  kp_default_repository:
  kp_default_repository_username:
  kp_default_repository_password:

cartographer_catalog:
  registry:
    server:
    repository:

# The namespace field below will configure the namespace for creating workloads.
developer_namespace:  #default value is default

# The excluded_packages field consists of packages you do not want to install.
# Below is an example of how you can provide the packages you want to exclude.

excluded_packages:
   # - contour.community.tanzu.vmware.com
   # - cert-manager.community.tanzu.vmwware.com
```

### Step 3: Install App-toolkit Package

Install the 0.2.0 version of Application Toolkit.

```shell
tanzu package install app-toolkit --package-name app-toolkit.community.tanzu.vmware.com --version 0.2.0 -f app-toolkit-values.yaml -n tanzu-package-repo-global
```

To check that all the packages have successfully reconciled, use the command:

```shell
tanzu package installed list -A
```

You should see output consisting of the below packages, among other packages you may have installed.

```shell
NAME                      PACKAGE-NAME                                         PACKAGE-VERSION  STATUS               NAMESPACE
  secretgen-controller      secretgen-controller.community.tanzu.vmware.com      0.8.0            Reconcile succeeded  default
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
```

## Options

### Package configuration values

App Toolkit is a meta-package that predominantly passes configuration down to the packages it installs.

|Value | Required/Optional | Default | Description |
| cartographer_catalog | Required | | [See Cartographer catalog documentation](https://tanzucommunityedition.io/docs/v0.12/package-readme-cartographer-catalog-0.3.0/#configuration)|
| contour | Required | | [See Contour documentation](https://tanzucommunityedition.io/docs/v0.11/package-readme-contour-1.20.1/#configuration-reference) |
| excluded_packages | Optional | None | Allows installers to skip deploying named packages, specified as an array of package names |
| knative_serving | Required | | [See Knative-Serving documentation](https://tanzucommunityedition.io/docs/v0.11/package-readme-knative-serving-1.0.0/#configuration) |
| kpack | Required |  | [See kpack documentation](https://tanzucommunityedition.io/docs/v0.11/package-readme-kpack-0.5.2/#kpack-configuration) |
| cert_manager | Optional | | [See cert-manager documentation](https://tanzucommunityedition.io/docs/v0.11/package-readme-cert-manager-1.6.1/#configuration)|
| developer_namespace | Optional | `default` | Configures the namespace with the required secret, service binding and role binding to create Tanzu workloads |
| kpack_dependencies | Optional | | [See kpack dependencies documentation](https://tanzucommunityedition.io/docs/v0.12/package-readme-kpack-dependencies-0.0.9/#kpack-dependencies-configuration) |

### Application configuration values

Not applicable for this package.

## What This Package Does

App Toolkit installs open source components of the Tanzu Application Platform that, after being setup, allow the user to simply push a git commit and have their code deployed and accessible on their cluster.

## Components

| Name                                                                                                                     | Description                                                                                                                                 | Version |
| ------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| [Cartographer](https://tanzucommunityedition.io/docs/v0.12/package-readme-cartographer-0.3.0/)                           | Allows you to create secure and reusable supply chains that define all of your application continuous integration and continuous delivery in one place, in-cluster. | 0.3.0   |
| [Cartographer-Catalog](https://tanzucommunityedition.io/docs/v0.12/package-readme-cartographer-catalog-0.3.0/)           | Reusable Cartographer supply chains and templates for driving workloads from source code to running Knative service in a cluster.           | 0.3.0   |
| [cert-manager](https://tanzucommunityedition.io/docs/v0.11/package-readme-cert-manager-1.6.1/)                           | Provides certificate management functionality.                                                                                 | 1.6.1   |
| [Contour](https://tanzucommunityedition.io/docs/v0.11/package-readme-contour-1.20.1/)                                    | Provides Ingress capabilities for Kubernetes clusters.                                                                               | 1.20.1  |
| [Flux CD Source Controller](https://tanzucommunityedition.io/docs/v0.11/package-readme-fluxcd-source-controller-0.21.2/) | Specializes in artifact acquisition from external sources such as Git, Helm repositories, and S3 buckets.                      | 0.21.2  |
| [Knative Serving](https://tanzucommunityedition.io/docs/v0.11/package-readme-knative-serving-1.0.0/)                     | Provides the ability for users to create serverless workloads from OCI images.                                                                                                 | 1.0.0   |
| [kpack](https://tanzucommunityedition.io/docs/v0.11/package-readme-kpack-0.5.2/)                                         | Provides a platform for building OCI images from source code.
| [kpack-dependencies](https://tanzucommunityedition.io/docs/v0.12/package-readme-kpack-dependencies-0.0.9/)               | Provides a curated set of buildpacks and stacks required by kpack.                                                       | 0.0.9   |

### Supported Providers

Application Toolkit is currently tested with Unmanaged Clusters on any of the below providers:

| Unmanaged Clusters| AWS | Azure | vSphere | Docker |
|-------------------|-----|-------|---------|--------|
| ✅                | ✅   | ✅     | ✅       | ✅      |

## Files

Not relevant for this particular package.

## Package Limitations

By default, we enable a single namespace with the credentials needed to enable the supply chain (configurable with the `developer_namespace` option). If you would like to enable App Toolkit in additional namespaces, you will need to make a `registry-credentials` secret in the namespace.

The most common issue we have seen with installing this package is authenticating against the image repositories and getting 403 errors. If this happens to you, we encourage you to double check the configuration against the docs.

If you come across any issues or bugs, you can report them to the team [here](https://github.com/vmware-tanzu/package-for-application-toolkit/issues).

## Usage Example

### Example App Toolkit Configuration

In this example we will deploy Contour, our traffic ingress mechanism, with a ClusterIP configuration, and Knative to serve as localhost. You will add these configurations to a `app-toolkit-values.yaml` file, and add a Docker registry for kpack to store buildpacks on.

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
  kp_default_repository: index.docker.io/my-dockerhub-username/my-repo
  kp_default_repository_username: my-dockerhub-username
  kp_default_repository_password: my-dockerhub-password

cartographer_catalog:
  registry:
    server: index.docker.io
    repository: my-dockerhub-username
```

1. Ensure you have followed the steps for Installing the App Toolkit Package.

2. If you want to create applications in a namespace different from the one you configured during App Toolkit installation, please [set up the developer namespace](#set-up-the-developer-namespace) before proceeding to create the Tanzu workload.

3. Create a Tanzu workload using the `tanzu apps workload create` command. You can supply your own git repo, or use the sample git repo below.

    ```shell
    tanzu apps workload create hello-world \
                --git-repo GIT-URL-TO-PROJECT-REPO \
                --git-branch main \
                --type web \
                --label app.kubernetes.io/part-of=hello-world \
                --yes \
                -n YOUR-DEVELOPER-NAMESPACE
    ```

    where GIT-URL-TO-PROJECT-REPO is your git repository and YOUR-DEVELOPER-NAMESPACE is the namespace configured while installing the package.

4. Watch the logs of the workload to see it build and deploy. You'll know it's complete when you see `Build successful`

    ```shell
    tanzu apps workload tail hello-world -n YOUR-DEVELOPER-NAMESPACE
    ```

5. After the workload is built and running, get its URL by running the command below:

    ```shell
    tanzu apps workload get hello-world --namespace YOUR-DEVELOPER-NAMESPACE
    ```

### Set up the Developer Namespace

If you want to create workloads in additional namespaces, follow these steps to set up each developer namespace:

1. Ensure you have access to kubectl.

2. Create the additional namespace using the below command:

    ```shell
    kubectl create namespace YOUR-DEVELOPER-NAMESPACE
    ```

3. Use this kubectl command to create the secret, a service account, and a role binding in the developer namespace:

    ```shell
    cat <<EOF | kubectl -n YOUR-DEVELOPER-NAMESPACE create -f -
    apiVersion: v1
    kind: Secret
    metadata:
      name: registry-credentials
      annotations:
        secretgen.carvel.dev/image-pull-secret: ""
    type: kubernetes.io/dockerconfigjson
    data:
      .dockerconfigjson: e30K
    ---
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: workload-user-sa
    secrets:
    - name: registry-credentials
    imagePullSecrets:
    - name: registry-credentials
    ```

## Troubleshooting

### Insufficient CPU or Memory Error

* Make sure the environment you're running your cluster in has enough resources allocated. You can find TCE's unmanaged-cluster specifications [here](https://tanzucommunityedition.io/docs/v0.12/support-matrix/).

### Sample App deploy fails with MissingValueAtPath error

* Double-check the formatting for the registry credentials provided in our [Usage Example](#usage-example). Different registry types expect different formats for each of the fields.

### Error when you curl the `tanzu-simple-web-app` url

* The service can sometimes take a minute or two to set up, even after the build shows a success with `tanzu app workload tail tanzu-simple-web-app`.
* You can also double-check that the Knative service for the tanzu-simple-web-app was created and is running by checking `kubectl get ksvc`.

If you come across any issues or bugs, you can report them to the team [here](https://github.com/vmware-tanzu/package-for-application-toolkit/issues).

## Additional Documentation

For a getting started guide on how to use App Toolkit, you can check out “[Getting Started with Tanzu Workloads](https://tanzucommunityedition.io/docs/v0.12/getting-started-tanzu-workloads/)”
