# (Experimental) Application Toolkit

Application Toolkit package is a package that installs a set of packages for creating, iterating and managing applications.

For a simple tutorial on how to use Application Toolkit, please refer to the "Getting Started with Creating a Tanzu Workload" guide.

## Supported Providers

Application Toolkit is currently tested with Unmanaged Clusters on any of the below providers.

| Unmanaged Clusters| AWS | Azure | vSphere | Docker |
|-------------------|-----|-------|---------|--------|
| ✅                | ✅   | ✅     | ✅       | ✅      |

## Components

### App Toolkit 0.2.0

| Name                                                                                                                     | Description                                                                                                                                 | Version |
| ------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| [Cartographer](https://tanzucommunityedition.io/docs/v0.11/package-readme-cartographer-0.3.0/)                           | Cartographer allows you to create secure and reusable supply chains that define all of your application CI and CD in one place, in cluster. | 0.3.0   |
| [Cartographer-Catalog](https://tanzucommunityedition.io/docs/v0.11/package-readme-cartographer-catalog-0.3.0/)           | Reusable Cartographer Supply Chains and templates for driving workloads from source code to running Knative service in a cluster.           | 0.3.0   |
| [cert-manager](https://tanzucommunityedition.io/docs/v0.11/package-readme-cert-manager-1.6.1/)                           | Cert Manager provides certificate management functionality.                                                                                 | 1.6.1   |
| [Contour](https://tanzucommunityedition.io/docs/v0.11/package-readme-contour-1.20.1/)                                    | Contour provides Ingress capabilities for Kubernetes clusters                                                                               | 1.20.1  |
| [Flux CD Source Controller](https://tanzucommunityedition.io/docs/v0.11/package-readme-fluxcd-source-controller-0.21.2/) | FluxCD Source specialises in artifact acquisition from external sources such as Git, Helm repositories and S3 buckets.                      | 0.21.2  |
| [Knative Serving](https://tanzucommunityedition.io/docs/v0.11/package-readme-knative-serving-1.0.0/)                     | Knative Serving provides the ability for users to create serverless workloads from OCI images                                                                                                 | 1.0.0   |
| [kpack](https://tanzucommunityedition.io/docs/v0.11/package-readme-kpack-0.5.2/)                                         | kpack provides a platform for building OCI images from source code.
| [kpack-dependencies](https://tanzucommunityedition.io/docs/v0.11/package-readme-kpack-dependencies-0.0.9/)               | kpack-dependencies provides a curated set of buildpacks and stacks required by kpack.                                                       | 0.0.9   |

## Configuration

| Config | Values | Description |
|--------|--------|-------------|
| cartographer-catalog | | [See cartographer catalog documentation](https://tanzucommunityedition.io/docs/package-readme-cartographer-catalog-0.3.0/#configuration)|
| cert_manager | | [See cert-manager documentation](https://tanzucommunityedition.io/docs/package-readme-cert-manager-1.6.1/#configuration)|
| contour | | [See contour documentation](https://tanzucommunityedition.io/docs/package-readme-contour-1.20.1/#configuration-reference) |
| excluded_packages  | Array of package names | Allows installers to skip deploying named packages |
| knative_serving | | [See knative documentation](https://tanzucommunityedition.io/docs/package-readme-knative-serving-1.0.0/#configuration) |
| kpack | | [See kpack documentation](https://tanzucommunityedition.io/docs/package-readme-kpack-0.5.2/#kpack-configuration) |
| kpack-dependencies | | [See kpack dependencies documentation](https://tanzucommunityedition.io/docs/package-readme-kpack-dependencies-0.5.2/#kpack-dependencies-configuration) |
| developer-namespace | (default value is `default`) | Configures the namespace with the required secret, service binding and role binding to create Tanzu workloads |

## Installing the App Toolkit Package

### Before you begin

* Ensure the Tanzu CLI is installed, see [Getting Started](https://tanzucommunityedition.io/docs/v0.12/getting-started/#install-tanzu-cli).
* Ensure you have created an unmanaged cluster.
* An OCI Compliant Container registry credentials to be used in kpack, kpack-dependencies, cartographer-catalog package configuration and for secret creation.
* A Load Balancer or Ingress configuration to be used in conjunction with Contour and Knative-Serving

### Step 1: Create the registry secret

 1. Install `secretgen-controller` based on the version in the secretgen-controller package docs. For example, if the version is 0.8.0

    ```shell
    tanzu package install secretgen-controller --package-name secretgen-controller.community.tanzu.vmware.com --version 0.8.0
    ```

 1. Create a registry secret `registry-credentials` using the below command

    ```shell
    tanzu secret registry add registry-credentials --server REGISTRY_URL --username REGISTRY_USER --password REGISTRY_PASS --export-to-all-namespaces`
    ```

      * `REGISTRY_URL` - URL for the registry you plan to upload your builds to.
        * For Dockerhub, it would be <https://index.docker.io/v1/>
        * For GCR, it would be <gcr.io>
        * For Harbor, it would be <myharbor.example.com>
      * `REGISTRY_USER`: the username for the account with write access to the registry specified with `REGISTRY_URL`
      * `REGISTRY_PASS`: the password for the same account. If you have special characters in your password, you'll want to double check that the credential is populated correctly. You can also use the `--pasword-env-var`, `--password-file`, or `--password-stdin` options to provide your password if you prefer

### Step 2: Prepare an app-toolkit-values.yaml

As mentioned in the pre-requisite, ensure you provide the image registry configuration and ingress configuration while installing App-Toolkit package.

```yaml
contour:

knative_serving:

kpack:
  kp_default_repository:
  kp_default_repository_username:
  kp_default_repository_password:

cartographer-catalog:
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

#### Example App Toolkit Configuration

In the following example, we will deploy our traffic ingress mechanism, Contour, with a ClusterIP configuration, and Knative to serve as localhost. You will add these configurations to a `app-toolkit-values.yaml` file, and add a Docker registry for kpack to store buildpacks on.

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
  kp_default_repository: docker.io/my-dockerhub-username/my-repo
  kp_default_repository_username: my-dockerhub-username
  kp_default_repository_password: my-dockerhub-password

cartographer-catalog:
  registry:
      server: index.docker.io
      repository: my-dockerhub-username
```

### Step 3: Install App-toolkit Package

1. Install the 0.2.0 version of Application Toolkit.

```shell
tanzu package install app-toolkit --package-name app-toolkit.community.tanzu.vmware.com --version 0.2.0 -f app-toolkit-values.yaml -n tanzu-package-repo-global
```

1. You can validate this by checking that all the packages have successfully reconciled using the command:

```shell
tanzu package installed list -A
```

You should see output consisting of the below packages among other packages you may have installed.

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

## Usage

1. Ensure you have followed the steps for Installing the App Toolkit Package.

2. In case you want to create applications in a namespace different than the one you configured during App Toolkit installation, please complete the "Set-up the Developer Namespace" before proceeding to create the Tanzu workload.

3. Please create a workload.yaml with the minimal information as provided below:

4. Create a Tanzu workload using the `tanzu apps workload create` command. You can supply your own git repo, or use the sample git repo below.

    ```shell
    tanzu apps workload create hello-world \
                --git-repo GIT-URL-TO-PROJECT-REPO \
                --git-branch main \
                --type web \
                --label app.kubernetes.io/part-of=hello-world \
                --yes \
                -n YOUR-DEVELOPER-NAMESPACE
    ```

    where GIT-URL-TO-PROJECT-REPO is your git repositoryand YOUR-DEVELOPER-NAMESPACE is the namespace configured while installing the package.

5. Watch the logs of the workload to see it build and deploy. You'll know it's complete when you see `Build successful`

    ```shell
    tanzu apps workload tail hello-world -n YOUR-DEVELOPER-NAMESPACE
    ```

6. After the workload is built and running, you can get the URL of the workload by running the command below.

    ```shell
    tanzu apps workload get hello-world --namespace YOUR-DEVELOPER-NAMESPACE
    ```

### Set-up the Developer Namespace

In case you want to create workloads in additional namespaces, follow the below procedure to set-up each developer namespace.

1. Ensure you have access to kubectl.

2. Create the additional namespace using the below command

    ```shell
    kubectl create namespace YOUR-DEVELOPER-NAMESPACE
    ```

3. Create the secret, a service account and a role binding in the developer namespace using the below kubectl command

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

* Make sure the environment you're running your cluster in has enough resources allocated. You can find TCE's unmanaged-cluster specifications [here](https://tanzucommunityedition.io/docs/v0.11/support-matrix/)

### Sample App deploy fails with MissingValueAtPath error

* Double check the formatting for the registry credentials provided in [Usage Example](#usage-example). Different registry types expect different formats for each of the fields.

### Error when you curl the `tanzu-simple-web-app` url

* The service can sometimes take a minute or two to setup, even after the build shows a success with `tanzu app workload tail tanzu-simple-web-app`
* You can also double check that the knative service for the tanzu-simple-web-app was created and is running by checking `kubectl get ksvc`
