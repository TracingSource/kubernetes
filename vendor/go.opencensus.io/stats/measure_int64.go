package stats

// Int64Measure is a measure for int64 values.
type Int64Measure struct {
	desc *measureDescriptor
}

// M creates a new int64 measurement.
// Use Record to record measurements.
func (m *Int64Measure) M(v int64) Measurement {
	return Measurement{
		m:    m,
		desc: m.desc,
		v:    float64(v),
	}
}

// Int64 creates a new measure for int64 values.
//
// See the documentation for interface Measure for more guidance on the
// parameters of this function.
func Int64(name, description, unit string) *Int64Measure {
	mi := registerMeasureHandle(name, description, unit)
	return &Int64Measure{mi}
}

// Name returns the name of the measure.
func (m *Int64Measure) Name() string {
	return m.desc.name
}

// Description returns the description of the measure.
func (m *Int64Measure) Description() string {
	return m.desc.description
}

// Unit returns the unit of the measure.
func (m *Int64Measure) Unit() string {
	return m.desc.unit
}
