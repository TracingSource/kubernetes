#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::verify_go_version

cd "${KUBE_ROOT}"

go run cmd/verifydependencies/verifydependencies.go "${KUBE_ROOT}"/build/dependencies.yaml
