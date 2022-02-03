module github.com/vmware-tanzu/community-edition/hack/builder

go 1.16

// -- Used to build TCE plugins with local codebase --
// replace github.com/vmware-tanzu/tanzu-framework => ../../../tanzu-framework

require (
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/golangplus/bytes v0.0.0-20160111154220-45c989fe5450 // indirect
	github.com/golangplus/fmt v0.0.0-20150411045040-2a5d6d7d2995 // indirect
	github.com/google/gopacket v1.1.17 // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/operator-framework/operator-sdk v0.0.7 // indirect
	github.com/spf13/cobra v1.3.0
	github.com/vmware-labs/yaml-jsonpath v0.3.2 // indirect
	github.com/vmware-tanzu/carvel-vendir v0.23.0 // indirect
	github.com/vmware-tanzu/tanzu-framework v0.16.0
	github.com/xlab/handysort v0.0.0-20150421192137-fb3537ed64a1 // indirect
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.13.0 // indirect
	sigs.k8s.io/cluster-api/test/infrastructure/docker v0.0.0-20210720023132-dfeb8d447bdc // indirect
	sigs.k8s.io/kustomize v2.0.3+incompatible // indirect
	vbom.ml/util v0.0.0-20160121211510-db5cfe13f5cc // indirect
)

replace sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v1.0.1

replace github.com/k14s/kbld => github.com/anujc25/carvel-kbld v0.31.0-update-vendir
