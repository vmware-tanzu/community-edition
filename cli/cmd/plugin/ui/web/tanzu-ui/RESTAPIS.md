# REST API Consumption

This React application consumes a set of REST APIs as defined by `/community-edition/cli/cmd/plugin/ui/api/swagger.yaml`

A set of typescript services and models are auto-generated from the swagger spec and can be found in
`/community-edition/cli/cmd/plugin/ui/web/tanzu-ui/src/swagger-api`

The methods found in the `services` folder should be used when making any REST API calls to the backend.

It is recommended that the types defined in the `models` folder are leveraged for ensuring integrity of
data when sending POST API payloads, or retrieving GET API responses.

## Generating the swagger-api folder and contents

The swagger API generated methods and types are created using the `openapi-typescript-codegen` library.

An npm script (generate-api) is defined in `/community-edition/cli/cmd/plugin/ui/web/tanzu-ui/package.json`
and is run as part of the `npm run start` script. However, you can also run `npm run generate-api` manually
at any time to re-generate the REST API types and methods.
