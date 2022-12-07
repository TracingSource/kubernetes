package validation

import (
	"k8s.io/kube-openapi/pkg/util/proto"
)

func ValidateModel(obj interface{}, schema proto.Schema, name string) []error {
	rootValidation, err := itemFactory(proto.NewPath(name), obj)
	if err != nil {
		return []error{err}
	}
	schema.Accept(rootValidation)
	return rootValidation.Errors()
}
