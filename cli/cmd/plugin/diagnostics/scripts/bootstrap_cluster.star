# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# extract diagnostic info from local kind boostrap cluster
def diagnose_bootstrap_clusters(kubeconfig, cluster, workdir, outputdir):
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
    wd = "{}/{}".format(workdir, cluster)
    log(prefix="Info", msg="Bootstrap cluster: {}: capturing node logs".format(cluster))

    # collect node logs and info
    nodes = get_kind_nodes(cluster)
    for node in nodes:
        log(prefix="Debug", msg="Capturing logs for node {}".format(node))
        capture_kind_logs(node, wd)

    # extract kubeconfig file for cluster
    control_planes = get_control_plane_nodes(cluster)
    if len(control_planes) == 0:
        log(prefix="Warn", msg="Cluster {}: has no control plane node:".format(cluster))

    conf = crashd_config(workdir=wd)
    k8sconf = kube_config(path=kubeconfig)

    capture_k8s_objects(k8sconf, cluster, nspaces)

    arc_file = "{}/bootstrap.{}.diagnostics.tar.gz".format(outputdir, cluster)
    log(prefix="Info", msg="Archiving: {}".format(arc_file))
    archive(output_file=arc_file, source_paths=[conf.workdir])


# returns a priviledged command string to launch a contained process
# docker exec --privileged <container> <command>
def docker_exec_cmd(container, cmd):
    return "docker exec --privileged {} {}".format(container, cmd)

# simulates command `kind export log`
# to collect logs from cluster_node
def capture_kind_logs(cluster_node, dir):
    # docker exec --privileged <node_name> cat /kind/version >> kubernetes-version.txt
    command = docker_exec_cmd(cluster_node, "cat /kind/version")
    capture_local(
        cmd=docker_exec_cmd(cluster_node, "cat /kind/version"),
        workdir=dir,
        file_name="kubernetes-version.txt"
    )

    # journalctl --no-pager >> journal.log
    capture_local(
        cmd=docker_exec_cmd(cluster_node,"journalctl --no-pager"),
        workdir=dir,
        file_name="journal.log"
    )

    # journalctl --no-pager -u kubelet.service >> kubelet.log
    capture_local(
        cmd=docker_exec_cmd(cluster_node, "journalctl --no-pager -u kubelet.service"),
        workdir=dir,
        file_name="kubelet.log"
    )

    # journalctl --no-pager -u containerd.service >> container.log
    capture_local(
        cmd=docker_exec_cmd(cluster_node, "journalctl --no-pager -u containerd.service"),
        workdir=dir,
        file_name="container.log"
    )

    # docker logs <node-name> >> serial.log
    capture_local(cmd="docker logs {}".format(cluster_node), workdir=dir, file_name="serial.log")

    # capture docker info `docker inspect >> inspect.json`
    capture_local(cmd="docker inspect {}".format(cluster_node), workdir=dir, file_name="inspect.log")

# returns a list of kind clusters in Docker
def get_kind_clusters():
    return run_local("""docker ps -a --filter 'label=io.x-k8s.kind.cluster' --format '{{.Label "io.x-k8s.kind.cluster"}}'""").split('\n')

# Returns all nodes in kind cluster
def get_kind_nodes(cluster_name):
    return run_local("""docker ps -a --filter 'label=io.x-k8s.kind.cluster={}' --format '{{{{.Names}}}}'""".format(cluster_name)).split('\n')

# Returns a list of control plane nodes for specified cluster
def get_control_plane_nodes(cluster_name):
    return run_local("""docker ps -a --filter 'label=io.x-k8s.kind.cluster={}' --filter 'label=io.x-k8s.kind.role=control-plane' --format '{{{{.Names}}}}'""".format(cluster_name)).split('\n')


def diagnose():
    # program pre-checks
    if not prog_checks():
        log(prefix="Error", msg="One or more required program(s) missing")
        return

    # argument validation
    name = None
    if not hasattr(args, "bootstrap_cluster_name") or len(args.bootstrap_cluster_name) == 0:
        log(prefix="Error", msg="No valid bootstrap cluster provided")
        return

    name = args.bootstrap_cluster_name

    if name != None and not name.startswith("tkg-kind"):
        log(prefix="Warn", msg="Bootstrap cluster: may not be a valid tanzu bootstrap cluster: {}".format(name))

    kubecfg = None
    if not hasattr(args, "bootstrap_kubeconfig") or len(args.bootstrap_kubeconfig) == 0:
        log(prefix="Error", msg="No valid bootstrap cluster kubeconfig provided")
        return

    kubecfg = args.bootstrap_kubeconfig

    workdir = "./diagnostics"
    if hasattr(args, "workdir") and len(args.workdir) > 0:
        workdir = args.workdir

    outputdir = "./"
    if hasattr(args, "outputdir") and len(args.outputdir) > 0:
        outputdir = args.outputdir


    # clusters=get_bootstrap_clusters(name)

    # diagnose boostrap cluster
    diagnose_bootstrap_clusters(
        kubeconfig=kubecfg,
        cluster=name,
        workdir=workdir,
        outputdir=outputdir,
    )


# starting point
diagnose()