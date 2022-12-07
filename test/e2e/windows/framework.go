package windows

import (
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/onsi/ginkgo"
)

// SIGDescribe annotates the test with the SIG label.
func SIGDescribe(text string, body func()) bool {
	return ginkgo.Describe("[sig-windows] "+text, func() {
		ginkgo.BeforeEach(func() {
			// all tests in this package are Windows specific
			framework.SkipUnlessNodeOSDistroIs("windows")
		})

		body()
	})
}
