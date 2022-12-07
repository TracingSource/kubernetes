#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

# create a nice clean place to put our new licenses
# must be in the user dir (e.g. KUBE_ROOT) in order for the docker volume mount
# to work with docker-machine on macs
mkdir -p "${KUBE_ROOT}/_tmp"
_tmpdir="$(mktemp -d "${KUBE_ROOT}/_tmp/kube-vendor-licenses.XXXXXX")"
#echo "Created workspace: ${_tmpdir}"
function cleanup {
  #echo "Removing workspace: ${_tmpdir}"
  rm -rf "${_tmpdir}"
}
kube::util::trap_add cleanup EXIT

cp -r "${KUBE_ROOT}/Godeps" "${_tmpdir}/Godeps"
ln -s "${KUBE_ROOT}/LICENSE" "${_tmpdir}"
ln -s "${KUBE_ROOT}/vendor" "${_tmpdir}"
ln -s "${KUBE_ROOT}/staging" "${_tmpdir}"

# Update vendor Licenses
LICENSE_ROOT="${_tmpdir}" "${KUBE_ROOT}/hack/update-vendor-licenses.sh"

# Compare vendor Licenses
if ! _out="$(diff -Naupr "${KUBE_ROOT}/Godeps/LICENSES" "${_tmpdir}/Godeps/LICENSES")"; then
  echo "Your vendor licenses file is out of date. Run hack/update-vendor-licenses.sh and commit the results." >&2
  echo "${_out}" >&2
  exit 1
fi
