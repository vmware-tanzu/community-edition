module github.com/vmware-tanzu/community-edition/hack/runner/webhook

go 1.17

require (
	github.com/aws/aws-sdk-go v1.41.1
	github.com/go-playground/webhooks/v6 v6.0.0-rc.1
	github.com/google/go-github/v39 v39.1.0
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1
	k8s.io/klog/v2 v2.50.2
)

require (
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20211011170408-caeb26a5c8c0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9 // indirect
)
