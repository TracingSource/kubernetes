package config

// ReplicationControllerConfiguration contains elements describing ReplicationController.
type ReplicationControllerConfiguration struct {
	// concurrentRCSyncs is the number of replication controllers that are
	// allowed to sync concurrently. Larger number = more responsive replica
	// management, but more CPU (and network) load.
	ConcurrentRCSyncs int32
}
