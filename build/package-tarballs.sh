#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

# Complete the release with the standard env
KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/build/common.sh"
source "${KUBE_ROOT}/build/lib/release.sh"

kube::build::ensure_tar
kube::version::get_version_vars
kube::release::package_tarballs
