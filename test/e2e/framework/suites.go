package framework

import (
	"fmt"
	"io/ioutil"
	"path"
	"time"

	// TODO: Remove the following imports (ref: https://github.com/kubernetes/kubernetes/issues/81245)
	e2emetrics "k8s.io/kubernetes/test/e2e/framework/metrics"
)

// CleanupSuite is the boilerplate that can be used after tests on ginkgo were run, on the SynchronizedAfterSuite step.
// Similar to SynchronizedBeforeSuite, we want to run some operations only once (such as collecting cluster logs).
// Here, the order of functions is reversed; first, the function which runs everywhere,
// and then the function that only runs on the first Ginkgo node.
func CleanupSuite() {
	// Run on all Ginkgo nodes
	Logf("Running AfterSuite actions on all nodes")
	RunCleanupActions()
}

// AfterSuiteActions are actions that are run on ginkgo's SynchronizedAfterSuite
func AfterSuiteActions() {
	// Run only Ginkgo on node 1
	Logf("Running AfterSuite actions on node 1")
	if TestContext.ReportDir != "" {
		CoreDump(TestContext.ReportDir)
	}
	if TestContext.GatherSuiteMetricsAfterTest {
		if err := gatherTestSuiteMetrics(); err != nil {
			Logf("Error gathering metrics: %v", err)
		}
	}
	if TestContext.NodeKiller.Enabled {
		close(TestContext.NodeKiller.NodeKillerStopCh)
	}
}

func gatherTestSuiteMetrics() error {
	Logf("Gathering metrics")
	c, err := LoadClientset()
	if err != nil {
		return fmt.Errorf("error loading client: %v", err)
	}

	// Grab metrics for apiserver, scheduler, controller-manager, kubelet (for non-kubemark case) and cluster autoscaler (optionally).
	grabber, err := e2emetrics.NewMetricsGrabber(c, nil, !ProviderIs("kubemark"), true, true, true, TestContext.IncludeClusterAutoscalerMetrics)
	if err != nil {
		return fmt.Errorf("failed to create MetricsGrabber: %v", err)
	}

	received, err := grabber.Grab()
	if err != nil {
		return fmt.Errorf("failed to grab metrics: %v", err)
	}

	metricsForE2E := (*e2emetrics.ComponentCollection)(&received)
	metricsJSON := metricsForE2E.PrintJSON()
	if TestContext.ReportDir != "" {
		filePath := path.Join(TestContext.ReportDir, "MetricsForE2ESuite_"+time.Now().Format(time.RFC3339)+".json")
		if err := ioutil.WriteFile(filePath, []byte(metricsJSON), 0644); err != nil {
			return fmt.Errorf("error writing to %q: %v", filePath, err)
		}
	} else {
		Logf("\n\nTest Suite Metrics:\n%s\n", metricsJSON)
	}

	return nil
}
