#!/usr/bin/env bash

# Generates `types_swagger_doc_generated.go` files for API group
# versions. That file contains functions on API structs that return
# the comments that should be surfaced for the corresponding API type
# in our API docs.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"
source "${KUBE_ROOT}/hack/lib/swagger.sh"

kube::golang::setup_env

IFS=" " read -r -a GROUP_VERSIONS <<< "meta/v1 meta/v1beta1 ${KUBE_AVAILABLE_GROUP_VERSIONS}"

# To avoid compile errors, remove the currently existing files.
for group_version in "${GROUP_VERSIONS[@]}"; do
  rm -f "$(kube::util::group-version-to-pkg-path "${group_version}")/types_swagger_doc_generated.go"
done
for group_version in "${GROUP_VERSIONS[@]}"; do
  kube::swagger::gen_types_swagger_doc "${group_version}" "$(kube::util::group-version-to-pkg-path "${group_version}")"
done
