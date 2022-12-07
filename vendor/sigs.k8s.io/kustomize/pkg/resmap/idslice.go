package resmap

import (
	"sort"

	"sigs.k8s.io/kustomize/pkg/resid"
)

// IdSlice implements the sort interface.
type IdSlice []resid.ResId

var _ sort.Interface = IdSlice{}

func (a IdSlice) Len() int      { return len(a) }
func (a IdSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a IdSlice) Less(i, j int) bool {
	if !a[i].Gvk().Equals(a[j].Gvk()) {
		return a[i].Gvk().IsLessThan(a[j].Gvk())
	}
	return a[i].String() < a[j].String()
}
