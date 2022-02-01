# Documentation

This should include a brief overview of the software components contained in the package, a description of configuration parameters, and general usage information. This documentation is not intended to replace, or be as extensive as the official documentation for the software.

The package documentation should highlight dependencies or considerations on other packages, software, Kubernetes distributions, or underlying infrastructure (e.g. AWS, GCP, Docker, vSphere, etc).

## Sample README

```text
# Example Package

This package provides << awesome functionality >> using [example](https://example.com).

## Supported Providers

The following tables shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ❌  | ❌  | ❌  |

## Components

* Example version: `1.0.0`

## Configuration

The following configuration values can be set to customize the Example package installation.

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `example.namespace` | Optional | namespace to install Example to |

## Usage Example

The following is a basic guide for getting started with the Example package.

Step 1...
```

## Bootstrap Command

```shell
cat <<EOF > README.md
# Example Package

This package provides << awesome functionality >> using [example](https://example.com).

## Supported Providers

The following tables shows the providers this package can work with.

| AWS  |  Azure  | vSphere  | Docker |
|:---:|:---:|:---:|:---:|
| ✅  |  ❌  | ❌  | ❌  |

## Components

* Example version: `1.0.0`

## Configuration

The following configuration values can be set to customize the Example package installation.

| Value | Required/Optional | Description |
|-------|-------------------|-------------|
| `example.namespace` | Optional | namespace to install Example to |

## Usage Example

The following is a basic guide for getting started with the Example package.
EOF
```
