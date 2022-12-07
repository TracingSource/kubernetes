#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

run_kubeadm_tests() {
  set -o nounset
  set -o errexit

  KUBEADM_PATH="${KUBEADM_PATH:=$(kube::realpath "${KUBE_ROOT}")/cluster/kubeadm.sh}"

  # If testing a different version of kubeadm than the current build, you can
  # comment this out to save yourself from needlessly building here.
  make -C "${KUBE_ROOT}" WHAT=cmd/kubeadm

  #TODO(runyontr): Remove the KUBE_TIMEOUT override when 
  # kubernetes/kubeadm/issues/1430 is fixed
  make -C "${KUBE_ROOT}" test \
  WHAT=k8s.io/kubernetes/cmd/kubeadm/test/cmd \
  KUBE_TEST_ARGS="--kubeadm-path '${KUBEADM_PATH}'" \
  KUBE_TIMEOUT="--timeout=600s"
  set +o nounset
  set +o errexit
}
