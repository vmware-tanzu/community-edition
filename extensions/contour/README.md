# Contour Extension

This package provides an ingress controller using [Contour](https://projectcontour.io/). The contour-operator is used.

## Components

* contour Namespace
* contour Custom Resources
* contour Deployment

## Usage Example

This walkthrough guides you through deploying Contour.

Push the package and push Contour to a registry 
```shell
make push-contour-extension
```

Deploy Contour to a cluster
```shell
make deploy-contour
```
