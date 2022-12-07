#!/usr/bin/env bash
set -x

NODE_INSTANCE_PREFIX=${NODE_INSTANCE_PREFIX:-"${INSTANCE_PREFIX}-node"}

source "${KUBE_ROOT}/cluster/gce/util.sh"
