#!/usr/bin/env bash

# Script that destroys the clusters used, namespace, and deployment.

KUBECTL=kubectl
KUBEMARK_DIRECTORY="${KUBE_ROOT}/test/kubemark"
RESOURCE_DIRECTORY="${KUBEMARK_DIRECTORY}/resources"

# Login to cloud services
complete-login

# Remove resources created for kubemark
# shellcheck disable=SC2154 # Color defined in sourced script
echo -e "${color_yellow}REMOVING RESOURCES${color_norm}"
spawn-config
"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/addons" &> /dev/null || true
"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/hollow-node.yaml" &> /dev/null || true
"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/kubemark-ns.json" &> /dev/null || true
rm -rf "${RESOURCE_DIRECTORY}/addons"
	"${RESOURCE_DIRECTORY}/hollow-node.yaml" &> /dev/null || true

# Remove clusters, namespaces, and deployments
delete-clusters
if [[ -f "${RESOURCE_DIRECTORY}/iks-namespacelist.sh" ]] ; then
  bash "${RESOURCE_DIRECTORY}/iks-namespacelist.sh"
  rm -f "${RESOURCE_DIRECTORY}/iks-namespacelist.sh"
fi
# shellcheck disable=SC2154 # Color defined in sourced script
echo -e "${color_blue}EXECUTION COMPLETE${color_norm}"
exit 0
