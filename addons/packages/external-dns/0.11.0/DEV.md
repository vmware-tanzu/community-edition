# Development

## Generate package `valuesSchema` from ytt schema

1. Generate OpenAPIv3 schema from `schema.yaml` using `ytt`.

```bash
ytt -f bundle/config/schema.yaml --data-values-schema-inspect --output openapi-v3
```

1. Copy contents of `components.schemas.dataValues` into
`spec.valuesSchema.openAPIv3` of `package.yaml`.
