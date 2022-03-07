module github.com/vmware-tanzu/community-edition/addons/packages/sriov-network-device-plugin/3.3.2/test

go 1.17

require (
	github.com/instrumenta/kubeval v0.16.1
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.18.1
	github.com/vmware-tanzu/community-edition/addons/packages/test/pkg v0.0.0-00010101000000-000000000000
)

require (
	github.com/fatih/color v1.10.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/spf13/cobra v0.0.0-20180820174524-ff0d02e85550 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v0.0.0-20180816142147-da425ebb7609 // indirect
	golang.org/x/net v0.0.0-20211105192438-b53810dc28af // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/utils v0.0.0-20210820185131-d34e5cb4466e // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../../../test/pkg
