---
title: Metadata
weight: 3
---

The author will need to provide general metadata about the package. This metadata will be exposed in the package repository that makes a package available to a cluster. The metadata.yaml specification is defined in the kapp-controller [documentation](https://carvel.dev/kapp-controller/docs/latest/packaging/#package-metadata).

## Metadata Items

* Category(s) of the type of functionality provided
* Display name
* Icon image in SVG base64 format
* Short description of the package
* Long Description of the package
* List of Package Maintainers
* Provider Name
* Where/how to find support for the package

## Sample metadata.yaml CR

```yaml
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: PACKAGE_NAME.community.tanzu.vmware.com
spec:
  displayName: "PACKAGE_NAME"
  longDescription: ""
  shortDescription: ""
  providerName: VMware
  maintainers:
    - name: ""
  categories:
    - ""
```