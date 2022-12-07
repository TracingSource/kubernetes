package install

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	wardlefuzzer "k8s.io/sample-apiserver/pkg/apis/wardle/fuzzer"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForAPIGroup(t, Install, wardlefuzzer.Funcs)
	// TODO: enable protobuf generation for the sample-apiserver
	// roundtrip.RoundTripProtobufTestForAPIGroup(t, Install, wardlefuzzer.Funcs)
}
