module github.com/chenliu1993/community-edition/addons/packages/sriov-cni/2.6.1/test

go 1.16

require (
	github.com/instrumenta/kubeval v0.16.1
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
)

replace github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../../../test/pkg
