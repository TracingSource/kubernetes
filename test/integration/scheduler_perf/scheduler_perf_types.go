package benchmark

// High Level Configuration for all predicates and priorities.
type schedulerPerfConfig struct {
	NodeCount    int // The number of nodes which will be seeded with metadata to match predicates and have non-trivial priority rankings.
	PodCount     int // The number of pods which will be seeded with metadata to match predicates and have non-trivial priority rankings.
	NodeAffinity *nodeAffinity
	// TODO: Other predicates and priorities to be added here.
}

// nodeAffinity priority configuration details.
type nodeAffinity struct {
	nodeAffinityKey string // Node Selection Key.
	LabelCount      int    // number of labels to be added to each node or pod.
}
