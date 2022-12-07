#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
export KUBE_ROOT
source "${KUBE_ROOT}/hack/lib/init.sh"

# Ensure that we find the binaries we build before anything else.
export GOBIN="${KUBE_OUTPUT_BINPATH}"
PATH="${GOBIN}:${PATH}"

# Install tools we need, but only from vendor/...
go install k8s.io/kubernetes/vendor/github.com/client9/misspell/cmd/misspell

# Spell checking
# All the skipping files are defined in hack/.spelling_failures
skipping_file="${KUBE_ROOT}/hack/.spelling_failures"
failing_packages=$(sed "s| | -e |g" "${skipping_file}")
git ls-files | grep -v -e "${failing_packages}" | xargs misspell -i "Creater,creater,ect" -error -o stderr
