package scheme

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	"k8s.io/kubernetes/pkg/kubelet/apis/config/fuzzer"
)

func TestRoundTripTypes(t *testing.T) {
	scheme, _, err := NewSchemeAndCodecs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	roundtrip.RoundTripTestForScheme(t, scheme, fuzzer.Funcs)
}
