#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# NOTE: All output from this script needs to be copied back to the calling
# source tree.  This is managed in kube::build::copy_output in build/common.sh.
# If the output set is changed update that function.

APIROOTS=${APIROOTS:-$(git grep --files-with-matches -e '// +k8s:protobuf-gen=package' cmd pkg staging | \
	xargs -n 1 dirname | \
	sed 's,^,k8s.io/kubernetes/,;s,k8s.io/kubernetes/staging/src/,,' | \
	sort | uniq
)}

"${KUBE_ROOT}/build/run.sh" hack/update-generated-protobuf-dockerized.sh "${APIROOTS}" "$@"

# ex: ts=2 sw=2 et filetype=sh
