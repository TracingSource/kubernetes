package v1alpha1

import (
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubectrlmgrconfigv1alpha1 "k8s.io/kube-controller-manager/config/v1alpha1"
)

func TestKubeControllerDefaultsRoundTrip(t *testing.T) {
	ks1 := &kubectrlmgrconfigv1alpha1.KubeControllerManagerConfiguration{}
	SetDefaults_KubeControllerManagerConfiguration(ks1)
	cm, err := convertObjToConfigMap("KubeControllerManagerConfiguration", ks1)
	if err != nil {
		t.Errorf("unexpected ConvertObjToConfigMap error %v", err)
	}

	ks2 := &kubectrlmgrconfigv1alpha1.KubeControllerManagerConfiguration{}
	if err = json.Unmarshal([]byte(cm.Data["KubeControllerManagerConfiguration"]), ks2); err != nil {
		t.Errorf("unexpected error unserializing controller manager config %v", err)
	}

	if !reflect.DeepEqual(ks2, ks1) {
		t.Errorf("Expected:\n%#v\n\nGot:\n%#v", ks1, ks2)
	}
}

// convertObjToConfigMap converts an object to a ConfigMap.
// This is specifically meant for ComponentConfigs.
func convertObjToConfigMap(name string, obj runtime.Object) (*v1.ConfigMap, error) {
	eJSONBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			name: string(eJSONBytes[:]),
		},
	}
	return cm, nil
}
