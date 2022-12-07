package v1

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ClusterRoleBinding(obj *rbacv1.ClusterRoleBinding) {
	if len(obj.RoleRef.APIGroup) == 0 {
		obj.RoleRef.APIGroup = GroupName
	}
}
func SetDefaults_RoleBinding(obj *rbacv1.RoleBinding) {
	if len(obj.RoleRef.APIGroup) == 0 {
		obj.RoleRef.APIGroup = GroupName
	}
}
func SetDefaults_Subject(obj *rbacv1.Subject) {
	if len(obj.APIGroup) == 0 {
		switch obj.Kind {
		case rbacv1.ServiceAccountKind:
			obj.APIGroup = ""
		case rbacv1.UserKind:
			obj.APIGroup = GroupName
		case rbacv1.GroupKind:
			obj.APIGroup = GroupName
		}
	}
}
