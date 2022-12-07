package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"k8s.io/klog"
	"k8s.io/kubernetes/test/utils"
)

func bazelBuild() error {
	targets := []string{
		"//vendor/github.com/onsi/ginkgo/ginkgo",
		"//test/e2e_kubeadm:e2e_kubeadm.test",
	}

	args := append([]string{"build"}, targets...)

	return execCommand("bazel", args...)
}

var ginkgoFlags = flag.String("ginkgo-flags", "", "Space-separated list of arguments to pass to Ginkgo test runner.")
var testFlags = flag.String("test-flags", "", "Space-separated list of arguments to pass to kubeadm e2e test.")
var build = flag.Bool("build", false, "use Bazel to build binaries before testing")

func main() {
	flag.Parse()

	if *build {
		if err := bazelBuild(); err != nil {
			klog.Exitf("couldn't build with bazel: %v", err)
		}
	}

	ginkgo, err := getBazelGinkgo()
	if err != nil {
		klog.Fatalf("Failed to get ginkgo binary: %v", err)
	}

	test, err := getBazelTestBin()
	if err != nil {
		klog.Fatalf("Failed to get test file: %v", err)
	}

	args := append(strings.Split(*ginkgoFlags, " "), test, "--")
	args = append(args, strings.Split(*testFlags, " ")...)

	if execCommand(ginkgo, args...); err != nil {
		klog.Exitf("Test failed: %v", err)
	}

}

func getBazelTestBin() (string, error) {
	k8sRoot, err := utils.GetK8sRootDir()
	if err != nil {
		return "", err
	}
	buildFile := filepath.Join(k8sRoot, "bazel-bin/test/e2e_kubeadm/e2e_kubeadm.test")
	if _, err := os.Stat(buildFile); err != nil {
		return "", err
	}
	return buildFile, nil

}

func getBazelGinkgo() (string, error) {
	k8sRoot, err := utils.GetK8sRootDir()
	if err != nil {
		return "", err
	}
	buildOutputDir := filepath.Join(k8sRoot, "bazel-bin", "vendor/github.com/onsi/ginkgo/ginkgo", fmt.Sprintf("%s_%s_stripped", runtime.GOOS, runtime.GOARCH), "ginkgo")
	if _, err := os.Stat(buildOutputDir); err != nil {
		return "", err
	}
	return buildOutputDir, nil
}

func execCommand(binary string, args ...string) error {
	fmt.Printf("Running command: %v %v\n", binary, strings.Join(args, " "))
	cmd := exec.Command("sh", "-c", strings.Join(append([]string{binary}, args...), " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
