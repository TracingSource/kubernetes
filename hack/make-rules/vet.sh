#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${KUBE_ROOT}/hack/lib/init.sh"

cd "${KUBE_ROOT}"

# If called directly, exit.
if [[ "${CALLED_FROM_MAIN_MAKEFILE:-""}" == "" ]]; then
    echo "ERROR: $0 should not be run directly." >&2
    echo >&2
    echo "Please run this command using \"make vet\""
    exit 1
fi

# This is required before we run govet for the results to be correct.
# See https://github.com/golang/go/issues/16086 for details.
go install ./cmd/...

# Filter out arguments that start with "-" and move them to goflags.
targets=()
for arg; do
  if [[ "${arg}" == -* ]]; then
    goflags+=("${arg}")
  else
    targets+=("${arg}")
  fi
done

if [[ ${#targets[@]} -eq 0 ]]; then
  # Do not run on third_party directories or generated client code or build tools.
  while IFS='' read -r line; do
    targets+=("${line}")
  done < <(go list -e ./... | grep -E -v "/(build|third_party|vendor|staging|clientset_generated)/")
fi

go vet "${goflags[@]:+${goflags[@]}}" "${targets[@]}"
