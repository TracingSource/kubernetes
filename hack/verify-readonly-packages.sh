#!/usr/bin/env bash

# verify-readonly-packages.sh checks whether between $KUBE_VERIFY_GIT_BRANCH and HEAD files
# in readonly directories were modified. A directory is readonly iff it contains a .readonly
# file. Being readonly DOES NOT apply recursively to subdirectories.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

readonly branch=${1:-${KUBE_VERIFY_GIT_BRANCH:-master}}

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
        -o -wholename './staging/src/k8s.io/client-go/*vendor/*' \
        -o -wholename './staging/src/k8s.io/client-go/pkg/*' \
      \) -prune \
    \) -name '.readonly'
}

conflicts=()
while IFS=$'\n' read -r dir; do
    dir=${dir#./}
    if kube::util::has_changes "${branch}" "^${dir}/[^/]*\$" '/\.readonly$|/BUILD$|/zz_generated|/\.generated\.|\.proto$|\.pb\.go$' >/dev/null; then
        conflicts+=("${dir}")
    fi
done < <(find_files | sed 's|/.readonly||')

if [ ${#conflicts[@]} -gt 0 ]; then
    exec 1>&2
    for dir in "${conflicts[@]}"; do
        echo "Found ${dir}/.readonly, but files changed compared to \"${branch}\" branch."
    done
    exit 1
else
    echo "Readonly packages verified."
fi
# ex: ts=2 sw=2 et filetype=sh
