SHELL:=/bin/bash
default: help

screenshots: ../cmd/demo/demo ../tmp/python_venv/bin/activate ## Generate screenshots
	source ../tmp/python_venv/bin/activate && ./termtosvg.py --command="$<" --lines=40 --columns=120 demo.svg

.PHONY: help
help::
	@# See https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
