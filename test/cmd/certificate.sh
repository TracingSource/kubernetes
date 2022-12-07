#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

run_certificates_tests() {
  set -o nounset
  set -o errexit

  kube::log::status "Testing certificates"

  # approve
  kubectl create -f hack/testdata/csr.yml "${kube_flags[@]:?}"
  kube::test::get_object_assert 'csr/foo' '{{range.status.conditions}}{{.type}}{{end}}' ''
  kubectl certificate approve foo "${kube_flags[@]}"
  kubectl get csr "${kube_flags[@]}" -o json
  kube::test::get_object_assert 'csr/foo' '{{range.status.conditions}}{{.type}}{{end}}' 'Approved'
  kubectl delete -f hack/testdata/csr.yml "${kube_flags[@]}"
  kube::test::get_object_assert csr "{{range.items}}{{${id_field:?}}}{{end}}" ''

  kubectl create -f hack/testdata/csr.yml "${kube_flags[@]}"
  kube::test::get_object_assert 'csr/foo' '{{range.status.conditions}}{{.type}}{{end}}' ''
  kubectl certificate approve -f hack/testdata/csr.yml "${kube_flags[@]}"
  kubectl get csr "${kube_flags[@]}" -o json
  kube::test::get_object_assert 'csr/foo' '{{range.status.conditions}}{{.type}}{{end}}' 'Approved'
  kubectl delete -f hack/testdata/csr.yml "${kube_flags[@]}"
  kube::test::get_object_assert csr "{{range.items}}{{$id_field}}{{end}}" ''

  # deny
  kubectl create -f hack/testdata/csr.yml "${kube_flags[@]}"
  kube::test::get_object_assert 'csr/foo' '{{range.status.conditions}}{{.type}}{{end}}' ''
  kubectl certificate deny foo "${kube_flags[@]}"
  kubectl get csr "${kube_flags[@]}" -o json
  kube::test::get_object_assert 'csr/foo' '{{range.status.conditions}}{{.type}}{{end}}' 'Denied'
  kubectl delete -f hack/testdata/csr.yml "${kube_flags[@]}"
  kube::test::get_object_assert csr "{{range.items}}{{$id_field}}{{end}}" ''

  kubectl create -f hack/testdata/csr.yml "${kube_flags[@]}"
  kube::test::get_object_assert 'csr/foo' '{{range.status.conditions}}{{.type}}{{end}}' ''
  kubectl certificate deny -f hack/testdata/csr.yml "${kube_flags[@]}"
  kubectl get csr "${kube_flags[@]}" -o json
  kube::test::get_object_assert 'csr/foo' '{{range.status.conditions}}{{.type}}{{end}}' 'Denied'
  kubectl delete -f hack/testdata/csr.yml "${kube_flags[@]}"
  kube::test::get_object_assert csr "{{range.items}}{{$id_field}}{{end}}" ''

  set +o nounset
  set +o errexit
}
