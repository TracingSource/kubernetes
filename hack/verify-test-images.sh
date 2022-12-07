#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

cd "${KUBE_ROOT}"
result=0

# Find mentions of untagged gcr.io images in test/e2e/*.go
find_e2e_test_untagged_gcr_images() {
    grep -o -E -e 'gcr.io/[-a-z0-9/_:.]+' test/e2e/*.go | grep -v -E "gcr.io/.*:" | cut -d ":" -f 1 | LC_ALL=C sort -u
}


# Find mentions of latest gcr.io images in test/e2e/*.go
find_e2e_test_latest_gcr_images() {
    grep -o -E -e 'gcr.io/.*:latest' test/e2e/*.go | cut -d ":" -f 1 | LC_ALL=C sort -u
}

if find_e2e_test_latest_gcr_images; then
  echo "!!! Found :latest gcr.io images in the above files"
  result=1
fi

if find_e2e_test_untagged_gcr_images; then
  echo "!!! Found untagged gcr.io images in the above files"
  result=1
fi

exit ${result}
