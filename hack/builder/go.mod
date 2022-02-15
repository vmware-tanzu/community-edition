module github.com/vmware-tanzu/community-edition/hack/builder

go 1.16

// -- Used to build TCE plugins with local codebase --
// replace github.com/vmware-tanzu/tanzu-framework => ../../../tanzu-framework

require (
	github.com/spf13/cobra v1.3.0
	github.com/vmware-tanzu/tanzu-framework v0.16.0
)

replace sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.0.1

replace github.com/k14s/kbld => github.com/anujc25/carvel-kbld v0.31.0-update-vendir
