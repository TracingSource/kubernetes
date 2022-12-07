#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../../../
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::golang::setup_env

DIR_BASENAME=$(dirname "${BASH_SOURCE[0]}")
pushd "${DIR_BASENAME}"

cleanup() {
  popd 2> /dev/null
  kube::etcd::cleanup
  kube::log::status "performance test cleanup complete"
}

trap cleanup EXIT

kube::etcd::start

# We are using the benchmark suite to do profiling. Because it only runs a few pods and
# theoretically it has less variance.
if ${RUN_BENCHMARK:-false}; then
  kube::log::status "performance test (benchmark) compiling"
  go test -c -o "perf.test"

  kube::log::status "performance test (benchmark) start"
  "./perf.test" -test.bench=. -test.run=xxxx -test.cpuprofile=prof.out -test.short=false
  kube::log::status "...benchmark tests finished"
fi
# Running density tests. It might take a long time.
kube::log::status "performance test (density) start"
go test -test.run=. -test.timeout=60m -test.short=false
kube::log::status "...density tests finished"
