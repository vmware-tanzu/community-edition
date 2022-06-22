# Metrics Server

> Metrics Server is a scalable, efficient source of container resource metrics for Kubernetes built-in autoscaling pipelines.

This package deploys Metrics Server with insecure connection between API server and Kubelet by default.

For more information, see the [GitHub page](https://github.com/kubernetes-sigs/metrics-server) of Metrics Server.

## Configuration

The following configuration values can be set to customize the Metrics Server installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which metrics-server is deployed. Default: kube-system |

### Metrics Server Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `metricsServer.createNamespace` | Optional | Whether to create namespace specified for metrics-server. Default value is `true`. |
| `metricsServer.namespace` | Optional | The namespace value used by older templates, will be overwriten if top level namespace is present, kept for backward compatibility. Default value is `null`. |
| `metricsServer.config.securePort` | Optional | TThe HTTPS secure port used by metrics-server. Default: `4443`. |
| `metricsServer.config.updateStrategy` | Optional | TThe update strategy of the metrics-server deployment. Default: `RollingUpdate` |
| `metricsServer.config.probe.failureThreshold` | Optional | Probe failureThreshold of metrics-server deployment. Default: `3`. |
| `metricsServer.config.probe.periodSeconds` | Optional | Probe period of metrics-server deployment. Default: `10` . |
| `metricsServer.config.apiServiceInsecureTLS`| Optional | Whether to enable insecure TLS for metrics-server api service. Default: `True`. |

## Usage Example

Once the package has been deployed, get the metrics of all pods `kubectl top pod -A`. You should see a the CPU and memory usage info of all pods.
