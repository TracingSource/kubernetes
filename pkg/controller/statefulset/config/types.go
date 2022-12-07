package config

// StatefulSetControllerConfiguration contains elements describing StatefulSetController.
type StatefulSetControllerConfiguration struct {
	// concurrentStatefulSetSyncs is the number of statefulset objects that are
	// allowed to sync concurrently. Larger number = more responsive statefulsets,
	// but more CPU (and network) load.
	ConcurrentStatefulSetSyncs int32
}
