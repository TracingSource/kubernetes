package kubeadm

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/gomega"
	"github.com/spf13/pflag"

	morereporters "github.com/onsi/ginkgo/reporters"
	"k8s.io/kubernetes/test/e2e/framework"
	e2econfig "k8s.io/kubernetes/test/e2e/framework/config"
)

func TestMain(m *testing.M) {
	// Copy go flags in TestMain, to ensure go test flags are registered (no longer available in init() as of go1.13)
	e2econfig.CopyFlags(e2econfig.Flags, flag.CommandLine)
	framework.RegisterCommonFlags(flag.CommandLine)
	framework.RegisterClusterFlags(flag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	framework.AfterReadingAllFlags(&framework.TestContext)
	os.Exit(m.Run())
}

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	reporters := []ginkgo.Reporter{}
	reportDir := framework.TestContext.ReportDir
	if reportDir != "" {
		// Create the directory if it doesn't already exists
		if err := os.MkdirAll(reportDir, 0755); err != nil {
			t.Fatalf("Failed creating report directory: %v", err)
		} else {
			// Configure a junit reporter to write to the directory
			junitFile := fmt.Sprintf("junit_%s_%02d.xml", framework.TestContext.ReportPrefix, config.GinkgoConfig.ParallelNode)
			junitPath := filepath.Join(reportDir, junitFile)
			reporters = append(reporters, morereporters.NewJUnitReporter(junitPath))
		}
	}
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "E2EKubeadm suite", reporters)
}
