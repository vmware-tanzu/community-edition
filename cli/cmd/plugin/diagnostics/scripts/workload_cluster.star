# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# diagnose_standalone_cluster retrieves cluster information
# from a non-management standalone cluster.
def diagnose_workload_cluster(workdir, infra, kubeconfig, cluster_name, context_name, outputdir):
    conf = crashd_config(workdir=workdir)
    k8sconfig = kube_config(path=kubeconfig, cluster_context=context_name)
    log(prefix="Info", msg="Retrieving workload cluster: cluster={}; context={}; kubeconfig={};".format(cluster_name, context_name, kubeconfig))

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

    if infra == "aws":
        nspaces.append("capa-system")
    else:
        nspaces.append("capv-system")

    capture_k8s_objects(k8sconfig, cluster_name, nspaces)

    arc_file = "{}/workload-cluster.{}.diagnostics.tar.gz".format(outputdir, cluster_name)
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

    if not hasattr(args, "workload_cluster_name") or len(args.workload_cluster_name) == 0:
        log(prefix="Error", msg="workload-cluster-name is required")
        return
    name = args.workload_cluster_name

    standalone = False
    if hasattr(args, "workload_cluster_standalone"):
        standalone = True

    kubeconfig = "{}/.kube/config".format(os.home)
    namespace = "default"
    workload_context = "{}-admin@{}".format(name, name)

    if not standalone:
        if hasattr(args, "workload_cluster_namespace") and len(args.workload_cluster_namespace) > 0:
            namespace = args.workload_cluster_namespace

        # merge workload kubeconfig/context in default config.
        log(prefix="Info", msg="Retrieving workload cluster credentials")
        run_local(cmd="tanzu cluster kubeconfig get {} --admin --namespace={}".format(name, namespace))
    else:
        if hasattr(args, "workload_kubeconfig") and len(args.workload_kubeconfig) > 0:
            kubeconfig = args.workload_kubeconfig

        if hasattr(args, "workload_context") and len(args.workload_context) > 0:
            workload_context = args.workload_context

    infra = "docker"
    if hasattr(args, "workload_infra") and len(args.workload_infra) > 0:
        infra = args.workload_infra

    # # collect nodes data
    # sshconfig = None
    # if hasattr(args, "ssh_user") and hasattr(args, "ssh_pk_file"):
    #     sshconfig = ssh_config(username=args.ssh_user, private_key_path=args.ssh_pk_file)
    #
    # if sshconfig != None:
    #     nodes=resources(provider=kube_nodes_provider(kube_config=k8sconfig, ssh_config=sshconfig))
    #     # capture_node_diagnostics(nodes)

    # diagnose cluster
    diagnose_workload_cluster(
        workdir=workdir,
        infra=infra,
        kubeconfig=kubeconfig,
        context_name=workload_context,
        cluster_name=name,
        outputdir=outputdir,
    )

# starting point
diagnose()