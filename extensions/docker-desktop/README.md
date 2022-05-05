# Tanzu Community Edition extension for Docker Desktop

This provides an extension integration with Docker Desktop to allow creating
a local Tanzu Community Edition unmanaged cluster quickly and easily through
the Docker Desktop interface.

## Prerequisites

In order to run this extension, you must have Docker Desktop 4.8.0 or later
installed.

If you would like to contribute or modify the extension, there are additional
requirements. This extension is comprised of Go and React JavaScript code.
Building the extension can be done using containerized build tools, but you may
want to install development environments for these tools.

Runtime Requirements:

- [Docker Desktop 4.8.0 or later](https://www.docker.com/products/docker-desktop/)

Development Recommendations:

- [Go programming language](https://go.dev/doc/install)
- [React reference](https://reactjs.org)
- [Docker Extensions CLI](https://github.com/docker/extensions-sdk)

## Building and Installing

The standard way to get the Tanzu Community Edition extension for Docker
Desktop is by using the Docker Marketplace. This will install the official
released version of the extension.

If you are making local changes and would like to try them out, you will need
to follow these steps:

1. In Docker Desktop, go to Preferences > Extensions and make sure
   "Enable Docker Extensions" is checked.
1. From a terminal, navigate to `$TCE_REPO/extensions/docker-desktop`.
1. Run the following commands to build and install the local extension:

   ```sh
   make extension
   make install
   ```

1. From the Docker Dashboard you can now navigate to the Extensions section.
   It should now list *Tanzu Community Edition* as one of the available
   extensions. Click on *Tanzu Community Edition* from the list and you should
   be presented with the UI for creating and deleting clusters.

### `Docker Extension` CLI Setup

Note: The build steps assume that the Docker Extensions CLI has been installed.
While `docker-extension` can be called directly, the installation target
assumes it has been added as a CLI plugin and can be called as
`docker extension`.

If you have downloaded the `docker-extension` binary from their Releases page,
follow these steps to have it recognized as a CLI plugin under `docker`:

```sh
mkdir -p ~/.docker/cli-plugins
cp docker-extension ~/.docker/cli-plugins/
```
