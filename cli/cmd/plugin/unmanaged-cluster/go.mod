module github.com/vmware-tanzu/community-edition/cli/cmd/plugin/unmanaged-cluster

go 1.17

replace (
	// This is only required until https://github.com/vmware-tanzu/carvel-ytt/issues/524 is resolved.
	// Until then, this patch should be carried forward to ensure ytt parsing can work on Windows.
	github.com/k14s/ytt => github.com/joshrosso/carvel-ytt v0.37.1-0.20211027005517-74085add68cc
	// We redirect this locally so go doesn't try to download a different copy
	// of the same repo.
	github.com/vmware-tanzu/community-edition => ../../../../
)

require (
	github.com/spf13/cobra v1.5.0
	github.com/vmware-tanzu/community-edition v0.9.1
)

require (
	github.com/cppforlife/go-cli-ui v0.0.0-20220428182907-73db60c7611a
	github.com/fatih/color v1.13.0
	github.com/google/go-containerregistry v0.11.0
	github.com/k14s/imgpkg v0.6.0
	github.com/k14s/ytt v0.37.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/vmware-tanzu/carvel-kapp-controller v0.38.4
	github.com/vmware-tanzu/carvel-kbld v0.34.0
	github.com/vmware-tanzu/carvel-vendir v0.30.0
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.25.2
	k8s.io/apiextensions-apiserver v0.24.4
	k8s.io/apimachinery v0.25.2
	k8s.io/client-go v0.25.2
	k8s.io/kube-aggregator v0.25.2
	sigs.k8s.io/kind v0.14.0
)

require (
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.12.0 // indirect
	github.com/cppforlife/cobrautil v0.0.0-20200514214827-bb86e6965d72 // indirect
	github.com/cppforlife/color v1.9.1-0.20200716202919-6706ac40b835 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/cli v20.10.17+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k14s/semver/v4 v4.0.1-0.20210701191048-266d47ac6115 // indirect
	github.com/k14s/starlark-go v0.0.0-20200720175618-3a5c849cc368 // indirect
	github.com/klauspost/compress v1.15.8 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/vbatts/tar-split v0.11.2 // indirect
	github.com/vito/go-interact v1.0.1 // indirect
	github.com/vmware-tanzu/carvel-imgpkg v0.28.0 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/oauth2 v0.0.0-20220718184931-c8730f7fcb92 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
