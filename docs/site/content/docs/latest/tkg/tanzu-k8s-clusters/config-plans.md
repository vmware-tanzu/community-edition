# Configure Tanzu Kubernetes Plans and Clusters

This topic explains where Tanzu Kubernetes (workload) cluster plan configuration values come from, and what is the order of precedence for their multiple sources.
It also explains how you can customize the `dev` and `prod` plans for workload clusters on each cloud infrastructure,
and how you can use `ytt` overlays to customize cluster plans and clusters, and create new custom plans, while preserving original configuration code.

## <a id="source"></a> Where Cluster Configuration Values Come From

When the `tanzu` CLI creates a cluster, it combines configuration values from the following:

  - Live input at invocation
     - CLI input
     - UI input, when deploying a management cluster with the installer
  - Environment variables
  - `~/.tanzu/tkg/cluster-config.yaml` or other file passed to the CLI `--file` option
  - Cluster plan YAML configuration files in `~/.tanzu/tkg/providers`, as described in [Plan Configuration Files](#plans) below.
  - Other, non-plan YAML configuration files under `~/.tanzu/tkg/providers`

Live input applies configuration values that are unique to each invocation, environment variables persist them over a terminal session, and configuration files and overlays persist them indefinitely.
You can customize clusters through any of these sources, with recommendations and caveats described below.

See [Configuration Precedence Order](#precedence) for how the `tanzu` CLI derives specific cluster configuration values from these multiple sources where they may conflict.

### <a id="plans"></a> Plan Configuration Files

The `~/.tanzu/tkg/providers` directory contains workload cluster plan configuration files in the following subdirectories, based on the cloud infrastructure that deploys the clusters:

| Clusters deployed by... | `~/.tanzu/tkg/providers` Directory  |
|-------------------------|-------------------------------|
| Management cluster on **vSphere** | `/infrastructure-vsphere` |
| **vSphere 7** Supervisor cluster | `/infrastructure-tkg-service-vsphere` |
| Management cluster on **Amazon EC2** | `/infrastructure-aws` |
| Management cluster on **Azure** | `/infrastructure-azure` |

These plan configuration files are named `cluster-template-definition-PLAN.yaml`.
The configuration values for each plan come from these files and from the files that they list under `spec.paths`:

   - Config files that ship with the `tanzu` CLI
   - Custom files that users create and add to the `spec.paths` list
   - [ytt Overlays](https://github.com/vmware-tanzu/carvel-ytt) that users create or edit to overwrite values in other configuration files

## <a id="files"></a> Files to Edit, Files to Leave Alone

To customize cluster plans via YAML, you edit files under `~/.tanzu/tkg/providers/`, but you should avoid changing other files.

**Files to Edit**

Workload cluster plan configuration file paths follow the form `~/.tanzu/tkg/providers/infrastructure-INFRASTRUCTURE/VERSION/cluster-template-definition-PLAN.yaml`, where:

   - `INFRASTRUCTURE` is `vsphere`, `aws`, or `azure`.
   - `VERSION` is the version of the Cluster API Provider module that the configuration uses.
   - `PLAN` is `dev`, `prod`, or a custom plan as created in the [New `nginx` Workload Plan](#nginx) example.

Each plan configuration file has a `spec.paths` section that lists source files and `ytt` directories that configure the cluster plan.  For example:

```
spec:
  paths:
    - path: providers/infrastructure-aws/v0.5.5/ytt
    - path: providers/infrastructure-aws/ytt
    - path: providers/ytt
    - path: bom
      filemark: text-plain
    - path: providers/config_default.yaml
```

These files are processed in the order listed.
If the same configuration field is set in multiple files, the last-processed setting is the one that the `tanzu` CLI uses.

To customize your cluster configuration, you can:

  - Create new configuration files and add them to the `spec.paths` list.
     - This is the easier method.
  - Modify existing `ytt` overlay files as described in [`ytt` Overlays](#examples) below.
     - This is the more powerful method, for people who are comfortable with `ytt`.

**Files to Leave Alone**

VMware discourages changing the following files under `~/.tanzu/tkg/providers`, except as directed:

- `base-template.yaml` files, in `ytt` directories
   - These configuration files use values from the Cluster API provider repos for vSphere, AWS, and Azure under [Kubernetes SIGs](https://github.com/kubernetes-sigs), and other upstream, open-source projects, and they are best kept intact.
   - Instead, create new configuration files or see [Clusters and Cluster Plans](../ytt.md#clusters-plans) in _Customizing Clusters, Plans, and Extensions with ytt Overlays_ below to set values in the `overlay.yaml` file in the same `ytt` directory.

- `~/.tanzu/tkg/providers/config_default.yaml` - Append only
   - This file contains system-wide defaults for Tanzu Kubernetes Grid on all cloud infrastructures.
   - Do not modify existing values in this file, but you can append a `User Customizations` section at the end, as in the [Privileged Containers for Workloads](#privileged) example below.
   - Instead of changing values in this file, customize cluster configurations in files that you pass to the `--file` option of `tanzu cluster create` and `tanzu management-cluster create`.

- `~/.tanzu/tkg/providers/config.yaml`
   - The `tanzu` CLI uses this file as a reference for all providers present in the `/providers` directory, and their default versions.

## <a id="precedence"></a> Configuration Precedence Order

When the `tanzu` CLI creates a cluster, it reads in configuration values from multiple sources that may conflict.
It resolves conflicts by using values in the following order of descending precedence:

<table>
  <tr>
    <th width="20%">Processing layers, ordered by descending precedence</td>
    <th width="20%">Source</td>
    <th width="20%">Examples</td>
    <th width="40%">Notes</td>
  </tr><tr>
    <td>8. User-specific data values, from or written to top-level config file</td>
    <td><code>AZURE_NODE_MACHINE_TYPE: Standard_D2s_v3</code></td>
    <td><p>The main source of workload (and management) cluster parameters is the file passed to CLI <code>--file</code> option, which defaults to <code>~/.tanzu/tkg/cluster-config.yaml</code>.</p></td>
  </tr>
  <tr>
    <td>7. Factory default data values</td>
    <td rowspan=2>Shipped with TKG</td>
    <td><code>config_default.yaml</code></td>
    <td>These are the supported cluster template configuration "knobs", with documentation and their default settings where applicable.</td>
  </tr><tr>
    <td>6. BOM metadata data values</td>
    <td><code>bom-1.3.1+vmware.1.yaml</code></td>
    <td>One per Kubernetes version released by TKG</td>
  </tr><tr>
    <td>5 (tie). User-provided customizations</td>
    <td rowspan=6>Customizable ytt</td>
    <td><code>myhacks.yaml</code> </td>
    <td>Topmost layer of <code>ytt</code> processing files before the Data Values layers; takes precedence over the layers below it </td>
  </tr><tr>
    <td>5 (tie). Additional processing YAMLs, not user-provided</td>
    <td><code>rm-bastion.yaml</code>, <code>rm-mhc.yaml</code>, <code>custom-resource-annotations.yaml</code></td>
    <td></td>
  </tr><tr>
    <td>4. Add-on YAMLs and overlays</td>
    <td><code>calico.yaml</code>, <code>antrea.yaml</code></td>
    <td>A specific class of customization representing one of more resources to be applied to the cluster post-creation.</td>
  </tr><tr>
    <td>3. Plan-specific processing YAMLs</td>
    <td><code>prod.yaml</code>, <code>dev.yaml</code></td>
    <td>Plan-specific customizations.</td>
  </tr><tr>
    <td>2. Overlay YAML</td>
    <td><code>ytt/overlay.yaml</code></td>
    <td>Defines what in the basic template is overridable, using legacy, <code>"KEY_NAME:value"</code> style entries.</td>
  </tr><tr>
    <td>1. Base Cluster template YAML</td>
    <td><code>ytt/base-template.yaml</code></td>
    <td>Base CAPI template with actual default values and no <code>ytt</code> annotations.</td>
  </tr>
</table>

## <a id="examples"></a> `ytt` Overlays

Tanzu Kubernetes Grid supports customizing workload cluster configurations by adding or modifying configuration files directly,
but using `ytt` overlays instead lets you customize configurations at different scopes and manage multiple, modular configuration files, without destructively editing upstream and inherited configuration values.

For more information, see [Clusters and Cluster Plans](../ytt.md#clusters-plans) in _Customizing Clusters, Plans, and Extensions with ytt Overlays_.

**IMPORTANT**: You can only use `ytt` overlays to modify workload clusters. Using `ytt` overlays to modify management clusters is not supported.

The following examples show how to use configuration overlay files to customize workload clusters and create a new cluster plan.

For an overlay that customizes cluster certificates, see [Trust Custom CA Certificates on Cluster Nodes](../cluster-lifecycle/secrets.md#custom-ca) in the _Tanzu Kubernetes Cluster Secrets_ topic.

### <a id="nameserver"></a> Nameservers on vSphere

This example adds one or more custom nameservers to worker and control plane nodes in Tanzu Kubernetes Grid clusters on vSphere.  It disables DNS resolution from DHCP so that the custom nameservers take precedence.

Two overlay files apply to control plane nodes, and the other two apply to worker nodes.
You add all four files into your `~/.tanzu/tkg/providers/infrastructure-vsphere/ytt/` directory.

The last line of each `overlay-dns` file sets the nameserver addresses.
The code below shows a single nameserver, but you can specify multiple nameservers as a list, for example `nameservers: ["1.2.3.4","5.6.7.8"]`.

File `vsphere-overlay-dns-control-plane.yaml`:

```
#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind":"VSphereMachineTemplate", "metadata": {"name": data.values.CLUSTER_NAME+"-control-plane"}})
---
spec:
  template:
    spec:
      network:
        devices:
        #@overlay/match by=overlay.all, expects="1+"
        -
          #@overlay/match missing_ok=True
          nameservers: ["8.8.8.8"]
```

File `vsphere-overlay-dhcp-control-plane.yaml`:

```
#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind":"KubeadmControlPlane"})
---
spec:
  kubeadmConfigSpec:
    preKubeadmCommands:
    #! disable dns from being emitted by dhcp client
    #@overlay/append
    - echo '[DHCPv4]' >> /etc/systemd/network/10-id0.network
    #@overlay/append
    - echo 'UseDNS=no' >> /etc/systemd/network/10-id0.network
    #@overlay/append
    - '/usr/bin/systemctl restart systemd-networkd 2>/dev/null'
```

File `vsphere-overlay-dns-workers.yaml`:

```
#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind":"VSphereMachineTemplate", "metadata": {"name": data.values.CLUSTER_NAME+"-worker"}})
---
spec:
  template:
    spec:
      network:
        devices:
        #@overlay/match by=overlay.all, expects="1+"
        -
          #@overlay/match missing_ok=True
          nameservers: ["8.8.8.8"]
```

File `vsphere-overlay-dhcp-workers.yaml`:

```
#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")

#@overlay/match by=overlay.subset({"kind":"KubeadmConfigTemplate"})
---
spec:
  template:
    spec:
      #@overlay/match missing_ok=True
      preKubeadmCommands:
      #! disable dns from being emitted by dhcp client
      #@overlay/append
      - echo '[DHCPv4]' >> /etc/systemd/network/10-id0.network
      #@overlay/append
      - echo 'UseDNS=no' >> /etc/systemd/network/10-id0.network
      #@overlay/append
      - '/usr/bin/systemctl restart systemd-networkd 2>/dev/null'
```

### <a id="privileged"></a> Privileged Containers for Workloads

This example lets you configure the `--allow-privileged` option for a workload cluster's [kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/) by setting the variable `ALLOW_PRIVILEGED_CONTAINERS` in a config file or local environment.

1. Modify `~/.tanzu/tkg/providers/config_default.yaml` by adding a new `User customizations` section at the end of the file, after the `Additional Internal Config Values` section:

  ```
  #! User customizations
  ALLOW_PRIVILEGED_CONTAINERS: false
  ```

1. In `~/.tanzu/tkg/providers/ytt/04_user_customizations/`, create a new file `allow_privileged_container.yaml` containing:

    ```
    #@ load("@ytt:overlay", "overlay")
    #@ load("@ytt:data", "data")

    #@ if data.values.TKG_CLUSTER_ROLE == "workload" and data.values.ALLOW_PRIVILEGED_CONTAINERS:

    #@overlay/match missing_ok=True,by=overlay.subset({"kind":"KubeadmControlPlane"})
    ---
    spec:
      kubeadmConfigSpec:
        clusterConfiguration:
          apiServer:
            extraArgs:
              #@overlay/match missing_ok=True
              allow-privileged: true

    #@ end
    ```

  If the `04_user_customizations` directory does not already exist under the top-level `ytt` directory, create it.

### <a id="no-bastion"></a> Disable Bastion Host on AWS

For an example overlay that disables the Bastion host for workload clusters on AWS, see [Disable Bastion Server on AWS](https://github.com/Tanzu-Solutions-Engineering/tkg-lab/blob/main/capi-overrides/disable-aws-bastion.yaml) in the _TKG Lab_ repository.

### <a id="nginx"></a> New Plan `nginx`

This example adds and configures a new workload cluster plan `nginx` that runs an [nginx](https://nginx.org/en/) server.
It uses the Cluster Resource Set (CRS) to deploy the nginx server to vSphere clusters created with the vSphere Cluster API provider version v0.7.6.

1. In `.tkg/providers/infrastructure-vsphere/v0.7.6/`, add a new file `cluster-template-definition-nginx.yaml` with contents identical to the `cluster-template-definition-dev.yaml` and `cluster-template-definition-prod.yaml` files:

    ```
    apiVersion: run.tanzu.vmware.com/v1alpha1
    kind: TemplateDefinition
    spec:
      paths:
        - path: providers/infrastructure-vsphere/v0.7.6/ytt
        - path: providers/infrastructure-vsphere/ytt
        - path: providers/ytt
        - path: bom
          filemark: text-plain
        - path: providers/config_default.yaml
    ```

  The presence of this file creates a new plan, and the `tanzu` CLI parses its filename to create the option `nginx` to pass to `tanzu cluster create --plan`.

1. In `~/.tanzu/tkg/providers/ytt/04_user_customizations/`, create a new file `deploy_service.yaml` containing:

    ```
    #@ load("@ytt:overlay", "overlay")
    #@ load("@ytt:data", "data")
    #@ load("@ytt:yaml", "yaml")

    #@ def nginx_deployment():
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx-deployment
    spec:
      selector:
        matchLabels:
          app: nginx
      replicas: 2
      template:
        metadata:
          labels:
            app: nginx
        spec:
          containers:
          - name: nginx
            image: nginx:1.14.2
            ports:
            - containerPort: 80
    #@ end

    #@ if data.values.TKG_CLUSTER_ROLE == "workload" and data.values.CLUSTER_PLAN == "nginx":

    ---
    apiVersion: addons.cluster.x-k8s.io/v1alpha3
    kind: ClusterResourceSet
    metadata:
      name: #@ "{}-nginx-deployment".format(data.values.CLUSTER_NAME)
      labels:
        cluster.x-k8s.io/cluster-name: #@ data.values.CLUSTER_NAME
    spec:
      strategy: "ApplyOnce"
      clusterSelector:
        matchLabels:
          tkg.tanzu.vmware.com/cluster-name: #@ data.values.CLUSTER_NAME
      resources:
      - name: #@ "{}-nginx-deployment".format(data.values.CLUSTER_NAME)
        kind: ConfigMap
    ---
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: #@ "{}-nginx-deployment".format(data.values.CLUSTER_NAME)
    type: addons.cluster.x-k8s.io/resource-set
    stringData:
      value: #@ yaml.encode(nginx_deployment())

    #@ end
    ```

  In this file, the conditional `#@ if data.values.TKG_CLUSTER_ROLE == "workload" and data.values.CLUSTER_PLAN == "nginx":` applies the overlay that follows to workload clusters with the plan `nginx`.

  If the `04_user_customizations` directory does not already exist under the top-level `ytt` directory, create it.
