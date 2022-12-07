package metricdata

// LabelKey represents key of a label. It has optional
// description attribute.
type LabelKey struct {
	Key         string
	Description string
}

// LabelValue represents the value of a label.
// The zero value represents a missing label value, which may be treated
// differently to an empty string value by some back ends.
type LabelValue struct {
	Value   string // string value of the label
	Present bool   // flag that indicated whether a value is present or not
}

// NewLabelValue creates a new non-nil LabelValue that represents the given string.
func NewLabelValue(val string) LabelValue {
	return LabelValue{Value: val, Present: true}
}
