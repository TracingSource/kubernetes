#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

make -C "${KUBE_ROOT}" WHAT=cmd/linkcheck

linkcheck=$(kube::util::find-binary "linkcheck")

kube::util::ensure-temp-dir
OUTPUT="${KUBE_TEMP}"/linkcheck-output
cleanup() {
	rm -rf "${OUTPUT}"
}
trap "cleanup" EXIT SIGINT
mkdir -p "$OUTPUT"

APIROOT="${KUBE_ROOT}/pkg/api/"
APISROOT="${KUBE_ROOT}/pkg/apis/"
DOCROOT="${KUBE_ROOT}/docs/"
ROOTS=("$APIROOT" "$APISROOT" "$DOCROOT")
found_invalid=false
for root in "${ROOTS[@]}"; do
  "${linkcheck}" "--root-dir=${root}" 2> >(tee -a "${OUTPUT}/error" >&2) && ret=0 || ret=$?
  if [[ $ret -eq 1 ]]; then
    echo "Failed: found invalid links in ${root}."
    found_invalid=true
  fi
  if [[ $ret -gt 1 ]]; then
    echo "Error running linkcheck"
    exit 1
  fi
done

if [ ${found_invalid} = true ]; then
  echo "Summary of invalid links:"
  cat "${OUTPUT}/error"
  exit 1
fi

# ex: ts=2 sw=2 et filetype=sh
