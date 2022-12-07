package types

import (
	"reflect"
	"sort"
	"testing"
)

func TestVersionSort(t *testing.T) {
	unsortedVersions := []string{"v4beta1", "v2beta1", "v2alpha1", "v3", "v1"}
	expected := []string{"v2alpha1", "v2beta1", "v4beta1", "v1", "v3"}
	sort.Sort(sortableSliceOfVersions(unsortedVersions))
	if !reflect.DeepEqual(unsortedVersions, expected) {
		t.Errorf("expected %#v\ngot %#v", expected, unsortedVersions)
	}
}
