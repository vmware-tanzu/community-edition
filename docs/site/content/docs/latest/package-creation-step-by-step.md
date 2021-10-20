# Package Creation, Step by Step

Getting started with a new workflow or process can be daunting at first. Lots of trial and error, reading documentation that doesn't apply to what you're actually trying to do, cryptic error messages. Sometimes it's nice just to have a walk through of a simple example to get you going. So, let's go through a simple exercise in creating a package that we can deploy to the Tanzu Community Edition using the Carvel suite of tools. We'll use cert-manager as the example application to package.

## Prerequisites

Before we can get started with creating a package, we'll need a few things: A Tanzu Community Edition cluster, the Carvel suite of tools and an OCI registry.

### Tanzu Community Edition

Tanzu Community Edition comes ready to go with support for packages. Any Kubernetes cluster can support packages, you just need to install [kapp-controller](https://carvel.dev/kapp-controller/). To install Tanzu Community Edition, follow the [instructions](https://github.com/vmware-tanzu/community-edition#installation). If you have Homebrew, you can run this to get going:

```shell
brew install vmware-tanzu/tanzu/tanzu-community-edition
```

### Carvel

The [Carvel](https://carvel.dev) tool suite is really where all the packaging magic happens. The `kapp-controller` is what adds the packaging API resources to your cluster. The other applications, vendir, kbld, ytt and imgpkg, are all command line applications that focus on a specific piece of the packaging process. Let us take a quick look at what each one does so we have a better understanding why they're needed.

#### vendir

vendir allows you to synchronize the contents of remote data sources into a consistent local directory. Wait, what? Basically, you define in yaml where data lives in a remote location and how you want to structure that data locally. Vendir will copy that data locally so that you can operate on it. To put it in our packaging context, vendir allows us to say, "I want to use the manifest for cert-manager, go to GitHub and retrieve release 1.5.3 and put it in my package upstream config directory."

#### kbld

kbld seamlessly incorporates image building, pushing, and resolution into your development and deployment workflows. Long story short with kbld is that it allows you to build your configuration with immutable image references. kbld scans our package configuration for image references and resolves those references to digests. This is really important in that it allows us to say, "I specified image `cert-manager:1.5.3` which is actually `quay.io/jetstack/cert-manager-controller@sha256:7b039d469ed739a652f3bb8a1ddc122942b66cceeb85bac315449724ee64287f`" This is the tool that allows you to ensure that you're using the correct versions of software.

#### ytt

Need to modify some YAML? Override a default value? Add some custom configuration? ytt lets you create templates and patches for YAML files. I like to think of it as XSLT is to XML files as ytt is to YAML files.

#### imgpkg

imgpkg allows you to package, distribute, and relocate your Kubernetes configuration and OCI images as a `bundle`. imgpkg performs operations similar to the docker and crane commands, allowing you to create, push, pull and operate on OCI images and bundles. A sha256 digest is created for the bundle based on its contents, allowing imgpkg to verify the bundle's integrity. Bundles are important in that they capture your configuration and image references as one discrete unit. As a unit, your configuration and images can be referenced and copied, which can allow for easy operation with air-gaped environments.

#### Installation

You can find installation instructions for the Carvel tools [here](https://carvel.dev/#install). If you're using Homebrew, just run:

```shell
brew tap vmware-tanzu/carvel
brew install vendir kbld ytt imgpkg
```

### OCI Registry

An OCI registry is where you will upload your package and package repository to. You can use any OCI compliant registry:

* Docker Hub
* GitHub Container Registry
* Google Container Registry
* Harbor
* ttl.sh

At the moment, the OCI registry that you choose will need to be public. Support for private registries is coming soon. For this walk through. I'll use ttl.sh. Once you've decided on a registry, be sure to authenticate locally so that you can push images (if necessary).

> While using Docker is easy, you can quickly run into rate limiting issues.

## Package Creation

We're going to need a place to put all the package files. Start by creating a directory to house the package and the package repository.

```shell
mkdir package-example
```

Change into that directory and let's create further directories that we'll need. The package contents bundle is documented [here](https://carvel.dev/kapp-controller/docs/latest/packaging-artifact-formats/#package-contents-bundle)

```shell
cd package-example
mkdir -p bundle/.imgpkg
mkdir -p bundle/config/overlays
mkdir -p bundle/config/upstream
```

The `bundle/.imgpkg` directory will contain the bundle's lock file. `bundle/config/overlays` will contain ytt templates and overlays. `bundle/config/upstream` will house the upstream manifest from cert-manager.

### vendir

With some directories in place, we can start to pull down some configuration. We'll start with the vendir.yml file. This file will tell vendir where to find the remote, upstream configuration for cert-manager. The important part of this file is in `directories`. Here we are telling vendir that we want to synchronize our `config/upstream` directory with the contents of the `v1.5.3` GitHub release located in the `jetstack/cert-manager` repository. From that release, we want the `cert-manager.yaml` file.

```shell
cat > bundle/vendir.yml <<EOF
apiVersion: vendir.k14s.io/v1alpha1
kind: Config
minimumRequiredVersion: 0.12.0
directories:
  - path: config/upstream
    contents:
      - path: .
        githubRelease:
          slug: jetstack/cert-manager
          tag: v1.5.3
          disableAutoChecksumValidation: true
        includePaths:
          - cert-manager.yaml
EOF
```

> For the full specification of the vendir.yml file, see the vendir [documentation](https://carvel.dev/vendir/docs/latest/vendir-spec/).

Run the vendir sync command to pull down the cert-manager manifest.

```shell
vendir sync --chdir bundle
```

Inspect the `bundle/config/upstream` directory. You'll see that the `cert-manager.yaml` file from the `v1.5.3` release is present.

```shell
ls -l bundle/config/upstream

-rw-r--r--  1 seemiller  staff  1442034 Oct 18 12:39 cert-manager.yaml
```

You'll also notice that a `bundle/vendir.lock.yml` file has been created. This lock file resolves the `v1.5.3` release tag to the specific GitHub release and declares that the `config/upstream` is the synchronization target path. If you inspect the file, the contents should look like this:

```yaml
apiVersion: vendir.k14s.io/v1alpha1
directories:
- contents:
  - githubRelease:
      url: https://api.github.com/repos/jetstack/cert-manager/releases/48370396
    path: .
  path: config/upstream
kind: LockConfig
```

### ytt

Second, we can add overlays, templates and custom configuration. ytt is complicated, and isn't this post's focus, so we'll try to keep it as simple as possible. A typically easy modification to make is to override the namespace that the package will be installed into. To do this, we'll need to make ytt overlays to replace the default namespace with one that we provide in a configuration `values.yaml` file. Searching through the `bundle/config/upstream/cert-manager.yaml` file, we see that namespace appears in annotations, the deployment, service account, and webhook manifests, as well as a few others. This means that we'll need to override in quite a few places. Since this post is not about ytt, I'll ask you to trust me that these files work and forgo detailed explanations. For more information about ytt, refer to the official [documentation](https://carvel.dev/ytt/docs/latest/).

We'll create 3 overlay files to modify the various places in the cert-manager manifest where the namespace is referenced. We could have used just 1 overlay, but it's convenient to have things separated at times.

```shell
cat > bundle/config/overlays/annotations.yaml <<EOF
#@ load("@ytt:data", "data")
#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.subset({"kind":"CustomResourceDefinition"}), expects=6
---
metadata:
  annotations:
    #@overlay/match missing_ok=True
    cert-manager.io/inject-ca-from-secret: #@ "{}/cert-manager-webhook-ca".format(data.values.namespace)

#@overlay/match by=overlay.subset({"kind":"MutatingWebhookConfiguration"})
---
metadata:
  annotations:
    #@overlay/match missing_ok=True
    cert-manager.io/inject-ca-from-secret: #@ "{}/cert-manager-webhook-ca".format(data.values.namespace)

#@overlay/match by=overlay.subset({"kind":"ValidatingWebhookConfiguration"})
---
metadata:
  annotations:
    #@overlay/match missing_ok=True
    cert-manager.io/inject-ca-from-secret: #@ "{}/cert-manager-webhook-ca".format(data.values.namespace)
EOF
```

```shell
cat > bundle/config/overlays/deployment.yaml <<EOF
#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind": "Deployment", "metadata": {"name": "cert-manager-webhook"}})
---
spec:
  template:
    spec:
      containers:
      #@overlay/match by="name"
      - name: cert-manager
        args:
          #@overlay/match by=lambda i,l,r: l.startswith("--dynamic-serving-dns-names=")
          - #@ "--dynamic-serving-dns-names=cert-manager-webhook,cert-manager-webhook.{},cert-manager-webhook.{}.svc".format(data.values.namespace, data.values.namespace)
EOF
```

```shell
cat > bundle/config/overlays/misc.yaml <<EOF
#@ load("@ytt:data", "data")
#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.subset({"kind":"Namespace", "metadata": {"name": "cert-manager"}})
---
apiVersion: v1
kind: Namespace
metadata:
  name: #@ data.values.namespace

#@overlay/match by=overlay.subset({"metadata": {"namespace": "cert-manager"}}), expects=10
---
metadata:
  namespace: #@ data.values.namespace

#@ crb=overlay.subset({"kind":"ClusterRoleBinding"})
#@ rb=overlay.subset({"kind":"RoleBinding"})
#@overlay/match by=overlay.or_op(crb, rb), expects=13
---
subjects:
#@overlay/match by=overlay.subset({"namespace": "cert-manager"})
- kind: ServiceAccount
  namespace: #@ data.values.namespace

#@ vwc=overlay.subset({"kind":"ValidatingWebhookConfiguration"})
#@ mwc=overlay.subset({"kind":"MutatingWebhookConfiguration"})
#@overlay/match by=overlay.or_op(vwc, mwc), expects=2
---
webhooks:
#@overlay/match by="name"
- name: webhook.cert-manager.io
  clientConfig:
    service:
      namespace: #@ data.values.namespace

#@overlay/match by=overlay.subset({"kind":"CustomResourceDefinition"}), expects=6
---
spec:
  conversion:
    webhook:
      clientConfig:
        #@overlay/match by="name"
        service:
          name: cert-manager-webhook
          namespace: #@ data.values.namespace
EOF
```

Finally, we'll need one more file to hold our configuration values. In this case, the only value that we can modify is the namespace, so we provide a data value for the namespace. The configuration parameters defined in this file will later be documented in the package CRD.

```shell
cat > bundle/config/values.yaml <<EOF
#@data/values
---

#! The namespace in which to deploy cert-manager.
namespace: custom-namespace
EOF
```

To test if everything is working, we can run ytt. If everything is correct, ytt will output the transformed YAML. If there's a problem, you'll see it in the console.

```shell
ytt --file bundle/config
```

### kbld

Now that we've defined the configuration for the package, we need to lock it down. kbld will search through our configs looking for any references to images and create a mapping of image tags to a URL with a `sha256` digest. Images with the same name and tag on different registries are not necessiarly the same images! By referring to an image with a digest, you're guaranteed to get the image that you're expecting. This is basically the same as providing a checksum file alongside an executable on a download site.

When kbld runs, it parses your configuration files and finds images. It will then lookup the images on their registries and get their `sha256` digest. This mapping will then be placed into an `images.yml` lock file in the `bundle/.imgpkg` directory. The mapping file can be used for a different scenarios in the future: one being the ability to copy a package to removable media for transfer to an air-gapped network and the second being the ultimate retrieval to a cluster by kapp-controller.

```shell
kbld --file bundle --imgpkg-lock-output bundle/.imgpkg/images.yml 1>> /dev/null
```

> I've piped stdout to `/dev/null` as kbld is a bit noisy outputing all the config that it parsed.

Here is what the `images.yml` file should look like.

```shell
cat bundle/.imgpkg/images.yml
---
apiVersion: imgpkg.carvel.dev/v1alpha1
images:
- annotations:
    kbld.carvel.dev/id: quay.io/jetstack/cert-manager-cainjector:v1.5.3
    kbld.carvel.dev/origins: |
      - resolved:
          tag: v1.5.3
          url: quay.io/jetstack/cert-manager-cainjector:v1.5.3
  image: quay.io/jetstack/cert-manager-cainjector@sha256:de02e3f445cfe7c035f2a9939b948c4d043011713389d9437311a62740f20bef
- annotations:
    kbld.carvel.dev/id: quay.io/jetstack/cert-manager-controller:v1.5.3
    kbld.carvel.dev/origins: |
      - resolved:
          tag: v1.5.3
          url: quay.io/jetstack/cert-manager-controller:v1.5.3
  image: quay.io/jetstack/cert-manager-controller@sha256:7b039d469ed739a652f3bb8a1ddc122942b66cceeb85bac315449724ee64287f
- annotations:
    kbld.carvel.dev/id: quay.io/jetstack/cert-manager-webhook:v1.5.3
    kbld.carvel.dev/origins: |
      - resolved:
          tag: v1.5.3
          url: quay.io/jetstack/cert-manager-webhook:v1.5.3
  image: quay.io/jetstack/cert-manager-webhook@sha256:ed6354190d259524d32ae74471f93bf46bfdcf4df6f73629eedf576cd87e10b8
kind: ImagesLock
```

### imgpkg

At this point, we are basically done creating the package. All that is left to do is to push the package to an OCI Registry. To do that, we will use imgpkg. There's not much to say here, simply tell imgpkg to push the `bundle` directory and what project name and tag to give it. If you want to learn more about imgpkg, you can refer to the [documentation](https://carvel.dev/imgpkg/docs/latest/).

```shell
imgpkg push --bundle ttl.sh/seemiller/cert-manager:6h --file bundle/

dir: .
dir: .imgpkg
file: .imgpkg/images.yml
dir: config
dir: config/overlays
file: config/overlays/annotations.yaml
file: config/overlays/deployment.yaml
file: config/overlays/misc.yaml
dir: config/upstream
file: config/upstream/cert-manager.yaml
file: config/values.yaml
file: vendir.lock.yml
file: vendir.yml
Pushed 'ttl.sh/seemiller/cert-manager@sha256:7335d2f20d000695e7880731ad24406c3d98dff5008edb04fa30f16e31abbb1a'
Succeeded
```

> If you don't epecify a full URL for the registry, e.g. registry.example.com/seemiller/cert-manager:1.5.3, imgpkg will default to DockerHub.

Notice in the output above that imgpkg reports that it pushed `ttl.sh/cert-manager@sha256:7335d2f2...`. Take note of that URL/digest as it will be needed for the package.yaml file.

### Package CRDs

The last step in our package creation process will be to create 2 custom resource files used by the packaging API, the `package.yaml` and `metadata.yaml` files.

#### package.yaml

The Package CR is created for every new version of a package. It carries information about how to fetch, template, and deploy the package. The important information captured in this CR follows. For the complete specification, refer to the [documentation](https://carvel.dev/kapp-controller/docs/latest/packaging/#package).

* Name
* Version
* License(s)
* Image URL/digest to fetch
* Paths for ytt/kbld files within the package
* Arguments for deployment to kapp-controller
* OpenAPI values schema

Recall from the previous step that you needed to save off the URL/digest that imgpkg reported after pushing the package. You'll need to place that in the `spec.template.spec.fetch.imgpkgBundle.image` field. Also of note is the `.metadata.name` field. This field must be a combination of the `spec.refName` and `spec.version` fields.

To aid users in configuring their package, the package CRD makes a valuesSchema available. Any configurable parameter defined in the values.yaml and used in the ytt overlays/templates should be documented here. When the package is deployed to a cluster in a package repository, a user will be able to query for these configuration parameters.

```shell
cat > package.yaml <<EOF
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: cert-manager.example.com.1.5.3
spec:
  refName: cert-manager.example.com
  version: 1.5.3
  releaseNotes: "relevant release notes for this version..."
  licenses:
    - "Apache 2.0"
  template:
    spec:
      fetch:
        - imgpkgBundle:
            image: ttl.sh/seemiller/cert-manager@sha256:7335d2f20d000695e7880731ad24406c3d98dff5008edb04fa30f16e31abbb1a
      template:
        - ytt:
            paths:
              - config/
        - kbld:
            paths:
              - "-"
              - .imgpkg/images.yml
      deploy:
      - kapp: {}
  valuesSchema:
    openAPIv3:
      title: cert-manager.example.com.1.5.3 values schema
      examples:
        - namespace: cert-manager
      properties:
        namespace:
          type: string
          description: The namespace in which to deploy cert-manager.
          default: cert-manager
EOF
```

#### metadata.yaml

The metadata.yaml file contains metadata for the package that is unlikely to change with versions. Note that the `.metadata.name` value should match the name in the package.yaml from the previous step.

* Name
* Descriptions, short and long
* Maintainers
* Categories

For the complete specification, refer to the [documentation](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-metadata).

```shell
cat > metadata.yaml <<EOF
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: cert-manager.example.com
spec:
  displayName: "cert-manager"
  longDescription: "A long description"
  shortDescription: "A short description"
  providerName: packages-r-us
  maintainers:
    - name: Seemiller
  categories:
    - certificate-management
EOF
```

And with that, we've completed creating a package for cert-manager. Our package has been pushed to the OCI Registry, and its CRDs are ready to be put into a package repository.

## Package Repository

A package repository is a collection of packages. The Tanzu Community Edition project provides a package repository that we feel represents a solid set of 3rd party software necessary to start building an application platform on Kubernetes. A package repository could be created by a software provider to distribute different versions of their software. Imagine JetStack creating a package repository that contains every version of cert-manager. You could install this package repository on a test cluster and easily swap out versions to check for compatibility with your applications. Or imagine a training class that has a repository with cert-manager, Contour and Prometheus pre-configured ready to go to teach deploying and monitoring web applications on Kubernetes. Whatever the use, a package repository makes for an easy way to distribute software.

### Creating A Package Repository

Creating a package repository is pretty straight forward and similar to creating a package. You just need to:

1. Copy your package's package and metadata CRD files to a directory
2. Run kbld
3. Push with imgpkg
4. Install to a cluster

We'll start by creating a new directory for the package repository. We will need a `packages` subdirectory as that is where the package repository expects the package CRDs to be located. A `.imgpkg` direcotry is also needed as this will be an imgpkg bundle.

```shell
mkdir -p repo/packages
mkdir -p repo/.imgpkg
```

Copy the package CRDs into that directory. If you had multiple versions of the same package, you would have to distinguish each package.yaml file with a version or concatenate them together.

```shell
cp metadata.yaml repo/packages
cp package.yaml repo/packages
```

Since a Package Repository is expected to be an imgpkg bundle, we will need to run kbld to create an image.yaml lock file, just like we did for the package.

```shell
kbld --file repo --imgpkg-lock-output repo/.imgpkg/images.yml 1>> /dev/null
```

Then push the package repository to the OCI Registry.

```shell
imgpkg push --bundle ttl.sh/seemiller/cert-manager-repo:6h --file repo/

dir: .
dir: .imgpkg
file: .imgpkg/images.yml
dir: packages
file: packages/metadata.yaml
file: packages/package.yaml
Pushed 'ttl.sh/seemiller/cert-manager-repo@sha256:179e9f10fd2393284eaefc34c3c95020922359caea8847d9392468d533615cf8'
Succeeded
```

Once again, notice the URL/digest that imgpkg reported that it pushed, `ttl.sh/seemiller/cert-manager-repo@sha256:179e9f10...`. This value will be used in the next step.

The final step in creating a package repository is to create the PackageRepository CR. This YAML file tells the cluster the name of the package repository and where to find it. For the complete specification of the PackageRepository CRD, see the [documentation](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-repository).

```shell
cat > pkgr.yaml <<EOF
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: cert-manager-repo
spec:
  fetch:
    imgpkgBundle:
      image: ttl.sh/seemiller/cert-manager-repo@sha256:179e9f10fd2393284eaefc34c3c95020922359caea8847d9392468d533615cf8
EOF
```

The PackageRepository CRD tells Kubernetes where to find the bundle for your package repository. Package repositories can be retrieved from imgpkg bundles, images, git repositories or a file via HTTP. For the complete specification, refer to the [documentation](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-repository).

## Deployment

With the package and package repository both created and pushed to an OCI Registry, we can test it out. Start by deploying the PackageRepository CR to your Tanzu Community Edition cluster.

```shell
kubectl apply --file pkgr.yaml
```

Here is the equivalent command using the Tanzu CLI.

```shell
tanzu package repository install cert-manager-repo --url ttl.sh/seemiller/cert-manager-repo@sha256:179e9f10fd2393284eaefc34c3c95020922359caea8847d9392468d533615cf8
```

After a few seconds, retrieve the PackageRepositories from your cluster. Verify that the reconciliation has succeeded.

```shell
kubectl get pkgr

NAME                AGE   DESCRIPTION
cert-manager-repo    3m   Reconcile succeeded
```

With the package repository successfully installed, you can view the packages provided by the repository.

```shell
kubectl get pkg

NAME                                                        PACKAGEMETADATA NAME                                  VERSION   AGE
cert-manager.example.com.1.5.3                              cert-manager.example.com                              1.5.3     13m44
```

Or with the Tanzu CLI.

```shell
tanzu package available list
```

And there you have it. Creating a package and making it available on your Tanzu Community Edition cluster, step by step.
