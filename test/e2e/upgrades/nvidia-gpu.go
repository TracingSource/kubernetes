package upgrades

import (
	"k8s.io/kubernetes/test/e2e/framework"
	jobutil "k8s.io/kubernetes/test/e2e/framework/job"
	"k8s.io/kubernetes/test/e2e/scheduling"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	completions = int32(1)
)

// NvidiaGPUUpgradeTest tests that gpu resource is available before and after
// a cluster upgrade.
type NvidiaGPUUpgradeTest struct {
}

// Name returns the tracking name of the test.
func (NvidiaGPUUpgradeTest) Name() string { return "nvidia-gpu-upgrade [sig-node] [sig-scheduling]" }

// Setup creates a job requesting gpu.
func (t *NvidiaGPUUpgradeTest) Setup(f *framework.Framework) {
	scheduling.SetupNVIDIAGPUNode(f, false)
	ginkgo.By("Creating a job requesting gpu")
	scheduling.StartJob(f, completions)
}

// Test waits for the upgrade to complete, and then verifies that the
// cuda pod started by the gpu job can successfully finish.
func (t *NvidiaGPUUpgradeTest) Test(f *framework.Framework, done <-chan struct{}, upgrade UpgradeType) {
	<-done
	ginkgo.By("Verifying gpu job success")
	scheduling.VerifyJobNCompletions(f, completions)
	if upgrade == MasterUpgrade || upgrade == ClusterUpgrade {
		// MasterUpgrade should be totally hitless.
		job, err := jobutil.GetJob(f.ClientSet, f.Namespace.Name, "cuda-add")
		framework.ExpectNoError(err)
		gomega.Expect(job.Status.Failed).To(gomega.BeZero(), "Job pods failed during master upgrade: %v", job.Status.Failed)
	}
}

// Teardown cleans up any remaining resources.
func (t *NvidiaGPUUpgradeTest) Teardown(f *framework.Framework) {
	// rely on the namespace deletion to clean up everything
}
