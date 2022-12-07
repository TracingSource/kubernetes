package configmap

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	api "k8s.io/kubernetes/pkg/apis/core"
)

func TestConfigMapStrategy(t *testing.T) {
	ctx := genericapirequest.NewDefaultContext()
	if !Strategy.NamespaceScoped() {
		t.Errorf("ConfigMap must be namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("ConfigMap should not allow create on update")
	}

	cfg := &api.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "valid-config-data",
			Namespace: metav1.NamespaceDefault,
		},
		Data: map[string]string{
			"foo": "bar",
		},
	}

	Strategy.PrepareForCreate(ctx, cfg)

	errs := Strategy.Validate(ctx, cfg)
	if len(errs) != 0 {
		t.Errorf("unexpected error validating %v", errs)
	}

	newCfg := &api.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "valid-config-data-2",
			Namespace:       metav1.NamespaceDefault,
			ResourceVersion: "4",
		},
		Data: map[string]string{
			"invalidKey": "updatedValue",
		},
	}

	Strategy.PrepareForUpdate(ctx, newCfg, cfg)

	errs = Strategy.ValidateUpdate(ctx, newCfg, cfg)
	if len(errs) == 0 {
		t.Errorf("Expected a validation error")
	}
}
