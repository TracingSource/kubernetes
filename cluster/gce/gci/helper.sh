#!/usr/bin/env bash

# A library of helper functions and constant for GCI distro

# Creates the GCI specific metadata files if they do not exit.
# Assumed var
#   KUBE_TEMP
function ensure-gci-metadata-files {
  if [[ ! -f "${KUBE_TEMP}/gci-update.txt" ]]; then
    echo -n "update_disabled" > "${KUBE_TEMP}/gci-update.txt"
  fi
  if [[ ! -f "${KUBE_TEMP}/gci-ensure-gke-docker.txt" ]]; then
    echo -n "true" > "${KUBE_TEMP}/gci-ensure-gke-docker.txt"
  fi
  if [[ ! -f "${KUBE_TEMP}/gci-docker-version.txt" ]]; then
    echo -n "${GCI_DOCKER_VERSION:-}" > "${KUBE_TEMP}/gci-docker-version.txt"
  fi
}
