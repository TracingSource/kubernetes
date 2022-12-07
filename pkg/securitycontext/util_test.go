package securitycontext

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
)

func TestAddNoNewPrivileges(t *testing.T) {
	pfalse := false
	ptrue := true

	tests := map[string]struct {
		sc     *v1.SecurityContext
		expect bool
	}{
		"allowPrivilegeEscalation nil security context nil": {
			sc:     nil,
			expect: false,
		},
		"allowPrivilegeEscalation nil": {
			sc: &v1.SecurityContext{
				AllowPrivilegeEscalation: nil,
			},
			expect: false,
		},
		"allowPrivilegeEscalation false": {
			sc: &v1.SecurityContext{
				AllowPrivilegeEscalation: &pfalse,
			},
			expect: true,
		},
		"allowPrivilegeEscalation true": {
			sc: &v1.SecurityContext{
				AllowPrivilegeEscalation: &ptrue,
			},
			expect: false,
		},
	}

	for k, v := range tests {
		actual := AddNoNewPrivileges(v.sc)
		if actual != v.expect {
			t.Errorf("%s failed, expected %t but received %t", k, v.expect, actual)
		}
	}
}

func TestConvertToRuntimeMaskedPaths(t *testing.T) {
	dPM := v1.DefaultProcMount
	uPM := v1.UnmaskedProcMount
	tests := map[string]struct {
		pm     *v1.ProcMountType
		expect []string
	}{
		"procMount nil": {
			pm:     nil,
			expect: defaultMaskedPaths,
		},
		"procMount default": {
			pm:     &dPM,
			expect: defaultMaskedPaths,
		},
		"procMount unmasked": {
			pm:     &uPM,
			expect: []string{},
		},
	}

	for k, v := range tests {
		actual := ConvertToRuntimeMaskedPaths(v.pm)
		if !reflect.DeepEqual(actual, v.expect) {
			t.Errorf("%s failed, expected %#v but received %#v", k, v.expect, actual)
		}
	}
}

func TestConvertToRuntimeReadonlyPaths(t *testing.T) {
	dPM := v1.DefaultProcMount
	uPM := v1.UnmaskedProcMount
	tests := map[string]struct {
		pm     *v1.ProcMountType
		expect []string
	}{
		"procMount nil": {
			pm:     nil,
			expect: defaultReadonlyPaths,
		},
		"procMount default": {
			pm:     &dPM,
			expect: defaultReadonlyPaths,
		},
		"procMount unmasked": {
			pm:     &uPM,
			expect: []string{},
		},
	}

	for k, v := range tests {
		actual := ConvertToRuntimeReadonlyPaths(v.pm)
		if !reflect.DeepEqual(actual, v.expect) {
			t.Errorf("%s failed, expected %#v but received %#v", k, v.expect, actual)
		}
	}
}
