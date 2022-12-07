#!/usr/bin/env bash

# This script will build a dev release and bring up a new cluster with that
# release.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# Build a dev release


if ! make -f "${KUBE_ROOT}"/Makefile quick-release
then
        echo "Building the release failed!"
        exit 1
fi

# Now bring a new cluster up with that release.
"${KUBE_ROOT}/cluster/kube-up.sh"
