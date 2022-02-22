# Create a Package and a Package Repository

This topic describes the procedure to create a new package and a package repository.

This procedure provides the steps necessary to create software packages using OCI images, that can be discovered, installed, and managed using the package manager built into Tanzu Community Edition. The examples in this procedure refer to cert-manager.

* A description of the Tanzu Community Edition packages architecture is available [here](architecture/#package-management).  
* A description of how to work with existing packages and package repositories is available [here](package-management).  
* A package creation example is available [here](packages-create-example).

## Prerequisites

{{% include "/docs/v0.10.0/assets/prereq-pkg-creation.md" %}}

* This example uses ttl.sh registry.

## 1. Create a local directory structure

1. Create a directory for the package, and the package repository.

    ```shell
    mkdir package-example
    ```

2. Change into the directory you created in the previous step, and create further directories as follows:

    ```shell
    cd package-example
    mkdir -p bundle/.imgpkg
    mkdir -p bundle/config/overlays
    mkdir -p bundle/config/upstream
    ```

   * The `bundle/.imgpkg` directory will contain the bundle's lock file.
   * The `bundle/config/overlays` directory will contain ytt templates and overlays.
   * The `bundle/config/upstream` directory will contain the upstream manifest for cert-manager.

    For more information about the filesystem structure used for package bundle creation, see the [Package Contents Bundle](https://carvel.dev/kapp-controller/docs/latest/packaging-artifact-formats/#package-contents-bundle) topic in the Carvel documentation.

## 2. Use Vendir to Synchronize the Upstream Content to a Local Directory

{{% include "/docs/v0.10.0/assets/vendir-desc.md" %}}

1. Create the following `vendir.yml` file. This example uses cert-manager. The `vendir.yml` file indicates where to find the remote, upstream configuration for cert-manager.  It indicates to vendir  to synchronize the `config/upstream` directory created in the previous step with the contents of the cert-manager v1.5.3 GitHub release located in the `jetstack/cert-manager` repository. From that release, we want the `cert-manager.yaml` file.

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

    For the full specification of the `vendir.yml` file, see the [vendir.yml spec](https://carvel.dev/vendir/docs/latest/vendir-spec/) in the vendir documentation.

2. Run the vendir sync command to pull down the cert-manager manifest.

    ```shell
    vendir sync --chdir bundle
    ```

3. Run the following command to inspect your local `bundle/config/upstream` directory.

    ```shell
    ls -l bundle/config/upstream

    ```

    You should see the `cert-manager.yaml` file from the `v1.5.3` cert-manger release is present.

    ```shell
    -rw-r--r--  1 seemiller  staff  1442034 Oct 18 12:39 cert-manager.yaml
    ```

    You should also see the `bundle/vendir.lock.yml` file has been created. This lock file resolves the `v1.5.3` release tag to the specific GitHub release and declares that the `config/upstream` is the synchronization target path. If you inspect the file, the contents should look like this:

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

## 3. Create a ytt-annotated Kubernetes Configuration

{{% include "/docs/v0.10.0/assets/ytt-desc.md" %}}

This example uses cert-manager. In the case of cert-manager, a typically modification to make is to override the namespace that the package will be installed into. To do this, you need to make ytt overlays to replace the default namespace with one that you will provide in a configuration `values.yaml` file.
In the `bundle/config/upstream/cert-manager.yaml` file, observe that namespace appears in multiple places including annotations, the deployment, service account, and webhook manifests, as well as a few others.

1. Create the following three overlay files to modify the various places in the cert-manager manifest where the namespace is referenced. This example could use just one overlay, but it's convenient to have things separated at times.

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

2. One more file is required to hold configuration values. In this case, the only value that we can modify is the namespace, so we provide a data value for the namespace. The configuration parameters defined in this file will later be documented in the package CRD.

    ```shell
    cat > bundle/config/values.yaml <<EOF
    #@data/values
    ---

    #! The namespace in which to deploy cert-manager.
    namespace: custom-namespace
    EOF
    ```

3. To test if everything is working, run ytt. If everything is correct, ytt will output the transformed YAML. If there's a problem, you'll see it in the console.

    ```shell
    ytt --file bundle/config
    ```

## 4. Use kbld to resolve the referenced container image

The package configuration is now complete, use kbld to lock it down.

{{% include "/docs/v0.10.0/assets/kbld-desc.md" %}}

When kbld runs, it parses your configuration files and finds images. It will then lookup the images on their registries and get their `sha256` digest. This mapping will then be placed into an `images.yml` lock file in your `bundle/.imgpkg` directory. The mapping file can be used for different scenarios in the future; one being the ability to copy a package to removable media for transfer to an air-gapped network, and the second being retrieval to a cluster by kapp-controller.

1. Run the following command to create the `images.yml` file:

    ```shell
    kbld --file bundle --imgpkg-lock-output bundle/.imgpkg/images.yml 1>> /dev/null
    ```

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

## 5. Use Imgpkg to Push the Package to an OCI Registry

{{% include "/docs/v0.10.0/assets/imgpkg-desc.md" %}}

1. Use imgpkg to push the `bundle` directory and indicate what project name and tag to give it.

    ```sh
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

    If you don't specify a full URL for the registry, for example, `registry.example.com/seemiller/cert-manager:1.5.3`, imgpkg will default to DockerHub.

    Notice in the output above that imgpkg reports that it pushed `ttl.sh/cert-manager@sha256:7335d2f2...`. Take note of that URL/digest as it will be needed for the `package.yaml` file.

## 6. Create Package CRDs

The last step in the package creation process is to create two custom resource files used by the packaging API: the `package.yaml` and `metadata.yaml` files.

1. Create `package.yaml`. The Package CR is created for every new version of a package. It carries information about how to fetch, template, and deploy the package. The important information captured in this CR is as follows:

    Name  
    Version  
    License(s)  
    Image URL/digest to fetch  
    Paths for ytt/kbld files within the package  
    Arguments for deployment to kapp-controller  
    OpenAPI values schema

    For the complete specification, refer to the [Package Management](https://carvel.dev/kapp-controller/docs/latest/packaging/#package) topic in the Carvel documentation.

    Place the URL/digest that imgpkg reported after pushing the package in the `spec.template.spec.fetch.imgpkgBundle.image` field (you noted this in the [previous step](package-creation-step-by-step/#5-use-imgpkg-to-push-the-package-to-an-oci-registry)). The `.metadata.name` field must be a combination of the `spec.refName` and `spec.version` fields.

    To aid users in configuring their package, the package CRD makes a valuesSchema available. Any configurable parameter defined in the `values.yaml` and used in the ytt overlays in [step 3](package-creation-step-by-step/#3-create-a-ytt-annotated-kubernetes-configuration) should be documented here. When the package is deployed to a cluster in a package repository, a user will be able to query for these configuration parameters.

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

1. Create `metadata.yaml`. The `metadata.yaml` file contains metadata for the package that is unlikely to change with versions. Note that the `.metadata.name` value should match the name in the `package.yaml` from the previous step.

    Name  
    Descriptions, short and long  
    Maintainers  
    Categories  

    For the complete specification, refer to the [Package Metadata](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-metadata) topic in the Carvel documentation.

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

    You've now created a package for cert-manager, pushed the package to the OCI Registry, and the package CRDs are ready to be put into a package repository.

## 7. Create a Package Repository

{{% include "/docs/v0.10.0/assets/package-repository.md" %}}

A package repository provides an easy way to distribute software. A package repository can be created by a software provider to distribute different versions of their software. For example, JetStack could create a package repository that contains every version of cert-manager. You could install this package repository on a test cluster and easily swap out versions to check for compatibility with your applications. Or a training class could have a repository with cert-manager, Contour and Prometheus pre-configured to teach deploying and monitoring web applications on Kubernetes.

The steps for creating a package repository are similar to creating a package. There are five main steps required:

* Create a directory for the package repository
* Copy your package's package and metadata CRD files to a directory
* Run kbld
* Push with imgpkg
* Install to a cluster

1. Create a directory for the package repository. You will need a `packages` subdirectory as that is where the package repository expects the package CRDs to be located. A `.imgpkg` directory is also needed as this will be an imgpkg bundle.

    ```shell
    mkdir -p repo/packages
    mkdir -p repo/.imgpkg
    ```

1. Copy the package CRDs created in the previous [Create a Package CRD](package-creation-step-by-step/#6-create-package-crd) step into the `/package` directory. If you have multiple versions of the same package, you have to distinguish each `package.yaml` file with a version or concatenate them together.

    ```shell
    cp metadata.yaml repo/packages
    cp package.yaml repo/packages
    ```

1. As a package repository is expected to be an imgpkg bundle, you must run kbld to create an `image.yaml` lock file:

    ```shell
    kbld --file repo --imgpkg-lock-output repo/.imgpkg/images.yml 1>> /dev/null
    ```

1. Push the package repository to the OCI Registry.

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

1. The final step in creating a package repository is to create the PackageRepository CR. This YAML file tells the cluster the name of the package repository and where to find it. For the complete specification of the PackageRepository CRD, see the [Package Reposiotry](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-repository) topic in the Carvel documentation.

    ```sh
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

## 8. Test the Package

1. With the package and package repository both created and pushed to an OCI Registry, you can test it out. Start by deploying the PackageRepository CR to your Tanzu Community Edition cluster.

    ```shell
    kubectl apply --file pkgr.yaml
    ```

    **or**

    Here is the equivalent command using the Tanzu CLI.

    ```shell
    tanzu package repository install cert-manager-repo --url ttl.sh/seemiller/cert-manager-repo@sha256:179e9f10fd2393284eaefc34c3c95020922359caea8847d9392468d533615cf8
    ```

2. Retrieve the PackageRepositories from your cluster. Verify that the reconciliation has succeeded.

    ```shell
    kubectl get pkgr

    NAME                AGE   DESCRIPTION
    cert-manager-repo    3m   Reconcile succeeded
    ```

3. With the package repository successfully installed, you can view the packages provided by the repository.

    ```sh
    kubectl get pkg

    NAME                                                        PACKAGEMETADATA NAME                                  VERSION   AGE
    cert-manager.example.com.1.5.3                              cert-manager.example.com                              1.5.3     13m44
    ```

    **or**

    Here is the equivalent command using the Tanzu CLI:

    ```shell
    tanzu package available list
    ```
