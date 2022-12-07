package resource

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"
)

// FakeCategoryExpander is for testing only
var FakeCategoryExpander restmapper.CategoryExpander = restmapper.SimpleCategoryExpander{
	Expansions: map[string][]schema.GroupResource{
		"all": {
			{Group: "", Resource: "pods"},
			{Group: "", Resource: "replicationcontrollers"},
			{Group: "", Resource: "services"},
			{Group: "apps", Resource: "statefulsets"},
			{Group: "autoscaling", Resource: "horizontalpodautoscalers"},
			{Group: "batch", Resource: "jobs"},
			{Group: "batch", Resource: "cronjobs"},
			{Group: "extensions", Resource: "daemonsets"},
			{Group: "extensions", Resource: "deployments"},
			{Group: "extensions", Resource: "replicasets"},
		},
	},
}
