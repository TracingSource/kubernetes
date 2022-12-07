#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# NOTE: All output from this script needs to be copied back to the calling
# source tree.  This is managed in kube::build::copy_output in build/common.sh.
# If the output set is changed update that function.

"${KUBE_ROOT}/build/run.sh" hack/update-generated-runtime-dockerized.sh "$@"

# ex: ts=2 sw=2 et filetype=sh
