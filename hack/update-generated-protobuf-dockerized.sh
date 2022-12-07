#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

BINS=(
	vendor/k8s.io/code-generator/cmd/go-to-protobuf
	vendor/k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo
)
make -C "${KUBE_ROOT}" WHAT="${BINS[*]}"

if [[ -z "$(which protoc)" || "$(protoc --version)" != "libprotoc 3."* ]]; then
  echo "Generating protobuf requires protoc 3.0.0-beta1 or newer. Please download and"
  echo "install the platform appropriate Protobuf package for your OS: "
  echo
  echo "  https://github.com/google/protobuf/releases"
  echo
  echo "WARNING: Protobuf changes are not being validated"
  exit 1
fi

gotoprotobuf=$(kube::util::find-binary "go-to-protobuf")

while IFS=$'\n' read -r line; do
  APIROOTS+=( "$line" );
done <<< "${1}"
shift

# requires the 'proto' tag to build (will remove when ready)
# searches for the protoc-gen-gogo extension in the output directory
# satisfies import of github.com/gogo/protobuf/gogoproto/gogo.proto and the
# core Google protobuf types
PATH="${KUBE_ROOT}/_output/bin:${PATH}" \
  "${gotoprotobuf}" \
  --proto-import="${KUBE_ROOT}/vendor" \
  --proto-import="${KUBE_ROOT}/third_party/protobuf" \
  --packages="$(IFS=, ; echo "${APIROOTS[*]}")" \
  --go-header-file "${KUBE_ROOT}/hack/boilerplate/boilerplate.generatego.txt" \
  "$@"
