# Documentation

This should include a brief overview of the software components contained in the package, a description of configuration parameters, and general usage information. This documentation is not intended to replace, or be as extensive as the official documentation for the software.

The package documentation should highlight dependencies or considerations on other packages, software, Kubernetes distributions, or underlying infrastructure (e.g. AWS, GCP, Docker, vSphere, etc).

## Sample README

```text
# Tanzu Community Edition Package Template

Every TCE package is required to follow this documentation so that we can ensure a consistent user experience. If you’re a package owner, please ensure your documentation follows this structure.

## Name

The name of this package page.

## Installation

Briefly describes deploying the package.

### Installation of dependencies

### Installation of package

## Options

### Package configuration values

Description of values that can be configured for a package and how they change its behavior. This section can be blank but must specify “No available options to configure” if so. Otherwise, the Files section (see below) must include a working example values file with required values. This section must include table listing configuration options including Value, Required/Optional, Default, Description.  

### Application configuration values

Description of values that when configured pass through native configuration options to the installed application. This section must include sub-sections with title as the value name and the following format:

#### <Value>

Description: include a description of why you would use this native configuration option and reference to configurable options in application documentation
Required: true/false
Default:
Example:

#### Multi-cloud configuration steps

## What This Package Does

Gives an explanation of what the package does. Describe the usual case. For information on options of the package use the options section.

## Components

Describe the version of this package, dependencies, and supported providers.

### Supported Providers

## Files

Lists the files or paths the package directly operates on. If a package describes the `values-file` option, the #FILE section should provide an example of this file. The file's contents are yaml and reference the configuration parameters. If not relevant, put: “Not relevant for this particular package.”

## Package Limitations

Lists the limitations, known defects or inconveniences, and other questionable activities. Include the link for where to report bugs.

## Usage Example

Provides one or more examples describing how this package is used.

## Troubleshooting

Provides information and steps to support investigation and reporting of issues.

## Additional Documentation
...
```

## Bootstrap Command

```shell
cat <<EOF > README.md
# Tanzu Community Edition Package Template

Every TCE package is required to follow this documentation so that we can ensure a consistent user experience. If you’re a package owner, please ensure your documentation follows this structure.

## Name

The name of this package page.

## Installation

Briefly describes deploying the package.

### Installation of dependencies

### Installation of package

## Options

### Package configuration values

Description of values that can be configured for a package and how they change its behavior. This section can be blank but must specify “No available options to configure” if so. Otherwise, the Files section (see below) must include a working example values file with required values. This section must include table listing configuration options including Value, Required/Optional, Default, Description.  

### Application configuration values

Description of values that when configured pass through native configuration options to the installed application. This section must include sub-sections with title as the value name and the following format:

#### <Value>

Description: include a description of why you would use this native configuration option and reference to configurable options in application documentation
Required: true/false
Default:
Example:

#### Multi-cloud configuration steps

## What This Package Does

Gives an explanation of what the package does. Describe the usual case. For information on options of the package use the options section.

## Components

Describe the version of this package, dependencies, and supported providers.

### Supported Providers

## Files

Lists the files or paths the package directly operates on. If a package describes the `values-file` option, the #FILE section should provide an example of this file. The file's contents are yaml and reference the configuration parameters. If not relevant, put: “Not relevant for this particular package.”

## Package Limitations

Lists the limitations, known defects or inconveniences, and other questionable activities. Include the link for where to report bugs.

## Usage Example

Provides one or more examples describing how this package is used.

## Troubleshooting

Provides information and steps to support investigation and reporting of issues.

## Additional Documentation
...
EOF
```
