# Package for Kubeapps

Package for Kubeapps will house the carvel tooling specific configuration and templating for a deployment of [Kubeapps](https://github.com/vmware-tanzu/kubeapps) that will be leveraged by VMware Tanzu Community Edition.

Kubeapps is an in-cluster web-based application that enables users with a one-time installation to deploy, manage, and upgrade applications on a Kubernetes cluster.

See the [Kubeapps project README](https://github.com/vmware-tanzu/kubeapps) for more details.

## Components

The Kubeapps application is itself comprised of a number of smaller components:

- The Kubeapps dashboard is the user-interface that runs in the browser.
- The Kubeapps APIs service is the backend which serves requests for the user interface.
- A number of other components, such as nginx, Postgres, oauth2-proxy and Redis are used depending on the configuration.

See [Kubeapps Components](https://github.com/vmware-tanzu/kubeapps/tree/main/docs/reference/developer) in our main documentation for more information.

## Supported Providers

The following table shows the providers this package can work with.

| AWS  | Azure | vSphere | Docker |
|------|-------|---------|--------|
| ✅   | ✅    | ✅      | ✅     |

Although please note that currently Kubeapps can only be run on TCE with token authentication which is appropriate for demonstration purposes only. For more information, please see the relevant [Contour issue #4290](https://github.com/projectcontour/contour/issues/4290)

## Configuration

The configuration for the Kubeapps Carvel package is currently identical to the related Bitnami Helm chart. Please refer to the [configuration options in the Chart readme](https://github.com/vmware-tanzu/kubeapps/tree/main/chart/kubeapps).

Although the configuration options are identical, with TCE the environment into which Kubeapps is installed is different. In particular, when TCE is installed with Contour, certain functionality of Kubeapps is not currently possible. In this environment, Kubeapps can only be used with service-account token authentication, which is suitable for demonstration purposes only. The recommended OpenIDConnect authentication for Kubeapps is not currently possible when using Contour until the fix for [Contour issue #4920](https://github.com/projectcontour/contour/issues/4290) is released.

When running Kubeapps on a cluster with Contour installed, it is possible to use Kubeapps with token authentication together with a [required Contour `HTTPProxy` custom resource](https://github.com/vmware-tanzu/kubeapps/issues/3716#issuecomment-1067532124) that ensures the requests to the Kubeapps backend are routed correctly.

## Installation

   ```shell
   tanzu package install kubeapps \
      --package-name kubeapps.community.tanzu.vmware.com \
      --version ${KUBEAPPS_PACKAGE_VERSION} \
      --values-file my-values.yaml
   ```

   > You can get the `${KUBEAPPS_PACKAGE_VERSION}` by running `tanzu
   > package available list kubeapps.community.tanzu.vmware.com`.
   > Specifying a namespace may be required depending on where your package
   > repository was installed.

## Documentation

For Kubeapps-specific documentation, check out
our the main repository
[vmware-tanzu/kubeapps](https://github.com/vmware-tanzu/kubeapps).

## Contributing

The Package for Kubeapps project team welcomes contributions from the community. If you wish to contribute code and you have not signed our contributor license agreement (CLA), our bot will update the issue when you open a Pull Request. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

See the [Apache License](./LICENSE)
