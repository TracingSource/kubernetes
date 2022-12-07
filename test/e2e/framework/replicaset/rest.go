package replicaset

import (
	appsv1 "k8s.io/api/apps/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	testutils "k8s.io/kubernetes/test/utils"
)

// UpdateReplicaSetWithRetries updates replicaset template with retries.
func UpdateReplicaSetWithRetries(c clientset.Interface, namespace, name string, applyUpdate testutils.UpdateReplicaSetFunc) (*appsv1.ReplicaSet, error) {
	return testutils.UpdateReplicaSetWithRetries(c, namespace, name, applyUpdate, framework.Logf, framework.Poll, framework.PollShortTimeout)
}
