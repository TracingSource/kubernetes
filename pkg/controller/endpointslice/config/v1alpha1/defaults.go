package v1alpha1

import (
	kubectrlmgrconfigv1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
)

// RecommendedDefaultEndpointSliceControllerConfiguration defaults a pointer to
// a EndpointSliceControllerConfiguration struct. This will set the recommended
// default values, but they may be subject to change between API versions. This
// function is intentionally not registered in the scheme as a "normal"
// `SetDefaults_Foo` function to allow consumers of this type to set whatever
// defaults for their embedded configs. Forcing consumers to use these defaults
// would be problematic as defaulting in the scheme is done as part of the
// conversion, and there would be no easy way to opt-out. Instead, if you want
// to use this defaulting method run it in your wrapper struct of this type in
// its `SetDefaults_` method.
func RecommendedDefaultEndpointSliceControllerConfiguration(obj *kubectrlmgrconfigv1alpha1.EndpointSliceControllerConfiguration) {
	if obj.ConcurrentServiceEndpointSyncs == 0 {
		obj.ConcurrentServiceEndpointSyncs = 5
	}

	if obj.MaxEndpointsPerSlice == 0 {
		obj.MaxEndpointsPerSlice = 100
	}
}
