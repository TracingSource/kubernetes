package sysctl

import (
	"testing"

	api "k8s.io/kubernetes/pkg/apis/core"
)

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		whitelist     []string
		forbiddenSafe []string
		allowedUnsafe []string
		allowed       []string
		disallowed    []string
	}{
		// no container requests
		"with allow all": {
			whitelist: []string{"foo"},
			allowed:   []string{"foo"},
		},
		"empty": {
			whitelist:     []string{"foo"},
			forbiddenSafe: []string{"*"},
			disallowed:    []string{"foo"},
		},
		"without wildcard": {
			whitelist:  []string{"a", "a.b"},
			allowed:    []string{"a", "a.b"},
			disallowed: []string{"b"},
		},
		"with catch-all wildcard and non-wildcard": {
			allowedUnsafe: []string{"a.b.c", "*"},
			allowed:       []string{"a", "a.b", "a.b.c", "b"},
		},
		"without catch-all wildcard": {
			allowedUnsafe: []string{"a.*", "b.*", "c.d.e", "d.e.f.*"},
			allowed:       []string{"a.b", "b.c", "c.d.e", "d.e.f.g.h"},
			disallowed:    []string{"a", "b", "c", "c.d", "d.e", "d.e.f"},
		},
	}

	for k, v := range tests {
		strategy := NewMustMatchPatterns(v.whitelist, v.allowedUnsafe, v.forbiddenSafe)

		pod := &api.Pod{}
		errs := strategy.Validate(pod)
		if len(errs) != 0 {
			t.Errorf("%s: unexpected validaton errors for empty sysctls: %v", k, errs)
		}

		testAllowed := func() {
			sysctls := []api.Sysctl{}
			for _, s := range v.allowed {
				sysctls = append(sysctls, api.Sysctl{
					Name:  s,
					Value: "dummy",
				})
			}
			pod.Spec.SecurityContext = &api.PodSecurityContext{
				Sysctls: sysctls,
			}
			errs = strategy.Validate(pod)
			if len(errs) != 0 {
				t.Errorf("%s: unexpected validaton errors for sysctls: %v", k, errs)
			}
		}
		testDisallowed := func() {
			for _, s := range v.disallowed {
				pod.Spec.SecurityContext = &api.PodSecurityContext{
					Sysctls: []api.Sysctl{
						{
							Name:  s,
							Value: "dummy",
						},
					},
				}
				errs = strategy.Validate(pod)
				if len(errs) == 0 {
					t.Errorf("%s: expected error for sysctl %q", k, s)
				}
			}
		}

		testAllowed()
		testDisallowed()
	}
}
