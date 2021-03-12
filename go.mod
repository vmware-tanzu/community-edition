module github.com/vmware-tanzu/tce

go 1.16

require (
	cloud.google.com/go/storage v1.12.0
	github.com/adrg/xdg v0.3.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/vmware-tanzu-private/core v1.3.0-rc.2.0.20210225204548-11aadc74602a
	github.com/vmware-tanzu/carvel-kapp-controller v0.16.1-0.20210311153426-867eef4d6913
	github.com/vmware-tanzu/carvel-vendir v0.16.0
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
	google.golang.org/api v0.36.0
	honnef.co/go/tools v0.1.3 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog/v2 v2.4.0
	sigs.k8s.io/controller-runtime v0.7.0
)

replace (
	github.com/vmware-tanzu-private/tkg-cli => github.com/vmware-tanzu-private/tkg-cli v1.3.0-pre-alpha-1.0.20210114003033-285a8c9131d4
	github.com/vmware-tanzu-private/tkg-providers => github.com/vmware-tanzu-private/tkg-providers v1.3.0-pre-alpha-1.0.20210113202657-eb07b4e0558d
)
