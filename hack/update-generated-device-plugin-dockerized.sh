#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../" && pwd -P)"
DEVICE_PLUGIN_ALPHA="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/deviceplugin/v1alpha/"
DEVICE_PLUGIN_V1BETA1="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1/"

source "${KUBE_ROOT}/hack/lib/protoc.sh"
kube::protoc::generate_proto "${DEVICE_PLUGIN_ALPHA}"
kube::protoc::generate_proto "${DEVICE_PLUGIN_V1BETA1}"
