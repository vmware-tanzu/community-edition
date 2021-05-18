# Kapp Controller Package

This package provides app/package lifecycle management using [kapp-controller](https://carvel.dev/kapp-controller/).

## Components

* kapp-controller

## Configuration

The following configuration values can be set to customize the kapp-controller installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `namespace` | Optional | The namespace in which to deploy kapp-controller. |

### Kapp Controller Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `kappController.createNamespace` | Optional | A boolean that indicates whether to create the namespace specified. Default value is `true`. |
| `kappController.namespace` | Optional | The namespace value used by older templates, will overwrite will top level namespace of present, keep for backward compatibility. Default value is `null`. |
| `kappController.deployment.hostNetwork` | Optional | HostNetwork of kapp-controller deployment. Default is `null`. |
| `kappController.deployment.priorityClassName` | Optional | priorityClassName of kapp-controller deployment. Default value is `null`. |
| `kappController.deployment.concurrency` | Optional | concurrency of kapp-controller deployment. Default is `2`. |
| `kappController.config.caCerts` | Optional | A cert chain of trusted ca certs. These will be added to the system-wide cert pool of trusted ca's. Default value is `""`. |
| `kappController.config.httpProxy` | Optional | The url/ip of a proxy for kapp controller to use when making network requests. Default is `""`. |
| `kappController.config.httpsProxy` | Optional | The url/ip of a tls capable proxy for kapp controller to use when making network requests. Default value is `""`. |
| `kappController.config.noProxy` | Optional | A comma delimited list of domain names which kapp controller should bypass the proxy for when making requests. Default is `""`. |
| `kappController.config.dangerousSkipTLSVerify` | Optional | A comma delimited list of hostnames for which kapp controller should skip TLS verification. Default value is `""`. |

## Usage Example

To learn more about how to use kapp-controller refer to [kapp-controller website](https://carvel.dev/kapp-controller/)
