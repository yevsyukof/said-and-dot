APP_NAME=snd-backend
VERSION?=v0.1.0
#BUILD=$(shell git rev-parse HEAD)

GO=go
LDFLAGS=-ldflags="-s -w -X"
BINARY_DIR=.bin

ENTRYPOINT=./services/monolith/cmd/app/main.go
ARTIFACT_NAME=snd-backend

#GOOSS=darwin linux windows freebsd netbsd openbsd dragonfly
#GOARCHS=386 arm arm64 amd64

.DEFAULT_GOAL := help

.PHONY: help
# Source: https://habr.com/ru/post/211751/
help: ## Displays all the available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Format go files
	@$(GO) fmt ./...

.PHONY: vet
vet: ## go vet
	@$(GO) vet ./...

.PHONY: test
test: ## Runs unit tests [cmd: go test -v -bench . -benchmem ./...]
	@$(GO) test -v -bench . -benchmem ./...

.PHONY: clean
clean: ## Deletes all compiled / executable files
	@find .bin -type f -name '*' -print0 | xargs -0 rm --
	@echo ">> Deleted all build files!"

.PHONY: install
install: ## Installs the package
	@$(GO) install ${LDFLAGS} ./...

.PHONY: install-deps
install-deps: ## Install dependencies
	@$(GO) mod download

.PHONY: run
run: ## Runs the backend server
	@$(GO) run $(ENTRYPOINT)

.PHONY: dev
dev: ## Runs the backend server with hot-reload (Must have air installed https://github.com/cosmtrek/air)
	@air -c configs/.air.toml

.PHONY: artifact-name
artifact-name: ## Return 'make build' artifact name
	@echo $(ARTIFACT_NAME)

.PHONY: build
build: ## Compiles the rest-api server
	@$(GO) build -o $(ARTIFACT_NAME) $(ENTRYPOINT)

.PHONY: build-rest
build-rest: ## Compiles the rest-api server
	@$(GO) build $(LDFLAGS) -v -o $(BINARY_DIR)/$(APP_NAME)-$(VERSION)_server $(ENTRYPOINT)

#.PHONY: build-rest-all
#build-rest-all: ## Cross-compiles the rest api server
#	@$(foreach GOOS, $(GOOSS),\
#		$(foreach GOARCH, $(GOARCHS),\
#			$(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH);\
#			$(GO) build $(LDFLAGS) -v -o $(BINARY_DIR)/$(APP_NAME)-$(VERSION)_server-$(GOOS)-$(GOARCH) $(ENTRYPOINT))))
