#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=${KUBE_ROOT:-$(dirname "${BASH_SOURCE[0]}")/..}
source "${KUBE_ROOT}/cluster/clientbin.sh"

# If KUBEADM_PATH isn't set, gather up the list of likely places and use ls
# to find the latest one.
if [[ -z "${KUBEADM_PATH:-}" ]]; then
  kubeadm=$( get_bin "kubeadm" "cmd/kubeadm" )

  if [[ ! -x "$kubeadm" ]]; then
    print_error "kubeadm"
    exit 1
  fi
elif [[ ! -x "${KUBEADM_PATH}" ]]; then
  {
    echo "KUBEADM_PATH environment variable set to '${KUBEADM_PATH}', but "
    echo "this doesn't seem to be a valid executable."
  } >&2
  exit 1
fi
kubeadm="${KUBEADM_PATH:-${kubeadm}}"

"${kubeadm}" "${@+$@}"
