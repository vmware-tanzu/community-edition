module github.com/vmware-tanzu/community-edition/test

go 1.16

require (
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.21.3
	k8s.io/client-go v0.21.3
	k8s.io/utils v0.0.0-20210820185131-d34e5cb4466e // indirect
	sigs.k8s.io/e2e-framework v0.0.4
)

replace github.com/vmware-tanzu/community-edition/addons/packages/test/pkg => ../addons/packages/test/pkg
