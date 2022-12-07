package config

// DaemonSetControllerConfiguration contains elements describing DaemonSetController.
type DaemonSetControllerConfiguration struct {
	// concurrentDaemonSetSyncs is the number of daemonset objects that are
	// allowed to sync concurrently. Larger number = more responsive daemonset,
	// but more CPU (and network) load.
	ConcurrentDaemonSetSyncs int32
}
