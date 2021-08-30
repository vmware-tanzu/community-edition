module github.com/vmware-tanzu/community-edition/hack/builder

go 1.16

// -- Used to build TCE plugins with local codebase --
// replace github.com/vmware-tanzu/tanzu-framework => ../../../tanzu-framework

require (
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/tanzu-framework v1.4.0-pre-alpha-2.0.20210819000359-75cfa0e3ada3
	sigs.k8s.io/controller-runtime v0.9.0 // indirect
)
