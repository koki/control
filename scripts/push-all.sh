#!/usr/bin/env bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
${SCRIPT_DIR}/./push-sidecar.sh
${SCRIPT_DIR}/./push-controller.sh
