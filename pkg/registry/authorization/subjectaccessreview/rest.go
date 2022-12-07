package subjectaccessreview

import (
	"context"
	"fmt"

	kapierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/registry/rest"
	authorizationapi "k8s.io/kubernetes/pkg/apis/authorization"
	authorizationvalidation "k8s.io/kubernetes/pkg/apis/authorization/validation"
	authorizationutil "k8s.io/kubernetes/pkg/registry/authorization/util"
)

type REST struct {
	authorizer authorizer.Authorizer
}

func NewREST(authorizer authorizer.Authorizer) *REST {
	return &REST{authorizer}
}

func (r *REST) NamespaceScoped() bool {
	return false
}

func (r *REST) New() runtime.Object {
	return &authorizationapi.SubjectAccessReview{}
}

func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	subjectAccessReview, ok := obj.(*authorizationapi.SubjectAccessReview)
	if !ok {
		return nil, kapierrors.NewBadRequest(fmt.Sprintf("not a SubjectAccessReview: %#v", obj))
	}
	if errs := authorizationvalidation.ValidateSubjectAccessReview(subjectAccessReview); len(errs) > 0 {
		return nil, kapierrors.NewInvalid(authorizationapi.Kind(subjectAccessReview.Kind), "", errs)
	}

	if createValidation != nil {
		if err := createValidation(ctx, obj.DeepCopyObject()); err != nil {
			return nil, err
		}
	}

	authorizationAttributes := authorizationutil.AuthorizationAttributesFrom(subjectAccessReview.Spec)
	decision, reason, evaluationErr := r.authorizer.Authorize(ctx, authorizationAttributes)

	subjectAccessReview.Status = authorizationapi.SubjectAccessReviewStatus{
		Allowed: (decision == authorizer.DecisionAllow),
		Denied:  (decision == authorizer.DecisionDeny),
		Reason:  reason,
	}
	if evaluationErr != nil {
		subjectAccessReview.Status.EvaluationError = evaluationErr.Error()
	}

	return subjectAccessReview, nil
}
