# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# extract diagnostic info from local kind boostrap cluster
def diagnose_bootstrap_clusters(workdir, clusters, outputdir):
    if len(clusters) == 0:
        log(prefix="Warn", msg="Bootstrap cluster: no cluster found: nothing will be collected")
        return

    nspaces=[
        "capi-kubeadm-bootstrap-system",
        "capi-kubeadm-control-plane-system",
        "capi-system",
        "capi-webhook-system",
        "capv-system",
        "capa-system",
        "cert-manager",
        "tkg-system",
    ]

    # for each tkg-kind cluster:
    #  - capture kind logs, export kubecfg, and api objects
    for kind_cluster in clusters:
        wd = "{}/{}".format(workdir, kind_cluster)
        log(prefix="Info", msg="Bootstrap cluster: {}: capturing node logs".format(kind_cluster))
        run_local("kind export logs --name {} {}/kind-logs".format(kind_cluster, wd))

        # extract kubeconfig file for cluster
        kind_cfg = capture_local(
            cmd="kind get kubeconfig --name {0}".format(kind_cluster),
            workdir=wd,
            file_name="{}.kubecfg".format(kind_cluster)
        )

        conf = crashd_config(workdir=wd)
        k8sconf = kube_config(path=kind_cfg)

        capture_k8s_objects(k8sconf, kind_cluster, nspaces)

        # remove kubeconfig before archiving
        run_local("rm {}".format(kind_cfg))

        arc_file = "{}/bootstrap.{}.diagnostics.tar.gz".format(outputdir, kind_cluster)
        log(prefix="Info", msg="Archiving: {}".format(arc_file))
        archive(output_file=arc_file, source_paths=[conf.workdir])

# return all bootstrap clusters in kind (tkg-kind-xxxx) or
# returns the cluster that matches name
def get_bootstrap_clusters(name):
    clusters = run_local("kind get clusters").split('\n')
    result = []

    for cluster in clusters:
        if name == cluster:
            result.append(cluster)
            break

        if cluster.startswith("tkg-kind"):
            result.append(cluster)

    if len(result) > 0:
        log(prefix="Info", msg="Found bootstrap cluster(s): {}".format(result))
    return result

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

    # argument validation
    name = None
    if hasattr(args, "bootstrap_cluster_name") and len(args.bootstrap_cluster_name) > 0:
        log(prefix="Info", msg="Bootstrap cluster: name={}".format(args.bootstrap_cluster_name))
        name = args.bootstrap_cluster_name

    if name != None and not name.startswith("tkg-kind"):
        log(prefix="Warn", msg="Bootstrap cluster: specified name may not be valid: {}".format(name))

    clusters=get_bootstrap_clusters(name)

    # diagnose boostrap cluster
    diagnose_bootstrap_clusters(
        workdir=workdir,
        clusters=clusters,
        outputdir=outputdir,
    )


# starting point
diagnose()