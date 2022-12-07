package container

import (
	"errors"
	"testing"
)

func TestPodSyncResult(t *testing.T) {
	okResults := []*SyncResult{
		NewSyncResult(StartContainer, "container_0"),
		NewSyncResult(SetupNetwork, "pod"),
	}
	errResults := []*SyncResult{
		NewSyncResult(KillContainer, "container_1"),
		NewSyncResult(TeardownNetwork, "pod"),
	}
	errResults[0].Fail(errors.New("error_0"), "message_0")
	errResults[1].Fail(errors.New("error_1"), "message_1")

	// If the PodSyncResult doesn't contain error result, it should not be error
	result := PodSyncResult{}
	result.AddSyncResult(okResults...)
	if result.Error() != nil {
		t.Errorf("PodSyncResult should not be error: %v", result)
	}

	// If the PodSyncResult contains error result, it should be error
	result = PodSyncResult{}
	result.AddSyncResult(okResults...)
	result.AddSyncResult(errResults...)
	if result.Error() == nil {
		t.Errorf("PodSyncResult should be error: %v", result)
	}

	// If the PodSyncResult is failed, it should be error
	result = PodSyncResult{}
	result.AddSyncResult(okResults...)
	result.Fail(errors.New("error"))
	if result.Error() == nil {
		t.Errorf("PodSyncResult should be error: %v", result)
	}

	// If the PodSyncResult is added an error PodSyncResult, it should be error
	errResult := PodSyncResult{}
	errResult.AddSyncResult(errResults...)
	result = PodSyncResult{}
	result.AddSyncResult(okResults...)
	result.AddPodSyncResult(errResult)
	if result.Error() == nil {
		t.Errorf("PodSyncResult should be error: %v", result)
	}
}
