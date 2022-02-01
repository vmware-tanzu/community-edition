# Maintenance

## Version Updates

As a bleeding-edge distribution of Tanzu, we recommend package maintainers stay
up to date with the versions of their underlying software.

As described in [versioning](../considerations/versioning.md), versions should track their
underlying software. Thus, package updates should track software updates and
conform to [guarantees provided in semver](https://semver.org/).

In cases where packages **must** be patched before a software update occurs,
build metadata can be added to signify and update. Syntactically, this is
represented by `${SOFTWARE_VERSION}+${BUILD_ITERATION}`. For example, consider a
patch that must occur to the existing `pkg-contour:v1.19.1` package. It would be
updated to `pkg-contour:v1.19.1+1`. Most times, packages will not need to add
this additional build metadata.
