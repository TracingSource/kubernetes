#!/usr/bin/env bash

# Run a bash script in the Docker build image.
#
# This container will have a snapshot of the current sources.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/build/common.sh"

KUBE_RUN_COPY_OUTPUT="${KUBE_RUN_COPY_OUTPUT:-n}" "${KUBE_ROOT}/build/run.sh" bash "$@"
