#!/usr/bin/env bash

# This file is not intended to be run automatically. It is meant to be run
# immediately before exporting docs. We do not want to check these documents in
# by default.

set -o errexit
set -o nounset
set -o pipefail

KUBE_HACK_ROOT=$(dirname "${BASH_SOURCE[0]}")

echo "WARNING: hack/generate-docs.sh is an alias for hack/update-generated-docs.sh"
echo "and will be removed in a future version."

"${KUBE_HACK_ROOT}"/update-generated-docs.sh
