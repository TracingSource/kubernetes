#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../" && pwd -P)"
KUBELET_PLUGIN_REGISTRATION_V1ALPHA="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/pluginregistration/v1alpha1/"
KUBELET_PLUGIN_REGISTRATION_V1BETA="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/pluginregistration/v1beta1/"
KUBELET_EXAMPLE_PLUGIN_V1BETA1="${KUBE_ROOT}/pkg/kubelet/pluginmanager/pluginwatcher/example_plugin_apis/v1beta1/"
KUBELET_EXAMPLE_PLUGIN_V1BETA2="${KUBE_ROOT}/pkg/kubelet/pluginmanager/pluginwatcher/example_plugin_apis/v1beta2/"

source "${KUBE_ROOT}/hack/lib/protoc.sh"
kube::protoc::generate_proto "${KUBELET_PLUGIN_REGISTRATION_V1ALPHA}"
kube::protoc::generate_proto "${KUBELET_PLUGIN_REGISTRATION_V1BETA}"
kube::protoc::generate_proto "${KUBELET_EXAMPLE_PLUGIN_V1BETA1}"
kube::protoc::generate_proto "${KUBELET_EXAMPLE_PLUGIN_V1BETA2}"
