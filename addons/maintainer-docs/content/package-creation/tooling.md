---
title: Tooling
---

To develop packages for Tanzu Community Edition, the following [Carvel](https://carvel.dev) tools are required. These tools should be installed on any development machine and available in any CI pipeline.

## imgpkg

[imgpkg](https://carvel.dev/imgpkg/) is a Carvel tool that enables you to package, distribute, and relocate your Kubernetes configuration and OCI images as a bundle. Imgpkg performs operations similar to the docker and crane commands, allowing you to create, push, pull, and operate on OCI images and bundles. A sha256 digest is created for the bundle based on its contents, allowing imgpkg to verify the bundle's integrity. Bundles are important because they capture your configuration and image references as one discrete unit. As a unit, your configuration and images can be referenced and copied. Referencing your configuration and images as a unit allows for easy operation with air-gaped environments.

## kbld

[kbld](https://carvel.dev/kbld/)  is a Carvel tool that enables you to ensure that you're using the correct versions of software when you are creating a package. It allows you to build your package configuration with immutable image references. kbld scans a package configuration for image references and resolves those references to digests. For example, it allows you to  specify image `cert-manager:1.5.3` which is actually `quay.io/jetstack/cert-manager-controller@sha256:7b039d469ed739a652f3bb8a1ddc122942b66cceeb85bac315449724ee64287f`.

kbld scans a package configuration for any references to images and creates a mapping of image tags to a URL with a `sha256` digest. As images with the same name and tag on different registries are not necessarily the same images, by referring to an image with a digest, you're guaranteed to get the image that you're expecting. This is similar to providing a checksum file alongside an executable on a download site.

## vendir

[vendir](https://carvel.dev/vendir/) is a Carvel tool used in the package creation process. Use vendir to synchronize the contents of remote data sources into a consistent local directory. Use a YAML file to define the remote data location and how you want to structure that data locally. Vendir will copy the data locally so that you can operate on it. For example, you can indicate in a YAML file that you want to retrieve the manifest for cert-manager v1.5.3 in GitHub, and put it in a local `/config/upstream` directory.

## ytt

[ytt](https://carvel.dev/ytt/) is a Carvel templating tool that dynamically overwrites values used in YAML files.  ytt is used to override default values and add custom configurations in yaml files. ytt lets you create templates and patches for YAML file. ytt interacts with YAML files similarly to how XSLT interacts with XML files.
