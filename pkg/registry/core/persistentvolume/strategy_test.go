package persistentvolume

import (
	"testing"

	apitesting "k8s.io/kubernetes/pkg/api/testing"
	api "k8s.io/kubernetes/pkg/apis/core"

	// install all api groups for testing
	_ "k8s.io/kubernetes/pkg/api/testapi"
)

func TestSelectableFieldLabelConversions(t *testing.T) {
	apitesting.TestSelectableFieldLabelConversionsOfKind(t,
		"v1",
		"PersistentVolume",
		PersistentVolumeToSelectableFields(&api.PersistentVolume{}),
		map[string]string{"name": "metadata.name"},
	)
}
