#!/usr/bin/env bash

# Common utility functions for build scripts

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

function kube::release::semantic_version() {
  # This takes:
  # Client Version: version.Info{Major:"1", Minor:"1+", GitVersion:"v1.1.0-alpha.0.2328+3c0a05de4a38e3", GitCommit:"3c0a05de4a38e355d147dbfb4d85bad6d2d73bb9", GitTreeState:"clean"}
  # and spits back the GitVersion piece in a way that is somewhat
  # resilient to the other fields changing (we hope)
  "${KUBE_ROOT}/cluster/kubectl.sh" version --client | sed "s/, */\\
/g" | grep -E "^GitVersion:" | cut -f2 -d: | cut -f2 -d\"
}

function kube::release::semantic_image_tag_version() {
    printf "%s" "$(kube::release::semantic_version)" | tr + _
}
