module github.com/vmware-tanzu/community-edition/addons/packages/harbor/test

go 1.16

require (
	github.com/containerd/containerd v1.5.7
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/opencontainers/image-spec v1.0.1
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	oras.land/oras-go v0.4.0
	rsc.io/letsencrypt v0.0.3 // indirect
)

replace github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../../../test/pkg
