#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::verify_go_version

cd "${KUBE_ROOT}"

ret=0
go run cmd/preferredimports/preferredimports.go --confirm "$@" || ret=$?
if [[ $ret -ne 0 ]]; then
  echo "!!! Unable to fix imports programmatically. Please see errors above." >&2
  exit 1
fi

