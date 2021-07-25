module github.com/vmware-tanzu/tce/addons/packages/prometheus/test/e2e

go 1.16

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/vmware-tanzu/tce/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
)

replace github.com/vmware-tanzu/tce/addons/packages/test/pkg => ../../../test/pkg
