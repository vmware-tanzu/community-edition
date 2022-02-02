module github.com/vmware-tanzu/community-edition/addons/packages/vsphere-cpi/1.23.0-alpha.1/test

go 1.16

require (
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
	github.com/vmware-tanzu/tanzu-framework/pkg/v1/providers/tests v0.0.0-20220113002410-90ed5669b49e
	k8s.io/api v0.23.1
	k8s.io/cloud-provider-vsphere v1.23.0-alpha.1
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../../../test/pkg
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.4.5
)
