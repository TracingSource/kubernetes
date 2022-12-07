// Package k8sdeps provides kustomize factory with k8s dependencies
package k8sdeps

import (
	"k8s.io/cli-runtime/pkg/kustomize/k8sdeps/kunstruct"
	"k8s.io/cli-runtime/pkg/kustomize/k8sdeps/transformer"
	"k8s.io/cli-runtime/pkg/kustomize/k8sdeps/validator"
	"sigs.k8s.io/kustomize/pkg/factory"
)

// NewFactory creates an instance of KustFactory using k8sdeps factories
func NewFactory() *factory.KustFactory {
	return factory.NewKustFactory(
		kunstruct.NewKunstructuredFactoryImpl(),
		validator.NewKustValidator(),
		transformer.NewFactoryImpl(),
	)
}
