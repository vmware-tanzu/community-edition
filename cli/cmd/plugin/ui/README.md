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
