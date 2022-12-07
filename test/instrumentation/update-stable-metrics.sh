#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
BAZEL_OUT_DIR="$KUBE_ROOT/bazel-bin"
BAZEL_GEN_DIR="$KUBE_ROOT/bazel-genfiles"
METRICS_LIST_PATH="test/instrumentation/stable-metrics-list.yaml"

bazel build //test/instrumentation:list_stable_metrics
if [ -d "$BAZEL_OUT_DIR" ]; then
  cp "$BAZEL_OUT_DIR/$METRICS_LIST_PATH" "$KUBE_ROOT/test/instrumentation/testdata/stable-metrics-list.yaml"
else
  # Handle bazel < 0.25
  # https://github.com/bazelbuild/bazel/issues/6761
  echo "$BAZEL_OUT_DIR not found trying $BAZEL_GEN_DIR"
  cp "$BAZEL_GEN_DIR/$METRICS_LIST_PATH" "$KUBE_ROOT/test/instrumentation/testdata/stable-metrics-list.yaml"
fi
