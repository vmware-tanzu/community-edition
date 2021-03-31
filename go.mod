module github.com/vmware-tanzu/tce

go 1.16

require (
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/adrg/xdg v0.3.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/joshrosso/image/v5 v5.10.2-0.20210331180716-71545f2b27af
	github.com/spf13/cobra v1.1.1
	github.com/vmware-tanzu-private/core v1.3.1-0.20210318014653-05bd34a5267a
	github.com/vmware-tanzu/carvel-kapp-controller v0.16.1-0.20210324160852-64ffdd7026ca
	github.com/vmware-tanzu/carvel-vendir v0.16.0
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	honnef.co/go/tools v0.1.3 // indirect
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	k8s.io/klog/v2 v2.4.0
	sigs.k8s.io/controller-runtime v0.7.0
)

replace github.com/containers/storage => github.com/joshrosso/storage v1.28.2-0.20210331182201-51e6dd05f861
