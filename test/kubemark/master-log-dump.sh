#!/usr/bin/env bash

REPORT_DIR="${1:-_artifacts}"
KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

source "${KUBE_ROOT}/test/kubemark/cloud-provider-config.sh"
source "${KUBE_ROOT}/cluster/kubemark/${CLOUD_PROVIDER}/config-default.sh"

export KUBEMARK_MASTER_NAME="${MASTER_NAME}"

echo "Dumping logs for kubemark master: ${KUBEMARK_MASTER_NAME}"
DUMP_ONLY_MASTER_LOGS=true "${KUBE_ROOT}/cluster/log-dump/log-dump.sh" "${REPORT_DIR}"
