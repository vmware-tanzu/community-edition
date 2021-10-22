module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster

go 1.16

require (
	github.com/cppforlife/go-cli-ui v0.0.0-20200716203538-1e47f820817f
	github.com/k14s/imgpkg v0.6.0
	github.com/k14s/ytt v0.37.0
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/carvel-kapp-controller v0.28.0
	github.com/vmware-tanzu/carvel-vendir v0.23.0
	github.com/vmware-tanzu/tanzu-framework v1.4.0-pre-alpha-2.0.20210915174701-14fe0fdf4f0b
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.21.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/kube-aggregator v0.19.2
	sigs.k8s.io/kind v0.11.1
)
