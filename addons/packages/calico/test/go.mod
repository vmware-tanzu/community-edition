module github.com/vmware-tanzu/community-edition/addons/packages/calico/test

go 1.16

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
	k8s.io/api v0.22.1
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../../test/pkg
