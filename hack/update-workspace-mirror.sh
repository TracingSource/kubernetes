#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [[ $# -ne 1 ]]; then
    echo 'use "bazel run //hack:update-mirror"'
    echo "(usage: $0 <file with list of URLs to mirror>)"
    exit 1
fi

BUCKET="gs://k8s-bazel-cache"

gsutil acl get "${BUCKET}" > /dev/null

tmpfile=$(mktemp bazel_workspace_mirror.XXXXXX)
trap 'rm ${tmpfile}' EXIT
while read -r url; do
  echo "${url}"
  if gsutil ls "${BUCKET}/${url}" &> /dev/null; then
    echo present
  else
    echo missing
    if curl -fLag "${url}" > "${tmpfile}"; then
        gsutil cp -a public-read "${tmpfile}" "${BUCKET}/${url}"
    fi
  fi
done < "$1"
