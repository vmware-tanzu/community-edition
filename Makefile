.DEFAULT_GOAL:=help

OCI_REGISTRY := projects.registry.vmware.com/tce

help: ## display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

deploy-kapp-controller: ## deploys the latest version of kapp-controller
	kubectl create ns kapp-controller || true
	kubectl -n kapp-controller apply -f ./release.yml

push-extension-velero: ## build and push extension templates
	imgpkg push --image $(OCI_REGISTRY)/velero-extension-templates -f extensions/velero/config/

push-extension-gatekeeper: ## build and push extension templates
	imgpkg push --bundle $(OCI_REGISTRY)/gatekeeper-extension-templates -f extensions/gatekeeper/bundle/config/

update-image-lockfiles: ## updates the ImageLock files in each extension
	kbld -f extensions/gatekeeper/bundle --imgpkg-lock-output extensions/gatekeeper/bundle/.imgpkg/images.yml

redeploy-gatekeeper: ## delete and redeploy the velero extension
	kubectl -n tanzu-extensions delete app gatekeeper || true
	kubectl apply -f extensions/gatekeeper/extension.yaml
