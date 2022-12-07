package fuzzer

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForAPIGroup(t, scheme.AddToScheme, Funcs)
}
