#!/usr/bin/env bash

# Smoke Tests a running Kubernetes cluster.
# Validates that the cluster was deployed, is accessible, and at least
# satisfies minimal functional requirements.
# Emphasis on speed and being non-destructive over thoroughness.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

SMOKE_TEST_FOCUS_REGEX="Guestbook.application"

exec "${KUBE_ROOT}/cluster/test-e2e.sh" -ginkgo.focus="${SMOKE_TEST_FOCUS_REGEX}" "$@"
