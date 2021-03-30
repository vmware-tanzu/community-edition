#!/bin/bash

# This script bundles an addons configuration / manifests and uploads a Package CR to the cluster.
# Lastly it applies the InstalledPackage, ClusterRoleBinding, and ServiceAccount YAMLs to validate the
# addon works as expected.
# This script is intended for those validating changes to packages that may be installed in a TCE-based
# cluster.

ADDON=$1
TEMPLATE_IMAGE=$2

if [ -z "$ADDON" ] || [ -z "$TEMPLATE_IMAGE" ];then
  echo "usage: ./hack/addon-dev/dev-deploy-addon.sh [PACKAGE_NAME] [REPO:TAG]"
  exit 1
fi

echo "building imgpkg bundle"
imgpkg push --bundle "$TEMPLATE_IMAGE" --file "addons/packages/$ADDON/bundle"

echo "deploying kapp controller"
kubectl create namespace tanzu-extensions || echo "namespace exists already"
kubectl apply -f https://raw.githubusercontent.com/vmware-tanzu/carvel-kapp-controller/dev-packaging/alpha-releases/v0.18.0-alpha.4.yml

echo "applying overlays and deploying dev package to current kubectl context..."
ytt -f addons/repos/main/packages/"$ADDON"* -f hack/addon-dev/addon-dev-overlay.yaml -f hack/addon-dev/values.yaml \
    --data-value dev_image="$TEMPLATE_IMAGE" | \
    kubectl apply -f -

kubectl apply -f addons/packages/"$ADDON"/installedpackage.yaml
kubectl apply -f addons/packages/"$ADDON"/*clusterrolebinding.yaml
kubectl apply -f addons/packages/"$ADDON"/*serviceaccount.yaml
