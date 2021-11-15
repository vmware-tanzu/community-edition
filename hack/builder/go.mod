module github.com/vmware-tanzu/community-edition/hack/builder

go 1.16

// -- Used to build TCE plugins with local codebase --
// replace github.com/vmware-tanzu/tanzu-framework => ../../../tanzu-framework

require (
	github.com/aunum/log v0.0.0-20200821225356-38d2e2c8b489 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/docker/docker v1.4.2-0.20190924003213-a8608b5b67c7 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/tanzu-framework v0.10.0
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	sigs.k8s.io/controller-runtime v0.9.0 // indirect
)
