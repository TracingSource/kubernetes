#!/usr/bin/env bash

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../../..

source "${KUBE_ROOT}/test/kubemark/common/util.sh"

# Leave the skeleton definition of execute-cmd-on-master-with-retries
# so only the pre-existing provider functions will target this.
function execute-cmd-on-pre-existing-master-with-retries() {
  IP_WITHOUT_PORT=$(echo "${MASTER_IP}" | cut -f 1 -d ':') || "${MASTER_IP}"

  RETRIES="${2:-1}" run-cmd-with-retries ssh kubernetes@"${IP_WITHOUT_PORT}" "${1}"
}
