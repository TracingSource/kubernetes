package kubeadm

import "k8s.io/kubernetes/test/e2e/framework"

// Describe annotates the test with the Kubeadm label.
func Describe(text string, body func()) bool {
	return framework.KubeDescribe("[sig-cluster-lifecycle] [area-kubeadm] "+text, body)
}
