package e2e_node

import "github.com/onsi/ginkgo"

func SIGDescribe(text string, body func()) bool {
	return ginkgo.Describe("[sig-node] "+text, body)
}
