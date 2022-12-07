#!/usr/bin/env bash

# Build the docker image necessary for building Kubernetes
#
# This script will package the parts of the repo that we need to build
# Kubernetes into a tar file and put it in the right place in the output
# directory.  It will then copy over the Dockerfile and build the kube-build
# image.
set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."
source "${KUBE_ROOT}/build/common.sh"

kube::build::verify_prereqs
kube::build::build_image
