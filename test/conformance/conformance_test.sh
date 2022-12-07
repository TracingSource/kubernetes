#!/usr/bin/env bash

# echo ${TEST_SRCDIR}
# pwd
# env | grep --color=always test/conformance
# find ${TEST_SRCDIR} -ls | grep --color=always test/conformance

set -o errexit

if diff -u test/conformance/testdata/conformance.txt test/conformance/conformance.txt; then
  echo PASS
  exit 0
fi
echo 'See instructions in test/conformance/README.md'
exit 1
