// Copyright 2013-2015 CoreOS, Inc.
//


package semver

import (
	"sort"
)

type Versions []*Version

func (s Versions) Len() int {
	return len(s)
}

func (s Versions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Versions) Less(i, j int) bool {
	return s[i].LessThan(*s[j])
}

// Sort sorts the given slice of Version
func Sort(versions []*Version) {
	sort.Sort(Versions(versions))
}
