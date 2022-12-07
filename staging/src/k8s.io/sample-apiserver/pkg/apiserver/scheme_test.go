package apiserver

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	wardlefuzzer "k8s.io/sample-apiserver/pkg/apis/wardle/fuzzer"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForScheme(t, Scheme, wardlefuzzer.Funcs)
}
