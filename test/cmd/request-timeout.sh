#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

run_kubectl_request_timeout_tests() {
  set -o nounset
  set -o errexit

  kube::log::status "Testing kubectl request timeout"
  ### Test global request timeout option
  # Pre-condition: no POD exists
  create_and_use_new_namespace
  kube::test::get_object_assert pods "{{range.items}}{{${id_field:?}}}:{{end}}" ''
  # Command
  kubectl create "${kube_flags[@]:?}" -f test/fixtures/doc-yaml/admin/limitrange/valid-pod.yaml
  # Post-condition: valid-pod POD is created
  kubectl get "${kube_flags[@]}" pods -o json
  kube::test::get_object_assert pods "{{range.items}}{{$id_field}}:{{end}}" 'valid-pod:'

  ## check --request-timeout on 'get pod'
  output_message=$(kubectl get pod valid-pod --request-timeout=1)
  kube::test::if_has_string "${output_message}" 'valid-pod'

  ## check --request-timeout on 'get pod' with --watch
  output_message=$(kubectl get pod valid-pod --request-timeout=1 --watch 2>&1)
  kube::test::if_has_string "${output_message}" 'Timeout exceeded while reading body'

  ## check --request-timeout value with no time unit
  output_message=$(kubectl get pod valid-pod --request-timeout=1 2>&1)
  kube::test::if_has_string "${output_message}" 'valid-pod'

  ## check --request-timeout value with invalid time unit
  output_message=$(! kubectl get pod valid-pod --request-timeout="1p" 2>&1)
  kube::test::if_has_string "${output_message}" 'Invalid timeout value'

  # cleanup
  kubectl delete pods valid-pod "${kube_flags[@]}"

  set +o nounset
  set +o errexit
}
