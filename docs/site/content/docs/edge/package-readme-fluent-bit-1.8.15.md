# Fluent Bit Package

This package forward and collector log using [fluent-bit](https://fluentbit.io/).

## Installation

### Installation of package

Install TCE Fluent-bit package through tanzu command:

```bash
tanzu package install fluent-bit --package-name fluent-bit.community.tanzu.vmware.com --version ${FLUENT-BIT_PACKAGE_VERSION}
```

> You can get the `${FLUENT-BIT_PACKAGE_VERSION}` from running `tanzu package
> available list fluent-bit.community.tanzu.vmware.com`. Specifying a
> namespace may be required depending on where your package repository was
> installed.

## Options

The following configuration values can be set to customize the Prometheus installation.

### Package configuration values

#### Global

| Value | Required/Optional | Default | Description |
|-------|-------------------|---------|-------------|
| `namespace` | Optional | `fluent-bit` | The namespace in which to deploy Fluent Bit. |

### Fluent Bit Configuration

> Fluent-bit's primary configuration interface is its config file, which is documented on Fluent's [documentation](https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/configuration-file) page.
> In order to ensure that any supported inputs, outputs, filters, parsers, or other capabilities of the deployed version
of Fluent Bit are available, the addon's configuration is intentionally a lightweight pass-through of the Fluent Bit config file format.

| Value | Required/Optional | Default | Description |
|-------|-------------------|---------|-------------|
|`fluent_bit.config.service`|Optional|fluent-bit.conf|custom configuration for Fluent Bit service, as a multiline `yaml` string.|
|`fluent_bit.config.outputs`|Optional|outputs.conf|configuration for Fluent Bit outputs, as a multiline `yaml` string.|
|`fluent_bit.config.inputs`|Optional|inputs.conf|configuration for Fluent Bit inputs, as a multiline `yaml` string.|
|`fluent_bit.config.filters`|Optional|filters.conf|configuration for Fluent Bit filters, as a multiline `yaml` string.|
|`fluent_bit.config.parsers`|Optional|parsers.conf|configuration for Fluent Bit parsers, as a multiline `yaml` string.|
|`fluent_bit.config.streams`|Optional|streams.conf|content for Fluent Bit streams file, as a multiline `yaml` string.|
|`fluent_bit.config.plugins`|Optional|plugins.conf|content for a Fluent Bit plugins configuration file, as a multiline `yaml` string.|
|`fluent_bit.daemonset.resources`|Optional|{}|resource limits and/or requests for the `fluent-bit` container within the daemonset's pods |
|`fluent_bit.daemonset.podAnnotations`|Optional|{}|metadata annotations for the daemonset pods|
|`fluent_bit.daemonset.podLabels`|Optional|{}|metadata labels for the daemonset pods|
|`fluent_bit.daemonset.env`|Optional|[]|Environment variables for the daemonset pods|

### Application configuration values

No available options to configure.

#### Multi-cloud configuration steps

There are currently no configuration steps necessary for installation of the Grafana package to any provider.

## What This Package Does

This package deploys Fluent Bit as a DaemonSet, with one Fluent Bit `pod` on each Kubernetes `node`, collecting and forwarding logs from each `pod` running on that `node`.
The provided default configuration will print forwarded logs to `stdout` on the `fluent-bit`.
The `stdout` output plugin is useful for testing, but typically not appropriate for production use.
Instead, Fluent Bit supports configurations for a number of different output types, such as open-source software like `postgresql` or `loki`, cloud platform-provided services like `CloudWatch`, `S3`, or `Stackdriver`, or generic protocols like `http` and `syslog`.

For more information on output types, see the [Output](https://docs.fluentbit.io/manual/pipeline/outputs) section of the Fluent Bit Documentation.

## Components

- Fluent-bit

## Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Files

Here is an example of the package configuration file [values.yaml](bundle/config/values.yaml).

## Package Limitations

- This package currently does not support users to add custom paths to collect logs.

## Usage Example

- Once the package has been deployed, tail the logs from the Fluent Bit DaemonSet using `kubectl logs daemonset/fluent-bit -n fluent-bit`. You should see a large volume of logs from all pods.

## Troubleshooting

Not applicable.

## Additional Documentation

See the [documentation](https://docs.fluentbit.io/manual/v/1.8/) for more information.
