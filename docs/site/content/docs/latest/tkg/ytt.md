# Customizing Clusters, Plans, and Extensions with ytt Overlays

This topic explains how you can use [ytt Overlays](https://github.com/vmware-tanzu/carvel-ytt) to customize Tanzu Kubernetes (workload) clusters, cluster plans, extensions, and shared services.

Tanzu Kubernetes Grid distributes configuration files for clusters and cluster plans in the `~/.tanzu/tkg/providers/` directory, and for extensions and shared services in the Tanzu Kubernetes Grid [Extensions Bundle](extensions/index.md#unpack-bundle).

You can customize these configurations by adding or modifying configuration files directly, or by using [`ytt`](https://github.com/vmware-tanzu/carvel-ytt) overlays.

Directly customizing configuration files is simpler, but if you are comfortable with `ytt` overlays, they let you customize configurations at different scopes and manage multiple, modular configuration files, without destructively editing upstream and inherited configuration values.

For more information about `ytt`, including overlay examples and an interactive validator tool, see:

- Carvel Tools > `ytt` > [Interactive Playground](https://carvel.dev/ytt/#example:example-overlay-files)
- The IT Hollow: [Using YTT to Customize TKG Deployments](https://theithollow.com/2020/11/09/using-ytt-to-customize-tkg-deployments/)

## <a id="clusters-plans"></a> Clusters and Cluster Plans

The `~/.tanzu/tkg/providers/` directory includes `ytt` directories and  `overlay.yaml` files at different levels, which lets you scope configuration settings at each level:

- Provider- and version-specific `ytt` directories. For example, `~/.tanzu/tkg/providers/infrastructure-aws/v0.6.4/ytt`.
   - Specific configurations for provider API version.
   - The `base-template.yaml` file contains all-caps placeholders such as `"${CLUSTER_NAME}"` and should not be edited.
   - The `overlay.yaml` file is tailored to overlay values into `base-template.yaml`.
- Provider-wide `ytt` directories. For example, `~/.tanzu/tkg/providers/infrastructure-aws/ytt`.
   - Provider-wide configurations that apply to all versions.
- Top-level `ytt` directory, `~/.tanzu/tkg/providers/ytt`.
   - Cross-provider configurations.
   - Organized into numbered directories, and processed in number order.
   - `ytt` traverses directories in alphabetical order and overwrites duplicate settings, so you can create a `/04_user_customizations` subdirectory for configurations that take precedence over any in lower-numbered `ytt` subdirectories.

**IMPORTANT**: You can only use `ytt` overlays to modify workload clusters. Using `ytt` overlays to modify management clusters is not supported.

### <a id="cluster-plan-examples"></a> Cluster and Plan Examples

Examples of `ytt` overlays for customizing workload clusters and cluster plans include:

- [Nameservers on vSphere](tanzu-k8s-clusters/config-plans.md#nameserver)
- [Privileged Containers for Workloads](tanzu-k8s-clusters/config-plans.md#privileged)
- [Trust Custom CA Certificates on Cluster Nodes](cluster-lifecycle/secrets.md#custom-ca)
- [Disable Bastion Server on AWS](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/capi-overrides/disable-aws-bastion.yaml) in the _TKG Lab_ repository
- [New `nginx` Workload Plan](tanzu-k8s-clusters/config-plans.md#nginx)

## <a id="extensions"></a> Extensions and Shared Services

The Tanzu Kubernetes Grid [Extensions Bundle](extensions/index.md#unpack-bundle) includes templates for `ytt` overlays to implement various customizations.

These `ytt` overlay templates are in `overlay` subdirectories in the following locations:

- [Ingress Control with Contour](extensions/ingress-contour.md): `ingress/contour/overlays`
- [Harbor Registry](extensions/harbor-registry.md): `registry/harbor/overlays`
- [Monitoring with Prometheus and Grafana](extensions/monitoring.md): `monitoring/prometheus/overlays` and  `monitoring/grafana/overlays`
- [Log Forwarding with Fluent Big](extensions/logging-fluentbit.md): `logging/fluent-bit/overlays`
- (Deprecated) User Authentication with Dex and Gangway: `authentication/dex/overlays` and `authentication/gangway/overlays`

Before deploying an extension, you can use these overlay templates or create and apply your own overlays to the extension as follows:

1. In the extensions bundle directory, under the extension's `/overlays` directory, modify an overlay template or create a new overlay to contain your custom values:

  - **Existing template**: Find and modify the template that fits your need. The template filenames indicate their use; for example `change-namespace.yaml` and `update-registry.yaml`.
  - **New overlay**: Create a new `ytt` overlay file. For example, to add an annotation to the Grafana extension's HTTP Proxy, create a `update-grafana-httproxy.yaml` overlay with contents:

      ```
      #@ load("@ytt:overlay", "overlay")
      #@ load("@ytt:yaml", "yaml")
      #@overlay/match by=overlay.subset({"kind": "HTTPProxy", "metadata": {"name": "grafana-httpproxy"}})
      ---
      metadata:
        #@overlay/match missing_ok=True
        annotations:
          #@overlay/match missing_ok=True
          dns.alpha.kubernetes.io/hostname: grafana.tkg.vclass.local
      ```

1. Save the overlay content as a secret in the extension's namespace. For example with `update-grafana-httproxy.yaml` above, run:

  ```
  kubectl create secret generic grafana-httpproxy --from-file=update-grafana-httpproxy.yaml=update-grafana-httpproxy.yaml -n tanzu-system-monitoring
  ```

1. In the extension's deployment file, under `extensions/` add a reference to the new secret. For example for the `grafana-httproxy` secret above, add the following to the file `/extensions/monitoring/grafana/grafana-extension.yaml` under `spec.template.ytt.inline.pathsFrom`, after the existing `grafana-data-values` setting:

  ```
          - secretRef:
              name: grafana-httpproxy
  ```

The examples below show some specific use cases for creating and applying custom overlays.

### <a id="extension-examples"></a> Extension and Shared Service Examples

Examples of applying `ytt` overlay files for customizing extensions and shared services include:

- Contour: [External DNS Annotation](extensions/ingress-contour.md#ytt)
- Harbor: [Clean Up S3 and Trust Let's Encrypt](extensions/harbor-registry.md#ytt)

For more examples, see the [TKG Lab repository](https://github.com/Tanzu-Solutions-Engineering/tkg-lab) and its [Step by Step setup guide](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/docs/baseline-lab-setup/step-by-step.md).
