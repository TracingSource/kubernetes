package framework

import (
	"fmt"
	"time"
)

// ResizeGroup resizes an instance group
func ResizeGroup(group string, size int32) error {
	if TestContext.ReportDir != "" {
		CoreDump(TestContext.ReportDir)
		defer CoreDump(TestContext.ReportDir)
	}
	return TestContext.CloudConfig.Provider.ResizeGroup(group, size)
}

// GetGroupNodes returns a node name for the specified node group
func GetGroupNodes(group string) ([]string, error) {
	return TestContext.CloudConfig.Provider.GetGroupNodes(group)
}

// GroupSize returns the size of an instance group
func GroupSize(group string) (int, error) {
	return TestContext.CloudConfig.Provider.GroupSize(group)
}

// WaitForGroupSize waits for node instance group reached the desired size
func WaitForGroupSize(group string, size int32) error {
	timeout := 30 * time.Minute
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(20 * time.Second) {
		currentSize, err := GroupSize(group)
		if err != nil {
			Logf("Failed to get node instance group size: %v", err)
			continue
		}
		if currentSize != int(size) {
			Logf("Waiting for node instance group size %d, current size %d", size, currentSize)
			continue
		}
		Logf("Node instance group has reached the desired size %d", size)
		return nil
	}
	return fmt.Errorf("timeout waiting %v for node instance group size to be %d", timeout, size)
}
