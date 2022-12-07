package resource

import (
	"encoding/json"
	"io"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// dynamicCodec is a codec that wraps the standard unstructured codec
// with special handling for Status objects.
// Deprecated only used by test code and its wrong
type dynamicCodec struct{}

func (dynamicCodec) Decode(data []byte, gvk *schema.GroupVersionKind, obj runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(data, gvk, obj)
	if err != nil {
		return nil, nil, err
	}

	if strings.ToLower(gvk.Kind) == "status" && gvk.Version == "v1" && (gvk.Group == "" || gvk.Group == "meta.k8s.io") {
		if _, ok := obj.(*metav1.Status); !ok {
			obj = &metav1.Status{}
			err := json.Unmarshal(data, obj)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return obj, gvk, nil
}

func (dynamicCodec) Encode(obj runtime.Object, w io.Writer) error {
	// There is no need to handle runtime.CacheableObject, as we only
	// fallback to other encoders here.
	return unstructured.UnstructuredJSONScheme.Encode(obj, w)
}

// Identifier implements runtime.Encoder interface.
func (dynamicCodec) Identifier() runtime.Identifier {
	return unstructured.UnstructuredJSONScheme.Identifier()
}

// UnstructuredPlusDefaultContentConfig returns a rest.ContentConfig for dynamic types.  It includes enough codecs to act as a "normal"
// serializer for the rest.client with options, status and the like.
func UnstructuredPlusDefaultContentConfig() rest.ContentConfig {
	// TODO: scheme.Codecs here should become "pkg/apis/server/scheme" which is the minimal core you need
	// to talk to a kubernetes server
	jsonInfo, _ := runtime.SerializerInfoForMediaType(scheme.Codecs.SupportedMediaTypes(), runtime.ContentTypeJSON)

	jsonInfo.Serializer = dynamicCodec{}
	jsonInfo.PrettySerializer = nil
	return rest.ContentConfig{
		AcceptContentTypes:   runtime.ContentTypeJSON,
		ContentType:          runtime.ContentTypeJSON,
		NegotiatedSerializer: serializer.NegotiatedSerializerWrapper(jsonInfo),
	}
}
