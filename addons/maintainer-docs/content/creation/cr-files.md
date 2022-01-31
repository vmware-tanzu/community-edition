---
title: Custom Resource Files
weight: 3
---

There are 2 custom resources required for packages, `Package` and `PackageMetadata`. Refer to the [directory structure](./directory-structure/) for where to place these.

## Package

A package is a combination of configuration metadata and OCI images that informs the package manager what software it holds and how to install itself onto a Kubernetes cluster.

Here is an annotated example from the Carvel [documentation](https://carvel.dev/kapp-controller/docs/latest/packaging/#package).

```yaml
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  # Must be of the form '<spec.refName>.<spec.version>' (Note the period)
  name: fluent-bit.carvel.dev.1.5.3
spec:
  # The name of the PackageMetadata associated with this version
  # Must be a valid PackageMetadata name (see PackageMetadata CR for details)
  # Cannot be empty
  refName: fluent-bit.carvel.dev
  # Package version; Referenced by PackageInstall;
  # Must be valid semver (required)
  # Cannot be empty
  version: 1.5.3
  # Version release notes (optional; string)
  releaseNotes: "Fixed some bugs"
  # System requirements needed to install the package.
  # Note: these requirements will not be verified by kapp-controller on
  # installation. (optional; string)
  capacityRequirementsDescription: "RAM: 10GB"
  # Description of the licenses that apply to the package software
  # (optional; Array of strings)
  licenses:
  - "Apache 2.0"
  - "MIT"
  # Timestamp of release (iso8601 formatted string; optional)
  releasedAt: 2021-05-05T18:57:06Z
  # valuesSchema can be used to show template values that
  # can be configured by users when a Package is installed.
  # These values should be specified in an OpenAPI schema format. (optional)
  valuesSchema:
    # openAPIv3 key can be used to declare template values in OpenAPIv3
    # format. Read more on using ytt to generate this schema: 
    # https://carvel.dev/kapp-controller/docs/latest/packaging-tutorial/#creating-the-custom-resources
    openAPIv3:
      title: fluent-bit.carvel.dev.1.5.3 values schema
      examples:
      - namespace: fluent-bit
      properties:
        namespace:
          type: string
          description: Namespace where fluent-bit will be installed.
          default: fluent-bit
          examples:
          - fluent-bit
  # App template used to create the underlying App CR.
  # See 'App CR Spec' docs for more info
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: registry.tkg.vmware.run/tkg-fluent-bit@sha256:...
      template:
      - ytt:
          paths:
          - config/
      - kbld:
          paths:
          # - must be quoted when included with paths
          - "-"
          - .imgpkg/images.yml
      deploy:
      - kapp: {}
```

> When populating the `Package` CR, you will need the package digest to populate the `template.spec.fetch.imgpkgBundle.image` field. This digest is obtained by pushing your package to an OCI Registry. 

## PackageMetadata

Package Metadata are attributes of a single package that do not change frequently and that are shared across multiple versions of a single package. It contains information similar to a project’s README.md. The Package Maintainer will need to provide this information about the package. This metadata will be exposed in the package repository that makes a package available to a cluster. The metadata.yaml specification is defined in the kapp-controller [documentation](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-metadata).

Here is the list of supported metadata.

* Category(s) of the type of functionality provided
* Display name
* Icon image in SVG base64 format
* Short description of the package
* Long Description of the package
* List of Package Maintainers
* Provider Name
* Where/how to find support for the package

### PackageMetadata Example

Here is an annotated example from the Carvel [documentation](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-metadata).

```yaml
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  # Must consist of at least three segments separated by a '.'
  # Cannot have a trailing '.'
  name: fluent-bit.vmware.com
  # The namespace this package metadata is available in
  namespace: my-ns
spec:
  # Human friendly name of the package (optional; string)
  displayName: "Fluent Bit"
  # Long description of the package (optional; string)
  longDescription: "Fluent bit is an open source..."
  # Short desription of the package (optional; string)
  shortDescription: "Log processing and forwarding"
  # Base64 encoded icon (optional; string)
  iconSVGBase64: YXNmZGdlcmdlcg==
  # Name of the entity distributing the package (optional; string)
  providerName: VMware
  # List of maintainer info for the package.
  # Currently only supports the name key. (optional; array of maintner info)
  maintainers:
  - name: "Person 1"
  - name: "Person 2"
  # Classifiers of the package (optional; Array of strings)
  categories:
  - "logging"
  - "daemon-set"
  # Description of the support available for the package (optional; string)
  supportDescription: "..."
```

### PackageMetadata Bootstrap Command

```shell
cat <<EOF > metadata.yaml
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: PACKAGE_NAME
spec:
  displayName: "PACKAGE_NAME"
  longDescription: ""
  shortDescription: ""
  providerName:
  maintainers:
    - name: ""
  categories:
    - ""
```
