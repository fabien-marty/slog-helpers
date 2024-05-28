#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

export LOG_FORMAT=json-gcp
export LOG_DESTINATION=stdout
"${SCRIPT_DIR}/../slogc-demo1/slogc-demo1" | "${SCRIPT_DIR}/../../docs/format-json-lines.py"