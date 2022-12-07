#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# generate the code with:
# - --output-base because this script should also be able to run inside the vendor dir of
#   k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#   instead of the $GOPATH directly. For normal projects this can be dropped.
"$(dirname "${BASH_SOURCE[0]}")"/../generate-internal-groups.sh all \
  k8s.io/code-generator/_examples/apiserver k8s.io/code-generator/_examples/apiserver/apis k8s.io/code-generator/_examples/apiserver/apis \
  "example:v1 example2:v1 example3.io:v1" \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
"$(dirname "${BASH_SOURCE[0]}")"/../generate-groups.sh all \
  k8s.io/code-generator/_examples/crd k8s.io/code-generator/_examples/crd/apis \
  "example:v1 example2:v1" \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
"$(dirname "${BASH_SOURCE[0]}")"/../generate-groups.sh all \
  k8s.io/code-generator/_examples/MixedCase k8s.io/code-generator/_examples/MixedCase/apis \
  "example:v1" \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
"$(dirname "${BASH_SOURCE[0]}")"/../generate-groups.sh all \
  k8s.io/code-generator/_examples/HyphenGroup k8s.io/code-generator/_examples/HyphenGroup/apis \
  "example:v1" \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"

