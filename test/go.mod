module github.com/vmware-tanzu/community-edition/test

go 1.16

require (
	github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96 // indirect
	github.com/gophercloud/gophercloud v0.0.0-20190126172459-c818fa66e4c8 // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3 // indirect
	k8s.io/client-go v0.21.3
	k8s.io/klog v0.3.1 // indirect
	sigs.k8s.io/e2e-framework v0.0.2
)

replace github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../addons/packages/test/pkg
