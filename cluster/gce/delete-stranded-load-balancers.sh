#!/usr/bin/env bash

# A utility for deleting target pools and forwarding rules that are unattached to VMs
PROJECT=${PROJECT:-kubernetes-jenkins}
REGION=${REGION:-us-central1}

LIST=$(gcloud --project="${PROJECT}" compute target-pools list --format='value(name)')

result=0
for x in ${LIST}; do
    if ! gcloud compute --project="${PROJECT}" target-pools get-health "${x}" --region="${REGION}" 2>/dev/null >/dev/null; then
	echo DELETING "${x}"
	gcloud compute --project="${PROJECT}" firewall-rules delete "k8s-fw-${x}" -q
	gcloud compute --project="${PROJECT}" forwarding-rules delete "${x}" --region="${REGION}" -q
	gcloud compute --project="${PROJECT}" addresses delete "${x}" --region="${REGION}" -q
	gcloud compute --project="${PROJECT}" target-pools delete "${x}" --region="${REGION}" -q
        result=1
    fi
done
exit ${result}

