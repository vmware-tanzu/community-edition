# kapp-controller Package

kapp-controller is a Kubernetes controller for installing configuration (helm charts, ytt templates, plain yaml) with kapp as described by App CRD

This package deploys [kapp-controller](https://carvel.dev/kapp-controller/).

## Components

## Configuration

The following configuration values can be set to customize the kapp-controller installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `imageInfo.imageRepository` | Optional | The repository from which to fetch kapp-controller image. |
| `imageInfo.imagePullPolicy` | Optional | The image pull policy. |
| `imageInfo.images.kappControllerImage.imagePath` | Optional | The path of kapp-controler image. |
| `imageInfo.images.kappControllerImage.tag` | Optional | The tag of kapp-controller image. |

### kapp-controller Configuration

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `kappController.namespace` | Optional | Namespace in which to deploy kapp-controller. |
| `kappController.createNamespace` | Optional | Boolean flag to create kapp-controller namespace or not. |
| `kappController.deployment.hostNetwork` | Optional | Boolean flag to indicate if kapp-controller needs to be deployed in host network. |
| `kappController.deployment.priorityClassName` | Optional | Priority class name of kapp-controller deployment. |
| `kappController.deployment.concurrency` | Optional | Concurrency of kapp-controller. |
| `kappController.config.caCerts` | Optional | CA certificates of image registry from where kapp-controller fetches images. |
| `kappController.config.httpProxy` | Optional | HTTP proxy to use. |
| `kappController.config.httpsProxy` | Optonal | HTTPS proxy to use. |
| `kappController.config.noProxy` | Optonal | no proxy to use. |
| `kappController.config.dangerousSkipTLSVerify` | Optional | skip tls verify. |

## Usage Example

For kapp-controller usage refer to [kapp-controller](https://carvel.dev/kapp-controller/).
