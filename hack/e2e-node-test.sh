#!/usr/bin/env bash

# This script is a vestigial redirection.  Please do not add "real" logic.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# For help output
ARGHELP=""
if [[ -n "${FOCUS:-}" ]]; then
    ARGHELP="FOCUS='${FOCUS}' "
fi
if [[ -n "${SKIP:-}" ]]; then
    ARGHELP="${ARGHELP}SKIP='${SKIP}'"
fi

echo "NOTE: $0 has been replaced by 'make test-e2e-node'"
echo
echo "This script supports a number of parameters passed as environment variables."
echo "Please see the Makefile for more details."
echo
echo "The equivalent of this invocation is: "
echo "    make test-e2e-node ${ARGHELP}"
echo
echo
make --no-print-directory -C "${KUBE_ROOT}" test-e2e-node FOCUS="${FOCUS:-}" SKIP="${SKIP:-}"
