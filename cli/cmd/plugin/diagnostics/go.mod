module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/diagnostics

go 1.16

// We redirect this locally so go doesn't try to download a different copy
// of the same repo.
replace github.com/vmware-tanzu/community-edition => ../../../../

require (
	github.com/spf13/cobra v1.2.1
	github.com/vladimirvivien/gexe v0.1.1
	github.com/vmware-tanzu/community-edition v0.9.1
	github.com/vmware-tanzu/crash-diagnostics v0.3.7
	github.com/vmware-tanzu/tanzu-framework v0.10.1
	sigs.k8s.io/controller-runtime v0.9.0 // indirect
	sigs.k8s.io/kind v0.11.1
)
