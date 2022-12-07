#!/usr/bin/env bash

# This script assumes that kubectl binary is present in PATH.
kubectl config set-cluster hollow-cluster --server=http://localhost:8080 --insecure-skip-tls-verify=true
kubectl config set-credentials "$(whoami)"
kubectl config set-context hollow-context --cluster=hollow-cluster --user="$(whoami)"
