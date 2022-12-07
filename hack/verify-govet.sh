#!/usr/bin/env bash

# This script is a vestigial redirection.  Please do not add "real" logic.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# For help output
ARGHELP=""
if [[ "$#" -gt 0 ]]; then
    ARGHELP="WHAT='$*'"
fi

echo "NOTE: $0 has been replaced by 'make vet'"
echo
echo "The equivalent of this invocation is: "
echo "    make vet ${ARGHELP}"
echo
echo

make --no-print-directory -C "${KUBE_ROOT}" vet WHAT="$*"
