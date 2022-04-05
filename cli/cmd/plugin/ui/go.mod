module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui

go 1.17

replace (
	// We redirect this locally so go doesn't try to download a different copy
	// of the same repo.
	github.com/vmware-tanzu/community-edition => ../../../../
	// We need to do this for tanzu-framework for now.
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.1.3
)

require (
	github.com/go-openapi/errors v0.20.2
	github.com/go-openapi/loads v0.21.1
	github.com/go-openapi/runtime v0.23.3
	github.com/go-openapi/spec v0.20.4
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-openapi/swag v0.21.1
	github.com/go-openapi/validate v0.21.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/community-edition v0.9.1
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-openapi/analysis v0.21.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.1-0.20210427113832-6241f9ab9942 // indirect
	go.mongodb.org/mongo-driver v1.8.3 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
