#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

cd "${KUBE_ROOT}"

rc=0

# find test files accessing the mutable global feature gate or interface
direct_sets=$(find -L . -name '*_test.go' -exec grep -Hn 'MutableFeatureGate' {} \; 2>/dev/null) || true
if [[ -n "${direct_sets}" ]]; then
  echo "Test files may not access mutable global feature gates directly:" >&2
  echo "${direct_sets}" >&2
  echo >&2
  echo "Use this invocation instead:" >&2
  echo "  defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.<FeatureName>, <value>)()" >&2
  echo >&2
  rc=1
fi

# find test files calling SetFeatureGateDuringTest and not calling the result
missing_defers=$(find -L . -name '*_test.go' -exec grep -Hn 'SetFeatureGateDuringTest' {} \; 2>/dev/null | grep -E -v "defer .*\\)\\(\\)$") || true
if [[ -n "${missing_defers}" ]]; then
  echo "Invalid invocations of featuregatetesting.SetFeatureGateDuringTest():" >&2
  echo "${missing_defers}" >&2
  echo >&2
  echo "Always make a deferred call to the returned function to ensure the feature gate is reset:" >&2
  echo "  defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.<FeatureName>, <value>)()" >&2
  echo >&2
  rc=1
fi

exit $rc
