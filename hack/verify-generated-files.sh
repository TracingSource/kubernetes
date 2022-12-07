#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
export KUBE_ROOT
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::util::ensure_clean_working_dir

_tmpdir="$(kube::realpath "$(mktemp -d -t verify-generated-files.XXXXXX)")"
kube::util::trap_add "rm -rf ${_tmpdir}" EXIT

_tmp_gopath="${_tmpdir}/go"
_tmp_kuberoot="${_tmp_gopath}/src/k8s.io/kubernetes"
mkdir -p "${_tmp_kuberoot}/.."
cp -a "${KUBE_ROOT}" "${_tmp_kuberoot}/.."

cd "${_tmp_kuberoot}"

# clean out anything from the temp dir that's not checked in
git clean -ffxd
# regenerate any generated code
make generated_files

changed_files=$(git status --porcelain)

if [[ -n "${changed_files}" ]]; then
  echo "!!! Generated code is out of date:" >&2
  echo "${changed_files}" >&2
  echo >&2
  echo "Please run make generated_files." >&2
  exit 1
fi
