---
title: Configuration
weight: 3
---

The package maintainer is responsible for creating a [schema](https://carvel.dev/ytt/docs/latest/how-to-write-schema/) that details the available configuration parameters. This schema can be used to generate package documentation and a default values.yaml file for use when installing the package.

For each configuration item listed in the schema, there should be a corresponding section in a ytt overlay that makes use of it.

## Sample schema.yaml

```yaml
#! schema.yaml

#@data/values-schema
#@schema/desc "OpenAPIv3 Schema for secretgen-controller"
---
#@schema/desc "Configuration for secretgen-controller"
secretgenController:
  #@schema/desc "The namespace in which to deploy secretgen-controller"
  namespace: secretgen-controller
  #@schema/desc "Whether to create namespace specified for secretgen-controller"
  createNamespace: true
```