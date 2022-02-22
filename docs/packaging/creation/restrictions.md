# Restrictions

This page documents items that are restricted from use by Tanzu Community Edition packages.

## Alpine Images

Due to licensing concerns, Alpine images are not allowed. If a package contains software that uses an Alpine image, that software will need to be rebuilt using a suitable alternative base image.

You can check your package images by running the `imagelinter` utility. This utility is located in the `hack/imagelinter` directory of the Tanzu Community Edition source repository. Instructions on how to build, configure and run `imagelinter` are located there as well.
