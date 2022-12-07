#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

BINS=(
	cmd/clicheck
)
make -C "${KUBE_ROOT}" WHAT="${BINS[*]}"

clicheck=$(kube::util::find-binary "clicheck")

if ! output=$($clicheck 2>&1)
then
	echo "$output"
	echo
	echo "FAILURE: CLI is not following one or more required conventions."
	exit 1
else
	echo "$output"
	echo
  echo "SUCCESS: CLI is following all tested conventions."
fi
