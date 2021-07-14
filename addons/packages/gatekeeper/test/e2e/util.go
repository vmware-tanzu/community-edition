package e2e

// GetAllUpCmds returns all commands used to run on tanzu or kubectl
// these commands are w.r.t gatekeeper addon e2e only
func GetAllUpCmds() map[string][]string {
	return map[string][]string{
		"tanzu-package-install-gatekeeper":          []string{"package", "install", "gatekeeper.tce.vmware.com"},
		"tanzu-package-delete-gatekeeper":           []string{"package", "delete", "gatekeeper.tce.vmware.com"},
		"kubectl-get-pods-by-namespace":             []string{"get", "pods", "gatekeeper-system"},
		"kubectl-apply":                             []string{"apply", "-f", "$"},
		"kubectl-get-crds-constraint-template":      []string{"get", "crds"},
		"kubectl-create-ns":                         []string{"create", "ns", "$"},
		"kubeclt-check-pod-ready-status":            []string{"get", "pods", "-l", "$", "-n", "$", "-o", `jsonpath={..status.conditions[?(@.type=="Ready")].status}`},
		"kubeclt-check-gatekeeper-deployment-ready": []string{"get", "deployment", "-n", "gatekeeper-system", "-o", `jsonpath={.status.conditions[?(@.type == 'Available')].status}`},
	}
}

// GetTearDownCmds returns all commands used to tear-down.
// contains tanzu or kubectl commands
// these commands are w.r.t gatekeeper addon e2e only
func GetTearDownCmds() map[string][]string {
	return map[string][]string{
		"tanzu-package-delete-gatekeeper": []string{"package", "delete", "gatekeeper.tce.vmware.com"},
		"kubectl-delete-ns":               []string{"delete", "ns", "$"},
	}
}
