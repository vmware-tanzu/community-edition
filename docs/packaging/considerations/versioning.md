# Versioning

This section describes the expectations of package versioning.

## Single Software Package

When a package contains a single piece of software, that version tag should
correspond to the underlying software. For example, when the contour package
[released
v1.19.1](https://github.com/vmware-tanzu/community-edition/blob/main/addons/packages/contour/1.19.1/bundle/config/upstream/contour.yaml#L4738),
it corresponded to the [contour package
v1.19.1](https://github.com/vmware-tanzu/community-edition/tree/main/addons/packages/contour/1.19.1).

## Multi-Software (Meta) Package

When a package contains multiple pieces of software, it should take on its own
semantic versioning. For example, a package called `tanzu-observability` may
contain:

* `prometheus:v2.33.0`
* `grafana:v8.3.4`

With the above, the package may be versioned `1.0.0`. The version should
increment overtime reflecting the changes (breaking, features, patches) based on
[semantic versioning](https://semver.org/).
