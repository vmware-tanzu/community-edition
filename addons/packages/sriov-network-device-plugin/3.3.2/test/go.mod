module github.com/vmware-tanzu/community-edition/addons/packages/sriov-network-device-plugin/3.3.2/test

go 1.16

require (
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/instrumenta/kubeval v0.16.1
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20211105192438-b53810dc28af // indirect
	golang.org/x/sys v0.0.0-20211107104306-e0b2ad06fe42 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../../../test/pkg
