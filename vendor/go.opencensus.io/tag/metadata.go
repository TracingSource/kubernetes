// Copyright 2019, OpenCensus Authors
//

//

package tag

const (
	// valueTTLNoPropagation prevents tag from propagating.
	valueTTLNoPropagation = 0

	// valueTTLUnlimitedPropagation allows tag to propagate without any limits on number of hops.
	valueTTLUnlimitedPropagation = -1
)

// TTL is metadata that specifies number of hops a tag can propagate.
// Details about TTL metadata is specified at https://github.com/census-instrumentation/opencensus-specs/blob/master/tags/TagMap.md#tagmetadata
type TTL struct {
	ttl int
}

var (
	// TTLUnlimitedPropagation is TTL metadata that allows tag to propagate without any limits on number of hops.
	TTLUnlimitedPropagation = TTL{ttl: valueTTLUnlimitedPropagation}

	// TTLNoPropagation is TTL metadata that prevents tag from propagating.
	TTLNoPropagation = TTL{ttl: valueTTLNoPropagation}
)

type metadatas struct {
	ttl TTL
}

// Metadata applies metadatas specified by the function.
type Metadata func(*metadatas)

// WithTTL applies metadata with provided ttl.
func WithTTL(ttl TTL) Metadata {
	return func(m *metadatas) {
		m.ttl = ttl
	}
}
