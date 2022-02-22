module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/conformance

go 1.16

// We redirect this locally so go doesn't try to download a different copy
// of the same repo.
replace github.com/vmware-tanzu/community-edition => ../../../../

require (
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/community-edition v0.9.1
	github.com/vmware-tanzu/sonobuoy v0.53.2
)
