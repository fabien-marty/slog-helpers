SHELL:=/bin/bash

FIX=1
COVER=0
CGO_ENABLED=1
TESTARGS=-race
COVERARGS=-cover -covermode=atomic -coverprofile=coverage.out
CMDS=cmd/slogc-demo1/slogc-demo1 cmd/stacktrace-demo1/stacktrace-demo1 cmd/human-demo1/human-demo1 cmd/external-demo1/external-demo1
BUILDARGS=
PKGSOURCES:=$(shell find pkg -type f -name '*.go' 2>/dev/null)
INTERNALSOURCES:=$(shell find internal -type f -name '*.go' 2>/dev/null)

default: help

.PHONY: build
build: $(CMDS) ## Build Go binaries

cmd/slogc-demo1/slogc-demo1: $(shell find cmd/slogc-demo1 -type f -name '*.go') $(PKGSOURCES) $(INTERNALSOURCES)
	cd `dirname $@` && go build $(BUILDARGS) -o `basename $@`

cmd/stacktrace-demo1/stacktrace-demo1: $(shell find cmd/stacktrace-demo1 -type f -name '*.go') $(PKGSOURCES) $(INTERNALSOURCES)
	cd `dirname $@` && go build $(BUILDARGS) -o `basename $@`

cmd/human-demo1/human-demo1: $(shell find cmd/human-demo1 -type f -name '*.go') $(PKGSOURCES) $(INTERNALSOURCES)
	cd `dirname $@` && go build $(BUILDARGS) -o `basename $@`

cmd/external-demo1/external-demo1: $(shell find cmd/external-demo1 -type f -name '*.go') $(PKGSOURCES) $(INTERNALSOURCES)
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
lint: govet gofmt golangcilint lint-python ## Lint the code (also fix the code if FIX=1, default)

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

.PHONY: venv
venv: tmp/python_venv/bin/activate 

tmp/python_venv/bin/activate: requirements.txt
	@mkdir -p tmp
	python3 -m venv tmp/python_venv
	source tmp/python_venv/bin/activate && pip install -r requirements.txt

.PHONY: freeze-requirements
freeze-requirements: tmp/python_venv/bin/activate ## Freeze the python (dev) requirements
	source tmp/python_venv/bin/activate && pip freeze > requirements.txt

.PHONY: lint-python
lint-python: tmp/python_venv/bin/activate ## Lint the python code
	@if test "$(FIX)" = "1"; then \
		source tmp/python_venv/bin/activate && set -x; ruff format .;\
	else \
		source tmp/python_venv/bin/activate && set -x; ruff format --diff .;\
	fi
	@if test "$(FIX)" = "1"; then \
		source tmp/python_venv/bin/activate && set -x; ruff check --fix .;\
	else \
		source tmp/python_venv/bin/activate && set -x; ruff check --no-fix .;\
	fi

doc: build tmp/python_venv/bin/activate ## Generate the documentation
	source tmp/python_venv/bin/activate && ./docs/termtosvg.py --command "./cmd/stacktrace-demo1/stacktrace-demo1" --lines 24 --columns 120 ./docs/stacktrace-demo1.svg
	source tmp/python_venv/bin/activate && ./docs/termtosvg.py --command "./cmd/human-demo1/human-demo1" --lines 10 --columns 120 ./docs/human-demo1.svg
	source tmp/python_venv/bin/activate && ./docs/termtosvg.py --command "./cmd/external-demo1/external-demo1" --lines 10 --columns 120 ./docs/external-demo1.svg
	source tmp/python_venv/bin/activate && ./docs/termtosvg.py --command "./cmd/slogc-demo1/slogc-demo1" --lines 28 --columns 120 ./docs/slogc-demo1.svg
	source tmp/python_venv/bin/activate && jinja-tree .

.PHONY: help
help::
	@# See https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
