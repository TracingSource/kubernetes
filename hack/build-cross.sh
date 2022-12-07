#!/usr/bin/env bash

# This script is a vestigial redirection.  Please do not add "real" logic.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

echo "NOTE: $0 has been replaced by 'make cross'"
echo
echo "The equivalent of this invocation is: "
echo "    make cross"
echo
echo
make --no-print-directory -C "${KUBE_ROOT}" cross
