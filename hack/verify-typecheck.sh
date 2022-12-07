#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::verify_go_version

cd "${KUBE_ROOT}"

make --no-print-directory -C "${KUBE_ROOT}" generated_files

ret=0
go run test/typecheck/main.go "$@" || ret=$?
if [[ $ret -ne 0 ]]; then
  echo "!!! Type Check has failed. This may cause cross platform build failures." >&2
  echo "!!! Please see https://git.k8s.io/kubernetes/test/typecheck for more information." >&2
  exit 1
fi
