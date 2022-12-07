package cronjob

import (
	"context"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/api/pod"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/batch/validation"
	corevalidation "k8s.io/kubernetes/pkg/apis/core/validation"
)

// cronJobStrategy implements verification logic for Replication Controllers.
type cronJobStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating CronJob objects.
var Strategy = cronJobStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// DefaultGarbageCollectionPolicy returns OrphanDependents for batch/v1beta1 and batch/v2alpha1 for backwards compatibility,
// and DeleteDependents for all other versions.
func (cronJobStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	var groupVersion schema.GroupVersion
	if requestInfo, found := genericapirequest.RequestInfoFrom(ctx); found {
		groupVersion = schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}
	}
	switch groupVersion {
	case batchv1beta1.SchemeGroupVersion, batchv2alpha1.SchemeGroupVersion:
		// for back compatibility
		return rest.OrphanDependents
	default:
		return rest.DeleteDependents
	}
}

// NamespaceScoped returns true because all scheduled jobs need to be within a namespace.
func (cronJobStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears the status of a scheduled job before creation.
func (cronJobStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	cronJob := obj.(*batch.CronJob)
	cronJob.Status = batch.CronJobStatus{}

	pod.DropDisabledTemplateFields(&cronJob.Spec.JobTemplate.Spec.Template, nil)
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (cronJobStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newCronJob := obj.(*batch.CronJob)
	oldCronJob := old.(*batch.CronJob)
	newCronJob.Status = oldCronJob.Status

	pod.DropDisabledTemplateFields(&newCronJob.Spec.JobTemplate.Spec.Template, &oldCronJob.Spec.JobTemplate.Spec.Template)
}

// Validate validates a new scheduled job.
func (cronJobStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	cronJob := obj.(*batch.CronJob)
	allErrs := validation.ValidateCronJob(cronJob)
	allErrs = append(allErrs, corevalidation.ValidateConditionalPodTemplate(&cronJob.Spec.JobTemplate.Spec.Template, nil, field.NewPath("spec.jobTemplate.spec.template"))...)
	return allErrs
}

// Canonicalize normalizes the object after validation.
func (cronJobStrategy) Canonicalize(obj runtime.Object) {
}

func (cronJobStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// AllowCreateOnUpdate is false for scheduled jobs; this means a POST is needed to create one.
func (cronJobStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (cronJobStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newCronJob := obj.(*batch.CronJob)
	oldCronJob := old.(*batch.CronJob)
	allErrs := validation.ValidateCronJobUpdate(newCronJob, oldCronJob)
	allErrs = append(allErrs, corevalidation.ValidateConditionalPodTemplate(
		&newCronJob.Spec.JobTemplate.Spec.Template,
		&oldCronJob.Spec.JobTemplate.Spec.Template,
		field.NewPath("spec.jobTemplate.spec.template"))...)
	return allErrs
}

type cronJobStatusStrategy struct {
	cronJobStrategy
}

// StatusStrategy is the default logic invoked when updating object status.
var StatusStrategy = cronJobStatusStrategy{Strategy}

func (cronJobStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newJob := obj.(*batch.CronJob)
	oldJob := old.(*batch.CronJob)
	newJob.Spec = oldJob.Spec
}

func (cronJobStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
