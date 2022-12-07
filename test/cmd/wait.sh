#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

run_wait_tests() {
    set -o nounset
    set -o errexit

    kube::log::status "Testing kubectl wait"

    create_and_use_new_namespace

    ### Wait for deletion using --all flag

    # create test data
    kubectl create deployment test-1 --image=busybox
    kubectl create deployment test-2 --image=busybox

    # Post-Condition: deployments exists
    kube::test::get_object_assert "deployments" "{{range .items}}{{.metadata.name}},{{end}}" 'test-1,test-2,'

    # Delete all deployments async to kubectl wait
    ( sleep 2 && kubectl delete deployment --all ) &

    # Command: Wait for all deployments to be deleted
    output_message=$(kubectl wait deployment --for=delete --all)

    # Post-Condition: Wait was successful
    kube::test::if_has_string "${output_message}" 'test-1 condition met'
    kube::test::if_has_string "${output_message}" 'test-2 condition met'

    set +o nounset
    set +o errexit
}