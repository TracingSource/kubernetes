#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
KUBE_REMOTE_RUNTIME_ROOT="${KUBE_ROOT}/staging/src/k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

function cleanup {
	rm -rf "${KUBE_REMOTE_RUNTIME_ROOT}/_tmp/"
}

trap cleanup EXIT

mkdir -p "${KUBE_REMOTE_RUNTIME_ROOT}/_tmp"
cp "${KUBE_REMOTE_RUNTIME_ROOT}/api.pb.go" "${KUBE_REMOTE_RUNTIME_ROOT}/_tmp/"

ret=0
KUBE_VERBOSE=3 "${KUBE_ROOT}/hack/update-generated-runtime.sh"
diff -I "gzipped FileDescriptorProto" -I "0x" -Naupr "${KUBE_REMOTE_RUNTIME_ROOT}/_tmp/api.pb.go" "${KUBE_REMOTE_RUNTIME_ROOT}/api.pb.go" || ret=$?
if [[ $ret -eq 0 ]]; then
    echo "Generated container runtime api is up to date."
    cp "${KUBE_REMOTE_RUNTIME_ROOT}/_tmp/api.pb.go" "${KUBE_REMOTE_RUNTIME_ROOT}/"
else
    echo "Generated container runtime api is out of date. Please run hack/update-generated-runtime.sh"
    exit 1
fi
