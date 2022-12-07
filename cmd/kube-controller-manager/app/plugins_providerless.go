// +build providerless

package app

import (
	"k8s.io/component-base/featuregate"

	"k8s.io/kubernetes/pkg/volume"
)

func appendAttachableLegacyProviderVolumes(allPlugins []volume.VolumePlugin, featureGate featuregate.FeatureGate) ([]volume.VolumePlugin, error) {
	// no-op when compiled without legacy cloud providers
	return allPlugins, nil
}

func appendExpandableLegacyProviderVolumes(allPlugins []volume.VolumePlugin, featureGate featuregate.FeatureGate) ([]volume.VolumePlugin, error) {
	// no-op when compiled without legacy cloud providers
	return allPlugins, nil
}

func appendLegacyProviderVolumes(allPlugins []volume.VolumePlugin, featureGate featuregate.FeatureGate) ([]volume.VolumePlugin, error) {
	// no-op when compiled without legacy cloud providers
	return allPlugins, nil
}
