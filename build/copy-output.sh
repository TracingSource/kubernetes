#!/usr/bin/env bash

# Copies any built binaries (and other generated files) out of the Docker build container.
set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/build/common.sh"

kube::build::verify_prereqs
kube::build::copy_output
