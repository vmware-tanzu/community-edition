#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

################################################################################
##                               GLOBAL VARIABLES                             ##
################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
CLUSTER_NAME=${CLUSTER_NAME:-"kind"}
CLUSTER_CONFIG=${CLUSTER_CONFIG:-"${SCRIPT_DIR}/../config/kind-cluster-config.yaml"}
KAPP_CONTROLLER_MANIFEST=${KAPP_CONTROLLER_MANIFEST:-"https://raw.githubusercontent.com/vmware-tanzu/carvel-kapp-controller/develop/alpha-releases/v0.20.0-rc.1.yml"}

################################################################################
##                               FUNCTIONS                                    ##
################################################################################
function prepare_cluster() {
  kind get clusters | grep "${CLUSTER_NAME}" && kind delete cluster --name="${CLUSTER_NAME}"
  kind create cluster --name="${CLUSTER_NAME}" --config="${CLUSTER_CONFIG}"
  kind get kubeconfig --name="${CLUSTER_NAME}" > ~/.kube/kind-kubeconfig
  export KUBECONFIG=~/.kube/kind-kubeconfig

  kubectl apply -f "${KAPP_CONTROLLER_MANIFEST}"
  kubectl rollout status --timeout=300s deployment/kapp-controller -n kapp-controller
}

function delete_cluster() {
  kind delete cluster --name="${CLUSTER_NAME}"
}

"$@"
