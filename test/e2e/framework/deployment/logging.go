package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	testutils "k8s.io/kubernetes/test/utils"
)

func logReplicaSetsOfDeployment(deployment *appsv1.Deployment, allOldRSs []*appsv1.ReplicaSet, newRS *appsv1.ReplicaSet) {
	testutils.LogReplicaSetsOfDeployment(deployment, allOldRSs, newRS, framework.Logf)
}

func logPodsOfDeployment(c clientset.Interface, deployment *appsv1.Deployment, rsList []*appsv1.ReplicaSet) {
	testutils.LogPodsOfDeployment(c, deployment, rsList, framework.Logf)
}
