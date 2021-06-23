module github.com/vmware-tanzu/tce

go 1.16

require (
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/adrg/xdg v0.3.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/joshrosso/image/v5 v5.10.2-0.20210331180716-71545f2b27af
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/vmware-tanzu-private/core v1.3.0
	github.com/vmware-tanzu-private/tkg-cli v1.3.0
	github.com/vmware-tanzu/carvel-kapp-controller v0.19.1-0.20210422224550-3c235246c149
	github.com/vmware-tanzu/carvel-vendir v0.19.0
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	honnef.co/go/tools v0.1.3 // indirect
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	k8s.io/klog/v2 v2.4.0
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920
	sigs.k8s.io/controller-runtime v0.5.14
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	github.com/containers/storage => github.com/joshrosso/storage v1.28.2-0.20210331182201-51e6dd05f861
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0 // indirect

	// toggle this between local for development (just set these dev paths based on GOPATH)
	// github.com/vmware-tanzu-private/core => ../../vmware-tanzu-private/core
	github.com/vmware-tanzu-private/core => github.com/vmware-tanzu-private/core v1.3.1-0.20210623184735-f219aaef1e1d
	// github.com/vmware-tanzu-private/tanzu-cli-tkg-plugins => ../../vmware-tanzu-private/tanzu-cli-tkg-plugins
	github.com/vmware-tanzu-private/tkg-cli => github.com/vmware-tanzu-private/tkg-cli v1.3.1-0.20210623184308-09830a241cd5
	// github.com/vmware-tanzu-private/tkg-cli => ../../vmware-tanzu-private/tkg-cli
	// github.com/vmware-tanzu-private/tkg-providers => ../../vmware-tanzu-private/tkg-providers
	github.com/vmware-tanzu-private/tkg-providers => github.com/vmware-tanzu-private/tkg-providers v1.3.1-0.20210422215837-027482ef8765

	k8s.io/api => k8s.io/api v0.17.11
	k8s.io/apimachinery => github.com/joshrosso/apimachinery v0.17.12-rc.0.0.20210402165939-550cad781ca6
	k8s.io/client-go => github.com/joshrosso/client-go v0.17.12-0.20210402163626-0c304a130f6a
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.4.0
	k8s.io/kubectl => k8s.io/kubectl v0.17.11

	// toggle this between local for development (just set these dev paths based on GOPATH)
	sigs.k8s.io/cluster-api => github.com/vmware-tanzu/cluster-api v0.3.15-0.20210609222148-e9e6c9d422e8
	// sigs.k8s.io/cluster-api => ../../vmware-tanzu/cluster-api

	sigs.k8s.io/cluster-api-provider-azure => sigs.k8s.io/cluster-api-provider-azure v0.4.11
)
