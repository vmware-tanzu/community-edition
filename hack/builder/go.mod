module github.com/vmware-tanzu-private/hack/builder

go 1.16

// -- Used to build TCE plugins with local codebase --
// replace github.com/vmware-tanzu-private/core => ../../../../vmware-tanzu-private/core

require (
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/tanzu-framework v1.4.0-pre-alpha-2.0.20210712202925-59ae6ee3afb1
)
