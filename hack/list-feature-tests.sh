#!/usr/bin/env bash

# A single script that lists all of the [Feature:.+] tests in our e2e suite.
set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
grep "\[Feature:\w+\]" "${KUBE_ROOT}"/test/e2e/**/*.go -Eoh | LC_ALL=C sort -u
