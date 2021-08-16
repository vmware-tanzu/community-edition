module github.com/vmware-tanzu-private/hack/builder

go 1.16

// -- Used to build TCE plugins with local codebase --
// replace github.com/vmware-tanzu/tanzu-framework => ../../../tanzu-framework

require (
	github.com/fabriziopandini/capi-conditions v0.0.0-20201102133039-7eb142d1b6d6 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/tanzu-framework v1.4.0-pre-alpha-2.0.20210817154238-f7e4fbdac647
)
