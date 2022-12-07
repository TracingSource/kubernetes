#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

make test \
    WHAT="$*" \
    KUBE_COVER="" \
    KUBE_RACE=" " \
    KUBE_TEST_ARGS="-- -test.run='^X' -benchtime=1s -bench=. -benchmem" \
