package v1beta1

import (
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/kubernetes/pkg/apis/abac"
)

// allAuthenticated matches k8s.io/apiserver/pkg/authentication/user.AllAuthenticated,
// but we don't want an client library (which must include types), depending on a server library
const allAuthenticated = "system:authenticated"

func Convert_v1beta1_Policy_To_abac_Policy(in *Policy, out *abac.Policy, s conversion.Scope) error {
	if err := autoConvert_v1beta1_Policy_To_abac_Policy(in, out, s); err != nil {
		return err
	}

	// In v1beta1, * user or group maps to all authenticated subjects
	if in.Spec.User == "*" || in.Spec.Group == "*" {
		out.Spec.Group = allAuthenticated
		out.Spec.User = ""
	}

	return nil
}
