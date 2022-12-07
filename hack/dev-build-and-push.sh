#!/usr/bin/env bash

# This script will build a dev release and push it to an existing cluster.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# Build a dev release
if ! make -f "${KUBE_ROOT}"/Makefile quick-release
then
        echo "Building a release failed!"
        exit 1
fi

# Now push this out to the cluster
"${KUBE_ROOT}/cluster/kube-push.sh"
