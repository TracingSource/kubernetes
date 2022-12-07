#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

make -C "${KUBE_ROOT}" WHAT=cmd/genswaggertypedocs

# Find binary
genswaggertypedocs=$(kube::util::find-binary "genswaggertypedocs")

gen_swagger_result=0
result=0

find_files() {
  find . -not \( \
      \( \
        -wholename './output' \
        -o -wholename './_output' \
        -o -wholename './_gopath' \
        -o -wholename './release' \
        -o -wholename './target' \
        -o -wholename '*/third_party/*' \
        -o -wholename '*/vendor/*' \
      \) -prune \
    \) \
    \( -wholename '*pkg/apis/*/v*/types.go' \
       -o -wholename '*pkg/api/unversioned/types.go' \
    \)
}

if [[ $# -eq 0 ]]; then
  versioned_api_files=$(find_files | grep -E "pkg/.[^/]*/((v.[^/]*)|unversioned)/types\.go") || true
else
  versioned_api_files="${*}"
fi

for file in $versioned_api_files; do
  $genswaggertypedocs -v -s "${file}" -f - || gen_swagger_result=$?
  if [[ "${gen_swagger_result}" -ne "0" ]]; then
    echo "API file: ${file} is missing: ${gen_swagger_result} descriptions"
    result=1
  fi
  if grep json: "${file}" | grep -v // | grep description: ; then
    echo "API file: ${file} should not contain descriptions in struct tags"
    result=1
  fi
  if grep json: "${file}" | grep -Ee ",[[:space:]]+omitempty|omitempty[[:space:]]+" ; then
    echo "API file: ${file} should not contain leading or trailing spaces for omitempty directive"
    result=1
  fi
done

internal_types_files="${KUBE_ROOT}/pkg/apis/core/types.go ${KUBE_ROOT}/pkg/apis/extensions/types.go"
for internal_types_file in $internal_types_files; do
  if [[ ! -e $internal_types_file ]]; then
    echo "Internal types file ${internal_types_file} does not exist"
    result=1
    continue
  fi

  if grep json: "${internal_types_file}" | grep -v // | grep description: ; then
    echo "Internal API types should not contain descriptions"
    result=1
  fi
done

exit ${result}
