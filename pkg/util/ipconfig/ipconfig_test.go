package ipconfig

import (
	"testing"

	"k8s.io/utils/exec"
)

func TestGetDNSSuffixSearchList(t *testing.T) {
	// Simple test
	ipconfigInterface := New(exec.New())

	_, err := ipconfigInterface.GetDNSSuffixSearchList()
	if err != nil {
		t.Errorf("expected success, got %v", err)
	}
}
