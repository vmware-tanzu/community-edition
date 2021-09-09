module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/standalone-cluster

go 1.16

require (
	github.com/spf13/cobra v1.2.0
	github.com/vmware-tanzu/tanzu-framework v1.4.0-pre-alpha-2.0.20210909135501-d143f231734a
	k8s.io/klog/v2 v2.8.0
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
