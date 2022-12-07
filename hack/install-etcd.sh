#!/usr/bin/env bash

# Convenience script to download and install etcd in third_party.
# Mostly just used by CI.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::etcd::install
