package util

import (
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	authorizationapi "k8s.io/kubernetes/pkg/apis/authorization"
)

// ResourceAttributesFrom combines the API object information and the user.Info from the context to build a full authorizer.AttributesRecord for resource access
func ResourceAttributesFrom(user user.Info, in authorizationapi.ResourceAttributes) authorizer.AttributesRecord {
	return authorizer.AttributesRecord{
		User:            user,
		Verb:            in.Verb,
		Namespace:       in.Namespace,
		APIGroup:        in.Group,
		APIVersion:      in.Version,
		Resource:        in.Resource,
		Subresource:     in.Subresource,
		Name:            in.Name,
		ResourceRequest: true,
	}
}

// NonResourceAttributesFrom combines the API object information and the user.Info from the context to build a full authorizer.AttributesRecord for non resource access
func NonResourceAttributesFrom(user user.Info, in authorizationapi.NonResourceAttributes) authorizer.AttributesRecord {
	return authorizer.AttributesRecord{
		User:            user,
		ResourceRequest: false,
		Path:            in.Path,
		Verb:            in.Verb,
	}
}

func convertToUserInfoExtra(extra map[string]authorizationapi.ExtraValue) map[string][]string {
	if extra == nil {
		return nil
	}
	ret := map[string][]string{}
	for k, v := range extra {
		ret[k] = []string(v)
	}

	return ret
}

// AuthorizationAttributesFrom takes a spec and returns the proper authz attributes to check it.
func AuthorizationAttributesFrom(spec authorizationapi.SubjectAccessReviewSpec) authorizer.AttributesRecord {
	userToCheck := &user.DefaultInfo{
		Name:   spec.User,
		Groups: spec.Groups,
		UID:    spec.UID,
		Extra:  convertToUserInfoExtra(spec.Extra),
	}

	var authorizationAttributes authorizer.AttributesRecord
	if spec.ResourceAttributes != nil {
		authorizationAttributes = ResourceAttributesFrom(userToCheck, *spec.ResourceAttributes)
	} else {
		authorizationAttributes = NonResourceAttributesFrom(userToCheck, *spec.NonResourceAttributes)
	}

	return authorizationAttributes
}
