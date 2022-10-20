# kpack dependencies

[kpack](https://github.com/pivotal/kpack) utilizes unprivileged Kubernetes
primitives to provide builds of OCI images as a platform implementation
of [Cloud Native Buildpacks](https://buildpacks.io) (CNB).

kpack requires dependencies in the form of buildpacks and stacks to keep app images patched and up-to-date.

This package provides a curated set of dependencies to be used with the [kpack package](https://github.com/vmware-tanzu/package-for-kpack).

## Components

* kpack dependencies

## Supported Providers

The following table shows the providers this package can work with.

| AWS | Azure | vSphere | Docker |
|-----|-------|---------|--------|
| ✅   | ✅     | ✅       | ✅      |

## Configuration

The following configuration values can be set to customize the kpack
installation.

### kpack Configuration

| Value                    | Required/Optional | Description                                                                                                                                                                     |
|--------------------------|-------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `kp_default_repository`  | Required          | Docker repository used for builder images (i.e. `gcr.io/my-project/my-repo` or `mydockerhubusername/my-repo`. **This must be the same value used during installation of kpack.**|

## Installation

### Prerequisites

* The [kpack package](https://github.com/vmware-tanzu/package-for-kpack) is installed

### Package Installation steps

You can install the kpack dependencies package using the command below:

`tanzu package install kpack-dependencies --package-name kpack-dependencies.community.tanzu.vmware.com --version 0.0.1 -f kpack-deps-values.yaml`

#### Verification

Once the package is installed you can view the resources that have been created:

**NOTE: These resources cannot be modified manually, they can only be upgraded via upgrades of the kpack dependencies package. If you wish to create custom ClusterStores, ClusterStacks, or ClusterBuilders you must create [new resources](https://github.com/pivotal/kpack/blob/main/docs/tutorial.md) and manage them manually.**

```bash
$ kubectl get clusterstore
NAME      READY
default   True

$ kubectl get clusterstack
NAME      READY
base      True
default   True

$ kubectl get clusterbuilder
NAME      LATESTIMAGE    READY
base      <some-image>   True
default   <some-image>   True
```

#### Troubleshooting

Currently, the kpack dependencies package will not immediately fail if the installation is in a bad state.

If your installation is reconciling for a long time or receives a timeout, check the status of the relevant resources:

```bash
kubectl describe clusterstore
kubectl describe clusterstack
kubectl describe clusterbuilder
```

## Configuration and Usage

After installing the kpack dependencies package, ClusterStores, ClusterStacks, and ClusterBuilders do not need to be manually installed, and you can immediately create source-to-image builds.

* [Creating an image using kp](#creating-an-image-using-kp-cli)
* [Creating an image using kubectl](#creating-an-image-using-kubectl)

### Creating an image using kp CLI

1. Create a secret with push credentials for the docker registry that you plan on publishing OCI images to with kpack.

   The easiest way to do that is with `kp secret save`

    ```bash
    kp secret save tutorial-registry-credentials \
       --registry <REGISTRY-HOSTNAME> \
       --registry-user <REGISTRY-USER>
    ```

   > Note: The `<REGISTRY-HOSTNAME>` must be the registry prefix for its corresponding registry
   * For [dockerhub](https://hub.docker.com/) this should be `https://index.docker.io/v1/`. `kp` also offers a simplified way to create a dockerhub secret with a `--dockerhub` flag.
   * For [GCR](https://cloud.google.com/container-registry/) this should be `gcr.io`. If you use GCR then the username can be `_json_key` and the password can be the JSON credentials you get from the GCP UI (under `IAM -> Service Accounts` create an account or edit an existing one and create a key with type JSON). `kp` also offers a simplified way to create a gcr secret with a `--gcr` flag.

   Your secret create should look something like this:

    ```bash
    kp secret save tutorial-registry-credentials \
       --registry https://index.docker.io/v1/ \
       --registry-user my-dockerhub-username
    ```

2. Create a kpack Image Resource.

   An Image Resource is the specification for an OCI image that kpack should build and manage.

   We will create a sample Image Resource that builds with the `base` builder installed with the kpack dependencies package.

   The example included here utilizes the [Spring Pet Clinic sample app](https://github.com/spring-projects/spring-petclinic).
   We encourage you to substitute it with your own application.

   Create an Image Resource:

    ```yaml
    kp image save tutorial-image \
      --tag <IMAGE-TAG> \
      --git https://github.com/spring-projects/spring-petclinic \
      --git-revision 82cb521d636b282340378d80a6307a08e3d4a4c4 \
      --builder base
    ```

   * Make sure to replace `<IMAGE-TAG>` with the tag in the registry of the
     secret you configured. Something like:
     `your-name/app` or `gcr.io/your-project/app`
   * If you are using your application source, replace `--git`
     & `--git-revision`.
   > Note: To use a private git repo follow the instructions in [secrets](https://github.com/pivotal/kpack/blob/main/docs/secrets.md)

   You can now check the status of the Image Resource.

   ```bash
   kp image status tutorial-image
   ```

   You should see that the Image Resource has a status Building as it is
   currently building.

    ```text
    Status:         Building
    Message:        --
    LatestImage:    --

    Source
    Type:        GitUrl
    Url:         https://github.com/spring-projects/spring-petclinic
    Revision:    82cb521d636b282340378d80a6307a08e3d4a4c4

    Builder Ref
    Name:    base
    Kind:    Builder

    Last Successful Build
    Id:              --
    Build Reason:    --

    Last Failed Build
    Id:              --
    Build Reason:    --
    ```

   You can tail the logs for Image Resource that is currently building using
   the [kp cli](https://github.com/vmware-tanzu/kpack-cli/blob/main/docs/kp_build_logs.md)

    ```bash
    kp build logs tutorial-image
    ```

   Once the Image Resource finishes building you can get the fully resolved
   built OCI image with `kp`

    ```bash
    kp image status tutorial-image
    ```

   The latest built OCI image is available to be used locally via `docker pull`
   and in a Kubernetes deployment.

3. Run the built application image locally.

   Download the latest built OCI image available and run it with
   Docker.

   ```bash
   docker run -p 8080:8080 <latest-image-with-digest>
   ```

   You should see the java app start up:

   ```text

              |\      _,,,--,,_
             /,`.-'`'   ._  \-;;,_
    _______ __|,4-  ) )_   .;.(__`'-'__     ___ __    _ ___ _______
    |       | '---''(_/._)-'(_\_)   |   |   |   |  |  | |   |       |
    |    _  |    ___|_     _|       |   |   |   |   |_| |   |       | __ _ _
    |   |_| |   |___  |   | |       |   |   |   |       |   |       | \ \ \ \
    |    ___|    ___| |   | |      _|   |___|   |  _    |   |      _|  \ \ \ \
    |   |   |   |___  |   | |     |_|       |   | | |   |   |     |_    ) ) ) )
    |___|   |_______| |___| |_______|_______|___|_|  |__|___|_______|  / / / /
    ==================================================================/_/_/_/

    :: Built with Spring Boot :: 2.2.2.RELEASE
   ```

4. Rebuilding kpack Image Resources.

   We recommend updating the kpack Image Resource with a CI/CD tool when new
   commits are ready to be built.
   > Note: You can also provide a branch or tag as the `spec.git.revision` and kpack will poll and rebuild on updates!

   We can simulate an update from a CI/CD tool by updating
   the `spec.git.revision` on the Image Resource.

   If you are using your own application please push an updated commit and use
   the new commit sha. If you are using Spring Pet Clinic you can update the
   revision to: `4e1f87407d80cdb4a5a293de89d62034fdcbb847`.

   Edit the Image Resource with:

   ```bash
   kp image save tutorial-image --git-revision 4e1f87407d80cdb4a5a293de89d62034fdcbb847
   ```

   You should see kpack schedule a new build by running:

   ```bash
   kp build list tutorial-image
   ```

   You should see a new build with

   ```text
   BUILD    STATUS     IMAGE                                            REASON
   1        SUCCESS    index.docker.io/your-name/app@sha256:6744b...    BUILDPACK
   2        BUILDING                                                    CONFIG
   ```

   You can tail the logs for the Image Resource with the kp cli.

   ```bash
   kp build logs tutorial-image
   ```

   > Note: This second build should be notably faster because the buildpacks can leverage the cache from the previous build.

5. Next steps.

   The next time new buildpacks are added to the store, kpack will automatically
   rebuild the builder. If the updated buildpacks were used by the tutorial
   Image Resource, kpack will automatically create a new build to rebuild your
   OCI image.

### Creating an image using kubectl

1. Create a secret with push credentials for the docker registry that you plan on publishing OCI images to with kpack.

   The easiest way to do that is with `kubectl create secret docker-registry`

    ```bash
    kubectl create secret docker-registry tutorial-registry-credentials \
        --docker-username=user \
        --docker-password=password \
        --docker-server=string
    ```

   > Note: The docker server must be the registry prefix for its corresponding registry. For [dockerhub](https://hub.docker.com/) this should be `https://index.docker.io/v1/`.
   For [GCR](https://cloud.google.com/container-registry/) this should be `gcr.io`. If you use GCR then the username can be `_json_key` and the password can be the JSON credentials you get from the GCP UI (under `IAM -> Service Accounts` create an account or edit an existing one and create a key with type JSON).

   Your secret create should look something like this:

    ```bash
    kubectl create secret docker-registry tutorial-registry-credentials \
        --docker-username=my-dockerhub-username \
        --docker-password=my-dockerhub-password \
        --docker-server=https://index.docker.io/v1/
    ```

   > Note: Learn more about kpack secrets with the [kpack secret documentation](secrets.md)

2. Create a service account that references the registry secret created above

    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: tutorial-service-account
    secrets:
    - name: tutorial-registry-credentials
    imagePullSecrets:
    - name: tutorial-registry-credentials
    ```

   Apply that service account to the cluster

     ```bash
     kubectl apply -f service-account.yaml
     ```

3. Create a kpack Image Resource.

   An Image Resource is the specification for an OCI image that kpack should
   build and manage.

   We will create a sample Image Resource that builds with the builder we created.

   The example included here utilizes
   the [Spring Pet Clinic sample app](https://github.com/spring-projects/spring-petclinic)
   . We encourage you to substitute it with your own application.

   Create an Image Resource:

    ```yaml
    apiVersion: kpack.io/v1alpha2
    kind: Image
    metadata:
      name: tutorial-image
    spec:
      tag: <DOCKER-IMAGE-TAG>
      serviceAccountName: tutorial-service-account
      builder:
        name: my-builder
        kind: Builder
      source:
        git:
          url: https://github.com/spring-projects/spring-petclinic
          revision: 82cb521d636b282340378d80a6307a08e3d4a4c4
    ```

   * Replace `<DOCKER-IMAGE-TAG>` with a valid image tag that exists in the
     registry you configured with the `--docker-server` flag when creating a
     Secret. Something like: `your-name/app`
     or `gcr.io/your-project/app`
   * If you are using your application source, replace `source.git.url`
     & `source.git.revision`.
   > Note: To use a private git repo follow the instructions in [secrets](https://github.com/pivotal/kpack/blob/main/docs/secrets.md)

   Apply that Image Resource to the cluster

    ```bash
    kubectl apply -f image.yaml
    ```

   You can now check the status of the Image Resource.

   ```bash
   kubectl get images
   ```

   You should see that the Image Resource has an unknown READY status as it is
   currently building.

   ```text
    NAME                  LATESTIMAGE   READY
    tutorial-image                      Unknown
    ```

   You can tail the logs for the image that is currently building using
   the [kp cli](https://github.com/vmware-tanzu/kpack-cli/blob/main/docs/kp_build_logs.md)

    ```bash
    kp build logs tutorial-image
    ```

   Once the Image Resource finishes building you can get the fully resolved
   built OCI image with `kubectl get`

    ```bash
    kubectl get image tutorial-image
    ```

   The output should look something like this:

    ```text
    NAMESPACE   NAME                  LATESTIMAGE                                        READY
    default     tutorial-image        index.docker.io/your-project/app@sha256:6744b...   True
    ```

   The latest OCI image is available to be used locally via `docker pull` and in
   a Kubernetes deployment.

4. Run the built application image locally.

   Download the latest OCI image available and run it with Docker.

   ```bash
   docker run -p 8080:8080 <latest-image-with-digest>
   ```

   You should see the java app start up:

   ```text

              |\      _,,,--,,_
             /,`.-'`'   ._  \-;;,_
    _______ __|,4-  ) )_   .;.(__`'-'__     ___ __    _ ___ _______
    |       | '---''(_/._)-'(_\_)   |   |   |   |  |  | |   |       |
    |    _  |    ___|_     _|       |   |   |   |   |_| |   |       | __ _ _
    |   |_| |   |___  |   | |       |   |   |   |       |   |       | \ \ \ \
    |    ___|    ___| |   | |      _|   |___|   |  _    |   |      _|  \ \ \ \
    |   |   |   |___  |   | |     |_|       |   | | |   |   |     |_    ) ) ) )
    |___|   |_______| |___| |_______|_______|___|_|  |__|___|_______|  / / / /
    ==================================================================/_/_/_/

    :: Built with Spring Boot :: 2.2.2.RELEASE
   ```

5. Rebuilding kpack Image Resources.

   We recommend updating the kpack Image Resource with a CI/CD tool when new
   commits are ready to be built.
   > Note: You can also provide a branch or tag as the `spec.git.revision` and kpack will poll and rebuild on updates!

   We can simulate an update from a CI/CD tool by updating
   the `spec.git.revision` on the Image Resource.

   If you are using your own application push an updated commit and use the new
   commit sha. If you are using Spring Pet Clinic you can update the revision
   to: `4e1f87407d80cdb4a5a293de89d62034fdcbb847`.

   Edit the Image Resource with:

   ```bash
   kubectl edit image tutorial-image
   ```

   You should see kpack schedule a new build by running:

   ```bash
   kubectl get builds
   ```

   You should see a new build with

   ```text
   NAME                                IMAGE                                          SUCCEEDED
   tutorial-image-build-1-8mqkc       index.docker.io/your-name/app@sha256:6744b...   True
   tutorial-image-build-2-xsf2l                                                       Unknown
   ```

   You can tail the logs for the image with the kp cli.

   ```bash
   kp build logs tutorial-image -b 2
   ```

   > Note: This second build should be notably faster because the buildpacks can leverage the cache from the previous build.

6. Next steps.

   The next time new buildpacks are added to the store, kpack will
   automatically rebuild the builder. If the updated buildpacks were used by
   the tutorial Image Resource, kpack will automatically create a new build to
   rebuild your OCI image.

## Additional Resources

* [kpack documentation](https://github.com/pivotal/kpack/tree/main/docs)
* [Slack channel](https://kubernetes.slack.com/channels/kpack)
