#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

BINS=(
	cmd/gendocs
	cmd/genkubedocs
	cmd/genman
	cmd/genyaml
)
make -C "${KUBE_ROOT}" WHAT="${BINS[*]}"

kube::util::ensure-temp-dir

# just verify the generation process
kube::util::gen-docs "${KUBE_TEMP}"
