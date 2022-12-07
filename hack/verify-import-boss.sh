#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

make -C "${KUBE_ROOT}" WHAT=vendor/k8s.io/code-generator/cmd/import-boss

$(kube::util::find-binary "import-boss") --verify-only "$@"
