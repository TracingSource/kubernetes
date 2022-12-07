package install

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	quota "k8s.io/kubernetes/pkg/quota/v1"
	core "k8s.io/kubernetes/pkg/quota/v1/evaluator/core"
	generic "k8s.io/kubernetes/pkg/quota/v1/generic"
)

// NewQuotaConfigurationForAdmission returns a quota configuration for admission control.
func NewQuotaConfigurationForAdmission() quota.Configuration {
	evaluators := core.NewEvaluators(nil)
	return generic.NewConfiguration(evaluators, DefaultIgnoredResources())
}

// NewQuotaConfigurationForControllers returns a quota configuration for controllers.
func NewQuotaConfigurationForControllers(f quota.ListerForResourceFunc) quota.Configuration {
	evaluators := core.NewEvaluators(f)
	return generic.NewConfiguration(evaluators, DefaultIgnoredResources())
}

// ignoredResources are ignored by quota by default
var ignoredResources = map[schema.GroupResource]struct{}{
	{Group: "", Resource: "events"}: {},
}

// DefaultIgnoredResources returns the default set of resources that quota system
// should ignore. This is exposed so downstream integrators can have access to them.
func DefaultIgnoredResources() map[schema.GroupResource]struct{} {
	return ignoredResources
}
