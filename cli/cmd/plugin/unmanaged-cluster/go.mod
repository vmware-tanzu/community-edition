module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster

go 1.16

replace (
	// This is only required until https://github.com/vmware-tanzu/carvel-ytt/issues/524 is resolved.
	// Until then, this patch should be carried forward to ensure ytt parsing can work on Windows.
	github.com/k14s/ytt => github.com/joshrosso/carvel-ytt v0.37.1-0.20211027005517-74085add68cc
	// We redirect this locally so go doesn't try to download a different copy
	// of the same repo.
	github.com/vmware-tanzu/community-edition => ../../../../
)

require (
	github.com/spf13/cobra v1.3.0
	github.com/vmware-tanzu/community-edition v0.9.1
)

require (
	github.com/cppforlife/go-cli-ui v0.0.0-20200716203538-1e47f820817f
	github.com/fatih/color v1.13.0
	github.com/k14s/imgpkg v0.6.0
	github.com/k14s/ytt v0.37.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/spf13/pflag v1.0.5
	github.com/vmware-tanzu/carvel-kapp-controller v0.28.0
	github.com/vmware-tanzu/carvel-kbld v0.32.1-0.20220207174123-dd5e71b95085
	github.com/vmware-tanzu/carvel-vendir v0.24.0
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.23.0
	k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery v0.23.2
	k8s.io/client-go v0.23.0
	k8s.io/kube-aggregator v0.19.2
	sigs.k8s.io/kind v0.11.1
)
