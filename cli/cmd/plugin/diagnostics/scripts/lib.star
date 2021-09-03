# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

def prog_checks():
    if prog_avail_local("kind") == "":
        log(prefix="Error", msg="kind program binary not found")
        return False

    if prog_avail_local("tanzu") == "":
        log(prefix="Error", msg="tanzu program binary not found")
        return False

    return True

def capture_node_diagnostics(nodes):
    log(prefix="Info", msg="Capturing information for {} nodes".format(len(nodes)))
    capture(cmd="sudo df -i", resources=nodes)
    capture(cmd="sudo crictl info", resources=nodes)
    capture(cmd="df -h /var/lib/containerd", resources=nodes)
    capture(cmd="sudo systemctl status kubelet", resources=nodes)
    capture(cmd="sudo systemctl status containerd", resources=nodes)
    capture(cmd="sudo journalctl -xeu kubelet", resources=nodes)
    capture(cmd="sudo cat /var/log/cloud-init-output.log", resources=nodes)
    capture(cmd="sudo cat /var/log/cloud-init.log", resources=nodes)

# extracts kubernetes object from cluster
def capture_k8s_objects(k8sconf,cluster_name,nspaces):
    log(prefix="Info", msg="Capturing pod logs: cluster={}".format(cluster_name))
    kube_capture(what="logs", namespaces=nspaces, kube_config=k8sconf)
    log(prefix="Info", msg="Capturing API objects: cluster={}".format(cluster_name))
    kube_capture(what="objects", kinds=["pods", "services"], namespaces=nspaces, kube_config=k8sconf)
    kube_capture(what="objects", kinds=["deployments", "replicasets"], groups=["apps"], namespaces=nspaces, kube_config=k8sconf)
    kube_capture(what="objects", kinds=["apps"], groups=["kappctrl.k14s.io"], namespaces=["tkg-system"], kube_config=k8sconf)
    kube_capture(what="objects", kinds=["tanzukubernetesreleases"], groups=["run.tanzu.vmware.com"], kube_config=k8sconf)
    kube_capture(what="objects", kinds=["configmaps"], namespaces=["tkr-system"], kube_config=k8sconf)
    kube_capture(what="objects", categories=["cluster-api"], kube_config=k8sconf)
