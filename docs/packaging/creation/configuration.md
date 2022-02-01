# Configuration

## Schema

The package maintainer is responsible for creating a [schema](https://carvel.dev/ytt/docs/latest/how-to-write-schema/) that details the available configuration parameters. This schema can be used to generate package documentation and a default `values.yaml` file for use when installing the package.

For each configuration item listed in the schema, there should be a corresponding section in a ytt overlay that makes use of it. Ideally, there is also a unit test to verify the overlay is working properly.

### Sample schema.yaml

```yaml
#! schema.yaml

#@data/values-schema
#@schema/desc "OpenAPIv3 Schema for example"
---
#@schema/desc "Configuration for example"
example:
  #@schema/desc "The namespace in which to deploy example"
  namespace: example-ns
```

### Bootstrap Command

Create a minimal schema.yaml file with the following command.

```shell
cat <<EOF > bundle/config/schema.yaml
#! schema.yaml

#@data/values-schema
#@schema/desc "OpenAPIv3 Schema for example"
---
#@schema/desc "Configuration for example"
example:
  #@schema/desc "The namespace in which to deploy example"
  namespace: example-ns
EOF
```

## Default Values

After declaring the configurable parameters in the schema, a `values.yaml` file should be created that contains the default values for all required fields.

### Sample values.yaml

```yaml
#@data/values
---

#! The namespace in which to deploy example.
namespace: example-ns
```

### Bootstrap Command

```shell
cat <<EOF > bundle/config/values.yaml
#@data/values
---

#! The namespace
namespace: example-ns
EOF
```
