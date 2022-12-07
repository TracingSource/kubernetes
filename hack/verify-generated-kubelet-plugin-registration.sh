#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
ERROR="Kubelet Plugin Registration api is out of date. Please run hack/update-generated-kubelet-plugin-registration.sh"
KUBELET_PLUGIN_REGISTRATION_V1ALPHA="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/pluginregistration/v1alpha1/"
KUBELET_PLUGIN_REGISTRATION_V1BETA="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/pluginregistration/v1beta1/"

source "${KUBE_ROOT}/hack/lib/protoc.sh"
kube::golang::setup_env

function cleanup {
	rm -rf "${KUBELET_PLUGIN_REGISTRATION_V1ALPHA}/_tmp/"
	rm -rf "${KUBELET_PLUGIN_REGISTRATION_V1BETA}/_tmp/"
}

trap cleanup EXIT

mkdir -p "${KUBELET_PLUGIN_REGISTRATION_V1ALPHA}/_tmp"
mkdir -p "${KUBELET_PLUGIN_REGISTRATION_V1BETA}/_tmp"

cp "${KUBELET_PLUGIN_REGISTRATION_V1ALPHA}/api.pb.go" "${KUBELET_PLUGIN_REGISTRATION_V1ALPHA}/_tmp/"
cp "${KUBELET_PLUGIN_REGISTRATION_V1BETA}/api.pb.go" "${KUBELET_PLUGIN_REGISTRATION_V1BETA}/_tmp/"

# Check V1Alpha
KUBE_VERBOSE=3 "${KUBE_ROOT}/hack/update-generated-kubelet-plugin-registration.sh"
kube::protoc::diff "${KUBELET_PLUGIN_REGISTRATION_V1ALPHA}/api.pb.go" "${KUBELET_PLUGIN_REGISTRATION_V1ALPHA}/_tmp/api.pb.go" "${ERROR}"
echo "Generated Kubelet Plugin Registration api is up to date."

# Check V1Beta
KUBE_VERBOSE=3 "${KUBE_ROOT}/hack/update-generated-kubelet-plugin-registration.sh"
kube::protoc::diff "${KUBELET_PLUGIN_REGISTRATION_V1BETA}/api.pb.go" "${KUBELET_PLUGIN_REGISTRATION_V1BETA}/_tmp/api.pb.go" "${ERROR}"
echo "Generated Kubelet Plugin Registration api is up to date."
