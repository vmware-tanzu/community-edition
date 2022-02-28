# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# diagnose_unmanaged_cluster retrieves cluster information
# from a non-management unmanaged cluster.
def diagnose_unmanaged_cluster(workdir, kubeconfig, cluster_name, context_name, outputdir):
    conf = crashd_config(workdir=workdir)
    k8sconfig = kube_config(path=kubeconfig, cluster_context=context_name)
    log(prefix="Info", msg="Retrieving unmanaged cluster: cluster={}; context={}; kubeconfig={};".format(cluster_name, context_name, kubeconfig))

    nspaces=[
        "tkg-system",
        "kube-system",
		"local-path-storage",
		"tanzu-package-repo-global"
    ]

    capture_k8s_objects(k8sconfig, cluster_name, nspaces)
    capture_pod_describe(workdir=workdir,k8sconf=kubeconfig,context=context_name,namespaces=nspaces)
    capture_node_describe(workdir=workdir,k8sconf=kubeconfig,context=context_name)
    capture_summary(workdir=workdir,k8sconf=kubeconfig,context=context_name)

    arc_file = "{}/unmanaged-cluster.{}.diagnostics.tar.gz".format(outputdir, cluster_name)
    log(prefix="Info", msg="Archiving: {}".format(arc_file))
    archive(output_file=arc_file, source_paths=[conf.workdir])


# starting point
diagnose_unmanaged_cluster(
    workdir=args.workdir,
    kubeconfig=args.unmanaged_kubeconfig,
    context_name=args.unmanaged_context,
    cluster_name=args.unmanaged_cluster_name,
    outputdir=args.outputdir,
) 
