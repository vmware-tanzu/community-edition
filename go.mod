module github.com/vmware-tanzu/tce

go 1.15

require (
	cloud.google.com/go/storage v1.12.0
	github.com/adrg/xdg v0.3.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/vmware-tanzu-private/core v1.3.0-rc.2.0.20210225181944-795041f9add9
	github.com/vmware-tanzu/carvel-kapp-controller v0.13.0
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
	golang.org/x/tools v0.1.0 // indirect
	google.golang.org/api v0.36.0
	honnef.co/go/tools v0.1.2 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog/v2 v2.4.0
	sigs.k8s.io/controller-runtime v0.7.0
)

replace (
	github.com/vmware-tanzu-private/tkg-cli => github.com/vmware-tanzu-private/tkg-cli v1.3.0-pre-alpha-1.0.20210114003033-285a8c9131d4
	github.com/vmware-tanzu-private/tkg-providers => github.com/vmware-tanzu-private/tkg-providers v1.3.0-pre-alpha-1.0.20210113202657-eb07b4e0558d
	github.com/vmware-tanzu/carvel-kapp-controller => github.com/alexbrand/carvel-kapp-controller v0.13.1-0.20210127180239-bf57d388b9b3
)
