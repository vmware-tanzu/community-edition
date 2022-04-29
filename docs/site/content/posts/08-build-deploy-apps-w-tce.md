---
title: "Build and Deploy an Application with VMware Tanzu Community Edition"
slug: build-deploy-application-tanzu-community-edition
date: 2022-04-28
author: Cora Iberkleid
image: /img/deploy-app-post-img11.png
excerpt: "We presented a demo that shows how to use Tanzu Community Edition to automate many of the steps involved in cloud native services delivery. You can watch the demo here, or if reading works better than watching for you, then check out the transcript following!"
tags: ['Cora Iberkleid']
---
I recently presented a demo that shows how to use VMware Tanzu Community Edition to automate many of the steps involved in cloud native services delivery.  You can watch the demo here, or if reading works better than watching for you, then check out the transcript following!

<!-- https://gohugo.io/content-management/shortcodes/#youtube -->
{{< youtube id="KdR_QGXdCzw" title="Build and Deploy an Application with VMware Tanzu Community Edition" >}}

## Introduction

In this demo, we’re going to use Tanzu Community Edition to build and deploy an application. We are using Tanzu Community Edition in a few different ways:

1. To provision a Kubernetes cluster
2. To enhance that cluster with additional tooling for building, publishing, and storing container images, and for managing applications at runtime
3. To actually run our application

Let’s get started.

## Cluster provisioning

You can see I actually have several clusters here:

```shell
$ kubectx
kind-tce-demo
tce-mgmt-cluster-admin@tce-mgmt-cluster
tce-workload-cluster-1-admin@tce-workload-cluster-1
```

All of these were provisioned with Tanzu Community Edition. This first one (kind-tce-demo) is an example of an unmanaged cluster. This particular one happens to be running on Docker, locally on my machine, similar to how you might run kind or minikube.

The next one is a management cluster (tce-mgmt-cluster), which is useful for use cases where you need to provision many clusters. I’ve used this management cluster, in fact, to provision the third “workload” cluster that you see here (tce-workload-cluster-1). These happen to be running on AWS.

You can see by the color of the output that my context is pointing to the workload cluster. This is the one I will be using for this demo.

## Documentation

