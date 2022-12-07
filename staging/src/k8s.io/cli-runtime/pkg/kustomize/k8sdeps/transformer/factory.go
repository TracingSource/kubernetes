// Package transformer provides transformer factory
package transformer

import (
	"k8s.io/cli-runtime/pkg/kustomize/k8sdeps/transformer/hash"
	"k8s.io/cli-runtime/pkg/kustomize/k8sdeps/transformer/patch"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/transformers"
)

// FactoryImpl makes patch transformer and name hash transformer
type FactoryImpl struct{}

// NewFactoryImpl makes a new factoryImpl instance
func NewFactoryImpl() *FactoryImpl {
	return &FactoryImpl{}
}

// MakePatchTransformer makes a new patch transformer
func (p *FactoryImpl) MakePatchTransformer(slice []*resource.Resource, rf *resource.Factory) (transformers.Transformer, error) {
	return patch.NewPatchTransformer(slice, rf)
}

// MakeHashTransformer makes a new name hash transformer
func (p *FactoryImpl) MakeHashTransformer() transformers.Transformer {
	return hash.NewNameHashTransformer()
}
