package metricdata

// Unit is a string encoded according to the case-sensitive abbreviations from the
// Unified Code for Units of Measure: http://unitsofmeasure.org/ucum.html
type Unit string

// Predefined units. To record against a unit not represented here, create your
// own Unit type constant from a string.
const (
	UnitDimensionless Unit = "1"
	UnitBytes         Unit = "By"
	UnitMilliseconds  Unit = "ms"
)
