# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

def prog_checks():
    if prog_avail_local("tanzu") == "":
        log(prefix="Error", msg="tanzu program binary not found")
        return False

    return True

def capture_summary(workdir,k8sconf,context):
    if not context:
        command = "kubectl get po,deploy,cluster,kubeadmcontrolplane,machine,machinedeployment -A --kubeconfig {}".format(k8sconf)
    else:
        command = "kubectl get po,deploy,cluster,kubeadmcontrolplane,machine,machinedeployment -A --kubeconfig {} --context {}".format(k8sconf,context)

    capture_local(cmd=command,workdir=workdir,file_name="cluster-summary.txt")


def capture_pod_describe(workdir,k8sconf,context,namespaces):
    wd = "{}/describe".format(workdir)

    for ns in namespaces:
        if not context:
            command = "kubectl describe po --kubeconfig {} -n {}".format(k8sconf, ns)
        else:
            command = "kubectl describe po --kubeconfig {} --context {} -n {}".format(k8sconf, context, ns)

        capture_local(cmd=command,workdir=wd,file_name="{}_pods.txt".format(ns))

def capture_node_describe(workdir, k8sconf, context):
    wd = "{}/describe".format(workdir)
    if not context:
        command = "kubectl describe nodes --kubeconfig {}".format(k8sconf)
    else:
        command = "kubectl describe nodes --kubeconfig {} --context {}".format(k8sconf, context)

    capture_local(cmd=command,workdir=wd,file_name="nodes.txt")

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
    log(prefix="Info", msg="Capturing pod logs: cluster={}; kubeconf={}".format(cluster_name, k8sconf))
    kube_capture(what="logs", namespaces=nspaces, kube_config=k8sconf)
    log(prefix="Info", msg="Capturing API objects: cluster={}".format(cluster_name))
    kube_capture(what="objects", kinds=["pods", "services"], namespaces=nspaces, kube_config=k8sconf)
    kube_capture(what="objects", kinds=["deployments", "replicasets"], groups=["apps"], namespaces=nspaces, kube_config=k8sconf)
    kube_capture(what="objects", kinds=["apps"], groups=["kappctrl.k14s.io"], namespaces=["tkg-system"], kube_config=k8sconf)
    kube_capture(what="objects", kinds=["tanzukubernetesreleases"], groups=["run.tanzu.vmware.com"], kube_config=k8sconf)
    kube_capture(what="objects", kinds=["configmaps"], namespaces=["tkr-system"], kube_config=k8sconf)
    kube_capture(what="objects", categories=["cluster-api"], kube_config=k8sconf)
