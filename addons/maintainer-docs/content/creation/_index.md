---
title: Creation
weight: 1
---

Creating a Package is quite straightforward. Here is a high level overview of the steps to create a package. Click through the links in each step for more details and examples.

1. Ensure that you have meet all the [prerequisites](../creation/#Prerequisites)
2. Create the [directory structure](../creation/directory-structure/#bootstrap-commands)
3. Use [vendir](../creation/tooling/#vendir) to synchronize [upstream content](../creation/upstream-content/#example-usage)
4. Use [kbld](../creation/tooling/#kbld) to create [immutable image references](../creation/image-refs/#example-usage)
5. Define [configurable parameters](../creation/configuration/) in the schema
6. Create overlays to apply configuration over the upstream content
7. Create unit tests
8. Create Documentation
9. Create Metadata and Package CR files
10. Push the package to an OCI Registry

## Prerequisites

In order to complete the creation of a package, you will need the following.

### Tanzu Community Edition Cluster

Access to an installed Tanzu Community Edition Cluster. The installer and instructions can be found [here](https://tanzucommunityedition.io/download/).

### Carvel Tools

The following Carvel tools are required to create a package. You can learn more about the individual tools [here](../creation/tooling/).

* imgpkg
* kbld
* vendir
* ytt

### OCI Registry

An [OCI Registry](../creation/oci-registry) is required to host the package images.

-----

content to rework

## 5. Use Imgpkg to Push the Package to an OCI Registry

include "/docs/assets/imgpkg-desc.md" %}}

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

 include "/docs/assets/package-repository.md" %}}

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

-----

## Upstream Dependencies

When relevant, point at upstream dependencies. For example, if the software being packaged is available as a release on GitHub, reference the release and version in the `vendir.yml` file. By running `vendir`, the upstream resources will be downloaded to your local package. Doing this ensures that you have the proper upstream resources.

## Consolidate Upstream Container Images

```shell
wget -q -O -
https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml | grep -i image:
          image: "quay.io/jetstack/cert-manager-cainjector:v1.6.1"
          image: "quay.io/jetstack/cert-manager-controller:v1.6.1"
          image: "quay.io/jetstack/cert-manager-webhook:v1.6.1"
```

```shell
$ crane copy quay.io/jetstack/cert-manager-controller:v1.6.1 projects.registry.vmware.com/tce/images/jetstack/cert-manager-controller:v1.6.1

2021/12/13 16:49:13 Copying from quay.io/jetstack/cert-manager-controller:v1.6.1 to projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller:v1.6.1
2021/12/13 16:49:15 existing manifest: sha256:41917b5d23b4abe3f5c34a156b1554e49e41185431361af46640580e4d6258fc
2021/12/13 16:49:16 existing blob: sha256:ec52731e927332d44613a9b1d70e396792d20a50bccfa06332a371e1c68d7785
2021/12/13 16:49:17 existing blob: sha256:dc34538f67ce001ae34667e7a528f5d7f1b7373b4c897cec96b54920a46cde65
2021/12/13 16:49:17 pushed blob: sha256:a6dbf7b27db03dd5a6e8d423d831a2574a72cc170d47fbae95318d3eeae32149
2021/12/13 16:49:57 pushed blob: sha256:29e5180199b812b0af5fe3d7cbe11787ba3234935537ec14ad0adf56847f005d
2021/12/13 16:49:58 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller@sha256:e2be0d9dfa684e1abf5ef9b24b601b1ca6b9dd6d725342b13c18b44156518b49: digest: sha256:e2be0d9dfa684e1abf5ef9b24b601b1ca6b9dd6d725342b13c18b44156518b49 size: 947
2021/12/13 16:49:59 existing blob: sha256:ec52731e927332d44613a9b1d70e396792d20a50bccfa06332a371e1c68d7785
2021/12/13 16:49:59 existing blob: sha256:dc34538f67ce001ae34667e7a528f5d7f1b7373b4c897cec96b54920a46cde65
2021/12/13 16:50:00 pushed blob: sha256:24882da6a70629e1639eb5bff873474c56a8c794a4adeca7cde9ed3fcda12102
2021/12/13 16:50:42 pushed blob: sha256:313817109359e805c69c3824ca6bc0a4a491e8b418399f0beea479d140541973
2021/12/13 16:50:43 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller@sha256:8898cc51a41a7848076cd7735e5a86feee734f13e802c563ef1deaafe6685040: digest: sha256:8898cc51a41a7848076cd7735e5a86feee734f13e802c563ef1deaafe6685040 size: 947
2021/12/13 16:50:44 existing blob: sha256:ec52731e927332d44613a9b1d70e396792d20a50bccfa06332a371e1c68d7785
2021/12/13 16:50:44 existing blob: sha256:dc34538f67ce001ae34667e7a528f5d7f1b7373b4c897cec96b54920a46cde65
2021/12/13 16:50:45 pushed blob: sha256:0714e6c1a7c35f6ea4fa848f83b7a8f341e3dcf44b5a5721fc569132d151a40c
2021/12/13 16:51:23 pushed blob: sha256:b68f7fa8b507c96446c17634e98eadacfac7b0473da27558ea4c9df64edd0fb6
2021/12/13 16:51:24 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller@sha256:7a60aca7f3c33e58f722229a139514b24cee45881b4c39428ae3cc252ef3190d: digest: sha256:7a60aca7f3c33e58f722229a139514b24cee45881b4c39428ae3cc252ef3190d size: 947
2021/12/13 16:51:25 existing blob: sha256:ec52731e927332d44613a9b1d70e396792d20a50bccfa06332a371e1c68d7785
2021/12/13 16:51:25 existing blob: sha256:dc34538f67ce001ae34667e7a528f5d7f1b7373b4c897cec96b54920a46cde65
2021/12/13 16:51:26 pushed blob: sha256:19542d9fe421c98aa84668010a0842501e30f6a99007846962ec1f2bcf6f6b37
2021/12/13 16:52:14 pushed blob: sha256:2a38dfa462ca3cb493a46809d9f587c3df314c96c62697a9a23aad9782f00990
2021/12/13 16:52:14 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller@sha256:1faa4c99e61db1e2227ca074de4e40c4e9008335f009fd6fd139c07ac4d5024b: digest: sha256:1faa4c99e61db1e2227ca074de4e40c4e9008335f009fd6fd139c07ac4d5024b size: 947
2021/12/13 16:52:15 projects.registry.vmware.com/tce/packages/jetstack/cert-manager-controller:v1.6.1: digest: sha256:fef465f62524ed89c27451752385ab69e5c35ea4bc48b62bf61f733916ea674c size: 1723
```

## Resolve Digest Values for Images

[kbld](https://carvel.dev/kbld/) is a Carvel tool that enables you to ensure that you're using the correct versions of software when you are creating a package. It allows you to build your package configuration with immutable image references. kbld scans a package configuration for image references and resolves those references to digests. For example, it allows you to  specify image `cert-manager:1.5.3` which is actually `quay.io/jetstack/cert-manager-controller@sha256:7b039d469ed739a652f3bb8a1ddc122942b66cceeb85bac315449724ee64287f`.

kbld scans a package configuration for any references to images and creates a mapping of image tags to a URL with a `sha256` digest. As images with the same name and tag on different registries are not necessarily the same images, by referring to an image with a digest, you're guaranteed to get the image that you're expecting. This is similar to providing a checksum file alongside an executable on a download site.

## User Configurable Values: Schema

Creating a schema...

## Overlays

Overlays provide a means for the package maintainer to modify or configure the behavior of the underlying software in the package. Overlays are processed by [ytt](https://carvel.dev/ytt/).

The package maintainer will create a `schema.yaml` file that contains the configuration values available in the package. For each configuration value there should be a template or overlay that modifies the underlying software's configuration.

## Run linting checks

* Validate Markdown for documentation
* Validate schema is present
* Validate overlays and templates
* TODO - more checks

## Push Package Bundle to an OCI Registry

Initial software components that are to be packaged
Configuration parameters that are to be exposed
Basic tests should be provided along with execution instructions.

## Open Pull Request in community-edition repo

When your package is complete and ready for acceptance, create a Pull Request in the [community-edition](https://github.com/vmware-tanzu/community-edition/pulls) GitHub repository. The PR should reference the GitHub issue created for introducing the package.
