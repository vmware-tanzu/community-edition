<!-- markdownlint-disable MD031 -->
<!-- TODO @randomvariable: Fix spacing to re-enable MD031 -->

# Implementing Log Forwarding with Fluent Bit

[Fluent Bit](https://fluentbit.io/) is a lightweight log processor and forwarder that allows you to collect data and logs from different sources, unify them, and send them to multiple destinations. Tanzu Kubernetes Grid includes signed binaries for Fluent Bit, that you can deploy on management clusters and on Tanzu Kubernetes clusters to provide a log-forwarding service.

The Fluent Bit implementation provided in this release of Tanzu Kubernetes Grid allows you to gather logs from management clusters or Tanzu Kubernetes clusters running in vSphere, Amazon EC2, and Azure, and to forward them to an [Elastic Search](https://www.elastic.co/), [Kafka](https://www.confluent.io/confluent-operator/), [Splunk](https://www.splunk.com/), or HTTP endpoint. The Fluent Bit deployment that Tanzu Kubernetes Grid provides is also pre-configured to expose certain metrics to Prometheus and Grafana.

You can deploy Fluent Bit on any management clusters or Tanzu Kubernetes clusters from which you want to collect logs. First, you configure an output plugin on the cluster from which you want to gather logs, depending on the endpoint that you use. Then, you deploy Fluent Bit on the cluster. The examples in this topic deploy Fluent Bit on a Tanzu Kubernetes cluster.

**IMPORTANT**:

- Tanzu Kubernetes Grid does not support IPv6 addresses. This is because upstream Kubernetes only provides alpha support for IPv6. Always provide IPv4 addresses in the procedures in this section.
- The extensions folder `tkg-extensions-v1.4.0+vmware.1
` contains subfolders for each type of extension, for example, `authentication`, `ingress`, `registry`, and so on. At the top level of the folder there is an additional subfolder named `extensions`. The `extensions` folder also contains subfolders for `authentication`, `ingress`, `registry`, and so on. Take care to run commands from the location provided in the instructions. Commands are usually run from within the `extensions` folder.

## <a id="prereqs"></a> Prerequisites

- You have deployed a management cluster on vSphere, Amazon EC2, or Azure.
- You have downloaded and unpacked the bundle of Tanzu Kubernetes Grid extensions. For information about where to obtain the bundle, see [Download and Unpack the Tanzu Kubernetes Grid Extensions Bundle](index.md#unpack-bundle).
- You have installed the Carvel tools. For information about installing the Carvel tools, see [Install the Carvel Tools](../install-cli.md#install-carvel).
- You have deployed one of the following logging management backends for storing and analyzing logs:
    - Elastic Search.
    - Kafka.
    - Splunk.
    - Another HTTP-based service. For example, VMware vRealize Log Insight.
    - Syslog. For example, VMware vRealize Log Insight.

## <a id="prepare-tkc"></a> Prepare the Cluster for Fluent Bit Deployment

Before you can deploy Fluent Bit on a management cluster or on a Tanzu Kubernetes cluster, you must install the tools that the Fluent Bit extension requires.

**NOTE**: Deploying Fluent Bit on a management cluster does not mean that logs will also be collected from the Tanzu Kubernetes clusters that the management cluster manages. You must deploy Fluent Bit on each cluster from which you want to gather logs.

This procedure applies to all clusters, running on vSphere, Amazon EC2, and Azure.

1. Create a cluster configuration YAML file for a management or workload cluster. For a cluster named `fluentbit-test`, for example:

   ```
   INFRASTRUCTURE_PROVIDER: vsphere
   CLUSTER_NAME: fluentbit-test
   CLUSTER_PLAN: dev
   ```

1. **vSphere**: To deploy the cluster to vSphere, add a line to the configuration file that sets `VSPHERE_CONTROL_PLANE_ENDPOINT` to a static virtual IP (VIP) address for the control plane of the Tanzu Kubernetes cluster. Ensure that this IP address is not in the DHCP range, but is in the same subnet as the DHCP range. If you mapped a fully qualified domain name (FQDN) to the VIP address, you can specify the FQDN instead of the VIP address. For example:

   ```
   VSPHERE_CONTROL_PLANE_ENDPOINT: 10.10.10.10
   ```

1. Deploy the cluster by passing the cluster configuration file to the `tanzu cluster create`:

   ```
   tanzu cluster create fluentbit-test --file fluentbit-test-config.yaml
   ```

1. In a terminal, navigate to the folder that contains the unpacked Tanzu Kubernetes Grid extension manifest files, `tkg-extensions-v1.4.0+vmware.1
/extensions`.

   ```
   cd <path>/tkg-extensions-v1.4.0+vmware.1
/extensions
   ```

   You should see folders for `authentication`, `ingress`, `logging`, `monitoring`, `registry`, and some YAML files. Run all of the commands in these procedures from this location.

1. Get the `admin` credentials of the Tanzu Kubernetes cluster on which to deploy Fluent Bit.

   ```
   tanzu cluster kubeconfig get fluentbit-test --admin
   ```

1. Set the context of `kubectl` to the Tanzu Kubernetes cluster.

   ```
   kubectl config use-context fluentbit-test-admin@fluentbit-test
   ```

1. Create a namespace for the Fluent Bit service on the cluster.

    ```sh
    kubectl apply -f logging/fluent-bit/namespace-role.yaml
    ```
    You should see confirmation that a `tanzu-system-logging` namespace, service account, and RBAC role bindings are created.
    ```
    namespace/tanzu-system-logging created
    serviceaccount/fluent-bit-extension-sa created
    role.rbac.authorization.k8s.io/fluent-bit-extension-role created
    rolebinding.rbac.authorization.k8s.io/fluent-bit-extension-rolebinding created
    clusterrole.rbac.authorization.k8s.io/fluent-bit-extension-cluster-role created
    clusterrolebinding.rbac.authorization.k8s.io/fluent-bit-extension-cluster-rolebinding created
    ```

The cluster is ready for you to deploy the Fluent Bit service. For the next steps, see the appropriate procedure for your logging endpoint:

- [Prepare the Fluent Bit Configuration File for an Elastic Search Output Plugin](#elasticsearch)
- [Prepare the Fluent Bit Configuration File for a Kafka Output Plugin ](#kafka)
- [Prepare the Fluent Bit Configuration File for a Splunk Output Plugin ](#splunk)
- [Prepare the Fluent Bit Configuration File for an HTTP Endpoint Output Plugin](#http)
- [Prepare the Fluent Bit Configuration File for a Syslog Endpoint Output Plugin](#syslog)

## <a id="elasticsearch"></a> Prepare the Fluent Bit Configuration File for an Elastic Search Output Plugin

After you have set up the cluster, you must create the configuration file that is used when you deploy Fluent Bit and use it to create a Kubernetes secret.

This procedure describes how to configure Elastic Search as the output plugin on a cluster for clusters that are running on vSphere, Amazon EC2, and Azure. You must already have deployed Elastic Search as the logging management backend for storing and analyzing logs.

1. Make a copy of the `fluent-bit-data-values.yaml.example` file and name it `fluent-bit-data-values.yaml`.

    ```sh
    cp logging/fluent-bit/elasticsearch/fluent-bit-data-values.yaml.example logging/fluent-bit/elasticsearch/fluent-bit-data-values.yaml
    ```

1. Open `fluent-bit-data-values.yaml` in a text editor.
1. Update `fluent-bit-data-values.yaml` with information about your Tanzu Kubernetes cluster and Elastic Search server.

   - `<TKG_INSTANCE_NAME>`: The name of Tanzu Kubernetes Grid instance. This is the same as the name of the management cluster.
   - `<TKG_CLUSTER_NAME>`: The name of the Tanzu Kubernetes cluster on which you are deploying Fluent Bit.
   - `<ELASTIC_SEARCH_HOST>`: The IP address or host name of the target Elastic Search server.
   - `<ELASTIC_SEARCH_PORT>`: The TCP port of the target Elastic Search server.

      **NOTE**: The host name and port number must be wrapped in quotes (`""`).

    For example:

    ```
    #@data/values
    #@overlay/match-child-defaults missing_ok=True
    ---
    logging:
      image:
        repository: "projects.registry.vmware.com/tkg"
    tkg:
      instance_name: "my-mgmt-cluster"
      cluster_name: "fluentbit-test"
    fluent_bit:
      output_plugin: "elasticsearch"
      elasticsearch:
        host: "elasticsearch"
        port: "9200"
    ```

1. Optionally customize your Fluent Bit deployment.

   For information about other values that you can configure to customize your Fluent Bit deployment, see [Optional Configuration](#configs) below.
1. Create a Kubernetes secret named `fluent-bit-data-values` with the values that you set in `fluent-bit-data-values.yaml`.

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/elasticsearch/fluent-bit-data-values.yaml -n tanzu-system-logging
   ```

   You should see the confirmation `secret/fluent-bit-data-values created`.

For the next steps, see [Deploy the Fluent Bit Extension](#deploy-fluentbit).

## <a id="kafka"></a> Prepare the Fluent Bit Configuration File for a Kafka Output Plugin

After you have set up the cluster, you must create the configuration file that is used when you deploy Fluent Bit and use it to create a Kubernetes secret.

This procedure describes how to configure Kafka as the output plugin on a cluster for clusters that are running on vSphere, Amazon EC2, and Azure. You must already have deployed Kafka as the logging management backend for storing and analyzing logs.

1. Make a copy of the `fluent-bit-data-values.yaml.example` file and name it `fluent-bit-data-values.yaml`.

    ```sh
    cp logging/fluent-bit/kafka/fluent-bit-data-values.yaml.example logging/fluent-bit/kafka/fluent-bit-data-values.yaml
    ```
1. Open `fluent-bit-data-values.yaml` in a text editor.
1. Update `fluent-bit-data-values.yaml` with information about your Tanzu Kubernetes cluster and Kafka server.

   - `<TKG_INSTANCE_NAME>`: The name of Tanzu Kubernetes Grid instance. This is the same as the name of the management cluster.
   - `<TKG_CLUSTER_NAME>`: The name of the Tanzu Kubernetes cluster on which you are deploying Fluent Bit.
   - `<BROKER_SERVICE_NAME>`: The name of the Kafka broker service.
   - `<TOPIC_NAME>`: The name of the topic that ingests the logs in Kafka.

      **NOTE**: The broker service name and topic name values must be wrapped in quotes (`""`).

    For example:

    ```
    #@data/values
    #@overlay/match-child-defaults missing_ok=True
    ---
    logging:
      image:
        repository: "projects.registry.vmware.com/tkg"
    tkg:
      instance_name: "my-mgmt-cluster"
      cluster_name: "fluentbit-test"
    fluent_bit:
      output_plugin: "kafka"
      kafka:
        broker_service_name: "my-kafka-broker"
        topic_name: "my-kafka-topic-name"
    ```

1. Optionally customize your Fluent Bit deployment.

   For information about other values that you can configure to customize your Fluent Bit deployment, see [Optional Configuration](#configs) below.

1. Create a Kubernetes secret named `fluent-bit-data-values` with the values that you set in `fluent-bit-data-values.yaml`.

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/kafka/fluent-bit-data-values.yaml -n tanzu-system-logging
   ```

   You should see the confirmation `secret/fluent-bit-data-values created`.

For the next steps, see [Deploy the Fluent Bit Extension](#deploy-fluentbit).

## <a id="splunk"></a> Prepare the Fluent Bit Configuration File for a Splunk Output Plugin

After you have set up the cluster, you must create the configuration file that is used when you deploy Fluent Bit and use it to create a Kubernetes secret.

This procedure describes how to configure Splunk as the output plugin on a cluster for clusters that are running on vSphere, Amazon EC2, and Azure. You must already have deployed Splunk as the logging management backend for storing and analyzing logs.

1. Make a copy of the `fluent-bit-data-values.yaml.example` file and name it `fluent-bit-data-values.yaml`.

    ```sh
    cp logging/fluent-bit/splunk/fluent-bit-data-values.yaml.example logging/fluent-bit/splunk/fluent-bit-data-values.yaml
    ```
1. Open `fluent-bit-data-values.yaml` in a text editor.
1. Update `fluent-bit-data-values.yaml` with information about your Tanzu Kubernetes cluster and Splunk server.

   - `<TKG_INSTANCE_NAME>`: The name of Tanzu Kubernetes Grid instance. This is the same as the name of the management cluster.
   - `<TKG_CLUSTER_NAME>`: The name of the Tanzu Kubernetes cluster on which you are deploying Fluent Bit.
   - `<SPLUNK_HOST>`: The IP address or host name of the target Splunk Server.
   - `<SPLUNK_PORT>`: The TCP port of the target Splunk Server. The default TCP port for Splunk Server is 8088.
   - `<SPLUNK_TOKEN>`: The authentication token for the HTTP event collector interface.

      **NOTE**: The host name, port number, and token values must be wrapped in quotes (`""`).

    For example:

    ```
    #@data/values
    #@overlay/match-child-defaults missing_ok=True
    ---
    logging:
      image:
        repository: "projects.registry.vmware.com/tkg"
    tkg:
      instance_name: "my-mgmt-cluster"
      cluster_name: "fluentbit-test"
    fluent_bit:
      output_plugin: "splunk"
      splunk:
        host: "mysplunkhost.example.com"
        port: "8088"
        token: "f61b7aecf75e95cd226234f4fe901ed450fa323648165a91bf02f0a07c5199eb"
    ```

1. Optionally customize your Fluent Bit deployment.

   For information about other values that you can configure to customize your Fluent Bit deployment, see [Optional Configuration](#configs) below.
1. Create a Kubernetes secret named `fluent-bit-data-values` with the values that you set in `fluent-bit-data-values.yaml`.

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/splunk/fluent-bit-data-values.yaml -n tanzu-system-logging
   ```

   You should see the confirmation `secret/fluent-bit-data-values created`.

For the next steps, see [Deploy the Fluent Bit Extension](#deploy-fluentbit).

## <a id="http"></a> Prepare the Fluent Bit Configuration File for an HTTP Endpoint Output Plugin

After you have set up the cluster, you must create the configuration file that is used when you deploy Fluent Bit and use it to create a Kubernetes secret.

This procedure describes how to configure an HTTP endpoint as the output plugin on a cluster for clusters that are running on vSphere, Amazon EC2, and Azure. You must already have deployed an HTTP endpoint as the logging management backend for storing and analyzing logs, such as VMware vRealize Log Insight Cloud.

**NOTE**: The provided implementation of Fluent Bit only works with vRealize Log Insight Cloud and not with the on-premises version of vRealize Log Insight. If you want to use the on-premises version of vRealize Log Insight, see [Prepare the Fluent Bit Configuration File for a Syslog Output Plugin](#syslog) below.

4. Make a copy of the `fluent-bit-data-values.yaml.example` file and name it `fluent-bit-data-values.yaml`.

    ```sh
    cp logging/fluent-bit/http/fluent-bit-data-values.yaml.example logging/fluent-bit/http/fluent-bit-data-values.yaml
    ```
1. Open `fluent-bit-data-values.yaml` in a text editor.
1. Update `fluent-bit-data-values.yaml` with information about your Tanzu Kubernetes cluster and HTTP logging endpoint.

   - `<TKG_INSTANCE_NAME>`: The name of Tanzu Kubernetes Grid instance. This is the same as the name of the management cluster.
   - `<TKG_CLUSTER_NAME>`: The name of the Tanzu Kubernetes cluster on which you are deploying Fluent Bit.
   - `<HTTP_HOST>`: The IP address or host name of the target HTTP server.
   - `<HTTP_PORT>`: The TCP port of the target HTTP server.
   - `<URI>`: The HTTP URI for the target Web server.
   - `<HEADER_KEY_VALUE>`: The HTTP header key/value pair. For example, for VMware vRealize Log Insight Cloud, set the `key` to the authorization bearer token and the `value` to the API token.
   - `json`: The data format to use in the HTTP request body. For example, for vRealize Log Insight Cloud, leave this value set to `json`.

      **NOTE**: The host name, port number, URI, header key value, and format values must be wrapped in quotes (`""`).

    For example:

    ```
    #@data/values
    #@overlay/match-child-defaults missing_ok=True
    ---
    logging:
      image:
        repository: projects.registry.vmware.com/tkg
    tkg:
      instance_name: "my-mgmt-cluster"
      cluster_name: "fluentbit-test"
    fluent_bit:
      output_plugin: "http"
      http:
        host: "data.mgmt.cloud.vmware.com"
        port: "443"
        uri: "/le-mans/v1/streams/ingestion-pipeline-stream"
        header_key_value: "Authorization Bearer H986qVcPlbQry1UWbY8QuSMXhAtAYsG6"
        format: "json"
    ```

1. Optionally customize your Fluent Bit deployment.

   For information about other values that you can configure to customize your Fluent Bit deployment, see [Optional Configuration](#configs) below.

1. Create a Kubernetes secret named `fluent-bit-data-values` with the values that you set in `fluent-bit-data-values.yaml`.

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/http/fluent-bit-data-values.yaml -n tanzu-system-logging
   ```

   You should see the confirmation `secret/fluent-bit-data-values created`.

For the next steps, see [Deploy the Fluent Bit Extension](#deploy-fluentbit).

## <a id="syslog"></a> Prepare the Fluent Bit Configuration File for a Syslog Output Plugin

This procedure describes how to configure a syslog endpoint as the output plugin on a cluster that is running on vSphere, Amazon EC2, or Azure. You must already have deployed a syslog endpoint as the logging management backend for storing and analyzing logs, such as VMware vRealize Log Insight.

**NOTE**: The provided implementation of Fluent Bit works only with the on-premises version of vRealize Log Insight and not with vRealize Log Insight Cloud. If you want to use vRealize Log Insight Cloud, see [Prepare the Fluent Bit Configuration File for an HTTP Endpoint Output Plugin](#http) above.

After you have set up the cluster, you must create the Fluent Bit configuration file and use it to create a Kubernetes secret:

1. Make a copy of the `fluent-bit-data-values.yaml.example` file and name it `fluent-bit-data-values.yaml`.

    ```sh
    cp logging/fluent-bit/syslog/fluent-bit-data-values.yaml.example logging/fluent-bit/syslog/fluent-bit-data-values.yaml
    ```
1. Open `fluent-bit-data-values.yaml` in a text editor.
1. Update `fluent-bit-data-values.yaml` with information about your cluster and syslog endpoint.

   - `<TKG_INSTANCE_NAME>`: The name of the Tanzu Kubernetes Grid instance. This is the same as the name of the management cluster.
   - `<TKG_CLUSTER_NAME>`: The name of the Tanzu Kubernetes cluster on which you are deploying Fluent Bit.
   - `<SYSLOG_HOST>`: The domain name or IP address of your syslog server.
   - `<SYSLOG_PORT>`: The TCP or UDP port of your syslog server.
   - `<SYSLOG_MODE>`: The transport protocol to use for syslog messages. Set to `tcp`, `udp`, or `tls`.
   - `<SYSLOG_FORMAT>`: The format to use for syslog messages. For example, `rfc5424`.

      **NOTE**: These values must be wrapped in quotes (`""`).

    For example:

    ```
    #@data/values
    #@overlay/match-child-defaults missing_ok=True
    ---
    logging:
      image:
        repository: projects.registry.vmware.com/tkg
    tkg:
      instance_name: "my-mgmt-cluster"
      cluster_name: "fluentbit-test"
    fluent_bit:
      output_plugin: "syslog"
      syslog:
        host: "10.182.176.50"
        port: "514"
        mode: "tcp"
        format: "rfc5424"
      ```

1. Optionally, customize your Fluent Bit deployment.

   For information about other values that you can configure to customize your Fluent Bit deployment, see [Optional Configuration](#configs) below.

1. Create a Kubernetes secret named `fluent-bit-data-values` with the values that you set in `fluent-bit-data-values.yaml`.

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/syslog/fluent-bit-data-values.yaml -n tanzu-system-logging
   ```

   You should see the confirmation `secret/fluent-bit-data-values created`.

For the next steps, see [Deploy the Fluent Bit Extension](#deploy-fluentbit).

## <a id="deploy-fluentbit"></a> Deploy the Fluent Bit Extension

After you have prepared the cluster and updated the appropriate configuration file for your logging endpoint, you can deploy Fluent Bit on the cluster.

1. Deploy the Fluent Bit extension.

   ```sh
   kubectl apply -f logging/fluent-bit/fluent-bit-extension.yaml
   ```

   You should see the confirmation `extension.clusters.tmc.cloud.vmware.com/fluent-bit created`.

1. View the status of the Fluent Bit service itself.

    ```
    kubectl get app fluent-bit -n tanzu-system-logging
    ```
    The status of the Fluent Bit app should show `Reconcile Succeeded` when Fluent Bit has deployed successfully.
    ```
    NAME         DESCRIPTION           SINCE-DEPLOY   AGE
    fluent-bit   Reconcile succeeded   54s            14m
    ```
1. If the status is not `Reconcile Succeeded`, view the full status details of the Fluent Bit service.

   Viewing the full status can help you to troubleshoot the problem.

    ```
    kubectl get app fluent-bit -n tanzu-system-logging -o yaml
    ```

1. Check that the new services are running by listing all of the pods that are running in the cluster.

    ```
    kubectl get pods -A
    ```

   In the `tanzu-system-logging` namespace, you should see the `fluent-bit` services running in a pod with names similar to `fluent-bit-6zn4l`.

    ```
    NAMESPACE              NAME                     READY   STATUS    RESTARTS   AGE
    [...]
    tanzu-system-logging   fluent-bit-6zn4l         1/1     Running   0          22m
    tanzu-system-logging   fluent-bit-g8tkp         1/1     Running   0          22m
    ```

Fluent Bit will now capture logs from all of the containers that run in pods in the cluster. Fluent Bit also captures logs from `systemd` for the `kubelet` and `containerd` services. The `kubelet` and `containerd` services do not run in containers, and they write their logs to `journald`. Fluent Bit implements a `systemd` input plugin to collect log messages from the `journald` daemon.

## <a id="view-logs"></a> Monitoring and Viewing Logs

You can use different dashboards to view all of the logs that Fluent Bit collects, depending on the logging endpoint that you configured. For example, if you configured an Elastic Search endpoint, you can view logs in the Kibana dashboard, or if you implemented an HTTP endpoint, you can view them in vRealize Log Insight dashboards. The image below shows an example of a Kibana dashboard.

![View Fluent Bit logs in Kibana](../images/fluent-bit-kibana.png)

The Fluent Bit deployment that Tanzu Kubernetes Grid provides is pre-configured to expose certain metrics to Prometheus. If you also deploy the Tanzu Kubernetes Grid implementation of Prometheus on your cluster, it includes a predefined job named `kubernetes_pods` that you can use to view Fluent Bit metrics in the Prometheus dashboard without performing any additional configuration.

![View Fluent Bit logs in Prometheus dashboard](../images/prometheus-dashboard-fluentbit.png)

The image below shows an example of the same data viewed in a Grafana dashboard.

![View Fluent Bit logs in Grafana](../images/fluent-bit-grafana.png)

For information about the Prometheus and Grafana extensions, see [Implementing Monitoring with Prometheus and Grafana](monitoring.md).

## <a id="configs"></a> Optional Configuration

In addition to the minimum configuration provided in the `fluent-bit-data-values.yaml` file for your logging endpoint, you can customize your configuration by adding values that you can copy from the file `tkg-extensions-v1.4.0+vmware.1
/logging/fluent-bit/values.yaml`. Note that this file is not located in the `tkg-extensions-v1.4.0+vmware.1
/extensions/logging/fluent-bit` folder, but in the `logging` folder that is at the same level as the `extensions` folder.

You can also customize your Fluent Bit logging setup using `ytt` overlays, as described in [Extensions and Shared Services](../ytt.md#extensions) in _Customizing Clusters, Plans, and Extensions with ytt Overlays_ and in the extensions mods examples in the [TKG Lab repository](https://github.com/Tanzu-Solutions-Engineering/tkg-lab).

The table below contains information on the values you can copy from the `tkg-extensions-v1.4.0+vmware.1
/logging/fluent-bit/values.yaml` file and how they can be used to modify the default behavior of Fluent Bit when deployed onto a Tanzu Kubernetes cluster.

**NOTE**: Where applicable, the settings that you configure in `fluent-bit-data-values.yaml` override any settings that you configure in `values.yaml`.

If you reconfigure your Fluent Bit settings after the initial deployment, you must follow the steps in [Update a Running Fluent Bit Deployment](#update) in order to apply the new configuration to the cluster.

| Parameter                              | Type and Description                                                                                                                                           | Default                            |
| -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------- |
| `logging.namespace`                    | String. Namespace in which to deploy Fluent Bit.                                                                                                               | `tanzu-system-logging`             |
| `logging.service_account_name`         | String. Name of the Fluent Bit service account.                                                                                                                | `fluent-bit`                       |
| `logging.cluster_role_name`            | String. Name of the cluster role which grants `get`, `watch`, and `list` permissions to Fluent Bit.                                                            | `fluent-bit-read`                  |
| `logging.image.name`                   | String. Name of the Fluent Bit image.                                                                                                                          | `fluent-bit`                       |
| `logging.image.tag`                    | String. Fluent bit image tag.                                                                                                                                  | `v1.5.3_vmware.1`                  |
| `logging.image.repository`             | Repository containing the Fluent Bit image String.                                                                                                             | `projects.registry.vmware.com/tkg` |
| `logging.image.pullPolicy`             | Fluent bit image pull policy string                                                                                                                            | `IfNotPresent`                     |
| `logging.update_strategy`              | String. Update strategy to be used when updating DaemonSet.                                                                                                    | `RollingUpdate`                    |
| `tkg.cluster_name`                     | String. Name of the Tanzu Kubernetes cluster. Set this value in `fluent-bit-data-values.yaml`.                                                                 | `null`                             |
| `tkg.instance_name`                    | String. Name of Tanzu Kubernetes Grid management cluster. Set this value in `fluent-bit-data-values.yaml`.                                                     | `null`                             |
| `fluent_bit.log_level`                 | String. Log level to use for Fluent Bit.                                                                                                                       | `info`                             |
| `fluent_bit.output_plugin`             | String. Set the back end to which Fluent Bit should flush the information it gathers. Set this value in `fluent-bit-data-values.yaml`.                         | `null`                             |
| `fluent_bit.elasticsearch.host`        | String. IP address or host name of the target Elastic Search instance. Set this value in `fluent-bit-data-values.yaml`.                                        | `null`                             |
| `fluent_bit.elasticsearch.port`        | Integer. TCP port of the target Elastic Search instance. Set this value in `fluent-bit-data-values.yaml`.                                                      | `null`                             |
| `fluent_bit.kafka.broker_service_name` | String. Single or multiple list of Kafka Brokers, for example 192.168.1.3:9092. Set this value in `fluent-bit-data-values.yaml`.                               | `null`                             |
| `fluent_bit.kafka.topic_name`          | String. Single entry or list of topics separated by (`,`) that Fluent Bit will use to send messages to Kafka. Set this value in `fluent-bit-data-values.yaml`. | `null`                             |
| `fluent_bit.splunk.host`               | String. IP address or host name of the target Splunk Server. Set this value in `fluent-bit-data-values.yaml`.                                                  | `null`                             |
| `fluent_bit.splunk.port`               | Integer. TCP port of the target Splunk Server. Set this value in `fluent-bit-data-values.yaml`.                                                                | `null`                             |
| `fluent_bit.splunk.token`              | String. Specify the Authentication Token for the HTTP Event Collector interface. Set this value in `fluent-bit-data-values.yaml`.                              | `null`                             |
| `fluent_bit.http.host`                 | String. IP address or host name of the target HTTP Server. Set this value in `fluent-bit-data-values.yaml`.                                                    | `null`                             |
| `fluent_bit.http.port`                 | Integer. TCP port of the target HTTP Server. Set this value in `fluent-bit-data-values.yaml`.                                                                  | `null`                             |
| `fluent_bit.http.header_key_value`     | String. HTTP header key/value pair. Multiple headers can be set. Set this value in `fluent-bit-data-values.yaml`.                                              | `null`                             |
| `fluent_bit.http.format`               | String. Specify the data format to be used in the HTTP request body. Set this value in `fluent-bit-data-values.yaml`.                                          | `null`                             |
| `fluent_bit.syslog.host`               | String. Domain name or IP address of the target syslog server. Set this value in `fluent-bit-data-values.yaml`.                                                | `null`                             |
| `fluent_bit.syslog.port`               | Integer. TCP or UDP port of the target syslog server. Set this value in `fluent-bit-data-values.yaml`.                                                         | `null`                             |
| `fluent_bit.syslog.mode`               | String. Transport protocol for syslog messages. This must be `tcp`, `udp` or `tls`. Set this value in `fluent-bit-data-values.yaml`.                           | `null`                             |
| `fluent_bit.syslog.format`             | String. Format for syslog messages. Set this value in `fluent-bit-data-values.yaml`.                                                                           | `null`                             |
| `host_path.volume_1`                   | String. Directory path from the host node's file system into the pod, for volume 1.                                                                            | `/var/log`                         |
| `host_path.volume_2`                   | String. Directory path from the host node's file system into the pod, for volume 2.                                                                            | `/var/lib/docker/containers`       |
| `host_path.volume_3`                   | String. Directory path from the host node's file system into the pod, for volume 3.                                                                            | `/run/log`                         |

## <a id="update"></a> Update a Running Fluent Bit Deployment

If you need to make changes to the configuration of the Fluent Bit extension, you can update the extension.

Perform the steps in this procedure if you modify either of the `extensions/logging/fluent-bit/<PROVIDER>/fluent-bit-data-values.yaml` or  `logging/fluent-bit/values.yaml` files after initial deployment of Fluent Bit.

1. Modify either or both of the `extensions/logging/fluent-bit/<PROVIDER>/fluent-bit-data-values.yaml` or  `logging/fluent-bit/values.yaml` files to update your configuration.
1. Update the Kubernetes secret.

   This command assumes that you are running it from `tkg-extensions-v1.4.0+vmware.1
/extensions`.

  **Elastic Search**:

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/elasticsearch/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```

  **Kafka**:

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/kafka/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```

  **Splunk**:

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/splunk/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```

  **HTTP**:

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/http/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```

   **Syslog**

   ```sh
   kubectl create secret generic fluent-bit-data-values --from-file=values.yaml=logging/fluent-bit/syslog/fluent-bit-data-values.yaml -n tanzu-system-logging -o yaml --dry-run | kubectl replace -f-
   ```

   Note that the final `-` on the `kubectl replace` command above is necessary to instruct `kubectl` to accept the input being piped to it from the `kubectl create secret` command.

   The Fluent Bit extension will be reconciled using the new values you just added. The changes should show up in five minutes or less. This is handled by the Kapp controller, which synchronizes every five minutes.
