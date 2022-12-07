#!/usr/bin/env bash

# verify-util-pkg.sh checks whether *.go except doc.go in pkg/util have been moved into
# sub-pkgs, see issue #15634.

set -o errexit
set -o nounset
set -o pipefail

BASH_DIR=$(dirname "${BASH_SOURCE[0]}")

find_go_files() {
  find . -maxdepth 1 -not \( \
      \( \
        -wholename './doc.go' \
      \) -prune \
    \) -name '*.go'
}

ret=0

pushd "${BASH_DIR}" > /dev/null
  for path in $(find_go_files); do
    file=$(basename "$path")
    echo "Found pkg/util/${file}, but should be moved into util sub-pkgs." 1>&2
    ret=1
  done
popd > /dev/null

if [[ ${ret} -gt 0 ]]; then
  exit ${ret}
fi

echo "Util Package Verified."
