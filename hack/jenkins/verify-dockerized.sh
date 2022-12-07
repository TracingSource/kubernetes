#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

retry() {
  for i in {1..5}; do
    if "$@"; then
      return 0
    else
      sleep "${i}"
    fi
  done
  "$@"
}

# This script is intended to be run from kubekins-test container with a
# kubernetes repo mapped in. See k8s.io/test-infra/scenarios/kubernetes_verify.py

export PATH=${GOPATH}/bin:${PWD}/third_party/etcd:/usr/local/go/bin:${PATH}

# Set artifacts directory
export ARTIFACTS=${ARTIFACTS:-"${WORKSPACE}/artifacts"}
# Produce a JUnit-style XML test report
export KUBE_JUNIT_REPORT_DIR="${ARTIFACTS}"

export LOG_LEVEL=4

cd "${GOPATH}/src/k8s.io/kubernetes"

make verify
