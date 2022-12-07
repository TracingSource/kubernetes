#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

# call verify on sub-project for now
#
# Note: these must be before the main script call because the later calls the sub-project's
#       update-codegen.sh scripts. We wouldn't see any error on changes then.
CODEGEN_PKG=./vendor/k8s.io/code-generator vendor/k8s.io/code-generator/hack/verify-codegen.sh
CODEGEN_PKG=./vendor/k8s.io/code-generator vendor/k8s.io/kube-aggregator/hack/verify-codegen.sh
CODEGEN_PKG=./vendor/k8s.io/code-generator vendor/k8s.io/sample-apiserver/hack/verify-codegen.sh
CODEGEN_PKG=./vendor/k8s.io/code-generator vendor/k8s.io/sample-controller/hack/verify-codegen.sh
CODEGEN_PKG=./vendor/k8s.io/code-generator vendor/k8s.io/apiextensions-apiserver/hack/verify-codegen.sh
CODEGEN_PKG=./vendor/k8s.io/code-generator vendor/k8s.io/metrics/hack/verify-codegen.sh
CODEGEN_PKG=./vendor/k8s.io/code-generator vendor/k8s.io/node-api/hack/verify-codegen.sh

"${KUBE_ROOT}/hack/update-codegen.sh" --verify-only
