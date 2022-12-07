#!/usr/bin/env bash

# This script is a vestigial redirection.  Please do not add "real" logic.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# For help output
ARGHELP=""
if [[ -n "${KUBE_VERIFY_GIT_BRANCH:-}" ]]; then
    ARGHELP="BRANCH=${KUBE_VERIFY_GIT_BRANCH}"
fi

echo "NOTE: $0 has been replaced by 'make verify'"
echo
echo "The equivalent of this invocation is: "
echo "    make verify ${ARGHELP}"
echo
echo
make --no-print-directory -C "${KUBE_ROOT}" verify BRANCH="${KUBE_VERIFY_GIT_BRANCH:-}"
