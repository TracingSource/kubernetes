package kubelet

import (
	"testing"

	"k8s.io/apimachinery/pkg/types"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

func TestReasonCache(t *testing.T) {
	// Create test sync result
	syncResult := kubecontainer.PodSyncResult{}
	results := []*kubecontainer.SyncResult{
		// reason cache should be set for SyncResult with StartContainer action and error
		kubecontainer.NewSyncResult(kubecontainer.StartContainer, "container_1"),
		// reason cache should not be set for SyncResult with StartContainer action but without error
		kubecontainer.NewSyncResult(kubecontainer.StartContainer, "container_2"),
		// reason cache should not be set for SyncResult with other actions
		kubecontainer.NewSyncResult(kubecontainer.KillContainer, "container_3"),
	}
	results[0].Fail(kubecontainer.ErrRunContainer, "message_1")
	results[2].Fail(kubecontainer.ErrKillContainer, "message_3")
	syncResult.AddSyncResult(results...)
	uid := types.UID("pod_1")

	reasonCache := NewReasonCache()
	reasonCache.Update(uid, syncResult)
	assertReasonInfo(t, reasonCache, uid, results[0], true)
	assertReasonInfo(t, reasonCache, uid, results[1], false)
	assertReasonInfo(t, reasonCache, uid, results[2], false)

	reasonCache.Remove(uid, results[0].Target.(string))
	assertReasonInfo(t, reasonCache, uid, results[0], false)
}

func assertReasonInfo(t *testing.T, cache *ReasonCache, uid types.UID, result *kubecontainer.SyncResult, found bool) {
	name := result.Target.(string)
	actualReason, ok := cache.Get(uid, name)
	if ok && !found {
		t.Fatalf("unexpected cache hit: %v, %q", actualReason.Err, actualReason.Message)
	}
	if !ok && found {
		t.Fatalf("corresponding reason info not found")
	}
	if !found {
		return
	}
	reason := result.Error
	message := result.Message
	if actualReason.Err != reason || actualReason.Message != message {
		t.Errorf("expected %v %q, got %v %q", reason, message, actualReason.Err, actualReason.Message)
	}
}
