# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# diagnose_standalone_cluster retrieves cluster information
# from a non-management standalone cluster.
def diagnose_management_cluster(workdir, kubeconfig, cluster_name, context_name, outputdir):
    conf = crashd_config(workdir=workdir)
    k8sconfig = kube_config(path=kubeconfig, cluster_context=context_name)
    log(prefix="Info", msg="Capturing management cluster diagnostics: cluster={}; context={}; kubeconfig={};".format(cluster_name, context_name, kubeconfig))

    nspaces=[
        "capi-kubeadm-bootstrap-system",
        "capi-kubeadm-control-plane-system",
        "capi-system",
        "capi-webhook-system",
        "cert-manager",
        "tkg-system",
        "kube-system",
        "tkr-system",
    ]

    capture_k8s_objects(k8sconfig, cluster_name, nspaces)

    arc_file = "{}/management-cluster.{}.diagnostics.tar.gz".format(outputdir, cluster_name)
    log(prefix="Info", msg="Archiving: {}".format(arc_file))
    archive(output_file=arc_file, source_paths=[conf.workdir])

def diagnose():
    # program pre-checks
    if not prog_checks():
        log(prefix="Error", msg="One or more required program(s) missing")
        return

    workdir = "./diagnostics"
    if hasattr(args, "workdir") and len(args.workdir) > 0:
        workdir = args.workdir

    outputdir = "./"
    if hasattr(args, "outputdir") and len(args.outputdir) > 0:
        outputdir = args.outputdir

    kubeconfig = "{}/.kube/config".format(os.home)
    if hasattr(args, "management_kubeconfig") and len(args.management_kubeconfig) > 0:
        kubeconfig = args.management_kubeconfig

    if not hasattr(args, "management_cluster_name") or len(args.management_cluster_name) == 0:
        log(prefix="Error", msg="management-cluster-name is required")
        return
    name = args.management_cluster_name

    management_context = "{}-admin@{}".format(name, name)
    if hasattr(args, "management_context") and len(args.management_context) > 0:
        management_context = args.management_context

    # diagnose cluster
    diagnose_management_cluster(
        workdir=workdir,
        kubeconfig=kubeconfig,
        context_name=management_context,
        cluster_name=name,
        outputdir=outputdir,
    )

# starting point
diagnose()