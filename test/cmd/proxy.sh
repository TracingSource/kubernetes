#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

run_kubectl_local_proxy_tests() {
  set -o nounset
  set -o errexit

  kube::log::status "Testing kubectl local proxy"

  start-proxy
  check-curl-proxy-code /api/kubernetes 404
  check-curl-proxy-code /api/v1/namespaces 200
  if kube::test::if_supports_resource "metrics" ; then
    check-curl-proxy-code /metrics 200
  fi
  if kube::test::if_supports_resource "static" ; then
    check-curl-proxy-code /static/ 200
  fi
  stop-proxy

  # Make sure the in-development api is accessible by default
  start-proxy
  check-curl-proxy-code /apis 200
  check-curl-proxy-code /apis/extensions/ 200
  stop-proxy

  # Custom paths let you see everything.
  start-proxy /custom
  check-curl-proxy-code /custom/api/kubernetes 404
  check-curl-proxy-code /custom/api/v1/namespaces 200
  if kube::test::if_supports_resource "metrics" ; then
    check-curl-proxy-code /custom/metrics 200
  fi
  check-curl-proxy-code /custom/api/v1/namespaces 200
  stop-proxy

  set +o nounset
  set +o errexit
}
