package testing

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/informers"
	fakeclient "k8s.io/client-go/kubernetes/fake"
	utiltesting "k8s.io/client-go/util/testing"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/pkg/volume/csi"
	volumetest "k8s.io/kubernetes/pkg/volume/testing"
)

// NewTestPlugin creates a plugin mgr to load plugins and setup a fake client
func NewTestPlugin(t *testing.T, client *fakeclient.Clientset) (*volume.VolumePluginMgr, *volume.VolumePlugin, string) {
	tmpDir, err := utiltesting.MkTmpdir("csi-test")
	if err != nil {
		t.Fatalf("can't create temp dir: %v", err)
	}

	if client == nil {
		client = fakeclient.NewSimpleClientset()
	}

	client.Tracker().Add(&v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fakeNode",
		},
		Spec: v1.NodeSpec{},
	})

	// Start informer for CSIDrivers.
	factory := informers.NewSharedInformerFactory(client, csi.CsiResyncPeriod)
	csiDriverInformer := factory.Storage().V1beta1().CSIDrivers()
	csiDriverLister := csiDriverInformer.Lister()
	go factory.Start(wait.NeverStop)

	host := volumetest.NewFakeVolumeHostWithCSINodeName(
		tmpDir,
		client,
		csi.ProbeVolumePlugins(),
		"fakeNode",
		csiDriverLister,
	)
	plugMgr := host.GetPluginMgr()

	plug, err := plugMgr.FindPluginByName(csi.CSIPluginName)
	if err != nil {
		t.Fatalf("can't find plugin %v", csi.CSIPluginName)
	}

	if utilfeature.DefaultFeatureGate.Enabled(features.CSIDriverRegistry) {
		// Wait until the informer in CSI volume plugin has all CSIDrivers.
		wait.PollImmediate(csi.TestInformerSyncPeriod, csi.TestInformerSyncTimeout, func() (bool, error) {
			return csiDriverInformer.Informer().HasSynced(), nil
		})
	}

	return plugMgr, &plug, tmpDir
}
