#!/usr/bin/env bash

# Clean out the output directory on the docker host.
set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/build/common.sh"

kube::build::verify_prereqs false
kube::build::clean
