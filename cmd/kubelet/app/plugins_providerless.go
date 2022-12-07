// +build providerless

package app

import (
	"k8s.io/component-base/featuregate"

	"k8s.io/kubernetes/pkg/volume"
)

func appendLegacyProviderVolumes(allPlugins []volume.VolumePlugin, featureGate featuregate.FeatureGate) ([]volume.VolumePlugin, error) {
	// no-op when we didn't compile in support for these
	return allPlugins, nil
}
