#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
KUBE_KMS_GRPC_ROOT="${KUBE_ROOT}/staging/src/k8s.io/apiserver/pkg/storage/value/encrypt/envelope/v1beta1/"
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

function cleanup {
	rm -rf "${KUBE_KMS_GRPC_ROOT}/_tmp/"
}

trap cleanup EXIT

mkdir -p "${KUBE_KMS_GRPC_ROOT}/_tmp"
cp "${KUBE_KMS_GRPC_ROOT}/service.pb.go" "${KUBE_KMS_GRPC_ROOT}/_tmp/"

ret=0
KUBE_VERBOSE=3 "${KUBE_ROOT}/hack/update-generated-kms.sh"
diff -I "gzipped FileDescriptorProto" -I "0x" -Naupr "${KUBE_KMS_GRPC_ROOT}/_tmp/service.pb.go" "${KUBE_KMS_GRPC_ROOT}/service.pb.go" || ret=$?
if [[ $ret -eq 0 ]]; then
    echo "Generated KMS gRPC is up to date."
    cp "${KUBE_KMS_GRPC_ROOT}/_tmp/service.pb.go" "${KUBE_KMS_GRPC_ROOT}/"
else
    echo "Generated KMS gRPC is out of date. Please run hack/update-generated-kms.sh"
    exit 1
fi
