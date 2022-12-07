#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

make -C "${KUBE_ROOT}" WHAT=cmd/importverifier

# Find binary
importverifier=$(kube::util::find-binary "importverifier")

if [[ ! -x "$importverifier" ]]; then
  {
    echo "It looks as if you don't have a compiled importverifier binary"
    echo
    echo "If you are running from a clone of the git repo, please run"
    echo "'make WHAT=cmd/importverifier'."
  } >&2
  exit 1
fi

"${importverifier}" "k8s.io/" "${KUBE_ROOT}/staging/publishing/import-restrictions.yaml"
