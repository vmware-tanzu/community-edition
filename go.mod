module github.com/vmware-tanzu/tce

go 1.16

require (
	github.com/Microsoft/go-winio v0.4.17-0.20210211115548-6eac466e5fa3 // indirect
	github.com/containerd/containerd v1.5.0-beta.1 // indirect
	github.com/docker/docker v1.4.2-0.20191219165747-a9416c67da9f // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/zapr v0.4.0 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.13.0
	github.com/opencontainers/image-spec v1.0.2-0.20190823105129-775207bd45b6 // indirect
	github.com/spf13/cobra v1.2.0
	github.com/vmware-tanzu/tanzu-framework v1.4.0-pre-alpha-2.0.20210712202925-59ae6ee3afb1
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.21.2
	k8s.io/klog/v2 v2.8.0
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1 // indirect
	k8s.io/api => k8s.io/api v0.17.11
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.11
	k8s.io/client-go => k8s.io/client-go v0.17.11
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29

	sigs.k8s.io/cluster-api => github.com/vmware-tanzu/cluster-api v0.3.15-0.20210609222148-e9e6c9d422e8
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.5.14
	sigs.k8s.io/kind => sigs.k8s.io/kind v0.11.1
)