It’s probably worth mentioning that all of the documentation is right here on the Tanzu Community Edition website at [tanzucommunityedition.io/docs](https://tanzucommunityedition.io/docs/). Here you’ll find instructions for provisioning both managed and unmanaged clusters.

Tanzu Community Edition also includes a repository of packages you can add, and this documentation includes installation and configuration for each package. This is really great because one of the challenges in building a more comprehensive platform out of a Kubernetes cluster is all of the research and learning and the choices you need to make about additional tooling you might need. Tanzu Community Edition simplifies this process by providing a curated repository of packages and documentation on dependencies, configuration, installation, and some basic usage.

## Cluster namespaces

Let’s take a look at the namespaces in our workload cluster.

```shell
$ kubens
cert-manager
default
external-dns
harbor
kube-node-lease
kube-public
kube-system
projectcontour
tanzu-package-repo-global
tkg-system
tkg-system-public
```

You can see some namespaces with tkg in the name. This stands for Tanzu Kubernetes Grid and these namespaces contain resources related to core functionality of Tanzu Community Edition, like provisioning clusters, for example.

There is also a namespace called tanzu-package-repo-global—that's where the package repository that I mentioned is installed.

And then, just based on the other namespaces, you can already guess that I've installed a couple of things ahead of time, just in the interest of time. Specifically I installed Harbor already, which is an image registry, and the dependencies it requires, including for example Contour, for ingress.

All of the components that are installed came from the package repository, and I installed all of them following the documentation that I showed you just a minute ago.

## Package repository

So let's look a little deeper into this package repository. We can use the Tanzu CLI to query the cluster directly to see the list of packages that are available.

```shell
$ tanzu package available list
- Retrieving available packages...
  NAME                                           DISPLAY-NAME        SHORT-DESCRIPTION
  cert-manager.community.tanzu.vmware.com        cert-manager        Certificate management
  contour.community.tanzu.vmware.com             contour             An ingress controller
  external-dns.community.tanzu.vmware.com        external-dns        This package provides DNS synchronization functionality.
  fluent-bit.community.tanzu.vmware.com          fluent-bit          Fluent Bit is a fast Log Processor and Forwarder
  gatekeeper.community.tanzu.vmware.com          gatekeeper          policy management
  grafana.community.tanzu.vmware.com             grafana             Visualization and analytics software
  harbor.community.tanzu.vmware.com              harbor              OCI Registry
  knative-serving.community.tanzu.vmware.com     knative-serving     Knative Serving builds on Kubernetes to support deploying and serving of applications and functions as serverless containers
  kpack.community.tanzu.vmware.com               kpack               kpack builds application source code into OCI compliant images using Cloud Native Buildpacks
  local-path-storage.community.tanzu.vmware.com  local-path-storage  This package provides local path node storage and primarily supports RWO AccessMode.
  multus-cni.community.tanzu.vmware.com          multus-cni          This package provides the ability for enabling attaching multiple network interfaces to pods in Kubernetes
  prometheus.community.tanzu.vmware.com          prometheus          A time series database for your metrics
  velero.community.tanzu.vmware.com              velero              Disaster recovery capabilities
  whereabouts.community.tanzu.vmware.com         whereabouts         A CNI IPAM plugin that assigns IP addresses cluster-wide
```

Things that stand out… Cert-Manager, Contour, Harbor, all of which I’ve already installed.

Together we're going to install Knative Serving and kpack.

## Installed packages

In fact, we can check specifically to see which packages have been installed.

```shell
$ tanzu package installed list
/ Retrieving installed packages...
  NAME          PACKAGE-NAME                             PACKAGE-VERSION  STATUS
  cert-manager  cert-manager.community.tanzu.vmware.com  1.6.1            Reconcile succeeded
  contour       contour.community.tanzu.vmware.com       1.19.1           Reconcile succeeded
  external-dns  external-dns.community.tanzu.vmware.com  0.8.0            Reconcile succeeded
  harbor        harbor.community.tanzu.vmware.com        2.3.3            Reconcile succeeded
```

## Create project in Harbor (image registry)

Since Harbor is already installed, we can go ahead and open the UI. We’ll just log in as admin, and let’s go ahead and create a new “demo” project to store the images we’re going to publish.

!["Create a new demo project dashboard"](/img/deploy-app-post-img1.png)

## Install kpack (image build service)

Okay so we have a cluster and we have a registry where we can publish images. The next thing we need is a tool that will turn our source code into images and publish them to the registry. That's exactly what kpack is for, so let's go ahead and install kpack.

Normally you would go to the kpack GitHub website, download the release.yaml file, and use “kubectl apply…”. But we're going to use the kpack package instead. We can use the Tanzu CLI to list the versions of kpack that are included in the repository.

```shell
$ tanzu package available list kpack.community.tanzu.vmware.com
- Retrieving package versions for kpack.community.tanzu.vmware.com...
  NAME                              VERSION  RELEASED-AT
  kpack.community.tanzu.vmware.com  0.5.0
```

I've already checked the documentation I showed you earlier, so I can tell you kpack doesn’t have any prerequisites and there are only a few configuration parameters we need to set in order to  give kpack access to publish to Harbor.

We need to set up authentication for the image registry. You can see that here in the first two lines of this configuration file—there's actually a third line with the password.

```shell
$ head -n 2 cfg/kpack-values.yaml
kp_default_repository: harbor.tanzu.coraiberkleid.site/demo/kp
kp_default_repository_username: admin
```

So we're ready to install kpack. You can see the syntax of the installation command is pretty straightforward, and you would use the same syntax for any of the packages in the package repository.

```shell
$ tanzu package install kpack --package-name kpack.community.tanzu.vmware.com --version 0.5.0 --values-file cfg/kpack-values.yaml
- Installing package 'kpack.community.tanzu.vmware.com'
| Getting namespace 'default'
/ Getting package metadata for 'kpack.community.tanzu.vmware.com'
| Creating service account 'kpack-default-sa'
| Creating cluster admin role 'kpack-default-cluster-role'
| Creating cluster role binding 'kpack-default-cluster-rolebinding'
| Creating secret 'kpack-default-values'
- Creating package resource
| Package install status: Reconciling

 Added installed package 'kpack' in namespace 'default'
```

OK, so kpack has been successfully installed. From here forward, we're just using kpack as we normally would. In other words, Tanzu Community Edition helps us install and upgrade the package, but it doesn't interfere with how we use it.

## Kpack operator perspective: configure ClusterStack

With kpack, you first need to create a builder image, and a builder comprises a stack, meaning the base images for the build step and for runtime, and a store, meaning all of the buildpacks you might need, depending on the kinds of applications you want to handle.

You can create all of these objects using YAML files and the kubectl CLI, but I’m going to do it using the kpack CLI. The kpack CLI also needs access to Harbor, so I need to do a local Docker login.

```shell
$ echo $REGISTRY_PASSWORD | docker login -u admin --password-stdin harbor.tanzu.coraiberkleid.site
Login Succeeded
```

This is the command to create the cluster stack.

```shell
$ kp clusterstack save base --build-image paketobuildpacks/build:base-cnb --run-image paketobuildpacks/run:base-cnb
Creating ClusterStack...
Uploading to 'harbor.tanzu.coraiberkleid.site/demo/kp'...
    Uploading 'harbor.tanzu.coraiberkleid.site/demo/kp@sha256:67c2e829b256aa5be2e9535d358b6f3874e4ac4edb91ad6a256e0e90acfc190a'
    Uploading 'harbor.tanzu.coraiberkleid.site/demo/kp@sha256:41ea15b4d591c2722543009fbf8267a13019ecdc8c6a2b4f437ed83ed29bf72c'
ClusterStack "base" created
```

You can see here that for the build and run base images, I am just borrowing from the open-source Paketo Buildpacks.

One interesting thing to notice is that the kpack CLI is actually doing a few different things. It's copying the two images from Docker Hub into my Harbor registry, and it’s resolving the “base-cnb” tag, which is mutable, to the actual digest. It also creates the ClusterStack resource in Kubernetes, and it uses the address on Harbor with the immutable digest as the tag in the ClusterStack configuration.

So now, not only do I have those images physically closer to where my build is going to take place, so that it's more efficient, but also I have greater governance and control over the precise version of the images that I'm using.

If we check back on Harbor, we can see these two images now on our registry. One of these is the build image and the other is the run image.

!["Build and run images shown on dashboard"](/img/deploy-app-post-img2.png)

## Kpack operator perspective: configure ClusterStore

Similarly, for the store, we can just borrow from existing buildpacks. I'm going to borrow the Java and Go buildpacks, again from Paketo.

```shell
$ kp clusterstore save default -b gcr.io/paketo-buildpacks/java:6.8 -b gcr.io/paketo-buildpacks/go:0.14
Creating ClusterStore...
    Uploading 'harbor.tanzu.coraiberkleid.site/demo/kp@sha256:b5bbfac1e4b534b81709547fde78595e5b3dd1809f11a992a8b10291a5188600'
    Uploading 'harbor.tanzu.coraiberkleid.site/demo/kp@sha256:88ef81ee4e442954c15465cdca232725e963dd8a966fdad973c15d571b5794a7'
ClusterStore "default" created
```

Again, we see that kpack is copying these images into my Harbor registry, and it will create the ClusterStore resource using the Harbor address and the digest rather than the gcr.io address and the mutable tags that I used in the kp command.

We can check Harbor again, and we see two more images here. One has the set of Java buildpacks, and the other has the set of Go buildpacks.

!["Java and Go buildpacks images shown on dashboard"](/img/deploy-app-post-img3.png)

## Kpack operator perspective: configure buildpack detect order

The last step before we can create the builder is to define the order in which we want to evaluate buildpacks during a build. We're going to do this through a file, and we’ll just list Go first.

```shell
$ yq cfg/kpack-builder-order.yaml
- group:
    - id: paketo-buildpacks/go
- group:
    - id: paketo-buildpacks/java
```

## Kpack operator perspective: configure ClusterBuilder

Now we have the building blocks for a builder. So this is the command:

```shell
$ kp clusterbuilder save builder --tag harbor.tanzu.coraiberkleid.site/demo/builder --stack base --store default --order cfg/kpack-builder-order.yaml
ClusterBuilder "builder" created
```

Essentially create a new ClusterBuilder using the stack called base, which is the name of the ClusterStack resource in Kubernetes that we created earlier, and also using the store called default, which is the name of the ClusterStore we created. We can set the tag for the builder image, and we need to reference the file to set the order.

In this case the kpack CLI is only creating the ClusterBuilder resource in Kubernetes. You can see from the output it’s not publishing anything to Harbor. Instead, kpack itself is going to use this configuration and generate a new image for us. That’s happening on the cluster. We can hop over to Harbor and make sure the image is there.

!["Generating new image shown on demo dashboard"](/img/deploy-app-post-img4.png)

!["Generating new image shown on builder dashboard"](/img/deploy-app-post-img5.png)

So that's everything that the kpack operator needs to do. At this point developers can start using the builder to create images for their applications.

## Kpack builder cross-platform compatibility

But before we build an application image, I just want to make the point that this builder can actually be used by any platform that is compliant with cloud-native buildpacks, not just kpack.

If you have developers using the pack CLI, for example, or the Spring Boot Maven or Gradle plugins, they could use this builder and achieve the same builds on their local machine that they can expect to have as part of an automated workflow with kpack.

## Kpack developer perspective

Okay, so on to the developer experience. We first need to create a secret so that kpack can publish images from our working namespace to Harbor on behalf of the developer. This enables you to set different credentials for developers.

```shell
$ kp secret create regcred --registry harbor.tanzu.coraiberkleid.site --registry-user admin # uses $REGISTRY_PASSWORD
Secret "regcred" created
```

And then we can use the kpack CLI to create the image resource in Kubernetes. We specify the builder that we created earlier, the tag where we want the image published, our source code, and the revision.

```shell
$ kp image create hello-world --tag harbor.tanzu.coraiberkleid.site/demo/hello-world --cluster-builder builder --git https://github.com/ciberkleid/go-sample-
app.git --git-revision edd446bc43d2042dfc0045766340ce81a0a4f33f
Creating Image Resource...
Image Resource "hello-world" created
```

The revision can be a branch or a specific commit ID. If it's a branch, kpack will build an image for every commit. So if you only want to build images for commits that have passed testing, you might want your testing job in your pipeline to update the revision in the image with a specific commit ID. Both approaches are valid.

So once we've done that, we can see the list of images on the cluster. The status is unknown, probably because it's still in progress.

```shell
$ kp image list    # enriches `kubectl get images`
NAME           READY      LATEST REASON    LATEST IMAGE    NAMESPACE
hello-world    Unknown    CONFIG                           default
```

We can get a little more information by listing the builds that correspond to this image. There's a one-to-one relationship between your source code and an image resource in Kubernetes, but there's a one-to-many relationship between the image resource and build resources. One build per commit ID.

```shell
$ kp build list    # enriches `kubectl get builds`
BUILD    STATUS      BUILT IMAGE    REASON    IMAGE RESOURCE
1        BUILDING                   CONFIG    hello-world
```

Now the build resource actually spawned a pod and that's where the build is actually taking place. So we can list the pods.

```shell
$ kubectl get pods
NAME                            READY   STATUS     RESTARTS   AGE
hello-world-build-1-build-pod   0/1     Init:4/6   0          24s
```

We can use the kp CLI to get the logs from the pod. If you have used Builpacks before, the output might look familiar to you. It basically shows the lifecycle stages of the build. You can see here it detects which buildpacks to use, and then it runs each buildpack, and then at the end it exports the image to Harbor.

```shell
$ kp build logs hello-world    # gets logs from all containers in pod, in order
===> PREPARE
Build reason(s): CONFIG
CONFIG:
    resources: {}
    - source: {}
    + source:
    +   git:
    +     revision: edd446bc43d2042dfc0045766340ce81a0a4f33f
    +     url: https://github.com/ciberkleid/go-sample-app.git
Loading secret for "harbor.tanzu.coraiberkleid.site" from secret "regcred" at location "/var/build-secrets/regcred"
Cloning "https://github.com/ciberkleid/go-sample-app.git" @ "edd446bc43d2042dfc0045766340ce81a0a4f33f"...
Successfully cloned "https://github.com/ciberkleid/go-sample-app.git" @ "edd446bc43d2042dfc0045766340ce81a0a4f33f" in path "/workspace"
===> ANALYZE
Previous image with name "harbor.tanzu.coraiberkleid.site/demo/hello-world" not found
===> DETECT
4 of 8 buildpacks participating
paketo-buildpacks/ca-certificates 3.0.2
paketo-buildpacks/go-dist         0.8.3
paketo-buildpacks/go-mod-vendor   0.4.0
paketo-buildpacks/go-build        0.7.0
===> RESTORE
===> BUILD

Paketo CA Certificates Buildpack 3.0.2
  https://github.com/paketo-buildpacks/ca-certificates
  Launch Helper: Contributing to layer
    Creating /layers/paketo-buildpacks_ca-certificates/helper/exec.d/ca-certificates-helper
Paketo Go Distribution Buildpack 0.8.3
  Resolving Go version
    Candidate version sources (in priority order):
      go.mod    -> ">= 1.14"
      <unknown> -> ""

    Selected Go version (using go.mod): 1.17.6

  Executing build process
    Installing Go 1.17.6
      Completed in 6.067s

Paketo Go Mod Vendor Buildpack 0.4.0
  Checking module graph
    Running 'go mod graph'
      Completed in 4ms

  Skipping build process: module graph is empty

Paketo Go Build Buildpack 0.7.0
  Executing build process
    Running 'go build -o /layers/paketo-buildpacks_go-build/targets/bin -buildmode pie -trimpath .'
      Completed in 11.186s

  Assigning launch processes:
    hello-server (default): /layers/paketo-buildpacks_go-build/targets/bin/hello-server

===> EXPORT
Adding layer 'paketo-buildpacks/ca-certificates:helper'
Adding layer 'paketo-buildpacks/go-build:targets'
Adding layer 'launch.sbom'
Adding 1/1 app layer(s)
Adding layer 'launcher'
Adding layer 'config'
Adding layer 'process-types'
Adding label 'io.buildpacks.lifecycle.metadata'
Adding label 'io.buildpacks.build.metadata'
Adding label 'io.buildpacks.project.metadata'
Setting default process type 'hello-server'
Saving harbor.tanzu.coraiberkleid.site/demo/hello-world...
*** Images (sha256:32c0ffc57c0c560913e24e6922d1f343a504b67f54ec1156ddcac38031180fd3):
      harbor.tanzu.coraiberkleid.site/demo/hello-world
      harbor.tanzu.coraiberkleid.site/demo/hello-world:b1.20220218.175627
Adding cache layer 'paketo-buildpacks/go-dist:go'
Adding cache layer 'paketo-buildpacks/go-build:gocache'
===> COMPLETION
Build successful
```

And in fact we can check Harbor and validate that the image is there.

```shell
$ kp image list
NAME           READY    LATEST REASON    LATEST IMAGE                                                                                                                NAMESPACE
hello-world    True     CONFIG           harbor.tanzu.coraiberkleid.site/demo/hello-world@sha256:32c0ffc57c0c560913e24e6922d1f343a504b67f54ec1156ddcac38031180fd3    default
```

!["kpack is tagging with a build number and timestamp dashboard"](/img/deploy-app-post-img6.png)

And you can see kpack is also tagging it with a build number and a timestamp.

!["See what buildpacks were used, what run image is dashboard"](/img/deploy-app-post-img7.png)

And it’s good to know that you can inspect this image directly with the pack CLI as well. You can see, for example, which buildpacks were used, what the run image is, etc.

```shell
$ pack inspect harbor.tanzu.coraiberkleid.site/demo/hello-world:latest
Inspecting image: harbor.tanzu.coraiberkleid.site/demo/hello-world:latest

REMOTE:

Stack: io.buildpacks.stacks.bionic

Base Image:
  Reference: harbor.tanzu.coraiberkleid.site/demo/kp@sha256:41ea15b4d591c2722543009fbf8267a13019ecdc8c6a2b4f437ed83ed29bf72c
  Top Layer: sha256:e132d7ea4ce76add13c0d7be8638d170df9322a1e5873608f941693fa6c7dead

Run Images:
  harbor.tanzu.coraiberkleid.site/demo/kp@sha256:41ea15b4d591c2722543009fbf8267a13019ecdc8c6a2b4f437ed83ed29bf72c

Buildpacks:
  ID                                       VERSION        HOMEPAGE
  paketo-buildpacks/ca-certificates        3.0.2          https://github.com/paketo-buildpacks/ca-certificates
  paketo-buildpacks/go-dist                0.8.3          https://github.com/paketo-buildpacks/go-dist
  paketo-buildpacks/go-mod-vendor          0.4.0          https://github.com/paketo-buildpacks/go-mod-vendor
  paketo-buildpacks/go-build               0.7.0          https://github.com/paketo-buildpacks/go-build

Processes:
  TYPE                          SHELL        COMMAND        ARGS
  hello-server (default)                     /layers/paketo-buildpacks_go-build/targets/bin/hello-server

LOCAL:
(not present)
```

This is showing lifecycle level information.

You can actually also extract a bill of materials that has more detailed information provided by individual buildpacks.

## Kpack: rebuild for source code changes

I want to show you rebuilding capabilities. Let's say that your testing pipeline validated some new code, then it could issue the following command to trigger a build for the new commit id.  Since we're doing this at the terminal, we can add this --wait option so that the logs are streamed directly to the console.  We can see a similar output as before showing an image build. It tells us the reason for the build: a new commit id. So we can just wait for that to finish building and exporting to Harbor.

```shell
$ kp image patch hello-world --wait --git-revision fe2c75eac61f52c5613103b025493efefdfbd408
Patching Image Resource...
Image Resource "hello-world" patched
===> PREPARE
Build reason(s): COMMIT
COMMIT:
    - edd446bc43d2042dfc0045766340ce81a0a4f33f
    + fe2c75eac61f52c5613103b025493efefdfbd408
Loading secret for "harbor.tanzu.coraiberkleid.site" from secret "regcred" at location "/var/build-secrets/regcred"
Cloning "https://github.com/ciberkleid/go-sample-app.git" @ "fe2c75eac61f52c5613103b025493efefdfbd408"...
Successfully cloned "https://github.com/ciberkleid/go-sample-app.git" @ "fe2c75eac61f52c5613103b025493efefdfbd408" in path "/workspace"
===> ANALYZE
===> DETECT
4 of 8 buildpacks participating
paketo-buildpacks/ca-certificates 3.0.2
paketo-buildpacks/go-dist         0.8.3
paketo-buildpacks/go-mod-vendor   0.4.0
paketo-buildpacks/go-build        0.7.0
===> RESTORE
Restoring metadata for "paketo-buildpacks/ca-certificates:helper" from app image
Restoring metadata for "paketo-buildpacks/go-dist:go" from cache
Restoring metadata for "paketo-buildpacks/go-build:targets" from app image
Restoring metadata for "paketo-buildpacks/go-build:gocache" from cache
Restoring data for "paketo-buildpacks/go-dist:go" from cache
Restoring data for "paketo-buildpacks/go-build:gocache" from cache
===> BUILD

Paketo CA Certificates Buildpack 3.0.2
  https://github.com/paketo-buildpacks/ca-certificates
  Launch Helper: Reusing cached layer
Paketo Go Distribution Buildpack 0.8.3
  Resolving Go version
    Candidate version sources (in priority order):
      go.mod    -> ">= 1.14"
      <unknown> -> ""

    Selected Go version (using go.mod): 1.17.6

  Reusing cached layer /layers/paketo-buildpacks_go-dist/go

Paketo Go Mod Vendor Buildpack 0.4.0
  Checking module graph
    Running 'go mod graph'
      Completed in 4ms

  Skipping build process: module graph is empty

Paketo Go Build Buildpack 0.7.0
  Executing build process
    Running 'go build -o /layers/paketo-buildpacks_go-build/targets/bin -buildmode pie -trimpath .'
      Completed in 425ms

  Assigning launch processes:
    hello-server (default): /layers/paketo-buildpacks_go-build/targets/bin/hello-server

===> EXPORT
Reusing layers from image 'harbor.tanzu.coraiberkleid.site/demo/hello-world@sha256:32c0ffc57c0c560913e24e6922d1f343a504b67f54ec1156ddcac38031180fd3'
Reusing layer 'paketo-buildpacks/ca-certificates:helper'
Reusing layer 'paketo-buildpacks/go-build:targets'
Reusing 1/1 app layer(s)
Reusing layer 'launcher'
Reusing layer 'config'
Reusing layer 'process-types'
Adding label 'io.buildpacks.lifecycle.metadata'
Adding label 'io.buildpacks.build.metadata'
Adding label 'io.buildpacks.project.metadata'
Setting default process type 'hello-server'
Saving harbor.tanzu.coraiberkleid.site/demo/hello-world...
*** Images (sha256:0313de8f942fa8aaee8cfc7997428a939d4523ab23d49d590aeedc7020b492fb):
      harbor.tanzu.coraiberkleid.site/demo/hello-world
      harbor.tanzu.coraiberkleid.site/demo/hello-world:b2.20220218.180419
Reusing cache layer 'paketo-buildpacks/go-dist:go'
Adding cache layer 'paketo-buildpacks/go-build:gocache'
===> COMPLETION
Build successful
```

You can list the builds here and see the reasons that triggered them, and the address of the image on the registry.

```shell
$ kp build list
BUILD    STATUS     BUILT IMAGE                                                                                                                 REASON    IMAGE RESOURCE
1        SUCCESS    harbor.tanzu.coraiberkleid.site/demo/hello-world@sha256:32c0ffc57c0c560913e24e6922d1f343a504b67f54ec1156ddcac38031180fd3    CONFIG    hello-world
2        SUCCESS    harbor.tanzu.coraiberkleid.site/demo/hello-world@sha256:0313de8f942fa8aaee8cfc7997428a939d4523ab23d49d590aeedc7020b492fb    COMMIT    hello-world
```

!["hello-world on dashboard"](/img/deploy-app-post-img8.png)

## Kpack: rebase for base image changes

We can also look at a rebase, which essentially swaps out operating system layers without touching the rest of the image. You can do this one image at a time using the pack CLI, but kpack will rebase all of your images automatically. So if you think back to those periods of high risk and high stress following the discovery of a vulnerability in an operating system, being able to patch the OS in all of your images within a matter of seconds without touching the app layers is a really powerful capability from a security perspective.

So let's say a vulnerability is discovered and our base images provider, in this case Paketo, releases updated images with a patch. We would update our stack with the patched images.

```shell
$ kp clusterstack save base --build-image paketobuildpacks/build:1.1.41-base-cnb --run-image paketobuildpacks/run:1.1.41-base-cnb
Updating ClusterStack...
Uploading to 'harbor.tanzu.coraiberkleid.site/demo/kp'...
    Uploading 'harbor.tanzu.coraiberkleid.site/demo/kp@sha256:c986d45f6e7be53e056a73e3229d3eaf09fd72a0bba85314d647ed13807f2f9b'
    Uploading 'harbor.tanzu.coraiberkleid.site/demo/kp@sha256:79185c8427ebfed9b7df3e0fa12e101ec8b24aa899bbc541648d5923fb494084'
ClusterStack "base" updated
```

Kpack automatically updates the builder, and then any images that use this builder will kick off a new build. We can see we have a third build.

!["See we have a third build on dashboard"](/img/deploy-app-post-img9.png)

The build reason is “STACK.”

```shell
$ kp build list
BUILD    STATUS     BUILT IMAGE                                                                                                                 REASON    IMAGE RESOURCE
1        SUCCESS    harbor.tanzu.coraiberkleid.site/demo/hello-world@sha256:32c0ffc57c0c560913e24e6922d1f343a504b67f54ec1156ddcac38031180fd3    CONFIG    hello-world
2        SUCCESS    harbor.tanzu.coraiberkleid.site/demo/hello-world@sha256:0313de8f942fa8aaee8cfc7997428a939d4523ab23d49d590aeedc7020b492fb    COMMIT    hello-world
3        SUCCESS    harbor.tanzu.coraiberkleid.site/demo/hello-world@sha256:7779bcfaa4d72a370b6c7af73098e0b68b83bce70de715411e0250c3ff1977b2    STACK     hello-world
```

That alone tells us this is a rebase, but we can validate it by looking at the log.

```shell
$ kp build logs hello-world -b 3
===> REBASE
Build reason(s): STACK
STACK:
    - sha256:41ea15b4d591c2722543009fbf8267a13019ecdc8c6a2b4f437ed83ed29bf72c
    + sha256:79185c8427ebfed9b7df3e0fa12e101ec8b24aa899bbc541648d5923fb494084
Loading secret for "harbor.tanzu.coraiberkleid.site" from secret "regcred" at location "/var/build-secrets/regcred"
*** Images (sha256:7779bcfaa4d72a370b6c7af73098e0b68b83bce70de715411e0250c3ff1977b2):
      harbor.tanzu.coraiberkleid.site/demo/hello-world
      harbor.tanzu.coraiberkleid.site/demo/hello-world:b3.20220218.180539

*** Digest: sha256:7779bcfaa4d72a370b6c7af73098e0b68b83bce70de715411e0250c3ff1977b2
===> COMPLETION
Build successful
```

You can see that the log file is actually quite different from the build log. It's also very fast.

## Install Knative Serving

So I think our image is ready for deployment. We could just create Deployment and Service yaml files, but since we’re working with Tanzu Community Edition, it’s easy to up our game.

Let’s install Knative Serving so that we can achieve a more sophisticated deployment for our application. Again, we can check for available versions.

```shell
$ tanzu package available list knative-serving.community.tanzu.vmware.com
\ Retrieving package versions for knative-serving.community.tanzu.vmware.com...
  NAME                                        VERSION  RELEASED-AT
  knative-serving.community.tanzu.vmware.com  0.22.0
  knative-serving.community.tanzu.vmware.com  0.26.0
  knative-serving.community.tanzu.vmware.com  1.0.0
```

I’ve already checked the docs and made sure we have the necessary prerequisites and configuration values, which basically just provide my AWS DNS address for routing. So we can use the same syntax we used for kpack to install Knative Serving.

```shell
$ tanzu package install knative-serving --package-name knative-serving.community.tanzu.vmware.com --version 1.0.0 -f
cfg/knative-serving-values.yaml
- Installing package 'knative-serving.community.tanzu.vmware.com'
| Getting namespace 'default'
/ Getting package metadata for 'knative-serving.community.tanzu.vmware.com'
| Creating service account 'knative-serving-default-sa'
| Creating cluster admin role 'knative-serving-default-cluster-role'
| Creating cluster role binding 'knative-serving-default-cluster-rolebinding'
| Creating secret 'knative-serving-default-values'
- Creating package resource
\ Package install status: Reconciling

 Added installed package 'knative-serving' in namespace 'default'
```

OK, looks like it completed installation. We can again list all of the packages we’ve installed so far. We see kpack and Knative Serving there, in addition to what I had already installed.

```shell
$ tanzu package installed list
- Retrieving installed packages...
  NAME             PACKAGE-NAME                                PACKAGE-VERSION  STATUS
  cert-manager     cert-manager.community.tanzu.vmware.com     1.6.1            Reconcile succeeded
  contour          contour.community.tanzu.vmware.com          1.19.1           Reconcile succeeded
  external-dns     external-dns.community.tanzu.vmware.com     0.8.0            Reconcile succeeded
  harbor           harbor.community.tanzu.vmware.com           2.3.3            Reconcile succeeded
  knative-serving  knative-serving.community.tanzu.vmware.com  1.0.0            Reconcile succeeded
  kpack            kpack.community.tanzu.vmware.com            0.5.0            Reconcile succeeded
```

OK, back to the developer perspective.

## Knative Serving developer perspective

Let’s create a namespace to deploy our application. Then all we need to do is create a Knative type Service. You can do this with a YAML file but I’m going to opt for the Knative CLI. You can see it’s a simple command.

```shell
$ kubectl create ns apps
namespace/apps created

$ kn service create hello-world-app --image harbor.tanzu.coraiberkleid.site/demo/hello-world -n apps
Creating service 'hello-world-app' in namespace 'apps':

  0.044s The Route is still working to reflect the latest desired specification.
  0.084s ...
  0.106s Configuration "hello-world-app" is waiting for a Revision to become ready.
  5.980s ...
  5.980s Ingress has not yet been reconciled.
  5.980s Waiting for Envoys to receive Endpoints data.
  6.177s Waiting for load balancer to be ready
  6.396s Ready to serve.

Service 'hello-world-app' created to latest revision 'hello-world-app-00001' is available at URL:
http://hello-world-app.apps.tanzu.coraiberkleid.site
```

We can make sure the app is working.

!["hello sunshine example shown"](/img/deploy-app-post-img10.png)

So that’s really it, but in case you’re curious about Knative Serving, let me show you what it generated in terms of Kubernetes resources.

```shell
$ kubectl get all -n apps
NAME                                                    READY   STATUS    RESTARTS   AGE
pod/hello-world-app-00001-deployment-67567f89b8-lwfsj   2/2     Running   0          12s

NAME                                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                                      AGE
service/hello-world-app                 ClusterIP   None            <none>        80/TCP                                       11s
service/hello-world-app-00001           ClusterIP   100.67.142.96   <none>        80/TCP                                       12s
service/hello-world-app-00001-private   ClusterIP   100.67.2.12     <none>        80/TCP,9090/TCP,9091/TCP,8022/TCP,8012/TCP   12s

NAME                                               READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello-world-app-00001-deployment   1/1     1            1           12s

NAME                                                          DESIRED   CURRENT   READY   AGE
replicaset.apps/hello-world-app-00001-deployment-67567f89b8   1         1         1       12s

NAME                                        URL                                                    READY   REASON
route.serving.knative.dev/hello-world-app   http://hello-world-app.apps.tanzu.coraiberkleid.site   True

NAME                                                LATESTCREATED           LATESTREADY             READY   REASON
configuration.serving.knative.dev/hello-world-app   hello-world-app-00001   hello-world-app-00001   True

NAME                                          URL                                                    LATESTCREATED           LATESTREADY             READY   REASON
service.serving.knative.dev/hello-world-app   http://hello-world-app.apps.tanzu.coraiberkleid.site   hello-world-app-00001   hello-world-app-00001   True

NAME                                                 CONFIG NAME       K8S SERVICE NAME   GENERATION   READY   REASON   ACTUAL REPLICAS   DESIRED REPLICAS
revision.serving.knative.dev/hello-world-app-00001   hello-world-app                      1            True             1                 1
```

Of course, the basic resources we need—Deployment and Service—are still there, and the Deployment creates a ReplicaSet and Pod as expected, but there are other resources that the Knative Service created.

Of course it configured the routing, as we verified by checking the app.

It also creates revisions, to manage versions of the application. Knative is able to run multiple versions simultaneously and split traffic between them. You can also use them to easily roll back.
Knative also scales your app to zero instances if it is idle for a while, and scales back up as needed.

So it really is a more sophisticated way of deploying an application as compared to a simple Deployment and Service, and with Tanzu Community Edition it’s so easy to install Knative, there’s very little reason not to take advantage of it.

Overall, you can see that it’s powerful to work with a project that can make it so much easier to provision Kubernetes clusters and build out a better platform on them.

## Join the Tanzu Community Edition Community

We are excited to hear from you and learn with you! Here are several ways you can get involved:

* Join Tanzu Community Edition's slack channel, [#tanzu-community-edition](https://kubernetes.slack.com/archives/C02GY94A8KT) Kubernetes workspace, and connect with maintainers and other Tanzu Community Edition users.
* Find us on [GitHub](https://github.com/vmware-tanzu/community-edition/). Suggest how we can improve the project, the docs, or share any other feedback.
* Attend our Community Meetings, with two options to choose from. Check out the [Community page](https://tanzucommunityedition.io/community/) for full details on how to attend.
