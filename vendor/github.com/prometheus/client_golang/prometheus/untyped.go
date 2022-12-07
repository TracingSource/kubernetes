package prometheus

// UntypedOpts is an alias for Opts. See there for doc comments.
type UntypedOpts Opts

// UntypedFunc works like GaugeFunc but the collected metric is of type
// "Untyped". UntypedFunc is useful to mirror an external metric of unknown
// type.
//
// To create UntypedFunc instances, use NewUntypedFunc.
type UntypedFunc interface {
	Metric
	Collector
}

// NewUntypedFunc creates a new UntypedFunc based on the provided
// UntypedOpts. The value reported is determined by calling the given function
// from within the Write method. Take into account that metric collection may
// happen concurrently. If that results in concurrent calls to Write, like in
// the case where an UntypedFunc is directly registered with Prometheus, the
// provided function must be concurrency-safe.
func NewUntypedFunc(opts UntypedOpts, function func() float64) UntypedFunc {
	return newValueFunc(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), UntypedValue, function)
}
