package v1

import "github.com/opencontainers/image-spec/specs-go"

// Manifest provides `application/vnd.oci.image.manifest.v1+json` mediatype structure when marshalled to JSON.
type Manifest struct {
	specs.Versioned

	// Config references a configuration object for a container, by digest.
	// The referenced configuration object is a JSON blob that the runtime uses to set up the container.
	Config Descriptor `json:"config"`

	// Layers is an indexed list of layers referenced by the manifest.
	Layers []Descriptor `json:"layers"`

	// Annotations contains arbitrary metadata for the image manifest.
	Annotations map[string]string `json:"annotations,omitempty"`
}
