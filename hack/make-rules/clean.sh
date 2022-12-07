#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${KUBE_ROOT}/hack/lib/util.sh"

CLEAN_PATTERNS=(
  "_tmp"
  "doc_tmp"
  "((?!staging\/src\/k8s\.io\/apiextensions-apiserver\/pkg\/generated\/openapi).)*/zz_generated.openapi.go"
  "test/e2e/generated/bindata.go"
)

for pattern in "${CLEAN_PATTERNS[@]}"; do
  while IFS=$'\n' read -r match; do
    echo "Removing ${match#${KUBE_ROOT}\/} .."
    rm -rf "${match#${KUBE_ROOT}\/}"
  done <   <(find "${KUBE_ROOT}" -iregex "^${KUBE_ROOT}/${pattern}$")
done

# ex: ts=2 sw=2 et filetype=sh
