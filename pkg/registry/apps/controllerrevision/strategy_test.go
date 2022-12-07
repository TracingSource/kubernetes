package controllerrevision

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/kubernetes/pkg/apis/apps"
	api "k8s.io/kubernetes/pkg/apis/core"
)

func TestStrategy_NamespaceScoped(t *testing.T) {
	if !Strategy.NamespaceScoped() {
		t.Error("ControllerRevision strategy must be namespace scoped")
	}
}

func TestStrategy_AllowCreateOnUpdate(t *testing.T) {
	if Strategy.AllowCreateOnUpdate() {
		t.Error("ControllerRevision should not be created on update")
	}
}

func TestStrategy_Validate(t *testing.T) {
	ctx := genericapirequest.NewDefaultContext()
	var (
		valid       = newControllerRevision("validname", "validns", newObject(), 0)
		badRevision = newControllerRevision("validname", "validns", newObject(), -1)
		emptyName   = newControllerRevision("", "validns", newObject(), 0)
		invalidName = newControllerRevision("NoUppercaseOrSpecialCharsLike=Equals", "validns", newObject(), 0)
		emptyNs     = newControllerRevision("validname", "", newObject(), 100)
		invalidNs   = newControllerRevision("validname", "NoUppercaseOrSpecialCharsLike=Equals", newObject(), 100)
		nilData     = newControllerRevision("validname", "validns", nil, 0)
	)

	tests := map[string]struct {
		history *apps.ControllerRevision
		isValid bool
	}{
		"valid":             {valid, true},
		"negative revision": {badRevision, false},
		"empty name":        {emptyName, false},
		"invalid name":      {invalidName, false},
		"empty namespace":   {emptyNs, false},
		"invalid namespace": {invalidNs, false},
		"nil data":          {nilData, false},
	}

	for name, tc := range tests {
		errs := Strategy.Validate(ctx, tc.history)
		if tc.isValid && len(errs) > 0 {
			t.Errorf("%v: unexpected error: %v", name, errs)
		}
		if !tc.isValid && len(errs) == 0 {
			t.Errorf("%v: unexpected non-error", name)
		}
	}
}

func TestStrategy_ValidateUpdate(t *testing.T) {
	ctx := genericapirequest.NewDefaultContext()
	var (
		valid       = newControllerRevision("validname", "validns", newObject(), 0)
		changedData = newControllerRevision("validname", "validns",
			func() runtime.Object {
				modified := newObject()
				ss := modified.(*apps.StatefulSet)
				ss.Name = "cde"
				return modified
			}(), 0)
		changedRevision = newControllerRevision("validname", "validns", newObject(), 1)
	)

	cases := []struct {
		name       string
		newHistory *apps.ControllerRevision
		oldHistory *apps.ControllerRevision
		isValid    bool
	}{
		{
			name:       "valid",
			newHistory: valid,
			oldHistory: valid,
			isValid:    true,
		},
		{
			name:       "changed data",
			newHistory: changedData,
			oldHistory: valid,
			isValid:    false,
		},
		{
			name:       "changed revision",
			newHistory: changedRevision,
			oldHistory: valid,
			isValid:    true,
		},
	}

	for _, tc := range cases {
		errs := Strategy.ValidateUpdate(ctx, tc.newHistory, tc.oldHistory)
		if tc.isValid && len(errs) > 0 {
			t.Errorf("%v: unexpected error: %v", tc.name, errs)
		}
		if !tc.isValid && len(errs) == 0 {
			t.Errorf("%v: unexpected non-error", tc.name)
		}
	}
}

func newControllerRevision(name, namespace string, data runtime.Object, revision int64) *apps.ControllerRevision {
	return &apps.ControllerRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			ResourceVersion: "1",
			Labels:          map[string]string{"foo": "bar"},
		},
		Data:     data,
		Revision: revision,
	}
}

func newObject() runtime.Object {
	return &apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
		Spec: apps.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
			Template: api.PodTemplateSpec{
				Spec: api.PodSpec{
					RestartPolicy: api.RestartPolicyAlways,
					DNSPolicy:     api.DNSClusterFirst,
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"foo": "bar"},
				},
			},
		},
	}
}
