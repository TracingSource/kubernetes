package node

import (
	"fmt"
	"strings"

	"k8s.io/kubernetes/test/e2e/framework"
	e2essh "k8s.io/kubernetes/test/e2e/framework/ssh"

	"github.com/onsi/ginkgo"
)

var _ = SIGDescribe("crictl", func() {
	f := framework.NewDefaultFramework("crictl")

	ginkgo.BeforeEach(func() {
		// `crictl` is not available on all cloud providers.
		framework.SkipUnlessProviderIs("gce", "gke")
		// The test requires $HOME/.ssh/id_rsa key to be present.
		framework.SkipUnlessSSHKeyPresent()
	})

	ginkgo.It("should be able to run crictl on the node", func() {
		// Get all nodes' external IPs.
		ginkgo.By("Getting all nodes' SSH-able IP addresses")
		hosts, err := e2essh.NodeSSHHosts(f.ClientSet)
		if err != nil {
			framework.Failf("Error getting node hostnames: %v", err)
		}

		testCases := []struct {
			cmd string
		}{
			{`sudo crictl version`},
			{`sudo crictl info`},
		}

		for _, testCase := range testCases {
			// Choose an arbitrary node to test.
			host := hosts[0]
			ginkgo.By(fmt.Sprintf("SSH'ing to node %q to run %q", host, testCase.cmd))

			result, err := e2essh.SSH(testCase.cmd, host, framework.TestContext.Provider)
			stdout, stderr := strings.TrimSpace(result.Stdout), strings.TrimSpace(result.Stderr)
			if err != nil {
				framework.Failf("Ran %q on %q, got error %v", testCase.cmd, host, err)
			}
			// Log the stdout/stderr output.
			// TODO: Verify the output.
			if len(stdout) > 0 {
				framework.Logf("Got stdout from %q:\n %s\n", host, strings.TrimSpace(stdout))
			}
			if len(stderr) > 0 {
				framework.Logf("Got stderr from %q:\n %s\n", host, strings.TrimSpace(stderr))
			}
		}
	})
})
