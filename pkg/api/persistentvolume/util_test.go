package persistentvolume

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/util/diff"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	featuregatetesting "k8s.io/component-base/featuregate/testing"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/features"
)

func TestDropDisabledFields(t *testing.T) {
	specWithMode := func(mode *api.PersistentVolumeMode) *api.PersistentVolumeSpec {
		return &api.PersistentVolumeSpec{VolumeMode: mode}
	}

	secretRef := &api.SecretReference{
		Name:      "expansion-secret",
		Namespace: "default",
	}

	modeBlock := api.PersistentVolumeBlock

	tests := map[string]struct {
		oldSpec             *api.PersistentVolumeSpec
		newSpec             *api.PersistentVolumeSpec
		expectOldSpec       *api.PersistentVolumeSpec
		expectNewSpec       *api.PersistentVolumeSpec
		blockEnabled        bool
		csiExpansionEnabled bool
	}{
		"disabled block clears new": {
			blockEnabled:  false,
			newSpec:       specWithMode(&modeBlock),
			expectNewSpec: specWithMode(nil),
			oldSpec:       nil,
			expectOldSpec: nil,
		},
		"disabled block clears update when old pv did not use block": {
			blockEnabled:  false,
			newSpec:       specWithMode(&modeBlock),
			expectNewSpec: specWithMode(nil),
			oldSpec:       specWithMode(nil),
			expectOldSpec: specWithMode(nil),
		},
		"disabled block does not clear new on update when old pv did use block": {
			blockEnabled:  false,
			newSpec:       specWithMode(&modeBlock),
			expectNewSpec: specWithMode(&modeBlock),
			oldSpec:       specWithMode(&modeBlock),
			expectOldSpec: specWithMode(&modeBlock),
		},

		"enabled block preserves new": {
			blockEnabled:  true,
			newSpec:       specWithMode(&modeBlock),
			expectNewSpec: specWithMode(&modeBlock),
			oldSpec:       nil,
			expectOldSpec: nil,
		},
		"enabled block preserves update when old pv did not use block": {
			blockEnabled:  true,
			newSpec:       specWithMode(&modeBlock),
			expectNewSpec: specWithMode(&modeBlock),
			oldSpec:       specWithMode(nil),
			expectOldSpec: specWithMode(nil),
		},
		"enabled block preserves update when old pv did use block": {
			blockEnabled:  true,
			newSpec:       specWithMode(&modeBlock),
			expectNewSpec: specWithMode(&modeBlock),
			oldSpec:       specWithMode(&modeBlock),
			expectOldSpec: specWithMode(&modeBlock),
		},
		"disabled csi expansion clears secrets": {
			csiExpansionEnabled: false,
			newSpec:             specWithCSISecrets(secretRef),
			expectNewSpec:       specWithCSISecrets(nil),
			oldSpec:             nil,
			expectOldSpec:       nil,
		},
		"enabled csi expansion preserve secrets": {
			csiExpansionEnabled: true,
			newSpec:             specWithCSISecrets(secretRef),
			expectNewSpec:       specWithCSISecrets(secretRef),
			oldSpec:             nil,
			expectOldSpec:       nil,
		},
		"enabled csi expansion preserve secrets when both old and new have it": {
			csiExpansionEnabled: true,
			newSpec:             specWithCSISecrets(secretRef),
			expectNewSpec:       specWithCSISecrets(secretRef),
			oldSpec:             specWithCSISecrets(secretRef),
			expectOldSpec:       specWithCSISecrets(secretRef),
		},
		"disabled csi expansion old pv had secrets": {
			csiExpansionEnabled: false,
			newSpec:             specWithCSISecrets(secretRef),
			expectNewSpec:       specWithCSISecrets(secretRef),
			oldSpec:             specWithCSISecrets(secretRef),
			expectOldSpec:       specWithCSISecrets(secretRef),
		},
		"enabled csi expansion preserves secrets when old pv did not had secrets": {
			csiExpansionEnabled: true,
			newSpec:             specWithCSISecrets(secretRef),
			expectNewSpec:       specWithCSISecrets(secretRef),
			oldSpec:             specWithCSISecrets(nil),
			expectOldSpec:       specWithCSISecrets(nil),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.BlockVolume, tc.blockEnabled)()
			defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.ExpandCSIVolumes, tc.csiExpansionEnabled)()

			DropDisabledFields(tc.newSpec, tc.oldSpec)
			if !reflect.DeepEqual(tc.newSpec, tc.expectNewSpec) {
				t.Error(diff.ObjectReflectDiff(tc.newSpec, tc.expectNewSpec))
			}
			if !reflect.DeepEqual(tc.oldSpec, tc.expectOldSpec) {
				t.Error(diff.ObjectReflectDiff(tc.oldSpec, tc.expectOldSpec))
			}
		})
	}
}

func specWithCSISecrets(secret *api.SecretReference) *api.PersistentVolumeSpec {
	pvSpec := &api.PersistentVolumeSpec{
		PersistentVolumeSource: api.PersistentVolumeSource{
			CSI: &api.CSIPersistentVolumeSource{
				Driver:       "com.google.gcepd",
				VolumeHandle: "foobar",
			},
		},
	}

	if secret != nil {
		pvSpec.CSI.ControllerExpandSecretRef = secret
	}
	return pvSpec
}
