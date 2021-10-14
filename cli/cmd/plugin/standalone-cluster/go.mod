module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster

go 1.16

require (
	github.com/cppforlife/go-cli-ui v0.0.0-20200506005011-4268990983cc
	github.com/ghodss/yaml v1.0.0
	github.com/k14s/imgpkg v0.6.0
	github.com/k14s/ytt v0.32.1-0.20210511155130-214258be2519
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/carvel-kapp-controller v0.20.0-rc.1
	github.com/vmware-tanzu/carvel-vendir v0.19.0
	github.com/vmware-tanzu/tanzu-framework v1.4.0-pre-alpha-2.0.20210915174701-14fe0fdf4f0b
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v1.5.2
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/kind v0.11.1
)

replace (
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1 // indirect

	k8s.io/api => k8s.io/api v0.17.11
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.11
	k8s.io/client-go => k8s.io/client-go v0.17.11
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29

	sigs.k8s.io/cluster-api => github.com/vmware-tanzu/cluster-api v0.3.23-0.20210722162135-d31e78c28159
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.5.14
	sigs.k8s.io/kind => sigs.k8s.io/kind v0.11.1
)
