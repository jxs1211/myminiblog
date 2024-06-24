# Default to execute the 'all' target
.DEFAULT_GOAL := all

# ==============================================================================
# Define the Makefile all phony target, which will be executed by default when `make` is run
.PHONY: all
all: gen.add-copyright go.format go.lint go.cover go.build

# ==============================================================================
# Includes

# Ensure `include common.mk` is on the first line, common.mk defines some variables that subsequent sub-makefiles depend on
include scripts/make-rules/common.mk 
include scripts/make-rules/tools.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/generate.mk
include scripts/make-rules/image.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build
                   Example: make build BINS="miniblog test"
  IMAGES           Backend images to make. Default is all of cmd.
                   This option is available when using: make image/push
                   Example: make image IMAGES="miniblog"
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

## --------------------------------------
## Generate / Manifests
## --------------------------------------

##@ generate:

.PHONY: add-copyright
add-copyright: ## Add copyright header information.
	@$(MAKE) gen.add-copyright

.PHONY: ca
ca: ## Generate CA files.
	@$(MAKE) gen.ca

.PHONY: protoc
protoc: ## Compile protobuf files.
	@$(MAKE) gen.protoc

.PHONY: deps
deps: ## Install dependencies, such as generating required code, installing necessary tools, etc.
	@$(MAKE) gen.deps

## --------------------------------------
## Binaries
## --------------------------------------

##@ build:

.PHONY: build
build: go.tidy  ## Compile source code, depends on the tidy target to automatically add/remove dependencies.
	@$(MAKE) go.build

.PHONY: image
image: ## Build Docker images.
	@$(MAKE) image.build

.PHONY: push
push: ## Build Docker images and push to the image repository.
	@$(MAKE) image.push

## --------------------------------------
## Cleanup
## --------------------------------------

##@ clean:

.PHONY: clean
clean: ## Clean up build artifacts, temporary files, etc.
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)


## --------------------------------------
## Lint / Verification
## --------------------------------------

##@ lint and verify:

.PHONY: lint
lint: ## Perform static code analysis.
	@$(MAKE) go.lint


## --------------------------------------
## Testing
## --------------------------------------

##@ test:

.PHONY: test 
test: ## Run unit tests.
	@$(MAKE) go.test

.PHONY: cover 
cover: ## Run unit tests and check coverage thresholds.
	@$(MAKE) go.cover


## --------------------------------------
## Hack / Tools
## --------------------------------------

##@ hack/tools:

.PHONY: format
format:  ## Format Go source code.
	@$(MAKE) go.format

.PHONY: swagger
swagger: tools.verify.swagger ## Start swagger online documentation (listening port: 65534).
	@swagger serve -F=swagger --no-open --port 65534 $(ROOT_DIR)/api/openapi/openapi.yaml

.PHONY: tidy
tidy: ## Automatically add/remove dependencies.
	@$(MAKE) go.tidy

.PHONY: help
help: Makefile ## Print Makefile help information.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGETS> <OPTIONS>\033[0m\n\n\033[35mTargets:\033[0m\n"} /^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' Makefile #$(MAKEFILE_LIST)
	@echo -e "$$USAGE_OPTIONS"