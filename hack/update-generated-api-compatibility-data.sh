#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

# UPDATE_COMPATIBILITY_FIXTURE_DATA=true regenerates fixture data if needed.
# -run //HEAD only runs the test cases comparing against testdata for HEAD.
# We suppress the output because we are expecting to have changes.
# We suppress the test failure that occurs when there are changes.
UPDATE_COMPATIBILITY_FIXTURE_DATA=true go test ./vendor/k8s.io/api -run //HEAD >/dev/null 2>&1 || true

# Now that we have regenerated data at HEAD, run the test without suppressing output or failures
go test ./vendor/k8s.io/api -run //HEAD -count=1
