package e2e

func GetClusterUpCmds() map[string][]string {
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

func GetClusterDownCmds() map[string][]string {
	return map[string][]string{
		"tanzu-package-delete-gatekeeper": []string{"package", "delete", "gatekeeper.tce.vmware.com"},
		"kubectl-delete-ns":               []string{"delete", "ns", "$"},
	}
}
