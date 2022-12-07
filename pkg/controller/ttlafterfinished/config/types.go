package config

// TTLAfterFinishedControllerConfiguration contains elements describing TTLAfterFinishedController.
type TTLAfterFinishedControllerConfiguration struct {
	// concurrentTTLSyncs is the number of TTL-after-finished collector workers that are
	// allowed to sync concurrently.
	ConcurrentTTLSyncs int32
}
