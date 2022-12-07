package hash

import (
	"fmt"

	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/transformers"
)

type nameHashTransformer struct{}

var _ transformers.Transformer = &nameHashTransformer{}

// NewNameHashTransformer construct a nameHashTransformer.
func NewNameHashTransformer() transformers.Transformer {
	return &nameHashTransformer{}
}

// Transform appends hash to generated resources.
func (o *nameHashTransformer) Transform(m resmap.ResMap) error {
	for _, res := range m {
		if res.NeedHashSuffix() {
			h, err := NewKustHash().Hash(res.Map())
			if err != nil {
				return err
			}
			res.SetName(fmt.Sprintf("%s-%s", res.GetName(), h))
		}
	}
	return nil
}
