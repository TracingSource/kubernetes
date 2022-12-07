#!/usr/bin/env bash

# This script builds some binaries and then the conformance image.
# REGISTRY and VERSION must be set.
# Example usage:
#   $ export REGISTRY=gcr.io/someone
#   $ export VERSION=v1.4.0-testfix
#   ./hack/dev-push-conformance.sh
# That will build and push gcr.io/someone/conformance-amd64:v1.4.0-testfix

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."
source "${KUBE_ROOT}/build/common.sh"

if [[ -z "${REGISTRY:-}" ]]; then
	echo "REGISTRY must be set"
	exit 1
fi
if [[ -z "${VERSION:-}" ]]; then
	echo "VERSION must be set"
	exit 1
fi

IMAGE="${REGISTRY}/conformance-amd64:${VERSION}"

kube::build::verify_prereqs
kube::build::build_image
kube::build::run_build_command make WHAT="vendor/github.com/onsi/ginkgo/ginkgo test/e2e/e2e.test cmd/kubectl cluster/images/conformance/go-runner"
kube::build::copy_output

make -C "${KUBE_ROOT}/cluster/images/conformance" build
docker push "${IMAGE}"
