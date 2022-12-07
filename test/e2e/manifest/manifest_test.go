package manifest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	networkingv1beta1 "k8s.io/api/networking/v1beta1"
)

func TestIngressToManifest(t *testing.T) {
	ing := &networkingv1beta1.Ingress{}
	// Create a temp dir.
	tmpDir, err := ioutil.TempDir("", "kubemci")
	if err != nil {
		t.Fatalf("unexpected error in creating temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)
	ingPath := filepath.Join(tmpDir, "ing.yaml")

	// Write the ingress to a file and ensure that there is no error.
	if err := IngressToManifest(ing, ingPath); err != nil {
		t.Fatalf("Error in creating file: %s", err)
	}
	// Writing it again should not return an error.
	if err := IngressToManifest(ing, ingPath); err != nil {
		t.Fatalf("Error in creating file: %s", err)
	}
}
