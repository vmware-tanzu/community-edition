module github.com/vmware-tanzu-private/hack/builder

go 1.16

// -- Used to build TCE plugins with local codebase --
// replace github.com/vmware-tanzu/tanzu-framework => ../../../tanzu-framework

require (
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/tanzu-framework v1.4.0-pre-alpha-2.0.20210818150133-ab306729ecd3
)
