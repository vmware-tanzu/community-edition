module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster

go 1.16

require (
	github.com/cppforlife/go-cli-ui v0.0.0-20200716203538-1e47f820817f
	github.com/fatih/color v1.13.0
	github.com/k14s/imgpkg v0.6.0
	github.com/k14s/ytt v0.37.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/vmware-tanzu/carvel-kapp-controller v0.28.0
	github.com/vmware-tanzu/carvel-vendir v0.23.0
	github.com/vmware-tanzu/tanzu-framework v0.10.0
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.21.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/kube-aggregator v0.19.2
	// This indirect dependency should be removed
	// once the tanzu-framework dependency has been decoupled:
	// https://github.com/vmware-tanzu/community-edition/issues/2811
	sigs.k8s.io/controller-runtime v0.9.0 // indirect
	sigs.k8s.io/kind v0.11.1
)

// This is only required until https://github.com/vmware-tanzu/carvel-ytt/issues/524 is resolved.
// Until then, this patch should be carried forward to ensure ytt parsing can work on Windows.
replace github.com/k14s/ytt => github.com/joshrosso/carvel-ytt v0.37.1-0.20211027005517-74085add68cc
