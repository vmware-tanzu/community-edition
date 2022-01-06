// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kapp

// This is entirely a workaround until we've got better plumbing

const (
	DefaultKappValues = `---
namespace: kapp-controller
kappController:
  namespace: tkg-system
  createNamespace: true
  globalNamespace: tanzu-package-repo-global
  deployment:
    hostNetwork: true
    priorityClassName: null
    concurrency: 4
    tolerations: 
    - key: "node.kubernetes.io/not-ready"
      operator: "Exists"
      effect: "NoSchedule"
    apiPort: 10400
  config:
    caCerts: ""
    httpProxy: ""
    httpsProxy: ""
    noProxy: ""
    dangerousSkipTLSVerify: ""
`
)
