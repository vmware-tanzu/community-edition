// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

// GetAllUpCmds returns all commands used to run on tanzu or k8s
// these commands are w.r.t fluent-bit addon e2e only
func GetAllUpCmds() map[string][]string {
	return map[string][]string{
		"k8s-apply":             {"apply", "-f", "$"},
		"k8s-get-crds":          {"get", "crds"},
		"k8s-create-ns":         {"create", "ns", "$"},
		"k8s-port-forward":      {"port-forward", `daemonset/fluent-bit`, "-n", "fluent-bit", `56017:2020`},
		"k8s-daemonset-state":   {"get", "daemonset.apps/fluent-bit", "-n", "$", "-o", "$"},
		"tanzu-package-install": {"package", "install", "$"},

		"fluent-bit-health-check": {"-I", "http://127.0.0.1:56017/api/v1/health"},
	}
}

// GetTearDownCmds returns all commands used to tear-down.
// contains tanzu or k8s commands
// these commands are w.r.t fluent-bit addon e2e only
func GetTearDownCmds() map[string][]string {
	return map[string][]string{
		"k8s-delete":    {"delete", "-f", "$"},
		"k8s-delete-ns": {"delete", "ns", "$"},

		"tanzu-package-delete": {"package", "delete", "$"},
	}
}
