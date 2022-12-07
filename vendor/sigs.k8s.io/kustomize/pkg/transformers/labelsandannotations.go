package transformers

import (
	"errors"
	"fmt"

	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/transformers/config"
)

// mapTransformer applies a string->string map to fieldSpecs.
type mapTransformer struct {
	m          map[string]string
	fieldSpecs []config.FieldSpec
}

var _ Transformer = &mapTransformer{}

// NewLabelsMapTransformer constructs a mapTransformer.
func NewLabelsMapTransformer(
	m map[string]string, fs []config.FieldSpec) (Transformer, error) {
	return NewMapTransformer(fs, m)
}

// NewAnnotationsMapTransformer construct a mapTransformer.
func NewAnnotationsMapTransformer(
	m map[string]string, fs []config.FieldSpec) (Transformer, error) {
	return NewMapTransformer(fs, m)
}

// NewMapTransformer construct a mapTransformer.
func NewMapTransformer(
	pc []config.FieldSpec, m map[string]string) (Transformer, error) {
	if m == nil {
		return NewNoOpTransformer(), nil
	}
	if pc == nil {
		return nil, errors.New("fieldSpecs is not expected to be nil")
	}
	return &mapTransformer{fieldSpecs: pc, m: m}, nil
}

// Transform apply each <key, value> pair in the mapTransformer to the
// fields specified in mapTransformer.
func (o *mapTransformer) Transform(m resmap.ResMap) error {
	for id := range m {
		objMap := m[id].Map()
		for _, path := range o.fieldSpecs {
			if !id.Gvk().IsSelected(&path.Gvk) {
				continue
			}
			err := mutateField(objMap, path.PathSlice(), path.CreateIfNotPresent, o.addMap)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *mapTransformer) addMap(in interface{}) (interface{}, error) {
	m, ok := in.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%#v is expected to be %T", in, m)
	}
	for k, v := range o.m {
		m[k] = v
	}
	return m, nil
}
