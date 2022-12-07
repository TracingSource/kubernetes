#!/usr/bin/env bash

# GoFmt apparently is changing @ head...

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

cd "${KUBE_ROOT}"

# Prefer bazel's gofmt.
gofmt="external/io_bazel_rules_go_toolchain/bin/gofmt"
if [[ ! -x "${gofmt}" ]]; then
  gofmt=$(which gofmt)
  kube::golang::verify_go_version
fi

find_files() {
  find . -not \( \
      \( \
        -wholename './output' \
        -o -wholename './.git' \
        -o -wholename './_output' \
        -o -wholename './_gopath' \
        -o -wholename './release' \
        -o -wholename './target' \
        -o -wholename '*/third_party/*' \
        -o -wholename '*/vendor/*' \
        -o -wholename './staging/src/k8s.io/client-go/*vendor/*' \
        -o -wholename '*/bindata.go' \
      \) -prune \
    \) -name '*.go'
}

# gofmt exits with non-zero exit code if it finds a problem unrelated to
# formatting (e.g., a file does not parse correctly). Without "|| true" this
# would have led to no useful error message from gofmt, because the script would
# have failed before getting to the "echo" in the block below.
diff=$(find_files | xargs "${gofmt}" -d -s 2>&1) || true
if [[ -n "${diff}" ]]; then
  echo "${diff}" >&2
  echo >&2
  echo "Run ./hack/update-gofmt.sh" >&2
  exit 1
fi
