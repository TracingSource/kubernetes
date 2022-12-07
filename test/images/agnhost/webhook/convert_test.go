package webhook

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	fuzz "github.com/google/gofuzz"

	v1 "k8s.io/api/admission/v1"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/diff"
	admissionfuzzer "k8s.io/kubernetes/pkg/apis/admission/fuzzer"
)

func TestConvertAdmissionRequestToV1(t *testing.T) {
	f := fuzzer.FuzzerFor(admissionfuzzer.Funcs, rand.NewSource(rand.Int63()), serializer.NewCodecFactory(runtime.NewScheme()))
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Run %d/100", i), func(t *testing.T) {
			orig := &v1beta1.AdmissionRequest{}
			f.Fuzz(orig)
			converted := convertAdmissionRequestToV1(orig)
			rt := convertAdmissionRequestToV1beta1(converted)
			if !reflect.DeepEqual(orig, rt) {
				t.Errorf("expected all request fields to be in converted object but found unaccounted for differences, diff:\n%s", diff.ObjectReflectDiff(orig, converted))
			}
		})
	}
}

func TestConvertAdmissionResponseToV1beta1(t *testing.T) {
	f := fuzz.New()
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Run %d/100", i), func(t *testing.T) {
			orig := &v1.AdmissionResponse{}
			f.Fuzz(orig)
			converted := convertAdmissionResponseToV1beta1(orig)
			rt := convertAdmissionResponseToV1(converted)
			if !reflect.DeepEqual(orig, rt) {
				t.Errorf("expected all fields to be in converted object but found unaccounted for differences, diff:\n%s", diff.ObjectReflectDiff(orig, converted))
			}
		})
	}
}
