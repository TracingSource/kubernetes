package v1beta1

import (
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_CronJob(obj *batchv1beta1.CronJob) {
	if obj.Spec.ConcurrencyPolicy == "" {
		obj.Spec.ConcurrencyPolicy = batchv1beta1.AllowConcurrent
	}
	if obj.Spec.Suspend == nil {
		obj.Spec.Suspend = new(bool)
	}
	if obj.Spec.SuccessfulJobsHistoryLimit == nil {
		obj.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*obj.Spec.SuccessfulJobsHistoryLimit = 3
	}
	if obj.Spec.FailedJobsHistoryLimit == nil {
		obj.Spec.FailedJobsHistoryLimit = new(int32)
		*obj.Spec.FailedJobsHistoryLimit = 1
	}
}
