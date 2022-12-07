#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

kubectl create -f conformance-e2e.yaml
while true; do
  STATUS=$(kubectl -n conformance get pods e2e-conformance-test -o jsonpath="{.status.phase}")
  timestamp=$(date +"[%H:%M:%S]")
  echo "$timestamp Pod status is: ${STATUS}"
  if [[ "$STATUS" == "Succeeded" ]]; then
    echo "$timestamp Done."
    break
  elif [[ "$STATUS" == "Failed" ]]; then
    echo "$timestamp Failed."
    break
  else
    sleep 5
  fi
done
echo "Please use 'kubectl logs -n conformance e2e-conformance-test' to view the results"
