#!/usr/bin/env bash

# Tear down a Kubernetes cluster.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

if [ -f "${KUBE_ROOT}/cluster/env.sh" ]; then
    source "${KUBE_ROOT}/cluster/env.sh"
fi

source "${KUBE_ROOT}/cluster/kube-util.sh"

echo "Bringing down cluster using provider: $KUBERNETES_PROVIDER"

echo "... calling verify-prereqs" >&2
verify-prereqs
echo "... calling verify-kube-binaries" >&2
verify-kube-binaries
echo "... calling kube-down" >&2
kube-down

echo "Done"
