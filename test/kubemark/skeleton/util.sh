#!/usr/bin/env bash

# This script contains the helper functions that each provider hosting
# Kubermark must implement to use test/kubemark/start-kubemark.sh and
# test/kubemark/stop-kubemark.sh scripts.

# This function should authenticate docker to be able to read/write to
# the right container registry (needed for pushing kubemark image).
function authenticate-docker {
	echo "Configuring registry authentication" 1>&2
}

# This function should create kubemark master and write kubeconfig to
# "${RESOURCE_DIRECTORY}/kubeconfig.kubemark".
function create-kubemark-master {
  echo "Creating cluster..."
}

# This function should delete kubemark master.
function delete-kubemark-master {
  echo "Deleting cluster..."
}

# Common colors used throughout the kubemark scripts
if [[ -z "${color_start-}" ]]; then
  declare -r color_start="\033["
  # shellcheck disable=SC2034
  declare -r color_red="${color_start}0;31m"
  # shellcheck disable=SC2034
  declare -r color_yellow="${color_start}0;33m"
  # shellcheck disable=SC2034
  declare -r color_green="${color_start}0;32m"
  # shellcheck disable=SC2034
  declare -r color_blue="${color_start}1;34m"
  # shellcheck disable=SC2034
  declare -r color_cyan="${color_start}1;36m"
  # shellcheck disable=SC2034
  declare -r color_norm="${color_start}0m"
fi
