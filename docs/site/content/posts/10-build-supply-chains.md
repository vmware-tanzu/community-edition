---
title: "Build Modern Software Supply Chains with Cartographer in VMware Tanzu Community Edition"
slug: build-modern-software-supply-chains-with-cartographer
date: 2022-06-09
author: Cora Iberkleid
image: /img/Build-blog-img0.jpg
excerpt: "Tanzu Community Edition offers Cartographer, advanced software supply chain tooling that can help you deliver applications more rapidly, securely, and efficiently at scale. In the embedded video you can see how you can use Cartographer to build and maintain paths to production from reusable building blocks and create pre-approved workflows that enable dev teams to focus on writing code and boost their productivity. A transcript of the demo video is provided in this post."
tags: ['Cora Iberkleid']
---
Tanzu Community Edition offers [Cartographer](https://cartographer.sh/), advanced software supply chain tooling that can help you deliver applications more rapidly, securely, and efficiently at scale. In the video embedded below you can see how you can use Cartographer to build and maintain paths to production from reusable building blocks and create pre-approved workflows that enable dev teams to focus on writing code and boost their productivity. A transcript of the demo video follows.

<!-- https://gohugo.io/content-management/shortcodes/#youtube -->
{{< youtube id="URgg2s4OClc" title="Create a Software Supply Chain with VMware Tanzu Community Edition" >}}

## Introduction - Packages

I'm going to take a few moments to show you Cartographer in action on VMware Tanzu Community Edition.

One of the great things about Tanzu Community Edition is that it provides you a curated repository of packages that you can use to enhance your cluster to make it more valuable and easier to use for both developers and operators.

You can see here the list of packages available in the repository.

```shell
$ tanzu package available list -o yaml | yq '.[].display-name'
App-Toolkit package for TCE
Cartographer
cert-injection-webhook
cert-manager
contour
external-dns
fluent-bit
Flux Source Controller
gatekeeper
grafana
harbor
knative-serving
kpack
local-path-storage
multus-cni
prometheus
velero
whereabouts
```

I have already installed some packages.

```shell
$ tanzu package installed list
- Retrieving installed packages...
  NAME                      PACKAGE-NAME                                         PACKAGE-VERSION  STATUS
  cartographer              cartographer.community.tanzu.vmware.com              0.2.2            Reconcile succeeded
  cert-manager              cert-manager.community.tanzu.vmware.com              1.6.1            Reconcile succeeded
  contour                   contour.community.tanzu.vmware.com                   1.20.1           Reconcile succeeded
  external-dns              external-dns.community.tanzu.vmware.com              0.10.0           Reconcile succeeded
  fluxcd-source-controller  fluxcd-source-controller.community.tanzu.vmware.com  0.21.2           Reconcile succeeded
  harbor                    harbor.community.tanzu.vmware.com                    2.3.3            Reconcile succeeded
  kpack                     kpack.community.tanzu.vmware.com                     0.5.1            Reconcile
```

Some of these are for managing ingress, DNS, and certificates automatically, and some of these are for defining our path to production, including, of course, Cartographer, which is the highlight of this demo.

## Workflow

In order to understand how all of these tools can be used together, let's start with a basic workflow.

```shell
# BASIC WORKFLOW: source (fluxcd) -> image (kpack) -> running app (Deployment, Service, HTTPProxy)

# Cartographer: Automation through Choreography of Kubernetes resources
```

We will use [Flux](https://fluxcd.io/) to poll source code repositories for new commits; [kpack](https://github.com/pivotal/kpack) to build an image; and then we’ll deploy the image using basic Kubernetes resources: a Deployment, Service, and HTTPProxy, nothing fancy there. It’s not shown in the workflow but kpack will be publishing the image to [Harbor](https://goharbor.io/), which is a container registry that was also installed in the cluster using the package repository.

The beauty of creating a path to production in this way is that you can leverage mature and powerful ecosystem tools, but that's not enough to create a path to production. That's where Cartographer comes in—Cartographer enables you to integrate these tools and resources into a meaningful workflow, and to make it available to many different applications.

## Workload

Let's start with the developer perspective.

```shell
$ # Developer Perspective: Workload

$ yq developer/workload-1.yaml
apiVersion: carto.run/v1alpha1
kind: Workload
metadata:
  name: hello-sunshine
  labels:
    app.tanzu.vmware.com/workload-type: web
spec:
  serviceAccountName: cartographer-example-workload-sa
  source:
    git:
      url: https://github.com/ciberkleid/hello-go
      ref:
        branch: main
```

The developer is responsible for providing information unique to their application, or their workload. That could be as simple as the git repo url and the branch. Of course, security is always top-of-mind, so it makes sense to include a service account for role based access control. And finally, in this example, you can see a label here. This label allows the Workload to express which Path to Production it should match with, because you may have more than one path to handle different types of applications.

This Workload abstraction provides a clear separation of concerns between the developer, providing their application details, and the application operator behind the scenes, creating the path that this application will follow in order to get securely to production.

As a developer, I'll go ahead and apply this Workload.

```shell
$ kubectl apply -f developer/workload-1.yaml
workload.carto.run/hello-sunshine created
```

And then it makes sense, since I've just deployed a Workload type resource, to check its status.

```shell
$ kubectl get workload hello-sunshine
NAME             SOURCE                                   SUPPLYCHAIN    READY     REASON
hello-sunshine   https://github.com/ciberkleid/hello-go   supply-chain   Unknown   MissingValueAtPath
```

We can see it has been matched with a Supply Chain called “supply-chain,” we'll get to that in a minute, and we can also see that it’s not quite ready yet.

We can get more detail by getting the status as YAML.

```shell
$ kubectl get workload hello-sunshine -o yaml | yq
apiVersion: carto.run/v1alpha1
kind: Workload
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"carto.run/v1alpha1","kind":"Workload","metadata":{"annotations":{},"labels":{"app.tanzu.vmware.com/workload-type":"web"},"name":"hello-sunshine","namespace":"default"},"spec":{"serviceAccountName":"cartographer-example-workload-sa","source":{"git":{"ref":{"branch":"main"},"url":"https://github.com/ciberkleid/hello-go"}}}}
  creationTimestamp: "2022-04-13T21:28:04Z"
  generation: 1
  labels:
    app.tanzu.vmware.com/workload-type: web
  name: hello-sunshine
  namespace: default
  resourceVersion: "6559536"
  uid: 24a506a3-2fba-4b0d-ac07-90972e0892f0
spec:
  serviceAccountName: cartographer-example-workload-sa
  source:
    git:
      ref:
        branch: main
      url: https://github.com/ciberkleid/hello-go
status:
  conditions:
    - lastTransitionTime: "2022-04-13T21:28:04Z"
      message: ""
      reason: Ready
      status: "True"
      type: SupplyChainReady
    - lastTransitionTime: "2022-04-13T21:28:07Z"
      message: waiting to read value [.status.latestImage] from resource [image.kpack.io/hello-sunshine] in namespace [default]
      reason: MissingValueAtPath
      status: Unknown
      type: ResourcesSubmitted
    - lastTransitionTime: "2022-04-13T21:28:07Z"
      message: waiting to read value [.status.latestImage] from resource [image.kpack.io/hello-sunshine] in namespace [default]
      reason: MissingValueAtPath
      status: Unknown
      type: Ready
  observedGeneration: 1
  supplyChainRef:
    kind: ClusterSupplyChain
    name: supply-chain
```

And it looks like Cartographer is waiting for kpack to finish building the image so that it can retrieve the latest tag.

We’ll just give kpack a few moments to finish building the image. If you want to learn more about how kpack builds images, definitely check out [the webinar](https://tanzu.vmware.com/content/webinars/feb-22-automate-container-builds-with-vmware-tanzu-community-edition) we did last month which showcased kpack.

As I mentioned earlier, I’ve configured kpack to publish the image to Harbor. Harbor will scan the image for vulnerabilities using a scanner called Trivy by default, and it can also sign images.

Let’s check the status again.

```shell
$ kubectl get workload hello-sunshine -o yaml | yq
apiVersion: carto.run/v1alpha1
kind: Workload
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"carto.run/v1alpha1","kind":"Workload","metadata":{"annotations":{},"labels":{"app.tanzu.vmware.com/workload-type":"web"},"name":"hello-sunshine","namespace":"default"},"spec":{"serviceAccountName":"cartographer-example-workload-sa","source":{"git":{"ref":{"branch":"main"},"url":"https://github.com/ciberkleid/hello-go"}}}}
  creationTimestamp: "2022-04-13T21:28:04Z"
  generation: 1
  labels:
    app.tanzu.vmware.com/workload-type: web
  name: hello-sunshine
  namespace: default
  resourceVersion: "6559880"
  uid: 24a506a3-2fba-4b0d-ac07-90972e0892f0
spec:
  serviceAccountName: cartographer-example-workload-sa
  source:
    git:
      ref:
        branch: main
      url: https://github.com/ciberkleid/hello-go
status:
  conditions:
    - lastTransitionTime: "2022-04-13T21:28:04Z"
      message: ""
      reason: Ready
      status: "True"
      type: SupplyChainReady
    - lastTransitionTime: "2022-04-13T21:28:58Z"
      message: ""
      reason: ResourceSubmissionComplete
      status: "True"
      type: ResourcesSubmitted
    - lastTransitionTime: "2022-04-13T21:28:58Z"
      message: ""
      reason: Ready
      status: "True"
      type: Ready
  observedGeneration: 1
  supplyChainRef:
    kind: ClusterSupplyChain
    name: supply-chain
```

Great - it looks like it’s done, meaning Cartographer has also already deployed the application.

Let's check on that.

```shell
$ kubectl get all,httpproxies
NAME                                   READY   STATUS      RESTARTS   AGE
pod/hello-sunshine-6965c68dcc-72ntq    1/1     Running     0          112s
pod/hello-sunshine-build-1-build-pod   0/1     Completed   0          2m43s

NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/hello-sunshine   ClusterIP   100.65.215.75   <none>        80/TCP    112s
service/kubernetes       ClusterIP   100.64.0.1      <none>        443/TCP   5d5h

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello-sunshine   1/1     1            1           113s

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/hello-sunshine-6965c68dcc   1         1         1       113s

NAME                                SOURCE                                   SUPPLYCHAIN    READY   REASON
workload.carto.run/hello-sunshine   https://github.com/ciberkleid/hello-go   supply-chain   True    Ready

NAME                                         FQDN                                      TLS SECRET                        STATUS   STATUS DESCRIPTION
httpproxy.projectcontour.io/hello-sunshine   hello-sunshine.tanzu.coraiberkleid.site   developer-certificates/wildcard   valid    Valid HTTPProxy
```

And, in fact, we can see the Deployment, the Service, and the HTTPProxy that we expected.

We can also just make sure the app is working.

!["Hello Sunshine example"](/img/Build-blog-img1.png)

Everything looks good.

Now, let's look at what actually happened behind the scenes and understand, as an operator, how you can put this workflow together.

## Putting it all together

Cartographer gives an operator a few different resources to work with. An operator creating a path to production would be interested in chaining together a sequence of activities. With Cartographer, this would be done using a ClusterSupplyChain and a set of Templates. You can see there are a few different types of templates to choose from. Let's take a look at what we used in the example we just ran.

```shell
$ kubectl api-resources --api-group carto.run
NAME                         SHORTNAMES   APIVERSION           NAMESPACED   KIND
clusterconfigtemplates                    carto.run/v1alpha1   false        ClusterConfigTemplate
clusterdeliveries                         carto.run/v1alpha1   false        ClusterDelivery
clusterdeploymenttemplates                carto.run/v1alpha1   false        ClusterDeploymentTemplate
clusterimagetemplates                     carto.run/v1alpha1   false        ClusterImageTemplate
clusterruntemplates                       carto.run/v1alpha1   false        ClusterRunTemplate
clustersourcetemplates                    carto.run/v1alpha1   false        ClusterSourceTemplate
clustersupplychains                       carto.run/v1alpha1   false        ClusterSupplyChain
clustertemplates                          carto.run/v1alpha1   false        ClusterTemplate
deliverables                              carto.run/v1alpha1   true         Deliverable
runnables                                 carto.run/v1alpha1   true         Runnable
workloads                                 carto.run/v1alpha1   true         Workload
```

Here's the Supply Chain that handled our Workload.

```shell
$ yq app-operator/supply-chain.yaml
apiVersion: carto.run/v1alpha1
kind: ClusterSupplyChain
metadata:
  name: supply-chain
spec:
  selector:
    app.tanzu.vmware.com/workload-type: web
  resources:
    - name: source-provider
      templateRef:
        kind: ClusterSourceTemplate
        name: source
    - name: image-builder
      templateRef:
        kind: ClusterImageTemplate
        name: image
      params:
        - name: image_prefix
          value: harbor.tanzu.coraiberkleid.site/demo/
      sources:
        - resource: source-provider
          name: source
    - name: deployer
      templateRef:
        kind: ClusterTemplate
        name: app-deploy
      images:
        - resource: image-builder
          name: image
```

You can see the selector value is "web," which matches the label that we used in the Workload. The name of the Supply Chain is simply supply-chain—which we saw in the status information for the Workload. And you can see a list of three resources that use three of the available templates: ClusterSourceTemplate, ClusterImageTemplate, and a generic ClusterTemplate. As you can imagine, these correspond to Flux, which provides the source, kpack, which produces the image, and the three deployment resources. You can infer that the choice of template is related to the expected output of that particular activity.

The order of the resources is also significant, which means that the Supply Chain enables an operator to express a sequence of activities. You wouldn't want to create a Deployment resource before you had an image available, for example.

And you can see a mapping of outputs to inputs here, which also enables an operator to chain resources together. The image-provider, for example, specifies that it needs input from the source-provider. In other words, kpack needs input from Flux. And similarly, the third resource is dependent on the output of our image-provider, meaning the Deployment needs the output from kpack. Cartographer will be responsible for retrieving outputs and providing them as inputs to the next resource in real-time.

## Templates

So let's dig into these three templates.

```shell
$ kubectl get clustersourcetemplate,clusterimagetemplate,clustertemplate
NAME                                     AGE
clustersourcetemplate.carto.run/source   50m

NAME                                   AGE
clusterimagetemplate.carto.run/image   50m

NAME                                   AGE
clustertemplate.carto.run/app-deploy   50m
```

You can see they exist in the cluster, and you can see that their kind and name match what was expressed in the supply chain. Let's look at the first two so that we can understand the pattern.

The first template is a ClusterSourceTemplate.

```shell
$ kubectl get clustersourcetemplate source -o yaml | yq 'del(.metadata)'
apiVersion: carto.run/v1alpha1
kind: ClusterSourceTemplate
spec:
  revisionPath: .status.artifact.revision
  template:
    apiVersion: source.toolkit.fluxcd.io/v1beta1
    kind: GitRepository
    metadata:
      name: $(workload.metadata.name)$
    spec:
      gitImplementation: libgit2
      ignore: ""
      interval: 1m0s
      ref: $(workload.spec.source.git.ref)$
      url: $(workload.spec.source.git.url)$
  urlPath: .status.artifact.url
```

You can see that it has a "template" section in the spec, and the value there is literally the configuration for a Flux GitRepository resource. This template field can be used for any arbitrary resource configuration. Cartographer will simply submit it to the Kubernetes API.

Of course we want to use this template for many different applications, so instead of hard-coding any workload values here, we have placeholders that map to the structure of the Workload we saw earlier. So for every Workload, Cartographer will instantiate a new Flux GitRepository resource. If you wanted to use a different resource, other than Flux, you could do that as well, just by replacing this with the YAML configuration for the resource of your choice.

The other thing to note here are the output paths. For a ClusterSourceTemplate, the outputs are a url and a revision. These values will be different for each source code commit, so we can’t know them ahead of time. Instead, what we do know is that the Flux GitRepository resource will store these values in its status, so we can provide the path to these values. This enables Cartographer to monitor the resource and pull out the correct values.

To summarize, this template gives Cartographer the ability to customize the configuration of the Flux GitRepository resource with Workload-specific details, it gives Cartographer the ability to instantiate new GitRepository resources by submitting the configuration to the Kubernetes API server, and it gives Cartographer the ability to continually monitor the resource and extract the right output values every time they change.

So let's see the instance that was created when we submitted the Workload. I'm filtering out the metadata for readability, but the name is "hello-sunshine", and you can see placeholders have been replaced with the values from the Workload.

```shell
$ kubectl get gitrepository hello-sunshine -o yaml | yq 'del(.metadata)'
apiVersion: source.toolkit.fluxcd.io/v1beta1
kind: GitRepository
spec:
  gitImplementation: libgit2
  ignore: ""
  interval: 1m0s
  ref:
    branch: main
  timeout: 60s
  url: https://github.com/ciberkleid/hello-go
status:
  artifact:
    checksum: 8adaa9bf76c8bdc0db2a8bd777572fa3d1f78da151b08cb70add2bf1977393f4
    lastUpdateTime: "2022-04-13T21:28:06Z"
    path: gitrepository/default/hello-sunshine/2afe4943221ad9b1bc6c2ac9e16b4a5c1e5c0a5d.tar.gz
    revision: main/2afe4943221ad9b1bc6c2ac9e16b4a5c1e5c0a5d
    url: http://source-controller.flux-system.svc.cluster.local./gitrepository/default/hello-sunshine/2afe4943221ad9b1bc6c2ac9e16b4a5c1e5c0a5d.tar.gz
  conditions:
    - lastTransitionTime: "2022-04-13T21:28:06Z"
      message: 'Fetched revision: main/2afe4943221ad9b1bc6c2ac9e16b4a5c1e5c0a5d'
      reason: GitOperationSucceed
      status: "True"
      type: Ready
  observedGeneration: 1
  url: http://source-controller.flux-system.svc.cluster.local./gitrepository/default/hello-sunshine/latest.tar.gz
```

You can also see that the status, which is where Kubernetes maintains information about the current state of resources, contains the paths that we used to specify where the output would be. For url, we specified the path .status.artifact.url, and in fact you can see under status→artifact→url, there is the value that we need to pass to kpack. This is how we can teach Cartographer to work with Flux, or any other resource. So you see how Cartographer can be used with arbitrary resources—it doesn't require special plugins or anything like that.

Let's take a look at the second resource, for kpack. The approach is the same, except in this case we are using a ClusterImageTemplate.

```shell
$ kubectl get clusterimagetemplate image -o yaml | yq 'del(.metadata)'
apiVersion: carto.run/v1alpha1
kind: ClusterImageTemplate
spec:
  imagePath: .status.latestImage
  params:
    - default: some-default-prefix-
      name: image_prefix
  template:
    apiVersion: kpack.io/v1alpha2
    kind: Image
    metadata:
      name: $(workload.metadata.name)$
    spec:
      build:
        env: $(workload.spec.build.env)$
      builder:
        kind: ClusterBuilder
        name: builder
      serviceAccountName: cartographer-example-registry-creds-sa
      source:
        blob:
          url: $(sources.source.url)$
      tag: $(params.image_prefix)$$(workload.metadata.name)$
```

In the template block, we provide the configuration for a kpack image. In this case, you can see Cartographer will inject the name from the Workload, the URL will come from the source that we linked in the Supply Chain, which is of course the Flux GitRepository, and it will also use something called params, which is a way of specifying global or default parameters across Workloads. In this case, we are using params to set the address of the Harbor registry.

You can also see this template has an output field called image. In this case, we provide the path to the status location where kpack places the new tag information.

Let's look at the kpack image that was created for our hello-sunshine Workload.

```shell
$ kubectl get cnbimage hello-sunshine -o yaml | yq 'del(.metadata)'
apiVersion: kpack.io/v1alpha2
kind: Image
spec:
  build:
    resources: {}
  builder:
    kind: ClusterBuilder
    name: builder
  cache:
    volume:
      size: 2G
  failedBuildHistoryLimit: 10
  imageTaggingStrategy: BuildNumber
  serviceAccountName: cartographer-example-registry-creds-sa
  source:
    blob:
      url: http://source-controller.flux-system.svc.cluster.local./gitrepository/default/hello-sunshine/2afe4943221ad9b1bc6c2ac9e16b4a5c1e5c0a5d.tar.gz
  successBuildHistoryLimit: 10
  tag: harbor.tanzu.coraiberkleid.site/demo/hello-sunshine
status:
  buildCacheName: hello-sunshine-cache
  buildCounter: 1
  conditions:
    - lastTransitionTime: "2022-04-13T21:28:56Z"
      status: "True"
      type: Ready
    - lastTransitionTime: "2022-04-13T21:28:56Z"
      status: "True"
      type: BuilderReady
  latestBuildImageGeneration: 1
  latestBuildReason: CONFIG
  latestBuildRef: hello-sunshine-build-1
  latestImage: harbor.tanzu.coraiberkleid.site/demo/hello-sunshine@sha256:04e90c1a785557f5f6fdb90d68c3b3626b5a3aff3d4adde8012b58bbb8295b26
  latestStack: io.buildpacks.stacks.bionic
  observedGeneration: 1
```

You can see the values that Cartographer injected in the spec, and you can see that the imagePath that we used in the template matches the location of the new tag in the status, .status.latestImage. So, any time Flux finds a new git commit, Cartographer will update the source URL here and resubmit this kpack image YAML to Kubernetes. That will trigger kpack to build a new image.

In the interest of time, I am not going to show the third template, but it follows the same pattern.

## Alternative Workflow

What if you decide you want to make an alternative workflow? For example, you want to deploy some of your applications using [Knative Serving](https://knative.dev/docs/serving/), which offers some useful features like scaling to zero and managing revisions, with very simple configuration. What would it take to create this alternative path?

```shell
# ALTERNATIVE WORKFLOW: source (fluxcd) -> image (kpack) -> running app (knative-serving)
```

First of all, let's use Tanzu Community Edition to install Knative Serving, since it is included as a package in the repository. We can check what versions are available.

```shell
$ tanzu package available list knative-serving.community.tanzu.vmware.com
\ Retrieving package versions for knative-serving.community.tanzu.vmware.com...
  NAME                                        VERSION  RELEASED-AT
  knative-serving.community.tanzu.vmware.com  0.22.0   0001-01-01 00:00:00 +0000 UTC
  knative-serving.community.tanzu.vmware.com  0.26.0   0001-01-01 00:00:00 +0000 UTC
  knative-serving.community.tanzu.vmware.com  1.0.0    0001-01-01 00:00:00 +0000 UTC
```

And we'll go ahead and install the latest version.

```shell
$ tanzu package install knative-serving --package-name knative-serving.community.tanzu.vmware.com --version 1.0.0 -f knative-values.yaml
| 'PackageInstall' resource install status: Reconciling
```

You can see that Tanzu Community Edition makes it very easy to add packages to a cluster.

Now that we've got Knative installed, we need a new template for the application deployment. I have that configured already. You can see the Knative Service configuration here, which is very simple, especially considering all of the features it will provide.

```shell
$ yq app-operator/app-deploy-template-kn.yaml
apiVersion: carto.run/v1alpha1
kind: ClusterTemplate
metadata:
  name: app-deploy-kn
spec:
  template:
    apiVersion: kappctrl.k14s.io/v1alpha1
    kind: App
    metadata:
      name: $(workload.metadata.name)$
    spec:
      serviceAccountName: cartographer-example-registry-creds-sa
      fetch:
        - inline:
            paths:
              manifest.yml: |
                ---
                apiVersion: kapp.k14s.io/v1alpha1
                kind: Config
                rebaseRules:
                  - path:
                      - metadata
                      - annotations
                      - serving.knative.dev/creator
                    type: copy
                    sources: [new, existing]
                    resourceMatchers: &matchers
                      - apiVersionKindMatcher:
                          apiVersion: serving.knative.dev/v1
                          kind: Service
                  - path:
                      - metadata
apiVersion: carto.run/v1alpha1
                      - annotations
                      - serving.knative.dev/lastModifier
                    type: copy
                    sources: [new, existing]
                    resourceMatchers: *matchers

                ---
                apiVersion: serving.knative.dev/v1
                kind: Service
                metadata:
                  name: $(workload.metadata.name)$
                spec:
                  template:
                    metadata:
                      annotations:
                        autoscaling.knative.dev/minScale: "1"
                    spec:
                      serviceAccountName: cartographer-example-registry-creds-sa
                      containers:
                        - name: workload
                          image: $(images.image.image)$
                          env: $(workload.spec.env)$
                          securityContext:
                            runAsUser: 1000
      template:
        - ytt: {}
      deploy:
        - kapp: {}
```

There are actually two other details in this template—one is wrapping the [Knative](https://knative.dev/docs/) configuration with a [kapp-controller](https://carvel.dev/kapp-controller/) resource. Kapp-controller is part of the [Carvel tool suite](https://carvel.dev/), and it provides some extra functionality when submitting resources to Kubernetes. That's a separate topic, but I would encourage you to check it out.

The second detail is a hook to allow developers to provide runtime environment variables. So the information in the Workload can be very simple, but it can also be enriched to express more configuration. The point is, you have a lot of freedom and control in how you choose to create these templates.

So next we need a new Supply Chain. Let’s copy the original Supply Chain. After all, we want to reuse the same Flux and kpack templates.

```shell
cp app-operator/supply-chain.yaml app-operator/supply-chain-kn.yaml
```

We'll just edit the name of the Supply Chain, the selector—so that Workloads can choose which chain to use—and we'll make sure this chain uses our new Knative app-deploy template.

```shell
$ vi app-operator/supply-chain.yaml
apiVersion: carto.run/v1alpha1
kind: ClusterSupplyChain
metadata:
  name: supply-chain-kn
spec:
  selector:
    app.tanzu.vmware.com/workload-type: web-kn
  resources:
    - name: source-provider
      templateRef:
        kind: ClusterSourceTemplate
        name: source
    - name: image-builder
      templateRef:
        kind: ClusterImageTemplate
        name: image
      params:
        - name: image_prefix
          value: harbor.tanzu.coraiberkleid.site/demo/
      sources:
        - resource: source-provider
          name: source
    - name: deployer
      templateRef:
        kind: ClusterTemplate
        name: app-deploy-kn
      images:
        - resource: image-builder
          name: image
```

Let's apply the new template and the new Supply Chain to the cluster.

```shell
$ kubectl apply -f app-operator/app-deploy-template-kn.yaml -f app-operator/supply-chain-kn.yaml
clustertemplate.carto.run/app-deploy-kn created
clustersupplychain.carto.run/supply-chain-kn created
```

So now it's ready for a developer to use.

Let's look at a second Workload.

```shell
$ yq developer/workload-2.yaml
apiVersion: carto.run/v1alpha1
kind: Workload
metadata:
  name: hello-tanzu
  labels:
    app.tanzu.vmware.com/workload-type: web-kn
spec:
  serviceAccountName: cartographer-example-workload-sa
  source:
    git:
      url: https://github.com/ciberkleid/hello-go
      ref:
        branch: main
  build:
    env:
      - name: CGO_ENABLED
        value: "0"
  env:
    - name: HELLO_MSG
      value: "tanzu"
```

I'm cheating a little by using the same git URL but let's pretend it is a different app. We're giving the application a different name, and taking advantage of the runtime environment variable. You can see you can even differentiate between build-time and run-time environment variables. And of course we are using the new selector as our label here.

We can apply the Workload.

```shell
$ kubectl apply -f developer/workload-2.yaml
workload.carto.run/hello-tanzu created
```

And we can see how it's doing.

```shell
$ watch kubectl tree workload hello-tanzu
NAMESPACE  NAME                                         READY    REASON               AGE
default    Workload/hello-tanzu                         Unknown  MissingValueAtPath   18s
default    ├─GitRepository/hello-tanzu                  True     GitOperationSucceed  16s
default    └─Image/hello-tanzu                          True                          14s
default      ├─Build/hello-tanzu-build-1                -                             14s
default      │ └─Pod/hello-tanzu-build-1-build-pod      False    PodCompleted         14s
default      ├─PersistentVolumeClaim/hello-tanzu-cache  -                             14s
default      └─SourceResolver/hello-tanzu-source        True                          14s
```

So you can see that already it created the Flux resource and the kpack image resource. Let's give kpack a moment to build the image. When the build is complete, we should see this MissingValueAtPath change to Ready, and we’ll see the Kapp Controller App show up.

```shell
$ watch kubectl tree workload hello-tanzu
NAMESPACE  NAME                                         READY  REASON               AGE
default    Workload/hello-tanzu                         True   Ready                59s
default    ├─App/hello-tanzu                            -                           3s
default    ├─GitRepository/hello-tanzu                  True   GitOperationSucceed  57s
default    └─Image/hello-tanzu                          True                        55s
default      ├─Build/hello-tanzu-build-1                -                           55s
default      │ └─Pod/hello-tanzu-build-1-build-pod      False  PodCompleted         55s
default      ├─PersistentVolumeClaim/hello-tanzu-cache  -                           55s
default      └─SourceResolver/hello-tanzu-source        True                        55s
```

Let's take a look at the resources that have been deployed. You can see that a simple Knative Serving configuration actually results in several resources, some manage scaling, some manage routing, and some manage revisions. I also would encourage you to look into Knative Serving if you are not already familiar with it.

```shell
$ kubectl get all --selector serving.knative.dev/service=hello-tanzu
NAME                                               READY   STATUS    RESTARTS   AGE
pod/hello-tanzu-00001-deployment-5fd8f7fc5-zrv84   2/2     Running   0          22s

NAME                                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                                      AGE
service/hello-tanzu                 ClusterIP   None             <none>        80/TCP                                       19s
service/hello-tanzu-00001           ClusterIP   100.67.19.194    <none>        80/TCP                                       21s
service/hello-tanzu-00001-private   ClusterIP   100.65.250.126   <none>        80/TCP,9090/TCP,9091/TCP,8022/TCP,8012/TCP   21s

NAME                                           READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello-tanzu-00001-deployment   1/1     1            1           22s

NAME                                                     DESIRED   CURRENT   READY   AGE
replicaset.apps/hello-tanzu-00001-deployment-5fd8f7fc5   1         1         1       22s

NAME                                            LATESTCREATED       LATESTREADY         READY   REASON
configuration.serving.knative.dev/hello-tanzu   hello-tanzu-00001   hello-tanzu-00001   True

NAME                                             CONFIG NAME   K8S SERVICE NAME   GENERATION   READY   REASON   ACTUAL REPLICAS   DESIRED REPLICAS
revision.serving.knative.dev/hello-tanzu-00001   hello-tanzu                      1            True             1                 1

NAME                                    URL                                                   READY   REASON
route.serving.knative.dev/hello-tanzu   http://hello-tanzu.default.tanzu.coraiberkleid.site   True
```

Knative has also exposed our app for us, so we can make sure it is working.

!["Hello Tanzu example"](/img/Build-blog-img2.png)

And indeed, this instance returns "hello tanzu" because we changed the message using a runtime environment variable.

Now in both of these cases we've seen the flow go from end to end, and for every new source code commit, the same flow will repeat. But it’s worth mentioning that if you update the stack, or the base OS, that kpack is using to build images, kpack will automatically rebase all of the images. The benefit of using Cartographer in conjunction with kpack is that Cartographer will notice the new images and trigger the next resource in the chain for redeployment. So supply chains with Cartographer are more flexible in terms of the triggers that they can respond to.

That's all I have. I hope this demo helps you see the value of using Cartographer as a platform for your path, or paths, to production, and I hope it shows how Tanzu Community Edition can make it easier to enhance a Kubernetes cluster.

## Join the Tanzu Community Edition Community

We are excited to hear from you and learn with you! Here are several ways you can get involved:

* Join Tanzu Community Edition's slack channel, [#tanzu-community-edition](https://kubernetes.slack.com/archives/C02GY94A8KT) on the Kubernetes workspace, and connect with maintainers and other Tanzu Community Edition users.
* Find us on [GitHub](https://github.com/vmware-tanzu/community-edition/). Suggest how we can improve the project, the docs, or share any other feedback.
* Attend our Community Meetings, with two options to choose from. Check out the [Community page](https://tanzucommunityedition.io/community/) for full details on how to attend.
