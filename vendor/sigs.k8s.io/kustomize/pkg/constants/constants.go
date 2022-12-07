// Package constants holds global constants for the kustomize tool.
package constants

// KustomizationFileNames is a list of filenames that can be recognized and consumbed
// by Kustomize.
// In each directory, Kustomize searches for file with the name in this list.
// Only one match is allowed.
var KustomizationFileNames = []string{
	"kustomization.yaml",
	"kustomization.yml",
	"Kustomization",
}
