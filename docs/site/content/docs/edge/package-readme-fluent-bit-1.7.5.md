# Fluent Bit

> Fluent Bit is an open source Log Processor and Forwarder which allows you to collect any data like metrics and logs from different sources, enrich them with filters and send them to multiple destinations.

This package deploys Fluent Bit as a DaemonSet, with one Fluent Bit `pod` on each Kubernetes `node`, collecting and forwarding logs from each `pod` running on that `node`.
The provided default configuration will print forwarded logs to `stdout` on the `fluent-bit`.
The `stdout` output plugin is useful for testing, but typically not appropriate for production use.
Instead, Fluent Bit supports configurations for a number of different output types, such as open-source software like `postgresql` or `loki`, cloud platform-provided services like `CloudWatch`, `S3`, or `Stackdriver`, or generic protocols like `http` and `syslog`.

For more information on output types, see the [Output](https://docs.fluentbit.io/manual/pipeline/outputs) section of the Fluent Bit Documentation.

## Supported Providers

The following table shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ✅  | ✅  | ✅ |

## Configuration

The following configuration values can be set to customize the Fluent Bit installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy Fluent Bit. |

### Fluent Bit Configuration

Fluent-bit's primary configuration interface is its config file, which is documented on Fluent's [documentation](https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/configuration-file) page.

In order to ensure that any supported inputs, outputs, filters, parsers, or other capabilities of the deployed version
of Fluent Bit are available, the addon's configuration is intentionally a lightweight pass-through of the Fluent Bit config file format.

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
|`fluent_bit.config.service`|Optional|custom configuration for Fluent Bit service, as a multiline `yaml` string.|
|`fluent_bit.config.outputs`|Optional|configuration for Fluent Bit outputs, as a multiline `yaml` string.|
|`fluent_bit.config.inputs`|Optional|configuration for Fluent Bit inputs, as a multiline `yaml` string.|
|`fluent_bit.config.filters`|Optional|configuration for Fluent Bit filters, as a multiline `yaml` string.|
|`fluent_bit.config.parsers`|Optional|configuration for Fluent Bit parsers, as a multiline `yaml` string.|
|`fluent_bit.config.streams`|Optional|content for Fluent Bit streams file, as a multiline `yaml` string.|
|`fluent_bit.config.plugins`|Optional|content for a Fluent Bit plugins configuration file, as a multiline `yaml` string.|
|`fluent_bit.daemonset.resources`|Optional|resource limits and/or requests for the `fluent-bit` container within the daemonset's pods  |
|`fluent_bit.daemonset.podAnnotations`|Optional|metadata annotations for the daemonset pods|
|`fluent_bit.daemonset.podLabels`|Optional|metadata labels for the daemonset pods|

## Usage Example

Once the package has been deployed, tail the logs from the Fluent Bit DaemonSet using `kubectl logs daemonset/fluent-bit -n fluent-bit`. You should see a large volume of logs from all pods.
