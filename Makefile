SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

PROJECT_NAME = cached

.PHONY: help
help: ## View help information
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Cleans up artifacts
	docker-compose -p $(PROJECT_NAME) -f hack/docker-compose.yml down --remove-orphans

.PHONY: dev
dev: ## Runs a local development environment
	docker-compose -p $(PROJECT_NAME) -f hack/docker-compose.yml up
