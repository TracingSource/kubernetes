package utils

import (
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	e2enode "k8s.io/kubernetes/test/e2e/framework/node"
)

// GetNodeIds returns the list of node names and panics in case of failure.
func GetNodeIds(cs clientset.Interface) []string {
	nodes, err := e2enode.GetReadySchedulableNodes(cs)
	framework.ExpectNoError(err)
	nodeIds := []string{}
	for _, n := range nodes.Items {
		nodeIds = append(nodeIds, n.Name)
	}
	return nodeIds
}
