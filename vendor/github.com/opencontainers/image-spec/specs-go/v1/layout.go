package v1

const (
	// ImageLayoutFile is the file name of oci image layout file
	ImageLayoutFile = "oci-layout"
	// ImageLayoutVersion is the version of ImageLayout
	ImageLayoutVersion = "1.0.0"
)

// ImageLayout is the structure in the "oci-layout" file, found in the root
// of an OCI Image-layout directory.
type ImageLayout struct {
	Version string `json:"imageLayoutVersion"`
}
