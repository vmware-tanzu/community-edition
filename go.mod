module github.com/vmware-tanzu/tce

go 1.16

require (
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/adrg/xdg v0.3.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/joshrosso/image/v5 v5.10.2-0.20210331180716-71545f2b27af
	github.com/spf13/cobra v1.1.1
	github.com/vmware-tanzu-private/core v1.3.0
	github.com/vmware-tanzu-private/tkg-cli v1.3.0
	github.com/vmware-tanzu/carvel-kapp-controller v0.16.1-0.20210324160852-64ffdd7026ca
	github.com/vmware-tanzu/carvel-vendir v0.16.0
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	golang.org/x/tools v0.1.0 // indirect
	honnef.co/go/tools v0.1.3 // indirect
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	k8s.io/klog/v2 v2.4.0
	sigs.k8s.io/controller-runtime v0.5.14
)

replace (
	github.com/containers/storage => github.com/joshrosso/storage v1.28.2-0.20210331182201-51e6dd05f861
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0 // indirect

	//github.com/vmware-tanzu-private/core => github.com/vmware-tanzu-private/core v1.3.0-rc.1.0.20210415014539-3f9cb357b7e9
	// toggle this between remote or local for development
	github.com/vmware-tanzu-private/tkg-cli => github.com/vmware-tanzu-private/tkg-cli v1.3.0-rc.1.0.20210422001449-797adde0cba1
	//github.com/vmware-tanzu-private/tkg-cli => /home/josh/d/tkg-cli
	github.com/vmware-tanzu-private/tkg-providers => github.com/vmware-tanzu-private/tkg-providers v1.3.0-rc.1.0.20210421170839-1c2d40890b38

	k8s.io/api => k8s.io/api v0.17.11
	k8s.io/apimachinery => github.com/joshrosso/apimachinery v0.17.12-rc.0.0.20210402165939-550cad781ca6
	k8s.io/client-go => github.com/joshrosso/client-go v0.17.12-0.20210402163626-0c304a130f6a
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.4.0
	k8s.io/kubectl => k8s.io/kubectl v0.17.11
	sigs.k8s.io/cluster-api-provider-azure => sigs.k8s.io/cluster-api-provider-azure v0.4.11
)
