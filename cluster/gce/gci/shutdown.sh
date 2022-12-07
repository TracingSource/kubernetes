#!/usr/bin/env bash

# A script that let's gci preemptible nodes gracefully terminate in the event of a VM shutdown.
preemptible=$(curl "http://metadata.google.internal/computeMetadata/v1/instance/scheduling/preemptible" -H "Metadata-Flavor: Google")
if [ "${preemptible}" == "TRUE" ]; then
    echo "Shutting down! Sleeping for a minute to let the node gracefully terminate"
    # https://cloud.google.com/compute/docs/instances/stopping-or-deleting-an-instance#delete_timeout
    sleep 30
fi
