package trace

import (
	"encoding/binary"
)

const defaultSamplingProbability = 1e-4

// Sampler decides whether a trace should be sampled and exported.
type Sampler func(SamplingParameters) SamplingDecision

// SamplingParameters contains the values passed to a Sampler.
type SamplingParameters struct {
	ParentContext   SpanContext
	TraceID         TraceID
	SpanID          SpanID
	Name            string
	HasRemoteParent bool
}

// SamplingDecision is the value returned by a Sampler.
type SamplingDecision struct {
	Sample bool
}

// ProbabilitySampler returns a Sampler that samples a given fraction of traces.
//
// It also samples spans whose parents are sampled.
func ProbabilitySampler(fraction float64) Sampler {
	if !(fraction >= 0) {
		fraction = 0
	} else if fraction >= 1 {
		return AlwaysSample()
	}

	traceIDUpperBound := uint64(fraction * (1 << 63))
	return Sampler(func(p SamplingParameters) SamplingDecision {
		if p.ParentContext.IsSampled() {
			return SamplingDecision{Sample: true}
		}
		x := binary.BigEndian.Uint64(p.TraceID[0:8]) >> 1
		return SamplingDecision{Sample: x < traceIDUpperBound}
	})
}

// AlwaysSample returns a Sampler that samples every trace.
// Be careful about using this sampler in a production application with
// significant traffic: a new trace will be started and exported for every
// request.
func AlwaysSample() Sampler {
	return func(p SamplingParameters) SamplingDecision {
		return SamplingDecision{Sample: true}
	}
}

// NeverSample returns a Sampler that samples no traces.
func NeverSample() Sampler {
	return func(p SamplingParameters) SamplingDecision {
		return SamplingDecision{Sample: false}
	}
}
