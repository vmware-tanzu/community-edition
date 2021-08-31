.DEFAULT_GOAL:=help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

lint: ## Run Go code linting
ifeq ($(origin GOLANGCI_LINT),undefined)
	@echo "Error! GOLANGCI_LINT env var not set"
else
	$(GOLANGCI_LINT) run -v --timeout=5m
endif

get-deps: ## Get all dependencies
	go mod download

test: ## Run unit testing suite
	echo "N/A: No unit tests for addons/packages/pinniped/hack"

e2e-test: ## Run e2e testing suite
	echo "N/A: No e2e tests for addons/packages/pinniped/hack"

build: ## Build the executable
	@echo "N/A No binary to build for addons/packages/pinniped/hack"

