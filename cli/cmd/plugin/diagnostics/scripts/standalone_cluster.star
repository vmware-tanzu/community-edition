# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# diagnose_standalone_cluster retrieves cluster information
# from a non-management standalone cluster.
def diagnose_standalone_cluster(workdir, kubeconfig, cluster_name, context_name, outputdir):
    conf = crashd_config(workdir=workdir)
    k8sconfig = kube_config(path=kubeconfig, cluster_context=context_name)
    log(prefix="Info", msg="Retrieving standalone cluster: cluster={}; context={}; kubeconfig={};".format(cluster_name, context_name, kubeconfig))

    nspaces=[
        "capi-kubeadm-bootstrap-system",
        "capi-kubeadm-control-plane-system",
        "capi-system",
        "capi-webhook-system",
        "cert-manager",
        "tkg-system",
        "kube-system",
        "tkr-system",
        "capa-system",
        "capd-system",
        "capv-system",
        "capz-system",
    ]

    capture_k8s_objects(k8sconfig, cluster_name, nspaces)
    capture_pod_describe(workdir=workdir,k8sconf=kubeconfig,context=context_name,namespaces=nspaces)
    capture_node_describe(workdir=workdir,k8sconf=kubeconfig,context=context_name)
    capture_summary(workdir=workdir,k8sconf=kubeconfig,context=context_name)

    arc_file = "{}/standalone-cluster.{}.diagnostics.tar.gz".format(outputdir, cluster_name)
    log(prefix="Info", msg="Archiving: {}".format(arc_file))
    archive(output_file=arc_file, source_paths=[conf.workdir])


# starting point
diagnose_standalone_cluster(
    workdir=args.workdir,
    kubeconfig=args.standalone_kubeconfig,
    context_name=args.standalone_context,
    cluster_name=args.standalone_cluster_name,
    outputdir=args.outputdir,
)