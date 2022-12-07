#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

# Complete the release with the standard env
KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# Check and error if not "in-a-container"
if [[ ! -f /.dockerenv ]]; then
  echo
  echo "'make release-in-a-container' can only be used from a docker container."
  echo
  exit 1
fi

# Other dependencies: Your container should contain docker
if ! type -p docker >/dev/null 2>&1; then
  echo
  echo "'make release-in-a-container' requires a container with" \
       "docker installed."
  echo
  exit 1
fi


# First run make cross-in-a-container
make cross-in-a-container

# at the moment only make test is supported.
if [[ $KUBE_RELEASE_RUN_TESTS =~ ^[yY]$ ]]; then
  make test
fi

"${KUBE_ROOT}/build/package-tarballs.sh"
