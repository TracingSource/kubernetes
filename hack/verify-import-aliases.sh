#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::verify_go_version

cd "${KUBE_ROOT}"

ret=0
go run cmd/preferredimports/preferredimports.go "$@" || ret=$?
if [[ $ret -ne 0 ]]; then
  echo "!!! Please see hack/.import-aliases for the preferred aliases for imports." >&2
  exit 1
fi

