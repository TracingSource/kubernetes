#!/usr/bin/env bash

# Script to fetch latest openapi spec.
# Puts the updated spec at api/openapi-spec/

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
OPENAPI_ROOT_DIR="${KUBE_ROOT}/api/openapi-spec"
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::util::require-jq
kube::golang::setup_env

make -C "${KUBE_ROOT}" WHAT=cmd/kube-apiserver

function cleanup()
{
    if [[ -n ${APISERVER_PID-} ]]; then
      kill "${APISERVER_PID}" 1>&2 2>/dev/null
      wait "${APISERVER_PID}" || true
    fi
    unset APISERVER_PID

    kube::etcd::cleanup

    kube::log::status "Clean up complete"
}

trap cleanup EXIT SIGINT

kube::golang::setup_env

TMP_DIR=$(mktemp -d /tmp/update-openapi-spec.XXXX)
ETCD_HOST=${ETCD_HOST:-127.0.0.1}
ETCD_PORT=${ETCD_PORT:-2379}
API_PORT=${API_PORT:-8050}
API_HOST=${API_HOST:-127.0.0.1}
API_LOGFILE=${API_LOGFILE:-/tmp/openapi-api-server.log}

kube::etcd::start

echo "dummy_token,admin,admin" > "${TMP_DIR}/tokenauth.csv"

# Start kube-apiserver
kube::log::status "Starting kube-apiserver"
"${KUBE_OUTPUT_HOSTBIN}/kube-apiserver" \
  --insecure-bind-address="${API_HOST}" \
  --bind-address="${API_HOST}" \
  --insecure-port="${API_PORT}" \
  --etcd-servers="http://${ETCD_HOST}:${ETCD_PORT}" \
  --advertise-address="10.10.10.10" \
  --cert-dir="${TMP_DIR}/certs" \
  --runtime-config="api/all=true,extensions/v1beta1/daemonsets=true,extensions/v1beta1/deployments=true,extensions/v1beta1/replicasets=true,extensions/v1beta1/networkpolicies=true,extensions/v1beta1/podsecuritypolicies=true,extensions/v1beta1/replicationcontrollers=true" \
  --token-auth-file="${TMP_DIR}/tokenauth.csv" \
  --service-account-issuer="https://kubernetes.devault.svc/" \
  --service-account-signing-key-file="${KUBE_ROOT}/staging/src/k8s.io/client-go/util/cert/testdata/dontUseThisKey.pem" \
  --logtostderr \
  --v=2 \
  --service-cluster-ip-range="10.0.0.0/24" >"${API_LOGFILE}" 2>&1 &
APISERVER_PID=$!

if ! kube::util::wait_for_url "${API_HOST}:${API_PORT}/healthz" "apiserver: "; then
  kube::log::error "Here are the last 10 lines from kube-apiserver (${API_LOGFILE})"
  kube::log::error "=== BEGIN OF LOG ==="
  tail -10 "${API_LOGFILE}" || :
  kube::log::error "=== END OF LOG ==="
  exit 1
fi

kube::log::status "Updating " "${OPENAPI_ROOT_DIR}"

curl -w "\n" -fs "${API_HOST}:${API_PORT}/openapi/v2" | jq -S '.info.version="unversioned"' > "${OPENAPI_ROOT_DIR}/swagger.json"

kube::log::status "SUCCESS"

# ex: ts=2 sw=2 et filetype=sh
