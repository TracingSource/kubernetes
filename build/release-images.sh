#!/usr/bin/env bash

# Build Kubernetes release images. This will build the server target binaries,
# and create wrap them in Docker images, see `make release` for full releases

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/build/common.sh"
source "${KUBE_ROOT}/build/lib/release.sh"

CMD_TARGETS="${KUBE_SERVER_IMAGE_TARGETS[*]}"
if [[ "${KUBE_BUILD_CONFORMANCE}" =~ [yY] ]]; then
    CMD_TARGETS="${CMD_TARGETS} ${KUBE_CONFORMANCE_IMAGE_TARGETS[*]}"
fi

# TODO(dims): Remove this when we get rid of hyperkube image
CMD_TARGETS="${CMD_TARGETS} cmd/kubelet cmd/kubectl"

kube::build::verify_prereqs
kube::build::build_image
kube::build::run_build_command make all WHAT="${CMD_TARGETS}" KUBE_BUILD_PLATFORMS="${KUBE_SERVER_PLATFORMS[*]}"

kube::build::copy_output

kube::release::build_server_images
