package node

import "k8s.io/kubernetes/test/e2e/framework"

// SIGDescribe annotates the test with the SIG label.
func SIGDescribe(text string, body func()) bool {
	return framework.KubeDescribe("[sig-node] "+text, body)
}
