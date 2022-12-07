package stats

// Float64Measure is a measure for float64 values.
type Float64Measure struct {
	desc *measureDescriptor
}

// M creates a new float64 measurement.
// Use Record to record measurements.
func (m *Float64Measure) M(v float64) Measurement {
	return Measurement{
		m:    m,
		desc: m.desc,
		v:    v,
	}
}

// Float64 creates a new measure for float64 values.
//
// See the documentation for interface Measure for more guidance on the
// parameters of this function.
func Float64(name, description, unit string) *Float64Measure {
	mi := registerMeasureHandle(name, description, unit)
	return &Float64Measure{mi}
}

// Name returns the name of the measure.
func (m *Float64Measure) Name() string {
	return m.desc.name
}

// Description returns the description of the measure.
func (m *Float64Measure) Description() string {
	return m.desc.description
}

// Unit returns the unit of the measure.
func (m *Float64Measure) Unit() string {
	return m.desc.unit
}
