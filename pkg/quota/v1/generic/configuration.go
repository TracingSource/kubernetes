package generic

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	quota "k8s.io/kubernetes/pkg/quota/v1"
)

// implements a basic configuration
type simpleConfiguration struct {
	evaluators       []quota.Evaluator
	ignoredResources map[schema.GroupResource]struct{}
}

// NewConfiguration creates a quota configuration
func NewConfiguration(evaluators []quota.Evaluator, ignoredResources map[schema.GroupResource]struct{}) quota.Configuration {
	return &simpleConfiguration{
		evaluators:       evaluators,
		ignoredResources: ignoredResources,
	}
}

func (c *simpleConfiguration) IgnoredResources() map[schema.GroupResource]struct{} {
	return c.ignoredResources
}

func (c *simpleConfiguration) Evaluators() []quota.Evaluator {
	return c.evaluators
}
