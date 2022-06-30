# UI

`ui` is a proof-of-concept for trying out user interface ideas for Tanzu
Community Edition. This is not an official project within TCE and anything in
this tree may be removed completely when done.

## Building

To build the UI plugin, from the root of this repo run the `make build-tce-cli-plugins`
command. This will produce a binary in
`artifacts/ui/${VERSION}/tanzu-ui-${OS}_${ARCH}`.

To speed up build time to only build the plugin for your local environment,
run `make build-tce-cli-plugins ENVS=darwin-amd64`, replacing the `ENVS` value
with the appropriate `${GOHOSTOS}-${GOHOSTARCH}` for your environment.

This plugin can then either be "installed" to be used with the root `tanzu`
command, or you can run it directly by calling
`./artifacts/ui/v0.10.0-dev.5/tanzu-ui-darwin_amd64`.

## REST API

The REST API uses OpenAPI to define its API endpoints and objects. The OpenAPI
specification can be found in the `api/swagger.yaml` file.

The server stubs are generated using the go-swagger tool. Installation and links
to more detailed documentation can be found
[here](https://goswagger.io/install.html).

For convenience, a `make` target is available that uses a containerized
go-swagger to generate the server stubs:

```sh
make generate-ui-swagger-api
```

## UI Development

To run the front end using the mock server (for faster UI development), clone the repo and follow these steps:
1. Install dependencies: in a shell, cd to `cli/cmd/plugin/ui/web/tanzu-ui` and run `npm ci`
2. Run the mock server: cd to `cli/cmd/plugin/ui/web/tanzu-ui/node-server` and run `npm run start`
3. Run the UI: in a separate shell, cd to `cli/cmd/plugin/ui/web/tanzu-ui` and run `npm run start`. This should bring up a browser window with the UI.
4. To understand what credentials the mock server is expecting (for mock-connecting to a cloud provider), look in
   `cli/cmd/plugin/ui/web/tanzu-ui/node-server/src/routes/api/endpoints/[provider].js`. This file will also contain any logic surrounding other mock endpoints for the given
   provider.
