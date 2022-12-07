#!/usr/bin/env bash

# This script sets up a go workspace locally and builds all go components.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
KUBE_VERBOSE="${KUBE_VERBOSE:-1}"
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::build_binaries "$@"
kube::golang::place_bins
