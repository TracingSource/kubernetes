#!/usr/bin/env bash

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

if [[ -n "${1:-}" ]]; then
  export KUBE_GCE_ZONE="${1}"
fi
if [[ -n "${2:-}" ]]; then
  export MULTIZONE="${2}"
fi
if [[ -n "${3:-}" ]]; then
  export KUBE_DELETE_NODES="${3}"
fi
if [[ -n "${4:-}" ]]; then
  export KUBE_USE_EXISTING_MASTER="${4}"
fi

source "${KUBE_ROOT}/hack/e2e-internal/e2e-down.sh"

