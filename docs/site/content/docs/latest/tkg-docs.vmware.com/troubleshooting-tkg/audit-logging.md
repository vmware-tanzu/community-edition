# Audit Logging

This topic describes audit logging in Tanzu Kubernetes Grid.

## Overview

In Tanzu Kubernetes Grid, you can access the following audit logs:

* Audit logs from the Kubernetes API server. See [Kubernetes Audit Logs](#api-logs) below.
* System audit logs for each node in a cluster, collected using `auditd`. See [System Audit Logs for Nodes](#system-logs) below.

## <a id="api-logs"></a> Kubernetes Audit Logs

Kubernetes audit logs record requests to the Kubernetes API server. To enable Kubernetes auditing on a management or Tanzu Kubernetes cluster, set the `ENABLE_AUDIT_LOGGING` variable to `true` before you deploy the cluster.

To access these logs in Tanzu Kubernetes Grid, navigate to `/var/log/kubernetes/audit.log` on the control plane node. If you deploy Fluent Bit on the cluster, it will forward the logs to your log destination. For instructions, see [Implementing Log Forwarding with Fluent Bit](../extensions/logging-fluentbit.md).

To view the audit policy and audit backend configuration, navigate to:

   * `/etc/kubernetes/audit-policy.yaml` on the control plane node
   * `~/.tanzu/tkg/providers/ytt/03_customizations/audit-logging/audit_logging.yaml` on your machine

## <a id="system-logs"></a> System Audit Logs for Nodes

When you deploy a management or Tanzu Kubernetes cluster, `auditd` is enabled on the
cluster by default. You can access your system audit logs on each node in the cluster by navigating to `/var/log/audit/audit.log`.

If you deploy Fluent Bit on the cluster, it will forward these audit logs to your log destination. For instructions, see [Implementing Log Forwarding with Fluent Bit](../extensions/logging-fluentbit.md).
