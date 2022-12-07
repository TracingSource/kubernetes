#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

: "${KUBECTL:=${KUBE_ROOT}/cluster/kubectl.sh}"
: "${KUBE_CONFIG_FILE:="config-test.sh"}"

export KUBECTL KUBE_CONFIG_FILE

source "${KUBE_ROOT}/cluster/kube-util.sh"

prepare-e2e

test-teardown
