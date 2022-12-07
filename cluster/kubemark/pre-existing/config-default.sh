#!/usr/bin/env bash
# Configuration for landing a Kubemark cluster on a pre-existing Kubernetes
# cluster.

# Pre-existing provider expects a MASTER_IP.
# If you need to specify a port that's not the default (443), add it to MASTER_IP.
#
# Example: Connect to the Master on the secure port 6443
#          MASTER_IP=192.168.122.5:6443
#
MASTER_IP="${MASTER_IP:-}"

# The container registry and project given to the kubemark container:
#   $CONTAINER_REGISTRY/$PROJECT/kubemark
#
CONTAINER_REGISTRY="${CONTAINER_REGISTRY:-}"
PROJECT="${PROJECT:-}"

NUM_NODES="${NUM_NODES:-1}"

TEST_CLUSTER_API_CONTENT_TYPE="${TEST_CLUSTER_API_CONTENT_TYPE:-}"
KUBELET_TEST_LOG_LEVEL="${KUBELET_TEST_LOG_LEVEL:-}"
KUBEPROXY_TEST_LOG_LEVEL="${KUBEPROXY_TEST_LOG_LEVEL:-}"
MASTER_NAME="${MASTER_NAME:-}"
