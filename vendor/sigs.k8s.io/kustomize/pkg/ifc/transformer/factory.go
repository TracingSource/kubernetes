// Package patch holds miscellaneous interfaces used by kustomize.
package transformer

import (
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/transformers"
)

// Factory makes transformers
type Factory interface {
	MakePatchTransformer(slice []*resource.Resource, rf *resource.Factory) (transformers.Transformer, error)
	MakeHashTransformer() transformers.Transformer
}
