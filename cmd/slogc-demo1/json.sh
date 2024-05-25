#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

export LOG_FORMAT=json
"${SCRIPT_DIR}/slogc-demo1"