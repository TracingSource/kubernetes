#!/usr/bin/env bash

# Script that destroys Kubemark cluster and deletes all master resources.

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

source "${KUBE_ROOT}/test/kubemark/skeleton/util.sh"
source "${KUBE_ROOT}/test/kubemark/cloud-provider-config.sh"
source "${KUBE_ROOT}/test/kubemark/${CLOUD_PROVIDER}/util.sh"
source "${KUBE_ROOT}/cluster/kubemark/${CLOUD_PROVIDER}/config-default.sh"

if [[ -f "${KUBE_ROOT}/test/kubemark/${CLOUD_PROVIDER}/shutdown.sh" ]] ; then
  source "${KUBE_ROOT}/test/kubemark/${CLOUD_PROVIDER}/shutdown.sh"
fi

source "${KUBE_ROOT}/cluster/kubemark/util.sh"

KUBECTL="${KUBE_ROOT}/cluster/kubectl.sh"
KUBEMARK_DIRECTORY="${KUBE_ROOT}/test/kubemark"
RESOURCE_DIRECTORY="${KUBEMARK_DIRECTORY}/resources"

detect-project &> /dev/null

"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/addons" &> /dev/null || true
"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/hollow-node.yaml" &> /dev/null || true
"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/kubemark-ns.json" &> /dev/null || true

rm -rf "${RESOURCE_DIRECTORY}/addons" \
	"${RESOURCE_DIRECTORY}/kubeconfig.kubemark" \
	"${RESOURCE_DIRECTORY}/hollow-node.yaml"  &> /dev/null || true

delete-kubemark-master
