package openapi

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"

	"k8s.io/apimachinery/pkg/util/diff"
)

func TestOpenAPIRoundtrip(t *testing.T) {
	dummyRef := func(name string) spec.Ref { return spec.MustCreateRef("#/definitions/dummy") }
	for name, value := range GetOpenAPIDefinitions(dummyRef) {
		t.Run(name, func(t *testing.T) {
			data, err := json.Marshal(value.Schema)
			if err != nil {
				t.Error(err)
				return
			}

			roundTripped := spec.Schema{}
			if err := json.Unmarshal(data, &roundTripped); err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(value.Schema, roundTripped) {
				t.Errorf("unexpected diff (a=expected,b=roundtripped):\n%s", diff.ObjectReflectDiff(value.Schema, roundTripped))
				return
			}
		})
	}
}
