package patch

import "sigs.k8s.io/kustomize/pkg/gvk"

// Json6902 represents a json patch for an object
// with format documented https://tools.ietf.org/html/rfc6902.
type Json6902 struct {
	// Target refers to a Kubernetes object that the json patch will be
	// applied to. It must refer to a Kubernetes resource under the
	// purview of this kustomization. Target should use the
	// raw name of the object (the name specified in its YAML,
	// before addition of a namePrefix and a nameSuffix).
	Target *Target `json:"target" yaml:"target"`

	// relative file path for a json patch file inside a kustomization
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

// Target represents the kubernetes object that the patch is applied to
type Target struct {
	gvk.Gvk   `json:",inline,omitempty" yaml:",inline,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Name      string `json:"name" yaml:"name"`
}
