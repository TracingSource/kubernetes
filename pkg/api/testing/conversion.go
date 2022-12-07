package testing

import (
	"testing"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

// TestSelectableFieldLabelConversionsOfKind verifies that given resource have field
// label conversion defined for each its selectable field.
// fields contains selectable fields of the resource.
// labelMap maps deprecated labels to their canonical names.
func TestSelectableFieldLabelConversionsOfKind(t *testing.T, apiVersion string, kind string, fields fields.Set, labelMap map[string]string) {
	badFieldLabels := []string{
		"name",
		".name",
		"bad",
		"metadata",
		"foo.bar",
	}

	value := "value"

	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		t.Errorf("kind=%s: got unexpected error: %v", kind, err)
		return
	}
	gvk := gv.WithKind(kind)

	if len(fields) == 0 {
		t.Logf("no selectable fields for kind %q, skipping", kind)
	}
	for label := range fields {
		if label == "name" {
			t.Logf("FIXME: \"name\" is deprecated by \"metadata.name\", it should be removed from selectable fields of kind=%s", kind)
			continue
		}
		newLabel, newValue, err := legacyscheme.Scheme.ConvertFieldLabel(gvk, label, value)
		if err != nil {
			t.Errorf("kind=%s label=%s: got unexpected error: %v", kind, label, err)
		} else {
			expectedLabel := label
			if l, exists := labelMap[label]; exists {
				expectedLabel = l
			}
			if newLabel != expectedLabel {
				t.Errorf("kind=%s label=%s: got unexpected label name (%q != %q)", kind, label, newLabel, expectedLabel)
			}
			if newValue != value {
				t.Errorf("kind=%s label=%s: got unexpected new value (%q != %q)", kind, label, newValue, value)
			}
		}
	}

	for _, label := range badFieldLabels {
		_, _, err := legacyscheme.Scheme.ConvertFieldLabel(gvk, label, "value")
		if err == nil {
			t.Errorf("kind=%s label=%s: got unexpected non-error", kind, label)
		}
	}
}
