module github.com/vmware-tanzu/community-edition/addons/packages/load-balancer-and-ingress-service/1.6.1/test

go 1.17

require (
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
	github.com/vmware-tanzu/tanzu-framework/pkg/v1/providers/tests v0.0.0-20220113002410-90ed5669b49e
	k8s.io/api v0.23.1
	sigs.k8s.io/yaml v1.3.0
)

require (
	golang.org/x/sys v0.0.0-20211029165221-6e7872819dc8 // indirect
	golang.org/x/tools v0.1.8 // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.0 // indirect
)

replace (
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../../../test/pkg
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.4.5
)
