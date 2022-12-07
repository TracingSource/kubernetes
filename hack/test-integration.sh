#!/usr/bin/env bash

# This script is a vestigial redirection.  Please do not add "real" logic.

set -o errexit
set -o nounset
set -o pipefail

echo "$0 has been replaced by 'make test-integration'"
echo
echo "The following invocation will run all integration tests: "
echo '    make test-integration'
echo
echo "The following invocation will run a specific test with the verbose flag set: "
echo '    make test-integration WHAT=./test/integration/pods GOFLAGS="-v" KUBE_TEST_ARGS="-run ^TestPodUpdateActiveDeadlineSeconds$"'
echo
exit 1
