package patch

// StrategicMerge represents a relative path to a
// stategic merge patch with the format
// https://github.com/kubernetes/community/blob/master/contributors/devel/strategic-merge-patch.md
type StrategicMerge string

// Append appends a slice of patch paths to a StrategicMerge slice
func Append(patches []StrategicMerge, paths ...string) []StrategicMerge {
	for _, p := range paths {
		patches = append(patches, StrategicMerge(p))
	}
	return patches
}

// Exist determines if a patch path exists in a slice of StrategicMerge
func Exist(patches []StrategicMerge, path string) bool {
	for _, p := range patches {
		if p == StrategicMerge(path) {
			return true
		}
	}
	return false
}
