#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../" && pwd -P)"
POD_RESOURCES_ALPHA="${KUBE_ROOT}/pkg/kubelet/apis/podresources/v1alpha1/"

source "${KUBE_ROOT}/hack/lib/protoc.sh"
kube::protoc::generate_proto "${POD_RESOURCES_ALPHA}"
