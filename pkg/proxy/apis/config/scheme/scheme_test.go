package scheme

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	"k8s.io/kubernetes/pkg/proxy/apis/config/fuzzer"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForScheme(t, Scheme, fuzzer.Funcs)
}
