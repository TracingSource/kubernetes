package upgrades

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/test/e2e/framework"
	jobutil "k8s.io/kubernetes/test/e2e/framework/job"
	"k8s.io/kubernetes/test/e2e/upgrades"

	"github.com/onsi/ginkgo"
)

// JobUpgradeTest is a test harness for batch Jobs.
type JobUpgradeTest struct {
	job       *batchv1.Job
	namespace string
}

// Name returns the tracking name of the test.
func (JobUpgradeTest) Name() string { return "[sig-apps] job-upgrade" }

// Setup starts a Job with a parallelism of 2 and 2 completions running.
func (t *JobUpgradeTest) Setup(f *framework.Framework) {
	t.namespace = f.Namespace.Name

	ginkgo.By("Creating a job")
	t.job = jobutil.NewTestJob("notTerminate", "foo", v1.RestartPolicyOnFailure, 2, 2, nil, 6)
	job, err := jobutil.CreateJob(f.ClientSet, t.namespace, t.job)
	t.job = job
	framework.ExpectNoError(err)

	ginkgo.By("Ensuring active pods == parallelism")
	err = jobutil.WaitForAllJobPodsRunning(f.ClientSet, t.namespace, job.Name, 2)
	framework.ExpectNoError(err)
}

// Test verifies that the Jobs Pods are running after the an upgrade
func (t *JobUpgradeTest) Test(f *framework.Framework, done <-chan struct{}, upgrade upgrades.UpgradeType) {
	<-done
	ginkgo.By("Ensuring active pods == parallelism")
	err := jobutil.EnsureAllJobPodsRunning(f.ClientSet, t.namespace, t.job.Name, 2)
	framework.ExpectNoError(err)
}

// Teardown cleans up any remaining resources.
func (t *JobUpgradeTest) Teardown(f *framework.Framework) {
	// rely on the namespace deletion to clean up everything
}
