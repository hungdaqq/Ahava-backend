SHELL := /bin/bash

.PHONY: all build test deps deps-cleancache

GOCMD=go
BUILD_DIR=build
BINARY_DIR=$(BUILD_DIR)/bin
CODE_COVERAGE=code-coverage

all: test build

${BINARY_DIR}:
	mkdir -p $(BINARY_DIR)

build: ${BINARY_DIR} ## Compile the code, build Executable File
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o $(BINARY_DIR) -v ./cmd/api
	scp ./build/bin/api ahavavps:/home/ahava/app/api

run: ## Start application
	$(GOCMD) run ./cmd/api

test: ## Run tests
	$(GOCMD) test ./... -cover

test-coverage: ## Run tests and generate coverage file
	$(GOCMD) test ./... -coverprofile=$(CODE_COVERAGE).out
	$(GOCMD) tool cover -html=$(CODE_COVERAGE).out

deps: ## Install dependencies
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	$(GOCMD) get -u -t -d -v ./...
	$(GOCMD) mod tidy
	$(GOCMD) mod vendor

deps-cleancache: ## Clear cache in Go module
	$(GOCMD) clean -modcache

wire: ## Generate wire_gen.go
	cd pkg/di && wire

mock: ##make mock files using mockgen
	mockgen -source=pkg/repository/interface/user.go -destination=pkg/mock/mockrepo/user_mock.go -package=mockrepo
	mockgen -source=pkg/service/interface/user.go -destination=pkg/mock/mockservice/user_mock.go -package=mockservice
	mockgen -source=pkg/repository/interface/product.go -destination=pkg/mock/mockrepo/product_mock.go -package=mockrepo
	mockgen -source=pkg/repository/interface/order.go -destination=pkg/mock/mockrepo/order_mock.go -package=mockrepo

swag: ## Generate swagger docs
		swag init -g pkg/api/handler/admin.go -o ./cmd/api/docs

lint: ## for linting go code
		golangci-lint run ./...

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

docker:
	docker build -t hungdaqq/ahava-backend:1.0.0 .