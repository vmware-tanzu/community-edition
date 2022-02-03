module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/diagnostics

go 1.16

// We redirect this locally so go doesn't try to download a different copy
// of the same repo.
replace github.com/vmware-tanzu/community-edition => ../../../../

require (
	github.com/spf13/cobra v1.3.0
	github.com/vladimirvivien/gexe v0.1.1
	github.com/vmware-tanzu/community-edition v0.9.1
	github.com/vmware-tanzu/crash-diagnostics v0.3.7
	github.com/vmware-tanzu/tanzu-framework v0.16.0
	sigs.k8s.io/kind v0.11.1
)

replace sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.0.1

replace github.com/k14s/kbld => github.com/anujc25/carvel-kbld v0.31.0-update-vendir
