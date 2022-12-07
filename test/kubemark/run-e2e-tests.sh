#!/usr/bin/env bash

export KUBERNETES_PROVIDER="kubemark"
export KUBE_CONFIG_FILE="config-default.sh"

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

# We need an absolute path to KUBE_ROOT
ABSOLUTE_ROOT=$(readlink -f "${KUBE_ROOT}")

source "${KUBE_ROOT}/cluster/kubemark/util.sh"

echo "Kubemark master name: ${MASTER_NAME}"

detect-master

export KUBE_MASTER_URL="https://${KUBE_MASTER_IP}"
export KUBECONFIG="${ABSOLUTE_ROOT}/test/kubemark/resources/kubeconfig.kubemark"
export E2E_MIN_STARTUP_PODS=0

if [[ -z "$*" ]]; then
	ARGS=('--ginkgo.focus=[Feature:Performance]')
else
	ARGS=("$@")
fi

if [[ "${ENABLE_KUBEMARK_CLUSTER_AUTOSCALER}" == "true" ]]; then
  ARGS+=("--kubemark-external-kubeconfig=${DEFAULT_KUBECONFIG}")
fi

# Running locally.
for ((i=0; i < ${#ARGS[@]}; i++)); do
  ARGS[$i]="$(echo "${ARGS[$i]}" | sed -e 's/\[/\\\[/g' -e 's/\]/\\\]/g' )"
done
"${KUBE_ROOT}/hack/ginkgo-e2e.sh" "--e2e-verify-service-account=false" "--dump-logs-on-failure=false" "${ARGS[@]}"
