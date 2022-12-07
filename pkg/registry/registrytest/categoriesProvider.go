package registrytest

import (
	"reflect"
	"testing"

	"k8s.io/apiserver/pkg/registry/rest"
)

func AssertCategories(t *testing.T, storage rest.CategoriesProvider, expected []string) {
	actual := storage.Categories()
	ok := reflect.DeepEqual(actual, expected)
	if !ok {
		t.Errorf("categories are not equal. expected = %v actual = %v", expected, actual)
	}
}
