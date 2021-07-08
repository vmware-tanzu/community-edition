# tanzu-addons-manager

This package provides addons lifecycle management capabilities on tanzu clusters.

## Components

* tanzu-addons-manager

## Configuration

The following configuration values can be set to customize the tanzu-addons-manager installation.

### Global

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `tanzuAddonsManager.namespace` | Optional | the namespace in which to deploy tanzu-addons-manager. |
| `tanzuAddonsManager.createNamespace` | Optional | boolean flag to indicate whether to create namespace or not. |
| `tanzuAddonsManager.deployment.hostNetwork` | Optional | boolean flag to indicate whether to deploy in hostnetwork or not. |
| `tanzuAddonsManager.deployment.priorityClassName` | Optional | priority class name for tanzu-addons-manager deployment. |
| `tanzuAddonsManager.deployment.tolerations` | Optional | tolerations for tanzu-addons-manager deployment. |
