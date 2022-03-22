#!/bin/bash

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

kubectl wait --timeout=10m --for=jsonpath='{.status.conditions[].type}'=ReconcileSucceeded packageinstalls/cartographer -n tanzu-package-repo-global
kubectl wait --timeout=10m --for=jsonpath='{.status.conditions[].type}'=ReconcileSucceeded packageinstalls/cert-manager -n tanzu-package-repo-global
kubectl wait --timeout=10m --for=jsonpath='{.status.conditions[].type}'=ReconcileSucceeded packageinstalls/contour -n tanzu-package-repo-global
kubectl wait --timeout=10m --for=jsonpath='{.status.conditions[].type}'=ReconcileSucceeded packageinstalls/fluxcd-source-controller -n tanzu-package-repo-global
kubectl wait --timeout=10m --for=jsonpath='{.status.conditions[].type}'=ReconcileSucceeded packageinstalls/knative-serving -n tanzu-package-repo-global
kubectl wait --timeout=10m --for=jsonpath='{.status.conditions[].type}'=ReconcileSucceeded packageinstalls/kpack -n tanzu-package-repo-global
kubectl wait --timeout=10m --for=jsonpath='{.status.conditions[].type}'=ReconcileSucceeded packageinstalls/app-toolkit -n tanzu-package-repo-global
