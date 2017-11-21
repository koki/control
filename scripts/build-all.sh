#!/usr/bin/env bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
${SCRIPT_DIR}/./build-sidecar.sh
${SCRIPT_DIR}/./build-controller.sh

