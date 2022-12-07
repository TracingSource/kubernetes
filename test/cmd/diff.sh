#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Runs tests for kubectl diff
run_kubectl_diff_tests() {
    set -o nounset
    set -o errexit

    create_and_use_new_namespace
    kube::log::status "Testing kubectl diff"

    # Test that it works when the live object doesn't exist
    output_message=$(! kubectl diff -f hack/testdata/pod.yaml)
    kube::test::if_has_string "${output_message}" 'test-pod'

    kubectl apply -f hack/testdata/pod.yaml

    output_message=$(! kubectl diff -f hack/testdata/pod-changed.yaml)
    kube::test::if_has_string "${output_message}" 'k8s.gcr.io/pause:3.0'

    kubectl delete -f hack/testdata/pod.yaml

    set +o nounset
    set +o errexit
}

run_kubectl_diff_same_names() {
    set -o nounset
    set -o errexit

    create_and_use_new_namespace
    kube::log::status "Test kubectl diff with multiple resources with the same name"

    output_message=$(KUBECTL_EXTERNAL_DIFF="find" kubectl diff -Rf hack/testdata/diff/)
    kube::test::if_has_string "${output_message}" 'v1\.Pod\..*\.test'
    kube::test::if_has_string "${output_message}" 'apps\.v1\.Deployment\..*\.test'
    kube::test::if_has_string "${output_message}" 'v1\.ConfigMap\..*\.test'
    kube::test::if_has_string "${output_message}" 'v1\.Secret\..*\.test'

    set +o nounset
    set +o errexit
}
