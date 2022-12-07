#!/usr/bin/env bash

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${KUBE_ROOT}/test/kubemark/cloud-provider-config.sh"
source "${KUBE_ROOT}/cluster/${CLOUD_PROVIDER}/util.sh"
source "${KUBE_ROOT}/cluster/kubemark/${CLOUD_PROVIDER}/config-default.sh"
