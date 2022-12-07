package v1alpha1

import (
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestCloudControllerManagerDefaultsRoundTrip(t *testing.T) {
	ks1 := &CloudControllerManagerConfiguration{}
	SetDefaults_CloudControllerManagerConfiguration(ks1)
	cm, err := convertObjToConfigMap("CloudControllerManagerConfiguration", ks1)
	if err != nil {
		t.Errorf("unexpected ConvertObjToConfigMap error %v", err)
	}

	ks2 := &CloudControllerManagerConfiguration{}
	if err = json.Unmarshal([]byte(cm.Data["CloudControllerManagerConfiguration"]), ks2); err != nil {
		t.Errorf("unexpected error unserializing cloud controller manager config %v", err)
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
