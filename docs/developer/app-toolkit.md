# TCE App Toolkit Adding Packages How-To

## Summary

This document describes how to add an additional TCE package to the App Toolkit meta-package.

## Assumptions

This document is intended for maintainers of official repositories, as defined in the [community-edition/docs/packaging](https://github.com/vmware-tanzu/community-edition/tree/main/docs/packaging) README. Currently the role is limied to VMware employees.

## Concepts

### Meta-package

A meta-package is defined as a TCE package that installs a set of other already present TCE packages. It is not intended to add new functionality on its own, but deploy a set of existing functionalities to achieve a particular outcome.

### App Toolkit

App Toolkit is an instance of a meta-package with a goal of allowing TCE users to deploy a set of TCE packages sufficient to deploy a workload. App Toolkit itself provides little/no novel functionality, but comprises an opinion about what TCE packages are required to achieve the outcome.

The OSS App Package concept is comparable to proprietary Tanzu Application Platform, although reduced in scope.

### Configuration of App Toolkit and Included TCE Packages

Configuration is deferred to the included TCE package. For instance, if end users are interested in changing the default behavior of, say, knative-serving, they will learn about and provide configuration for knative-serving, not app-toolkit.

### Similarities and Differences to TCE Packaging of Upstream Projects

TCE provides guidance about how to include other upstream projects into TCE - e.g. kpack, knative, etc., as described in the [TCE Documentation](https://tanzucommunityedition.io/docs/package-creation-step-by-step/). Principally it describes how to

- Use Vendir to syncronize upstream content to a local directory.
- Import/pin a version of the upstream's manifests
- Create ytt overlays to alter via templating those manifests for use in TCE.

These importation steps are NOT APPLICABLE or required for App Toolkit. As a meta-package, it refers to _packages that have already been imported into TCE_, so any packaging concerns are dealt with in the included package; and _adds no new functionality itself_ so does not have a _vmware-tanzu/package-for-${PACKAGE_NAME}_ -style repo of its own. App Toolkit is wholly an artifact within TCE.

App Toolkit _does_ leverage the other structures of the TCE packaging process, to take advantage of the machinery as well as be more easily understandable within the TCE context.

## OSS Contribution process

This may have been completed already, or be the responsibility of a PM, not a developer. The end goal of this process is to have a pull request accepted into TCE. The process is documented in the [TCE Proposal Process](https://github.com/vmware-tanzu/community-edition/tree/main/docs/designs). This document will focus on the developer work required.

## Add a Package to App Toolkit- Walkthrough

1. Get the repository

    Follow your preferred process to collaborate on a git project. For instance,

    1. Fork the [TCE repo](https://github.com/vmware-tanzu/community-edition)
    1. git checkout the fork
    1. Create a new feature branch and switch to it

1. Add a new minor version of app-toolkit.
    Follow semver principles - adding a new package likely introduces new capabilities while introducing no breaking changes, so is probably a minor release.

    ```shell
    cd addons/packages/app-toolkit
    cp -r [a.N.c] [a.N+1.0]
    ```

    Where a is the major version, N is the current minor, and N+1 is the new minor version. Incrementing minor version resets the c/bugfix version to 0.

1. Add/edit boilerplate files.
    1. [version]/bundle/.imgpkg/bundle.yml

        Consider adding yourself as an author.
    1. [version]/bundle/.imgpkg/images.yml

        Nothing to add or change here! This file tracks what upstream images are required for the project; app-toolkit has no upstream packages.
    1. [version]/config/rbac

        Likely nothing to add or change here - this template adds service accounts under which to add the package's components to TCE. You may need to add a role here if you are interacting with a different part of the system.
    1. [version]/README.md

        Edit to describe your new package, and any other changes introduced.

    1. hack/

        This directory contains bash scripts to help install, uninstall etc packages during development and are not used after release. Use these or add additional functionality as convenient.
    1. metadata.yaml

        Applies to the overall App Toolkit package, so unlikely to change. Review if the maintainer, description or other data is changing.

1. config/[your-new-package]

    Create a new template file to contain the PackageInstall and other CRDs to tell `kapp` how to deploy your package within App Toolkit. The existing files in config may serve as examples. Below is the kpack.yaml that includes kpack into app-toolkit.

    ```yaml
    #@ load("@ytt:data", "data")
    #@ load("@ytt:yaml", "yaml")

    ---
    apiVersion: packaging.carvel.dev/v1alpha1
    kind: PackageInstall
    metadata:
      name: kpack
      namespace: app-toolkit-install
      annotations:
        kapp.k14s.io/change-group: "kpack"
        kapp.k14s.io/change-rule: "delete before deleting serviceaccount"
    spec:
      serviceAccountName: app-toolkit-install-sa
      packageRef:
        refName: kpack.community.tanzu.vmware.com
        versionSelection:
        constraints: 0.5.0
        prereleases: {}
      values:
      - secretRef:
          name: kpack-values
    ---
    apiVersion: v1
    kind: Secret
    metadata:
      name: kpack-values
      namespace: app-toolkit-install
    stringData:
      values.yaml: #@ yaml.encode(data.values.kpack)
    ```

    1. `metadata.name` The name of your-new-package as it will be referenced within the cluster
    1. `metadata.annotations.kapp.*` Describe how Kapp should lifecycle your-new-app. Can be simple as shown, although presumably new-apps could have specific requirements.
    1. `spec.packageRef` Points to the package that will be installed. Note that the package referred to here must be separately made visible to kapp

1. package.yaml

    1. `metadata.name` Update as appropriate for version changes
    1. `spec.version` Update as appropriate for version changes
    1. `spec.template.spec.fetch.imgpkgBundle.image` Update as approprite, see next step
    1. spec.valuesSchema

        This is a passthrough section for configuration details that will be overlaid to the sub-package's configuration. Add a section for your-new-package.

        ```yaml
        ...
        properties:
          your-new-package:
            type: object
            default: {}
            description: "contour values"
        ...
        ```

1. Push a new app-toolkit bundle.

    ```shell
    imgpkg push -b [registry]/app-toolkit-package-bundle:0.1.0 -f bundle/
    ```

    A public docker hub example:

    ```shell
    imgpkg push -b index.docker.io/csamp/app-toolkit-package-bundle:0.1.0 -f bundle/
    ```

    A Tanzu dev example:

    ```shell
    imgpkg push -b dev.repository.tanzu.vmware.com/app-toolkit-package-bundle:0.1.0 -f bundle/ -u [username] -p [password]
    ```

1. Run kapp deploy, until the upstream repository is updated.

1. `tanzu package install app-toolkit -p app-toolkit.community.tanzu.vmware.com -v 0.1.0 -n app-toolkit-install`
