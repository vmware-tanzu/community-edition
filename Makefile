.DEFAULT_GOAL:=help

OCI_REGISTRY := projects.registry.vmware.com/tce

help: ## display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

deploy-kapp-controller: ## deploys the latest version of kapp-controller
	kubectl create ns kapp-controller || true
	kubectl -n kapp-controller apply -f https://github.com/k14s/kapp-controller/releases/latest/download/release.yml


push-extensions: ## build and push extension templates
	imgpkg push --bundle $(OCI_REGISTRY)/velero-extension-templates:dev --file extensions/velero/bundle/


redeploy-velero: ## delete and redeploy the velero extension
	kubectl --namespace tanzu-extensions delete app velero
	kubectl apply --filename extensions/velero/extension.yaml
