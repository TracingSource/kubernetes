package config

// EndpointSliceControllerConfiguration contains elements describing
// EndpointSliceController.
type EndpointSliceControllerConfiguration struct {
	// concurrentServiceEndpointSyncs is the number of service endpoint syncing
	// operations that will be done concurrently. Larger number = faster
	// endpoint slice updating, but more CPU (and network) load.
	ConcurrentServiceEndpointSyncs int32

	// maxEndpointsPerSlice is the maximum number of endpoints that will be
	// added to an EndpointSlice. More endpoints per slice will result in fewer
	// and larger endpoint slices, but larger resources.
	MaxEndpointsPerSlice int32
}
