SHELL:=/bin/bash

FIX=1
COVER=0
CGO_ENABLED=1
TESTARGS=-race
COVERARGS=-cover -covermode=atomic -coverprofile=coverage.out
CMDS=cmd/demo/demo
BUILDARGS=
PKGSOURCES:=$(shell find ../../pkg -type f -name '*.go' 2>/dev/null)
INTERNALSOURCES:=$(shell find ../../internal -type f -name '*.go' 2>/dev/null)

default: help

.PHONY: build
build: $(CMDS) ## Build Go binaries

cmd/demo/demo: $(shell find cmd/demo -type f -name '*.go') $(PKGSOURCES) $(INTERNALSOURCES)
	cd `dirname $@` && go build $(BUILDARGS) -o `basename $@`

.PHONY: gofmt
gofmt:
	@if test "$(FIX)" = "1"; then \
		set -x ; gofmt -s -w . ;\
	else \
		set -x ; gofmt -s -d . ;\
	fi

.PHONY: golangcilint
golangcilint: tmp/bin/golangci-lint
	@if test "$(FIX)" = "1"; then \
		set -x ; $< run --fix --timeout 10m;\
	else \
		set -x ; $< run --timeout 10m;\
	fi

.PHONY: no-dirty
no-dirty: ## Test if there are some dirty files
	git diff --exit-code

.PHONY: govet
govet:
	go vet ./...

.PHONY: gomodtidy
gomodtidy:
	go mod tidy -v

.PHONY: unit_test
unit_test: ## Execute all unit tests
	@if test "$(COVER)" = "1"; then \
		go test $(TESTARGS) $(COVERARGS) ./...;\
	else \
		go test $(TESTARGS) ./...;\
	fi

.PHONY: test
test: unit_test ## Execute all tests 

.PHONY: html-coverage
html-coverage: ## Build html coverage
	$(MAKE) COVER=1 test
	go tool cover -html coverage.out -o cover.html

.PHONY: lint
lint: govet gofmt golangcilint ## Lint the code (also fix the code if FIX=1, default)

tmp/bin/golangci-lint:
	@mkdir -p tmp/bin
	cd tmp/bin && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b . v1.56.2 && chmod +x `basename $@`

.PHONY: clean
clean: _cmd_clean ## Clean the repo
	rm -f coverage.out cover.html
	rm -Rf tmp build

.PHONY: _cmd_clean
_cmd_clean:
	rm -f $(CMDS)

.PHONY: help
help::
	@# See https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
