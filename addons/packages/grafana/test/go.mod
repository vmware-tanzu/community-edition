module github.com/vmware-tanzu/tce/addons/packages/grafana/test

go 1.16

replace github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../../test/pkg

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
)
