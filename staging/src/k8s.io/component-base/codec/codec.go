package codec

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// NewLenientSchemeAndCodecs constructs a CodecFactory with strict decoding
// disabled, that has only the Schemes registered into it which are passed
// and added via AddToScheme functions. This can be used to skip strict decoding
// a specific version only.
func NewLenientSchemeAndCodecs(addToSchemeFns ...func(s *runtime.Scheme) error) (*runtime.Scheme, *serializer.CodecFactory, error) {
	lenientScheme := runtime.NewScheme()
	for _, s := range addToSchemeFns {
		if err := s(lenientScheme); err != nil {
			return nil, nil, fmt.Errorf("unable to add API to lenient scheme: %v", err)
		}
	}
	lenientCodecs := serializer.NewCodecFactory(lenientScheme, serializer.DisableStrict)
	return lenientScheme, &lenientCodecs, nil
}
