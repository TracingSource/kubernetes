package v1

import "github.com/opencontainers/image-spec/specs-go"

// Index references manifests for various platforms.
// This structure provides `application/vnd.oci.image.index.v1+json` mediatype when marshalled to JSON.
type Index struct {
	specs.Versioned

	// Manifests references platform specific manifests.
	Manifests []Descriptor `json:"manifests"`

	// Annotations contains arbitrary metadata for the image index.
	Annotations map[string]string `json:"annotations,omitempty"`
}
