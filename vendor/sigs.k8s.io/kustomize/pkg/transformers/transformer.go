// Package transformers has implementations of resmap.ResMap transformers.
package transformers

import "sigs.k8s.io/kustomize/pkg/resmap"

// A Transformer modifies an instance of resmap.ResMap.
type Transformer interface {
	// Transform modifies data in the argument, e.g. adding labels to resources that can be labelled.
	Transform(m resmap.ResMap) error
}
