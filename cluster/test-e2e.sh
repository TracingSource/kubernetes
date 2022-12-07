#!/usr/bin/env bash

# Tests a running Kubernetes cluster.
# TODO: move code from hack/ginkgo-e2e.sh to here

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/cluster/kube-util.sh"

echo "Testing cluster with provider: ${KUBERNETES_PROVIDER}" 1>&2

echo "Running e2e tests:" 1>&2
echo "./hack/ginkgo-e2e.sh $*" 1>&2
exec "${KUBE_ROOT}/hack/ginkgo-e2e.sh" "$@"
