package framework

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
)

// NewSingleContentTypeSerializer wraps a serializer in a NegotiatedSerializer that handles one content type
func NewSingleContentTypeSerializer(scheme *runtime.Scheme, info runtime.SerializerInfo) runtime.StorageSerializer {
	return &wrappedSerializer{
		scheme: scheme,
		info:   info,
	}
}

type wrappedSerializer struct {
	scheme *runtime.Scheme
	info   runtime.SerializerInfo
}

var _ runtime.StorageSerializer = &wrappedSerializer{}

func (s *wrappedSerializer) SupportedMediaTypes() []runtime.SerializerInfo {
	return []runtime.SerializerInfo{s.info}
}

func (s *wrappedSerializer) UniversalDeserializer() runtime.Decoder {
	return s.info.Serializer
}

func (s *wrappedSerializer) EncoderForVersion(encoder runtime.Encoder, gv runtime.GroupVersioner) runtime.Encoder {
	return versioning.NewCodec(encoder, nil, s.scheme, s.scheme, s.scheme, s.scheme, gv, nil, s.scheme.Name())
}

func (s *wrappedSerializer) DecoderToVersion(decoder runtime.Decoder, gv runtime.GroupVersioner) runtime.Decoder {
	return versioning.NewCodec(nil, decoder, s.scheme, s.scheme, s.scheme, s.scheme, nil, gv, s.scheme.Name())
}
