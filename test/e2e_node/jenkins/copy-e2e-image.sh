#!/usr/bin/env bash

# Usage: copy-e2e-image.sh <image-name> <from-project-name> <to-project-name>

# See *.properties for list of images to copy,
# typically from kubernetes-node-e2e-images

set -e

print_usage() {
    echo "This script helps copy a GCE image from a source to a target project"
    echo -e "\nUsage:\n$0 <from-image-name> <from-project-name> <to-project-name> <to-image-name>\n"
}

if [  $# -ne 4 ]; then
    print_usage
    exit 1
fi

FROM_IMAGE=$1
FROM_PROJECT=$2
TO_PROJECT=$3
TO_IMAGE=$4

echo "Copying image $FROM_IMAGE from project $FROM_PROJECT to project $TO_PROJECT as image $TO_IMAGE..."
gcloud compute --project "${TO_PROJECT}" disks create "${TO_IMAGE}" --image=https://www.googleapis.com/compute/v1/projects/"${FROM_PROJECT}/global/images/""${FROM_IMAGE}"
gcloud compute --project "${TO_PROJECT}" images create "${TO_IMAGE}" \
  --source-disk="${TO_IMAGE}" \
  --description="Cloned from projects/$2/global/images/$1 by $USER on $(date)"
gcloud -q compute --project "${TO_PROJECT}" disks delete "${TO_IMAGE}"
