package clusterrolebinding

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/kubernetes/pkg/apis/rbac"
	rbacv1helpers "k8s.io/kubernetes/pkg/apis/rbac/v1"
)

// Registry is an interface for things that know how to store ClusterRoleBindings.
type Registry interface {
	ListClusterRoleBindings(ctx context.Context, options *metainternalversion.ListOptions) (*rbacv1.ClusterRoleBindingList, error)
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Lister
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s rest.StandardStorage) Registry {
	return &storage{s}
}

func (s *storage) ListClusterRoleBindings(ctx context.Context, options *metainternalversion.ListOptions) (*rbacv1.ClusterRoleBindingList, error) {
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}

	ret := &rbacv1.ClusterRoleBindingList{}
	if err := rbacv1helpers.Convert_rbac_ClusterRoleBindingList_To_v1_ClusterRoleBindingList(obj.(*rbac.ClusterRoleBindingList), ret, nil); err != nil {
		return nil, err
	}
	return ret, nil
}

// AuthorizerAdapter adapts the registry to the authorizer interface
type AuthorizerAdapter struct {
	Registry Registry
}

func (a AuthorizerAdapter) ListClusterRoleBindings() ([]*rbacv1.ClusterRoleBinding, error) {
	list, err := a.Registry.ListClusterRoleBindings(genericapirequest.NewContext(), &metainternalversion.ListOptions{})
	if err != nil {
		return nil, err
	}

	ret := []*rbacv1.ClusterRoleBinding{}
	for i := range list.Items {
		ret = append(ret, &list.Items[i])
	}
	return ret, nil
}
