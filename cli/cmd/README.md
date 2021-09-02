# Developing Command Line Plugins for Tanzu Community Edition

Developing plugins for the `tanzu` command line utility is relatively straightforward.
However, there are a few conventions to follow in order to make the process even easier.

## Directory Structure and File Layout

When creating a plugin, it is important to use the standard layout expected by the [`tanzu builder`][builder] plugin.
That layout looks like the following:

```shell
my-plugin
  |-- pkg/               # implementation of the CLI plugin
  |   | foo-command.go   # and it's various commands as a single go package.
  |   | bar-command.go
  |
  |-- other-things/      # any other Go packages or shared libraries
  |   | secret.go        # (or scripts, tooling, etc) that should be
  |                      # separate from the core cobra command package.
  |
  | README.md            # Documentation, required by `tanzu builder`.
  | main.go              # The main cobra interface. Adds commands from "pkg"
  | go.mod               # A single go module at the root of the plugin directory
  | go.sum
  | Makefile             # It is expected that each go module has it's own Makefile
  |                      # with common, expected targets for building, testing, and running
  | 
  |-- test/              # Test files for your CLI plugin. A valid Go package is required.
  |   test.go
```

## The `pkg` Directory

Within the `pkg` directory, the [Cobra][cobra] sub-commands for the plugin should be defined.
Name each file for the command it contains.

## The `main.go` File

`main.go` imports the Cobra sub-commands from the `pkg` package, then adds them to a Tanzu `PluginDescriptor`, allowing them to be discovered by the `tanzu` CLI tool.

## The `go.mod` File

In order to isolate plugin dependencies and avoid conflicts, each plugin must define it's own `go.mod` file.

## `gopls` and Editor Support

Since each plugin has its own `go.mod` file, editor sessions should start at the top level directory in order for `gopls` and other tools to properly work.
Separating plugins into their own git repositories makes this trivial, though exceptions exist in the Tanzu Community Edition repository.

## The `Makefile`

Each go module is expected to have a `Makefile` that provides common, expected targets for building, testing, and running modules.
This enables individual plugin authors to define how their module should be ran and built while still being accessible from
the top level Makefile.
Please refer to the [`CONTRIBUTING.md` section on Makefiles](../../CONTRIBUTING.md#nested-makefiles) for reference of expected targets.

## Code Re-use

Plugins should _not_ import the `pkg` package from another plugin.
Other directories/packages may be shared between plugins, though code used between multiple plugins should be factored out into shared libraries.

## Code Repositories

Plugins should live in their own repository whenever possible.
Some plugins live in the main Tanzu Community Edition repository, but these are exceptions and not the standard.

[builder]: https://github.com/vmware-tanzu/tanzu-framework/tree/main/cmd/cli/plugin-admin/builder
[cobra]: https://github.com/spf13/cobra
