package image

// DeprecatedImage contains an image and a new tag,
// which will replace the original tag.
// Deprecated, instead use Image.
type DeprecatedImage struct {
	// Name is a tag-less image name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// NewTag is the value to use in replacing the original tag.
	NewTag string `json:"newTag,omitempty" yaml:"newTag,omitempty"`

	// Digest is the value used to replace the original image tag.
	// If digest is present NewTag value is ignored.
	Digest string `json:"digest,omitempty" yaml:"digest,omitempty"`
}
