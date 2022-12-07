#!/usr/bin/env bash

# Run a command in the docker build container.  Typically this will be one of
# the commands in `hack/`.  When running in the build container the user is sure
# to have a consistent reproducible build environment.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "$KUBE_ROOT/build/common.sh"

KUBE_RUN_COPY_OUTPUT="${KUBE_RUN_COPY_OUTPUT:-y}"

kube::build::verify_prereqs
kube::build::build_image

if [[ ${KUBE_RUN_COPY_OUTPUT} =~ ^[yY]$ ]]; then
  kube::log::status "Output from this container will be rsynced out upon completion. Set KUBE_RUN_COPY_OUTPUT=n to disable."
else
  kube::log::status "Output from this container will NOT be rsynced out upon completion. Set KUBE_RUN_COPY_OUTPUT=y to enable."
fi

kube::build::run_build_command "$@"

if [[ ${KUBE_RUN_COPY_OUTPUT} =~ ^[yY]$ ]]; then
  kube::build::copy_output
fi
