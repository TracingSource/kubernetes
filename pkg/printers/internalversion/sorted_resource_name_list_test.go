package internalversion

import (
	api "k8s.io/kubernetes/pkg/apis/core"
	"reflect"
	"sort"
	"testing"
)

func TestSortableResourceNamesSorting(t *testing.T) {
	want := SortableResourceNames{
		api.ResourceName(""),
		api.ResourceName("42"),
		api.ResourceName("bar"),
		api.ResourceName("foo"),
		api.ResourceName("foo"),
		api.ResourceName("foobar"),
	}

	in := SortableResourceNames{
		api.ResourceName("foo"),
		api.ResourceName("42"),
		api.ResourceName("foobar"),
		api.ResourceName("foo"),
		api.ResourceName("bar"),
		api.ResourceName(""),
	}

	sort.Sort(in)
	if !reflect.DeepEqual(in, want) {
		t.Errorf("got %v, want %v", in, want)
	}
}
