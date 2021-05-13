# Metrics Server

> Metrics Server is a scalable, efficient source of container resource metrics for Kubernetes built-in autoscaling pipelines.

This package deploys Metrics Server with insecure connection between API server and Kubelet by default.

For more information, see the [Github page](https://github.com/kubernetes-sigs/metrics-server) of Metrics Server.

## Configuration

The following configuration values can be set to customize the Metrics Server installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy resources. Default: kube-system |

### Metrics Server Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `metricsServer.createNamespace` | Optional | A boolean that indicates whether to create the namespace specified. Default value is `true`. |
| `metricsServer.config.updateStrategy` | Optional | The update strategy of the deployment. Default: `RollingUpdate` |
| `metricsServer.config.probe.failureThreshold` | Optional | Probe failure threshold. Default: `3`. |
| `metricsServer.config.probe.periodSeconds` | Optional | Probe period. Default: `10` . |
| `metricsServer.config.nodeSelector.key` | Optional | Select which node should Metrics-server pod runs on. Default: `null`. |
| `metricsServer.config.nodeSelector.value` | Optional | Select which node should Metrics-server pod runs on. Default: `null`. |
| `metricsServer.config.apiServiceInsecureTLS`| Optional | Insecure connection between API service. Default: `True`. |

## Usage Example

Once the package has been deployed, get the metrics of all pods `kubectl top pod -A`. You should see a the CPU and memory usage info of all pods.
