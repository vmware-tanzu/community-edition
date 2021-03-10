# Fluent Bit Addon

> Fluent Bit is an open source Log Processor and Forwarder which allows you to collect any data like metrics and logs from different sources, enrich them with filters and send them to multiple destinations.

This addon deploys Fluent Bit as a DaemonSet, with one Fluent Bit `pod` on each Kubernetes `node`, collecting logs from `systemd` and from each `pod` on the `node`.
The provided default configuration uses the `stdout` output plugin to forward `systemd` logs to the `fluent-bit` pod for illustrative purposes.
The `stdout` output plugin is useful for testing, but typically not appropriate for production use.
Instead, Fluent Bit supports configurations for a number of different output types, such as open-source software like `postgresql` or `loki`, cloud platform-provided services like `CloudWatch`, `S3`, or `Stackdriver`, or generic protocols like `http` and `syslog`.

For more information on output types, see the [Output](https://docs.fluentbit.io/manual/pipeline/outputs) section of the Fluent Bit Documentation.

## Configuration

The following configuration values can be set to customize the Fluent Bit installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy fluent-bit. |

### Fluent Bit Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
|`fluent_bit.outputs`|Optional|configuration for fluent-bit outputs, as a multiline `yaml` string. See the fluent-bit config file [documentation](https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/configuration-file#config_output) for detailed configuration guidance.|

## Usage Example

Once the addon has been deployed, tail the logs from the fluent-bit DaemonSet using `kubectl logs daemonset/fluent-bit -n fluent-bit`. You should see a handful of log messages from `systemd`.
