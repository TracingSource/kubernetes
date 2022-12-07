package transformers

import "sigs.k8s.io/kustomize/pkg/resmap"

// noOpTransformer contains a no-op transformer.
type noOpTransformer struct{}

var _ Transformer = &noOpTransformer{}

// NewNoOpTransformer constructs a noOpTransformer.
func NewNoOpTransformer() Transformer {
	return &noOpTransformer{}
}

// Transform does nothing.
func (o *noOpTransformer) Transform(_ resmap.ResMap) error {
	return nil
}
