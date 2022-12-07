package v1

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_Job(obj *batchv1.Job) {
	// For a non-parallel job, you can leave both `.spec.completions` and
	// `.spec.parallelism` unset.  When both are unset, both are defaulted to 1.
	if obj.Spec.Completions == nil && obj.Spec.Parallelism == nil {
		obj.Spec.Completions = new(int32)
		*obj.Spec.Completions = 1
		obj.Spec.Parallelism = new(int32)
		*obj.Spec.Parallelism = 1
	}
	if obj.Spec.Parallelism == nil {
		obj.Spec.Parallelism = new(int32)
		*obj.Spec.Parallelism = 1
	}
	if obj.Spec.BackoffLimit == nil {
		obj.Spec.BackoffLimit = new(int32)
		*obj.Spec.BackoffLimit = 6
	}
	labels := obj.Spec.Template.Labels
	if labels != nil && len(obj.Labels) == 0 {
		obj.Labels = labels
	}
}
