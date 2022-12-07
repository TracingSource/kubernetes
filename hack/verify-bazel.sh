#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
export KUBE_ROOT
source "${KUBE_ROOT}/hack/lib/init.sh"

if [[ ! -f "${KUBE_ROOT}/vendor/BUILD" ]]; then
  echo "${KUBE_ROOT}/vendor/BUILD does not exist." >&2
  echo >&2
  echo "Run ./hack/update-bazel.sh" >&2
  exit 1
fi

# Remove generated files prior to running kazel.
# TODO(spxtr): Remove this line once Bazel is the only way to build.
rm -f "${KUBE_ROOT}/{pkg/generated,staging/src/k8s.io/apiextensions-apiserver/pkg/client,staging/src/k8s.io/kube-aggregator/pkg/client}/openapi/zz_generated.openapi.go"

_tmpdir="$(kube::realpath "$(mktemp -d -t verify-bazel.XXXXXX)")"
kube::util::trap_add "rm -rf ${_tmpdir}" EXIT

_tmp_gopath="${_tmpdir}/go"
_tmp_kuberoot="${_tmp_gopath}/src/k8s.io/kubernetes"
mkdir -p "${_tmp_kuberoot}/.."
cp -a "${KUBE_ROOT}" "${_tmp_kuberoot}/.."

cd "${_tmp_kuberoot}"
GOPATH="${_tmp_gopath}" PATH="${_tmp_gopath}/bin:${PATH}" ./hack/update-bazel.sh

diff=$(diff -Naupr -x '_output' "${KUBE_ROOT}" "${_tmp_kuberoot}" || true)

if [[ -n "${diff}" ]]; then
  echo "${diff}" >&2
  echo >&2
  echo "Run ./hack/update-bazel.sh" >&2
  exit 1
fi
