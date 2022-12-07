package storageobjectinuseprotection

import (
	"context"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	api "k8s.io/kubernetes/pkg/apis/core"
	volumeutil "k8s.io/kubernetes/pkg/volume/util"
)

func TestAdmit(t *testing.T) {
	claim := &api.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind: "PersistentVolumeClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "claim",
			Namespace: "ns",
		},
	}

	pv := &api.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind: "PersistentVolume",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "pv",
		},
	}
	claimWithFinalizer := claim.DeepCopy()
	claimWithFinalizer.Finalizers = []string{volumeutil.PVCProtectionFinalizer}

	pvWithFinalizer := pv.DeepCopy()
	pvWithFinalizer.Finalizers = []string{volumeutil.PVProtectionFinalizer}

	tests := []struct {
		name           string
		resource       schema.GroupVersionResource
		object         runtime.Object
		expectedObject runtime.Object
		featureEnabled bool
		namespace      string
	}{
		{
			"create -> add finalizer",
			api.SchemeGroupVersion.WithResource("persistentvolumeclaims"),
			claim,
			claimWithFinalizer,
			true,
			claim.Namespace,
		},
		{
			"finalizer already exists -> no new finalizer",
			api.SchemeGroupVersion.WithResource("persistentvolumeclaims"),
			claimWithFinalizer,
			claimWithFinalizer,
			true,
			claimWithFinalizer.Namespace,
		},
		{
			"disabled feature -> no finalizer",
			api.SchemeGroupVersion.WithResource("persistentvolumeclaims"),
			claim,
			claim,
			false,
			claim.Namespace,
		},
		{
			"create -> add finalizer",
			api.SchemeGroupVersion.WithResource("persistentvolumes"),
			pv,
			pvWithFinalizer,
			true,
			pv.Namespace,
		},
		{
			"finalizer already exists -> no new finalizer",
			api.SchemeGroupVersion.WithResource("persistentvolumes"),
			pvWithFinalizer,
			pvWithFinalizer,
			true,
			pvWithFinalizer.Namespace,
		},
		{
			"disabled feature -> no finalizer",
			api.SchemeGroupVersion.WithResource("persistentvolumes"),
			pv,
			pv,
			false,
			pv.Namespace,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := newPlugin()
			ctrl.storageObjectInUseProtection = test.featureEnabled

			obj := test.object.DeepCopyObject()
			attrs := admission.NewAttributesRecord(
				obj,                  // new object
				obj.DeepCopyObject(), // old object, copy to be sure it's not modified
				schema.GroupVersionKind{},
				test.namespace,
				"foo",
				test.resource,
				"", // subresource
				admission.Create,
				&metav1.CreateOptions{},
				false, // dryRun
				nil,   // userInfo
			)

			err := ctrl.Admit(context.TODO(), attrs, nil)
			if err != nil {
				t.Errorf("Test %q: got unexpected error: %v", test.name, err)
			}
			if !reflect.DeepEqual(test.expectedObject, obj) {
				t.Errorf("Test %q: Expected object:\n%s\ngot:\n%s", test.name, spew.Sdump(test.expectedObject), spew.Sdump(obj))
			}
		})
	}
}
