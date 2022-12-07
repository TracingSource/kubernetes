#!/usr/bin/env bash

# A library of helper functions and constant for GCI distro
source "${KUBE_ROOT}/cluster/gce/gci/helper.sh"

function get-node-instance-metadata-from-file {
  local metadata=""
  metadata+="kube-env=${KUBE_TEMP}/node-kube-env.yaml,"
  metadata+="kubelet-config=${KUBE_TEMP}/node-kubelet-config.yaml,"
  metadata+="user-data=${KUBE_ROOT}/cluster/gce/gci/node.yaml,"
  metadata+="configure-sh=${KUBE_ROOT}/cluster/gce/gci/configure.sh,"
  metadata+="cluster-location=${KUBE_TEMP}/cluster-location.txt,"
  metadata+="cluster-name=${KUBE_TEMP}/cluster-name.txt,"
  metadata+="gci-update-strategy=${KUBE_TEMP}/gci-update.txt,"
  metadata+="gci-ensure-gke-docker=${KUBE_TEMP}/gci-ensure-gke-docker.txt,"
  metadata+="gci-docker-version=${KUBE_TEMP}/gci-docker-version.txt,"
  metadata+="shutdown-script=${KUBE_ROOT}/cluster/gce/gci/shutdown.sh,"
  metadata+="${NODE_EXTRA_METADATA}"
  echo "${metadata}"
}

# Assumed vars:
#   scope_flags
# Parameters:
#   $1: template name (required).
function create-linux-node-instance-template {
  local template_name="$1"
  local machine_type="${2:-$NODE_SIZE}"

  ensure-gci-metadata-files
  # shellcheck disable=2154 # 'scope_flags' is assigned by upstream
  create-node-template "${template_name}" "${scope_flags[*]}" "$(get-node-instance-metadata-from-file)" "" "linux" "${machine_type}"
}
