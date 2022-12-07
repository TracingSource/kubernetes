#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Runs tests related to kubectl delete --all-namespaces.
run_kubectl_delete_allnamespaces_tests() {
  set -o nounset
  set -o errexit

  ns_one="namespace-$(date +%s)-${RANDOM}"
  ns_two="namespace-$(date +%s)-${RANDOM}"
  kubectl create namespace "${ns_one}"
  kubectl create namespace "${ns_two}"

  kubectl create configmap "one" --namespace="${ns_one}"
  kubectl create configmap "two" --namespace="${ns_two}"
  kubectl label configmap "one" --namespace="${ns_one}" deletetest=true
  kubectl label configmap "two" --namespace="${ns_two}" deletetest=true
  kubectl delete configmap -l deletetest=true --all-namespaces

  # no configmaps should be in either of those namespaces
  kubectl config set-context "${CONTEXT}" --namespace="${ns_one}"
  kube::test::get_object_assert configmap "{{range.items}}{{${id_field:?}}}:{{end}}" ''
  kubectl config set-context "${CONTEXT}" --namespace="${ns_two}"
  kube::test::get_object_assert configmap "{{range.items}}{{${id_field:?}}}:{{end}}" ''

  set +o nounset
  set +o errexit
}
