#!/usr/bin/env bash

# This script sets up a go workspace locally and builds all for all appropriate
# platforms.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${KUBE_ROOT}/hack/lib/init.sh"

# NOTE: Using "${array[*]}" here is correct.  [@] becomes distinct words (in
# bash parlance).

make all WHAT="${KUBE_SERVER_TARGETS[*]}" KUBE_BUILD_PLATFORMS="${KUBE_SERVER_PLATFORMS[*]}"

make all WHAT="${KUBE_NODE_TARGETS[*]}" KUBE_BUILD_PLATFORMS="${KUBE_NODE_PLATFORMS[*]}"

make all WHAT="${KUBE_CLIENT_TARGETS[*]}" KUBE_BUILD_PLATFORMS="${KUBE_CLIENT_PLATFORMS[*]}"

make all WHAT="${KUBE_TEST_TARGETS[*]}" KUBE_BUILD_PLATFORMS="${KUBE_TEST_PLATFORMS[*]}"

make all WHAT="${KUBE_TEST_SERVER_TARGETS[*]}" KUBE_BUILD_PLATFORMS="${KUBE_TEST_SERVER_PLATFORMS[*]}"
