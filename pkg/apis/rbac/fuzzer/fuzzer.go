package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/rbac"
)

// Funcs returns the fuzzer functions for the rbac api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(r *rbac.RoleRef, c fuzz.Continue) {
			c.FuzzNoCustom(r) // fuzz self without calling this function again

			// match defaulter
			if len(r.APIGroup) == 0 {
				r.APIGroup = rbac.GroupName
			}
		},
		func(r *rbac.Subject, c fuzz.Continue) {
			switch c.Int31n(3) {
			case 0:
				r.Kind = rbac.ServiceAccountKind
				r.APIGroup = ""
				c.FuzzNoCustom(&r.Name)
				c.FuzzNoCustom(&r.Namespace)
			case 1:
				r.Kind = rbac.UserKind
				r.APIGroup = rbac.GroupName
				c.FuzzNoCustom(&r.Name)
				// user "*" won't round trip because we convert it to the system:authenticated group. try again.
				for r.Name == "*" {
					c.FuzzNoCustom(&r.Name)
				}
			case 2:
				r.Kind = rbac.GroupKind
				r.APIGroup = rbac.GroupName
				c.FuzzNoCustom(&r.Name)
			}
		},
	}
}
